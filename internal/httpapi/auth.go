package httpapi

import (
	"context"
	"net/http"
	"slices"
	"strings"
)

// Authentication e authorization são etapas separadas, como na spec 002:
//
//   - authenticate (middleware, todo request) responde "quem está chamando?" e
//     só guarda um subject no context quando o bearer é válido;
//   - requirePermission (wrapper, por rota) responde "este subject pode fazer
//     isto?" e escreve 401 (sem subject) ou 403 (sem permissão).
//
// A task 002-13 só tem a fixture de contrato Local abaixo. A task 002-14 a
// substitui pela verificação Firebase e move o subject para internal/identity.

// Tokens da fixture de contrato Local. Provam o pipeline 401/403/200 e não são
// validação Firebase; config.Validate os rejeita fora de Local.
const (
	contractAuthenticatedToken = "contract-authenticated-token"
	contractAuthorizedToken    = "contract-authorized-token"
	contractSubject            = "contract-firebase-uid"
	permissionBootstrapRead    = "bootstrap:read"
)

type subject struct {
	// externalID é o subject verificado do provider. Nunca é identidade de
	// domínio e nunca vem de body, query ou headers arbitrários.
	externalID  string
	permissions []string
}

func (s subject) has(permission string) bool {
	return slices.Contains(s.permissions, permission)
}

type subjectKey struct{}

func subjectFrom(ctx context.Context) (subject, bool) {
	s, ok := ctx.Value(subjectKey{}).(subject)
	return s, ok
}

// authenticate guarda um subject para bearers de fixture válidos. Com fixtures
// desligadas nenhum bearer é aceito, então rotas protegidas falham fechado com
// 401.
func authenticate(contractFixtures bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r)
		if ok && contractFixtures {
			if s, valid := contractFixtureSubject(token); valid {
				r = r.WithContext(context.WithValue(r.Context(), subjectKey{}, s))
			}
		}
		next.ServeHTTP(w, r)
	})
}

// requirePermission envolve uma única rota. Declará-lo no registro mantém a
// segurança de cada rota visível num só lugar (veja as rotas em server.go).
func requirePermission(permission string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, ok := subjectFrom(r.Context())
		if !ok {
			w.Header().Set("WWW-Authenticate", "Bearer")
			writeProblem(w, r, http.StatusUnauthorized, ProblemCodeAuthenticationRequired)
			return
		}
		if !s.has(permission) {
			writeProblem(w, r, http.StatusForbidden, ProblemCodeAuthorizationForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// bearerToken extrai o token de exatamente um header "Authorization: Bearer x".
// Headers repetidos são rejeitados em vez de escolher um deles.
func bearerToken(r *http.Request) (string, bool) {
	values := r.Header.Values("Authorization")
	if len(values) != 1 {
		return "", false
	}
	scheme, token, found := strings.Cut(values[0], " ")
	token = strings.TrimSpace(token)
	if !found || !strings.EqualFold(scheme, "Bearer") || token == "" {
		return "", false
	}
	return token, true
}

func contractFixtureSubject(token string) (subject, bool) {
	switch token {
	case contractAuthenticatedToken:
		return subject{externalID: contractSubject}, true
	case contractAuthorizedToken:
		return subject{externalID: contractSubject, permissions: []string{permissionBootstrapRead}}, true
	default:
		return subject{}, false
	}
}
