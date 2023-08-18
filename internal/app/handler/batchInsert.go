// Package handler provides HTTP request handlers for managing URL-related operations.
package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/gin-gonic/gin"
)

// batchURLinsert handles the HTTP POST request to insert URLs in batch.
func (h Handler) batchURLinsert(c *gin.Context) {
	body := c.Request.Body
	userID, exists := c.Get("user_id")
	if !exists {
		fmt.Println("Failed to get user_id")
	}
	userIDStr, _ := userID.(string)

	logger.Log.Infof("batchURLinsert user ID: %s", userIDStr)

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			fmt.Printf("Failed to close body: %v", err)
		}
	}(body)

	res, err := h.services.URL.HandleBatchInsert(body, userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	logger.Log.Infof("batchURLinsert response: %v", res)
	c.JSON(http.StatusCreated, res)
}
