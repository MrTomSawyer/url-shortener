package middlewares

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

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
