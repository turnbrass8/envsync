package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/yourorg/envsync/internal/envfile"
	"github.com/yourorg/envsync/internal/interpolate"
)

func runInterpolate(args []string) error {
	fs := flag.NewFlagSet("interpolate", flag.ContinueOnError)
	envPath := fs.String("env", ".env", "Path to the .env file to interpolate")
	outPath := fs.String("out", "", "Output file (defaults to stdout)")
	strict := fs.Bool("strict", true, "Fail on unresolved references")

	if err := fs.Parse(args); err != nil {
		return err
	}

	env, err := envfile.Parse(*envPath)
	if err != nil {
		return fmt.Errorf("interpolate: parse %q: %w", *envPath, err)
	}

	if *strict {
		if err := interpolate.ResolveAll(env); err != nil {
			return fmt.Errorf("interpolate: %w", err)
		}
	} else {
		for k, v := range env {
			resolved, _ := interpolate.Resolve(v, env)
			env[k] = resolved
		}
	}

	out := os.Stdout
	if *outPath != "" {
		f, err := os.Create(*outPath)
		if err != nil {
			return fmt.Errorf("interpolate: create output: %w", err)
		}
		defer f.Close()
		out = f
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(out, "%s=%s\n", k, env[k])
	}
	return nil
}
