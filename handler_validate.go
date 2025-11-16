package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func ValidateChirp(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(res, 500, "Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(res, 400, "Chirp is too long")
		return
	}
	respondWithJson(res, 200, returnVals{Valid: true})
}
