package utils

import (
	"compress/gzip"
	"io"
	"net/http"
)

// GzipWriter provides the same functionality as http.ResponseWriter
// It compresses data using compress/zip writer and provides headers
// from given http.ResponseWriter
type GzipWriter struct {
	GzipWriter *gzip.Writer
	Writer     http.ResponseWriter
}

// Header implements Header() API of http.ResponseWriter interface
func (g GzipWriter) Header() http.Header {
	return g.Writer.Header()
}

// Write implements Write API of http.ResponseWriter interface
func (g GzipWriter) Write(b []byte) (int, error) {
	return g.GzipWriter.Write(b)
}

// WriteHeader implements WriteHeader API of http.ResponseWriter interface
func (g GzipWriter) WriteHeader(statusCode int) {
	g.Writer.WriteHeader(statusCode)
}

// GzipReader struct to handle GZip'ed content of HTTP requests
type GzipReader struct {
	*gzip.Reader
	io.Closer
}

// Close function closes gzip reader
func (gz GzipReader) Close() error {
	return gz.Closer.Close()
}
