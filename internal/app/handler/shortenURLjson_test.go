package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
	"github.com/MrTomSawyer/url-shortener/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/jmoiron/sqlx"
)

func TestShortenURLjson(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	cfg := config.AppConfig{}
	cfg.Server.DefaultAddr = "http://localhost:8080"
	cfg.Server.ServerAddr = ":8080"
	cfg.Server.TempFolder = "/tmp/short-url-db.json"

	err := logger.InitLogger()
	if err != nil {
		panic(err)
	}

	type want struct {
		code     int
		response map[string]string
	}

	tests := []struct {
		name   string
		url    string
		method string
		body   map[string]string
		want   want
	}{
		{
			name:   "Test #1 - Regular URL",
			url:    "http://localhost:8080",
			method: "POST",
			body: map[string]string{
				"url": "https://yandex.ru",
			},
			want: want{
				code: 201,
				response: map[string]string{
					"result": "http://localhost:8080/e9db20b2",
				},
			},
		},
		{
			name:   "Test #2 - Empty Body",
			url:    "http://localhost:8080",
			method: "POST",
			body:   map[string]string{},
			want: want{
				code:     400,
				response: map[string]string{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyStr, err := json.Marshal(test.body)
			if err != nil {
				fmt.Println("Error marshalling body:", err)
			}

			c.Request, _ = http.NewRequest(test.method, test.url, strings.NewReader(string(bodyStr)))

			var db *sqlx.DB
			urlRepo, err := repository.InitRepository(context.Background(), cfg, db)
			if err != nil {
				fmt.Printf("Error creating initializing repo: %v", err)
			}
			repo, err := repository.NewRepositoryContainer(db, urlRepo)
			if err != nil {
				fmt.Printf("Error creating repo container: %v", err)
			}

			serviceContainer, err := service.NewServiceContainer(repo, cfg)
			if err != nil {
				fmt.Printf("Error creating service container: %v", err)
			}
			h := Handler{
				services: serviceContainer,
			}
			h.ShortenURLjson(c)

			if c.Writer.Status() != test.want.code {
				t.Errorf("got status code %d, want %d", w.Code, test.want.code)
			}
		})
	}
}
