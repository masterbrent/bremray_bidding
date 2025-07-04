package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/masterbrent/electrical-bidding-app/internal/models"
	"github.com/masterbrent/electrical-bidding-app/internal/repository"
)

// CustomerHandler handles HTTP requests for customers
type CustomerHandler struct {
	customerRepo repository.CustomerRepository
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(customerRepo repository.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{
		customerRepo: customerRepo,
	}
}

// RegisterRoutes registers all customer routes
func (h *CustomerHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/customers", h.List).Methods("GET", "OPTIONS")
	router.HandleFunc("/customers", h.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/customers/{id}", h.Get).Methods("GET", "OPTIONS")
	router.HandleFunc("/customers/{id}", h.Update).Methods("PUT", "OPTIONS")
	router.HandleFunc("/customers/{id}", h.Delete).Methods("DELETE", "OPTIONS")
}

// List returns all customers
func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Parse query parameters
	limit := 50
	offset := 0
	
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}
	
	customers, err := h.customerRepo.List(ctx, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, customers)
}

// Get returns a single customer by ID
func (h *CustomerHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	
	customer, err := h.customerRepo.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "customer not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, customer)
}

// Create creates a new customer
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	customer, err := models.NewCustomer(req.Name, req.Email, req.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.customerRepo.Create(ctx, customer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	respondJSON(w, customer)
}

// Update updates an existing customer
func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	
	// Get existing customer
	customer, err := h.customerRepo.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "customer not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Update fields
	if req.Name != "" {
		customer.Name = req.Name
	}
	if req.Email != "" {
		customer.Email = req.Email
	}
	if req.Phone != "" {
		customer.Phone = req.Phone
	}
	
	if err := h.customerRepo.Update(ctx, customer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, customer)
}

// Delete deletes a customer
func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	
	if err := h.customerRepo.Delete(ctx, id); err != nil {
		if err.Error() == "customer not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}