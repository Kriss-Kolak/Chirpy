package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Kriss-Kolak/Chirpy/internal/database"
	"github.com/google/uuid"
)

var Forbiden = map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}

func (cfg *apiConfig) CreateChirp(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type parameters struct {
		Body   string        `json:"body"`
		UserID uuid.NullUUID `json:"user_id"`
	}

	type returnVals struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserId    string `json:"user_id"`
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
	post, err := cfg.dbqueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   result,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(res, 400, "Failue to create post")
		return
	}
	respondWithJson(res, 201, returnVals{
		Id:        post.ID.String(),
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		Body:      post.Body,
		UserId:    post.UserID.UUID.String(),
	})

}

func (cfg *apiConfig) GetAllChirps(res http.ResponseWriter, req *http.Request) {

	chirps, err := cfg.dbqueries.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(res, 400, "Failed to obtain all chirps")
		return
	}
	respondWithJson(res, 200, chirps)
}

func (cfg *apiConfig) GetChirpWithId(res http.ResponseWriter, req *http.Request) {

	ID := req.PathValue("chirpID")

	parsedID := uuid.MustParse(ID)

	chirp, err := cfg.dbqueries.GetChirpWithId(req.Context(), parsedID)
	if err != nil {
		respondWithError(res, 404, "Failed to obtain chirp")
		return
	}
	respondWithJson(res, 200, chirp)

}
