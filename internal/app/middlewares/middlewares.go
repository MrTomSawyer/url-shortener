// Package middlewares provides middleware functions for handling various aspects of HTTP requests and responses.
package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/gin-gonic/gin"
)

// LogReqResInfo is a middleware that logs information about incoming requests and outgoing responses.
func LogReqResInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		latency := time.Since(startTime)

		logger.Log.Infof("Method: %s URI: %s Latency: %v", c.Request.Method, c.Request.URL, latency)
	}
}

// DataCompressor is a middleware that handles request and response data compression.
func DataCompressor() gin.HandlerFunc {
	return func(c *gin.Context) {
		encoding := c.Request.Header.Get("Content-Encoding")

		if strings.Contains(encoding, "gzip") {
			gzipBodyReader, err := newGzipBodyReader(c.Request.Body)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			defer func() {
				if err := gzipBodyReader.Close(); err != nil {
					fmt.Printf("Failed to close gzip body reader: %v", err)
				}
			}()
			c.Request.Body = gzipBodyReader
		}

		acceptEncoding := c.GetHeader("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			switch c.Writer.Header().Get("Content-Type") {
			case "application/json", "text/html":
				compressWriter := newGzipBodyWriter(c.Writer)
				compressWriter.writer.Header().Set("Content-Type", "gzip")
				defer func() {
					if err := compressWriter.Close(); err != nil {
						fmt.Printf("Failed to close gzip body reader: %v", err)
					}
				}()
			}
			c.Next()
		} else {
			c.Next()
		}
	}
}
