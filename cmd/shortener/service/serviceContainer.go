package service

import "github.com/MrTomSawyer/url-shortener/cmd/shortener/config"

type ServiceContainer struct {
	URL URLservice
}

func NewServiceContainer(repo map[string]string, config config.AppConfig) *ServiceContainer {
	return &ServiceContainer{
		URL: URLservice{
			repo:   repo,
			config: config,
		},
	}
}
