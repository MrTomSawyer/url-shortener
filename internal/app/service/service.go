package service

import (
	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
)

type ServiceContainer struct {
	URL urlService
	DB  *repository.Repository
}

func NewServiceContainer(repo map[string]string, config config.AppConfig, storage *Storage, db *repository.Repository) (*ServiceContainer, error) {
	URLService := urlService{
		repo:    repo,
		db:      db,
		config:  config,
		storage: storage,
	}

	return &ServiceContainer{
		URL: URLService,
		DB:  db,
	}, nil
}
