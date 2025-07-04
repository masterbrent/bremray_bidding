package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/masterbrent/electrical-bidding-app/internal/repository"
)

// CompanyHandler handles HTTP requests for company settings
type CompanyHandler struct {
	companyRepo repository.CompanyRepository
}

// NewCompanyHandler creates a new company handler
func NewCompanyHandler(companyRepo repository.CompanyRepository) *CompanyHandler {
	return &CompanyHandler{
		companyRepo: companyRepo,
	}
}

// RegisterRoutes registers all company routes
func (h *CompanyHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/company", h.Get).Methods("GET", "OPTIONS")
	router.HandleFunc("/company", h.Update).Methods("PUT", "OPTIONS")
}

// Get returns the company settings
func (h *CompanyHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	company, err := h.companyRepo.Get(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, company)
}

// Update updates the company settings
func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get current company settings
	company, err := h.companyRepo.Get(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	var req struct {
		Name    string  `json:"name"`
		Logo    *string `json:"logo"`
		Address string  `json:"address"`
		City    string  `json:"city"`
		State   string  `json:"state"`
		Zip     string  `json:"zip"`
		Phone   string  `json:"phone"`
		Email   string  `json:"email"`
		License string  `json:"license"`
		Website string  `json:"website"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Update company fields
	if err := company.Update(
		req.Name,
		req.Email,
		req.Address,
		req.City,
		req.State,
		req.Zip,
		req.Phone,
		req.License,
		req.Website,
	); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Handle logo separately
	if req.Logo != nil {
		if *req.Logo == "" {
			company.RemoveLogo()
		} else {
			company.SetLogo(*req.Logo)
		}
	}
	
	// Save to database
	if err := h.companyRepo.Update(ctx, company); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, company)
}