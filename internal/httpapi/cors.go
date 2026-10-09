package httpapi

import (
	"net/http"
	"slices"
	"strings"
)

var (
	corsMethods = []string{"GET", "HEAD", "PUT", "DELETE", "OPTIONS"}
	corsHeaders = []string{"Accept", "Authorization", "Content-Type"}
)

// cors aplica a política CORS de origens exatas da spec 002. Preflights são
// respondidos com 204 sem chamar next; requests normais sempre seguem para next.
func cors(allowedOrigins []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Origem não permitida ainda chega ao handler: quem bloqueia a leitura é o browser.
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Add("Vary", "Origin")

		originAllowed := slices.Contains(allowedOrigins, origin)

		// Verifica se é uma requisição preflight.
		requestMethod := r.Header.Get("Access-Control-Request-Method")

		if r.Method == http.MethodOptions && requestMethod != "" {
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")

			// Valida o método solicitado.
			methodAllowed := slices.Contains(corsMethods, requestMethod)

			// Valida os headers solicitados.
			headersAllowed := true
			requestHeaders := r.Header.Get("Access-Control-Request-Headers")

			// SplitSeq("") produz um item vazio; sem esta guarda todo preflight sem headers seria recusado.
			if requestHeaders != "" {
				for header := range strings.SplitSeq(requestHeaders, ",") {
					header = strings.TrimSpace(header)

					if !slices.ContainsFunc(corsHeaders, func(a string) bool { return strings.EqualFold(a, header) }) {
						headersAllowed = false
						break
					}
				}
			}

			// Só permite se tudo estiver autorizado.
			if originAllowed && methodAllowed && headersAllowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set(
					"Access-Control-Allow-Methods",
					strings.Join(corsMethods, ", "),
				)
				w.Header().Set(
					"Access-Control-Allow-Headers",
					strings.Join(corsHeaders, ", "),
				)
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Requisições normais seguem para o handler.
		if originAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		next.ServeHTTP(w, r)
	})
}
