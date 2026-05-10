package generator

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// Config controls how a service project is generated.
type Config struct {
	Name       string
	ModulePath string
	OutputDir  string
	Force      bool
}

type data struct {
	Name       string
	ModulePath string
	BinaryName string
	GRPCPort   string
}

var validName = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// Generate writes a service bootstrap project to Config.OutputDir.
func Generate(cfg Config) error {
	cfg.Name = strings.TrimSpace(cfg.Name)
	cfg.ModulePath = strings.TrimSpace(cfg.ModulePath)
	cfg.OutputDir = strings.TrimSpace(cfg.OutputDir)

	if !validName.MatchString(cfg.Name) {
		return errors.New("name must use lowercase kebab-case, for example grpc-sandbox")
	}
	if cfg.ModulePath == "" {
		return errors.New("module is required")
	}
	if cfg.OutputDir == "" {
		return errors.New("out is required")
	}

	d := data{
		Name:       cfg.Name,
		ModulePath: cfg.ModulePath,
		BinaryName: cfg.Name,
		GRPCPort:   "50051",
	}

	files := make(map[string]string, len(commonTemplates)+len(grpcTemplates))
	for path, body := range commonTemplates {
		files[path] = body
	}
	for path, body := range grpcTemplates {
		files[path] = body
	}

	for rel, body := range files {
		if err := writeTemplate(cfg.OutputDir, rel, body, d, cfg.Force); err != nil {
			return err
		}
	}

	return nil
}

func writeTemplate(root, rel, body string, d data, force bool) error {
	target := filepath.Join(root, filepath.FromSlash(rel))
	if !force {
		if _, err := os.Stat(target); err == nil {
			return fmt.Errorf("%s already exists; pass -force to overwrite", target)
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stat %s: %w", target, err)
		}
	}

	tmpl, err := template.New(rel).Parse(body)
	if err != nil {
		return fmt.Errorf("parse %s: %w", rel, err)
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, d); err != nil {
		return fmt.Errorf("render %s: %w", rel, err)
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("create parent directory for %s: %w", target, err)
	}
	//nolint:gosec // Generated project files should be editable by normal repository tooling.
	if err := os.WriteFile(target, out.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", target, err)
	}
	return nil
}
