package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/usuario/apishub/internal/domain"
)

func TestAtualizarURLEPreservaSlug(t *testing.T) {
	repo := NewMemoryRouteRepository()
	rota := &domain.APIRoute{ID: "1", Slug: "pagamento/pix", URLOriginal: "https://antiga.example"}
	if err := repo.Criar(context.Background(), rota); err != nil {
		t.Fatal(err)
	}

	anterior, err := repo.ObterPorSlug(context.Background(), rota.Slug)
	if err != nil {
		t.Fatal(err)
	}
	atualizada, err := repo.AtualizarURL(context.Background(), rota.Slug, "https://nova.example")
	if err != nil {
		t.Fatal(err)
	}
	if atualizada.Slug != rota.Slug || atualizada.URLOriginal != "https://nova.example" {
		t.Fatalf("rota atualizada inesperada: %+v", atualizada)
	}
	if anterior.URLOriginal != "https://antiga.example" {
		t.Fatalf("rota previamente lida foi modificada: %+v", anterior)
	}
}

func TestRemoverEOperacoesComRotaInexistente(t *testing.T) {
	repo := NewMemoryRouteRepository()
	rota := &domain.APIRoute{ID: "1", Slug: "pagamento/pix"}
	if err := repo.Criar(context.Background(), rota); err != nil {
		t.Fatal(err)
	}
	if err := repo.Remover(context.Background(), rota.Slug); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.ObterPorSlug(context.Background(), rota.Slug); !errors.Is(err, domain.ErrRotaNaoEncontrada) {
		t.Fatalf("erro esperado para rota removida, recebido %v", err)
	}
	if _, err := repo.AtualizarURL(context.Background(), rota.Slug, "https://nova.example"); !errors.Is(err, domain.ErrRotaNaoEncontrada) {
		t.Fatalf("erro esperado ao atualizar rota inexistente, recebido %v", err)
	}
	if err := repo.Remover(context.Background(), rota.Slug); !errors.Is(err, domain.ErrRotaNaoEncontrada) {
		t.Fatalf("erro esperado ao remover rota inexistente, recebido %v", err)
	}
}
