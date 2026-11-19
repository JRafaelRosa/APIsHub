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
	ID          string
	Nome        string
	Descricao   string
	Slug        string
	URLOriginal string
	Ativa       bool
	CriadaEm    time.Time
}
