package models

import (
	"testing"
	"time"
)

func TestNewJob(t *testing.T) {
	tests := []struct {
		name         string
		customerID   string
		templateID   string
		status       JobStatus
		scheduledDate time.Time
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "valid job",
			customerID:   "cust123",
			templateID:   "template123",
			status:       JobStatusScheduled,
			scheduledDate: time.Now().Add(24 * time.Hour),
			wantErr:      false,
		},
		{
			name:         "empty customer ID",
			customerID:   "",
			templateID:   "template123",
			status:       JobStatusScheduled,
			scheduledDate: time.Now().Add(24 * time.Hour),
			wantErr:      true,
			errMsg:       "customer ID is required",
		},
		{
			name:         "empty template ID",
			customerID:   "cust123",
			templateID:   "",
			status:       JobStatusScheduled,
			scheduledDate: time.Now().Add(24 * time.Hour),
			wantErr:      true,
			errMsg:       "template ID is required",
		},
		{
			name:         "invalid status",
			customerID:   "cust123",
			templateID:   "template123",
			status:       "invalid",
			scheduledDate: time.Now().Add(24 * time.Hour),
			wantErr:      true,
			errMsg:       "invalid job status",
		},
		{
			name:         "past scheduled date",
			customerID:   "cust123",
			templateID:   "template123",
			status:       JobStatusScheduled,
			scheduledDate: time.Now().Add(-24 * time.Hour),
			wantErr:      true,
			errMsg:       "scheduled date cannot be in the past",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job, err := NewJob(tt.customerID, tt.templateID, tt.status, tt.scheduledDate)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.errMsg {
					t.Errorf("expected error message %q but got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if job == nil {
					t.Fatal("expected job but got nil")
				}

				if job.CustomerID != tt.customerID {
					t.Errorf("expected customer ID %q but got %q", tt.customerID, job.CustomerID)
				}

				if job.TemplateID != tt.templateID {
					t.Errorf("expected template ID %q but got %q", tt.templateID, job.TemplateID)
				}

				if job.Status != tt.status {
					t.Errorf("expected status %q but got %q", tt.status, job.Status)
				}

				if job.ID == "" {
					t.Error("expected ID to be generated")
				}

				if job.TotalAmount != 0 {
					t.Error("expected initial total amount to be 0")
				}

				if job.PermitRequired {
					t.Error("expected permit required to be false by default")
				}

				if len(job.Photos) != 0 {
					t.Error("expected no photos initially")
				}
			}
		})
	}
}

func TestJob_AddPhoto(t *testing.T) {
	job, _ := NewJob("cust123", "template123", JobStatusScheduled, time.Now().Add(24*time.Hour))
	
	photo := JobPhoto{
		ID:        "photo123",
		URL:       "https://example.com/photo.jpg",
		Caption:   "Test photo",
		UploadedAt: time.Now(),
	}
	
	job.AddPhoto(photo)
	
	if len(job.Photos) != 1 {
		t.Errorf("expected 1 photo but got %d", len(job.Photos))
	}
	
	if job.Photos[0].ID != photo.ID {
		t.Errorf("expected photo ID %q but got %q", photo.ID, job.Photos[0].ID)
	}
}

func TestJob_UpdateStatus(t *testing.T) {
	job, _ := NewJob("cust123", "template123", JobStatusScheduled, time.Now().Add(24*time.Hour))
	
	tests := []struct {
		name    string
		status  JobStatus
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid status change to in progress",
			status:  JobStatusInProgress,
			wantErr: false,
		},
		{
			name:    "valid status change to completed",
			status:  JobStatusCompleted,
			wantErr: false,
		},
		{
			name:    "invalid status",
			status:  "invalid",
			wantErr: true,
			errMsg:  "invalid job status",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := job.UpdateStatus(tt.status)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.errMsg {
					t.Errorf("expected error message %q but got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				
				if job.Status != tt.status {
					t.Errorf("expected status %q but got %q", tt.status, job.Status)
				}
			}
		})
	}
}

func TestJob_CalculateTotal(t *testing.T) {
	job, _ := NewJob("cust123", "template123", JobStatusScheduled, time.Now().Add(24*time.Hour))
	
	items := []JobItem{
		{
			ItemID:   "item1",
			Quantity: 5,
			Price:    10.50,
		},
		{
			ItemID:   "item2",
			Quantity: 3,
			Price:    25.00,
		},
	}
	
	job.Items = items
	job.CalculateTotal()
	
	expectedTotal := (5 * 10.50) + (3 * 25.00) // 52.50 + 75.00 = 127.50
	
	if job.TotalAmount != expectedTotal {
		t.Errorf("expected total amount %.2f but got %.2f", expectedTotal, job.TotalAmount)
	}
}