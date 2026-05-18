package main

import (
	"net/http"

	"github.com/davicbtoliveira/http-server/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIdString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error to convert chirpID to uuid", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Can't fetch chirp info", err)
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Can't authenticate user", err)
		return
	}

	jwtID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, 403, "Can't validate user auth", err)
		return
	}

	if chirp.UserID != jwtID {
		respondWithError(w, 403, "You can't delete chirps that you doesn't own", err)
		return
	}

	if err := cfg.db.DeleteChirp(r.Context(), chirpID); err != nil {
		respondWithError(w, http.StatusBadRequest, "Can't Delete chirp", err)
	}

	respondWithJSON(w, 204, nil)
}
