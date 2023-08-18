// Package handler provides HTTP request handlers for managing URL-related operations.
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// pingDB handles the HTTP GET request to check the database connection.
func (h *Handler) pingDB(c *gin.Context) {
	err := h.services.DB.Ping()
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
