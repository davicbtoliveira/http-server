package main

import (
	"encoding/json"
	"log"
	"net/http"
)


func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email	string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	body := params{}
	if err := decoder.Decode(&body); err != nil {
		log.Printf("Error decoding the request body: %s", err)
	}

	user, err := cfg.db.CreateUser(r.Context(), body.Email)
	if err != nil {
		log.Fatalf("Error when creating the new user: %s", err)
	}

	respondWithJSON(w, 201, user)
}
