package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
	"github.com/MrTomSawyer/url-shortener/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func BenchmarkShortenURLjson(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	cfg := config.AppConfig{}
	cfg.Server.DefaultAddr = "http://localhost:8080"
	cfg.Server.ServerAddr = ":8080"
	cfg.Server.TempFolder = ""
	cfg.DataBase.ConnectionStr = ""

	err := logger.InitLogger()
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	var db *sqlx.DB
	urlRepo, err := repository.InitRepository(context.Background(), cfg, db)
	if err != nil {
		b.Fatalf("Error initializing repo: %v", err)
	}

	repo, err := repository.NewRepositoryContainer(db, urlRepo)
	if err != nil {
		b.Fatalf("Error creating repo container: %v", err)
	}

	serviceContainer, err := service.NewServiceContainer(repo, cfg)
	if err != nil {
		b.Fatalf("Error creating service container: %v", err)
	}

	h := Handler{
		services: serviceContainer,
	}

	b.ResetTimer()

	logger.Log.Infof("Starting ShortenURLjson benchmark")
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		bodyJSON := `{"url": "https://yandex.ru"}`
		c.Request, _ = http.NewRequest("POST", "http://localhost:8080", strings.NewReader(bodyJSON))
		c.Set("user_id", 1)
		b.StartTimer()

		h.ShortenURLjson(c)
	}
}
