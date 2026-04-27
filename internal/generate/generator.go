// Package generate provides utilities for generating .env files from
// a manifest template, optionally pre-populating values from an existing
// environment file or from system environment variables.
package generate

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/user/envsync/internal/manifest"
)

// Options controls how the generator produces output.
type Options struct {
	// IncludeDefaults pre-fills keys that have a default value defined in the
	// manifest.
	IncludeDefaults bool

	// IncludeSystem pulls values from the current process environment for any
	// key that is not already resolved from a source file.
	IncludeSystem bool

	// Overwrite controls whether an existing destination file is replaced.
	Overwrite bool

	// CommentMissing adds a comment marker next to keys whose values could not
	// be resolved, making it easy to spot gaps.
	CommentMissing bool
}

// Result describes the outcome of a generation run.
type Result struct {
	// Written is the number of keys written to the output.
	Written int

	// Missing contains the key names that had no resolved value.
	Missing []string
}

// Generate creates a new .env file at dst based on the provided manifest and
// an optional source environment map. Keys are written in the order they appear
// in the manifest.
func Generate(mf *manifest.Manifest, src map[string]string, dst string, opts Options) (*Result, error) {
	if !opts.Overwrite {
		if _, err := os.Stat(dst); err == nil {
			return nil, fmt.Errorf("destination file already exists: %s (use overwrite option to replace)", dst)
		}
	}

	var sb strings.Builder
	result := &Result{}

	// Sort keys for deterministic output.
	keys := make([]string, 0, len(mf.Keys))
	for k := range mf.Keys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		entry := mf.Keys[key]
		value, resolved := resolve(key, entry, src, opts)

		if entry.Description != "" {
			fmt.Fprintf(&sb, "# %s\n", entry.Description)
		}

		if !resolved {
			result.Missing = append(result.Missing, key)
			if opts.CommentMissing {
				fmt.Fprintf(&sb, "# MISSING: %s=\n", key)
			} else {
				fmt.Fprintf(&sb, "%s=\n", key)
			}
			continue
		}

		fmt.Fprintf(&sb, "%s=%s\n", key, quoteIfNeeded(value))
		result.Written++
	}

	if err := os.WriteFile(dst, []byte(sb.String()), 0o644); err != nil {
		return nil, fmt.Errorf("writing generated file: %w", err)
	}

	return result, nil
}

// resolve determines the value for a single key using the following priority:
//  1. Explicit value in src map.
//  2. System environment variable (when IncludeSystem is enabled).
//  3. Default value from the manifest entry (when IncludeDefaults is enabled).
func resolve(key string, entry manifest.Entry, src map[string]string, opts Options) (string, bool) {
	if v, ok := src[key]; ok {
		return v, true
	}

	if opts.IncludeSystem {
		if v, ok := os.LookupEnv(key); ok {
			return v, true
		}
	}

	if opts.IncludeDefaults && entry.Default != "" {
		return entry.Default, true
	}

	return "", false
}

// quoteIfNeeded wraps the value in double quotes when it contains whitespace
// or special shell characters.
func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t\n#$\"\'\\`") {
		escaped := strings.ReplaceAll(v, `"`, `\"`)
		return `"` + escaped + `"`
	}
	return v
}
