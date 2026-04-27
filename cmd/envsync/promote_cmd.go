package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envsync/internal/manifest"
	"github.com/yourorg/envsync/internal/promote"
)

func runPromote(args []string) error {
	fs := flag.NewFlagSet("promote", flag.ContinueOnError)
	manifestPath := fs.String("manifest", ".envsync", "path to manifest file")
	overwrite := fs.Bool("overwrite", false, "overwrite existing keys in destination")
	dryRun := fs.Bool("dry-run", false, "preview changes without writing")

	if err := fs.Parse(args); err != nil {
		return err
	}

	positional := fs.Args()
	if len(positional) < 2 {
		return fmt.Errorf("usage: envsync promote [flags] <src.env> <dst.env>")
	}

	srcPath := positional[0]
	dstPath := positional[1]

	mf, err := manifest.Parse(*manifestPath)
	if err != nil {
		return fmt.Errorf("promote: load manifest: %w", err)
	}

	opts := promote.Options{
		DryRun:    *dryRun,
		Overwrite: *overwrite,
	}

	results, err := promote.Promote(srcPath, dstPath, mf, opts)
	if err != nil {
		return err
	}

	promoted, skipped := 0, 0
	for _, r := range results {
		switch {
		case r.Promoted:
			promoted++
			if *dryRun {
				fmt.Fprintf(os.Stdout, "[dry-run] would promote: %s\n", r.Key)
			} else {
				fmt.Fprintf(os.Stdout, "promoted: %s\n", r.Key)
			}
		case r.Skipped:
			skipped++
			fmt.Fprintf(os.Stdout, "skipped:  %s (%s)\n", r.Key, r.Reason)
		}
	}

	fmt.Fprintf(os.Stdout, "\npromote complete: %d promoted, %d skipped\n", promoted, skipped)
	return nil
}
