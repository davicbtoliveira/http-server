package main

import (
	"net/http"

	"github.com/davicbtoliveira/http-server/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 403, "Can't fetch the Bearer Token from request header", err)
		return
	}

	if err := cfg.db.RevokeRefreshToken(r.Context(), token); err != nil {
		respondWithError(w, 403, "Can't revoke the refresh token", err)
		return
	}

	respondWithJSON(w, 204, nil)
}
