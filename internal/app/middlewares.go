package app

import (
	"log"
	"net/http"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func withLoggingMiddleware(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		lw := &loggingResponseWriter{w, http.StatusOK}
		log.Printf("%s %s", r.Method, r.URL.String())
		h(lw, r)
		log.Printf("%d", lw.statusCode)
	}
}
