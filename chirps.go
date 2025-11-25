package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Kriss-Kolak/Chirpy/internal/auth"
	"github.com/Kriss-Kolak/Chirpy/internal/database"
	"github.com/google/uuid"
)

var Forbiden = map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}

func (cfg *apiConfig) CreateChirp(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserId    string `json:"user_id"`
	}

	userToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting user token: %s", err)
		respondWithError(res, 401, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(userToken, cfg.secretToken)
	if err != nil {
		log.Printf("Error during token validation: %s", err)
		respondWithError(res, 401, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
		UserID: userID,
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
		UserId:    post.UserID.String(),
	})

}

func (cfg *apiConfig) GetAllChirps(res http.ResponseWriter, req *http.Request) {

	authorID := req.URL.Query().Get("author_id")
	sortMethod := req.URL.Query().Get("sort")
	var chirps []database.Chirp
	var err error
	if authorID != "" {

		parsedAuthorID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(res, 500, "Something went wrong")
			return
		}

		chirps, err = cfg.dbqueries.GetChripsFromAuthorID(req.Context(), parsedAuthorID)
		if err != nil {
			respondWithError(res, 400, "Failed to obtain all chirps")
			return
		}

	} else {
		chirps, err = cfg.dbqueries.GetAllChirps(req.Context())
		if err != nil {
			respondWithError(res, 400, "Failed to obtain all chirps")
			return
		}
	}

	if sortMethod == "" || sortMethod == "asc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Sub(chirps[j].CreatedAt) < 0 })
	} else {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Sub(chirps[j].CreatedAt) > 0 })
	}
	respondWithJson(res, 200, chirps)

}

func (cfg *apiConfig) GetChirpWithId(res http.ResponseWriter, req *http.Request) {

	ID := req.PathValue("chirpID")

	parsedID, err := uuid.Parse(ID)
	if err != nil {
		respondWithError(res, 500, "Something went wrong")
		return
	}

	chirp, err := cfg.dbqueries.GetChirpWithId(req.Context(), parsedID)
	if err != nil {
		respondWithError(res, 404, "Failed to obtain chirp")
		return
	}
	respondWithJson(res, 200, chirp)

}

func (cfg *apiConfig) DeleteChripWithId(res http.ResponseWriter, req *http.Request) {

	chirpID := req.PathValue("chirpID")
	parsedChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		log.Printf("Error during parsing chirpID: %v", err)
		respondWithError(res, 500, "Something went wrong")
		return
	}

	userToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error during getting userToken: %v", err)
		respondWithError(res, 401, "Unauthorized")
		return
	}

	userIDToken, err := auth.ValidateJWT(userToken, cfg.secretToken)
	if err != nil {
		log.Printf("Error during validating userToken: %v", err)
		respondWithError(res, 403, "Unauthorized")
		return
	}

	chrip, err := cfg.dbqueries.GetChirpWithId(req.Context(), parsedChirpID)
	if err != nil {
		log.Printf("Error during getting chirp from db: %v", err)
		respondWithError(res, 404, "Chirp not found")
		return
	}

	if userIDToken != chrip.UserID {
		log.Printf("Error during comparing userID: %v", err)
		respondWithError(res, 403, "Unauthorized")
		return
	}

	err = cfg.dbqueries.DeleteChirpWithId(req.Context(), chrip.ID)
	if err != nil {
		log.Printf("Error during deleting chirp: %v", err)
		respondWithError(res, 500, "Something went wrong")
		return
	}

	respondWithJson(res, 204, []byte{})
}
