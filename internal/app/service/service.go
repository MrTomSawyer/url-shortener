package service

import (
	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
	"github.com/jmoiron/sqlx"
)

type ServiceContainer struct {
	URL *urlService
	DB  *sqlx.DB
}

func NewServiceContainer(repo *repository.RepositoryContainer, config config.AppConfig) (*ServiceContainer, error) {
	URLService := urlService{
		Repo:   repo.URLrepo,
		config: config,
	}

	return &ServiceContainer{
		URL: &URLService,
		DB:  repo.Postgres,
	}, nil
}
