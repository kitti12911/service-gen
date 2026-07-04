package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratePatterns(t *testing.T) {
	cases := []struct {
		pattern   string
		mustHave  []fileCheck
		mustMiss  []string // relative paths that must NOT exist
		fileCount int      // minimum expected file count
	}{
		{
			pattern: PatternGRPC,
			mustHave: []fileCheck{
				{"go.mod", "module github.com/kitti12911/demo-grpc"},
				{"go.mod", "go 1.26.4"},
				{"go.mod", "github.com/kitti12911/lib-orm/v3"},
				{"internal/server/grpc.go", "NewGRPCServer"},
				{"internal/server/grpc.go", "RegisterHealthServer"},
				{"internal/feature/starter/starter.go", "func Ping"},
				{"internal/database/database.go", "github.com/kitti12911/lib-orm/v3"},
				{"buf.gen.yaml", "directory: proto"},
				{"config.example.yml", "name: demo-grpc"},
				{".github/workflows/go-ci.yml", "actions/checkout"},
				{".golangci.yml", "github.com/kitti12911/demo-grpc"},
			},
			mustMiss: []string{
				".gitlab-ci.yml",
				"internal/feature/user",
				"internal/feature/worker",
				"internal/database/migrations",
				"internal/database/seeders",
			},
			fileCount: 40,
		},
		{
			pattern: PatternWorker,
			mustHave: []fileCheck{
				{"go.mod", "module github.com/kitti12911/demo-worker"},
				{"internal/worker/handler.go", "package worker"},
				{".github/workflows/go-ci.yml", "actions/checkout"},
			},
			mustMiss:  []string{".gitlab-ci.yml"},
			fileCount: 30,
		},
		{
			pattern: PatternOAS,
			mustHave: []fileCheck{
				{"go.mod", "module github.com/kitti12911/demo-oas"},
				{"go.mod", "github.com/danielgtaylor/huma/v2"},
				{"internal/api/system/api.go", "/health"},
				{"internal/server/http.go", "NewHTTPServer"},
				{"cmd/gen-oas/main.go", "OpenAPI"},
				{".github/workflows/go-ci.yml", "actions/checkout"},
			},
			mustMiss: []string{
				".gitlab-ci.yml",
				"internal/api/users",
				"internal/api/worker",
				"cmd/gen-patch",
				"buf.gen.yaml",
			},
			fileCount: 50,
		},
	}

	for _, tc := range cases {
		t.Run(tc.pattern, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "demo-"+tc.pattern)
			err := Generate(Config{
				Name:       "demo-" + tc.pattern,
				ModulePath: "github.com/kitti12911/demo-" + tc.pattern,
				OutputDir:  dir,
				Pattern:    tc.pattern,
				LibPath:    "github.com/kitti12911",
				NoTidy:     true,
				NoGit:      true,
			})
			if err != nil {
				t.Fatal(err)
			}

			count := countFiles(t, dir)
			if count < tc.fileCount {
				t.Errorf("got %d files, expected at least %d", count, tc.fileCount)
			}

			for _, check := range tc.mustHave {
				assertFileContains(t, dir, check.path, check.want)
			}
			for _, miss := range tc.mustMiss {
				if _, err := os.Stat(filepath.Join(dir, miss)); err == nil {
					t.Errorf("path %s should not exist", miss)
				}
			}
		})
	}
}

func TestGenerateRejectsBadInputs(t *testing.T) {
	cases := []struct {
		name string
		cfg  Config
		want string
	}{
		{
			name: "bad name",
			cfg:  Config{Name: "BadName", ModulePath: "x", OutputDir: t.TempDir(), Pattern: PatternGRPC, LibPath: "github.com/kitti12911"},
			want: "kebab-case",
		},
		{
			name: "missing module",
			cfg:  Config{Name: "ok", OutputDir: t.TempDir(), Pattern: PatternGRPC, LibPath: "github.com/kitti12911"},
			want: "module",
		},
		{
			name: "missing pattern",
			cfg:  Config{Name: "ok", ModulePath: "x", OutputDir: t.TempDir(), LibPath: "github.com/kitti12911"},
			want: "pattern",
		},
		{
			name: "unknown pattern",
			cfg:  Config{Name: "ok", ModulePath: "x", OutputDir: t.TempDir(), Pattern: "rest", LibPath: "github.com/kitti12911"},
			want: "unknown pattern",
		},
		{
			name: "missing lib-path",
			cfg:  Config{Name: "ok", ModulePath: "x", OutputDir: t.TempDir(), Pattern: PatternGRPC},
			want: "lib-path is required",
		},
		{
			name: "lib-path with trailing lib segment",
			cfg:  Config{Name: "ok", ModulePath: "x", OutputDir: t.TempDir(), Pattern: PatternGRPC, LibPath: "github.com/kitti12911/lib-util"},
			want: "lib-*",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.cfg.NoTidy = true
			tc.cfg.NoGit = true
			err := Generate(tc.cfg)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected error containing %q, got %v", tc.want, err)
			}
		})
	}
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
		Pattern:    PatternGRPC,
		LibPath:    "github.com/kitti12911",
		NoTidy:     true,
		NoGit:      true,
	})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected existing file error, got %v", err)
	}
}

func TestGenerateUsesConfiguredCodeOwner(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "demo")
	err := Generate(Config{
		Name:       "demo",
		ModulePath: "github.com/kitti12911/demo",
		OutputDir:  dir,
		CodeOwner:  "@example/platform",
		Pattern:    PatternWorker,
		LibPath:    "github.com/kitti12911",
		NoTidy:     true,
		NoGit:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertFileContains(t, dir, ".github/CODEOWNERS", "@example/platform")
}

func TestGenerateSubstitutesLibPath(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "demo")
	err := Generate(Config{
		Name:       "demo",
		ModulePath: "github.com/kitti12911/demo",
		OutputDir:  dir,
		Pattern:    PatternWorker,
		LibPath:    "github.com/kitti12911",
		NoTidy:     true,
		NoGit:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertFileContains(t, dir, "go.mod", "github.com/kitti12911/lib-util/v3 "+DefaultLibUtilVersion)
	assertFileContains(t, dir, "internal/config/config.go", "github.com/kitti12911/lib-monitor")
}

func TestGenerateOverridesLibVersions(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "demo")
	err := Generate(Config{
		Name:              "demo",
		ModulePath:        "github.com/kitti12911/demo",
		OutputDir:         dir,
		Pattern:           PatternGRPC,
		LibPath:           "github.com/kitti12911",
		LibUtilVersion:    "v3.99.0",
		LibMonitorVersion: "v2.0.0",
		LibOrmVersion:     "v3.42.0",
		NoTidy:            true,
		NoGit:             true,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertFileContains(t, dir, "go.mod", "lib-util/v3 v3.99.0")
	assertFileContains(t, dir, "go.mod", "lib-monitor v2.0.0")
	assertFileContains(t, dir, "go.mod", "lib-orm/v3 v3.42.0")
}

type fileCheck struct {
	path string
	want string
}

func assertFileContains(t *testing.T, root, rel, want string) {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	if !strings.Contains(string(body), want) {
		t.Fatalf("%s does not contain %q\n---\n%s", rel, want, body)
	}
}

func countFiles(t *testing.T, root string) int {
	t.Helper()
	n := 0
	err := filepath.Walk(root, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			n++
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return n
}
