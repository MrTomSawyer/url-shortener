package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MrTomSawyer/url-shortener/cmd/shortener/config"
	"github.com/MrTomSawyer/url-shortener/cmd/shortener/service"
	"github.com/gin-gonic/gin"
)

func TestExpandURL(t *testing.T) {
	appConfig := config.AppConfig{}
	appConfig.InitAppConfig()
	var testVault = make(map[string]string)

	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name   string
		path   string
		method string
		want   want
	}{
		{
			name:   "Test #3 - Get Original URL",
			path:   "/e9db20b2",
			method: "GET",
			want: want{
				code:     307,
				response: "https://yandex.ru",
			},
		},
		{
			name:   "Test #4 - Wrong code",
			path:   "/fff",
			method: "GET",
			want: want{
				code:     404,
				response: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest(test.method, test.path, strings.NewReader(""))
			h := Handler{
				services: service.NewServiceContainer(testVault, appConfig),
			}
			h.ExpandURL(c)
			if c.Writer.Status() != test.want.code {
				t.Errorf("got status code %d, want %d", w.Code, test.want.code)
			}
			if location := w.Header().Get("Location"); location != test.want.response {
				t.Errorf("got location header %s, want %s", location, test.want.response)
			}
		})
	}
}
