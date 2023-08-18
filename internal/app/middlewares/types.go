// Package middlewares provides middleware functions for handling various aspects of HTTP requests and responses.
package middlewares

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

// gzipBodyReader wraps a gzip.Reader and implements io.ReadCloser interface.
type gzipBodyReader struct {
	reader     io.ReadCloser
	gzipReader *gzip.Reader
}

// Read reads compressed data using gzip.Reader.
func (g *gzipBodyReader) Read(p []byte) (n int, err error) {
	return g.gzipReader.Read(p)
}

// Close closes both the original body and the gzip.Reader.
func (g *gzipBodyReader) Close() error {
	if err := g.reader.Close(); err != nil {
		return err
	}
	if err := g.gzipReader.Close(); err != nil {
		return err
	}
	return nil
}

// newGzipBodyReader creates a new gzipBodyReader from an existing io.ReadCloser.
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

// gzipBodyWriter wraps a gzip.Writer and provides a Close method.
type gzipBodyWriter struct {
	writer     http.ResponseWriter
	gzipWriter *gzip.Writer
}

// Close closes the gzip.Writer.
func (c *gzipBodyWriter) Close() error {
	if err := c.gzipWriter.Close(); err != nil {
		fmt.Printf("Failed to close gzip body writer: %v", err)
	}
	return nil
}

// newGzipBodyWriter creates a new gzipBodyWriter from an existing http.ResponseWriter.
func newGzipBodyWriter(w http.ResponseWriter) *gzipBodyWriter {
	return &gzipBodyWriter{
		writer:     w,
		gzipWriter: gzip.NewWriter(w),
	}
}
