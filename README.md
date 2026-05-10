# service-gen

Small bootstrap generator for homelab gRPC services.

It creates a starter project for an internal gRPC service. OpenAPI/gateway
projects stay separate because they are the outside-world contract boundary.
The generator is intentionally simple: it renders whole template files to the
target directory and avoids hidden framework behavior.

## Installation

Install the latest released command:

```sh
go install github.com/kitti12911/service-gen/cmd/service-gen@latest
```

## Usage

Generate a gRPC service:

```sh
go run ./cmd/service-gen \
    -name demo-grpc \
    -module github.com/kitti12911/demo-grpc \
    -out ../demo-grpc
```

Use `-force` when regenerating over an existing generated target.

## Generated Files

Generated services include:

- `.github` workflows, CODEOWNERS, and Renovate config.
- `.vscode` workspace recommendations and formatting settings.
- `.markdownlint-cli2.jsonc`, `.gitignore`, `Makefile`, `Dockerfile`, and
  `README.md`.
- `config.example.yml` and Go config loading from environment variables.
- A starter database connection helper using `lib-orm`.

The generated service starts a gRPC server with the standard gRPC health
service.

## Local Verification

```sh
make format
make lint
make test
```
