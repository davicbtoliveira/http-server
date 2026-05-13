package main

import (
	"encoding/json"
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

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, err := json.Marshal(payload)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	w.WriteHeader(code)
	w.Write(resp)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnError struct {
		Error string `json:"error"`
	}
	errorStruct := returnError{
		Error: msg,
	}
	resp, err := json.Marshal(errorStruct)
	if err != nil {
		log.Print(err)
		return
	}
	w.WriteHeader(code)
	w.Write(resp)
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

	httpMux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	apiCfg := apiConfig{}
	fileServer := http.FileServer(http.Dir("./app"))
	httpMux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServer)))

	httpMux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
		`, apiCfg.readHits())))
	})

	httpMux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.resetHits()
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
	})

	httpMux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		type parameters struct {
			Body string `json:"body"`
		}

		type responseJSON struct {
			Valid bool `json:"valid"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, 400, fmt.Sprint(err))
			return
		}

		if len(params.Body) > 140 {
			respondWithError(w, 400, "Chirp is too long")
			return
		}

		resp := responseJSON{
			Valid: true,
		}
		respondWithJSON(w, 200, resp)
	})

	if err := httpServer.ListenAndServe(); err != nil {
		log.Printf("error occured while serving the server: %v", err)
		return
	}
}
