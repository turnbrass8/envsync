package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/your-org/envsync/internal/patch"
)

func runPatch(args []string) error {
	fs := flag.NewFlagSet("patch", flag.ContinueOnError)
	var (
		rulesFile = fs.String("rules", "", "file containing patch operations (KEY=VALUE or -KEY)")
		set       = fs.String("set", "", "comma-separated KEY=VALUE pairs to set")
		del       = fs.String("delete", "", "comma-separated keys to delete")
		dryRun    = fs.Bool("dry-run", false, "print summary without modifying the file")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("patch: usage: envsync patch [flags] <envfile>")
	}
	envPath := fs.Arg(0)

	var ops []patch.Op

	// Load from rules file.
	if *rulesFile != "" {
		f, err := os.Open(*rulesFile)
		if err != nil {
			return fmt.Errorf("patch: open rules file: %w", err)
		}
		defer f.Close()
		parsed, err := patch.ParseRules(f)
		if err != nil {
			return err
		}
		ops = append(ops, parsed...)
	}

	// Inline --set flag.
	if *set != "" {
		for _, pair := range strings.Split(*set, ",") {
			pair = strings.TrimSpace(pair)
			idx := strings.IndexByte(pair, '=')
			if idx < 0 {
				return fmt.Errorf("patch: --set value %q must be KEY=VALUE", pair)
			}
			ops = append(ops, patch.Op{Key: pair[:idx], Value: pair[idx+1:]})
		}
	}

	// Inline --delete flag.
	if *del != "" {
		for _, k := range strings.Split(*del, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				ops = append(ops, patch.Op{Key: k, Delete: true})
			}
		}
	}

	res, err := patch.Patch(envPath, ops, *dryRun)
	if err != nil {
		return err
	}

	prefix := ""
	if *dryRun {
		prefix = "[dry-run] "
	}
	fmt.Printf("%s%s\n", prefix, res.Summary())
	return nil
}
