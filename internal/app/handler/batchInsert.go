package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/MrTomSawyer/url-shortener/internal/app/models"
	"github.com/gin-gonic/gin"
)

func (h Handler) batchURLinsert(c *gin.Context) {
	body := c.Request.Body

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			fmt.Printf("Failed to close body: %v", err)
		}
	}(body)

	var parsedReq []models.BatchURLRequest

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&parsedReq)
	if err != nil {
		logger.Log.Infof("Failed to decode json")
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	var response []models.BatchURLResponce
	var shortURL string

	type tempURLRequest struct {
		correlationID string
		shortURL      string
		originalURL   string
	}
	var tempURLRequests []tempURLRequest

	for _, req := range parsedReq {
		shortURL, err = h.services.URL.ShortenURL(req.OriginalURL)
		if err != nil {
			logger.Log.Infof("Failed to shorten URL")
			continue
		}
		tempURLRequests = append(tempURLRequests, tempURLRequest{req.CorrelationID, shortURL, req.OriginalURL})
	}

	if h.Cfg.DataBase.ConnectionStr != "" {
		tx, err := h.services.DB.Begin()
		if err != nil {
			logger.Log.Infof("Failed to start transaction")
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		for _, req := range tempURLRequests {
			query := "INSERT INTO urls (correlationid, shorturl, originalurl) VALUES ($1, $2, $3)"
			_, err = tx.Exec(query, req.correlationID, req.shortURL, req.originalURL)

			if err != nil {
				logger.Log.Infof("Failed to insert a shortened URL", err)
				tx.Rollback()
				continue
			}

			response = append(response, models.BatchURLResponce{
				CorrelationID: req.correlationID,
				ShortURL:      fmt.Sprintf("%s/%s", h.Cfg.Server.DefaultAddr, shortURL),
			})
		}
		err = tx.Commit()
		if err != nil {
			logger.Log.Infof("Failed to commit a transaction")
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else {
		for _, req := range tempURLRequests {
			err := h.services.URL.Repo.Create(req.shortURL, req.originalURL)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}

	c.JSON(http.StatusCreated, response)
}
