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
	flag.StringVar(&cfg.CI, "ci", "", "CI flavor to emit (required): github | gitlab")
	flag.StringVar(&cfg.LibPath, "lib-path", "", "base path for lib-* dependencies (required), for example github.com/kitti12911 or gitlab.bu8-sd.com/sdo/pharse-3")
	flag.StringVar(&cfg.LibUtilVersion, "lib-util-version", generator.DefaultLibUtilVersion, "version of lib-util/v3 to require in go.mod")
	flag.StringVar(&cfg.LibMonitorVersion, "lib-monitor-version", generator.DefaultLibMonitorVersion, "version of lib-monitor to require in go.mod")
	flag.StringVar(&cfg.LibOrmVersion, "lib-orm-version", generator.DefaultLibOrmVersion, "version of lib-orm/v3 to require in go.mod (grpc pattern only)")
	flag.StringVar(&cfg.LibAsyncVersion, "lib-async-version", generator.DefaultLibAsyncVersion, "version of lib-async to require in go.mod (worker pattern only)")
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
