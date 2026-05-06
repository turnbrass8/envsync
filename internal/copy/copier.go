// Package copy provides functionality to copy specific keys between .env files.
package copy

import (
	"fmt"
	"os"

	"github.com/user/envsync/internal/envfile"
)

// Options controls Copy behaviour.
type Options struct {
	Keys      []string // keys to copy; empty means all
	Overwrite bool     // replace existing keys in dst
	DryRun    bool     // report changes without writing
}

// Result summarises the outcome of a Copy operation.
type Result struct {
	Copied  []string
	Skipped []string
}

// Summary returns a human-readable one-liner.
func (r Result) Summary() string {
	return fmt.Sprintf("%d copied, %d skipped", len(r.Copied), len(r.Skipped))
}

// Copy reads keys from src and writes them into dst according to opts.
func Copy(src, dst string, opts Options) (Result, error) {
	srcEnv, err := envfile.Parse(src)
	if err != nil {
		return Result{}, fmt.Errorf("copy: read src %q: %w", src, err)
	}

	dstEnv, err := envfile.Parse(dst)
	if err != nil && !os.IsNotExist(err) {
		return Result{}, fmt.Errorf("copy: read dst %q: %w", dst, err)
	}
	if dstEnv == nil {
		dstEnv = &envfile.Env{}
	}

	keys := opts.Keys
	if len(keys) == 0 {
		keys = srcEnv.Keys()
	}

	var res Result
	for _, k := range keys {
		v, ok := srcEnv.Get(k)
		if !ok {
			continue
		}
		if _, exists := dstEnv.Get(k); exists && !opts.Overwrite {
			res.Skipped = append(res.Skipped, k)
			continue
		}
		if !opts.DryRun {
			dstEnv.Set(k, v)
		}
		res.Copied = append(res.Copied, k)
	}

	if !opts.DryRun && len(res.Copied) > 0 {
		if err := envfile.Write(dst, dstEnv); err != nil {
			return res, fmt.Errorf("copy: write dst %q: %w", dst, err)
		}
	}
	return res, nil
}
