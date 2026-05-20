# Files outside the business logic surface (main, DB setup, gRPC handlers,
# repositories, query glue, server bootstrap) are dropped from coverage so the
# reported % reflects code worth testing. Patterns are awk regexes matched
# against the file:line column of coverage.out.
GO_COVERAGE_EXCLUDE_REGEX = /cmd/|/internal/database/|/internal/server/grpc\.go:

# ____________________ Go Command ____________________
air:
	air

tidy:
	go mod tidy

run:
	go run ./cmd/server/main.go

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
	GO_COVERAGE_EXCLUDE_REGEX='$(GO_COVERAGE_EXCLUDE_REGEX)' ./scripts/ci/go-test.sh
	go tool cover -html=coverage.out

fix: 
	go fix ./...

# ____________________ Generate Command ____________________
# Once you have created a bunprotogen.yaml and uncommented the gen-mappers
# target below, add it to this list so `make gen` regenerates mappers too.
gen: gen-proto gen-go

gen-proto:
	rm -rf gen/grpc
	buf generate

gen-go:
	rm -rf gen/database
	go run ___LIB_PATH___/lib-orm/v3/cmd/fieldmapgen@___LIB_ORM_VERSION___ -model-dir internal/database -root Example -out gen/database/fieldmap_generated.go -package database

# gen-mappers regenerates proto<->bun mapper functions with bunprotogen.
# Customisation (enum bridging, relation mapping) happens via converters
# in bunprotogen.yaml + named functions in the target package — the
# generated file itself should never be hand-edited. Uncomment once you
# have defined a proto and created a bunprotogen.yaml at the repo root
# describing each (bun, proto) pair. See lib-orm/cmd/bunprotogen for the
# config schema.
#
# gen-mappers:
#	go run ___LIB_PATH___/lib-orm/v3/cmd/bunprotogen@___LIB_ORM_VERSION___ -config bunprotogen.yaml
