package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.ListChirps(r.Context())
	if err != nil {
		log.Printf("Error when fetching all chirps: %s", err)
		return
	}

	respList := []Chirp{}
	for _, v := range chirps {
		respList = append(respList, Chirp{
			ID: v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Body: v.Body,
			UserID: v.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, respList)
}
