package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

const sentinel = "PRIVATE-KEY-SENTINEL"

// TestSecretIsNeverPrinted cobre as formas comuns de um secret vazar por
// acidente: fmt em todos os verbos, a Config inteira impressa, slog e JSON.
func TestSecretIsNeverPrinted(t *testing.T) {
	t.Parallel()

	s := NewSecret(sentinel)
	cfg := Config{Firebase: Firebase{PrivateKey: s}}

	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	logger.Info("secret", "value", s, "config", cfg)

	encoded, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	outputs := map[string]string{
		"%v":      fmt.Sprintf("%v", s),
		"%s":      fmt.Sprintf("%s", s), //nolint:staticcheck // o teste é justamente o verbo %s
		"%+v":     fmt.Sprintf("%+v", s),
		"%#v":     fmt.Sprintf("%#v", s),
		"config":  fmt.Sprintf("%+v", cfg),
		"config#": fmt.Sprintf("%#v", cfg),
		"slog":    logs.String(),
		"json":    string(encoded),
		"Println": fmt.Sprintln(s),
		"Sprint":  fmt.Sprint(s),
		"Errorf":  fmt.Errorf("wrapped: %v", s).Error(),
	}
	for name, out := range outputs {
		if strings.Contains(out, sentinel) {
			t.Errorf("%s leaks the secret: %s", name, out)
		}
		if !strings.Contains(out, redacted) {
			t.Errorf("%s does not show %s: %s", name, redacted, out)
		}
	}

	if s.Reveal() != sentinel {
		t.Errorf("Reveal() = %q, want the original value", s.Reveal())
	}
	if !NewSecret("").IsZero() || s.IsZero() {
		t.Error("IsZero() is wrong")
	}
}
