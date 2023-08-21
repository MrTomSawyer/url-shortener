package service

import (
	"context"
	"fmt"
	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Example_shortenURL() {
	gin.SetMode(gin.ReleaseMode)

	cfg := config.AppConfig{}
	cfg.Server.DefaultAddr = "http://localhost:8080"
	cfg.Server.ServerAddr = ":8080"
	cfg.Server.TempFolder = ""
	cfg.DataBase.ConnectionStr = ""

	err := logger.InitLogger()
	if err != nil {
		panic(err)
	}

	var db *sqlx.DB
	urlRepo, err := repository.InitRepository(context.Background(), cfg, db)
	if err != nil {
		fmt.Printf("Error creating initializing repo: %v", err)
	}

	repo, err := repository.NewRepositoryContainer(db, urlRepo)
	if err != nil {
		fmt.Printf("Error creating repo container: %v", err)
	}
	serviceContainer, err := NewServiceContainer(repo, cfg)
	if err != nil {
		fmt.Printf("Error creating service container: %v", err)
	}

	_, err = serviceContainer.URL.ShortenURL("https://yandex.ru")
	fmt.Printf("Error: %v\n", err)

	// Output:
	// Error: <nil>
}
