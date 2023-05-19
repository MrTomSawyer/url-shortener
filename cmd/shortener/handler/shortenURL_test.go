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

func init() {

}

func TestShortenURL(t *testing.T) {
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
