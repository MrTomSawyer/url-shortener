package handler

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
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
		decompressTypes := [2]string{"application/json", "text/html"}
		if c.Request.Header.Get("Content-Encoding") != "gzip" {
			c.Next()
			return
		}

		isCorrectType := false
		contentType := c.Request.Header.Get("Content-Type")
		for _, val := range decompressTypes {
			if val == contentType {
				isCorrectType = true
			}
		}
		if !isCorrectType {
			c.Next()
			return
		}

		var buffer bytes.Buffer
		gzipReader, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer func() {
			err := gzipReader.Close()
			if err != nil {
				fmt.Println("Failed to close gzip reader", err)
			}
		}()

		_, err = buffer.ReadFrom(gzipReader)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		c.Request.Body = io.NopCloser(&buffer)
		c.Next()
	}
}
