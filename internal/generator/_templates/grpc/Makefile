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
# mapgen is zero-config: it discovers models, params structs, and proto types
# by convention. See lib-orm/cmd/mapgen for the conventions and directives.
gen: gen-proto gen-go gen-mappers

gen-proto:
	rm -rf gen/grpc
	buf generate

gen-go:
	rm -rf gen/database
	go run ___LIB_PATH___/lib-orm/v4/cmd/mapgen@___LIB_ORM_VERSION___ fields
	go run ___LIB_PATH___/lib-orm/v4/cmd/mapgen@___LIB_ORM_VERSION___ patch
	go run ___LIB_PATH___/lib-orm/v4/cmd/mapgen@___LIB_ORM_VERSION___ filter

gen-mappers:
	go run ___LIB_PATH___/lib-orm/v4/cmd/mapgen@___LIB_ORM_VERSION___ map
