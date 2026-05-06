package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envsync/internal/diff3"
	"github.com/user/envsync/internal/envfile"
)

func runDiff3(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("diff3", flag.ContinueOnError)
	showUnchanged := fs.Bool("unchanged", false, "also print unchanged keys")
	fs.SetOutput(out)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 3 {
		return fmt.Errorf("usage: envsync diff3 [--unchanged] <base> <left> <right>")
	}

	baseFile := fs.Arg(0)
	leftFile := fs.Arg(1)
	rightFile := fs.Arg(2)

	base, err := loadEnvMapDiff3(baseFile)
	if err != nil {
		return fmt.Errorf("base: %w", err)
	}
	left, err := loadEnvMapDiff3(leftFile)
	if err != nil {
		return fmt.Errorf("left: %w", err)
	}
	right, err := loadEnvMapDiff3(rightFile)
	if err != nil {
		return fmt.Errorf("right: %w", err)
	}

	result := diff3.Diff(base, left, right)

	for _, e := range result.Entries {
		if e.Kind == diff3.Unchanged && !*showUnchanged {
			continue
		}
		fmt.Fprintln(out, e.String())
	}

	if result.Conflicts > 0 {
		fmt.Fprintf(out, "\n%d conflict(s) detected\n", result.Conflicts)
		return fmt.Errorf("%d conflict(s) detected", result.Conflicts)
	}
	return nil
}

func loadEnvMapDiff3(path string) (map[string]string, error) {
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
	for _, kv := range env {
		m[kv.Key] = kv.Value
	}
	return m, nil
}
