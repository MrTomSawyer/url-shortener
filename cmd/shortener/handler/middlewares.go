package handler

import (
	"compress/gzip"
	"fmt"
	"net/http"
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

func (h *Handler) decompressData() gin.HandlerFunc {
	return func(c *gin.Context) {
		decompressTypes := [2]string{"application/json", "text/plain"}
		if c.Request.Header.Get("Content-Encoding") != "gzip" {
			c.Next()
			return
		}

		isCorrectType := false
		contentType := c.Request.Header.Get("Content-Type")
		for _, val := range decompressTypes {
			if val == contentType {
				isCorrectType = true
				break
			}
		}
		if !isCorrectType {
			c.Next()
			return
		}

		gzipReader, err := gzip.NewReader(c.Request.Body)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		defer func() {
			err := gzipReader.Close()
			if err != nil {
				fmt.Println("Failed to close gzip reader", err)
			}
		}()

		c.Request.Body = http.MaxBytesReader(c.Writer, gzipReader, c.Request.ContentLength)
		c.Next()
	}
}
