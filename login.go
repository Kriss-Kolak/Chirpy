package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Kriss-Kolak/Chirpy/internal/auth"
	"github.com/Kriss-Kolak/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) Login(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type returnVals struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		database.User
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(res, 500, "Something went wrong ")
		return
	}

	user, err := cfg.dbqueries.GetUserWithEmail(req.Context(), params.Email)
	if err != nil {
		log.Printf("Error during getting user from database: %s", err)
		respondWithError(res, 500, "Something went wrong ")
		return
	}

	valid, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error during checking password: %s", err)
		respondWithError(res, 500, "Something went wrong ")
		return
	}

	if !valid {
		respondWithError(res, 401, "Incorrect email or password")
		return
	}

	exprirationTime := time.Duration(time.Second * 3600)

	userToken, err := auth.MakeJWT(user.ID, cfg.secretToken, exprirationTime)
	if err != nil {
		log.Printf("error during making jwt: %v", err)
		respondWithError(res, 500, "Something went wrong")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("error during making refresh token: %v", err)
		respondWithError(res, 500, "Something went wrong")
	}
	_, err = cfg.dbqueries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{Token: refreshToken, UserID: uuid.NullUUID{UUID: user.ID, Valid: true}})
	if err != nil {
		log.Printf("error during making refresh token: %v", err)
		respondWithError(res, 500, "Something went wrong")
	}

	respondWithJson(res, 200, returnVals{User: user, Token: userToken, RefreshToken: refreshToken})

}
