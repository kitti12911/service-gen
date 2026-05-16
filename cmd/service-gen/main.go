package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kitti12911/service-gen/internal/generator"
)

func main() {
	var cfg generator.Config

	flag.StringVar(&cfg.Name, "name", "", "project name, for example user-service")
	flag.StringVar(&cfg.ModulePath, "module", "", "Go module path, for example github.com/kitti12911/user-service")
	flag.StringVar(&cfg.OutputDir, "out", "", "output directory (defaults to -name)")
	flag.StringVar(&cfg.CodeOwner, "code-owner", "@kitti12911", "CODEOWNERS owner, for example @team/service-owners")
	flag.StringVar(&cfg.Pattern, "pattern", "", "service pattern: grpc | oas | worker")
	flag.StringVar(&cfg.CI, "ci", generator.CIBoth, "CI flavor to emit: github | gitlab | both")
	flag.BoolVar(&cfg.Force, "force", false, "overwrite existing generated files")
	flag.BoolVar(&cfg.NoTidy, "no-tidy", false, "skip running 'go mod tidy' after generation")
	flag.BoolVar(&cfg.NoGit, "no-git", false, "skip 'git init' and initial commit after generation")
	flag.Parse()

	if cfg.OutputDir == "" && cfg.Name != "" {
		cfg.OutputDir = filepath.Clean(cfg.Name)
	}

	if err := generator.Generate(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(1)
	}

	log.Printf("generated %s service in %s", cfg.Pattern, cfg.OutputDir)
}
