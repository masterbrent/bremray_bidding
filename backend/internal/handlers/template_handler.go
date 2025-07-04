package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/masterbrent/electrical-bidding-app/internal/models"
	"github.com/masterbrent/electrical-bidding-app/internal/repository"
)

// TemplateHandler handles HTTP requests for job templates
type TemplateHandler struct {
	templateRepo repository.JobTemplateRepository
	itemRepo     repository.ItemRepository
}

// NewTemplateHandler creates a new template handler
func NewTemplateHandler(templateRepo repository.JobTemplateRepository, itemRepo repository.ItemRepository) *TemplateHandler {
	return &TemplateHandler{
		templateRepo: templateRepo,
		itemRepo:     itemRepo,
	}
}

// RegisterRoutes registers all template routes
func (h *TemplateHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/templates", h.List).Methods("GET", "OPTIONS")
	router.HandleFunc("/templates", h.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/templates/{id}", h.Get).Methods("GET", "OPTIONS")
	router.HandleFunc("/templates/{id}", h.Update).Methods("PUT", "OPTIONS")
	router.HandleFunc("/templates/{id}", h.Delete).Methods("DELETE", "OPTIONS")
	
	// Template items
	router.HandleFunc("/templates/{id}/items", h.AddItem).Methods("POST", "OPTIONS")
	router.HandleFunc("/templates/{id}/items/{itemId}", h.UpdateItem).Methods("PUT", "OPTIONS")
	router.HandleFunc("/templates/{id}/items/{itemId}", h.RemoveItem).Methods("DELETE", "OPTIONS")
}

// List returns all templates
func (h *TemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Check if only active templates are requested
	if active := r.URL.Query().Get("active"); active == "true" {
		templates, err := h.templateRepo.ListActive(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondJSON(w, templates)
		return
	}
	
	// Parse pagination parameters
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
	
	templates, err := h.templateRepo.List(ctx, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, templates)
}

// Get returns a single template by ID
func (h *TemplateHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	
	template, err := h.templateRepo.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "template not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, template)
}

// Create creates a new template
func (h *TemplateHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Items       []struct {
			ItemID          string  `json:"itemId"`
			DefaultQuantity float64 `json:"defaultQuantity"`
		} `json:"items"`
		Phases      []struct {
			Name        string `json:"name"`
			Order       int    `json:"order"`
			Description string `json:"description,omitempty"`
		} `json:"phases"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate items exist
	templateItems := make([]models.TemplateItem, 0, len(req.Items))
	for _, reqItem := range req.Items {
		// Verify item exists
		if _, err := h.itemRepo.GetByID(ctx, reqItem.ItemID); err != nil {
			if err.Error() == "item not found" {
				http.Error(w, "Item "+reqItem.ItemID+" not found", http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		templateItems = append(templateItems, models.TemplateItem{
			ItemID:          reqItem.ItemID,
			DefaultQuantity: reqItem.DefaultQuantity,
		})
	}
	
	// Convert phases
	templatePhases := make([]models.TemplatePhase, 0, len(req.Phases))
	for _, reqPhase := range req.Phases {
		templatePhases = append(templatePhases, models.TemplatePhase{
			Name:        reqPhase.Name,
			Order:       reqPhase.Order,
			Description: reqPhase.Description,
		})
	}
	
	// Create template
	template, err := models.NewJobTemplate(req.Name, req.Description, templateItems, templatePhases)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.templateRepo.Create(ctx, template); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	respondJSON(w, template)
}

// Update updates an existing template
func (h *TemplateHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	
	// Get existing template
	template, err := h.templateRepo.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "template not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsActive    *bool  `json:"isActive"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Update fields
	if req.Name != "" {
		template.Name = req.Name
	}
	
	if req.Description != "" {
		template.Description = req.Description
	}
	
	if req.IsActive != nil {
		if *req.IsActive {
			template.Activate()
		} else {
			template.Deactivate()
		}
	}
	
	if err := h.templateRepo.Update(ctx, template); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, template)
}

// Delete deletes a template
func (h *TemplateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	
	if err := h.templateRepo.Delete(ctx, id); err != nil {
		if err.Error() == "template not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// AddItem adds an item to a template
func (h *TemplateHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	templateID := vars["id"]
	
	// Verify template exists
	template, err := h.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		if err.Error() == "template not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	var req struct {
		ItemID          string  `json:"itemId"`
		DefaultQuantity float64 `json:"defaultQuantity"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Verify item exists
	if _, err := h.itemRepo.GetByID(ctx, req.ItemID); err != nil {
		if err.Error() == "item not found" {
			http.Error(w, "Item not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Create template item
	templateItem := &models.TemplateItem{
		ID:              uuid.New().String(),
		TemplateID:      templateID,
		ItemID:          req.ItemID,
		DefaultQuantity: req.DefaultQuantity,
	}
	
	if err := h.templateRepo.AddTemplateItem(ctx, templateID, templateItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	template.AddItem(*templateItem)
	
	w.WriteHeader(http.StatusCreated)
	respondJSON(w, templateItem)
}

// UpdateItem updates a template item
func (h *TemplateHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	templateID := vars["id"]
	itemID := vars["itemId"]
	
	// Get template to verify it exists
	template, err := h.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		if err.Error() == "template not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	var req struct {
		DefaultQuantity float64 `json:"defaultQuantity"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Find and update the item
	var updatedItem *models.TemplateItem
	for i := range template.Items {
		if template.Items[i].ItemID == itemID {
			template.Items[i].DefaultQuantity = req.DefaultQuantity
			updatedItem = &template.Items[i]
			break
		}
	}
	
	if updatedItem == nil {
		http.Error(w, "Template item not found", http.StatusNotFound)
		return
	}
	
	if err := h.templateRepo.UpdateTemplateItem(ctx, updatedItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, updatedItem)
}

// RemoveItem removes an item from a template
func (h *TemplateHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	templateID := vars["id"]
	itemID := vars["itemId"]
	
	// Get template
	template, err := h.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		if err.Error() == "template not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Use template's RemoveItem method to validate business rules
	if err := template.RemoveItem(itemID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.templateRepo.RemoveTemplateItem(ctx, templateID, itemID); err != nil {
		if err.Error() == "template item not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}