package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) readHits() int32 {
	x := cfg.fileserverHits.Load()
	return x
}

func (cfg *apiConfig) resetHits() {
	cfg.fileserverHits.Store(0)
}

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

	apiCfg := apiConfig{}
	fileServer := http.FileServer(http.Dir("./app"))
	httpMux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServer)))
	httpMux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("Hits: %v", apiCfg.readHits())))
	})
	httpMux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.resetHits()
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
	})

	if err := httpServer.ListenAndServe(); err != nil {
		log.Printf("error occured while serving the server: %v", err)
		return
	}
}
