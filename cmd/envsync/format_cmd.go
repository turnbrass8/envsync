package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envsync/internal/format"
)

func runFormat(args []string) error {
	fs := flag.NewFlagSet("format", flag.ContinueOnError)
	quoteStyle := fs.String("quote", "double", "quoting style: none, single, double")
	spaceEq := fs.Bool("space-eq", false, "add spaces around '='")
	dryRun := fs.Bool("dry-run", false, "print result without writing")

	if err := fs.Parse(args); err != nil {
		return err
	}

	paths := fs.Args()
	if len(paths) == 0 {
		return fmt.Errorf("format: at least one .env file required")
	}

	opts := format.Options{
		QuoteStyle:        *quoteStyle,
		SpaceAroundEquals: *spaceEq,
		DryRun:            *dryRun,
	}

	exitCode := 0
	for _, p := range paths {
		res, err := format.Format(p, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			exitCode = 1
			continue
		}
		if *dryRun {
			fmt.Print(res.Formatted)
		} else if res.Changed {
			fmt.Printf("formatted: %s\n", res.Path)
		} else {
			fmt.Printf("ok: %s\n", res.Path)
		}
	}

	if exitCode != 0 {
		return fmt.Errorf("format: one or more files failed")
	}
	return nil
}
