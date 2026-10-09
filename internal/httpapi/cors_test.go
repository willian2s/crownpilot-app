package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Estes testes especificam o comportamento do middleware CORS. Cada caso passa
// um request por cors(...) envolvendo um handler que registra se foi chamado.
func TestCORS(t *testing.T) {
	t.Parallel()

	allowed := []string{"http://localhost:5173", "https://app.example.com"}

	tests := []struct {
		name    string
		origins []string
		method  string
		headers map[string]string

		wantNextCalled bool
		wantStatus     int    // verificado só quando next não é chamado
		wantAllowed    string // Access-Control-Allow-Origin esperado ("" = ausente)
		wantVaryOrigin bool
		wantMethods    []string // precisam aparecer em Access-Control-Allow-Methods
		wantHeaders    []string // precisam aparecer em Access-Control-Allow-Headers
	}{
		{
			name:           "request without Origin is untouched",
			origins:        allowed,
			method:         http.MethodGet,
			wantNextCalled: true,
		},
		{
			name:           "allowed origin gets its exact origin back",
			origins:        allowed,
			method:         http.MethodGet,
			headers:        map[string]string{"Origin": "https://app.example.com"},
			wantNextCalled: true,
			wantAllowed:    "https://app.example.com",
			wantVaryOrigin: true,
		},
		{
			name:           "disallowed origin still reaches the handler without CORS headers",
			origins:        allowed,
			method:         http.MethodGet,
			headers:        map[string]string{"Origin": "https://evil.example.com"},
			wantNextCalled: true,
			wantVaryOrigin: true,
		},
		{
			name:           "matching is exact, not case-insensitive",
			origins:        allowed,
			method:         http.MethodGet,
			headers:        map[string]string{"Origin": "https://APP.example.com"},
			wantNextCalled: true,
			wantVaryOrigin: true,
		},
		{
			name:           "null origin is never allowed",
			origins:        allowed,
			method:         http.MethodGet,
			headers:        map[string]string{"Origin": "null"},
			wantNextCalled: true,
			wantVaryOrigin: true,
		},
		{
			name:           "empty allowlist allows nothing",
			origins:        nil,
			method:         http.MethodGet,
			headers:        map[string]string{"Origin": "https://app.example.com"},
			wantNextCalled: true,
			wantVaryOrigin: true,
		},
		{
			name:    "allowed preflight is answered with 204",
			origins: allowed,
			method:  http.MethodOptions,
			headers: map[string]string{
				"Origin":                         "http://localhost:5173",
				"Access-Control-Request-Method":  "PUT",
				"Access-Control-Request-Headers": "authorization, Content-Type",
			},
			wantStatus:     http.StatusNoContent,
			wantAllowed:    "http://localhost:5173",
			wantVaryOrigin: true,
			wantMethods:    []string{"GET", "HEAD", "PUT", "DELETE", "OPTIONS"},
			wantHeaders:    []string{"Accept", "Authorization", "Content-Type"},
		},
		{
			name:    "preflight from disallowed origin gets 204 without CORS headers",
			origins: allowed,
			method:  http.MethodOptions,
			headers: map[string]string{
				"Origin":                        "https://evil.example.com",
				"Access-Control-Request-Method": "GET",
			},
			wantStatus:     http.StatusNoContent,
			wantVaryOrigin: true,
		},
		{
			name:    "preflight for a method outside the allowlist is refused",
			origins: allowed,
			method:  http.MethodOptions,
			headers: map[string]string{
				"Origin":                        "https://app.example.com",
				"Access-Control-Request-Method": "PATCH",
			},
			wantStatus:     http.StatusNoContent,
			wantVaryOrigin: true,
		},
		{
			name:    "preflight for a header outside the allowlist is refused",
			origins: allowed,
			method:  http.MethodOptions,
			headers: map[string]string{
				"Origin":                         "https://app.example.com",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "Authorization, X-Debug",
			},
			wantStatus:     http.StatusNoContent,
			wantVaryOrigin: true,
		},
		{
			name:           "OPTIONS without Access-Control-Request-Method is not a preflight",
			origins:        allowed,
			method:         http.MethodOptions,
			headers:        map[string]string{"Origin": "https://app.example.com"},
			wantNextCalled: true,
			wantAllowed:    "https://app.example.com",
			wantVaryOrigin: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusTeapot)
			})

			req := httptest.NewRequestWithContext(t.Context(), tt.method, "/api/v1/bootstrap", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			rec := httptest.NewRecorder()

			cors(tt.origins, next).ServeHTTP(rec, req)

			if nextCalled != tt.wantNextCalled {
				t.Fatalf("next called = %v, want %v", nextCalled, tt.wantNextCalled)
			}
			if !tt.wantNextCalled && rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			h := rec.Header()
			if got := h.Get("Access-Control-Allow-Origin"); got != tt.wantAllowed {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, tt.wantAllowed)
			}
			if got := h.Get("Access-Control-Allow-Credentials"); got != "" {
				t.Errorf("Access-Control-Allow-Credentials = %q, want absent (bearer, no cookies)", got)
			}
			if got := containsToken(h.Values("Vary"), "Origin"); got != tt.wantVaryOrigin {
				t.Errorf("Vary contains Origin = %v, want %v (Vary = %q)", got, tt.wantVaryOrigin, h.Values("Vary"))
			}
			for _, m := range tt.wantMethods {
				if !containsToken(h.Values("Access-Control-Allow-Methods"), m) {
					t.Errorf("Access-Control-Allow-Methods = %q, missing %s", h.Values("Access-Control-Allow-Methods"), m)
				}
			}
			for _, hd := range tt.wantHeaders {
				if !containsToken(h.Values("Access-Control-Allow-Headers"), hd) {
					t.Errorf("Access-Control-Allow-Headers = %q, missing %s", h.Values("Access-Control-Allow-Headers"), hd)
				}
			}
			if tt.wantAllowed == "" && h.Get("Access-Control-Allow-Methods") != "" {
				t.Errorf("Access-Control-Allow-Methods set for a refused request")
			}
		})
	}
}

// containsToken informa se uma lista de header separada por vírgulas contém
// token, ignorando maiúsculas e espaços, inclusive em linhas repetidas.
func containsToken(values []string, token string) bool {
	for _, v := range values {
		for item := range strings.SplitSeq(v, ",") {
			if strings.EqualFold(strings.TrimSpace(item), token) {
				return true
			}
		}
	}
	return false
}
