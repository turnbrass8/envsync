package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsync/internal/split"
)

func runSplit(args []string) error {
	fs := flag.NewFlagSet("split", flag.ContinueOnError)
	var (
		src         = fs.String("src", ".env", "source .env file to split")
		mappings    = fs.String("map", "", "prefix:file mappings, comma-separated (e.g. APP_:app.env,DB_:db.env)")
		stripPrefix = fs.Bool("strip-prefix", false, "strip matched prefix from keys in output files")
		dryRun      = fs.Bool("dry-run", false, "print what would be written without modifying files")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *mappings == "" {
		return fmt.Errorf("split: --map is required (e.g. APP_:app.env,DB_:db.env)")
	}

	prefixes, err := parseSplitMappings(*mappings)
	if err != nil {
		return err
	}

	res, err := split.Split(*src, split.Options{
		Prefixes:    prefixes,
		StripPrefix: *stripPrefix,
		DryRun:      *dryRun,
	})
	if err != nil {
		return err
	}

	for outPath, count := range res.Written {
		if *dryRun {
			fmt.Fprintf(os.Stdout, "[dry-run] %s: %d key(s)\n", outPath, count)
		} else {
			fmt.Fprintf(os.Stdout, "wrote %s: %d key(s)\n", outPath, count)
		}
	}
	if len(res.Unmatched) > 0 {
		fmt.Fprintf(os.Stdout, "unmatched keys: %s\n", strings.Join(res.Unmatched, ", "))
	}
	return nil
}

func parseSplitMappings(raw string) (map[string]string, error) {
	prefixes := make(map[string]string)
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		idx := strings.Index(part, ":")
		if idx <= 0 {
			return nil, fmt.Errorf("split: invalid mapping %q (expected PREFIX:file)", part)
		}
		prefix := part[:idx]
		file := part[idx+1:]
		if file == "" {
			return nil, fmt.Errorf("split: empty file path for prefix %q", prefix)
		}
		prefixes[prefix] = file
	}
	if len(prefixes) == 0 {
		return nil, fmt.Errorf("split: no valid mappings found in %q", raw)
	}
	return prefixes, nil
}
