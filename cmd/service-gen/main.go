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
	flag.StringVar(&cfg.OutputDir, "out", "", "output directory")
	flag.BoolVar(&cfg.Force, "force", false, "overwrite existing generated files")
	flag.Parse()

	if cfg.OutputDir == "" && cfg.Name != "" {
		cfg.OutputDir = filepath.Clean(cfg.Name)
	}

	if err := generator.Generate(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(1)
	}

	log.Printf("generated grpc service in %s", cfg.OutputDir)
}
