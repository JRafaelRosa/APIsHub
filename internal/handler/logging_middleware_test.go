package handler

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddlewareRegistraRequisicao(t *testing.T) {
	var output bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&output, nil)))
	defer slog.SetDefault(previous)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	handler := LoggingMiddleware(next)
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/apihub/clima/regiao", nil))

	logLine := output.String()
	for _, field := range []string{
		"method=GET",
		"path=/apihub/clima/regiao",
		"slug=clima/regiao",
		"status=404",
		"latency_ms=",
	} {
		if !strings.Contains(logLine, field) {
			t.Errorf("log não contém %q: %s", field, logLine)
		}
	}
}
