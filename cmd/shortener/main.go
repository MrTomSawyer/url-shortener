package main

import (
	"flag"

	"github.com/MrTomSawyer/url-shortener/cmd/shortener/config"
	"github.com/MrTomSawyer/url-shortener/cmd/shortener/handler"
	"github.com/MrTomSawyer/url-shortener/cmd/shortener/server"
	"github.com/MrTomSawyer/url-shortener/cmd/shortener/service"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.AppConfig{}
	appConfig.InitAppConfig()
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	repo := make(map[string]string)
	storage := service.NewStorage(appConfig.Server.TempFolder)

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
