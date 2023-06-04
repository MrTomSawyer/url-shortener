package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) pingDB(c *gin.Context) {
	err := h.services.DB.DB.Ping()
	if err != nil {
		fmt.Printf("ERR DB %v\n", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
