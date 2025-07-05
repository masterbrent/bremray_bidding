package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// JobStatus represents the current status of a job
type JobStatus string

const (
	JobStatusScheduled  JobStatus = "scheduled"
	JobStatusInProgress JobStatus = "in_progress"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusCancelled  JobStatus = "cancelled"
)

// Job represents an electrical job
type Job struct {
	ID             string      `json:"id" db:"id"`
	CustomerID     string      `json:"customerId" db:"customer_id"`
	TemplateID     string      `json:"templateId" db:"template_id"`
	Address        string      `json:"address" db:"address"`
	Status         JobStatus   `json:"status" db:"status"`
	CurrentPhaseID *string     `json:"currentPhaseId,omitempty" db:"current_phase_id"`
	ScheduledDate  time.Time   `json:"scheduledDate" db:"scheduled_date"`
	StartDate      *time.Time  `json:"startDate,omitempty" db:"start_date"`
	EndDate        *time.Time  `json:"endDate,omitempty" db:"end_date"`
	PermitRequired bool        `json:"permitRequired" db:"permit_required"`
	PermitNumber   string      `json:"permitNumber,omitempty" db:"permit_number"`
	TotalAmount    float64     `json:"totalAmount" db:"total_amount"`
	Items          []JobItem   `json:"items"`
	Photos         []JobPhoto  `json:"photos"`
	Notes          string      `json:"notes,omitempty" db:"notes"`
	WaveInvoiceID  string      `json:"waveInvoiceId,omitempty" db:"wave_invoice_id"`
	WaveInvoiceURL string      `json:"waveInvoiceUrl,omitempty" db:"wave_invoice_url"`
	CreatedAt      time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time   `json:"updatedAt" db:"updated_at"`
}

// JobItem represents an item used in a job
type JobItem struct {
	ID       string  `json:"id" db:"id"`
	JobID    string  `json:"jobId" db:"job_id"`
	ItemID   string  `json:"itemId" db:"item_id"`
	Name     string  `json:"name" db:"name"`
	Nickname string  `json:"nickname,omitempty" db:"nickname"`
	Quantity float64 `json:"quantity" db:"quantity"`
	Price    float64 `json:"price" db:"price"`
	Total    float64 `json:"total" db:"total"`
}

// JobPhoto represents a photo associated with a job
type JobPhoto struct {
	ID         string    `json:"id" db:"id"`
	JobID      string    `json:"jobId" db:"job_id"`
	URL        string    `json:"url" db:"url"`
	Caption    string    `json:"caption,omitempty" db:"caption"`
	UploadedAt time.Time `json:"uploadedAt" db:"uploaded_at"`
}

// ValidateJobStatus checks if a job status is valid
func ValidateJobStatus(status JobStatus) bool {
	switch status {
	case JobStatusScheduled, JobStatusInProgress, JobStatusCompleted, JobStatusCancelled:
		return true
	default:
		return false
	}
}

// NewJob creates a new Job with validation
func NewJob(customerID, templateID string, status JobStatus, scheduledDate time.Time) (*Job, error) {
	if customerID == "" {
		return nil, errors.New("customer ID is required")
	}
	if templateID == "" {
		return nil, errors.New("template ID is required")
	}
	if !ValidateJobStatus(status) {
		return nil, errors.New("invalid job status")
	}
	if scheduledDate.Before(time.Now()) {
		return nil, errors.New("scheduled date cannot be in the past")
	}

	now := time.Now()
	return &Job{
		ID:             uuid.New().String(),
		CustomerID:     customerID,
		TemplateID:     templateID,
		Status:         status,
		ScheduledDate:  scheduledDate,
		PermitRequired: false,
		TotalAmount:    0,
		Items:          []JobItem{},
		Photos:         []JobPhoto{},
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// AddPhoto adds a photo to the job
func (j *Job) AddPhoto(photo JobPhoto) {
	photo.JobID = j.ID
	j.Photos = append(j.Photos, photo)
	j.UpdatedAt = time.Now()
}

// UpdateStatus updates the job status with validation
func (j *Job) UpdateStatus(status JobStatus) error {
	if !ValidateJobStatus(status) {
		return errors.New("invalid job status")
	}
	
	j.Status = status
	j.UpdatedAt = time.Now()
	
	// Set start/end dates based on status
	now := time.Now()
	switch status {
	case JobStatusInProgress:
		if j.StartDate == nil {
			j.StartDate = &now
		}
	case JobStatusCompleted:
		if j.StartDate == nil {
			j.StartDate = &now
		}
		j.EndDate = &now
	}
	
	return nil
}

// CalculateTotal calculates the total amount for the job based on items
func (j *Job) CalculateTotal() {
	total := 0.0
	for _, item := range j.Items {
		item.Total = item.Quantity * item.Price
		total += item.Total
	}
	j.TotalAmount = total
	j.UpdatedAt = time.Now()
}

// UpdatePhase updates the current phase of the job
func (j *Job) UpdatePhase(phaseID string) {
	if phaseID == "" {
		j.CurrentPhaseID = nil
	} else {
		j.CurrentPhaseID = &phaseID
	}
	j.UpdatedAt = time.Now()
}