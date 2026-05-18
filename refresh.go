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

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Cannot fetch the token from request header", err)
		return
	}

	refreshToken, err := cfg.db.GetRfByToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, 401, "Failed when fetching refresh token", err)
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, 401, "revoked token", err)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, 401, "expired token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, 401, "cannot get user by refresh token", err)
		return
	}

	token, err := auth.MakeJWT(
		user.ID,
		cfg.secret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, 400, "Cannot create user token", err)
		return
	}

	respondWithJSON(w, 200, payload{
		AccessToken: token,
	})
}
