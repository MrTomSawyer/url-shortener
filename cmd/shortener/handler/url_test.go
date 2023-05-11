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

var testVault = make(map[string]string)
var appConfig config.AppConfig

func init() {
	config.InitAppConfig(&appConfig)
}

type want struct {
	code     int
	response string
}

func TestShortenURL(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		method string
		body   string
		want   want
	}{
		{
			name:   "Test #1 - Regular URL",
			path:   "/",
			method: "POST",
			body:   "https://yandex.ru",
			want: want{
				code:     201,
				response: "http://localhost:8080/e9db20b2",
			},
		},
		{
			name:   "Test #2 - Empty Body",
			path:   "/",
			method: "POST",
			body:   "",
			want: want{
				code:     400,
				response: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			с, _ := gin.CreateTestContext(w)

			с.Request, _ = http.NewRequest(test.method, test.path, strings.NewReader(test.body))
			h := Handler{
				services: service.NewServiceContainer(testVault, appConfig),
			}
			h.ShortenURL(с)

			if с.Writer.Status() != test.want.code {
				t.Errorf("got status code %d, want %d", w.Code, test.want.code)
			}

			if body := strings.TrimSpace(w.Body.String()); body != test.want.response {
				t.Errorf("got response body '%s', want '%s'", body, test.want.response)
			}
		})
	}
}

func TestGetOriginalURL(t *testing.T) {
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
			h.GetOriginalURL(c)
			if c.Writer.Status() != test.want.code {
				t.Errorf("got status code %d, want %d", w.Code, test.want.code)
			}
			if location := w.Header().Get("Location"); location != test.want.response {
				t.Errorf("got location header %s, want %s", location, test.want.response)
			}
		})
	}
}
