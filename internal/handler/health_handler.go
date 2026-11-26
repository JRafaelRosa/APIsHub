package handler

import (
	"net/http"
)

// HealthHandler responde ao health check da aplicação.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
