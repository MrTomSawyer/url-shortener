package repository

import (
	"context"

	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/jmoiron/sqlx"
)

type RepoHandler interface {
	Create(shortURL, originalURL string) error
	OriginalURL(shortURL string) (string, error)
}

type RepositoryContainer struct {
	Postgres *sqlx.DB
	URLrepo  RepoHandler
}

func NewRepositoryContainer(db *sqlx.DB, urlRepo RepoHandler) (*RepositoryContainer, error) {

	return &RepositoryContainer{
		Postgres: db,
		URLrepo:  urlRepo,
	}, nil
}

func InitRepository(ctx context.Context, cfg config.AppConfig, db *sqlx.DB) (RepoHandler, error) {
	switch {
	case cfg.DataBase.ConnectionStr != "":
		logger.Log.Infof("Initializing postgres repository. Connection string: %s", cfg.DataBase.ConnectionStr)

		createTableQuery := `
			CREATE TABLE IF NOT EXISTS urls (
				id SERIAL PRIMARY KEY,
				correlationid TEXT,
				shorturl TEXT,
				originalurl TEXT
			);`

		if _, err := db.ExecContext(ctx, createTableQuery); err != nil {
			return nil, err
		}

		uniqueIndexQuery := "CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_originalurl ON urls (originalurl);"
		if _, err := db.ExecContext(ctx, uniqueIndexQuery); err != nil {
			return nil, err
		}

		return NewPostgresURLrepo(ctx, db), nil

	case cfg.Server.TempFolder != "":
		logger.Log.Infof("Initializing file repository")

		fileRepo, err := NewFileURLrepo(cfg.Server.TempFolder)
		if err != nil {
			return nil, err
		}

		return fileRepo, nil

	default:
		logger.Log.Infof("Initializing in-memory repository")

		return NewInMemoryURLRepo(), nil
	}
}
