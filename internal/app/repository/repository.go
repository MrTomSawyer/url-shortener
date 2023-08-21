// Package repository provides implementations for data storage and retrieval.
package repository

import (
	"context"

	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/MrTomSawyer/url-shortener/internal/app/models"
	"github.com/jmoiron/sqlx"
)

// RepoHandler defines the interface for URL repository operations.
type RepoHandler interface {
	Create(shortURL, originalURL, userID string) error
	OriginalURL(shortURL string) (string, error)
	BatchCreate(data []models.TempURLBatchRequest, userID string) ([]models.BatchURLResponse, error)
	GetAll(userID string) ([]models.URLJsonResponse, error)
	DeleteAll(shortURLs []string, userID string) error
}

// RepositoryContainer holds the database connection and the URL repository implementation.
type RepositoryContainer struct {
	Postgres *sqlx.DB
	URLRepo  RepoHandler
}

// NewRepositoryContainer creates a new RepositoryContainer instance.
func NewRepositoryContainer(db *sqlx.DB, urlRepo RepoHandler) (*RepositoryContainer, error) {
	return &RepositoryContainer{
		Postgres: db,
		URLRepo:  urlRepo,
	}, nil
}

// InitRepository initializes the appropriate repository based on the configuration.
func InitRepository(ctx context.Context, cfg config.AppConfig, db *sqlx.DB) (RepoHandler, error) {
	switch {
	case cfg.DataBase.ConnectionStr != "":
		logger.Log.Infof("Initializing postgres repository. Connection string: %s", cfg.DataBase.ConnectionStr)

		createTableQuery := `
			CREATE TABLE IF NOT EXISTS urls (
				id SERIAL PRIMARY KEY,
				correlationid TEXT,
				shorturl TEXT,
				userid TEXT,
				originalurl TEXT,
				isdeleted BOOLEAN DEFAULT false
			);`

		if _, err := db.ExecContext(ctx, createTableQuery); err != nil {
			return nil, err
		}

		uniqueIndexQuery := "CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_originalurl ON urls (originalurl);"
		if _, err := db.ExecContext(ctx, uniqueIndexQuery); err != nil {
			return nil, err
		}

		urlsCh := make(chan models.UserURLs)
		pgRepo := NewPostgresURLrepo(ctx, db, cfg, urlsCh)
		go WorkerDeleteURLs(urlsCh, db)
		return pgRepo, nil

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
