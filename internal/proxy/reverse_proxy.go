package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/usuario/apishub/internal/domain"
	"github.com/usuario/apishub/internal/service"
)

const defaultTimeout = 5 * time.Second

// ReverseProxyAdapter encaminha requisições mascaradas para APIs externas.
type ReverseProxyAdapter struct {
	routeService service.RouteService
	proxy        *httputil.ReverseProxy
	timeout      time.Duration
}

// NewReverseProxyAdapter cria um adaptador. Timeout <= 0 usa o padrão de 5 segundos.
func NewReverseProxyAdapter(routeService service.RouteService, timeout time.Duration) *ReverseProxyAdapter {
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	adapter := &ReverseProxyAdapter{
		routeService: routeService,
		timeout:      timeout,
	}

	adapter.proxy = &httputil.ReverseProxy{
		Director:       adapter.director,
		ModifyResponse: adapter.modifyResponse,
		ErrorHandler:   adapter.errorHandler,
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	return adapter
}

func (a *ReverseProxyAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rota, subPath, encontrada := a.resolverRota(r.Context(), r.URL.Path)
	if !encontrada {
		responderJSON(w, http.StatusNotFound, "rota não encontrada")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), a.timeout)
	defer cancel()
	ctx = context.WithValue(ctx, rotaContextKey{}, rota)
	ctx = context.WithValue(ctx, subPathContextKey{}, subPath)
	a.proxy.ServeHTTP(w, r.WithContext(ctx))
}

type rotaContextKey struct{}
type subPathContextKey struct{}

func (a *ReverseProxyAdapter) resolverRota(ctx context.Context, path string) (*domain.APIRoute, string, bool) {
	prefixo := strings.TrimPrefix(path, "/apihub/")
	segmentos := strings.Split(strings.Trim(prefixo, "/"), "/")
	for i := len(segmentos); i >= 1; i-- {
		candidato := strings.Join(segmentos[:i], "/")
		rota, err := a.routeService.ObterPorSlug(ctx, candidato)
		if err == nil && rota != nil && rota.Ativa {
			sobra := strings.Join(segmentos[i:], "/")
			if sobra != "" {
				return rota, "/" + sobra, true
			}
			return rota, "", true
		}
	}
	return nil, "", false
}

func (a *ReverseProxyAdapter) director(req *http.Request) {
	rota, _ := req.Context().Value(rotaContextKey{}).(*domain.APIRoute)
	subPath, _ := req.Context().Value(subPathContextKey{}).(string)
	if rota == nil {
		return
	}

	target, err := url.Parse(rota.URLOriginal)
	if err != nil || target.Scheme == "" || target.Host == "" {
		return
	}

	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path = concatenarPaths(target.Path, subPath)
	req.URL.RawPath = ""
	req.Host = target.Host
}

func concatenarPaths(basePath, subPath string) string {
	basePath = "/" + strings.Trim(basePath, "/")
	if basePath == "/" {
		basePath = ""
	}

	subPath = "/" + strings.Trim(subPath, "/")
	if subPath == "/" {
		subPath = ""
	}

	return basePath + subPath
}

func (a *ReverseProxyAdapter) modifyResponse(res *http.Response) error {
	res.Header.Del("Server")
	res.Header.Del("X-Powered-By")
	return nil
}

func (a *ReverseProxyAdapter) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(r.Context().Err(), context.DeadlineExceeded) {
		responderJSON(w, http.StatusGatewayTimeout, "gateway timeout")
		return
	}
	responderJSON(w, http.StatusBadGateway, "bad gateway")
}

func responderJSON(w http.ResponseWriter, status int, erro string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"erro": erro})
}

var _ http.Handler = (*ReverseProxyAdapter)(nil)
var _ = domain.ErrRotaNaoEncontrada
