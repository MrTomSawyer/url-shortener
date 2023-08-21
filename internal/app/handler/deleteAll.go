// Package handler provides HTTP request handlers for managing URL-related operations.
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// deleteAll handles the HTTP DELETE request to delete multiple URLs.
func (h Handler) deleteAll(c *gin.Context) {
	body := c.Request.Body

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			fmt.Printf("Error closing body: %v", err)
		}
	}(body)

	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(fmt.Errorf("failed to get user_id"))
		return
	}
	userIDStr, _ := userID.(string)

	bytes, err := io.ReadAll(body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error reading body"))
		return
	}

	var data []string
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error unmarshalling body"))
		return
	}

	err = h.services.URL.DeleteAll(data, userIDStr)
	if err != nil {
		fmt.Printf("Failed to delete URL: %v", err)
	}

	c.Writer.WriteHeader(http.StatusAccepted)
}
