package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/usuario/apishub/internal/domain"
	"github.com/usuario/apishub/internal/proxy"
	"github.com/usuario/apishub/internal/repository"
	"github.com/usuario/apishub/internal/service"
)

func TestAtualizarRemoverERetornar404NoProxy(t *testing.T) {
	antigo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "destino antigo")
	}))
	defer antigo.Close()
	novo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = io.WriteString(w, "destino novo")
	}))
	defer novo.Close()

	repo := repository.NewMemoryRouteRepository()
	routeService := service.NewRouteService(repo)
	if err := routeService.Criar(context.Background(), &domain.APIRoute{
		ID: "1", Slug: "pagamento/pix", URLOriginal: antigo.URL, Ativa: true,
	}); err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	RegisterRouteHandlers(mux, routeService)
	mux.Handle("/apihub/", proxy.NewReverseProxyAdapter(routeService, 0))

	put := httptest.NewRequest(http.MethodPut, "/api/rotas/pagamento/pix", strings.NewReader(`{"url_original":"`+novo.URL+`"}`))
	put.Header.Set("Content-Type", "application/json")
	putResponse := httptest.NewRecorder()
	mux.ServeHTTP(putResponse, put)
	if putResponse.Code != http.StatusOK {
		t.Fatalf("status PUT esperado %d, recebido %d: %s", http.StatusOK, putResponse.Code, putResponse.Body)
	}
	if !strings.Contains(putResponse.Body.String(), `"slug":"pagamento/pix"`) ||
		!strings.Contains(putResponse.Body.String(), `"url_original":"`+novo.URL+`"`) {
		t.Fatalf("resposta PUT inesperada: %s", putResponse.Body)
	}

	proxied := httptest.NewRecorder()
	mux.ServeHTTP(proxied, httptest.NewRequest(http.MethodGet, "/apihub/pagamento/pix", nil))
	if proxied.Code != http.StatusAccepted || proxied.Body.String() != "destino novo" {
		t.Fatalf("proxy não encaminhou ao destino atualizado: status=%d body=%q", proxied.Code, proxied.Body.String())
	}

	deleteRequest := httptest.NewRequest(http.MethodDelete, "/api/rotas/pagamento/pix", nil)
	deleteResponse := httptest.NewRecorder()
	mux.ServeHTTP(deleteResponse, deleteRequest)
	if deleteResponse.Code != http.StatusNoContent || deleteResponse.Body.Len() != 0 {
		t.Fatalf("resposta DELETE inesperada: status=%d body=%q", deleteResponse.Code, deleteResponse.Body.String())
	}

	proxied = httptest.NewRecorder()
	mux.ServeHTTP(proxied, httptest.NewRequest(http.MethodGet, "/apihub/pagamento/pix", nil))
	if proxied.Code != http.StatusNotFound {
		t.Fatalf("proxy após remoção esperado %d, recebido %d", http.StatusNotFound, proxied.Code)
	}
}

func TestEndpointsDeRotaRetornamErrosDeValidacaoENaoEncontrada(t *testing.T) {
	routeService := service.NewRouteService(repository.NewMemoryRouteRepository())
	mux := http.NewServeMux()
	RegisterRouteHandlers(mux, routeService)

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		status int
	}{
		{name: "URL inválida", method: http.MethodPut, path: "/api/rotas/inexistente", body: `{"url_original":"ftp://destino.example"}`, status: http.StatusBadRequest},
		{name: "payload inválido", method: http.MethodPut, path: "/api/rotas/inexistente", body: `{`, status: http.StatusBadRequest},
		{name: "rota inexistente no PUT", method: http.MethodPut, path: "/api/rotas/inexistente", body: `{"url_original":"https://destino.example"}`, status: http.StatusNotFound},
		{name: "rota inexistente no DELETE", method: http.MethodDelete, path: "/api/rotas/inexistente", status: http.StatusNotFound},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			mux.ServeHTTP(response, request)
			if response.Code != test.status {
				t.Fatalf("status esperado %d, recebido %d: %s", test.status, response.Code, response.Body)
			}
		})
	}
}
