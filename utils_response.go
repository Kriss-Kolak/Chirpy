package main

import (
	"encoding/json"
	"net/http"
)

func respondWithJson(res http.ResponseWriter, code int, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	res.Write(data)
	return nil
}

func respondWithError(res http.ResponseWriter, code int, message string) error {
	return respondWithJson(res, code, map[string]string{"error": message})
}
