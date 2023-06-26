package handler

import (
	"bytes"
	"fmt"
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
	"github.com/stretchr/testify/assert"
)

func TestdeleteAll(t *testing.T) {
	cfg := config.AppConfig{}
	cfg.Server.DefaultAddr = "http://localhost:8080"
	cfg.Server.ServerAddr = ":8080"
	cfg.Server.TempFolder = ""
	cfg.DataBase.ConnectionStr = "string"

	err := logger.InitLogger()
	if err != nil {
		panic(err)
	}

	type want struct {
		code int
	}

	tests := []struct {
		name   string
		url    string
		body   []byte
		method string
		want   want
	}{
		{
			name:   "Test #1 - Batch delete",
			url:    "localhost:8080/api/user/urls",
			body:   []byte(`["e98192e1"]`),
			method: "DELETE",
			want: want{
				code: 202,
			},
		},
	}

	for _, test := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepoHandler(ctrl)

		m.EXPECT().DeleteAll(gomock.Any(), gomock.Any()).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := bytes.NewBuffer(test.body)
		c.Request, _ = http.NewRequest(test.method, test.url, body)

		var db *sqlx.DB
		repo, err := repository.NewRepositoryContainer(db, m)
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
		h.batchURLinsert(c)

		assert.Equal(t, test.want.code, c.Writer.Status())
	}
}
