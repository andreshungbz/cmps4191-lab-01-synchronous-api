package data

import (
	"database/sql"
	"testing"
)

// TestNewModelsWiresDB verifies each model receives the provided DB handle.
func TestNewModelsWiresDB(t *testing.T) {
	db := &sql.DB{}
	models := NewModels(db)

	if models.Consumer.DB != db {
		t.Fatal("Consumer model DB was not wired correctly")
	}

	if models.APIKey.DB != db {
		t.Fatal("APIKey model DB was not wired correctly")
	}

	if models.Job.DB != db {
		t.Fatal("Job model DB was not wired correctly")
	}

	if models.ConsumerActivityReport.DB != db {
		t.Fatal("ConsumerActivityReport model DB was not wired correctly")
	}
}
