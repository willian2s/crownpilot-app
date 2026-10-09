// Package httpapi é o adapter HTTP: rotas, middlewares, DTOs e ProblemDetails.
// Traduz HTTP em chamadas da aplicação e vice-versa; regras de negócio nunca
// moram aqui.
package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/willian2s/crownpilot-app/internal/config"
)

// NewHandler monta o handler HTTP completo para cfg.
//
// Os middlewares rodam do mais externo para o mais interno. A ordem importa:
//
//	withRequestID   toda etapa seguinte, inclusive erros, consegue ler o ID
//	recoverPanic    um panic em qualquer ponto abaixo vira 500 ProblemDetails
//	problemFallback os 404/405 em texto puro do ServeMux viram ProblemDetails
//	cors            preflights são respondidos antes da authentication, porque
//	                browsers nunca enviam Authorization num preflight
//	authenticate    guarda o subject para a authorization por rota
//	mux             roteamento por método e path
func NewHandler(cfg config.Config, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	// Liveness é só do processo: nunca pode chamar provider ou banco, senão uma
	// queda em outro lugar faria a plataforma reiniciar a API.
	mux.HandleFunc("GET /health/live", handleLive)

	// Rotas não registradas caem no 404 do mux, que vira ProblemDetails. Assim
	// Preview e Production respondem 404 sem nenhum código condicional extra.
	if cfg.OpenAPI.ExposeJSON {
		mux.HandleFunc("GET /openapi/v1.json", handleOpenAPI)
	}
	if cfg.OpenAPI.ExposeUI {
		mux.HandleFunc("GET /docs", handleDocs)
	}

	mux.Handle("GET /api/v1/bootstrap",
		requirePermission(permissionBootstrapRead, http.HandlerFunc(handleBootstrap)))

	var h http.Handler = mux
	h = authenticate(cfg.Auth.ContractFixtures, h)
	h = cors(cfg.CORS.AllowedOrigins, h)
	h = problemFallback(h)
	h = recoverPanic(logger, h)
	h = withRequestID(h)
	return h
}

func handleLive(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("Healthy"))
}

func handleBootstrap(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, BootstrapStatus{Status: "ready", Stage: "bootstrap"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
