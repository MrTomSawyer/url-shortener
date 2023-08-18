// Package handler provides HTTP request handlers for managing URL-related operations.
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

// ShortenURLjson handles the HTTP POST request to shorten a URL using JSON input.
func (h *Handler) ShortenURLjson(c *gin.Context) {
	var req models.ShortenRequest
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

	dec := json.NewDecoder(body)
	if err := dec.Decode(&req); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if req.URL == "" {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	shortenURL, err := h.services.URL.ShortenURLHandler(req.URL, userIDStr)
	if err != nil {
		var urlConflictError *apperrors.URLConflict
		if errors.As(err, &urlConflictError) {
			res := models.ShortenResponse{
				Result: shortenURL,
			}
			c.JSON(http.StatusConflict, res)
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := models.ShortenResponse{
		Result: shortenURL,
	}

	c.JSON(http.StatusCreated, res)
}
