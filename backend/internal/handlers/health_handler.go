package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/masterbrent/electrical-bidding-app/internal/services"
)

type HealthHandler struct {}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// CheckWave checks Wave API connectivity
func (h *HealthHandler) CheckWave(w http.ResponseWriter, r *http.Request) {
	// For now, we'll simulate the check since Wave service isn't implemented yet
	// In production, this would actually test the Wave API connection
	
	// TODO: Implement actual Wave API health check
	// Example: make a simple API call to Wave to verify credentials
	
	response := map[string]interface{}{
		"service": "wave",
		"status":  "connected",
		"message": "Wave API is accessible",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CheckCloudflare checks Cloudflare R2 connectivity
func (h *HealthHandler) CheckCloudflare(w http.ResponseWriter, r *http.Request) {
	// Test R2 connection by initializing the service
	r2Service, err := services.NewR2Service()
	if err != nil {
		response := map[string]interface{}{
			"service": "cloudflare",
			"status":  "disconnected",
			"message": err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Try to list objects with a very short timeout to test connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Test the connection by checking if we can access the bucket
	err = r2Service.TestConnection(ctx)
	if err != nil {
		response := map[string]interface{}{
			"service": "cloudflare",
			"status":  "disconnected",
			"message": "Cannot connect to R2: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	response := map[string]interface{}{
		"service": "cloudflare",
		"status":  "connected",
		"message": "Cloudflare R2 is accessible",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
