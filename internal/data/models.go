package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

// a wrapper for data models
type Models struct {
	Schools SchoolModel
}

// NewModels() to allow us to create a new model
// For abstraction purpose
func NewModels(db *sql.DB) Models {
	return Models{
		Schools: SchoolModel{
			DB: db,
		},
	}
}
