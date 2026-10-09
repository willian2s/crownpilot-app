package httpapi

import (
	"crypto/rand"
	"log/slog"
	"net/http"
	"strings"

	"github.com/willian2s/crownpilot-app/internal/observability"
)

// Um middleware é qualquer func(http.Handler) http.Handler: recebe o próximo
// handler e devolve um novo, que executa código antes e/ou depois dele.

// withRequestID dá a cada request um ID de correlação novo. Um X-Request-ID
// recebido é ignorado de propósito: é entrada do cliente e poderia ser usado
// para forjar ou injetar entradas de log.
func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := observability.WithRequestID(r.Context(), rand.Text())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// recoverPanic transforma um panic de handler num 500 ProblemDetails genérico.
// Só o request ID é logado junto: o valor do panic pode conter dados do request.
func recoverPanic(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				// http.ErrAbortHandler é a forma oficial do net/http de abortar
				// uma resposta; ele precisa continuar propagando.
				if v == http.ErrAbortHandler { //nolint:errorlint // recover() devolve o valor exato
					panic(v)
				}
				logger.Error("handler panic", "requestId", observability.RequestID(r.Context()))
				writeProblem(w, r, http.StatusInternalServerError, ProblemCodeInternalError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// problemFallback converte respostas de erro em texto puro em ProblemDetails.
//
// O ServeMux responde paths desconhecidos (404) e métodos errados (405) com
// http.Error, que escreve text/plain. Em vez de reimplementar o matching do
// mux, este middleware envolve o ResponseWriter e troca essas respostas por um
// corpo ProblemDetails, preservando headers como Allow no 405.
func problemFallback(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(&problemWriter{ResponseWriter: w, r: r}, r)
	})
}

type problemWriter struct {
	http.ResponseWriter
	r           *http.Request
	wroteHeader bool
	replaced    bool // true quando o corpo em texto puro está sendo descartado
}

func (pw *problemWriter) WriteHeader(status int) {
	if pw.wroteHeader {
		return
	}
	pw.wroteHeader = true
	if status >= 400 && strings.HasPrefix(pw.Header().Get("Content-Type"), "text/plain") {
		pw.replaced = true
		pw.Header().Del("X-Content-Type-Options")
		writeProblem(pw.ResponseWriter, pw.r, status, defaultCode(status))
		return
	}
	pw.ResponseWriter.WriteHeader(status)
}

func (pw *problemWriter) Write(b []byte) (int, error) {
	if !pw.wroteHeader {
		pw.WriteHeader(http.StatusOK)
	}
	if pw.replaced {
		// Finge que a escrita deu certo; o texto original não pode vazar.
		return len(b), nil
	}
	return pw.ResponseWriter.Write(b)
}

// Unwrap permite que http.ResponseController alcance o writer real (flush,
// deadlines) através deste wrapper.
func (pw *problemWriter) Unwrap() http.ResponseWriter {
	return pw.ResponseWriter
}
