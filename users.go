package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func (cfg *apiConfig) AddUser(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type parameters struct {
		Email string `json:"email"`
	}

	type returnVals struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(res, 400, "Error during reading email")
		return
	}

	user, err := cfg.dbqueries.CreateUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(res, 400, "Error during user creation")
		return
	}

	respondWithJson(res, 201, returnVals{
		Id:        user.ID.String(),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Email:     user.Email,
	})

}

func (cfg *apiConfig) ResetUsers(res http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(res, http.StatusForbidden, "forbidden")
		return
	}
	if cfg.dbqueries == nil {
		respondWithError(res, http.StatusInternalServerError, "db not initialized")
		return
	}
	if err := cfg.dbqueries.DeleteUsers(req.Context()); err != nil {
		respondWithError(res, http.StatusInternalServerError, "failed to reset")
		return
	}
	respondWithJson(res, http.StatusOK, map[string]string{"status": "ok"})

}
