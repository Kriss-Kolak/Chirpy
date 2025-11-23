package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	internal "github.com/Kriss-Kolak/Chirpy/internal/auth"
	"github.com/Kriss-Kolak/Chirpy/internal/database"
)

func (cfg *apiConfig) Login(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type returnVals struct {
		Token string `json:"token"`
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

	valid, err := internal.CheckPasswordHash(params.Password, user.HashedPassword)
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

	if params.ExpiresInSeconds != 0 {
		exprirationTime = time.Duration(params.ExpiresInSeconds)
	}

	userToken, err := internal.MakeJWT(user.ID, cfg.secretToken, exprirationTime)
	if err != nil {
		respondWithError(res, 500, "Something went wrong")
		return
	}

	respondWithJson(res, 200, returnVals{User: user, Token: userToken})

}
