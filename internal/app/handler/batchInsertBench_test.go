package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
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

func BenchmarkBatchURLInsert(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	cfg := config.AppConfig{}
	cfg.Server.DefaultAddr = "http://localhost:8080"
	cfg.Server.ServerAddr = ":8080"
	cfg.Server.TempFolder = ""
	cfg.DataBase.ConnectionStr = "string"

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

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m.EXPECT().OriginalURL(gomock.Any()).Return("", nil)
		m.EXPECT().BatchCreate(gomock.Any(), gomock.Any()).Return(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := bytes.NewBuffer([]byte(`[
			{
				"correlation_id": "1",
				"original_url":   "https://ya.ru"
			}
		]`))
		c.Request, _ = http.NewRequest("POST", "http://localhost:8080/api/shorten/batch", body)

		b.StartTimer()

		h.batchURLinsert(c)
	}
}
