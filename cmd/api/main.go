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
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// payload := []byte(`{"country":"IND","domain":"nic2004:52110"}`)
	// digest := auth.CreateDigest(payload)
	// fmt.Println("Digest:", digest)

	// Initialize PostgreSQL (GORM)
	database.ConnectDB()
	log.Println("Connected to database")

	// Initialize Redis
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
	mux.HandleFunc("/vlookup", auth.AuthMiddlewareVlookup(handlers.VlookupHandler))
	mux.HandleFunc("/signature", handlers.SignHandler)
	mux.HandleFunc("/sign", handlers.SignHandlers)
	// mux.HandleFunc("/search", handlers.SearchHandler)
	mux.HandleFunc("/search", auth.AuthMiddlewareSearch(handlers.SearchHandler)) 

	// Wrap mux with LoggingMiddleware
	loggedMux := auth.LoggingMiddleware(mux)

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), loggedMux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
