package repository

import (
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

func NewRepositoryContainer(cfg config.AppConfig) (*RepositoryContainer, error) {
	var ur RepoHandler
	var db *sqlx.DB

	if cfg.DataBase.ConnectionStr != "" {
		logger.Log.Infof("Initializing postgres repository. Connection string: %s", cfg.DataBase.ConnectionStr)
		db, err := NewPostgresDB(cfg.DataBase.ConnectionStr)
		if err != nil {
			return nil, err
		}
		query := `CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			shortURL TEXT,
			OriginalURL TEXT
		)`
		if _, err := db.Exec(query); err != nil {
			return nil, err
		}
		defer db.Close()
		ur = NewPostgresURLrepo(db)

	} else if cfg.Server.TempFolder != "" {
		logger.Log.Infof("Initializing file repository")
		fileRepo, err := NewFileURLrepo(cfg.Server.TempFolder)
		if err != nil {
			return nil, err
		}
		ur = fileRepo

	} else {
		logger.Log.Infof("Initializing in-memory repository")
		ur = NewInMemoryURLRepo()
	}

	return &RepositoryContainer{
		Postgres: db,
		URLrepo:  ur,
	}, nil
}
