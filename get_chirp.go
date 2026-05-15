package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	uuidConverted, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Cannot parse the uuid: %s", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), uuidConverted)
	if err != nil {
		respondWithError(w, 404, "Chirp not found", err)
		return
	}

	resp := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	respondWithJSON(w, 200, resp)
}
