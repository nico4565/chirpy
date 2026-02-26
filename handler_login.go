package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nico4565/chirpy/internal/auth"
	"github.com/nico4565/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {

	type response struct {
		Id           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request!", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't fetch user! Are you sure a user with email: [%s] exists?", params.Email), err)
		return
	}

	isPasswordCorrect, hashingParams, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Couldn't check password for user with email: [%s]!", params.Email), err)
		return
	}
	if !isPasswordCorrect {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: password is wrong!", err)
		return
	}

	newHashedPassword, err := auth.OpportunisticRehashing(hashingParams, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't optimize password hashing!", err)
		return

	} else if newHashedPassword != "" {
		user, err = cfg.db.UpdateUserPassword(req.Context(), database.UpdateUserPasswordParams{
			HashedPassword: newHashedPassword,
			Email:          params.Email,
		})
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't update password with new hashing!", err)
			return
		}
	}

	expireDuration := time.Duration(3600) * time.Second
	accesToken, err := auth.MakeJWT(user.ID, cfg.secret, expireDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to make jwt", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to make refresh token", err)
		return
	}
	_, err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
		RevokedAt: sql.NullTime{},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to store refresh token into database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{

		Id:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        accesToken,
		RefreshToken: refreshToken,
	})
}
