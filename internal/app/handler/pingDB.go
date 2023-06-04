package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) pingDB(c *gin.Context) {
	err := h.services.DB.Postgres.Ping()
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
