// Package normalize provides utilities for normalizing .env file key-value pairs
// by applying consistent casing, trimming, and value transformations.
package normalize

import (
	"fmt"
	"strings"
)

// Options controls normalization behaviour.
type Options struct {
	// UppercaseKeys converts all keys to UPPER_CASE.
	UppercaseKeys bool
	// TrimValues strips leading and trailing whitespace from values.
	TrimValues bool
	// RemoveEmptyValues drops entries whose value is the empty string.
	RemoveEmptyValues bool
	// DryRun reports what would change without modifying the map.
	DryRun bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys: true,
		TrimValues:    true,
	}
}

// Change records a single normalization action.
type Change struct {
	Key    string
	Reason string
}

func (c Change) String() string {
	return fmt.Sprintf("%s: %s", c.Key, c.Reason)
}

// Result holds the output of a Normalize call.
type Result struct {
	Env     map[string]string
	Changes []Change
}

// Normalize applies opts to env and returns a Result.
// The original map is never mutated.
func Normalize(env map[string]string, opts Options) (*Result, error) {
	if env == nil {
		return nil, fmt.Errorf("normalize: env map must not be nil")
	}

	out := make(map[string]string, len(env))
	var changes []Change

	for k, v := range env {
		newKey := k
		newVal := v

		if opts.UppercaseKeys {
			up := strings.ToUpper(k)
			if up != k {
				changes = append(changes, Change{Key: k, Reason: fmt.Sprintf("key renamed to %s", up)})
				newKey = up
			}
		}

		if opts.TrimValues {
			trimmed := strings.TrimSpace(v)
			if trimmed != v {
				changes = append(changes, Change{Key: newKey, Reason: "value trimmed"})
				newVal = trimmed
			}
		}

		if opts.RemoveEmptyValues && newVal == "" {
			changes = append(changes, Change{Key: newKey, Reason: "removed empty value"})
			continue
		}

		if !opts.DryRun {
			out[newKey] = newVal
		} else {
			out[k] = v // preserve original in dry-run
		}
	}

	return &Result{Env: out, Changes: changes}, nil
}
