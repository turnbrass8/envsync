package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envsync/internal/env"
)

func runEnv(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("env", flag.ContinueOnError)
	envFile := fs.String("f", ".env", "path to .env file")
	showOrigin := fs.Bool("origin", false, "show which source each key came from")
	filterKey := fs.String("key", "", "show value for a specific key only")
	noOS := fs.Bool("no-os", false, "exclude OS environment variables")

	if err := fs.Parse(args); err != nil {
		return err
	}

	var sources []env.Source

	if _, err := os.Stat(*envFile); err == nil {
		src, err := env.FileSource("file", *envFile)
		if err != nil {
			return fmt.Errorf("parse %s: %w", *envFile, err)
		}
		sources = append(sources, src)
	}

	if !*noOS {
		sources = append(sources, env.OSSource())
	}

	loader := env.NewLoader(sources...)

	if *filterKey != "" {
		key := strings.ToUpper(*filterKey)
		resolved := loader.Resolve()
		val, ok := resolved[key]
		if !ok {
			return fmt.Errorf("key %q not found", key)
		}
		if *showOrigin {
			fmt.Fprintf(stdout, "%s=%s\t[%s]\n", key, val, loader.Origin(key))
		} else {
			fmt.Fprintln(stdout, val)
		}
		return nil
	}

	resolved := loader.Resolve()
	for _, k := range loader.Keys() {
		if *showOrigin {
			fmt.Fprintf(stdout, "%s=%s\t[%s]\n", k, resolved[k], loader.Origin(k))
		} else {
			fmt.Fprintf(stdout, "%s=%s\n", k, resolved[k])
		}
	}
	return nil
}
