package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"___MODULE___/internal/config"
	"___MODULE___/internal/database"
	"___MODULE___/internal/server"

	"___LIB_PATH___/lib-monitor/profiling"
	"___LIB_PATH___/lib-monitor/tracing"
	libconfig "___LIB_PATH___/lib-util/v3/config"
	"___LIB_PATH___/lib-util/v3/logger"

	"github.com/dromara/carbon/v2"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()

	carbon.SetDefault(carbon.Default{
		Layout:       carbon.RFC3339Format,
		Timezone:     carbon.Bangkok,
		WeekStartsAt: carbon.Sunday,
		Locale:       "en",
	})

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

	db, err := database.New(ctx, cfg)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init database", "error", err)
		return 1
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			slog.ErrorContext(ctx, "failed to close database", "error", closeErr)
		}
	}()

	srv, err := server.NewGRPCServer(ctx, cfg.Service.Port)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create gRPC server", "error", err)
		return 1
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.Start()
	}()

	slog.InfoContext(ctx, "gRPC server started", "port", cfg.Service.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case err := <-serverErr:
		slog.ErrorContext(ctx, "gRPC server error", "error", err)
		return 1
	}

	slog.InfoContext(ctx, "shutting down gRPC server")

	shutdownCtx, cancel := context.WithTimeout(ctx, cfg.Service.ShutdownTimeout)
	defer cancel()

	srv.Stop(shutdownCtx)

	slog.InfoContext(ctx, "server stopped")

	return 0
}
