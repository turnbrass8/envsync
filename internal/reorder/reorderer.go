// Package reorder provides functionality to reorder keys in a .env file
// according to a specified ordering manifest or reference file.
package reorder

import (
	"fmt"
	"os"

	"github.com/user/envsync/internal/envfile"
)

// Options controls the behaviour of Reorder.
type Options struct {
	// Order defines the desired key sequence. Keys not listed appear at the end.
	Order []string
	// DryRun reports what would change without writing to disk.
	DryRun bool
	// Append controls whether unlisted keys are appended (true) or dropped (false).
	Append bool
}

// Result summarises the outcome of a reorder operation.
type Result struct {
	Reordered int
	Appended  int
	Dropped   int
}

// Summary returns a human-readable one-liner.
func (r Result) Summary() string {
	return fmt.Sprintf("reordered=%d appended=%d dropped=%d", r.Reordered, r.Appended, r.Dropped)
}

// Reorder reads path, reorders its keys according to opts.Order and writes the
// result back. When opts.DryRun is true the file is not modified.
func Reorder(path string, opts Options) (Result, error) {
	env, err := envfile.Parse(path)
	if err != nil {
		return Result{}, fmt.Errorf("reorder: parse %q: %w", path, err)
	}

	ordered := make([]string, 0, len(opts.Order))
	seen := make(map[string]bool)

	for _, key := range opts.Order {
		if _, ok := env.Get(key); ok {
			ordered = append(ordered, key)
			seen[key] = true
		}
	}

	var appended []string
	if opts.Append {
		for _, key := range env.Keys() {
			if !seen[key] {
				appended = append(appended, key)
			}
		}
	}

	dropped := 0
	for _, key := range env.Keys() {
		if !seen[key] && !opts.Append {
			dropped++
		}
	}

	res := Result{
		Reordered: len(ordered),
		Appended:  len(appended),
		Dropped:   dropped,
	}

	if opts.DryRun {
		return res, nil
	}

	final := append(ordered, appended...)
	out := make(map[string]string, len(final))
	for _, k := range final {
		v, _ := env.Get(k)
		out[k] = v
	}

	if err := envfile.Write(path, out, final); err != nil {
		return Result{}, fmt.Errorf("reorder: write %q: %w", path, err)
	}

	_ = os.Stdout // satisfy import if needed
	return res, nil
}
