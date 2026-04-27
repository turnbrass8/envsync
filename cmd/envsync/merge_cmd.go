package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/merge"
)

func runMerge(args []string) error {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	output := fs.String("out", "", "output file path (default: stdout)")
	strategyStr := fs.String("strategy", "last", "conflict resolution strategy: first|last|error")
	showConflicts := fs.Bool("show-conflicts", false, "print conflict summary to stderr")

	if err := fs.Parse(args); err != nil {
		return err
	}

	sources := fs.Args()
	if len(sources) < 2 {
		return fmt.Errorf("merge requires at least two source .env files")
	}

	var strategy merge.Strategy
	switch strings.ToLower(*strategyStr) {
	case "first":
		strategy = merge.StrategyFirst
	case "last":
		strategy = merge.StrategyLast
	case "error":
		strategy = merge.StrategyError
	default:
		return fmt.Errorf("unknown strategy %q: choose first, last, or error", *strategyStr)
	}

	result, err := merge.Merge(sources, strategy)
	if err != nil {
		return err
	}

	if *showConflicts && len(result.Conflicts) > 0 {
		fmt.Fprintf(os.Stderr, "merge conflicts (%d):\n", len(result.Conflicts))
		for _, c := range result.Conflicts {
			fmt.Fprintf(os.Stderr, "  %s: defined in [%s], chose %s\n",
				c.Key, strings.Join(c.Files, ", "), c.Chosen)
		}
	}

	var dest string
	if *output != "" {
		dest = *output
	}

	if dest == "" {
		for k, v := range result.Env {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
		return nil
	}

	return envfile.Write(dest, result.Env)
}
