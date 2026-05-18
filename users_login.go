package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/davicbtoliveira/http-server/internal/auth"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	expiration := time.Hour

	type parameters struct {
		Password  string `json:"password"`
		Email     string `json:"email"`
		ExpiresIn int    `json:"expires_in_seconds"`
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

	if params.ExpiresIn > 0 {
		expiration = time.Duration(params.ExpiresIn) * time.Second
	}

	if expiration > time.Hour {
		expiration = time.Hour
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
		expiration,
	)
	if err != nil {
		respondWithError(w, 400, "Cannot create user token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		},
	})
}
