package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/MrTomSawyer/url-shortener/internal/app/apperrors"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ShortenURL(c *gin.Context) {
	body := c.Request.Body
	userID, exists := c.Get("user_id")
	if !exists {
		fmt.Println("Failed to get user_id")
	}
	userIDStr, _ := userID.(string)

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			fmt.Printf("Failed to close body: %v", err)
		}
	}(body)

	data, err := io.ReadAll(body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if string(data) == "" {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyStr := string(data)
	shortURL, err := h.services.URL.ShortenURLHandler(bodyStr, userIDStr)
	if err != nil {
		var urlConflictError *apperrors.URLConflict
		if errors.As(err, &urlConflictError) {
			c.String(http.StatusConflict, shortURL)
			return
		}
		fmt.Printf("Failed to shorten a url: %v\n", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}

	c.String(http.StatusCreated, shortURL)
}
