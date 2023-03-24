package main

import (
	"fmt"
	"net/http"
	"time"
)

// struct with field fileServerHits
type apiConfig struct {
	fileServerHits int
}

func main() {
	//create new server instance to handle requests
	corsMux := http.NewServeMux()

	//initialize the file server
	fileServer := http.FileServer(http.Dir("."))

	// create an instance of apiConfig to hold our stateful data
	cfg := &apiConfig{
		fileServerHits: 0,
	}

	//serve the static index.html file
	corsMux.Handle("/", cfg.middlewareFileHits(fileServer))
	//serve the assets folder containing the chirpy logo
	corsMux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	corsMux.HandleFunc("/metrics", cfg.handlerMetrics)
	corsMux.HandleFunc("/healthz", handlerHealthCheck)

	//pass in the handler func to the middleware
	handler := corsMiddleware(corsMux)

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
