package main

import (
	"log/slog"
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
	routeHandler := handler.NewRouteHandler(routeService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.HealthHandler)
	mux.HandleFunc("/", handler.WebHandler)
	mux.HandleFunc("GET /api/catalogo", routeHandler.Catalogo)
	mux.Handle("/api/rotas/", routeHandler)
	mux.Handle("/api/rotas", routeHandler)
	mux.Handle("/apihub/", reverseProxy)

	slog.Info("servidor iniciado", "porta", porta)
	if err := http.ListenAndServe(porta, handler.LoggingMiddleware(mux)); err != nil {
		slog.Error("servidor encerrado", "erro", err)
	}
}
