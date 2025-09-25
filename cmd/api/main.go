package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"ondc-poc/internal/auth"
	"ondc-poc/internal/database"
	"ondc-poc/internal/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	payload := []byte(`{"country":"IND","domain":"nic2004:52110"}`)
	digest := auth.CreateDigest(payload)
	fmt.Println("Digest:", digest)

	if err := database.ConnectDB(); err != nil {
		log.Fatalf("Could not connect to database: %s\n", err)
	}
	log.Println("Connected to database")

	if err := database.ConnectRedis(); err != nil {
		log.Fatalf("Could not connect to Redis: %s\n", err)
	}
	log.Println("Connected to Redis")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create mux and register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/lookup", handlers.LookupHandler)
	mux.HandleFunc("/vlookup", auth.AuthMiddleware(handlers.LookupHandler))
	mux.HandleFunc("/sign", handlers.SignHandler)

	// Wrap mux with LoggingMiddleware
	loggedMux := auth.LoggingMiddleware(mux)

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), loggedMux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
