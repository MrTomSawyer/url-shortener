package main

import (
	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/handler"
	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
	"github.com/MrTomSawyer/url-shortener/internal/app/server"
	"github.com/MrTomSawyer/url-shortener/internal/app/service"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.AppConfig{}
	err := appConfig.InitAppConfig()
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	repo, err := repository.NewRepositoryContainer(appConfig)
	if err != nil {
		panic(err)
	}
	
	services, err := service.NewServiceContainer(repo, appConfig)
	if err != nil {
		panic(err)
	}
	handler := handler.NewHandler(services)
	server := new(server.Server)

	if err := server.Run(appConfig.Server.ServerAddr, handler.InitRoutes(sugar)); err != nil {
		panic(err)
	}
}
