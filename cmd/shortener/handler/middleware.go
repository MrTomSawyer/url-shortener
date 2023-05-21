package handler

import (
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) logReqResInfo(lg logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		latency := time.Since(startTime)

		lg.Infof("Method: %s URI: %s Latency: %v", c.Request.Method, c.Request.URL, latency)
	}
}
