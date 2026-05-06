// Package split provides functionality to split a .env file into multiple
// files based on key prefixes.
package split

import (
	"fmt"
	"strings"

	"github.com/user/envsync/internal/envfile"
)

// Options controls how the split is performed.
type Options struct {
	// Prefixes maps a prefix string to an output file path.
	// Keys matching a prefix are written to the corresponding file.
	Prefixes map[string]string
	// StripPrefix removes the prefix from keys before writing.
	StripPrefix bool
	// DryRun reports what would be written without modifying files.
	DryRun bool
}

// Result holds the outcome of a split operation.
type Result struct {
	// Written maps output file path to the number of keys written.
	Written map[string]int
	// Unmatched holds keys that did not match any prefix.
	Unmatched []string
}

// Summary returns a human-readable description of the result.
func (r Result) Summary() string {
	var sb strings.Builder
	for path, count := range r.Written {
		fmt.Fprintf(&sb, "%s: %d key(s)\n", path, count)
	}
	if len(r.Unmatched) > 0 {
		fmt.Fprintf(&sb, "unmatched: %s\n", strings.Join(r.Unmatched, ", "))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Split reads src and distributes keys into output files according to opts.
func Split(src string, opts Options) (Result, error) {
	if len(opts.Prefixes) == 0 {
		return Result{}, fmt.Errorf("split: at least one prefix mapping is required")
	}

	env, err := envfile.Parse(src)
	if err != nil {
		return Result{}, fmt.Errorf("split: parse %s: %w", src, err)
	}

	// Bucket keys by prefix.
	buckets := make(map[string]map[string]string) // outPath -> key -> value
	for outPath := range opts.Prefixes {
		buckets[outPath] = make(map[string]string)
	}

	result := Result{Written: make(map[string]int)}

	for _, key := range env.Keys() {
		val, _ := env.Get(key)
		matched := false
		for prefix, outPath := range opts.Prefixes {
			if strings.HasPrefix(key, prefix) {
				outKey := key
				if opts.StripPrefix {
					outKey = strings.TrimPrefix(key, prefix)
				}
				buckets[outPath][outKey] = val
				matched = true
				break
			}
		}
		if !matched {
			result.Unmatched = append(result.Unmatched, key)
		}
	}

	if opts.DryRun {
		for outPath, keys := range buckets {
			result.Written[outPath] = len(keys)
		}
		return result, nil
	}

	for outPath, keys := range buckets {
		if err := envfile.Write(outPath, keys); err != nil {
			return result, fmt.Errorf("split: write %s: %w", outPath, err)
		}
		result.Written[outPath] = len(keys)
	}

	return result, nil
}
