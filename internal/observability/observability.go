// Package observability é dono do logging estruturado e da correlação de
// requests.
//
// A task 002-13 o mantém mínimo: um logger slog JSON e o request ID carregado
// no context. Métricas, logs de request e regras de redaction pertencem à
// task 002-10.
package observability

import (
	"context"
	"io"
	"log/slog"
)

// NewLogger devolve um logger JSON. Quem chama nunca deve passar secrets,
// tokens, Firebase UIDs ou Player Tags como atributos.
func NewLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, nil))
}

// requestIDKey não é exportado, então nenhum outro pacote consegue ler ou
// sobrescrever o valor por acidente; o acesso passa só pelas funções abaixo.
type requestIDKey struct{}

// WithRequestID devolve uma cópia de ctx carregando id.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestID devolve o request ID guardado em ctx, ou "" quando ausente.
func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}
