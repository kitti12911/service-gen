# service-gen

Bootstrap generator for homelab services.

It scaffolds a batteries-included starter project for one of three patterns —
**gRPC**, **OpenAPI (REST)**, or **worker** — wired with the workspace's
logging, tracing, profiling, linting, and CI conventions. The generator is
intentionally simple: it renders whole template files to the target directory
with flat string substitution and no hidden framework behavior.

Generated projects build, test, and lint cleanly out of the box.

## Prerequisites

Install the third-party CLIs this repo expects. Match `go.mod` for the Go
version.

### macOS (Homebrew)

```sh
brew install go golangci-lint prettier markdownlint-cli2
```

## Installation

Install the latest released command:

```sh
go install github.com/kitti12911/service-gen/v2/cmd/service-gen@latest
```

Install a specific release directly from GitHub without using the Go module
proxy:

```sh
GOPROXY=direct go install github.com/kitti12911/service-gen/v2/cmd/service-gen@v2.2.1
```

When working from a local checkout, install the current source instead:

```sh
go install ./cmd/service-gen
```

## Usage

```sh
go run ./cmd/service-gen \
    -name demo-grpc \
    -module github.com/kitti12911/demo-grpc \
    -pattern grpc \
    -lib-path github.com/kitti12911 \
    -out ../demo-grpc \
    -code-owner @kitti12911
```

### Flags

| Flag          | Required | Default       | Description                                                                                                                                                |
| ------------- | -------- | ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `-name`       | yes      |               | Project name in lowercase kebab-case                                                                                                                       |
| `-module`     | yes      |               | Go module path                                                                                                                                             |
| `-pattern`    | yes      |               | `grpc`, `oas`, or `worker`                                                                                                                                 |
| `-lib-path`   | yes      |               | Base path for `lib-*` dependencies, without the trailing `lib-*` segment, for example `github.com/kitti12911`.                                             |
| `-out`        | no       | `-name`       | Output directory                                                                                                                                           |
| `-code-owner` | no       | `@kitti12911` | CODEOWNERS owner                                                                                                                                           |
| `-force`      | no       | `false`       | Overwrite existing generated files                                                                                                                         |
| `-no-tidy`    | no       | `false`       | Skip `go mod tidy` after generation                                                                                                                        |
| `-no-git`     | no       | `false`       | Skip `git init` and initial commit                                                                                                                         |

By default the generator runs `go mod tidy` and creates an initial Git commit
so the project is immediately usable.

## Patterns

| Pattern  | What you get                                                                                                                                                                                                                       |
| -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `grpc`   | gRPC server with health + reflection, `lib-orm/v3` DB connection wiring, a `feature/starter` package, and a local `proto/` directory ready for `make gen`. No example CRUD or migrations — migrations live in `migration-sandbox`. |
| `oas`    | Huma-based OpenAPI/REST service with a `/health` endpoint, OpenAPI JSON/YAML serving, embedded Swagger UI, and an OpenAPI diff/report CI tool. No DB.                                                                              |
| `worker` | Minimal `lib-async` worker with an event loop and handler.                                                                                                                                                                         |

Every pattern ships GitHub Actions, Renovate, CODEOWNERS, golangci-lint,
markdownlint, Prettier, `.air.toml`, a multi-stage `Dockerfile`, and
semantic-release.

The `internal/feature` / `internal/api` example is deliberately a single simple
function. Add real business logic and (for gRPC) `.proto` files yourself; run
`make gen` to regenerate protobuf code.

## Local Verification

```sh
make format
make lint
make test
```
