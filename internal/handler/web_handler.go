package handler

import (
	"embed"
	"net/http"
)

//go:embed web/index.html
var webFiles embed.FS

func WebHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data, err := webFiles.ReadFile("web/index.html")
	if err != nil {
		http.Error(w, "página indisponível", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}
