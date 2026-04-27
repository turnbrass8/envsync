package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envsync/internal/convert"
	"github.com/yourorg/envsync/internal/envfile"
)

func runConvert(args []string) error {
	fs := flag.NewFlagSet("convert", flag.ContinueOnError)
	format := fs.String("format", "dotenv", "output format: dotenv, json, yaml, toml")
	sorted := fs.Bool("sorted", true, "sort keys in output")
	prefix := fs.String("strip-prefix", "", "strip this prefix from key names before output")
	output := fs.String("output", "", "write output to file instead of stdout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	positional := fs.Args()
	if len(positional) < 1 {
		return fmt.Errorf("usage: envsync convert [flags] <envfile>")
	}

	env, err := envfile.Parse(positional[0])
	if err != nil {
		return fmt.Errorf("parsing env file: %w", err)
	}

	opts := convert.Options{
		Format:  convert.Format(*format),
		Sorted:  *sorted,
		Prefix:  *prefix,
	}

	result, err := convert.Convert(env, opts)
	if err != nil {
		return fmt.Errorf("converting: %w", err)
	}

	if *output != "" {
		if err := os.WriteFile(*output, []byte(result), 0o644); err != nil {
			return fmt.Errorf("writing output file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "written to %s\n", *output)
		return nil
	}

	fmt.Print(result)
	return nil
}
