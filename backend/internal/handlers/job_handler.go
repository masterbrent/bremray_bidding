package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/masterbrent/electrical-bidding-app/internal/models"
	"github.com/masterbrent/electrical-bidding-app/internal/repository"
)

// JobHandler handles HTTP requests for jobs
type JobHandler struct {
	jobRepo      repository.JobRepository
	customerRepo repository.CustomerRepository
	templateRepo repository.JobTemplateRepository
	itemRepo     repository.ItemRepository
}

// NewJobHandler creates a new job handler
func NewJobHandler(
	jobRepo repository.JobRepository,
	customerRepo repository.CustomerRepository,
	templateRepo repository.JobTemplateRepository,
	itemRepo repository.ItemRepository,
) *JobHandler {
	return &JobHandler{
		jobRepo:      jobRepo,
		customerRepo: customerRepo,
		templateRepo: templateRepo,
		itemRepo:     itemRepo,
	}
}

// RegisterRoutes registers all job routes
func (h *JobHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/jobs", h.List).Methods("GET", "OPTIONS")
	router.HandleFunc("/jobs", h.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/jobs/{id}", h.Get).Methods("GET", "OPTIONS")
	router.HandleFunc("/jobs/{id}", h.Update).Methods("PUT", "OPTIONS")
	router.HandleFunc("/jobs/{id}", h.Delete).Methods("DELETE", "OPTIONS")
	
	// Job items
	router.HandleFunc("/jobs/{id}/items", h.AddItem).Methods("POST", "OPTIONS")
	router.HandleFunc("/jobs/{id}/items/{itemId}", h.UpdateItem).Methods("PUT", "OPTIONS")
	router.HandleFunc("/jobs/{id}/items/{itemId}", h.RemoveItem).Methods("DELETE", "OPTIONS")
	
	// Job photos
	router.HandleFunc("/jobs/{id}/photos", h.AddPhoto).Methods("POST", "OPTIONS")
	router.HandleFunc("/jobs/{id}/photos/{photoId}", h.RemovePhoto).Methods("DELETE", "OPTIONS")
}

// List returns all jobs
func (h *JobHandler) List(w http.ResponseWriter, r *http.Request) {
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
	
	// Check for status filter
	if status := r.URL.Query().Get("status"); status != "" {
		jobs, err := h.jobRepo.GetByStatus(ctx, models.JobStatus(status))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondJSON(w, jobs)
		return
	}
	
	// Check for customer filter
	if customerID := r.URL.Query().Get("customerId"); customerID != "" {
		jobs, err := h.jobRepo.GetByCustomerID(ctx, customerID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondJSON(w, jobs)
		return
	}
	
	// Get all jobs with pagination
	jobs, err := h.jobRepo.List(ctx, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, jobs)
}

// Get returns a single job by ID
func (h *JobHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	
	job, err := h.jobRepo.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "job not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, job)
}

// Create creates a new job
func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var req struct {
		CustomerID    string     `json:"customerId"`
		TemplateID    string     `json:"templateId"`
		Address       string     `json:"address"`
		ScheduledDate *time.Time `json:"scheduledDate,omitempty"`
		Notes         string     `json:"notes"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding job creation request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	log.Printf("Creating job - CustomerID: %s, TemplateID: %s, Address: %s", 
		req.CustomerID, req.TemplateID, req.Address)
	
	// Validate required fields
	if req.CustomerID == "" || req.TemplateID == "" || req.Address == "" {
		log.Printf("Missing required fields - CustomerID: %s, TemplateID: %s, Address: %s",
			req.CustomerID, req.TemplateID, req.Address)
		http.Error(w, "Customer ID, Template ID, and Address are required", http.StatusBadRequest)
		return
	}
	
	// Validate customer exists
	if _, err := h.customerRepo.GetByID(ctx, req.CustomerID); err != nil {
		if err.Error() == "customer not found" {
			log.Printf("Customer not found: %s", req.CustomerID)
			http.Error(w, "Customer not found", http.StatusBadRequest)
			return
		}
		log.Printf("Error fetching customer %s: %v", req.CustomerID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Validate template exists and get it
	template, err := h.templateRepo.GetByID(ctx, req.TemplateID)
	if err != nil {
		if err.Error() == "template not found" {
			log.Printf("Template not found: %s", req.TemplateID)
			http.Error(w, "Template not found", http.StatusBadRequest)
			return
		}
		log.Printf("Error fetching template %s: %v", req.TemplateID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Create job
	// Default to tomorrow if no scheduled date provided to avoid "past date" validation error
	scheduledDate := time.Now().Add(24 * time.Hour)
	if req.ScheduledDate != nil {
		scheduledDate = *req.ScheduledDate
	}
	
	job, err := models.NewJob(req.CustomerID, req.TemplateID, models.JobStatusScheduled, scheduledDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	job.Address = req.Address
	job.Notes = req.Notes
	
	// Create job in database
	if err := h.jobRepo.Create(ctx, job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Add items from template
	for _, templateItem := range template.Items {
		// Get item details
		item, err := h.itemRepo.GetByID(ctx, templateItem.ItemID)
		if err != nil {
			continue // Skip if item not found
		}
		
		jobItem := &models.JobItem{
			ID:       uuid.New().String(),
			JobID:    job.ID,
			ItemID:   item.ID,
			Name:     item.Name,
			Quantity: 0, // Always start with 0 quantity - techs will increment as they install
			Price:    item.UnitPrice,
			Total:    0, // Total is 0 since quantity is 0
		}
		
		if err := h.jobRepo.AddJobItem(ctx, job.ID, jobItem); err != nil {
			// Log error but continue
			continue
		}
		
		job.Items = append(job.Items, *jobItem)
	}
	
	// Calculate total
	job.CalculateTotal()
	h.jobRepo.Update(ctx, job)
	
	w.WriteHeader(http.StatusCreated)
	respondJSON(w, job)
}

// Update updates an existing job
func (h *JobHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	
	// Get existing job
	job, err := h.jobRepo.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "job not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	var req struct {
		Address        string     `json:"address"`
		Status         string     `json:"status"`
		ScheduledDate  *time.Time `json:"scheduledDate"`
		PermitRequired *bool      `json:"permitRequired"`
		PermitNumber   string     `json:"permitNumber"`
		Notes          string     `json:"notes"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Update fields
	if req.Address != "" {
		job.Address = req.Address
	}
	
	if req.Status != "" {
		if err := job.UpdateStatus(models.JobStatus(req.Status)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	
	if req.ScheduledDate != nil {
		job.ScheduledDate = *req.ScheduledDate
	}
	
	if req.PermitRequired != nil {
		job.PermitRequired = *req.PermitRequired
	}
	
	if req.PermitNumber != "" {
		job.PermitNumber = req.PermitNumber
	}
	
	if req.Notes != "" {
		job.Notes = req.Notes
	}
	
	job.UpdatedAt = time.Now()
	
	if err := h.jobRepo.Update(ctx, job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, job)
}

// Delete deletes a job
func (h *JobHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	
	if err := h.jobRepo.Delete(ctx, id); err != nil {
		if err.Error() == "job not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// AddItem adds an item to a job
func (h *JobHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	jobID := vars["id"]
	
	// Verify job exists
	job, err := h.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		if err.Error() == "job not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	var req struct {
		ItemID   string  `json:"itemId"`
		Quantity float64 `json:"quantity"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Get item details
	item, err := h.itemRepo.GetByID(ctx, req.ItemID)
	if err != nil {
		if err.Error() == "item not found" {
			http.Error(w, "Item not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Create job item
	jobItem := &models.JobItem{
		ID:       uuid.New().String(),
		JobID:    jobID,
		ItemID:   item.ID,
		Name:     item.Name,
		Quantity: req.Quantity,
		Price:    item.UnitPrice,
		Total:    req.Quantity * item.UnitPrice,
	}
	
	if err := h.jobRepo.AddJobItem(ctx, jobID, jobItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Update job total
	job.Items = append(job.Items, *jobItem)
	job.CalculateTotal()
	h.jobRepo.Update(ctx, job)
	
	w.WriteHeader(http.StatusCreated)
	respondJSON(w, jobItem)
}

// UpdateItem updates a job item
func (h *JobHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	jobID := vars["id"]
	itemID := vars["itemId"]
	
	// Get job to verify it exists
	job, err := h.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		if err.Error() == "job not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	var req struct {
		Quantity float64 `json:"quantity"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Find and update the item
	var updatedItem *models.JobItem
	for i := range job.Items {
		if job.Items[i].ID == itemID {
			job.Items[i].Quantity = req.Quantity
			job.Items[i].Total = req.Quantity * job.Items[i].Price
			updatedItem = &job.Items[i]
			break
		}
	}
	
	if updatedItem == nil {
		http.Error(w, "Job item not found", http.StatusNotFound)
		return
	}
	
	if err := h.jobRepo.UpdateJobItem(ctx, updatedItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Update job total
	job.CalculateTotal()
	h.jobRepo.Update(ctx, job)
	
	respondJSON(w, updatedItem)
}

// RemoveItem removes an item from a job
func (h *JobHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	jobID := vars["id"]
	itemID := vars["itemId"]
	
	// Get job
	job, err := h.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		if err.Error() == "job not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if err := h.jobRepo.RemoveJobItem(ctx, jobID, itemID); err != nil {
		if err.Error() == "job item not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Update job items and total
	newItems := make([]models.JobItem, 0, len(job.Items)-1)
	for _, item := range job.Items {
		if item.ID != itemID {
			newItems = append(newItems, item)
		}
	}
	job.Items = newItems
	job.CalculateTotal()
	h.jobRepo.Update(ctx, job)
	
	w.WriteHeader(http.StatusNoContent)
}

// AddPhoto adds a photo to a job
func (h *JobHandler) AddPhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	jobID := vars["id"]
	
	// Verify job exists
	job, err := h.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		if err.Error() == "job not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	var req struct {
		URL     string `json:"url"`
		Caption string `json:"caption"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	photo := &models.JobPhoto{
		ID:         uuid.New().String(),
		JobID:      jobID,
		URL:        req.URL,
		Caption:    req.Caption,
		UploadedAt: time.Now(),
	}
	
	if err := h.jobRepo.AddPhoto(ctx, jobID, photo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	job.AddPhoto(*photo)
	
	w.WriteHeader(http.StatusCreated)
	respondJSON(w, photo)
}

// RemovePhoto removes a photo from a job
func (h *JobHandler) RemovePhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	jobID := vars["id"]
	photoID := vars["photoId"]
	
	if err := h.jobRepo.RemovePhoto(ctx, jobID, photoID); err != nil {
		if err.Error() == "photo not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}