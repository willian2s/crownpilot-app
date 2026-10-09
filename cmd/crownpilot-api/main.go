// Command crownpilot-api é o composition root do backend CrownPilot.
//
// É o único pacote autorizado a ligar configuração, HTTP e, em tasks
// posteriores, adapters concretos de providers. Módulos de negócio em
// internal/ nunca o importam.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/willian2s/crownpilot-app/internal/config"
	"github.com/willian2s/crownpilot-app/internal/httpapi"
	"github.com/willian2s/crownpilot-app/internal/observability"
)

// shutdownTimeout limita quanto tempo os requests em andamento têm para
// terminar após o SIGTERM. Plataformas de hosting costumam matar o container
// alguns segundos depois.
const shutdownTimeout = 10 * time.Second

func main() {
	// ctx é cancelado no Ctrl+C (SIGINT) ou no SIGTERM enviado pela plataforma.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, os.LookupEnv, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		stop() // os.Exit pula os defer, então libera o sinal agora.
		os.Exit(1)
	}
}

// run mantém o main sem lógica: cada etapa do startup devolve um erro em vez
// de encerrar, então o processo tem um único caminho de saída e run é testável
// com seu próprio context, ambiente e saída de log.
func run(ctx context.Context, lookup config.LookupFunc, logOutput io.Writer) error {
	// Fail-closed: a configuração é validada antes de existir qualquer listener.
	cfg, err := config.Load(lookup)
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	logger := observability.NewLogger(logOutput)

	// Fazer o bind antes, separado do Serve, faz um conflito de porta aparecer
	// como erro de startup antes de o processo anunciar que está escutando.
	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp", cfg.Addr())
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	server := &http.Server{
		Handler: httpapi.NewHandler(cfg, logger),
		// Timeouts protegem contra clientes lentos segurando conexões abertas.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Serve bloqueia, então roda na sua própria goroutine e reporta por um
	// channel com buffer; o buffer deixa a goroutine terminar mesmo que ninguém
	// leia.
	serveErr := make(chan error, 1)
	go func() { serveErr <- server.Serve(listener) }()

	// Só valores não secretos são logados.
	logger.Info("listening", "environment", cfg.Environment, "addr", listener.Addr().String())

	select {
	case err := <-serveErr:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
	}

	// O ctx pai já foi cancelado, então o shutdown recebe um prazo novo.
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
	defer cancel()
	logger.Info("shutting down")
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	// Depois do Shutdown, Serve devolve http.ErrServerClosed, que é esperado.
	if err := <-serveErr; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve: %w", err)
	}
	return nil
}
