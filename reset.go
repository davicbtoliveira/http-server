package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(403)
		w.Write([]byte("Cannot reset outside dev environment"))
	}

	if err := cfg.db.ResetUsers(r.Context()); err != nil {
		log.Fatalf("Error when deleting rows from `users` table: %s", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted all rows from `users` table"))
}
