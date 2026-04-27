package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/lint"
)

func runLint(args []string) error {
	fs := flag.NewFlagSet("lint", flag.ContinueOnError)
	errOnly := fs.Bool("errors-only", false, "only report error-level findings")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: envsync lint [options] <envfile>")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		fs.Usage()
		return fmt.Errorf("envfile argument required")
	}

	envPath := fs.Arg(0)
	env, err := envfile.Parse(envPath)
	if err != nil {
		return fmt.Errorf("parsing env file: %w", err)
	}

	findings := lint.Lint(env)

	hasErrors := false
	printed := 0
	for _, f := range findings {
		if *errOnly && f.Severity != lint.SeverityError {
			continue
		}
		fmt.Println(f.String())
		printed++
		if f.Severity == lint.SeverityError {
			hasErrors = true
		}
	}

	if printed == 0 {
		fmt.Println("No lint findings.")
	}

	if hasErrors {
		return fmt.Errorf("lint found %d error(s)", countErrors(findings))
	}
	return nil
}

func countErrors(findings []lint.Finding) int {
	n := 0
	for _, f := range findings {
		if f.Severity == lint.SeverityError {
			n++
		}
	}
	return n
}
