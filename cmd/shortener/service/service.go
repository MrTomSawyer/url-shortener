package service

import "github.com/MrTomSawyer/url-shortener/cmd/shortener/config"

type ServiceContainer struct {
	URL urlService
}

func NewServiceContainer(repo map[string]string, config config.AppConfig, storage *Storage) *ServiceContainer {
	return &ServiceContainer{
		URL: urlService{
			repo:    repo,
			config:  config,
			storage: storage,
		},
	}
}
