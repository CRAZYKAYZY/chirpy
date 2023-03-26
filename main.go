package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// struct with field fileServerHits
type apiConfig struct {
	fileServerHits int
}

func main() {
	//create new server instance to handle requests
	r := chi.NewRouter()

	//initialize the file server
	fileServer := http.FileServer(http.Dir("."))

	// create an instance of apiConfig to hold our stateful data
	cfg := &apiConfig{
		fileServerHits: 0,
	}

	//serve the static index.html file
	r.Mount("/", cfg.MiddlewareFileHits(fileServer))
	//serve the assets folder containing the chirpy logo
	r.Mount("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	r.Get("/metrics", cfg.HandlerMetrics)
	r.Get("/healthz", HandlerHealthCheck)

	//pass in the handler func to the middleware
	handler := CorsMiddleware(r)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//start and listen to incoming server requests
	fmt.Printf("Server listening on port %v...\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
