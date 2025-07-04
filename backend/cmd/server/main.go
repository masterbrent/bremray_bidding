package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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
	dbURL := getEnv("DATABASE_URL", "postgres://postgres@localhost/bremray_dev?sslmode=disable")
	
	// Optional configuration
	maxOpenConns := getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	maxIdleConns := getEnvAsInt("DB_MAX_IDLE_CONNS", 5)
	connMaxLifetime := getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute)

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	// Test database connection with retry
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err != nil {
			log.Printf("Failed to ping database (attempt %d/%d): %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			log.Fatalf("Failed to connect to database after %d attempts", maxRetries)
		}
		break
	}

	log.Println("Connected to database")

	// Initialize repositories
	itemRepo := repository.NewItemRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	templateRepo := repository.NewJobTemplateRepository(db)
	jobRepo := repository.NewJobRepository(db)
	companyRepo := repository.NewCompanyRepository(db)

	// Initialize services
	itemService := services.NewItemService(itemRepo)

	// Initialize handlers
	itemHandler := handlers.NewItemHandler(itemService)
	customerHandler := handlers.NewCustomerHandler(customerRepo)
	templateHandler := handlers.NewTemplateHandler(templateRepo, itemRepo)
	jobHandler := handlers.NewJobHandler(jobRepo, customerRepo, templateRepo, itemRepo)
	companyHandler := handlers.NewCompanyHandler(companyRepo)

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
	
	// Customer routes
	customerHandler.RegisterRoutes(api)
	
	// Template routes
	templateHandler.RegisterRoutes(api)
	
	// Job routes
	jobHandler.RegisterRoutes(api)
	
	// Company routes
	companyHandler.RegisterRoutes(api)
	
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

func getEnvAsInt(key string, defaultValue int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}
	if intValue, err := strconv.Atoi(strValue); err == nil {
		return intValue
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}
	if duration, err := time.ParseDuration(strValue); err == nil {
		return duration
	}
	return defaultValue
}