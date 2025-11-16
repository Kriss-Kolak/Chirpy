package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const root = "."
	const port = "8080"

	apiCfg := apiConfig{}
	NewServeMux := http.NewServeMux()

	NewServeMux.HandleFunc("GET /api/healthz", ServeReadiness)
	NewServeMux.HandleFunc("POST /api/validate_chirp", ValidateChirp)

	NewServeMux.HandleFunc("GET /admin/metrics", apiCfg.ServeMetrics)
	NewServeMux.HandleFunc("POST /admin/reset", apiCfg.ResetMetrics)

	NewServeMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(root)))))

	MyServer := http.Server{Handler: NewServeMux,
		Addr: ":" + port}
	log.Printf("Serving files from %s on port: %s\n", root, port)
	log.Fatal(MyServer.ListenAndServe())
}
