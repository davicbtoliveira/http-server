package main

import (
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sorting := r.URL.Query().Get("sort")

	if authorID != "" {
		authorUUID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Can't convert author_id to uuid", err)
			return
		}

		chirps, err := cfg.db.ListChirpsByID(r.Context(), authorUUID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Can't found chirps for the given author_id", err)
			return
		}

		respList := []Chirp{}
		for _, v := range chirps {
			respList = append(respList, Chirp{
				ID:        v.ID,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
				Body:      v.Body,
				UserID:    v.UserID,
			})
		}
		sort.Slice(respList, func(i, j int) bool { return respList[i].CreatedAt.Before(respList[j].CreatedAt) })

		if sorting == "desc" {
			sort.Slice(respList, func(i, j int) bool { return respList[i].CreatedAt.After(respList[j].CreatedAt) })
		}
		respondWithJSON(w, 200, respList)
		return
	}

	chirps, err := cfg.db.ListChirps(r.Context())
	if err != nil {
		log.Printf("Error when fetching all chirps: %s", err)
		return
	}

	respList := []Chirp{}
	for _, v := range chirps {
		respList = append(respList, Chirp{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Body:      v.Body,
			UserID:    v.UserID,
		})
	}
	sort.Slice(respList, func(i, j int) bool { return respList[i].CreatedAt.Before(respList[j].CreatedAt) })
	if sorting == "desc" {
		sort.Slice(respList, func(i, j int) bool { return respList[i].CreatedAt.After(respList[j].CreatedAt) })
	}
	respondWithJSON(w, http.StatusOK, respList)
}
