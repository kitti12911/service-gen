package generator

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed all:_templates
var templatesFS embed.FS

// Config controls how a service project is generated.
type Config struct {
	Name              string
	ModulePath        string
	OutputDir         string
	CodeOwner         string
	Pattern           string
	CI                string
	LibPath           string
	LibUtilVersion    string
	LibMonitorVersion string
	LibOrmVersion     string
	LibAsyncVersion   string
	Force             bool
	NoTidy            bool
	NoGit             bool
}

const (
	PatternGRPC   = "grpc"
	PatternOAS    = "oas"
	PatternWorker = "worker"

	CIGitHub = "github"
	CIGitLab = "gitlab"

	DefaultLibUtilVersion    = "v3.15.0"
	DefaultLibMonitorVersion = "v1.12.0"
	DefaultLibOrmVersion     = "v3.0.1"
	DefaultLibAsyncVersion   = "v1.5.1"
)

var (
	validName     = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	validPatterns = map[string]struct{}{PatternGRPC: {}, PatternOAS: {}, PatternWorker: {}}
	validCIs      = map[string]struct{}{CIGitHub: {}, CIGitLab: {}}

	// goPrivateBlockRE matches a fenced README section that should only render
	// when private-module setup applies. The fence uses HTML comments so the
	// markers don't break Markdown rendering when intentionally kept.
	goPrivateBlockRE  = regexp.MustCompile(`(?s)[ \t]*<!-- IF_GOPRIVATE -->\n.*?<!-- END_GOPRIVATE -->\n?`)
	goPrivateMarkerRE = regexp.MustCompile(`[ \t]*<!-- (?:IF_GOPRIVATE|END_GOPRIVATE) -->\n`)
	// collapseBlankLinesRE collapses runs of 3+ newlines into 2 so stripping a
	// fenced section never leaves stacked blank lines (would trip markdownlint
	// MD012).
	collapseBlankLinesRE = regexp.MustCompile(`\n{3,}`)
)

// stripGoPrivateBlocks removes IF_GOPRIVATE/END_GOPRIVATE fences from rendered
// templates. When keep is true the fences are stripped but content is kept;
// when keep is false the entire fenced section is removed.
func stripGoPrivateBlocks(body string, keep bool) string {
	if keep {
		body = goPrivateMarkerRE.ReplaceAllString(body, "")
	} else {
		body = goPrivateBlockRE.ReplaceAllString(body, "")
	}
	return collapseBlankLinesRE.ReplaceAllString(body, "\n\n")
}

// Generate writes a service bootstrap project to Config.OutputDir.
func Generate(cfg Config) error {
	if err := normalizeAndValidate(&cfg); err != nil {
		return err
	}

	goPrivate := goPrivateFromLibPath(cfg.LibPath)
	replacer := strings.NewReplacer(
		"___MODULE___", cfg.ModulePath,
		"___NAME___", cfg.Name,
		"___CODE_OWNER___", cfg.CodeOwner,
		"___LIB_PATH___", cfg.LibPath,
		"___LIB_NAMESPACE___", libNamespaceFromLibPath(cfg.LibPath),
		"___LIB_UTIL_VERSION___", cfg.LibUtilVersion,
		"___LIB_MONITOR_VERSION___", cfg.LibMonitorVersion,
		"___LIB_ORM_VERSION___", cfg.LibOrmVersion,
		"___LIB_ASYNC_VERSION___", cfg.LibAsyncVersion,
		"___GOPRIVATE___", goPrivate,
	)

	ctx := context.Background()

	if err := walkAndWrite(cfg, "_templates/"+cfg.Pattern, replacer, goPrivate != ""); err != nil {
		return err
	}

	if !cfg.NoTidy {
		if err := runIn(ctx, cfg.OutputDir, "go", "mod", "tidy"); err != nil {
			return fmt.Errorf("go mod tidy: %w", err)
		}
	}
	if !cfg.NoGit {
		if err := initGit(ctx, cfg.OutputDir); err != nil {
			return fmt.Errorf("git init: %w", err)
		}
	}

	return nil
}

func normalizeAndValidate(cfg *Config) error {
	cfg.Name = strings.TrimSpace(cfg.Name)
	cfg.ModulePath = strings.TrimSpace(cfg.ModulePath)
	cfg.OutputDir = strings.TrimSpace(cfg.OutputDir)
	cfg.CodeOwner = strings.TrimSpace(cfg.CodeOwner)
	cfg.Pattern = strings.TrimSpace(cfg.Pattern)
	cfg.CI = strings.TrimSpace(cfg.CI)
	cfg.LibPath = strings.Trim(strings.TrimSpace(cfg.LibPath), "/")
	cfg.LibUtilVersion = defaultIfEmpty(strings.TrimSpace(cfg.LibUtilVersion), DefaultLibUtilVersion)
	cfg.LibMonitorVersion = defaultIfEmpty(strings.TrimSpace(cfg.LibMonitorVersion), DefaultLibMonitorVersion)
	cfg.LibOrmVersion = defaultIfEmpty(strings.TrimSpace(cfg.LibOrmVersion), DefaultLibOrmVersion)
	cfg.LibAsyncVersion = defaultIfEmpty(strings.TrimSpace(cfg.LibAsyncVersion), DefaultLibAsyncVersion)

	if !validName.MatchString(cfg.Name) {
		return errors.New("name must use lowercase kebab-case, for example user-service")
	}
	if cfg.ModulePath == "" {
		return errors.New("module is required")
	}
	if cfg.OutputDir == "" {
		return errors.New("out is required")
	}
	if cfg.CodeOwner == "" {
		cfg.CodeOwner = "@kitti12911"
	}
	if cfg.Pattern == "" {
		return errors.New("pattern is required (grpc, oas, or worker)")
	}
	if _, ok := validPatterns[cfg.Pattern]; !ok {
		return fmt.Errorf("unknown pattern %q (must be grpc, oas, or worker)", cfg.Pattern)
	}
	if cfg.CI == "" {
		return errors.New("ci is required (github or gitlab)")
	}
	if _, ok := validCIs[cfg.CI]; !ok {
		return fmt.Errorf("unknown ci %q (must be github or gitlab)", cfg.CI)
	}
	if cfg.LibPath == "" {
		return errors.New("lib-path is required, for example github.com/kitti12911 or gitlab.bu8-sd.com/sdo/pharse-3")
	}
	if last := cfg.LibPath[strings.LastIndex(cfg.LibPath, "/")+1:]; strings.HasPrefix(last, "lib-") {
		return fmt.Errorf("lib-path %q must not include the trailing lib-* segment", cfg.LibPath)
	}
	return nil
}

