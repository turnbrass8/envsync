package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envsync/internal/envfile"
	"github.com/yourorg/envsync/internal/rotate"
)

func runRotate(args []string) error {
	fs := flag.NewFlagSet("rotate", flag.ContinueOnError)
	envPath := fs.String("env", ".env", "path to .env file")
	rulesPath := fs.String("rules", ".rotate", "path to rotation rules file")
	dryRun := fs.Bool("dry-run", false, "preview rotations without writing")
	outPath := fs.String("out", "", "output path (defaults to --env)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Parse env file
	env, err := envfile.Parse(*envPath)
	if err != nil {
		return fmt.Errorf("rotate: reading env file: %w", err)
	}

	// Parse rules file
	rf, err := os.Open(*rulesPath)
	if err != nil {
		return fmt.Errorf("rotate: opening rules file: %w", err)
	}
	defer rf.Close()

	rules, err := rotate.ParseRules(rf)
	if err != nil {
		return fmt.Errorf("rotate: parsing rules: %w", err)
	}

	// Perform rotation
	updated, results, err := rotate.Rotate(env.Values(), rules, *dryRun)
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	// Print results
	for _, r := range results {
		status := "rotated"
		if !r.Rotated {
			status = "dry-run"
		}
		fmt.Printf("  [%s] %s: %s -> %s\n", status, r.Key,
			redactValue(r.OldValue), redactValue(r.NewValue))
	}

	if *dryRun {
		fmt.Println("dry-run: no changes written")
		return nil
	}

	dest := *outPath
	if dest == "" {
		dest = *envPath
	}

	if err := envfile.Write(dest, updated); err != nil {
		return fmt.Errorf("rotate: writing output: %w", err)
	}

	fmt.Printf("rotated %d key(s) -> %s\n", len(results), dest)
	return nil
}

func redactValue(v string) string {
	if len(v) <= 4 {
		return strings.Repeat("*", len(v))
	}
	return v[:2] + strings.Repeat("*", len(v)-4) + v[len(v)-2:]
}
