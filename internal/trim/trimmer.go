// Package trim provides utilities to remove unused or stale keys
// from a .env file based on a manifest of expected keys.
package trim

import (
	"fmt"
	"sort"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/manifest"
)

// Result holds the outcome of a trim operation.
type Result struct {
	Removed []string
	Retained map[string]string
	DryRun bool
}

// Summary returns a human-readable summary of the trim result.
func (r *Result) Summary() string {
	if len(r.Removed) == 0 {
		return "no stale keys found"
	}
	verb := "would remove"
	if !r.DryRun {
		verb = "removed"
	}
	return fmt.Sprintf("%s %d stale key(s): %v", verb, len(r.Removed), r.Removed)
}

// Options controls the behaviour of Trim.
type Options struct {
	// DryRun prevents any writes when true.
	DryRun bool
}

// Trim reads envPath, removes keys not listed in mf, and (unless DryRun)
// writes the result back to envPath. It returns a Result describing what
// was changed.
func Trim(envPath string, mf *manifest.Manifest, opts Options) (*Result, error) {
	env, err := envfile.Parse(envPath)
	if err != nil {
		return nil, fmt.Errorf("trim: parse env: %w", err)
	}

	// Build a set of known keys from the manifest.
	known := make(map[string]struct{}, len(mf.Keys))
	for _, k := range mf.Keys {
		known[k.Name] = struct{}{}
	}

	retained := make(map[string]string)
	var removed []string

	for k, v := range env.Values {
		if _, ok := known[k]; ok {
			retained[k] = v
		} else {
			removed = append(removed, k)
		}
	}
	sort.Strings(removed)

	res := &Result{
		Removed:  removed,
		Retained: retained,
		DryRun:   opts.DryRun,
	}

	if !opts.DryRun && len(removed) > 0 {
		if err := envfile.Write(envPath, retained); err != nil {
			return nil, fmt.Errorf("trim: write env: %w", err)
		}
	}

	return res, nil
}
