package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateGRPC(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "demo-grpc")

	err := Generate(Config{
		Name:       "demo-grpc",
		ModulePath: "github.com/kitti12911/demo-grpc",
		OutputDir:  dir,
	})
	if err != nil {
		t.Fatal(err)
	}

	assertFileContains(t, dir, "internal/server/grpc.go", "RegisterHealthServer")
	assertFileContains(t, dir, "internal/server/grpc.go", "RegisterStarterServiceServer")
	assertFileContains(t, dir, "internal/server/grpc.go", "NewGRPCServer")
	assertFileContains(t, dir, "internal/feature/starter/handler.go", "Ping")
	assertFileContains(t, dir, "internal/database/database.go", "github.com/kitti12911/lib-orm/v2")
	assertFileContains(t, dir, "internal/config/config.go", "github.com/kitti12911/lib-util/v3/logger")
	assertFileContains(t, dir, "config.example.yml", "port: 50051")
	assertFileContains(t, dir, "buf.gen.yaml", "github.com/kitti12911/demo-grpc/gen/grpc")
	assertFileContains(t, dir, "proto/demo_grpc/v1/starter.proto", "service StarterService")
	assertFileContains(t, dir, "README.md", "Internal gRPC service")
	assertFileContains(t, dir, ".github/workflows/go-ci.yml", "actions/setup-go")
	assertFileContains(t, dir, ".github/workflows/go-ci.yml", "Update Helm Values")
	assertFileContains(t, dir, ".github/workflows/go-ci.yml", "golangci-lint-action")
	assertFileContains(t, dir, ".github/workflows/release.yml", "Build and Scan Image")
	assertFileContains(t, dir, ".releaserc.json", `"prerelease": "beta"`)
	assertFileContains(t, dir, ".prettierrc.json", `"tabWidth": 4`)
	assertFileContains(t, dir, ".prettierignore", "tmp/")
	assertFileContains(t, dir, ".golangci.yml", "github.com/kitti12911/demo-grpc")
	assertFileContains(t, dir, "go.mod", "github.com/kitti12911/lib-orm/v2")
	assertFileContains(t, dir, "go.mod", "google.golang.org/grpc")
	assertFileContains(t, dir, "go.mod", "go 1.26.3")
}

func TestGenerateUsesConfiguredCodeOwner(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "demo-grpc")

	err := Generate(Config{
		Name:       "demo-grpc",
		ModulePath: "github.com/kitti12911/demo-grpc",
		OutputDir:  dir,
		CodeOwner:  "@example/platform",
	})
	if err != nil {
		t.Fatal(err)
	}

	assertFileContains(t, dir, ".github/CODEOWNERS", "/.github/ @example/platform")
	assertFileContains(t, dir, ".github/CODEOWNERS", "/Makefile @example/platform")
}

func TestGenerateRefusesExistingFile(t *testing.T) {
	dir := t.TempDir()
	//nolint:gosec // Test fixture is a regular repository file, not a secret.
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Generate(Config{
		Name:       "demo-grpc",
		ModulePath: "github.com/kitti12911/demo-grpc",
		OutputDir:  dir,
	})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected existing file error, got %v", err)
	}
}

func assertFileContains(t *testing.T, root, rel, want string) {
	t.Helper()

	body, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), want) {
		t.Fatalf("%s does not contain %q", rel, want)
	}
}
