package main

import (
	"log"
	"net/http"
)

func main() {
	httpMux := http.NewServeMux()
	httpServer := http.Server{
		Handler: httpMux,
		Addr:    ":8080",
	}

	httpMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	fileServer := http.FileServer(http.Dir("./app"))
	httpMux.Handle("/app/", http.StripPrefix("/app", fileServer))

	if err := httpServer.ListenAndServe(); err != nil {
		log.Printf("error occured while serving the server: %v", err)
		return
	}
}