func walkAndWrite(cfg Config, root string, replacer *strings.Replacer, keepGoPrivate bool) error {
	err := fs.WalkDir(templatesFS, root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		rel := strings.TrimPrefix(path, root+"/")
		if !ciAllows(cfg.CI, rel) {
			return nil
		}
		body, err := fs.ReadFile(templatesFS, path)
		if err != nil {
			return fmt.Errorf("read embedded %s: %w", path, err)
		}
		return writeFile(cfg.OutputDir, rel, string(body), replacer, cfg.Force, keepGoPrivate)
	})
	if err != nil {
		return fmt.Errorf("walk templates %s: %w", root, err)
	}
	return nil
}

func defaultIfEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

// libNamespaceFromLibPath returns the path portion of LibPath after the host,
// e.g. "gitlab.example.com/team/sub" -> "team/sub". Used in clone examples in
// the README where the host is rendered separately as ___GOPRIVATE___.
func libNamespaceFromLibPath(libPath string) string {
	if i := strings.Index(libPath, "/"); i >= 0 {
		return libPath[i+1:]
	}
	return ""
}

// goPrivateFromLibPath returns the GOPRIVATE host derived from LibPath.
// LibPath rooted at github.com yields an empty string because public github
// modules are served by the default proxy and need no private-host handling.
// Any other host is returned verbatim and ends up in `GOPRIVATE`, the
// Dockerfile env, and the GitLab CI variables block.
func goPrivateFromLibPath(libPath string) string {
	if libPath == "" {
		return ""
	}
	host := libPath
	if i := strings.Index(libPath, "/"); i >= 0 {
		host = libPath[:i]
	}
	if host == "github.com" {
		return ""
	}
	return host
}

func ciAllows(ci, rel string) bool {
	switch ci {
	case CIGitLab:
		return !isGitHubTemplate(rel)
	default:
		return rel != ".gitlab-ci.yml"
	}
}

func writeFile(root, rel, body string, replacer *strings.Replacer, force, keepGoPrivate bool) error {
	// Templates may carry a .tmpl suffix so files like go.mod don't make Go's
	// embed see a nested module. Strip the suffix from the output path.
	rel = strings.TrimSuffix(rel, ".tmpl")
	rel = outputRel(rel)
	target := filepath.Join(root, filepath.FromSlash(replacer.Replace(rel)))
	if !force {
		if _, statErr := os.Stat(target); statErr == nil {
			return fmt.Errorf("%s already exists; pass -force to overwrite", target)
		} else if !errors.Is(statErr, os.ErrNotExist) {
			return fmt.Errorf("stat %s: %w", target, statErr)
		}
	}
	rendered := stripGoPrivateBlocks(replacer.Replace(body), keepGoPrivate)
	rendered = strings.TrimRight(rendered, " \t\r\n") + "\n"
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("create parent directory for %s: %w", target, err)
	}
	mode := os.FileMode(0o644)
	if strings.HasSuffix(rel, ".sh") {
		mode = 0o755
	}
	//nolint:gosec // Generated repository files; permission chosen above per file type.
	if err := os.WriteFile(target, []byte(rendered), mode); err != nil {
		return fmt.Errorf("write %s: %w", target, err)
	}
	return nil
}

func isGitHubTemplate(rel string) bool {
	return rel == "github" || strings.HasPrefix(rel, "github/")
}

func outputRel(rel string) string {
	if rel == "github" {
		return ".github"
	}
	if after, ok := strings.CutPrefix(rel, "github/"); ok {
		return ".github/" + after
	}
	return rel
}

func runIn(ctx context.Context, dir, name string, args ...string) error {
	//nolint:gosec // name and args are constructed from internal callers, not user input.
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run %s: %w", name, err)
	}
	return nil
}

func initGit(ctx context.Context, dir string) error {
	if err := runIn(ctx, dir, "git", "init", "--initial-branch=main"); err != nil {
		return err
	}
	if err := runIn(ctx, dir, "git", "add", "."); err != nil {
		return err
	}
	return runIn(ctx, dir,
		"git",
		"-c", "user.email=service-gen@local",
		"-c", "user.name=service-gen",
		"commit", "-m", "chore: initial scaffold from service-gen",
	)
}
