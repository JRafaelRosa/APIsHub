package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/usuario/apishub/internal/domain"
	"github.com/usuario/apishub/internal/proxy"
	"github.com/usuario/apishub/internal/repository"
	"github.com/usuario/apishub/internal/service"
)

func novoMuxTeste() http.Handler {
	r := repository.NewMemoryRouteRepository()
	s := service.NewRouteService(r)
	h := NewRouteHandler(s)
	mux := http.NewServeMux()
	mux.HandleFunc("/", WebHandler)
	mux.HandleFunc("GET /api/catalogo", h.Catalogo)
	mux.Handle("/api/rotas", h)
	mux.Handle("/api/rotas/", h)
	mux.Handle("/apihub/", proxy.NewReverseProxyAdapter(s, 0))
	return mux
}

func TestPUTEDLETEERotaEProxyRetorna404(t *testing.T) {
	mux := novoMuxTeste()
	criar := httptest.NewRequest(http.MethodPost, "/api/rotas", strings.NewReader(`{"nome":"Teste","slug":"teste","url_original":"http://127.0.0.1:1"}`))
	criar.Header.Set("Content-Type", "application/json")
	criarRec := httptest.NewRecorder()
	mux.ServeHTTP(criarRec, criar)
	if criarRec.Code != http.StatusCreated {
		t.Fatalf("cadastro: status %d", criarRec.Code)
	}

	alterar := httptest.NewRequest(http.MethodPut, "/api/rotas/teste", strings.NewReader(`{"url_original":"https://novo.exemplo"}`))
	alterarRec := httptest.NewRecorder()
	mux.ServeHTTP(alterarRec, alterar)
	if alterarRec.Code != http.StatusOK {
		t.Fatalf("PUT: status %d, corpo %s", alterarRec.Code, alterarRec.Body)
	}
	var rota domain.APIRoute
	if err := json.Unmarshal(alterarRec.Body.Bytes(), &rota); err != nil || rota.Slug != "teste" || rota.URLOriginal != "https://novo.exemplo" {
		t.Fatalf("resposta do PUT inesperada: %+v, erro %v", rota, err)
	}

	excluir := httptest.NewRequest(http.MethodDelete, "/api/rotas/teste", nil)
	excluirRec := httptest.NewRecorder()
	mux.ServeHTTP(excluirRec, excluir)
	if excluirRec.Code != http.StatusNoContent {
		t.Fatalf("DELETE: status %d", excluirRec.Code)
	}
	proxyReq := httptest.NewRequest(http.MethodGet, "/apihub/teste", nil)
	proxyRec := httptest.NewRecorder()
	mux.ServeHTTP(proxyRec, proxyReq)
	if proxyRec.Code != http.StatusNotFound {
		t.Fatalf("proxy após exclusão: status %d", proxyRec.Code)
	}
}

func TestPaginaWebEmbutida(t *testing.T) {
	rec := httptest.NewRecorder()
	WebHandler(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "APIHub") || !strings.Contains(rec.Body.String(), "/api/catalogo") {
		t.Fatalf("página web não entregue corretamente: status %d", rec.Code)
	}
}
