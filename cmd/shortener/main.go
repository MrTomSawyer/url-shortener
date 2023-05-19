package main

import (
	"flag"

	"github.com/MrTomSawyer/url-shortener/cmd/shortener/config"
	"github.com/MrTomSawyer/url-shortener/cmd/shortener/handler"
	"github.com/MrTomSawyer/url-shortener/cmd/shortener/server"
	"github.com/MrTomSawyer/url-shortener/cmd/shortener/service"
)

func main() {
	flag.Parse()

	appConfig := config.AppConfig{}
	appConfig.InitAppConfig()
	repo := make(map[string]string)

	services := service.NewServiceContainer(repo, appConfig)
	handler := handler.NewHandler(services)
	server := new(server.Server)

	if err := server.Run(appConfig.Server.ServerAddr, handler.InitRoutes()); err != nil {
		panic(err)
	}
}
