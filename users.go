package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Kriss-Kolak/Chirpy/internal/auth"
	internal "github.com/Kriss-Kolak/Chirpy/internal/auth"
	"github.com/Kriss-Kolak/Chirpy/internal/database"
)

func (cfg *apiConfig) AddUser(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	hashed, err := internal.HashPassword(params.Password)
	if err != nil {
		respondWithError(res, 400, "Error during hashing password")
		return
	}

	user, err := cfg.dbqueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed})
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

func (cfg *apiConfig) UpdateUserData(res http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	userToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Could not restore user token from request header: %v", err)
		respondWithError(res, 401, "Something went wrong")
		return
	}
	userID, err := auth.ValidateJWT(userToken, cfg.secretToken)
	if err != nil {
		log.Printf("Error during validation of user's token: %v", err)
		respondWithError(res, 401, "Something went wrong")
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error during decoding parameters: %v", err)
		respondWithError(res, 401, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error during hashing new user password: %v", err)
		respondWithError(res, 401, "Something went wrong")
		return
	}

	updatedUser, err := cfg.dbqueries.ChangeUserPasswordEmail(req.Context(), database.ChangeUserPasswordEmailParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		log.Printf("Error during updating user data: %v", err)
		respondWithError(res, 401, "Something went wrong")
		return
	}

	respondWithJson(res, 200, updatedUser)

}
