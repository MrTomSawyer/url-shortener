package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handler) batchURLinsert(c *gin.Context) {
	body := c.Request.Body
	userId, exists := c.Get("user_id")
	if !exists {
		fmt.Println("Failed to get user_id")
	}
	userIdStr, _ := userId.(string)
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			fmt.Printf("Failed to close body: %v", err)
		}
	}(body)

	res, err := h.services.URL.HandleBatchInsert(body, userIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusCreated, res)
}
