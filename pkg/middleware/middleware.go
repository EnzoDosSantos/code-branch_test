package middleware

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
    lrw.statusCode = code
    lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        lrw := &loggingResponseWriter{w, http.StatusOK}
        
        next.ServeHTTP(lrw, r)
        
        log.Printf(
            "%s %s %d %s",
            r.Method,
            r.URL.Path,
            lrw.statusCode,
            time.Since(start),
        )
    })
}

func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                log.Printf("panic: %v", err)
            }
        }()
        next.ServeHTTP(w, r)
    })
}