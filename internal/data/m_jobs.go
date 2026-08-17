package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/andreshungbz/cmps4191-lab-01-synchronous-api/internal/validator"
	"github.com/google/uuid"
)

// JobStatus represents the enumerated statuses of a job.
type JobStatus string

const (
	JobStatusQueued     JobStatus = "queued"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
	JobStatusCancelled  JobStatus = "cancelled"
)

// Job maps the jobs entity.
type Job struct {
	ID           uuid.UUID        `json:"id"`
	PublicID     uuid.UUID        `json:"public_id"`
	ConsumerID   uuid.UUID        `json:"consumer_id"`
	JobType      string           `json:"job_type"`
	Status       JobStatus        `json:"status"`
	Payload      json.RawMessage  `json:"payload"`
	Result       *json.RawMessage `json:"result,omitempty"`
	ErrorMessage *string          `json:"error_message,omitempty"`
	StartedAt    *time.Time       `json:"started_at,omitempty"`
	CompletedAt  *time.Time       `json:"completed_at,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// ValidateJob performs validation checks for a job record.
func ValidateJob(v *validator.Validator, j *Job) {
	v.Check(j.ConsumerID != uuid.Nil, "consumer_id", "must be provided")
	v.Check(j.JobType != "", "job_type", "must be provided")
	v.Check(j.Payload != nil && json.Valid(j.Payload), "payload", "must be a valid JSON object")
	if j.Result != nil {
		v.Check(json.Valid(*j.Result), "result", "must be a valid JSON object")
	}
}

// ValidateJobStatus performs validation checks for a job status.
func ValidateJobStatus(v *validator.Validator, j *Job) {
	v.Check(j.Status == JobStatusQueued ||
		j.Status == JobStatusProcessing ||
		j.Status == JobStatusCompleted ||
		j.Status == JobStatusFailed ||
		j.Status == JobStatusCancelled,
		"status", "must be a valid status",
	)
}

// JobModel holds the database handler.
type JobModel struct {
	DB *sql.DB
}

// Insert creates a job record.
func (m JobModel) Insert(j *Job) error {
	query := `
		INSERT INTO jobs (consumer_id, job_type, payload)
		VALUES ($1, $2, $3)
		RETURNING id, public_id, status, result, error_message, started_at, completed_at, created_at, updated_at`

	// fallback to empty JSON object if payload is empty
	payload := j.Payload
	if payload == nil {
		payload = json.RawMessage("{}")
	}

	args := []any{j.ConsumerID, j.JobType, payload}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&j.ID, &j.PublicID, &j.Status, &j.Result, &j.ErrorMessage, &j.StartedAt, &j.CompletedAt, &j.CreatedAt, &j.UpdatedAt)
}

// GetByID retrieves a single job record by its internal UUID.
func (m JobModel) GetByID(id uuid.UUID) (*Job, error) {
	query := `
		SELECT id, public_id, consumer_id, job_type, status, payload, result, error_message, started_at, completed_at, created_at, updated_at
		FROM jobs
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var j Job
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&j.ID,
		&j.PublicID,
		&j.ConsumerID,
		&j.JobType,
		&j.Status,
		&j.Payload,
		&j.Result,
		&j.ErrorMessage,
		&j.StartedAt,
		&j.CompletedAt,
		&j.CreatedAt,
		&j.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &j, nil
}

// GetByPublicID retrieves a single job record by its public UUID.
func (m JobModel) GetByPublicID(publicID uuid.UUID) (*Job, error) {
	query := `
		SELECT id, public_id, consumer_id, job_type, status, payload, result, error_message, started_at, completed_at, created_at, updated_at
		FROM jobs
		WHERE public_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var j Job
	err := m.DB.QueryRowContext(ctx, query, publicID).Scan(
		&j.ID,
		&j.PublicID,
		&j.ConsumerID,
		&j.JobType,
		&j.Status,
		&j.Payload,
		&j.Result,
		&j.ErrorMessage,
		&j.StartedAt,
		&j.CompletedAt,
		&j.CreatedAt,
		&j.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &j, nil
}

// GetMostRecent retrieves the single most recently created queued job for a consumer.
func (m JobModel) GetMostRecent(consumerID uuid.UUID) (*Job, error) {
	query := `
		SELECT id, public_id, consumer_id, job_type, status, payload, result, error_message, started_at, completed_at, created_at, updated_at
		FROM jobs
		WHERE consumer_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT 1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var j Job
	err := m.DB.QueryRowContext(ctx, query, consumerID, JobStatusQueued).Scan(
		&j.ID,
		&j.PublicID,
		&j.ConsumerID,
		&j.JobType,
		&j.Status,
		&j.Payload,
		&j.Result,
		&j.ErrorMessage,
		&j.StartedAt,
		&j.CompletedAt,
		&j.CreatedAt,
		&j.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &j, nil
}

// Update modifies a job record by id.
func (m JobModel) Update(j *Job) error {
	query := `
		UPDATE jobs
		SET status = $1, result = $2, error_message = $3, started_at = $4, completed_at = $5
		WHERE id = $6`

	args := []any{j.Status, j.Result, j.ErrorMessage, j.StartedAt, j.CompletedAt, j.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Delete removes a job record by id.
func (m JobModel) Delete(id uuid.UUID) error {
	query := `
		DELETE FROM jobs
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
