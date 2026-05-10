package generator

var commonTemplates = map[string]string{
	".gitignore": `# test coverage
coverage.out

# generated code
gen/

# local config
config.yml
config.yaml

# air
tmp/
`,
	".markdownlint-cli2.jsonc": `{
    "config": {
        "MD013": false
    },
    "globs": ["**/*.{md,markdown}"],
    "ignores": [".github/CODEOWNERS", "CHANGELOG.md"]
}
`,
	".prettierrc.json": `{
    "tabWidth": 4,
    "trailingComma": "none",
    "useTabs": false
}
`,
	".prettierignore": `CHANGELOG.md
.codex/
.cursor/
gen/
tmp/
`,
	".golangci.yml": `version: "2"

run:
    timeout: 5m

linters:
    enable:
        - bodyclose
        - cyclop
        - errcheck
        - errorlint
        - exhaustive
        - gocritic
        - godot
        - gosec
        - govet
        - ineffassign
        - misspell
        - nilerr
        - noctx
        - prealloc
        - revive
        - sqlclosecheck
        - staticcheck
        - unconvert
        - unparam
        - unused
        - wrapcheck

    settings:
        cyclop:
            max-complexity: 15

        exhaustive:
            default-signifies-exhaustive: true

        gocritic:
            enabled-tags:
                - diagnostic
                - performance
                - style
            disabled-checks:
                - hugeParam
                - rangeValCopy

        gosec:
            excludes:
                - G404

        govet:
            enable-all: true
            disable:
                - fieldalignment

        misspell:
            locale: US

        revive:
            rules:
                - name: exported
                  arguments:
                      - disableStutteringCheck
                - name: var-naming
                - name: blank-imports
                - name: context-as-argument
                - name: error-return
                - name: error-strings
                - name: increment-decrement
                - name: range
                - name: receiver-naming
                - name: unused-parameter
                  disabled: true

        unparam:
            check-exported: false

        wrapcheck:
            ignore-sigs:
                - .Errorf(
                - errors.New(
                - errors.Unwrap(
                - errors.Join(
                - .Wrap(
                - .Wrapf(
                - .WithMessage(
                - .WithMessagef(
                - .WithStack(

    exclusions:
        generated: lax
        presets:
            - comments
            - common-false-positives
            - legacy
            - std-error-handling
        rules:
            - path: _test\.go$
              linters:
                  - errcheck
                  - godot
                  - wrapcheck
            - path: ^internal/(server|config|database)/.*\.go$
              linters:
                  - wrapcheck

formatters:
    enable:
        - gofmt
        - goimports
    settings:
        goimports:
            local-prefixes:
                - {{ .ModulePath }}

issues:
    max-issues-per-linter: 0
    max-same-issues: 0
`,
	".vscode/settings.json": `{
    "[markdown]": {
        "editor.defaultFormatter": "esbenp.prettier-vscode",
        "editor.tabSize": 4,
        "editor.insertSpaces": true,
        "editor.detectIndentation": false
    },
    "[yaml]": {
        "editor.defaultFormatter": "esbenp.prettier-vscode",
        "editor.tabSize": 4,
        "editor.insertSpaces": true,
        "editor.detectIndentation": false
    },
    "[json]": {
        "editor.defaultFormatter": "esbenp.prettier-vscode",
        "editor.tabSize": 4,
        "editor.insertSpaces": true,
        "editor.detectIndentation": false
    },
    "[jsonc]": {
        "editor.defaultFormatter": "esbenp.prettier-vscode",
        "editor.tabSize": 4,
        "editor.insertSpaces": true,
        "editor.detectIndentation": false
    }
}
`,
	".vscode/extensions.json": `{
    "recommendations": [
        "davidanson.vscode-markdownlint",
        "esbenp.prettier-vscode",
        "github.vscode-github-actions",
        "golang.go",
        "ms-azuretools.vscode-containers",
        "ms-azuretools.vscode-docker",
        "ms-vscode.makefile-tools"
    ]
}
`,
	".github/CODEOWNERS": `/.github/ @kitti12911
/.vscode/ @kitti12911
/.golangci.yml @kitti12911
/.markdownlint-cli2.jsonc @kitti12911
/.prettierrc.json @kitti12911
/Makefile @kitti12911
`,
	".github/renovate.json": `{
    "$schema": "https://docs.renovatebot.com/renovate-schema.json",
    "extends": ["config:recommended"],
    "timezone": "Asia/Bangkok",
    "schedule": ["* 0-4 1 * *"],
    "updateNotScheduled": false,
    "enabledManagers": ["gomod", "github-actions"],
    "postUpdateOptions": ["gomodTidy"],
    "reviewersFromCodeOwners": true,
    "assigneesFromCodeOwners": true,
    "assignAutomerge": true
}
`,
	".github/workflows/markdownlint.yml": `name: Markdownlint

on:
    push:
        branches:
            - main
        paths:
            - ".github/workflows/markdownlint.yml"
            - ".markdownlint-cli2.*"
            - ".markdownlint.*"
            - "**/*.md"
            - "**/*.markdown"
    pull_request:
        paths:
            - ".github/workflows/markdownlint.yml"
            - ".markdownlint-cli2.*"
            - ".markdownlint.*"
            - "**/*.md"
            - "**/*.markdown"

permissions:
    contents: read

jobs:
    markdownlint:
        name: Markdownlint
        runs-on: ubuntu-latest
        steps:
            # actions/checkout v6.0.2
            - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd

            # DavidAnson/markdownlint-cli2-action v23.1.0
            - uses: DavidAnson/markdownlint-cli2-action@6b51ade7a9e4a75a7ad929842dd298a3804ebe8b
`,
	".github/workflows/go-ci.yml": `name: Go CI

on:
    push:
        branches:
            - develop
            - uat
            - main
        paths:
            - ".github/workflows/go-ci.yml"
            - ".github/workflows/markdownlint.yml"
            - ".golangci.yml"
            - "Dockerfile"
            - "Makefile"
            - "config.example.yml"
            - "go.mod"
            - "go.sum"
            - "**/*.go"
            - "**/*.yml"
            - "**/*.yaml"
            - "!README.md"
    pull_request:
        paths:
            - ".github/workflows/go-ci.yml"
            - ".github/workflows/markdownlint.yml"
            - ".golangci.yml"
            - "Dockerfile"
            - "Makefile"
            - "config.example.yml"
            - "go.mod"
            - "go.sum"
            - "**/*.go"
            - "**/*.yml"
            - "**/*.yaml"
            - "!README.md"

permissions:
    contents: read

jobs:
    lint:
        name: Lint
        runs-on: ubuntu-latest
        steps:
            # actions/checkout v6.0.2
            - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd

            # actions/setup-go v6.4.0
            - uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c
              with:
                  go-version-file: go.mod
                  cache: false

            - name: Go vet
              run: |
                  go mod tidy
                  go vet ./...

            # golangci/golangci-lint-action v9.2.0
            - uses: golangci/golangci-lint-action@1e7e51e771db61008b38414a730f564565cf7c20
              with:
                  version: v2.12.1
                  args: --timeout=5m

    test:
        name: Test
        runs-on: ubuntu-latest
        steps:
            # actions/checkout v6.0.2
            - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd

            # actions/setup-go v6.4.0
            - uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c
              with:
                  go-version-file: go.mod
                  cache: false

            - name: Test with race detector and coverage
              run: go test -race -coverprofile=coverage.out -covermode=atomic ./...

            - name: Show coverage
              run: |
                  coverage_total="$(go tool cover -func=coverage.out | tail -n 1)"
                  echo "${coverage_total}"
                  {
                      echo "## Go coverage"
                      echo
                      echo '` + "```text" + `'
                      echo "${coverage_total}"
                      echo '` + "```" + `'
                  } >> "${GITHUB_STEP_SUMMARY}"
	`,
	"Dockerfile": `FROM zot.kittiaccess.work/kitti12911/image-toolchain@sha256:47355f96a059465947c38aa956da1c4502c11d1e8f53eb2c8b3980ba58983d42 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
	go build -trimpath -ldflags="-s -w" -o /out/{{ .BinaryName }} ./cmd/server

FROM alpine:3.22@sha256:310c62b5e7ca5b08167e4384c68db0fd2905dd9c7493756d356e893909057601

RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S app \
	&& adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/{{ .BinaryName }} /app/{{ .BinaryName }}
COPY --chown=app:app config.example.yml /app/config.yml

USER app

EXPOSE {{ .GRPCPort }}

ENTRYPOINT ["/app/{{ .BinaryName }}"]
`,
	"Makefile": `# ____________________ Go Command ____________________
air:
	air

tidy:
	go mod tidy

run:
	go run ./cmd/server

lint: vet golangci-lint markdownlint

vet:
	go vet ./...

golangci-lint:
	golangci-lint run --timeout=5m

markdownlint:
	markdownlint-cli2

fmt:
	go fmt ./...

pretty:
	prettier --write "**/*.{md,markdown,yml,yaml,json,jsonc}"

format: fmt pretty

test:
	env CGO_ENABLED=1 go test --race -v ./...

cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

fix:
	go fix ./...
`,
	"README.md": `# {{ .Name }}

Internal gRPC service bootstrap generated by service-gen.

## Local Development

Copy the example config when you want a local file:

` + "```sh" + `
cp config.example.yml config.yml
` + "```" + `

Run the service:

` + "```sh" + `
make tidy
make run
` + "```" + `

Run verification:

` + "```sh" + `
make format
make lint
make test
` + "```" + `

## Configuration

The starter code reads configuration from environment variables first. The
` + "`config.example.yml`" + ` file documents the values expected by deployment
systems and humans.
`,
	"config.example.yml": `service:
    name: {{ .Name }}
    port: 50051
    shutdown_timeout: 10s

logging:
    level: warn
    add_source: false
    include_trace_id: true

tracing:
    enabled: true
    endpoint: alloy.lan:4317
    protocol: grpc
    insecure: false
    sample_ratio: 1.0

profiling:
    enabled: true
    server_address: https://pyroscope-ingest.lan
    namespace: {{ .Name }}
    basic_auth_user: ""
    basic_auth_password: ""
    tenant_id: ""

database:
    host: postgres.lan
    port: "5432"
    user: example_user
    password: example_password
    database: example
    run_migrations: false
    run_seeders: false
    pool:
        max_conns: 20
        min_conns: 5
        max_conn_life_time: 1h
        max_conn_idle_time: 6h
`,
	"internal/config/config.go": `package config

import (
	"time"

	"github.com/kitti12911/lib-monitor/profiling"
	"github.com/kitti12911/lib-monitor/tracing"
	liborm "github.com/kitti12911/lib-orm/v2"
	"github.com/kitti12911/lib-util/v3/logger"
)

// Config contains runtime settings for the service.
type Config struct {
	Service   Service          ` + "`mapstructure:\"service\" validate:\"required\"`" + `
	Logging   logger.Config    ` + "`mapstructure:\"logging\"`" + `
	Tracing   tracing.Config   ` + "`mapstructure:\"tracing\"`" + `
	Profiling profiling.Config ` + "`mapstructure:\"profiling\"`" + `
	Database  liborm.Config    ` + "`mapstructure:\"database\" validate:\"required\"`" + `
}

// Service contains service identity and listener settings.
type Service struct {
	Name            string        ` + "`mapstructure:\"name\"             env:\"SERVICE_NAME\"      validate:\"required\"`" + `
	Port            int           ` + "`mapstructure:\"port\"             env:\"PORT\"              validate:\"required,gte=1,lte=65535\"`" + `
	ShutdownTimeout time.Duration ` + "`mapstructure:\"shutdown_timeout\" env:\"SHUTDOWN_TIMEOUT\"`" + `
}
`,
	"internal/database/database.go": `package database

import (
	"context"

	"{{ .ModulePath }}/internal/config"

	orm "github.com/kitti12911/lib-orm/v2"
)

// New creates the service-owned database connection.
func New(ctx context.Context, cfg *config.Config) (*orm.DB, error) {
	db, err := orm.New(
		ctx,
		cfg.Database,
		orm.WithApplicationName(cfg.Service.Name),
		orm.WithModels(models()...),
		orm.WithTracing(cfg.Tracing.Enabled),
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
`,
	"internal/database/model.go": `package database

func models() []any {
	return nil
}
`,
}

