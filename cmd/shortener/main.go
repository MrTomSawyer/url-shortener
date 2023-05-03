package main

import (
	"github.com/MrTomSawyer/url-shortener/cmd/shortener/config"
	"github.com/MrTomSawyer/url-shortener/cmd/shortener/handler"
	"github.com/MrTomSawyer/url-shortener/cmd/shortener/server"
	"github.com/MrTomSawyer/url-shortener/cmd/shortener/service"
)

var repo = make(map[string]string)

func main() {
	config.ParseFlags()

	services := service.NewService(repo)
	handler := handler.NewHandler(services)
	server := new(server.Server)

	if err := server.Run(config.ServerAddr, handler.InitRoutes()); err != nil {
		panic(err)
	}
}
