# ____________________ Go Command ____________________
tidy:
	go mod tidy

run:
	go run ./cmd/service-gen

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

# ____________________ Generate Examples ____________________
examples: example-grpc example-oas example-worker

example-grpc:
	go run ./cmd/service-gen -name demo-grpc -module github.com/kitti12911/demo-grpc -pattern grpc -out tmp/demo-grpc -force -no-tidy -no-git

example-oas:
	go run ./cmd/service-gen -name demo-oas -module github.com/kitti12911/demo-oas -pattern oas -out tmp/demo-oas -force -no-tidy -no-git

example-worker:
	go run ./cmd/service-gen -name demo-worker -module github.com/kitti12911/demo-worker -pattern worker -out tmp/demo-worker -force -no-tidy -no-git

# force update #1