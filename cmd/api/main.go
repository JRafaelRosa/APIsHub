package main

import (
	"log"
	"net/http"
	"os"

	"github.com/usuario/apishub/internal/handler"
)

func main() {
	porta := os.Getenv("PORT")
	if porta == "" {
		porta = ":8080"
	}
	if porta[0] != ':' {
		porta = ":" + porta
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.HealthHandler)

	log.Printf("servidor iniciado em %s", porta)
	if err := http.ListenAndServe(porta, mux); err != nil {
		log.Fatal(err)
	}
}
