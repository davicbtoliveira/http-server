package main

import (
	"encoding/json"
	"net/http"

	"github.com/davicbtoliveira/http-server/internal/auth"
	"github.com/davicbtoliveira/http-server/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID string `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Can't read the API Key from the request header", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithJSON(w, 401, nil)
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Can't decode request body", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, 204, nil)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error when converting user_id to uuid", err)
		return
	}

	if err := cfg.db.UpdateChirpRed(r.Context(), database.UpdateChirpRedParams{
		ID:          userID,
		IsChirpyRed: true,
	}); err != nil {
		respondWithError(w, 404, "User not found", err)
		return
	}

	respondWithJSON(w, 204, nil)
}
