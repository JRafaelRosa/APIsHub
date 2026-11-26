package proxy

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/usuario/apishub/internal/domain"
)

type routeServiceStub struct {
	routes map[string]*domain.APIRoute
}

func (s *routeServiceStub) Criar(ctx context.Context, rota *domain.APIRoute) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.routes[rota.Slug] = rota
	return nil
}

func (s *routeServiceStub) ObterPorID(ctx context.Context, id string) (*domain.APIRoute, error) {
	return nil, domain.ErrRotaNaoEncontrada
}

func (s *routeServiceStub) ObterPorSlug(ctx context.Context, slug string) (*domain.APIRoute, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	rota, existe := s.routes[slug]
	if !existe || !rota.Ativa {
		return nil, domain.ErrRotaNaoEncontrada
	}
	return rota, nil
}

func (s *routeServiceStub) AtualizarURL(ctx context.Context, slug string, novaURL string) (*domain.APIRoute, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	rota, existe := s.routes[slug]
	if !existe {
		return nil, domain.ErrRotaNaoEncontrada
	}
	rota.URLOriginal = novaURL
	return rota, nil
}

func (s *routeServiceStub) Remover(ctx context.Context, slug string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, existe := s.routes[slug]; !existe {
		return domain.ErrRotaNaoEncontrada
	}
	delete(s.routes, slug)
	return nil
}

func novoProxy(t *testing.T, routes map[string]*domain.APIRoute, timeout time.Duration) http.Handler {
	t.Helper()
	adapter := NewReverseProxyAdapter(&routeServiceStub{routes: routes}, timeout)
	return adapter
}

func TestReverseProxyEncaminhaGETQuerySubrotaHostECabecalhos(t *testing.T) {
	external := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("método esperado GET, recebido %s", r.Method)
		}
		if r.URL.Path != "/clima/regiao/hoje" {
			t.Errorf("path esperado /clima/regiao/hoje, recebido %s", r.URL.Path)
		}
		if r.URL.RawQuery != "uf=PR&cidade=curitiba" {
			t.Errorf("query inesperada: %s", r.URL.RawQuery)
		}
		if r.Host != strings.TrimPrefix(externalURLPlaceholder, "http://") {
			t.Errorf("Host inesperado: %s", r.Host)
		}
		w.Header().Set("Server", "api-externa")
		w.Header().Set("X-Powered-By", "infra-externa")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte(`{"resultado":"ok"}`))
	}))
	defer external.Close()
	externalURLPlaceholder = external.URL

	routes := map[string]*domain.APIRoute{
		"clima": {Slug: "clima", URLOriginal: external.URL + "/clima", Ativa: true},
	}
	proxy := novoProxy(t, routes, 500*time.Millisecond)
	req := httptest.NewRequest(http.MethodGet, "/apihub/clima/regiao/hoje?uf=PR&cidade=curitiba", nil)
	req.Host = "localhost:8080"
	rec := httptest.NewRecorder()

	proxy.ServeHTTP(rec, req)

	if rec.Code != http.StatusTeapot {
		t.Fatalf("status esperado %d, recebido %d", http.StatusTeapot, rec.Code)
	}
	if rec.Header().Get("Server") != "" {
		t.Errorf("Server não foi removido")
	}
	if rec.Header().Get("X-Powered-By") != "" {
		t.Errorf("X-Powered-By não foi removido")
	}
	if rec.Body.String() != `{"resultado":"ok"}` {
		t.Errorf("corpo inesperado: %s", rec.Body.String())
	}
}

var externalURLPlaceholder string

func TestReverseProxyEncaminhaPOSTBody(t *testing.T) {
	external := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if r.Method != http.MethodPost {
			t.Errorf("método esperado POST, recebido %s", r.Method)
		}
		if r.URL.Path != "/v2/pix" {
			t.Errorf("path inesperado: %s", r.URL.Path)
		}
		if string(body) != `{"valor":10}` {
			t.Errorf("body inesperado: %s", string(body))
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"123"}`))
	}))
	defer external.Close()

	routes := map[string]*domain.APIRoute{
		"pagamento/pix": {Slug: "pagamento/pix", URLOriginal: external.URL + "/v2/pix", Ativa: true},
	}
	proxy := novoProxy(t, routes, 500*time.Millisecond)
	req := httptest.NewRequest(http.MethodPost, "/apihub/pagamento/pix", strings.NewReader(`{"valor":10}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	proxy.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status esperado %d, recebido %d", http.StatusCreated, rec.Code)
	}
	if rec.Body.String() != `{"id":"123"}` {
		t.Fatalf("corpo inesperado: %s", rec.Body.String())
	}
}

func TestReverseProxySlugInexistente(t *testing.T) {
	proxy := novoProxy(t, map[string]*domain.APIRoute{}, time.Second)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/apihub/inexistente", nil)

	proxy.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status esperado %d, recebido %d", http.StatusNotFound, rec.Code)
	}
	if rec.Body.String() != "{\"erro\":\"rota não encontrada\"}\n" {
		t.Fatalf("corpo inesperado: %s", rec.Body.String())
	}
}

func TestReverseProxyTimeout(t *testing.T) {
	external := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer external.Close()

	routes := map[string]*domain.APIRoute{
		"lento": {Slug: "lento", URLOriginal: external.URL, Ativa: true},
	}
	proxy := novoProxy(t, routes, 100*time.Millisecond)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/apihub/lento", nil)

	proxy.ServeHTTP(rec, req)

	if rec.Code != http.StatusGatewayTimeout {
		t.Fatalf("status esperado %d, recebido %d", http.StatusGatewayTimeout, rec.Code)
	}
	if rec.Body.String() != "{\"erro\":\"gateway timeout\"}\n" {
		t.Fatalf("corpo inesperado: %s", rec.Body.String())
	}
}

func TestReverseProxyFalhaConexao(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	_ = listener.Close()

	routes := map[string]*domain.APIRoute{
		"fora": {Slug: "fora", URLOriginal: "http://" + address, Ativa: true},
	}
	proxy := novoProxy(t, routes, 500*time.Millisecond)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/apihub/fora", nil)

	proxy.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status esperado %d, recebido %d", http.StatusBadGateway, rec.Code)
	}
	if rec.Body.String() != "{\"erro\":\"bad gateway\"}\n" {
		t.Fatalf("corpo inesperado: %s", rec.Body.String())
	}
}

var _ = errors.New
