package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	handlerHome := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handlerHome))

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting server from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())))
}
