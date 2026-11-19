package service

import (
	"context"
	"testing"

	"github.com/usuario/apishub/internal/domain"
	"github.com/usuario/apishub/internal/repository"
)

func TestAtualizarURLMantemSlugEValidaURL(t *testing.T) {
	r := repository.NewMemoryRouteRepository()
	s := NewRouteService(r)
	rota := &domain.APIRoute{ID: "clima", Slug: "clima", URLOriginal: "https://antiga.exemplo"}
	if err := s.Criar(context.Background(), rota); err != nil {
		t.Fatal(err)
	}
	atualizada, err := s.AtualizarURL(context.Background(), "clima", "http://nova.exemplo/api")
	if err != nil {
		t.Fatal(err)
	}
	if atualizada.Slug != "clima" || atualizada.URLOriginal != "http://nova.exemplo/api" {
		t.Fatalf("rota atualizada incorretamente: %+v", atualizada)
	}
	if _, err := s.AtualizarURL(context.Background(), "clima", "ftp://invalida"); err != domain.ErrURLInvalida {
		t.Fatalf("erro esperado para URL inválida, recebido %v", err)
	}
}

func TestRemoverRota(t *testing.T) {
	r := repository.NewMemoryRouteRepository()
	s := NewRouteService(r)
	if err := s.Criar(context.Background(), &domain.APIRoute{ID: "x", Slug: "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Remover(context.Background(), "x"); err != nil {
		t.Fatal(err)
	}
	if err := s.Remover(context.Background(), "x"); err != domain.ErrRotaNaoEncontrada {
		t.Fatalf("erro esperado para rota ausente, recebido %v", err)
	}
}
