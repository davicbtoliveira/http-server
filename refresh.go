package main

import (
	"net/http"
	"time"

	"github.com/davicbtoliveira/http-server/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		AccessToken string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Cannot fetch the token from request header", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "cannot get user by refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.secret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, 400, "Cannot create user token", err)
		return
	}

	respondWithJSON(w, 200, payload{
		AccessToken: accessToken,
	})
}
