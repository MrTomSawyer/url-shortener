package repository

import (
	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/jmoiron/sqlx"
)

type repoHandler interface {
	Create(shortURL, originalURL string) error
	OriginalURL(shortURL string) (string, error)
}

type RepositoryContainer struct {
	Postgres *sqlx.DB
	URLrepo  repoHandler
}

type repoMode string

const (
	InMemory repoMode = "InMemory"
	File     repoMode = "File"
	DB       repoMode = "DB"
)

func NewRepositoryContainer(db *sqlx.DB, cfg config.AppConfig) *RepositoryContainer {
	var ur repoHandler

	if cfg.DataBase.ConnectionStr != "" {
		ur = NewPostgresURLrepo(db)
	} else if cfg.Server.TempFolder != "" {
		ur = NewFileURLrepo(cfg.Server.TempFolder)
	} else {
		ur = NewInMemoryURLRepo()
	}

	return &RepositoryContainer{
		Postgres: db,
		URLrepo:  ur,
	}
}

func (r RepositoryContainer) InitTables() error {
	query := `CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		shortURL TEXT,
		OriginalURL TEXT
	)`

	if _, err := r.Postgres.Exec(query); err != nil {
		return err
	}
	return nil
}