var grpcTemplates = map[string]string{
	"go.mod": `module {{ .ModulePath }}

go 1.26.3

require (
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.3
	github.com/kitti12911/lib-monitor v1.8.0
	github.com/kitti12911/lib-orm/v2 v2.5.0
	github.com/kitti12911/lib-util/v3 v3.7.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.68.0
	google.golang.org/grpc v1.80.0
)
`,
	"cmd/server/main.go": `package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{ .ModulePath }}/internal/config"
	"{{ .ModulePath }}/internal/database"
	"{{ .ModulePath }}/internal/server"

	"github.com/kitti12911/lib-monitor/profiling"
	"github.com/kitti12911/lib-monitor/tracing"
	libconfig "github.com/kitti12911/lib-util/v3/config"
	"github.com/kitti12911/lib-util/v3/logger"
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
`,
	"internal/server/grpc.go": `package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"{{ .ModulePath }}/internal/server/interceptor"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// GRPCServer owns the gRPC listener and health state.
type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
	health   *health.Server
}

// NewGRPCServer creates the gRPC server and base middleware.
func NewGRPCServer(ctx context.Context, port int) (*GRPCServer, error) {
	listenConfig := net.ListenConfig{}
	listener, err := listenConfig.Listen(ctx, "tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	recoveryOpt := recovery.WithRecoveryHandler(func(p any) error {
		return status.Errorf(codes.Internal, "internal server error: %v", p)
	})

	srv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithFilter(interceptor.TraceableRPC),
		)),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpt),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recoveryOpt),
		),
	)

	healthServer := health.NewServer()
	healthv1.RegisterHealthServer(srv, healthServer)
	healthServer.SetServingStatus("", healthv1.HealthCheckResponse_SERVING)

	reflection.Register(srv)

	return &GRPCServer{
		server:   srv,
		listener: listener,
		health:   healthServer,
	}, nil
}

// Start begins serving gRPC requests.
func (s *GRPCServer) Start() error {
	slog.Info("gRPC server listening", "addr", s.listener.Addr().String())
	return s.server.Serve(s.listener)
}

// Stop gracefully stops the gRPC server.
func (s *GRPCServer) Stop(ctx context.Context) {
	s.health.Shutdown()

	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		slog.WarnContext(ctx, "graceful shutdown timed out, forcing stop")
		s.server.Stop()
	}
}
`,
	"internal/server/interceptor/interceptor.go": `package interceptor

import (
	"strings"

	"google.golang.org/grpc/stats"
)

// TraceableRPC reports whether a gRPC method should emit tracing spans.
func TraceableRPC(info *stats.RPCTagInfo) bool {
	if info == nil {
		return true
	}

	return !IsHealthCheck(info.FullMethodName)
}

// IsHealthCheck reports whether the full gRPC method belongs to the health service.
func IsHealthCheck(fullMethod string) bool {
	return strings.HasPrefix(fullMethod, "/grpc.health.v1.Health/")
}
`,
}
