package repository

import "github.com/jmoiron/sqlx"

func NewPostgresDB(connection string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	return db, nil
}
