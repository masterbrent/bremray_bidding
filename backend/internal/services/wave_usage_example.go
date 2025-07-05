package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Example handler showing EXACT usage of Wave integration
func CreateWaveInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	// Parse job ID from URL params
	jobID := r.URL.Query().Get("jobId")
	if jobID == "" {
		http.Error(w, "jobId required", http.StatusBadRequest)
		return
	}
	
	// Check Wave configuration
	if !IsWaveConfigured() {
		http.Error(w, "Wave API not configured. Please add Wave credentials.", http.StatusBadRequest)
		return
	}
	
	// TODO: Fetch job from your database
	// This is a mock job structure matching the production format
	job := map[string]interface{}{
		"id":           jobID,
		"customerName": "John Smith",
		"address":      "123 Main St, Toronto, ON",
		"permitRequired": true,
		"items": []interface{}{
			// Template item example
			map[string]interface{}{
				"id":                1,
				"itemId":            101,
				"completedQuantity": float64(2),
				"isCustomItem":      false,
				"item": map[string]interface{}{
					"name":        "Outlet Installation",
					"description": "Standard 15A outlet",
					"price":       "75.00",
				},
			},
			// Custom item example
			map[string]interface{}{
				"id":                2,
				"completedQuantity": float64(1),
				"isCustomItem":      true,
				"customName":        "Troubleshoot living room circuit",
				"customDescription": "Diagnosed and repaired faulty circuit breaker",
				"customPrice":       "150.00",
			},
			// Item with zero quantity (will be skipped)
			map[string]interface{}{
				"id":                3,
				"completedQuantity": float64(0),
				"isCustomItem":      false,
				"item": map[string]interface{}{
					"name":  "Switch Installation",
					"price": "65.00",
				},
			},
		},
	}
	
	// Get Wave credentials
	credentials, err := GetWaveCredentials()
	if err != nil {
		log.Printf("Wave credentials error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Create Wave service
	waveService, err := NewWaveAPIService(*credentials)
	if err != nil {
		log.Printf("Failed to create Wave service: %v", err)
		http.Error(w, "Failed to initialize Wave service", http.StatusInternalServerError)
		return
	}
	defer waveService.Close()
	
	// Prepare job for Wave invoice
	customerID, lineItems, err := PrepareJobForWaveInvoice(job)
	if err != nil {
		log.Printf("Failed to prepare job: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Format PO number
	customerName, _ := job["customerName"].(string)
	address, _ := job["address"].(string)
	poNumber := FormatPONumber(customerName, address)
	
	log.Printf("Creating invoice with PO number: %s, Customer ID: %s", poNumber, customerID)
	
	// Create invoice in Wave
	invoice, err := waveService.CreateInvoice(customerID, lineItems, poNumber, "")
	if err != nil {
		log.Printf("Wave invoice creation failed: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create invoice in Wave: %v", err), http.StatusInternalServerError)
		return
	}
	
	log.Printf("Invoice created successfully: %+v", invoice)
	
	// TODO: Update your job in database with Wave invoice details:
	// - status: 'invoiced'
	// - waveInvoiceId: invoice.ID
	// - waveInvoiceNumber: invoice.InvoiceNumber
	// - waveInvoiceUrl: invoice.ViewURL
	
	// Return success response
	response := map[string]interface{}{
		"message": "Invoice created successfully in Wave",
		"invoice": map[string]interface{}{
			"id":     invoice.ID,
			"number": invoice.InvoiceNumber,
			"url":    invoice.ViewURL,
		},
		"jobId": jobID,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Example of testing Wave connection
func TestWaveConnectionHandler(w http.ResponseWriter, r *http.Request) {
	credentials, err := GetWaveCredentials()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	waveService, err := NewWaveAPIService(*credentials)
	if err != nil {
		http.Error(w, "Failed to initialize Wave service", http.StatusInternalServerError)
		return
	}
	defer waveService.Close()
	
	// Test by fetching products
	products, err := waveService.GetProducts()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch products: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Test by finding Skyview customer
	customer, err := waveService.FindCustomerByName("Skyview")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to search for customer: %v", err), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success":      true,
		"productCount": len(products),
		"customer":     customer,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}