package main

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

// recoverPanic ensures that in the case of a panic, a Connection header of
// 'close' is sent to the client.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			pv := recover()
			if pv != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%v", pv))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// requestLogger logs the HTTP request's method and URL path.
func (app *application) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("Request received", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// gzipResponseWriter is a light wrapper around http.ResponseWriter that
// compresses responses written in the gzip format.
type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

// Write uses a *gzip.Writer instead of the http.ResponseWriter.
func (gzw gzipResponseWriter) Write(b []byte) (int, error) {
	return gzw.writer.Write(b)
}

// gzip compress responses if the client accepts the gzip encoding in its HTTP request.
func (app *application) gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check that the client set the Accept-Encoding header to "gzip"
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// set Content-Encoding and account for caching
		w.Header().Add("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")

		// create a new *gzip.Writer and gzipResponseWriter
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := gzipResponseWriter{
			ResponseWriter: w,
			writer:         gz,
		}

		// use the gzipResponseWriter in the next handler
		next.ServeHTTP(gzw, r)
	})
}
