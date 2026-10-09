package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/willian2s/crownpilot-app/api/openapi"
	"github.com/willian2s/crownpilot-app/internal/config"
)

func TestOpenAPIIsServedUnchanged(t *testing.T) {
	t.Parallel()
	h := NewHandler(localConfig(), slog.New(slog.DiscardHandler))

	rec := serve(t, h, http.MethodGet, "/openapi/v1.json", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	if rec.Body.String() != openapi.V1() {
		t.Error("served document differs from the embedded api/openapi/v1.json")
	}
}

func TestDocumentationExposure(t *testing.T) {
	t.Parallel()

	withDocs := func(env config.Environment, json, ui bool) config.Config {
		cfg := localConfig()
		cfg.Environment = env
		cfg.Auth.ContractFixtures = env == config.Local
		cfg.OpenAPI = config.OpenAPI{ExposeJSON: json, ExposeUI: ui}
		return cfg
	}

	tests := []struct {
		name     string
		cfg      config.Config
		wantJSON int
		wantUI   int
	}{
		{"Local exposes JSON and UI", withDocs(config.Local, true, true), http.StatusOK, http.StatusOK},
		{"Staging exposes JSON and UI", withDocs(config.Staging, true, true), http.StatusOK, http.StatusOK},
		{"Staging may expose JSON only", withDocs(config.Staging, true, false), http.StatusOK, http.StatusNotFound},
		{"Preview hides both", withDocs(config.Preview, false, false), http.StatusNotFound, http.StatusNotFound},
		{"Production hides both", withDocs(config.Production, false, false), http.StatusNotFound, http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewHandler(tt.cfg, slog.New(slog.DiscardHandler))

			if got := serve(t, h, http.MethodGet, "/openapi/v1.json", nil).Code; got != tt.wantJSON {
				t.Errorf("/openapi/v1.json = %d, want %d", got, tt.wantJSON)
			}
			rec := serve(t, h, http.MethodGet, "/docs", nil)
			if rec.Code != tt.wantUI {
				t.Fatalf("/docs = %d, want %d", rec.Code, tt.wantUI)
			}
			if rec.Code == http.StatusNotFound {
				decodeProblem(t, rec)
			}
		})
	}
}

func TestDocsPage(t *testing.T) {
	t.Parallel()
	rec := serve(t, NewHandler(localConfig(), slog.New(slog.DiscardHandler)), http.MethodGet, "/docs", nil)

	body := rec.Body.String()
	for _, want := range []string{`url: "/openapi/v1.json"`, swaggerUICSSSRI, swaggerUIJSSRI} {
		if !strings.Contains(body, want) {
			t.Errorf("/docs page does not contain %q", want)
		}
	}
	csp := rec.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "'sha256-") || strings.Contains(csp, "script-src 'unsafe-inline'") {
		t.Errorf("Content-Security-Policy = %q, want inline script allowed only by hash", csp)
	}
}

// TestRoutesMatchContract liga a spec ao ServeMux. Como o oapi-codegen gera só
// tipos, é este teste que garante que cada operação documentada existe no
// servidor, e que operações marcadas como planejadas ainda não existem.
func TestRoutesMatchContract(t *testing.T) {
	t.Parallel()

	var doc struct {
		OpenAPI string                                `json:"openapi"`
		Paths   map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal([]byte(openapi.V1()), &doc); err != nil {
		t.Fatalf("api/openapi/v1.json is not valid JSON: %v", err)
	}
	if !strings.HasPrefix(doc.OpenAPI, "3.") {
		t.Fatalf("openapi = %q, want 3.x", doc.OpenAPI)
	}
	if len(doc.Paths) == 0 {
		t.Fatal("spec has no paths")
	}

	h := NewHandler(localConfig(), slog.New(slog.DiscardHandler))
	for path, operations := range doc.Paths {
		for method, raw := range operations {
			var op struct {
				Planned bool `json:"x-crownpilot-planned"`
			}
			if err := json.Unmarshal(raw, &op); err != nil {
				t.Fatalf("%s %s: %v", method, path, err)
			}

			name := strings.ToUpper(method) + " " + path
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				req := httptest.NewRequestWithContext(t.Context(), strings.ToUpper(method), path, nil)
				req.Header.Set("Authorization", "Bearer "+contractAuthorizedToken)
				rec := httptest.NewRecorder()
				h.ServeHTTP(rec, req)

				routed := rec.Code != http.StatusNotFound && rec.Code != http.StatusMethodNotAllowed
				switch {
				case op.Planned && routed:
					t.Errorf("planned operation is already served (HTTP %d); remove x-crownpilot-planned", rec.Code)
				case !op.Planned && !routed:
					t.Errorf("documented operation is not served (HTTP %d)", rec.Code)
				}
			})
		}
	}
}
