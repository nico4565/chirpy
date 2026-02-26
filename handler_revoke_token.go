package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/nico4565/chirpy/internal/auth"
	"github.com/nico4565/chirpy/internal/database"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get bearer token", err)
		return
	}

	refreshTokenEntity, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshToken)
	if refreshTokenEntity.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", err)
		return
	}
	if refreshTokenEntity.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token already revoked", err)
		return
	}

	cfg.db.RevokeRefreshToken(req.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		Token: refreshToken,
	})

	respondWithJSON(w, http.StatusNoContent, response{})
}
