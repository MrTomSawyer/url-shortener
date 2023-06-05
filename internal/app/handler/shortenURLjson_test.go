package handler

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/MrTomSawyer/url-shortener/internal/app/config"
// 	"github.com/MrTomSawyer/url-shortener/internal/app/service"
// 	"github.com/gin-gonic/gin"
// 	"github.com/goccy/go-json"
// )

// func TestShortenURLjson(t *testing.T) {
// 	cfg := config.AppConfig{}
// 	cfg.Server.DefaultAddr = "http://localhost:8080"
// 	cfg.Server.ServerAddr = ":8080"
// 	cfg.Server.TempFolder = "/tmp/short-url-db.json"

// 	testVault := make(map[string]string)
// 	storage, err := service.NewStorage(cfg.Server.TempFolder)
// 	if err != nil {
// 		fmt.Printf("Failed to create test storage: %v", err)
// 		return
// 	}

// 	type want struct {
// 		code     int
// 		response map[string]string
// 	}

// 	tests := []struct {
// 		name   string
// 		url    string
// 		method string
// 		body   map[string]string
// 		want   want
// 	}{
// 		{
// 			name:   "Test #1 - Regular URL",
// 			url:    "http://localhost:8080",
// 			method: "POST",
// 			body: map[string]string{
// 				"url": "https://yandex.ru",
// 			},
// 			want: want{
// 				code: 201,
// 				response: map[string]string{
// 					"result": "http://localhost:8080/e9db20b2",
// 				},
// 			},
// 		},
// 		{
// 			name:   "Test #2 - Empty Body",
// 			url:    "http://localhost:8080",
// 			method: "POST",
// 			body:   map[string]string{},
// 			want: want{
// 				code:     400,
// 				response: map[string]string{},
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(w)

// 			bodyStr, err := json.Marshal(test.body)
// 			if err != nil {
// 				fmt.Println("Error marshalling body:", err)
// 			}

// 			c.Request, _ = http.NewRequest(test.method, test.url, strings.NewReader(string(bodyStr)))
// 			serviceContainer, err := service.NewServiceContainer(testVault, cfg, storage, nil)
// 			if err != nil {
// 				fmt.Printf("Error creating service container: %v", err)
// 			}
// 			h := Handler{
// 				services: serviceContainer,
// 			}
// 			h.ShortenURLjson(c)

// 			if c.Writer.Status() != test.want.code {
// 				t.Errorf("got status code %d, want %d", w.Code, test.want.code)
// 			}
// 		})
// 	}
// }
