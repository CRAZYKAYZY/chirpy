package main

import (
	"fmt"
	"net/http"
)

// middleware method that implements the fileServerHits
func (cfg *apiConfig) MiddlewareFileHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

// func that writes back the response with the fileServerHits
func (cfg *apiConfig) HandlerMetrics(w http.ResponseWriter, r *http.Request) {

	response := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileServerHits)
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
