package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yourorg/envsync/internal/diff2"
	"github.com/yourorg/envsync/internal/envfile"
)

func runDiff2(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("diff2", flag.ContinueOnError)
	changesOnly := fs.Bool("changes-only", false, "only print changed/added/removed lines")
	fs.SetOutput(out)

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return fmt.Errorf("usage: envsync diff2 <left.env> <right.env>")
	}

	leftPath := fs.Arg(0)
	rightPath := fs.Arg(1)

	leftEnv, err := loadEnvMap(leftPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", leftPath, err)
	}
	rightEnv, err := loadEnvMap(rightPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", rightPath, err)
	}

	result := diff2.Diff(leftEnv, rightEnv)

	for _, line := range result.Lines {
		if *changesOnly && line.Kind == diff2.KindEqual {
			continue
		}
		fmt.Fprintln(out, line.String())
	}

	fmt.Fprintf(out, "\nsummary: %s\n", result.Summary())

	if result.HasChanges() {
		return fmt.Errorf("environments differ")
	}
	return nil
}

func loadEnvMap(path string) (map[string]string, error) {
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
