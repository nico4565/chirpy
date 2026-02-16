package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		ChirpBody string `json:"body"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.ChirpBody) > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	splitBody := strings.Split(params.ChirpBody, " ")
	cleanedSplitBody := []string{}
	for _, s := range splitBody {
		toLowS := strings.ToLower(s)
		if toLowS == "kerfuffle" || toLowS == "sharbert" || toLowS == "fornax" {
			cleanedSplitBody = append(cleanedSplitBody, "****")
			continue
		}
		cleanedSplitBody = append(cleanedSplitBody, s)
	}

	type validateResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}
	payload := validateResponse{
		CleanedBody: strings.Join(cleanedSplitBody, " "),
	}
	respondWithJSON(w, 200, payload)

}
