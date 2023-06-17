package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/MrTomSawyer/url-shortener/internal/app/apperrors"
	"github.com/MrTomSawyer/url-shortener/internal/app/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ShortenURLjson(c *gin.Context) {
	var req models.ShortenRequest
	body := c.Request.Body

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			fmt.Printf("Failed to close body: %v", err)
		}
	}(body)

	dec := json.NewDecoder(body)
	if err := dec.Decode(&req); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}

	if req.URL == "" {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	shortenURL, err := h.services.URL.ShortenURLHandler(req.URL)
	if err != nil {
		var urlConflictError *apperrors.URLConflict
		if errors.As(err, &urlConflictError) {
			res := models.ShortenResponce{
				Result: shortenURL,
			}
			c.JSON(http.StatusConflict, res)
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}

	res := models.ShortenResponce{
		Result: shortenURL,
	}

	c.JSON(http.StatusCreated, res)
}
