package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/davicbtoliveira/http-server/internal/auth"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
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
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Print("Incorrect email or password")
		respondWithError(w, 401, "Incorrect email or password", err)
	}

	isValidPassword, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Print("Incorrect email or password")
		respondWithError(w, 401, "Incorrect email or password", err)
	}

	if isValidPassword {
		respondWithJSON(w, http.StatusOK, response{
			User: User{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email:     user.Email,
			},
		})
	} else {
		respondWithError(w, 401, "Incorrect email or password", err)
	}
}
