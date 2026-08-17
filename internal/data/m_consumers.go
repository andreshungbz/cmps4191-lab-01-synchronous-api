package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/andreshungbz/cmps4191-lab-01-synchronous-api/internal/validator"
	"github.com/google/uuid"
)

// ConsumerStatus represents the enumerated statuses of a consumer.
type ConsumerStatus string

const (
	ConsumerStatusActive     ConsumerStatus = "active"
	ConsumerStatusSuspended  ConsumerStatus = "suspended"
	ConsumerStatusTerminated ConsumerStatus = "terminated"
)

// Consumer maps the consumer entity.
type Consumer struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Status    ConsumerStatus `json:"status"`
	Version   int            `json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// ValidateContactEmail validates an email against the regular expression.
func ValidateContactEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

// ValidateConsumer perorms validation checks for a consumer record.
func ValidateConsumer(v *validator.Validator, c *Consumer) {
	v.Check(c.Name != "", "name", "must be provided")
	ValidateContactEmail(v, c.Email)
	v.Check(c.Status == ConsumerStatusActive ||
		c.Status == ConsumerStatusSuspended ||
		c.Status == ConsumerStatusTerminated,
		"status", "must be a valid status",
	)
}

// ConsumerModel holds the database handler.
type ConsumerModel struct {
	DB *sql.DB
}

// Insert creates a consumer record.
func (m ConsumerModel) Insert(c *Consumer) error {
	query := `
		INSERT INTO consumers (name, email)
		VALUES ($1, $2)
		RETURNING id, status, version, created_at, updated_at`

	args := []any{c.Name, c.Email}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&c.ID, &c.Status, &c.Version, &c.CreatedAt, &c.UpdatedAt)
}

// GetByID retrieves a single consumer record by its ID.
func (m ConsumerModel) GetByID(id uuid.UUID) (*Consumer, error) {
	query := `
		SELECT id, name, email, status, version, created_at, updated_at
		FROM consumers
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var c Consumer
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.Email, &c.Status, &c.Version, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &c, nil
}

// Update modifies a consumer record by id.
func (m ConsumerModel) Update(c *Consumer) error {
	query := `
		UPDATE consumers
		SET name = $1, email = $2, status = $3, version = version + 1
		WHERE id = $4
		RETURNING version, updated_at`

	args := []any{c.Name, c.Email, c.Status, c.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&c.Version, &c.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

// Delete removes a consumer record by id.
func (m ConsumerModel) Delete(id uuid.UUID) error {
	query := `
		DELETE FROM consumers
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
