package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	envcopy "github.com/user/envsync/internal/copy"
)

func runCopy(args []string) error {
	fs := flag.NewFlagSet("copy", flag.ContinueOnError)
	keys := fs.String("keys", "", "comma-separated list of keys to copy (default: all)")
	overwrite := fs.Bool("overwrite", false, "overwrite existing keys in destination")
	dryRun := fs.Bool("dry-run", false, "print changes without writing")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return fmt.Errorf("copy: usage: copy [flags] <src> <dst>")
	}

	src := fs.Arg(0)
	dst := fs.Arg(1)

	var keyList []string
	if *keys != "" {
		for _, k := range strings.Split(*keys, ",") {
			if k = strings.TrimSpace(k); k != "" {
				keyList = append(keyList, k)
			}
		}
	}

	res, err := envcopy.Copy(src, dst, envcopy.Options{
		Keys:      keyList,
		Overwrite: *overwrite,
		DryRun:    *dryRun,
	})
	if err != nil {
		return err
	}

	for _, k := range res.Copied {
		if *dryRun {
			fmt.Fprintf(os.Stdout, "[dry-run] would copy %s\n", k)
		} else {
			fmt.Fprintf(os.Stdout, "copied %s\n", k)
		}
	}
	for _, k := range res.Skipped {
		fmt.Fprintf(os.Stdout, "skipped %s (already exists)\n", k)
	}
	fmt.Fprintln(os.Stdout, res.Summary())
	return nil
}
