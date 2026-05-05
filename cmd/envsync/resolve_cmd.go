package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/resolve"
)

func runResolve(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("resolve", flag.ContinueOnError)
	var (
		files    = fs.String("files", "", "comma-separated list of .env files in precedence order (first = highest)")
		keys     = fs.String("keys", "", "comma-separated list of keys to resolve (empty = all)")
		strict   = fs.Bool("strict", false, "fail if any key is unresolved")
		fallback = fs.String("fallback", "", "value to use when a key is not found (non-strict mode)")
		origin   = fs.Bool("origin", false, "show source file for each resolved value")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *files == "" {
		return fmt.Errorf("resolve: --files is required")
	}

	filePaths := splitTrimmed(*files)
	sources := make([]resolve.Source, 0, len(filePaths))
	for _, path := range filePaths {
		env, err := envfile.Parse(path)
		if err != nil {
			return fmt.Errorf("resolve: reading %s: %w", path, err)
		}
		sources = append(sources, resolve.Source{Name: path, Values: env})
	}

	opts := resolve.Options{Strict: *strict, Fallback: *fallback}

	var results []resolve.Result
	if *keys == "" {
		results = resolve.ResolveAll(sources, opts)
	} else {
		ks := splitTrimmed(*keys)
		var err error
		results, err = resolve.Resolve(ks, sources, opts)
		if err != nil {
			return err
		}
	}

	for _, r := range results {
		if *origin {
			src := r.Source
			if src == "" {
				src = "<none>"
			}
			fmt.Fprintf(stdout, "%s=%s  # %s\n", r.Key, r.Value, src)
		} else {
			fmt.Fprintf(stdout, "%s=%s\n", r.Key, r.Value)
		}
	}
	return nil
}

func splitTrimmed(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

var _ = os.Stderr // ensure os import used
