// Package handler provides HTTP request handlers for managing URL-related operations.
package handler

import (
	"net/http"

	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/gin-gonic/gin"
)

// GetAll handles the HTTP GET request to retrieve all URLs associated with a user.
func (h Handler) GetAll(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		logger.Log.Infof("GetAll: no user ID found")
	}
	userIDStr, ok := userID.(string)
	if !ok {
		logger.Log.Infof("Failed to convert user ID to string")
	}

	logger.Log.Infof("GetAll user ID: %s", userIDStr)

	response, err := h.services.URL.GetAll(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if len(response) == 0 {
		c.JSON(http.StatusNoContent, nil)
		return
	}
	c.JSON(http.StatusOK, response)
}
