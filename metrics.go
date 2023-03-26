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
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileServerHits)))
}
