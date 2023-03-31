package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/CRAZYKAYZY/chirpy/internal/database"
	"github.com/CRAZYKAYZY/chirpy/internal/handlers"
	"github.com/go-chi/chi/v5"
)

// struct with field fileServerHits
type apiConfig struct {
	fileServerHits int
}

func main() {

	// Create a new Database
	db, err := database.NewDB("database.json")
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("Failed to open database file")
	}
	defer os.Remove("database.json")

	//create new server instance to handle requests
	r := chi.NewRouter()

	//initialize the file server
	fileServer := http.FileServer(http.Dir("."))

	// create an instance of apiConfig to hold our stateful data
	cfg := &apiConfig{
		fileServerHits: 0,
	}

	//create admin handler and mount to /admin namespace
	adminApi := chi.NewRouter()
	r.Mount("/admin", adminApi)
	adminApi.Get("/metrics", cfg.HandlerMetrics)

	//serve the static index.html file
	r.Mount("/", cfg.MiddlewareFileHits(fileServer))
	//serve the assets folder containing the chirpy logo
	r.Mount("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	//create new handler for the /api endpoints and mount on /api namespace
	apiRoute := chi.NewRouter()
	r.Mount("/api", apiRoute)
	apiRoute.Get("/healthz", HandlerHealthCheck)
	apiRoute.Post("/chirps", handlers.CreateChirpHandler(db))
	apiRoute.Get("/chirps", handlers.GetChirpsHandler(db))
	apiRoute.Get("/chirps/{id}", handlers.GetChirpHandler(db))
	apiRoute.Post("/users", handlers.CreateUserHandler(db))
	apiRoute.Get("/users", handlers.GetAllUsersHandler(db))
	apiRoute.Get("/users/{id}", handlers.GetUsersHandler(db))

	apiRoute.Post("/login", handlers.Login_handler(db))

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
