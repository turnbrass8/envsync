package main

import (
	"flag"
	"fmt"
	"os"
)

// config holds all parsed CLI flags and arguments for the envsync command.
type config struct {
	// manifestPath is the path to the .envsync manifest file.
	manifestPath string

	// envPath is the path to the target .env file to diff/sync.
	envPath string

	// dryRun indicates whether to preview changes without writing to disk.
	dryRun bool

	// strict causes envsync to exit non-zero if any keys are missing or extra.
	strict bool

	// quiet suppresses informational output; only errors are printed.
	quiet bool

	// command is the subcommand to run: "diff" or "sync".
	command string
}

// usage prints a formatted help message to stderr.
func usage() {
	fmt.Fprintf(os.Stderr, `envsync — diff and sync .env files against a manifest

Usage:
  envsync <command> [flags]

Commands:
  diff    Compare a .env file against the manifest and report differences
  sync    Apply missing keys and defaults from the manifest to the .env file

Flags:
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
Examples:
  envsync diff --manifest .envsync --env .env
  envsync sync --manifest .envsync --env .env
  envsync sync --manifest .envsync --env .env --dry-run
`)
}

// parseFlags parses os.Args and returns a populated config.
// It exits with a non-zero status on invalid input.
func parseFlags() config {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	if cmd == "-h" || cmd == "--help" || cmd == "help" {
		usage()
		os.Exit(0)
	}

	switch cmd {
	case "diff", "sync":
		// valid subcommands
	default:
		fmt.Fprintf(os.Stderr, "error: unknown command %q\n\n", cmd)
		usage()
		os.Exit(1)
	}

	// Define a new FlagSet scoped to the subcommand so that --help works
	// correctly per subcommand without conflicting with global flag.CommandLine.
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	fs.Usage = usage

	var cfg config
	cfg.command = cmd

	fs.StringVar(&cfg.manifestPath, "manifest", ".envsync", "path to the envsync manifest file")
	fs.StringVar(&cfg.envPath, "env", ".env", "path to the .env file")
	fs.BoolVar(&cfg.dryRun, "dry-run", false, "preview changes without writing to disk (sync only)")
	fs.BoolVar(&cfg.strict, "strict", false, "exit non-zero if any keys are missing or extra")
	fs.BoolVar(&cfg.quiet, "quiet", false, "suppress informational output")

	// Parse everything after the subcommand token.
	if err := fs.Parse(os.Args[2:]); err != nil {
		// ExitOnError handles this, but be explicit.
		os.Exit(1)
	}

	if cfg.manifestPath == "" {
		fmt.Fprintln(os.Stderr, "error: --manifest path must not be empty")
		os.Exit(1)
	}
	if cfg.envPath == "" {
		fmt.Fprintln(os.Stderr, "error: --env path must not be empty")
		os.Exit(1)
	}

	return cfg
}
