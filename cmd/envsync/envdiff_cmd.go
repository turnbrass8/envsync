package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envsync/internal/envdiff"
	"github.com/user/envsync/internal/envfile"
)

func runEnvDiff(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)
	showUnchanged := fs.Bool("unchanged", false, "also print unchanged keys")
	colorize := fs.Bool("color", false, "colorize output (added=green, removed=red, modified=yellow)")
	fs.SetOutput(stdout)

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return fmt.Errorf("usage: envsync envdiff [--unchanged] [--color] <base.env> <target.env>")
	}

	baseFile := fs.Arg(0)
	targetFile := fs.Arg(1)

	baseEnv, err := loadEnvMapEnvDiff(baseFile)
	if err != nil {
		return fmt.Errorf("loading base file: %w", err)
	}
	targetEnv, err := loadEnvMapEnvDiff(targetFile)
	if err != nil {
		return fmt.Errorf("loading target file: %w", err)
	}

	result := envdiff.Diff(baseEnv, targetEnv, *showUnchanged)

	for _, line := range result.Lines {
		if *colorize {
			fmt.Fprintln(stdout, colorLine(line))
		} else {
			fmt.Fprintln(stdout, line.String())
		}
	}

	fmt.Fprintf(stdout, "\nsummary: %s\n", result.Summary())

	if result.HasChanges() {
		return fmt.Errorf("files differ")
	}
	return nil
}

func loadEnvMapEnvDiff(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	env, err := envfile.Parse(f)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, len(env))
	for _, e := range env {
		m[e.Key] = e.Value
	}
	return m, nil
}

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
)

func colorLine(l envdiff.Line) string {
	switch l.Kind {
	case envdiff.Added:
		return colorGreen + l.String() + colorReset
	case envdiff.Removed:
		return colorRed + l.String() + colorReset
	case envdiff.Modified:
		return colorYellow + l.String() + colorReset
	default:
		return l.String()
	}
}
