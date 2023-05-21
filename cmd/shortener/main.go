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

	services := service.NewServiceContainer(repo, appConfig)
	handler := handler.NewHandler(services)
	server := new(server.Server)

	if err := server.Run(appConfig.Server.ServerAddr, handler.InitRoutes(sugar)); err != nil {
		panic(err)
	}
}
