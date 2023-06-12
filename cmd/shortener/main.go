package main

import (
	"context"

	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/handler"
	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
	"github.com/MrTomSawyer/url-shortener/internal/app/server"
	"github.com/MrTomSawyer/url-shortener/internal/app/service"
	"github.com/jmoiron/sqlx"

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

	var db *sqlx.DB
	if appConfig.DataBase.ConnectionStr != "" {
		db, err = repository.NewPostgresDB(appConfig.DataBase.ConnectionStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()
	}
	repo, err := repository.NewRepositoryContainer(context.Background(), appConfig, db)
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
