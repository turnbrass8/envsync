package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envsync/internal/envfile"
	"github.com/yourorg/envsync/internal/manifest"
	"github.com/yourorg/envsync/internal/sync"
)

func main() {
	manifestPath := flag.String("manifest", ".envsync", "path to manifest file")
	targetPath := flag.String("env", ".env", "path to target .env file")
	dryRun := flag.Bool("dry-run", false, "print changes without writing")
	flag.Parse()

	if err := run(*manifestPath, *targetPath, *dryRun); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(manifestPath, targetPath string, dryRun bool) error {
	man, err := manifest.Parse(manifestPath)
	if err != nil {
		return fmt.Errorf("reading manifest: %w", err)
	}

	target := envfile.Env{}
	if _, statErr := os.Stat(targetPath); statErr == nil {
		target, err = envfile.Parse(targetPath)
		if err != nil {
			return fmt.Errorf("reading env file: %w", err)
		}
	}

	s := sync.New(dryRun)
	res, err := s.Sync(man, target, targetPath)
	if err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	res.Print(os.Stdout)

	if len(res.Errors) > 0 {
		os.Exit(2)
	}

	return nil
}
