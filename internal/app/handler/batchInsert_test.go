package handler

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/MrTomSawyer/url-shortener/internal/app/config"
// 	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
// 	"github.com/MrTomSawyer/url-shortener/internal/app/models"
// 	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
// 	"github.com/MrTomSawyer/url-shortener/internal/app/repository/mocks"
// 	"github.com/MrTomSawyer/url-shortener/internal/app/service"
// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/mock/gomock"
// 	"github.com/jmoiron/sqlx"
// 	"github.com/stretchr/testify/assert"
// )

// func TestBatchURLinsert(t *testing.T) {
// 	cfg := config.AppConfig{}
// 	cfg.Server.DefaultAddr = "http://localhost:8080"
// 	cfg.Server.ServerAddr = ":8080"
// 	cfg.Server.TempFolder = ""
// 	cfg.DataBase.ConnectionStr = "string"

// 	err := logger.InitLogger()
// 	if err != nil {
// 		panic(err)
// 	}

// 	type want struct {
// 		code     int
// 		response []models.BatchURLResponce
// 	}

// 	tests := []struct {
// 		name      string
// 		url       string
// 		body      []byte
// 		method    string
// 		shortPath string
// 		want      want
// 	}{
// 		{
// 			name: "Test #1 - Batch insert",
// 			url:  "localhost:8080/api/shorten/batch",
// 			body: []byte(`[
// 				{
// 					"correlation_id": "1",
// 					"original_url":   "https://ya.ru"
// 				}
// 			]`),
// 			shortPath: "e98192e1",
// 			method:    "POST",
// 			want: want{
// 				code: 201,
// 				response: []models.BatchURLResponce{
// 					{
// 						CorrelationID: "1",
// 						ShortURL:      "http://localhost:8080/e98192e1",
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		m := mocks.NewMockRepoHandler(ctrl)

// 		m.EXPECT().OriginalURL(test.shortPath).Return("", nil)
// 		m.EXPECT().BatchCreate(gomock.Any()).Return(test.want.response, nil)

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		body := bytes.NewBuffer(test.body)
// 		c.Request, _ = http.NewRequest(test.method, test.url, body)

// 		var db *sqlx.DB
// 		repo, err := repository.NewRepositoryContainer(db, m)
// 		if err != nil {
// 			fmt.Printf("Error creating repo container: %v", err)
// 		}

// 		serviceContainer, err := service.NewServiceContainer(repo, cfg)
// 		if err != nil {
// 			fmt.Printf("Error creating service container: %v", err)
// 		}

// 		h := Handler{
// 			services: serviceContainer,
// 		}
// 		h.batchURLinsert(c)

// 		assert.Equal(t, test.want.code, c.Writer.Status())

// 		wantBody, _ := json.Marshal(test.want.response)
// 		assert.JSONEq(t, string(wantBody), strings.TrimSpace(w.Body.String()), "unexpected response body")
// 	}
// }
