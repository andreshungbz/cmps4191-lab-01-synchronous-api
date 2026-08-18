package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("Record not found")
	ErrEditConflict   = errors.New("Edit conflict")
)

// Models groups all database models used in the application.
type Models struct {
	Consumer               ConsumerModel
	APIKey                 APIKeyModel
	Job                    JobModel
	ConsumerActivityReport ConsumerActivityReportModel
}

// NewModels returns all Models configured with the database handler.
func NewModels(db *sql.DB) Models {
	return Models{
		Consumer:               ConsumerModel{DB: db},
		APIKey:                 APIKeyModel{DB: db},
		Job:                    JobModel{DB: db},
		ConsumerActivityReport: ConsumerActivityReportModel{DB: db},
	}
}
