package handler

import (
	"compress/gzip"
	"fmt"
	"io"
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

type gzipBodyReader struct {
	reader     io.ReadCloser
	gzipReader *gzip.Reader
}

func (g *gzipBodyReader) Read(p []byte) (n int, err error) {
	return g.gzipReader.Read(p)
}

func (g *gzipBodyReader) Close() error {
	if err := g.reader.Close(); err != nil {
		return err
	}
	if err := g.gzipReader.Close(); err != nil {
		return err
	}
	return nil
}

func newGzipBodyReader(body io.ReadCloser) (*gzipBodyReader, error) {
	gzipReader, err := gzip.NewReader(body)
	if err != nil {
		return nil, err
	}

	return &gzipBodyReader{
		reader:     body,
		gzipReader: gzipReader,
	}, nil
}

type gzipBodyWriter struct {
	writer     http.ResponseWriter
	gzipWriter *gzip.Writer
}

func (c *gzipBodyWriter) Close() error {
	if err := c.gzipWriter.Close(); err != nil {
		fmt.Printf("Failed to close gzip body writer: %v", err)
	}
	return nil
}

func newGzipBodyWriter(w http.ResponseWriter) *gzipBodyWriter {
	return &gzipBodyWriter{
		writer:     w,
		gzipWriter: gzip.NewWriter(w),
	}
}
