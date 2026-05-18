package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/davicbtoliveira/http-server/internal/auth"
	"github.com/davicbtoliveira/http-server/internal/database"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	exp := time.Hour
	refExp := time.Now().AddDate(0, 0, 60)

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error when fetching request: %s", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Print("Incorrect email or password")
		respondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	isValidPassword, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Print("Incorrect email or password")
		respondWithError(w, 401, "Error when validating password", err)
		return
	}
	if !isValidPassword {
		respondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(
		user.ID,
		cfg.secret,
		exp,
	)
	if err != nil {
		respondWithError(w, 400, "Cannot create user token", err)
		return
	}

	rt := auth.MakeRefreshToken()
	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     rt,
		UserID:    user.ID,
		ExpiresAt: refExp,
	})
	if err != nil {
		respondWithError(w, 401, "Couldn't create refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        token,
			RefreshToken: refreshToken.Token,
		},
	})
}
