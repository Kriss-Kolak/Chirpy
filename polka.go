package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kriss-Kolak/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) PolkaWebhook(res http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		}
	}

	requestAPIKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(res, 401, "No API key found")
		return
	}

	if requestAPIKey != cfg.polkaKey {
		respondWithError(res, 401, "API key does not match")
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error during polka parameters decoding: %v", err)
		respondWithError(res, 500, "Something went wrong")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJson(res, 204, []byte{})
		return
	}

	parsedUserID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(res, 500, "Something went wrong")
		return
	}

	err = cfg.dbqueries.UpgradeUserToRed(req.Context(), parsedUserID)
	if err != nil {
		respondWithError(res, 404, "User not found")
		return
	}

	respondWithJson(res, 204, []byte{})

}
