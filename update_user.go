package main

import (
	"encoding/json"
	"net/http"

	"github.com/davicbtoliveira/http-server/internal/auth"
	"github.com/davicbtoliveira/http-server/internal/database"
)

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 400, "Error when decoding request body", err)
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Can't fetch user token from header", err)
		return
	}

	jwtID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Cannot validate access token", err)
		return
	}

	hashedNewPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, "Error when hashing new password", err)
		return
	}

	updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedNewPass,
		ID:             jwtID,
	})

	respondWithJSON(w, 200, User{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	})
}
