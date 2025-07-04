package repository

import (
	"database/sql"
	"time"

	"github.com/masterbrent/electrical-bidding-app/internal/models"
)

type PhotoRepository struct {
	db *sql.DB
}

func NewPhotoRepository(db *sql.DB) *PhotoRepository {
	return &PhotoRepository{db: db}
}

// Create inserts a new photo
func (r *PhotoRepository) Create(photo *models.JobPhoto) error {
	query := `
		INSERT INTO job_photos (id, job_id, url, caption, uploaded_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	photo.UploadedAt = time.Now()
	
	_, err := r.db.Exec(query, 
		photo.ID, 
		photo.JobID, 
		photo.URL, 
		photo.Caption,
		photo.UploadedAt,
	)
	
	return err
}

// GetByID retrieves a photo by ID
func (r *PhotoRepository) GetByID(id string) (*models.JobPhoto, error) {
	query := `
		SELECT id, job_id, url, caption, uploaded_at
		FROM job_photos
		WHERE id = $1
	`
	
	var photo models.JobPhoto
	err := r.db.QueryRow(query, id).Scan(
		&photo.ID,
		&photo.JobID,
		&photo.URL,
		&photo.Caption,
		&photo.UploadedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &photo, nil
}

// GetByJobID retrieves all photos for a job
func (r *PhotoRepository) GetByJobID(jobID string) ([]models.JobPhoto, error) {
	query := `
		SELECT id, job_id, url, caption, uploaded_at
		FROM job_photos
		WHERE job_id = $1
		ORDER BY uploaded_at DESC
	`
	
	rows, err := r.db.Query(query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var photos []models.JobPhoto
	for rows.Next() {
		var photo models.JobPhoto
		err := rows.Scan(
			&photo.ID,
			&photo.JobID,
			&photo.URL,
			&photo.Caption,
			&photo.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}
	
	return photos, nil
}

// Update updates a photo's caption
func (r *PhotoRepository) UpdateCaption(id, caption string) error {
	query := `
		UPDATE job_photos
		SET caption = $2
		WHERE id = $1
	`
	
	_, err := r.db.Exec(query, id, caption)
	return err
}

// Delete removes a photo
func (r *PhotoRepository) Delete(id string) error {
	query := `DELETE FROM job_photos WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// DeleteByJobID removes all photos for a job
func (r *PhotoRepository) DeleteByJobID(jobID string) error {
	query := `DELETE FROM job_photos WHERE job_id = $1`
	_, err := r.db.Exec(query, jobID)
	return err
}
