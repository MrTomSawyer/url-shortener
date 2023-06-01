package main

import (
	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/handler"
	"github.com/MrTomSawyer/url-shortener/internal/app/server"
	"github.com/MrTomSawyer/url-shortener/internal/app/service"
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

	repo := make(map[string]string)

	storage, err := service.NewStorage(appConfig.Server.TempFolder)
	if err != nil {
		panic(err)
	}

	err = storage.Read(&repo)
	if err != nil {
		panic(err)
	}

	services, err := service.NewServiceContainer(repo, appConfig, storage)
	if err != nil {
		panic(err)
	}
	handler := handler.NewHandler(services)
	server := new(server.Server)

	if err := server.Run(appConfig.Server.ServerAddr, handler.InitRoutes(sugar)); err != nil {
		panic(err)
	}
}
