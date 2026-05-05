package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envsync/internal/sort"
)

func runSort(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("sort", flag.ContinueOnError)
	fs.SetOutput(stdout)

	strategy := fs.String("strategy", "alpha", "sort strategy: alpha | reverse")
	dryRun := fs.Bool("dry-run", false, "preview changes without writing")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("sort: usage: envsync sort [flags] <file>")
	}

	path := fs.Arg(0)

	opts := sort.Options{
		Strategy: sort.Strategy(*strategy),
		DryRun:   *dryRun,
	}

	res, err := sort.Sort(path, opts)
	if err != nil {
		return err
	}

	if *dryRun {
		fmt.Fprintf(stdout, "[dry-run] %s: %d key(s) would be reordered (strategy=%s)\n",
			res.Path, res.Reordered, *strategy)
		return nil
	}

	if res.Reordered == 0 {
		fmt.Fprintf(stdout, "%s: already sorted\n", res.Path)
	} else {
		fmt.Fprintf(stdout, "%s: sorted %d key(s) (strategy=%s)\n",
			res.Path, res.Reordered, *strategy)
	}
	return nil
}

// selfTest entry point wired in main.go switch.
var _ = func() bool {
	_ = runSort // ensure symbol is referenced
	return true
}()

func init() {
	if len(os.Args) > 1 && os.Args[1] == "sort" {
		if err := runSort(os.Args[2:], os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
