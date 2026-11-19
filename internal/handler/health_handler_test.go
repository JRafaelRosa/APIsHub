package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)

	HealthHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status esperado %d, recebido %d", http.StatusOK, recorder.Code)
	}

	desejado := `{"status":"ok"}`
	recebido := recorder.Body.String()
	if recebido != desejado+"\n" {
		t.Fatalf("corpo esperado %q, recebido %q", desejado+"\n", recebido)
	}
}
