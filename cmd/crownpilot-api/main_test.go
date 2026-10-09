package main

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func lookup(vars map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := vars[key]
		return v, ok
	}
}

func TestRunFailsBeforeListeningOnInvalidConfig(t *testing.T) {
	t.Parallel()
	port := freePort(t)

	var logs bytes.Buffer
	err := run(context.Background(), lookup(map[string]string{
		"CROWNPILOT_ENVIRONMENT":            "Production",
		"PORT":                              strconv.Itoa(port),
		"CROWNPILOT_AUTH_CONTRACT_FIXTURES": "true",
	}), &logs)
	if err == nil || !strings.Contains(err.Error(), "invalid configuration") {
		t.Fatalf("run() error = %v, want invalid configuration", err)
	}
	if logs.Len() != 0 {
		t.Errorf("logs written before failing: %s", logs.String())
	}
	// Nada pode ter feito bind na porta: ainda conseguimos ocupá-la.
	ln, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		t.Fatalf("port %d was bound by a failed start: %v", port, err)
	}
	_ = ln.Close()
}

func TestRunServesAndShutsDownGracefully(t *testing.T) {
	t.Parallel()
	port := freePort(t)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- run(ctx, lookup(map[string]string{
			"CROWNPILOT_ENVIRONMENT": "Production",
			"PORT":                   strconv.Itoa(port),
		}), &bytes.Buffer{})
	}()

	url := "http://127.0.0.1:" + strconv.Itoa(port) + "/health/live"
	deadline := time.Now().Add(5 * time.Second)
	for {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				break
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("server did not become live: %v", err)
		}
		time.Sleep(20 * time.Millisecond)
	}

	cancel() // o que o SIGTERM faz em produção
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("run() after shutdown = %v, want nil", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("run() did not return after context cancellation")
	}
}

// freePort pede ao SO uma porta livre. Outro processo pode ocupá-la antes de o
// teste fazer o bind de novo; essa corrida é aceitável num teste local.
func freePort(t *testing.T) int {
	t.Helper()
	ln, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()
	return ln.Addr().(*net.TCPAddr).Port
}
