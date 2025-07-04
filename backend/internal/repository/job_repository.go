package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/masterbrent/electrical-bidding-app/internal/models"
)

// JobRepository defines the interface for job database operations
type JobRepository interface {
	Create(ctx context.Context, job *models.Job) error
	GetByID(ctx context.Context, id string) (*models.Job, error)
	Update(ctx context.Context, job *models.Job) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*models.Job, error)
	GetByCustomerID(ctx context.Context, customerID string) ([]*models.Job, error)
	GetByStatus(ctx context.Context, status models.JobStatus) ([]*models.Job, error)
	
	// Job items operations
	AddJobItem(ctx context.Context, jobID string, item *models.JobItem) error
	GetJobItems(ctx context.Context, jobID string) ([]models.JobItem, error)
	UpdateJobItem(ctx context.Context, item *models.JobItem) error
	RemoveJobItem(ctx context.Context, jobID, itemID string) error
	
	// Job photos operations
	AddPhoto(ctx context.Context, jobID string, photo *models.JobPhoto) error
	GetPhotos(ctx context.Context, jobID string) ([]models.JobPhoto, error)
	RemovePhoto(ctx context.Context, jobID, photoID string) error
}

type jobRepository struct {
	db *sql.DB
}

// NewJobRepository creates a new job repository
func NewJobRepository(db *sql.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(ctx context.Context, job *models.Job) error {
	query := `
		INSERT INTO jobs (
			id, customer_id, template_id, address, status,
			scheduled_date, start_date, end_date, permit_required,
			permit_number, total_amount, notes, wave_invoice_id,
			wave_invoice_url, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		job.ID, job.CustomerID, job.TemplateID, job.Address, job.Status,
		job.ScheduledDate, job.StartDate, job.EndDate, job.PermitRequired,
		job.PermitNumber, job.TotalAmount, job.Notes, job.WaveInvoiceID,
		job.WaveInvoiceURL, job.CreatedAt, job.UpdatedAt,
	)
	
	return err
}

func (r *jobRepository) GetByID(ctx context.Context, id string) (*models.Job, error) {
	query := `
		SELECT 
			id, customer_id, template_id, address, status,
			scheduled_date, start_date, end_date, permit_required,
			permit_number, total_amount, notes, wave_invoice_id,
			wave_invoice_url, created_at, updated_at
		FROM jobs
		WHERE id = $1
	`
	
	job := &models.Job{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&job.ID, &job.CustomerID, &job.TemplateID, &job.Address, &job.Status,
		&job.ScheduledDate, &job.StartDate, &job.EndDate, &job.PermitRequired,
		&job.PermitNumber, &job.TotalAmount, &job.Notes, &job.WaveInvoiceID,
		&job.WaveInvoiceURL, &job.CreatedAt, &job.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("job not found")
	}
	if err != nil {
		return nil, err
	}
	
	// Load items and photos
	job.Items, err = r.GetJobItems(ctx, id)
	if err != nil {
		return nil, err
	}
	
	job.Photos, err = r.GetPhotos(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return job, nil
}

func (r *jobRepository) Update(ctx context.Context, job *models.Job) error {
	query := `
		UPDATE jobs SET
			customer_id = $2, template_id = $3, address = $4, status = $5,
			scheduled_date = $6, start_date = $7, end_date = $8, permit_required = $9,
			permit_number = $10, total_amount = $11, notes = $12, wave_invoice_id = $13,
			wave_invoice_url = $14, updated_at = $15
		WHERE id = $1
	`
	
	result, err := r.db.ExecContext(ctx, query,
		job.ID, job.CustomerID, job.TemplateID, job.Address, job.Status,
		job.ScheduledDate, job.StartDate, job.EndDate, job.PermitRequired,
		job.PermitNumber, job.TotalAmount, job.Notes, job.WaveInvoiceID,
		job.WaveInvoiceURL, job.UpdatedAt,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("job not found")
	}
	
	return nil
}

func (r *jobRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM jobs WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("job not found")
	}
	
	return nil
}

func (r *jobRepository) List(ctx context.Context, limit, offset int) ([]*models.Job, error) {
	query := `
		SELECT 
			id, customer_id, template_id, address, status,
			scheduled_date, start_date, end_date, permit_required,
			permit_number, total_amount, notes, wave_invoice_id,
			wave_invoice_url, created_at, updated_at
		FROM jobs
		ORDER BY scheduled_date DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	jobs := make([]*models.Job, 0)
	for rows.Next() {
		job := &models.Job{}
		err := rows.Scan(
			&job.ID, &job.CustomerID, &job.TemplateID, &job.Address, &job.Status,
			&job.ScheduledDate, &job.StartDate, &job.EndDate, &job.PermitRequired,
			&job.PermitNumber, &job.TotalAmount, &job.Notes, &job.WaveInvoiceID,
			&job.WaveInvoiceURL, &job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Load items and photos for each job
		job.Items, err = r.GetJobItems(ctx, job.ID)
		if err != nil {
			return nil, err
		}
		
		job.Photos, err = r.GetPhotos(ctx, job.ID)
		if err != nil {
			return nil, err
		}
		
		jobs = append(jobs, job)
	}
	
	return jobs, nil
}

func (r *jobRepository) GetByCustomerID(ctx context.Context, customerID string) ([]*models.Job, error) {
	query := `
		SELECT 
			id, customer_id, template_id, address, status,
			scheduled_date, start_date, end_date, permit_required,
			permit_number, total_amount, notes, wave_invoice_id,
			wave_invoice_url, created_at, updated_at
		FROM jobs
		WHERE customer_id = $1
		ORDER BY scheduled_date DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanJobs(ctx, rows)
}

func (r *jobRepository) GetByStatus(ctx context.Context, status models.JobStatus) ([]*models.Job, error) {
	query := `
		SELECT 
			id, customer_id, template_id, address, status,
			scheduled_date, start_date, end_date, permit_required,
			permit_number, total_amount, notes, wave_invoice_id,
			wave_invoice_url, created_at, updated_at
		FROM jobs
		WHERE status = $1
		ORDER BY scheduled_date DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanJobs(ctx, rows)
}

// Job Items operations
func (r *jobRepository) AddJobItem(ctx context.Context, jobID string, item *models.JobItem) error {
	query := `
		INSERT INTO job_items (id, job_id, item_id, name, quantity, price, total, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		item.ID, jobID, item.ItemID, item.Name, item.Quantity, item.Price, item.Total,
	)
	
	return err
}

func (r *jobRepository) GetJobItems(ctx context.Context, jobID string) ([]models.JobItem, error) {
	query := `
		SELECT ji.id, ji.job_id, ji.item_id, ji.name, ji.quantity, ji.price, ji.total,
		       COALESCE(i.nickname, '') as nickname
		FROM job_items ji
		LEFT JOIN items i ON ji.item_id = i.id
		WHERE ji.job_id = $1
		ORDER BY ji.name
	`
	
	rows, err := r.db.QueryContext(ctx, query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	items := make([]models.JobItem, 0)
	for rows.Next() {
		var item models.JobItem
		err := rows.Scan(
			&item.ID, &item.JobID, &item.ItemID, &item.Name,
			&item.Quantity, &item.Price, &item.Total, &item.Nickname,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	
	return items, nil
}

func (r *jobRepository) UpdateJobItem(ctx context.Context, item *models.JobItem) error {
	query := `
		UPDATE job_items SET
			quantity = $2, price = $3, total = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	result, err := r.db.ExecContext(ctx, query,
		item.ID, item.Quantity, item.Price, item.Total,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("job item not found")
	}
	
	return nil
}

func (r *jobRepository) RemoveJobItem(ctx context.Context, jobID, itemID string) error {
	query := `DELETE FROM job_items WHERE job_id = $1 AND id = $2`
	
	result, err := r.db.ExecContext(ctx, query, jobID, itemID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("job item not found")
	}
	
	return nil
}

// Photo operations
func (r *jobRepository) AddPhoto(ctx context.Context, jobID string, photo *models.JobPhoto) error {
	query := `
		INSERT INTO job_photos (id, job_id, url, caption, uploaded_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		photo.ID, jobID, photo.URL, photo.Caption, photo.UploadedAt,
	)
	
	return err
}

func (r *jobRepository) GetPhotos(ctx context.Context, jobID string) ([]models.JobPhoto, error) {
	query := `
		SELECT id, job_id, url, caption, uploaded_at
		FROM job_photos
		WHERE job_id = $1
		ORDER BY uploaded_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	photos := make([]models.JobPhoto, 0)
	for rows.Next() {
		var photo models.JobPhoto
		err := rows.Scan(
			&photo.ID, &photo.JobID, &photo.URL,
			&photo.Caption, &photo.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}
	
	return photos, nil
}

func (r *jobRepository) RemovePhoto(ctx context.Context, jobID, photoID string) error {
	query := `DELETE FROM job_photos WHERE job_id = $1 AND id = $2`
	
	result, err := r.db.ExecContext(ctx, query, jobID, photoID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("photo not found")
	}
	
	return nil
}

// Helper function to scan job rows
func (r *jobRepository) scanJobs(ctx context.Context, rows *sql.Rows) ([]*models.Job, error) {
	jobs := make([]*models.Job, 0)
	for rows.Next() {
		job := &models.Job{}
		err := rows.Scan(
			&job.ID, &job.CustomerID, &job.TemplateID, &job.Address, &job.Status,
			&job.ScheduledDate, &job.StartDate, &job.EndDate, &job.PermitRequired,
			&job.PermitNumber, &job.TotalAmount, &job.Notes, &job.WaveInvoiceID,
			&job.WaveInvoiceURL, &job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Load items and photos for each job
		job.Items, err = r.GetJobItems(ctx, job.ID)
		if err != nil {
			return nil, err
		}
		
		job.Photos, err = r.GetPhotos(ctx, job.ID)
		if err != nil {
			return nil, err
		}
		
		jobs = append(jobs, job)
	}
	
	return jobs, nil
}