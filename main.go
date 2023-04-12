package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CRAZYKAYZY/chirpy/internal/config"
	"github.com/CRAZYKAYZY/chirpy/internal/database"
	"github.com/CRAZYKAYZY/chirpy/internal/handlers"
	"github.com/CRAZYKAYZY/chirpy/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

// struct with field fileServerHits

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

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

	//pass in the handler func to the middleware
	handler := CorsMiddleware(r)

	//initialize the file server
	fileServer := http.FileServer(http.Dir("."))

	// create an instance of apiConfig to hold our stateful data
	cfg := &config.ApiConfig{
		FileServerHits: 0,
		JwtSecret:      jwtSecret,
	}

	cfgMiddle := &middleware.MyConfig{
		ApiConfig: config.ApiConfig{
			JwtSecret: jwtSecret,
		},
	}

	//public routes
	public := chi.NewRouter()
	r.Mount("/public", public)

	public.Post("/login", handlers.Login_handler(db, cfg))
	public.Post("/users", handlers.CreateUserHandler(db))

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
	apiRoute.Use(cfgMiddle.AuthMiddleware)
	r.Mount("/api", apiRoute)
	apiRoute.Post("/chirps", handlers.CreateChirpHandler(db))
	apiRoute.Get("/chirps", handlers.GetChirpsHandler(db))
	apiRoute.Get("/chirps/{id}", handlers.GetChirpHandler(db))

	apiRoute.Get("/users", handlers.GetAllUsersHandler(db))
	apiRoute.Get("/users/{id}", handlers.GetUsersHandler(db))
	apiRoute.Put("/users", handlers.UpdateUsersHandler(db, cfg))

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
