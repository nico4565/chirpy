package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nico4565/chirpy/internal/auth"
	"github.com/nico4565/chirpy/internal/database"
)

func (cfg *apiConfig) handlerPolkaUpdate(w http.ResponseWriter, req *http.Request) {
	type response struct{}

	apiKey, err := auth.GetApiKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get ApiKey from headers!", err)
		return
	}

	if apiKey != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "ApiKey doesn't match!", err)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request!", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, response{})
		return
	}

	userId, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse userId from request into uuid!", err)
		return
	}

	_, err = cfg.db.UpdateUserIsChirpyRed(req.Context(), database.UpdateUserIsChirpyRedParams{
		IsChirpyRed: true,
		ID:          userId,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("User with id:[%v] not found!", userId), err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, response{})
}
