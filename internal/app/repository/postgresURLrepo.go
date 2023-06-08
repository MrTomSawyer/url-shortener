package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PostgresURLrepo struct {
	Table    string
	Postgres *sqlx.DB
}

func NewPostgresURLrepo(db *sqlx.DB) *PostgresURLrepo {
	return &PostgresURLrepo{
		Table:    "urls",
		Postgres: db,
	}
}

func (u PostgresURLrepo) Create(shortURL, originalURL string) error {
	query := fmt.Sprintf("INSERT INTO %s (shortURL, OriginalURL) VALUES ($1, $2) RETURNING id", u.Table)
	row := u.Postgres.QueryRow(query, shortURL, originalURL)
	var createdRow string
	err := row.Scan(&createdRow)
	if err != nil {
		return err
	}
	return nil
}

func (u PostgresURLrepo) OriginalURL(shortURL string) (string, error) {
	query := fmt.Sprintf("SELECT originalurl from %s WHERE shorturl=$1", u.Table)
	row := u.Postgres.QueryRow(query, shortURL)
	var originalURL string
	err := row.Scan(&originalURL)
	if err != nil {
		return "", nil
	}
	return originalURL, nil
}
