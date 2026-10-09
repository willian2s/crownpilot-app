package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/willian2s/crownpilot-app/internal/config"
)

func localConfig() config.Config {
	return config.Config{
		Environment: config.Local,
		Port:        config.DefaultPort,
		CORS:        config.CORS{AllowedOrigins: []string{"http://localhost:5173"}},
		OpenAPI:     config.OpenAPI{ExposeJSON: true, ExposeUI: true},
		Auth:        config.Auth{ContractFixtures: true},
	}
}

func serve(t *testing.T, h http.Handler, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequestWithContext(t.Context(), method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// decodeProblem verifica uma resposta ProblemDetails e devolve o corpo.
func decodeProblem(t *testing.T, rec *httptest.ResponseRecorder) ProblemDetails {
	t.Helper()
	if ct := rec.Header().Get("Content-Type"); ct != "application/problem+json" {
		t.Fatalf("Content-Type = %q, want application/problem+json (body %q)", ct, rec.Body.String())
	}
	var p ProblemDetails
	dec := json.NewDecoder(rec.Body)
	dec.DisallowUnknownFields() // detail/instance ou qualquer outro campo não pode vazar
	if err := dec.Decode(&p); err != nil {
		t.Fatalf("decode ProblemDetails: %v", err)
	}
	if p.Type != problemType || p.Title == "" || p.TraceId == "" || p.Status != rec.Code {
		t.Errorf("incomplete ProblemDetails: %+v (HTTP %d)", p, rec.Code)
	}
	return p
}

func TestHealthLive(t *testing.T) {
	t.Parallel()
	h := NewHandler(localConfig(), slog.New(slog.DiscardHandler))

	for _, method := range []string{http.MethodGet, http.MethodHead} {
		rec := serve(t, h, method, "/health/live", nil)
		if rec.Code != http.StatusOK {
			t.Errorf("%s /health/live = %d, want 200", method, rec.Code)
		}
	}
	if body := serve(t, h, http.MethodGet, "/health/live", nil).Body.String(); body != "Healthy" {
		t.Errorf("body = %q, want Healthy", body)
	}
}

func TestBootstrapAuthentication(t *testing.T) {
	t.Parallel()

	nonLocal := localConfig()
	nonLocal.Environment = config.Staging
	nonLocal.Auth.ContractFixtures = false

	tests := []struct {
		name       string
		cfg        config.Config
		auth       []string // valores do header Authorization
		wantStatus int
		wantCode   ProblemCode
	}{
		{"missing bearer", localConfig(), nil, http.StatusUnauthorized, ProblemCodeAuthenticationRequired},
		{"unknown token", localConfig(), []string{"Bearer nope"}, http.StatusUnauthorized, ProblemCodeAuthenticationRequired},
		{"wrong scheme", localConfig(), []string{"Basic " + contractAuthorizedToken}, http.StatusUnauthorized, ProblemCodeAuthenticationRequired},
		{"empty token", localConfig(), []string{"Bearer "}, http.StatusUnauthorized, ProblemCodeAuthenticationRequired},
		{"repeated header", localConfig(), []string{"Bearer " + contractAuthorizedToken, "Bearer " + contractAuthorizedToken}, http.StatusUnauthorized, ProblemCodeAuthenticationRequired},
		{"authenticated without permission", localConfig(), []string{"Bearer " + contractAuthenticatedToken}, http.StatusForbidden, ProblemCodeAuthorizationForbidden},
		{"authorized", localConfig(), []string{"Bearer " + contractAuthorizedToken}, http.StatusOK, ""},
		{"scheme is case-insensitive", localConfig(), []string{"bearer " + contractAuthorizedToken}, http.StatusOK, ""},
		{"fixtures disabled fail closed", nonLocal, []string{"Bearer " + contractAuthorizedToken}, http.StatusUnauthorized, ProblemCodeAuthenticationRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewHandler(tt.cfg, slog.New(slog.DiscardHandler))
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/bootstrap", nil)
			for _, v := range tt.auth {
				req.Header.Add("Authorization", v)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body %q)", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantStatus == http.StatusOK {
				var got BootstrapStatus
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("decode: %v", err)
				}
				if got != (BootstrapStatus{Status: "ready", Stage: "bootstrap"}) {
					t.Errorf("body = %+v", got)
				}
				return
			}
			if p := decodeProblem(t, rec); p.Code != tt.wantCode {
				t.Errorf("code = %q, want %q", p.Code, tt.wantCode)
			}
			if tt.wantStatus == http.StatusUnauthorized && rec.Header().Get("WWW-Authenticate") != "Bearer" {
				t.Errorf("WWW-Authenticate = %q, want Bearer", rec.Header().Get("WWW-Authenticate"))
			}
		})
	}
}

func TestMuxErrorsBecomeProblemDetails(t *testing.T) {
	t.Parallel()
	h := NewHandler(localConfig(), slog.New(slog.DiscardHandler))

	rec := serve(t, h, http.MethodGet, "/nope", nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if p := decodeProblem(t, rec); p.Code != ProblemCodeResourceNotFound {
		t.Errorf("code = %q, want %q", p.Code, ProblemCodeResourceNotFound)
	}

	rec = serve(t, h, http.MethodPost, "/health/live", nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
	if p := decodeProblem(t, rec); p.Code != ProblemCodeMethodNotAllowed {
		t.Errorf("code = %q, want %q", p.Code, ProblemCodeMethodNotAllowed)
	}
	if allow := rec.Header().Get("Allow"); !strings.Contains(allow, http.MethodGet) {
		t.Errorf("Allow = %q, want it to list GET", allow)
	}
}

func TestTraceIDIsGeneratedNotTrusted(t *testing.T) {
	t.Parallel()
	h := NewHandler(localConfig(), slog.New(slog.DiscardHandler))

	first := decodeProblem(t, serve(t, h, http.MethodGet, "/nope", map[string]string{"X-Request-ID": "forged"}))
	second := decodeProblem(t, serve(t, h, http.MethodGet, "/nope", nil))
	if first.TraceId == "forged" {
		t.Error("traceId was taken from the request")
	}
	if first.TraceId == second.TraceId {
		t.Errorf("traceId repeated across requests: %q", first.TraceId)
	}
}

func TestPanicBecomesGenericProblem(t *testing.T) {
	t.Parallel()
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	h := withRequestID(recoverPanic(logger, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("secret-token-value")
	})))

	rec := serve(t, h, http.MethodGet, "/", nil)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
	if p := decodeProblem(t, rec); p.Code != ProblemCodeInternalError {
		t.Errorf("code = %q, want %q", p.Code, ProblemCodeInternalError)
	}
	for name, out := range map[string]string{"response": rec.Body.String(), "log": logs.String()} {
		if strings.Contains(out, "secret-token-value") {
			t.Errorf("panic value leaked into %s: %s", name, out)
		}
	}
}

func TestPreflightIsNotBlockedByAuthentication(t *testing.T) {
	t.Parallel()
	h := NewHandler(localConfig(), slog.New(slog.DiscardHandler))

	rec := serve(t, h, http.MethodOptions, "/api/v1/bootstrap", map[string]string{
		"Origin":                         "http://localhost:5173",
		"Access-Control-Request-Method":  "GET",
		"Access-Control-Request-Headers": "authorization",
	})
	if rec.Code != http.StatusNoContent {
		body, _ := io.ReadAll(rec.Body)
		t.Fatalf("preflight status = %d, want 204 (body %q)", rec.Code, body)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Errorf("Access-Control-Allow-Origin = %q", got)
	}
}
