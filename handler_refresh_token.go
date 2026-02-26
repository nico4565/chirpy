package main

import (
	"net/http"
	"time"

	"github.com/nico4565/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get bearer token", err)
		return
	}

	refreshTokenEntity, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't fetch refresh token from db", err)
		return
	}
	if refreshTokenEntity.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", err)
		return
	}
	if refreshTokenEntity.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked", err)
		return
	}

	expireDuration := time.Duration(3600) * time.Second
	accesToken, err := auth.MakeJWT(refreshTokenEntity.UserID, cfg.jwt_secret, expireDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to make jwt", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accesToken,
	})
}
