package service

import (
	"fmt"

	"github.com/MrTomSawyer/url-shortener/cmd/shortener/config"
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
	err := URLService.initializeLastUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to init UUID: %v", err)
	}

	return &ServiceContainer{
		URL: URLService,
	}, nil
}
