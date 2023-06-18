package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/gin-gonic/gin"
)

func (h Handler) batchURLinsert(c *gin.Context) {
	body := c.Request.Body
	userID, exists := c.Get("user_id")
	if !exists {
		fmt.Println("Failed to get user_id")
	}
	userIDStr, _ := userID.(string)

	logger.Log.Infof("batchURLinsert user id: %s", userID)

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			fmt.Printf("Failed to close body: %v", err)
		}
	}(body)

	res, err := h.services.URL.HandleBatchInsert(body, userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	logger.Log.Infof("batchURLinsert response: ", res)
	c.JSON(http.StatusCreated, res)
}
