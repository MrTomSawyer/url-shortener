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
	tx, err := h.services.DB.Begin()
	if err != nil {
		logger.Log.Infof("Failed to start transaction")
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	body := c.Request.Body

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			fmt.Printf("Failed to close body: %v", err)
		}
	}(body)

	var parsedReq []models.BatchURLRequest

	decoder := json.NewDecoder(body)
	err = decoder.Decode(&parsedReq)
	if err != nil {
		logger.Log.Infof("Failed to decode json")
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	var response []models.BatchURLResponce

	for _, req := range parsedReq {
		shortURL, err := h.services.URL.ShortenURL(req.OriginalURL)
		if err != nil {
			logger.Log.Infof("Failed to shorten URL")
			continue
		}
		query := "INSERT INTO urls (correlationid, shorturl, originalurl) VALUES ($1, $2, $3)"
		_, err = tx.Exec(query, req.CorrelationID, shortURL, req.OriginalURL)

		if err != nil {
			logger.Log.Infof("Failed to insert a shortened URL")
			continue
		}

		response = append(response, models.BatchURLResponce{
			CorrelationID: req.CorrelationID,
			ShortURL:      shortURL,
		})
	}

	err = tx.Commit()
	if err != nil {
		logger.Log.Infof("Failed to commit a transaction")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}
