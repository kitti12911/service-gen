# **_NAME_**

OpenAPI (REST) service scaffolded by
[`service-gen`](https://github.com/kitti12911/service-gen).

## What's included

- [Huma](https://huma.rocks)-based HTTP server with a `/health` endpoint.
- OpenAPI document served at `/openapi.json` and `/openapi.yaml` (plus
  `/download` variants) and an offline Swagger UI at `/docs`.
- Configuration via `config.yml` / environment variables (`lib-util/v3`).
- Structured logging, OpenTelemetry tracing, and Pyroscope profiling.
- Access-log, gzip, and panic-recovery middleware.
- `cmd/gen-oas` to print the OpenAPI document and `cmd/openapi-report` to
  diff/report OpenAPI changes in CI.
- GitHub Actions, Renovate, CODEOWNERS, golangci-lint,
  markdownlint, Prettier, `.air.toml`, and a multi-stage `Dockerfile`.

## Getting started

```sh
cp config.example.yml config.yml   # then edit values
make run
```

Then open <http://localhost:8080/docs>.

## Prerequisites

Install the third-party CLIs this repo expects. Versions are not pinned —
match `go.mod` for Go itself.

### macOS (Homebrew)

```sh
brew install go golangci-lint buf prettier markdownlint-cli2
go install github.com/air-verse/air@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Common commands

| Command            | Description                          |
| ------------------ | ------------------------------------ |
| `make run`         | Start the HTTP server locally        |
| `make air`         | Run with live reload                 |
| `make gen-openapi` | Print the OpenAPI document to stdout |
| `make test`        | Run tests with the race detector     |
| `make lint`        | Run Go and Markdown linting          |
| `make format`      | Format Go, Markdown, YAML, JSON      |

## Adding endpoints

Add a package under `internal/api/`, register it in
`internal/server/http.go`, and follow `internal/api/system` as a template.

## Deployment

A starter Helm chart is included in [`helm/`](helm/) so the service is
deployable out of the box. Deployment state should not live in the service
repository long-term — move the chart to your `helm-sandbox` repository (or
your preferred chart repository), then delete `helm/` from here:

```sh
mv helm /path/to/helm-sandbox/charts/___NAME___
```

Keep app-specific values (image, env, scaling, config) in the chart
repository, not in the service repo.
