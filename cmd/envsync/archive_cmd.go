package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsync/internal/archive"
)

func runArchive(args []string) error {
	fs := flag.NewFlagSet("archive", flag.ContinueOnError)
	dest := fs.String("out", "envsync-archive.zip", "destination zip file")
	labelsRaw := fs.String("labels", "", "comma-separated key=value labels (e.g. env=prod,version=1)")
	dryRun := fs.Bool("dry-run", false, "show what would be archived without writing")

	if err := fs.Parse(args); err != nil {
		return err
	}

	envFiles := fs.Args()
	if len(envFiles) == 0 {
		return fmt.Errorf("archive: at least one .env file must be specified")
	}

	labels := parseLabels(*labelsRaw)

	opts := archive.Options{
		Labels: labels,
		DryRun: *dryRun,
	}

	meta, err := archive.Archive(*dest, envFiles, opts)
	if err != nil {
		return err
	}

	if *dryRun {
		fmt.Fprintf(os.Stdout, "[dry-run] would archive %d file(s) to %s\n", len(meta.Files), *dest)
		for _, f := range meta.Files {
			fmt.Fprintf(os.Stdout, "  - %s\n", f)
		}
		return nil
	}

	fmt.Fprintf(os.Stdout, "archived %d file(s) to %s\n", len(meta.Files), *dest)
	for _, f := range meta.Files {
		fmt.Fprintf(os.Stdout, "  + %s\n", f)
	}
	return nil
}

func parseLabels(raw string) map[string]string {
	if raw == "" {
		return nil
	}
	m := map[string]string{}
	for _, pair := range strings.Split(raw, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) == 2 {
			m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return m
}
