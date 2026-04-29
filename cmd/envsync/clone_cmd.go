package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsync/internal/clone"
)

// runClone implements the `clone` sub-command.
// Usage: envsync clone -src=.env.prod -dst=.env.staging FOO BAR:NEW_BAR
func runClone(args []string) error {
	fs := flag.NewFlagSet("clone", flag.ContinueOnError)
	src := fs.String("src", "", "source env file (required)")
	dst := fs.String("dst", "", "destination env file (required)")
	overwrite := fs.Bool("overwrite", false, "overwrite existing keys in destination")
	dryRun := fs.Bool("dry-run", false, "print what would change without writing")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *src == "" || *dst == "" {
		return fmt.Errorf("clone: -src and -dst are required")
	}

	if fs.NArg() == 0 {
		return fmt.Errorf("clone: at least one KEY or KEY:NEW_KEY argument required")
	}

	var rules []clone.Rule
	for _, arg := range fs.Args() {
		parts := strings.SplitN(arg, ":", 2)
		r := clone.Rule{SrcKey: parts[0]}
		if len(parts) == 2 {
			r.DstKey = parts[1]
		}
		rules = append(rules, r)
	}

	res, err := clone.Clone(*src, *dst, clone.Options{
		Rules:     rules,
		Overwrite: *overwrite,
		DryRun:    *dryRun,
	})
	if err != nil {
		return err
	}

	for _, k := range res.Copied {
		fmt.Fprintf(os.Stdout, "copied:  %s\n", k)
	}
	for _, k := range res.Skipped {
		fmt.Fprintf(os.Stdout, "skipped: %s (already exists)\n", k)
	}

	if *dryRun {
		fmt.Fprintln(os.Stdout, "dry-run: no changes written")
	}

	return nil
}
