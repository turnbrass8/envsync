// Package dedupe removes duplicate keys from .env files,
// keeping either the first or last occurrence of each key.
package dedupe

import (
	"fmt"
	"os"

	"github.com/user/envsync/internal/envfile"
)

// Strategy controls which occurrence of a duplicate key is retained.
type Strategy string

const (
	StrategyFirst Strategy = "first" // keep the first occurrence
	StrategyLast  Strategy = "last"  // keep the last occurrence (default)
)

// Options configures the Dedupe operation.
type Options struct {
	Strategy Strategy
	DryRun   bool
}

// Result describes the outcome of a Dedupe run.
type Result struct {
	Removed []string // keys that had duplicates removed
}

// Dedupe reads the .env file at path, removes duplicate keys according to
// opts, and writes the result back unless DryRun is set.
func Dedupe(path string, opts Options) (*Result, error) {
	if opts.Strategy == "" {
		opts.Strategy = StrategyLast
	}

	env, err := envfile.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("dedupe: parse %q: %w", path, err)
	}

	seen := make(map[string]int)   // key -> index in ordered slice
	type kv struct{ k, v string }
	var ordered []kv
	duplicates := make(map[string]bool)

	for _, key := range env.Keys() {
		val, _ := env.Get(key)
		if idx, exists := seen[key]; exists {
			duplicates[key] = true
			if opts.Strategy == StrategyLast {
				ordered[idx].v = val
			}
			// StrategyFirst: ignore the new value
		} else {
			seen[key] = len(ordered)
			ordered = append(ordered, kv{key, val})
		}
	}

	result := &Result{}
	for k := range duplicates {
		result.Removed = append(result.Removed, k)
	}

	if opts.DryRun || len(duplicates) == 0 {
		return result, nil
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("dedupe: create %q: %w", path, err)
	}
	defer f.Close()

	for _, pair := range ordered {
		if _, err := fmt.Fprintf(f, "%s=%s\n", pair.k, pair.v); err != nil {
			return nil, fmt.Errorf("dedupe: write: %w", err)
		}
	}

	return result, nil
}
