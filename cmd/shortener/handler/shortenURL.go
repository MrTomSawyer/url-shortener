package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ShortenURL(c *gin.Context) {
	body := c.Request.Body

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			fmt.Printf("Failed to close body: %v", err)
		}
	}(body)

	data, err := io.ReadAll(body)
	fmt.Printf("EEEEEEERRRRRRR: %v %v\n", string(data) == "", len(data) == 0)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if string(data) == "" {
		fmt.Printf("fail 212121212121")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyStr := string(data)
	shortURL, err := h.services.URL.ShortenURL(bodyStr)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}

	c.String(http.StatusCreated, shortURL)
}
