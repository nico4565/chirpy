package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type responseGetChirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	dBchirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		log.Fatalf("Couldn't get chirps from db!\n%v", err)
	}
	log.Print("All chirps fetched")

	chirps := []responseGetChirp{}
	for _, dBchirp := range dBchirps {
		chirps = append(chirps, responseGetChirp{
			Id:        dBchirp.ID,
			CreatedAt: dBchirp.CreatedAt,
			UpdatedAt: dBchirp.UpdatedAt,
			Body:      dBchirp.Body,
			UserId:    dBchirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", nil)
		return
	}

	dBchirp, err := cfg.db.GetChirpById(req.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't fetch chirp by ID", nil)
		return
	}
	log.Printf("Fetched chirp with id: %v", id)

	respondWithJSON(w, http.StatusOK, responseGetChirp{
		Id:        dBchirp.ID,
		CreatedAt: dBchirp.CreatedAt,
		UpdatedAt: dBchirp.UpdatedAt,
		Body:      dBchirp.Body,
		UserId:    dBchirp.UserID,
	})
}
