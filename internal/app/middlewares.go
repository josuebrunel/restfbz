package app

import (
	"database/sql"
	"log"
	"net/http"
	"restfbz/pkg/stats"
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

func withStatsMiddleware(h http.HandlerFunc, db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
		sr := stats.New(db)
		url := r.URL.String()
		err := sr.CreateRecord(url)
		if err != nil {
			log.Printf("Failed to create record for %s", url)
		}
	}
}
