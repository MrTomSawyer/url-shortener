package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
	"github.com/MrTomSawyer/url-shortener/internal/app/repository/mocks"
	"github.com/MrTomSawyer/url-shortener/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
)

func BenchmarkExpandURL(b *testing.B) {
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

	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	m := mocks.NewMockRepoHandler(ctrl)

	var db *sqlx.DB
	repo, err := repository.NewRepositoryContainer(db, m)
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

	logger.Log.Infof("Starting ExpandURL benchmark")
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "http://localhost:8080", strings.NewReader(""))
		c.Set("user_id", 1)
		c.AddParam("id", "1")
		m.EXPECT().OriginalURL(gomock.Any()).Return("https://yandex.ru", nil)
		b.StartTimer()

		h.ExpandURL(c)
	}
}
