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
	// Initialize application configuration.
	appConfig := config.AppConfig{}
	err := appConfig.InitAppConfig()
	if err != nil {
		panic(err)
	}

	// Initialize logger.
	err = logger.InitLogger()
	if err != nil {
		panic(err)
	}

	var db *sqlx.DB
	if appConfig.DataBase.ConnectionStr != "" {
		// Create and set up a database connection.
		db, err = repository.NewPostgresDB(appConfig.DataBase.ConnectionStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()
	}

	// Initialize URL repository.
	urlRepo, err := repository.InitRepository(context.Background(), appConfig, db)
	if err != nil {
		panic(err)
	}

	// Create and set up the repository container.
	repo, err := repository.NewRepositoryContainer(db, urlRepo)
	if err != nil {
		panic(err)
	}

	// Create and set up service container.
	services, err := service.NewServiceContainer(repo, appConfig)
	if err != nil {
		panic(err)
	}

	// Create and set up request handlers.
	h := handler.NewHandler(services, appConfig)

	// Create a new HTTP server instance.
	s := new(server.Server)

	// Start the HTTP server with configured routes.
	if err := s.Run(appConfig.Server.ServerAddr, h.InitRoutes()); err != nil {
		panic(err)
	}
}
