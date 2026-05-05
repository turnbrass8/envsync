// Package defaults applies missing env keys from a set of default values,
// optionally overwriting existing ones.
package defaults

import (
	"fmt"
	"sort"
)

// Options controls how defaults are applied.
type Options struct {
	// Overwrite replaces existing keys with default values.
	Overwrite bool
	// DryRun reports what would change without modifying the map.
	DryRun bool
}

// Result describes the outcome of applying defaults.
type Result struct {
	Applied  []string
	Skipped  []string
	Overwritten []string
}

// Summary returns a human-readable one-liner.
func (r Result) Summary() string {
	return fmt.Sprintf("applied=%d skipped=%d overwritten=%d",
		len(r.Applied), len(r.Skipped), len(r.Overwritten))
}

// Apply merges defaults into env according to opts.
// env is modified in-place unless DryRun is true.
func Apply(env map[string]string, defaults map[string]string, opts Options) (Result, error) {
	if env == nil {
		return Result{}, fmt.Errorf("env map must not be nil")
	}

	keys := make([]string, 0, len(defaults))
	for k := range defaults {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var res Result
	for _, k := range keys {
		v := defaults[k]
		existing, exists := env[k]
		switch {
		case !exists:
			if !opts.DryRun {
				env[k] = v
			}
			res.Applied = append(res.Applied, k)
		case exists && opts.Overwrite && existing != v:
			if !opts.DryRun {
				env[k] = v
			}
			res.Overwritten = append(res.Overwritten, k)
		default:
			res.Skipped = append(res.Skipped, k)
		}
	}
	return res, nil
}
