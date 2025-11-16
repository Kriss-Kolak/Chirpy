package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}

func (cfg *apiConfig) ServeMetrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	response_text := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())

	res.Write([]byte(response_text))
}

func (cfg *apiConfig) ResetMetrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	res.Write([]byte("OK"))
}

func ServeReadiness(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))
}

func ValidateChirp(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type parameters struct {
		Body string `json:"body"`
	}

	type responseError struct {
		Error string `json:"error"`
	}

	type responseSuccess struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		jsonResponse := responseError{Error: "Something went wrong"}
		dat, err := json.Marshal(jsonResponse)
		if err != nil {
			return
		}
		res.WriteHeader(500)
		res.Write(dat)
		return
	}
	if len(params.Body) > 140 {
		jsonResponse := responseError{Error: "Chirp is too long"}
		dat, err := json.Marshal(jsonResponse)
		if err != nil {
			return
		}
		res.WriteHeader(400)
		res.Write(dat)
		return
	}
	jsonResponse := responseSuccess{Valid: true}
	dat, err := json.Marshal(jsonResponse)
	if err != nil {
		return
	}
	res.WriteHeader(200)
	res.Write(dat)
}

func main() {
	apiCfg := apiConfig{}
	NewServeMux := http.NewServeMux()
	NewServeMux.HandleFunc("GET /api/healthz", ServeReadiness)
	NewServeMux.HandleFunc("POST /api/validate_chirp", ValidateChirp)
	NewServeMux.HandleFunc("GET /admin/metrics", apiCfg.ServeMetrics)
	NewServeMux.HandleFunc("POST /admin/reset", apiCfg.ResetMetrics)
	NewServeMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	MyServer := http.Server{Handler: NewServeMux,
		Addr: ":8080"}
	log.Printf("Serving files from %s on port: %s\n", "/", "8080")
	log.Fatal(MyServer.ListenAndServe())
}
