package main

import (
	"context"

	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/handler"
	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
	"github.com/MrTomSawyer/url-shortener/internal/app/server"
	"github.com/MrTomSawyer/url-shortener/internal/app/service"

	_ "github.com/lib/pq"
)

func main() {
	appConfig := config.AppConfig{}
	err := appConfig.InitAppConfig()
	if err != nil {
		panic(err)
	}

	err = logger.InitLogger()
	if err != nil {
		panic(err)
	}

	repo, err := repository.NewRepositoryContainer(context.Background(), appConfig)
	if err != nil {
		panic(err)
	}

	services, err := service.NewServiceContainer(repo, appConfig)
	if err != nil {
		panic(err)
	}
	handler := handler.NewHandler(services, appConfig)
	server := new(server.Server)

	if err := server.Run(appConfig.Server.ServerAddr, handler.InitRoutes()); err != nil {
		panic(err)
	}
}
