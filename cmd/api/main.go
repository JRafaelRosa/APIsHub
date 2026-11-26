package main

import (
	"log"
	"net/http"
	"os"

	"github.com/usuario/apishub/internal/handler"
	"github.com/usuario/apishub/internal/proxy"
	"github.com/usuario/apishub/internal/repository"
	"github.com/usuario/apishub/internal/service"
)

func main() {
	porta := os.Getenv("PORT")
	if porta == "" {
		porta = ":8080"
	}
	if porta[0] != ':' {
		porta = ":" + porta
	}

	repository := repository.NewMemoryRouteRepository()
	routeService := service.NewRouteService(repository)
	reverseProxy := proxy.NewReverseProxyAdapter(routeService, 0)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.HealthHandler)
	handler.RegisterRouteHandlers(mux, routeService)
	mux.Handle("/apihub/", reverseProxy)

	log.Printf("servidor iniciado em %s", porta)
	if err := http.ListenAndServe(porta, handler.LoggingMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
