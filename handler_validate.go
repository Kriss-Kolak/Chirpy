package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var Forbiden = map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}

func ValidateChirp(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
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

	splited := strings.Fields(params.Body)
	filtered := []string{}
	for _, word := range splited {
		if Forbiden[strings.ToLower(word)] {
			filtered = append(filtered, "****")
		} else {
			filtered = append(filtered, word)
		}
	}
	result := strings.Join(filtered, " ")

	respondWithJson(res, 200, returnVals{CleanedBody: result})
}
