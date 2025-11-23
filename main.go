package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Kriss-Kolak/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbqueries      *database.Queries
	platform       string
	secretToken    string
}

func main() {
	const root = "."
	const port = "8080"
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	apiCfg := apiConfig{
		dbqueries:   database.New(db),
		platform:    os.Getenv("PLATFORM"),
		secretToken: os.Getenv("SECRET_TOKEN")}

	NewServeMux := http.NewServeMux()

	NewServeMux.HandleFunc("GET /api/healthz", ServeReadiness)
	NewServeMux.HandleFunc("POST /api/chirps", apiCfg.CreateChirp)
	NewServeMux.HandleFunc("POST /api/users", apiCfg.AddUser)
	NewServeMux.HandleFunc("PUT /api/users", apiCfg.UpdateUserData)
	NewServeMux.HandleFunc("POST /api/login", apiCfg.Login)
	NewServeMux.HandleFunc("POST /api/revoke", apiCfg.InvokeRefreshToken)
	NewServeMux.HandleFunc("POST /api/refresh", apiCfg.RefreshToken)
	NewServeMux.HandleFunc("GET /api/chirps", apiCfg.GetAllChirps)
	NewServeMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirpWithId)

	NewServeMux.HandleFunc("GET /admin/metrics", apiCfg.ServeMetrics)
	NewServeMux.HandleFunc("POST /admin/reset", apiCfg.ResetUsers)

	NewServeMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(root)))))

	MyServer := http.Server{Handler: NewServeMux,
		Addr: ":" + port}
	log.Printf("Serving files from %s on port: %s\n", root, port)
	log.Fatal(MyServer.ListenAndServe())
}
