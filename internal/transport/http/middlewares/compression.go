package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

func Compression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the client accepts gzip encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Create a gzip writer
		w.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		// wrap the response writer
		gzw := &gzipResponseWriter{
			ResponseWriter: w,
			writer:         gz,
		}
		next.ServeHTTP(gzw, r)

		fmt.Println("Sent response with gzip encoding from middleware")
	})
}

// Response writer
type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (gz *gzipResponseWriter) Write(b []byte) (int, error) {
	return gz.writer.Write(b)
}
