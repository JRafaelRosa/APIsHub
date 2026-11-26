package service

import (
	"context"

	"github.com/usuario/apishub/internal/domain"
	"github.com/usuario/apishub/internal/repository"
)

// RouteService define a camada de regras de negócio para rotas.
type RouteService interface {
	Criar(ctx context.Context, rota *domain.APIRoute) error
	ObterPorID(ctx context.Context, id string) (*domain.APIRoute, error)
	ObterPorSlug(ctx context.Context, slug string) (*domain.APIRoute, error)
}

// routeService é a implementação injetável do RouteService.
type routeService struct {
	repository repository.RouteRepository
}

func NewRouteService(repo repository.RouteRepository) RouteService {
	return &routeService{repository: repo}
}

func (s *routeService) Criar(ctx context.Context, rota *domain.APIRoute) error {
	return s.repository.Criar(ctx, rota)
}

func (s *routeService) ObterPorID(ctx context.Context, id string) (*domain.APIRoute, error) {
	return s.repository.ObterPorID(ctx, id)
}

func (s *routeService) ObterPorSlug(ctx context.Context, slug string) (*domain.APIRoute, error) {
	return s.repository.ObterPorSlug(ctx, slug)
}
