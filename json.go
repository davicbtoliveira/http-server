package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
