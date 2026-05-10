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

	if err := httpServer.ListenAndServe(); err != nil {
		log.Printf("error occured while serving the server: %v", err)
		return
	}
}
