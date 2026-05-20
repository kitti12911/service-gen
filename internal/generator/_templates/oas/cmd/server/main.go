package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"___LIB_PATH___/lib-monitor/profiling"
	"___LIB_PATH___/lib-monitor/tracing"
	libconfig "___LIB_PATH___/lib-util/v3/config"
	"___LIB_PATH___/lib-util/v3/logger"

	"___MODULE___/internal/config"
	"___MODULE___/internal/server"

	// Register the round_robin balancer so gRPC clients (tracing exporter,
	// future downstream service calls) spread requests across multiple
	// addresses from a headless Kubernetes service.
	_ "google.golang.org/grpc/balancer/roundrobin"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()

	cfg, err := libconfig.Load[config.Config]("config.yml")
	if err != nil {
		slog.ErrorContext(ctx, "failed to load config", "error", err)
		return 1
	}

	if cfg.Service.ShutdownTimeout == 0 {
		cfg.Service.ShutdownTimeout = 10 * time.Second
	}

	logger.NewFromConfig(cfg.Logging, cfg.Service.Name)

	profiler, err := profiling.NewFromConfig(cfg.Service.Name, cfg.Profiling)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init profiling", "error", err)
		return 1
	}
	defer func() {
		if shutdownErr := profiling.Shutdown(profiler); shutdownErr != nil {
			slog.ErrorContext(ctx, "failed to stop profiling", "error", shutdownErr)
		}
	}()

	tp, err := tracing.NewFromConfig(ctx, cfg.Service.Name, cfg.Tracing)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init tracing", "error", err)
		return 1
	}
	defer func() {
		if shutdownErr := tracing.Shutdown(ctx, tp); shutdownErr != nil {
			slog.ErrorContext(ctx, "failed to stop tracing", "error", shutdownErr)
		}
	}()

	srv := server.NewHTTPServer(cfg.Service.Port, cfg.Service.Name)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.Start()
	}()

	slog.InfoContext(ctx, "HTTP server started", "port", cfg.Service.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "HTTP server error", "error", err)
			return 1
		}
	}

	slog.InfoContext(ctx, "shutting down HTTP server")

	shutdownCtx, cancel := context.WithTimeout(ctx, cfg.Service.ShutdownTimeout)
	defer cancel()

	srv.Stop(shutdownCtx)

	slog.InfoContext(ctx, "server stopped")

	return 0
}
