package repository

import (
	"context"
	"sync"

	"github.com/usuario/apishub/internal/domain"
)

// RouteRepository define as operações de persistência de rotas.
type RouteRepository interface {
	Criar(ctx context.Context, rota *domain.APIRoute) error
	ObterPorID(ctx context.Context, id string) (*domain.APIRoute, error)
	ObterPorSlug(ctx context.Context, slug string) (*domain.APIRoute, error)
	AtualizarURL(ctx context.Context, slug string, novaURL string) (*domain.APIRoute, error)
	Remover(ctx context.Context, slug string) error
}

// MemoryRouteRepository mantém rotas em memória com acesso protegido por mutex.
type MemoryRouteRepository struct {
	mu     sync.RWMutex
	routes map[string]*domain.APIRoute
}

func NewMemoryRouteRepository() *MemoryRouteRepository {
	return &MemoryRouteRepository{routes: make(map[string]*domain.APIRoute)}
}

func (r *MemoryRouteRepository) Criar(ctx context.Context, rota *domain.APIRoute) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existente := range r.routes {
		if existente.Slug == rota.Slug {
			return domain.ErrSlugDuplicado
		}
	}

	r.routes[rota.ID] = rota
	return nil
}

func (r *MemoryRouteRepository) ObterPorID(ctx context.Context, id string) (*domain.APIRoute, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	rota, existe := r.routes[id]
	if !existe {
		return nil, domain.ErrRotaNaoEncontrada
	}
	return rota, nil
}

func (r *MemoryRouteRepository) ObterPorSlug(ctx context.Context, slug string) (*domain.APIRoute, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, rota := range r.routes {
		if rota.Slug == slug {
			return rota, nil
		}
	}
	return nil, domain.ErrRotaNaoEncontrada
}

func (r *MemoryRouteRepository) AtualizarURL(ctx context.Context, slug string, novaURL string) (*domain.APIRoute, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for id, rota := range r.routes {
		if rota.Slug == slug {
			atualizada := *rota
			atualizada.URLOriginal = novaURL
			r.routes[id] = &atualizada
			return &atualizada, nil
		}
	}
	return nil, domain.ErrRotaNaoEncontrada
}

func (r *MemoryRouteRepository) Remover(ctx context.Context, slug string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for id, rota := range r.routes {
		if rota.Slug == slug {
			delete(r.routes, id)
			return nil
		}
	}
	return domain.ErrRotaNaoEncontrada
}
