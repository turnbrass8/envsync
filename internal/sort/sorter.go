// Package sort provides utilities for sorting .env file keys
// alphabetically or by custom ordering rules.
package sort

import (
	"fmt"
	"os"
	gosort "sort"
	"strings"

	"github.com/user/envsync/internal/envfile"
)

// Strategy controls how keys are ordered.
type Strategy string

const (
	StrategyAlpha  Strategy = "alpha"  // ascending alphabetical
	StrategyReverse Strategy = "reverse" // descending alphabetical
)

// Options configures a Sort operation.
type Options struct {
	Strategy Strategy
	DryRun   bool
}

// Result holds the outcome of a sort operation.
type Result struct {
	Path     string
	Reordered int // number of keys that moved position
}

// Sort reads the .env file at path, reorders its keys according to opts,
// and writes the result back (unless DryRun is set).
func Sort(path string, opts Options) (Result, error) {
	if opts.Strategy == "" {
		opts.Strategy = StrategyAlpha
	}

	env, err := envfile.Parse(path)
	if err != nil {
		return Result{}, fmt.Errorf("sort: parse %q: %w", path, err)
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	originalOrder := make([]string, len(keys))
	copy(originalOrder, keys)
	gosort.Strings(originalOrder)

	switch opts.Strategy {
	case StrategyAlpha:
		gosort.Strings(keys)
	case StrategyReverse:
		gosort.Sort(gosort.Reverse(gosort.StringSlice(keys)))
	default:
		return Result{}, fmt.Errorf("sort: unknown strategy %q", opts.Strategy)
	}

	reordered := countReordered(originalOrder, keys)

	if opts.DryRun {
		return Result{Path: path, Reordered: reordered}, nil
	}

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
	}

	if err := os.WriteFile(path, []byte(sb.String()), 0o644); err != nil {
		return Result{}, fmt.Errorf("sort: write %q: %w", path, err)
	}

	return Result{Path: path, Reordered: reordered}, nil
}

func countReordered(before, after []string) int {
	count := 0
	for i := range before {
		if i < len(after) && before[i] != after[i] {
			count++
		}
	}
	return count
}
