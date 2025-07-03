package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/masterbrent/electrical-bidding-app/internal/handlers"
	"github.com/masterbrent/electrical-bidding-app/internal/middleware"
	"github.com/masterbrent/electrical-bidding-app/internal/repository"
	"github.com/masterbrent/electrical-bidding-app/internal/services"
	
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Get configuration from environment
	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgres://user:password@localhost/bremray_dev?sslmode=disable")

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database")

	// Initialize repositories
	itemRepo := repository.NewItemRepository(db)

	// Initialize services
	itemService := services.NewItemService(itemRepo)

	// Initialize handlers
	itemHandler := handlers.NewItemHandler(itemService)

	// Setup routes
	router := mux.NewRouter()
	
	// API routes
	api := router.PathPrefix("/api").Subrouter()
	
	// Apply middleware
	api.Use(middleware.CORS)
	api.Use(middleware.JSONContentType)
	api.Use(middleware.Logger)

	// Item routes
	api.HandleFunc("/items", itemHandler.Create).Methods("POST")
	api.HandleFunc("/items", itemHandler.List).Methods("GET")
	api.HandleFunc("/items/{id}", itemHandler.GetByID).Methods("GET")
	api.HandleFunc("/items/{id}", itemHandler.Update).Methods("PUT")
	api.HandleFunc("/items/{id}", itemHandler.Delete).Methods("DELETE")
	
	// Handle OPTIONS for all routes
	api.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}).Methods("OPTIONS")

	// Health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}