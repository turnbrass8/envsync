package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsync/internal/reorder"
)

func runReorder(args []string) error {
	fs := flag.NewFlagSet("reorder", flag.ContinueOnError)
	orderFlag := fs.String("order", "", "comma-separated list of keys in desired order (required)")
	dryRun := fs.Bool("dry-run", false, "print what would change without writing")
	append_ := fs.Bool("append", true, "append unlisted keys at the end; set false to drop them")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *orderFlag == "" {
		return fmt.Errorf("reorder: --order flag is required")
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("reorder: target .env file argument is required")
	}

	keys := splitTrimmedReorder(*orderFlag)
	path := fs.Arg(0)

	res, err := reorder.Reorder(path, reorder.Options{
		Order:  keys,
		DryRun: *dryRun,
		Append: *append_,
	})
	if err != nil {
		return err
	}

	if *dryRun {
		fmt.Fprintf(os.Stdout, "[dry-run] %s\n", res.Summary())
	} else {
		fmt.Fprintf(os.Stdout, "%s\n", res.Summary())
	}
	return nil
}

func splitTrimmedReorder(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
