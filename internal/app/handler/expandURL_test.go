package handler

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/MrTomSawyer/url-shortener/internal/app/config"
// 	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
// 	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
// 	"github.com/MrTomSawyer/url-shortener/internal/app/repository/mocks"
// 	"github.com/MrTomSawyer/url-shortener/internal/app/service"
// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/mock/gomock"
// )

// func TestExpandURL(t *testing.T) {
// 	cfg := config.AppConfig{}
// 	cfg.Server.DefaultAddr = "http://localhost:8080"
// 	cfg.Server.ServerAddr = ":8080"
// 	cfg.Server.TempFolder = ""
// 	cfg.DataBase.ConnectionStr = ""

// 	err := logger.InitLogger()
// 	if err != nil {
// 		panic(err)
// 	}

// 	type want struct {
// 		code     int
// 		response string
// 	}

// 	tests := []struct {
// 		name   string
// 		url    string
// 		id     string
// 		method string
// 		want   want
// 	}{
// 		{
// 			name:   "Test #3 - Get Original URL",
// 			url:    "http://localhost:8080",
// 			id:     "e9db20b2",
// 			method: "GET",
// 			want: want{
// 				code:     307,
// 				response: "https://yandex.ru",
// 			},
// 		},
// 		{
// 			name:   "Test #4 - Wrong code",
// 			url:    "http://localhost:8080",
// 			id:     "fff",
// 			method: "GET",
// 			want: want{
// 				code:     404,
// 				response: "",
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			m := mocks.NewMockRepoHandler(ctrl)

// 			w := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(w)

// 			c.Request, _ = http.NewRequest(test.method, test.url, strings.NewReader(""))
// 			c.AddParam("id", test.id)
// 			m.EXPECT().OriginalURL("e9db20b2").Return("https://yandex.ru", nil)
// 			m.EXPECT().OriginalURL("fff").Return("", nil)

// 			repo, err := repository.NewRepositoryContainer(context.Background(), cfg)
// 			if err != nil {
// 				fmt.Printf("Error creating repo container: %v", err)
// 			}
// 			serviceContainer, err := service.NewServiceContainer(repo, cfg)
// 			if err != nil {
// 				fmt.Printf("Error creating service container: %v", err)
// 			}
// 			h := Handler{
// 				services: serviceContainer,
// 			}
// 			h.ExpandURL(c)
// 			if c.Writer.Status() != test.want.code {
// 				t.Errorf("got status code %d, want %d", w.Code, test.want.code)
// 			}
// 			if location := w.Header().Get("Location"); location != test.want.response {
// 				t.Errorf("got location header %s, want %s", location, test.want.response)
// 			}
// 		})
// 	}
// }
