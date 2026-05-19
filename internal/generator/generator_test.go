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
				{"go.mod", "go 1.26.3"},
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
				CI:         CIGitHub,
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

func TestCIFilter(t *testing.T) {
	cases := []struct {
		ci       string
		wantGH   bool
		wantGitL bool
	}{
		{CIGitHub, true, false},
		{CIGitLab, false, true},
	}

	for _, tc := range cases {
		t.Run(tc.ci, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "demo")
			err := Generate(Config{
				Name:       "demo",
				ModulePath: "github.com/kitti12911/demo",
				OutputDir:  dir,
				Pattern:    PatternWorker,
				CI:         tc.ci,
				LibPath:    "github.com/kitti12911",
				NoTidy:     true,
				NoGit:      true,
			})
			if err != nil {
				t.Fatal(err)
			}

			ghExists := exists(filepath.Join(dir, ".github", "workflows", "go-ci.yml"))
			glExists := exists(filepath.Join(dir, ".gitlab-ci.yml"))

			if ghExists != tc.wantGH {
				t.Errorf("github workflows: got exists=%v, want %v", ghExists, tc.wantGH)
			}
			if glExists != tc.wantGitL {
				t.Errorf("gitlab-ci.yml: got exists=%v, want %v", glExists, tc.wantGitL)
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
			cfg:  Config{Name: "BadName", ModulePath: "x", OutputDir: t.TempDir(), Pattern: PatternGRPC, CI: CIGitHub, LibPath: "github.com/kitti12911"},
			want: "kebab-case",
		},
		{
			name: "missing module",
			cfg:  Config{Name: "ok", OutputDir: t.TempDir(), Pattern: PatternGRPC, CI: CIGitHub, LibPath: "github.com/kitti12911"},
			want: "module",
		},
		{
			name: "missing pattern",
			cfg:  Config{Name: "ok", ModulePath: "x", OutputDir: t.TempDir(), CI: CIGitHub, LibPath: "github.com/kitti12911"},
			want: "pattern",
		},
		{
			name: "unknown pattern",
			cfg:  Config{Name: "ok", ModulePath: "x", OutputDir: t.TempDir(), Pattern: "rest", CI: CIGitHub, LibPath: "github.com/kitti12911"},
			want: "unknown pattern",
		},
		{
			name: "missing ci",
			cfg:  Config{Name: "ok", ModulePath: "x", OutputDir: t.TempDir(), Pattern: PatternGRPC, LibPath: "github.com/kitti12911"},
			want: "ci is required",
		},
		{
			name: "unknown ci",
			cfg:  Config{Name: "ok", ModulePath: "x", OutputDir: t.TempDir(), Pattern: PatternGRPC, CI: "circle", LibPath: "github.com/kitti12911"},
			want: "unknown ci",
		},
		{
			name: "rejects both ci",
			cfg:  Config{Name: "ok", ModulePath: "x", OutputDir: t.TempDir(), Pattern: PatternGRPC, CI: "both", LibPath: "github.com/kitti12911"},
			want: "unknown ci",
		},
		{
			name: "missing lib-path",
			cfg:  Config{Name: "ok", ModulePath: "x", OutputDir: t.TempDir(), Pattern: PatternGRPC, CI: CIGitHub},
			want: "lib-path is required",
		},
		{
			name: "lib-path with trailing lib segment",
			cfg:  Config{Name: "ok", ModulePath: "x", OutputDir: t.TempDir(), Pattern: PatternGRPC, CI: CIGitHub, LibPath: "github.com/kitti12911/lib-util"},
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
		CI:         CIGitHub,
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
		CI:         CIGitHub,
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
		ModulePath: "gitlab.bu8-sd.com/sdo/pharse-3/demo",
		OutputDir:  dir,
		Pattern:    PatternWorker,
		CI:         CIGitLab,
		LibPath:    "gitlab.bu8-sd.com/sdo/pharse-3",
		NoTidy:     true,
		NoGit:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertFileContains(t, dir, "go.mod", "gitlab.bu8-sd.com/sdo/pharse-3/lib-util/v3 "+DefaultLibUtilVersion)
	assertFileContains(t, dir, "internal/config/config.go", "gitlab.bu8-sd.com/sdo/pharse-3/lib-monitor")
	if exists(filepath.Join(dir, ".github")) {
		t.Errorf(".github should not exist when ci=gitlab")
	}
}

func TestGoPrivateFromLibPath(t *testing.T) {
	cases := []struct {
		libPath string
		want    string
	}{
		{"", ""},
		{"github.com/kitti12911", ""},
		{"github.com", ""},
		{"gitlab.bu8-sd.com/sdo/pharse-3", "gitlab.bu8-sd.com"},
		{"gitlab.example.com", "gitlab.example.com"},
		{"git.internal.corp/team", "git.internal.corp"},
	}
	for _, tc := range cases {
		got := goPrivateFromLibPath(tc.libPath)
		if got != tc.want {
			t.Errorf("goPrivateFromLibPath(%q) = %q, want %q", tc.libPath, got, tc.want)
		}
	}
}

func TestGeneratePrivateLibPathEmitsGoPrivate(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "demo")
	err := Generate(Config{
		Name:       "demo",
		ModulePath: "gitlab.bu8-sd.com/sdo/pharse-3/demo",
		OutputDir:  dir,
		Pattern:    PatternGRPC,
		CI:         CIGitLab,
		LibPath:    "gitlab.bu8-sd.com/sdo/pharse-3",
		NoTidy:     true,
		NoGit:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertFileContains(t, dir, "Dockerfile", "ARG GOPRIVATE=gitlab.bu8-sd.com")
	assertFileContains(t, dir, "Dockerfile", "--mount=type=secret,id=netrc")
	assertFileContains(t, dir, ".gitlab-ci.yml", `GOPRIVATE: "gitlab.bu8-sd.com"`)
	assertFileContains(t, dir, ".gitlab-ci.yml", "CI_JOB_TOKEN")
	assertFileContains(t, dir, ".gitlab-ci.yml", "BUILD_SECRETS_NETRC=/tmp/netrc")
	assertFileContains(t, dir, "scripts/ci/build-image.sh", "id=netrc,src=${BUILD_SECRETS_NETRC}")
	assertFileContains(t, dir, "README.md", "## Private Go modules")
	assertFileContains(t, dir, "README.md", "gitlab.bu8-sd.com/sdo/pharse-3")
	// SSH clone example uses host + namespace-only path (no host duplication).
	assertFileContains(t, dir, "README.md", "git clone git@gitlab.bu8-sd.com:sdo/pharse-3/lib-util.git")
	// HTTPS clone example combines host + namespace.
	assertFileContains(t, dir, "README.md", "git clone https://gitlab.bu8-sd.com/sdo/pharse-3/lib-util.git")
	assertFileExcludes(t, dir, "README.md", "<!-- IF_GOPRIVATE -->")
	assertFileExcludes(t, dir, "README.md", "<!-- END_GOPRIVATE -->")
	// Catch host duplication regressions (e.g. ___GOPRIVATE___:___LIB_PATH___).
	assertFileExcludes(t, dir, "README.md", "gitlab.bu8-sd.com:gitlab.bu8-sd.com")
	assertFileExcludes(t, dir, "Dockerfile", "ARG GOPRIVATE=github.com")
}

func TestLibNamespaceFromLibPath(t *testing.T) {
	cases := []struct {
		libPath string
		want    string
	}{
		{"", ""},
		{"github.com", ""},
		{"github.com/kitti12911", "kitti12911"},
		{"gitlab.bu8-sd.com/sdo/pharse-3", "sdo/pharse-3"},
	}
	for _, tc := range cases {
		got := libNamespaceFromLibPath(tc.libPath)
		if got != tc.want {
			t.Errorf("libNamespaceFromLibPath(%q) = %q, want %q", tc.libPath, got, tc.want)
		}
	}
}

func TestGeneratePublicLibPathOmitsGoPrivateBlock(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "demo")
	err := Generate(Config{
		Name:       "demo",
		ModulePath: "github.com/kitti12911/demo",
		OutputDir:  dir,
		Pattern:    PatternGRPC,
		CI:         CIGitHub,
		LibPath:    "github.com/kitti12911",
		NoTidy:     true,
		NoGit:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertFileContains(t, dir, "Dockerfile", "ARG GOPRIVATE=")
	assertFileExcludes(t, dir, "README.md", "## Private Go modules")
	assertFileExcludes(t, dir, "README.md", "<!-- IF_GOPRIVATE -->")
	assertFileExcludes(t, dir, "README.md", "<!-- END_GOPRIVATE -->")
}

func TestGenerateOverridesLibVersions(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "demo")
	err := Generate(Config{
		Name:              "demo",
		ModulePath:        "gitlab.bu8-sd.com/sdo/pharse-3/demo",
		OutputDir:         dir,
		Pattern:           PatternGRPC,
		CI:                CIGitLab,
		LibPath:           "gitlab.bu8-sd.com/sdo/pharse-3",
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

func assertFileExcludes(t *testing.T, root, rel, unwanted string) {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	if strings.Contains(string(body), unwanted) {
		t.Fatalf("%s unexpectedly contains %q\n---\n%s", rel, unwanted, body)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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
