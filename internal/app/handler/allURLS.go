package handler

import (
	"net/http"

	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/gin-gonic/gin"
)

func (h Handler) GetAll(c *gin.Context) {
	userid, exists := c.Get("user_id")
	if !exists {
		logger.Log.Infof("GetAll: no userid found")
	}
	useridStr, ok := userid.(string)
	if !ok {
		logger.Log.Infof("FAiled to convert userid to string")
	}

	logger.Log.Infof("GetAll user id: %s", useridStr)

	response, err := h.services.URL.GetAll(useridStr)
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
