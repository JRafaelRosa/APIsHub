package handler

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
func (w *statusWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}

// LoggingMiddleware registra as informações essenciais de cada requisição.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		sw := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r)
		status := sw.status
		if status == 0 {
			status = http.StatusOK
		}
		slug := strings.TrimPrefix(r.URL.Path, "/apihub/")
		if i := strings.IndexByte(slug, '/'); i >= 0 {
			slug = slug[:i]
		}
		slog.Default().Info("requisição HTTP", "method", r.Method, "path", r.URL.Path, "slug", slug, "status", status, "latency_ms", time.Since(inicio).Milliseconds())
	})
}
