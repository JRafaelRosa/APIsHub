package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/usuario/apishub/internal/domain"
	"github.com/usuario/apishub/internal/service"
)

// RouteHandler expõe as operações administrativas de rotas.
type RouteHandler struct{ service service.RouteService }

func NewRouteHandler(s service.RouteService) *RouteHandler { return &RouteHandler{service: s} }

type rotaPayload struct {
	Nome        string `json:"nome"`
	Descricao   string `json:"descricao"`
	Slug        string `json:"slug"`
	URLOriginal string `json:"url_original"`
}

type atualizarURLPayload struct {
	URLOriginal string `json:"url_original"`
}

func (h *RouteHandler) Catalogo(w http.ResponseWriter, r *http.Request) {
	rotas, err := h.service.Listar(r.Context())
	if err != nil {
		responderErro(w, http.StatusInternalServerError, err.Error())
		return
	}
	tipo := make([]map[string]any, 0, len(rotas))
	for _, rota := range rotas {
		tipo = append(tipo, map[string]any{"id": rota.ID, "nome": rota.Nome, "descricao": rota.Descricao,
			"slug": rota.Slug, "url_mascarada": "/apihub/" + strings.TrimPrefix(rota.Slug, "/"),
			"url_original": rota.URLOriginal, "ativa": rota.Ativa, "criada_em": rota.CriadaEm})
	}
	responderJSON(w, http.StatusOK, tipo)
}

func (h *RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/rotas" && r.Method == http.MethodPost {
		h.criar(w, r)
		return
	}
	slug := strings.TrimPrefix(r.URL.Path, "/api/rotas/")
	slug = strings.Trim(slug, "/")
	if slug == "" {
		responderErro(w, http.StatusNotFound, "rota não encontrada")
		return
	}
	switch r.Method {
	case http.MethodPut:
		h.atualizar(w, r, slug)
	case http.MethodDelete:
		h.remover(w, r, slug)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *RouteHandler) criar(w http.ResponseWriter, r *http.Request) {
	var p rotaPayload
	if json.NewDecoder(r.Body).Decode(&p) != nil || p.Slug == "" || p.URLOriginal == "" {
		responderErro(w, http.StatusBadRequest, "payload inválido")
		return
	}
	rota := &domain.APIRoute{ID: p.Slug, Nome: p.Nome, Descricao: p.Descricao, Slug: strings.Trim(p.Slug, "/"), URLOriginal: p.URLOriginal, Ativa: true, CriadaEm: time.Now()}
	if err := h.service.Criar(r.Context(), rota); err != nil {
		statusErro(w, err)
		return
	}
	responderJSON(w, http.StatusCreated, rota)
}

func (h *RouteHandler) atualizar(w http.ResponseWriter, r *http.Request, slug string) {
	var p atualizarURLPayload
	if json.NewDecoder(r.Body).Decode(&p) != nil || p.URLOriginal == "" {
		responderErro(w, http.StatusBadRequest, "url_original é obrigatória")
		return
	}
	rota, err := h.service.AtualizarURL(r.Context(), slug, p.URLOriginal)
	if err != nil {
		statusErro(w, err)
		return
	}
	responderJSON(w, http.StatusOK, rota)
}

func (h *RouteHandler) remover(w http.ResponseWriter, r *http.Request, slug string) {
	if err := h.service.Remover(r.Context(), slug); err != nil {
		statusErro(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func statusErro(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrRotaNaoEncontrada) {
		responderErro(w, http.StatusNotFound, err.Error())
		return
	}
	if errors.Is(err, domain.ErrURLInvalida) || errors.Is(err, domain.ErrSlugDuplicado) {
		responderErro(w, http.StatusBadRequest, err.Error())
		return
	}
	responderErro(w, http.StatusInternalServerError, err.Error())
}

func responderErro(w http.ResponseWriter, status int, mensagem string) {
	responderJSON(w, status, map[string]string{"erro": mensagem})
}
func responderJSON(w http.ResponseWriter, status int, valor any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(valor)
}
