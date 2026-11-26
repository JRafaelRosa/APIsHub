package domain

import (
	"errors"
	"time"
)

var (
	ErrRotaNaoEncontrada = errors.New("rota não encontrada")
	ErrSlugDuplicado     = errors.New("slug duplicado")
	ErrURLInvalida       = errors.New("URL inválida")
)

// APIRoute representa os metadados de uma API cadastrada no APIHub.
type APIRoute struct {
	ID          string    `json:"id"`
	Nome        string    `json:"nome"`
	Descricao   string    `json:"descricao"`
	Slug        string    `json:"slug"`
	URLOriginal string    `json:"url_original"`
	Ativa       bool      `json:"ativa"`
	CriadaEm    time.Time `json:"criada_em"`
}
