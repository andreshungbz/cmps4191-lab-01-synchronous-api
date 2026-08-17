package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/andreshungbz/cmps4191-lab-01-synchronous-api/internal/validator"
	"github.com/google/uuid"
)

// KeyStatus represents the enumerated statuses of an API key.
type KeyStatus string

const (
	KeyStatusActive   KeyStatus = "active"
	KeyStatusRotating KeyStatus = "rotating"
	KeyStatusRevoked  KeyStatus = "revoked"
)

// APIKey maps the api_keys entity.
type APIKey struct {
	ID         uuid.UUID  `json:"id"`
	ConsumerID uuid.UUID  `json:"consumer_id"`
	KeyHash    string     `json:"-"`
	KeyPrefix  string     `json:"key_prefix"`
	Status     KeyStatus  `json:"status"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ValidateAPIKey performs validation checks for an API key record.
func ValidateAPIKey(v *validator.Validator, key *APIKey) {
	v.Check(key.ConsumerID != uuid.Nil, "consumer_id", "must be provided")
	v.Check(key.KeyHash != "", "key_hash", "must be provided")
	v.Check(key.KeyPrefix != "", "key_prefix", "must be provided")
	v.Check(key.Status == KeyStatusActive ||
		key.Status == KeyStatusRotating ||
		key.Status == KeyStatusRevoked,
		"status", "must be a valid status",
	)
	if key.ExpiresAt != nil {
		v.Check(key.ExpiresAt.After(time.Now()), "expires_at", "must be in the future")
	}
}

// APIKeyModel holds the database handler.
type APIKeyModel struct {
	DB *sql.DB
}

// Insert creates an API key record.
func (m APIKeyModel) Insert(key *APIKey) error {
	query := `
		INSERT INTO api_keys (consumer_id, key_hash, key_prefix, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, status, created_at`

	args := []any{
		key.ConsumerID,
		key.KeyHash,
		key.KeyPrefix,
		key.ExpiresAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&key.ID, &key.Status, &key.CreatedAt)
}

// GetByID retrieves a single API key record by its ID.
func (m APIKeyModel) GetByID(id uuid.UUID) (*APIKey, error) {
	query := `
		SELECT id, consumer_id, key_hash, key_prefix, status, last_used_at, expires_at, created_at
		FROM api_keys
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var key APIKey
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&key.ID,
		&key.ConsumerID,
		&key.KeyHash,
		&key.KeyPrefix,
		&key.Status,
		&key.LastUsedAt,
		&key.ExpiresAt,
		&key.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &key, nil
}

// GetByHash retrieves a single API key record by its hash (used for authentication).
func (m APIKeyModel) GetByHash(keyHash string) (*APIKey, error) {
	query := `
		SELECT id, consumer_id, key_hash, key_prefix, status, last_used_at, expires_at, created_at
		FROM api_keys
		WHERE key_hash = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var key APIKey
	err := m.DB.QueryRowContext(ctx, query, keyHash).Scan(
		&key.ID,
		&key.ConsumerID,
		&key.KeyHash,
		&key.KeyPrefix,
		&key.Status,
		&key.LastUsedAt,
		&key.ExpiresAt,
		&key.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &key, nil
}

// Update modifies an API key record by id.
func (m APIKeyModel) Update(key *APIKey) error {
	query := `
		UPDATE api_keys
		SET status = $1
		WHERE id = $2`

	args := []any{key.Status, key.ID}

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

// Delete removes an API key record by id.
func (m APIKeyModel) Delete(id uuid.UUID) error {
	query := `
		DELETE FROM api_keys
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
