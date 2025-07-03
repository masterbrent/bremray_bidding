package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/masterbrent/electrical-bidding-app/internal/models"
	"github.com/masterbrent/electrical-bidding-app/internal/repository"
)

// ItemService defines the interface for item business logic
type ItemService interface {
	Create(ctx context.Context, name, unit string, unitPrice float64, category string) (*models.Item, error)
	GetByID(ctx context.Context, id string) (*models.Item, error)
	List(ctx context.Context) ([]*models.Item, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) (*models.Item, error)
	Delete(ctx context.Context, id string) error
}

// ItemHandler handles HTTP requests for items
type ItemHandler struct {
	service ItemService
}

// NewItemHandler creates a new item handler
func NewItemHandler(service ItemService) *ItemHandler {
	return &ItemHandler{
		service: service,
	}
}

// Create handles POST /api/items
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name      string  `json:"name"`
		Unit      string  `json:"unit"`
		UnitPrice float64 `json:"unitPrice"`
		Category  string  `json:"category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.service.Create(r.Context(), req.Name, req.Unit, req.UnitPrice, req.Category)
	if err != nil {
		if errors.Is(err, models.ErrValidation) || err.Error() == "item name is required" || 
		   err.Error() == "unit is required" || err.Error() == "unit price must be positive" {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	respondWithJSON(w, http.StatusCreated, item)
}

// GetByID handles GET /api/items/:id
func (h *ItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		respondWithError(w, http.StatusNotFound, "Item not found")
		return
	}

	item, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Item not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to get item")
		return
	}

	respondWithJSON(w, http.StatusOK, item)
}

// List handles GET /api/items
func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	respondWithJSON(w, http.StatusOK, items)
}

// Update handles PUT /api/items/:id
func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		respondWithError(w, http.StatusNotFound, "Item not found")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.service.Update(r.Context(), id, updates)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Item not found")
			return
		}
		// Check for validation errors
		errMsg := err.Error()
		if errMsg == "item name is required" || errMsg == "unit is required" || 
		   errMsg == "unit price must be positive" || errMsg == "unit price cannot be negative" {
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to update item")
		return
	}

	respondWithJSON(w, http.StatusOK, item)
}

// Delete handles DELETE /api/items/:id
func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		respondWithError(w, http.StatusNotFound, "Item not found")
		return
	}

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Item not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to delete item")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper functions
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to marshal response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}