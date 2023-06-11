package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MrTomSawyer/url-shortener/internal/app/apperrors"
	"github.com/jmoiron/sqlx"
)

type PostgresURLrepo struct {
	Table    string
	Postgres *sqlx.DB
	ctx      context.Context
}

func NewPostgresURLrepo(ctx context.Context, db *sqlx.DB) *PostgresURLrepo {
	return &PostgresURLrepo{
		Table:    "urls",
		Postgres: db,
		ctx:      ctx,
	}
}

func (u PostgresURLrepo) Create(shortURL, originalURL string) error {
	cxt, cancel := context.WithTimeout(u.ctx, 2000*time.Millisecond)
	defer cancel()

	query := fmt.Sprintf("INSERT INTO %s (shortURL, OriginalURL) VALUES ($1, $2) ON CONFLICT (OriginalURL) DO NOTHING RETURNING id", u.Table)
	row := u.Postgres.QueryRowContext(cxt, query, shortURL, originalURL)
	var res string
	err := row.Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			query := fmt.Sprintf("SELECT shortURL FROM %s WHERE OriginalURL=$1", u.Table)
			row := u.Postgres.QueryRowContext(cxt, query, originalURL)
			err := row.Scan(&res)
			if err != nil {
				return err
			}
			return apperrors.NewURLConflict(res, err)
		}
		return err
	}
	return nil
}

func (u PostgresURLrepo) OriginalURL(shortURL string) (string, error) {
	cxt, cancel := context.WithTimeout(u.ctx, 2000*time.Millisecond)
	defer cancel()

	query := fmt.Sprintf("SELECT originalurl from %s WHERE shorturl=$1", u.Table)
	row := u.Postgres.QueryRowContext(cxt, query, shortURL)
	var originalURL string
	err := row.Scan(&originalURL)
	if err != nil {
		return "", nil
	}
	return originalURL, nil
}
