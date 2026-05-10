package generator

import (
	"bytes"
	"errors"
	"fmt"
	"maps"
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
	CodeOwner  string
	Force      bool
}

type data struct {
	Name             string
	ModulePath       string
	BinaryName       string
	GRPCPort         string
	CodeOwner        string
	ProtoPackage     string
	ProtoPackagePath string
}

var validName = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// Generate writes a service bootstrap project to Config.OutputDir.
func Generate(cfg Config) error {
	cfg.Name = strings.TrimSpace(cfg.Name)
	cfg.ModulePath = strings.TrimSpace(cfg.ModulePath)
	cfg.OutputDir = strings.TrimSpace(cfg.OutputDir)
	cfg.CodeOwner = strings.TrimSpace(cfg.CodeOwner)

	if !validName.MatchString(cfg.Name) {
		return errors.New("name must use lowercase kebab-case, for example grpc-sandbox")
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

	d := data{
		Name:             cfg.Name,
		ModulePath:       cfg.ModulePath,
		BinaryName:       cfg.Name,
		GRPCPort:         "50051",
		CodeOwner:        cfg.CodeOwner,
		ProtoPackage:     strings.ReplaceAll(cfg.Name, "-", "_"),
		ProtoPackagePath: strings.ReplaceAll(cfg.Name, "-", "_"),
	}

	files := make(map[string]string, len(commonTemplates)+len(grpcTemplates))
	maps.Copy(files, commonTemplates)
	maps.Copy(files, grpcTemplates)

	for rel, body := range files {
		if err := writeTemplate(cfg.OutputDir, rel, body, d, cfg.Force); err != nil {
			return err
		}
	}

	return nil
}

func writeTemplate(root, rel, body string, d data, force bool) error {
	renderedRel, err := renderTemplate(rel, rel, d, false)
	if err != nil {
		return err
	}

	target := filepath.Join(root, filepath.FromSlash(renderedRel))
	if !force {
		if _, statErr := os.Stat(target); statErr == nil {
			return fmt.Errorf("%s already exists; pass -force to overwrite", target)
		} else if !errors.Is(statErr, os.ErrNotExist) {
			return fmt.Errorf("stat %s: %w", target, statErr)
		}
	}

	rendered, err := renderTemplate(rel, body, d, true)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("create parent directory for %s: %w", target, err)
	}
	//nolint:gosec // Generated project files should be editable by normal repository tooling.
	if err := os.WriteFile(target, []byte(rendered), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", target, err)
	}
	return nil
}

func renderTemplate(name, body string, d data, normalizeFile bool) (string, error) {
	tmpl, err := template.New(name).Parse(body)
	if err != nil {
		return "", fmt.Errorf("parse %s: %w", name, err)
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, d); err != nil {
		return "", fmt.Errorf("render %s: %w", name, err)
	}
	rendered := strings.NewReplacer(
		"__GHA_OPEN__", "{{",
		"__GHA_CLOSE__", "}}",
		"__TPL_OPEN__", "{{",
		"__TPL_CLOSE__", "}}",
	).Replace(out.String())
	if normalizeFile {
		rendered = strings.TrimRight(rendered, " \t\r\n") + "\n"
	}
	return rendered, nil
}
