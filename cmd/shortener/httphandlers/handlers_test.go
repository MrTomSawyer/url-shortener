package httphandlers

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testVault = map[string]string{
	"e9db20b2": "https://yandex.ru",
}

type want struct {
	code    int
	resBody string
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
				code:    201,
				resBody: "http://localhost:8080/e9db20b2",
			},
		},
		{
			name:   "Test #2 - Empry Body",
			path:   "/",
			method: "GET",
			body:   "",
			want: want{
				code:    400,
				resBody: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reqBody := bytes.NewReader([]byte(test.body))
			r := httptest.NewRequest(test.method, test.path, reqBody)
			w := httptest.NewRecorder()

			ShortenURL(w, r)

			assert.Equal(t, test.want.code, w.Code, "Wrong response code")

			if test.want.resBody != "" {
				assert.Equal(t, test.want.resBody, w.Body.String(), "Wrong response body")
			}
		})
	}
}

func TestGetOriginalURL(t *testing.T) {
	vault = testVault

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
				code:    307,
				resBody: "https://yandex.ru",
			},
		},
		{
			name:   "Test #4 - Wrong code",
			path:   "/fff",
			method: "GET",
			want: want{
				code:    404,
				resBody: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()

			GetOriginalURL(w, r)

			fmt.Println(vault)

			assert.Equal(t, test.want.code, w.Code, "Wrong response code")

			if test.want.resBody != "" {
				assert.Equal(t, test.want.resBody, w.Header().Get("Location"), "Wrong original URL")
			}

		})
	}

}
