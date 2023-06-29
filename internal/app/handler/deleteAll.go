package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handler) deleteAll(c *gin.Context) {
	body := c.Request.Body

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			fmt.Printf("error closing bosy %v", err)
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
		fmt.Printf("failed to delete url: %v", err)
	}

	c.Writer.WriteHeader(http.StatusAccepted)
}
