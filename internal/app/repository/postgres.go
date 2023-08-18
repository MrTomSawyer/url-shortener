// Package repository provides implementations for data storage and retrieval.
package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPostgresDB creates a new instance of a PostgreSQL database connection.
func NewPostgresDB(connection string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
