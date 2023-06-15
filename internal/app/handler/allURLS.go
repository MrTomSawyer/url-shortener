package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetAll(c *gin.Context) {
	response, err := h.services.URL.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
