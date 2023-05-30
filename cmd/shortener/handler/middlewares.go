package handler

import (
	"net/http"
	"strings"
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

func (h *Handler) DataCompressor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Принимаем запросы в сжатом формате
		encoding := c.GetHeader("Content-Encoding")
		if strings.Contains(encoding, "gzip") {
			compressReader, err := newCompressReader(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			defer compressReader.Close()
			c.Request.Body = compressReader
		}

		acceptEncoding := c.GetHeader("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			switch c.Writer.Header().Get("Content-Type") {
			case "application/json", "text/html":
				compressWriter := newCompressWriter(c.Writer)
				compressWriter.writer.Header().Set("Content-Type", "gzip")
				defer compressWriter.Close()
			}
			// Продолжаем обработку запроса
			c.Next()
		} else {
			// Клиент не поддерживает сжатие, просто продолжаем обработку запроса
			c.Next()
		}
	}
}
