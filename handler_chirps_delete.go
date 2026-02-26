package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nico4565/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", nil)
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token!", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwt_secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token!", err)
		return
	}

	dBchirp, err := cfg.db.GetChirpById(req.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Chirp with ID: [%v] not found", id), nil)
		return
	}

	if dBchirp.UserID != userId {
		respondWithError(w, http.StatusForbidden, "Forbidden! User is not the author of this chirp, delete operation not allowed", nil)
		return
	}

	err = cfg.db.DeleteChirpById(req.Context(), id)

	type response struct{}

	respondWithJSON(w, http.StatusNoContent, response{})
}
