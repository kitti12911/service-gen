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
	Name       string
	ModulePath string
	OutputDir  string
	CodeOwner  string
	Pattern    string
	CI         string
	Force      bool
	NoTidy     bool
	NoGit      bool
}

const (
	PatternGRPC   = "grpc"
	PatternOAS    = "oas"
	PatternWorker = "worker"

	CIGitHub = "github"
	CIGitLab = "gitlab"
	CIBoth   = "both"
)

var (
	validName     = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	validPatterns = map[string]struct{}{PatternGRPC: {}, PatternOAS: {}, PatternWorker: {}}
	validCIs      = map[string]struct{}{CIGitHub: {}, CIGitLab: {}, CIBoth: {}}
)

// Generate writes a service bootstrap project to Config.OutputDir.
func Generate(cfg Config) error {
	cfg.Name = strings.TrimSpace(cfg.Name)
	cfg.ModulePath = strings.TrimSpace(cfg.ModulePath)
	cfg.OutputDir = strings.TrimSpace(cfg.OutputDir)
	cfg.CodeOwner = strings.TrimSpace(cfg.CodeOwner)
	cfg.Pattern = strings.TrimSpace(cfg.Pattern)
	cfg.CI = strings.TrimSpace(cfg.CI)

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
		cfg.CI = CIBoth
	}
	if _, ok := validCIs[cfg.CI]; !ok {
		return fmt.Errorf("unknown ci %q (must be github, gitlab, or both)", cfg.CI)
	}

	replacer := strings.NewReplacer(
		"___MODULE___", cfg.ModulePath,
		"___NAME___", cfg.Name,
		"___CODE_OWNER___", cfg.CodeOwner,
	)

	ctx := context.Background()

	if err := walkAndWrite(cfg, "_templates/"+cfg.Pattern, replacer); err != nil {
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

func walkAndWrite(cfg Config, root string, replacer *strings.Replacer) error {
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
		return writeFile(cfg.OutputDir, rel, string(body), replacer, cfg.Force)
	})
	if err != nil {
		return fmt.Errorf("walk templates %s: %w", root, err)
	}
	return nil
}

func ciAllows(ci, rel string) bool {
	switch ci {
	case CIGitHub:
		return rel != ".gitlab-ci.yml"
	case CIGitLab:
		return !strings.HasPrefix(rel, ".github/")
	default:
		return true
	}
}

func writeFile(root, rel, body string, replacer *strings.Replacer, force bool) error {
	// Templates may carry a .tmpl suffix so files like go.mod don't make Go's
	// embed see a nested module. Strip the suffix from the output path.
	rel = strings.TrimSuffix(rel, ".tmpl")
	target := filepath.Join(root, filepath.FromSlash(replacer.Replace(rel)))
	if !force {
		if _, statErr := os.Stat(target); statErr == nil {
			return fmt.Errorf("%s already exists; pass -force to overwrite", target)
		} else if !errors.Is(statErr, os.ErrNotExist) {
			return fmt.Errorf("stat %s: %w", target, statErr)
		}
	}
	rendered := replacer.Replace(body)
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
