package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/usuario/apishub/internal/domain"
	"github.com/usuario/apishub/internal/service"
)

func RegisterRouteHandlers(mux *http.ServeMux, routeService service.RouteService) {
	mux.HandleFunc("PUT /api/rotas/{slug...}", func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			URLOriginal string `json:"url_original"`
		}
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
		if err := decoder.Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "payload inválido")
			return
		}
		if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
			writeError(w, http.StatusBadRequest, "payload inválido")
			return
		}

		rota, err := routeService.AtualizarURL(r.Context(), r.PathValue("slug"), payload.URLOriginal)
		switch {
		case errors.Is(err, domain.ErrURLInvalida):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, domain.ErrRotaNaoEncontrada):
			writeError(w, http.StatusNotFound, err.Error())
		case err != nil:
			writeError(w, http.StatusInternalServerError, "erro interno")
		default:
			writeJSON(w, http.StatusOK, rota)
		}
	})

	mux.HandleFunc("DELETE /api/rotas/{slug...}", func(w http.ResponseWriter, r *http.Request) {
		err := routeService.Remover(r.Context(), r.PathValue("slug"))
		switch {
		case errors.Is(err, domain.ErrRotaNaoEncontrada):
			writeError(w, http.StatusNotFound, err.Error())
		case err != nil:
			writeError(w, http.StatusInternalServerError, "erro interno")
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	})
}
