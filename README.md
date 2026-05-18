# service-gen

Bootstrap generator for homelab services.

It scaffolds a batteries-included starter project for one of three patterns —
**gRPC**, **OpenAPI (REST)**, or **worker** — wired with the workspace's
logging, tracing, profiling, linting, and CI conventions. The generator is
intentionally simple: it renders whole template files to the target directory
with flat string substitution and no hidden framework behavior.

Generated projects build, test, and lint cleanly out of the box.

## Installation

Install the latest released command:

```sh
go install github.com/kitti12911/service-gen/cmd/service-gen@latest
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
    -ci github \
    -lib-path github.com/kitti12911 \
    -out ../demo-grpc \
    -code-owner @kitti12911
```

### Flags

| Flag                   | Required | Default       | Description                                                                                                                                                |
| ---------------------- | -------- | ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `-name`                | yes      |               | Project name in lowercase kebab-case                                                                                                                       |
| `-module`              | yes      |               | Go module path                                                                                                                                             |
| `-pattern`             | yes      |               | `grpc`, `oas`, or `worker`                                                                                                                                 |
| `-ci`                  | yes      |               | `github` or `gitlab`                                                                                                                                       |
| `-lib-path`            | yes      |               | Base path for `lib-*` dependencies, without the trailing `lib-*` segment. For example `github.com/kitti12911` or `gitlab.bu8-sd.com/sdo/pharse-3`.         |
| `-lib-util-version`    | no       | `v3.15.0`     | Version of `lib-util/v3` to require in `go.mod`                                                                                                            |
| `-lib-monitor-version` | no       | `v1.12.0`     | Version of `lib-monitor` to require in `go.mod`                                                                                                            |
| `-lib-orm-version`     | no       | `v3.0.1`      | Version of `lib-orm/v3` to require in `go.mod` (`grpc` pattern only)                                                                                       |
| `-lib-async-version`   | no       | `v1.5.1`      | Version of `lib-async` to require in `go.mod` (`worker` pattern only)                                                                                      |
| `-out`                 | no       | `-name`       | Output directory                                                                                                                                           |
| `-code-owner`          | no       | `@kitti12911` | CODEOWNERS owner                                                                                                                                           |
| `-force`               | no       | `false`       | Overwrite existing generated files                                                                                                                         |
| `-no-tidy`             | no       | `false`       | Skip `go mod tidy` after generation                                                                                                                        |
| `-no-git`              | no       | `false`       | Skip `git init` and initial commit                                                                                                                         |

The `lib-*` version defaults match what is currently published at
`github.com/kitti12911`. Override them when generating against a fork or a
mirror with different tags.

By default the generator runs `go mod tidy` and creates an initial Git commit
so the project is immediately usable.

## Patterns

| Pattern  | What you get                                                                                                                                                                                                                       |
| -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `grpc`   | gRPC server with health + reflection, `lib-orm/v3` DB connection wiring, a `feature/starter` package, and a local `proto/` directory ready for `make gen`. No example CRUD or migrations — migrations live in `migration-sandbox`. |
| `oas`    | Huma-based OpenAPI/REST service with a `/health` endpoint, OpenAPI JSON/YAML serving, embedded Swagger UI, and an OpenAPI diff/report CI tool. No DB.                                                                              |
| `worker` | Minimal `lib-async` worker with an event loop and handler.                                                                                                                                                                         |

Every pattern ships GitHub Actions or GitLab CI (per `-ci`), Renovate,
CODEOWNERS, golangci-lint, markdownlint, Prettier, `.air.toml`, a multi-stage
`Dockerfile`, and semantic-release.

The `internal/feature` / `internal/api` example is deliberately a single simple
function. Add real business logic and (for gRPC) `.proto` files yourself; run
`make gen` to regenerate protobuf code.

## Local Verification

```sh
make format
make lint
make test
```
