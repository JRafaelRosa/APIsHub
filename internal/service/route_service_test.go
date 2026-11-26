package service

import (
	"context"
	"errors"
	"testing"

	"github.com/usuario/apishub/internal/domain"
	"github.com/usuario/apishub/internal/repository"
)

func TestAtualizarURLValidaDestinoEPreservaSlug(t *testing.T) {
	service := NewRouteService(repository.NewMemoryRouteRepository())
	rota := &domain.APIRoute{ID: "1", Slug: "pagamento/pix", URLOriginal: "https://antiga.example"}
	if err := service.Criar(context.Background(), rota); err != nil {
		t.Fatal(err)
	}

	if _, err := service.AtualizarURL(context.Background(), rota.Slug, "ftp://nova.example"); !errors.Is(err, domain.ErrURLInvalida) {
		t.Fatalf("erro esperado para URL inválida, recebido %v", err)
	}
	atualizada, err := service.AtualizarURL(context.Background(), rota.Slug, "https://nova.example/v2")
	if err != nil {
		t.Fatal(err)
	}
	if atualizada.Slug != rota.Slug || atualizada.URLOriginal != "https://nova.example/v2" {
		t.Fatalf("rota atualizada inesperada: %+v", atualizada)
	}
	if _, err := service.AtualizarURL(context.Background(), "inexistente", "https://nova.example"); !errors.Is(err, domain.ErrRotaNaoEncontrada) {
		t.Fatalf("erro esperado para rota inexistente, recebido %v", err)
	}
}
