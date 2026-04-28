package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		userAgent := r.UserAgent()
		if userAgent == "" {
			userAgent = "-"
		}

		s.logger.Info(formatLogEntry(
			r.RemoteAddr,
			start,
			r.Method,
			r.RequestURI,
			r.Proto,
			rw.statusCode,
			duration,
			userAgent,
		))
	})
}

func formatLogEntry(ip string, timestamp time.Time, method, uri, proto string, status int, latency time.Duration, userAgent string) string {
	return ip + " " +
		"[" + timestamp.Format("02/Jan/2006:15:04:05 -0700") + "] " +
		method + " " +
		uri + " " +
		proto + " " +
		fmt.Sprintf("%d", status) + " " +
		fmt.Sprintf("%d", latency.Microseconds()) + " " +
		"\"" + userAgent + "\""
}
