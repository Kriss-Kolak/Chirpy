package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Kriss-Kolak/Chirpy/internal/auth"
)

func (cfg *apiConfig) RefreshToken(res http.ResponseWriter, req *http.Request) {

	type returnVals struct {
		Token string `json:"token"`
	}

	userRefreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error during getting refresh_token from request header: %v", err)
		respondWithError(res, 500, "Something went wrong")
	}

	DBuserRefreshToken, err := cfg.dbqueries.GetRefreshTokenFromToken(req.Context(), userRefreshToken)
	if err != nil {
		log.Printf("error during getting refresh_token from database: %v", err)
		respondWithError(res, 401, "Something went wrong")
		return
	}

	if time.Since(DBuserRefreshToken.ExpiresAt) > 0 || DBuserRefreshToken.RevokedAt.Valid {
		log.Printf("user refresh token has expired")
		respondWithError(res, 401, "Something went wrong")
		return
	}

	exprirationTime := time.Duration(60 * time.Minute)

	newToken, err := auth.MakeJWT(DBuserRefreshToken.UserID.UUID, cfg.secretToken, exprirationTime)
	if err != nil {
		log.Printf("user token was not created properly")
		respondWithError(res, 500, "Something went wrong")
		return
	}
	respondWithJson(res, 200, returnVals{Token: newToken})
}

func (cfg *apiConfig) InvokeRefreshToken(res http.ResponseWriter, req *http.Request) {

	userToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("user token was not created properly")
		respondWithError(res, 500, "Something went wrong")
		return
	}
}
