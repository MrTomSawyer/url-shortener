package service

import (
	"github.com/MrTomSawyer/url-shortener/internal/app/config"
)

type ServiceContainer struct {
	URL urlService
}

func NewServiceContainer(repo map[string]string, config config.AppConfig, storage *Storage) (*ServiceContainer, error) {
	URLService := urlService{
		repo:    repo,
		config:  config,
		storage: storage,
	}

	return &ServiceContainer{
		URL: URLService,
	}, nil
}
