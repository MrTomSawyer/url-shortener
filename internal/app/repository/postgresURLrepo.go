package repository

import (
	"github.com/jmoiron/sqlx"
)

type PostgresURLrepo struct {
	table    string
	Postgres *sqlx.DB
}

func NewPostgresURLrepo(db *sqlx.DB) *PostgresURLrepo {
	return &PostgresURLrepo{
		table:    "urls",
		Postgres: db,
	}
}

func (u PostgresURLrepo) Create(shortURL, originalURL string) error {
	_, err := u.Postgres.Exec("INSERT INTO %s (shortURL, OriginalURL) VALUES ($1, $2)", u.table)
	if err != nil {
		return err
	}
	return nil
}

func (u PostgresURLrepo) OriginalURL(shortURL string) (string, error) {
	row := u.Postgres.QueryRow("SELECT originalURL from %s WHERE shortURL=$1", u.table)
	var originalURL string
	err := row.Scan(&originalURL)
	if err != nil {
		return "", nil
	}
	return originalURL, nil
}
