package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envsync/internal/envfile"
	"github.com/yourorg/envsync/internal/export"
)

// runExport handles the `envsync export` sub-command.
// It reads an env file and writes it to stdout in the requested format.
func runExport(args []string) error {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	formatFlag := fs.String("format", "dotenv", "output format: dotenv | shell | json")
	sortedFlag := fs.Bool("sorted", true, "sort keys alphabetically")
	envFlag := fs.String("env", ".env", "path to the .env file to export")

	if err := fs.Parse(args); err != nil {
		return err
	}

	f, err := os.Open(*envFlag)
	if err != nil {
		return fmt.Errorf("opening env file %q: %w", *envFlag, err)
	}
	defer f.Close()

	env, err := envfile.Parse(f)
	if err != nil {
		return fmt.Errorf("parsing env file: %w", err)
	}

	var fmt_ export.Format
	switch *formatFlag {
	case "shell":
		fmt_ = export.FormatShell
	case "json":
		fmt_ = export.FormatJSON
	default:
		fmt_ = export.FormatDotenv
	}

	return export.Export(env, export.Options{
		Format:  fmt_,
		Sorted:  *sortedFlag,
		Writer:  os.Stdout,
	})
}
