// Package strip removes keys from a .env file based on a list of patterns or exact names.
package strip

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/user/envsync/internal/envfile"
)

// Options controls Strip behaviour.
type Options struct {
	// Keys holds exact key names to remove.
	Keys []string
	// Patterns holds regular expressions; any matching key is removed.
	Patterns []string
	// DryRun reports what would be removed without writing.
	DryRun bool
}

// Result describes what was stripped.
type Result struct {
	Removed []string
}

// Strip reads src, removes matching keys, and writes the result to dst.
// When Options.DryRun is true the destination file is not modified.
func Strip(src, dst string, opts Options) (Result, error) {
	env, err := envfile.Parse(src)
	if err != nil {
		return Result{}, fmt.Errorf("strip: parse %q: %w", src, err)
	}

	matchers, err := compilePatterns(opts.Patterns)
	if err != nil {
		return Result{}, fmt.Errorf("strip: compile patterns: %w", err)
	}

	exact := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		exact[strings.TrimSpace(k)] = true
	}

	var removed []string
	filtered := make(map[string]string)

	for _, k := range env.Keys() {
		v, _ := env.Get(k)
		if exact[k] || matchesAny(k, matchers) {
			removed = append(removed, k)
			continue
		}
		filtered[k] = v
	}

	if !opts.DryRun {
		if err := envfile.Write(dst, filtered, env.Keys()); err != nil {
			return Result{}, fmt.Errorf("strip: write %q: %w", dst, err)
		}
	}

	_ = os.Stderr // satisfy import if needed
	return Result{Removed: removed}, nil
}

func compilePatterns(patterns []string) ([]*regexp.Regexp, error) {
	var out []*regexp.Regexp
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", p, err)
		}
		out = append(out, re)
	}
	return out, nil
}

func matchesAny(key string, patterns []*regexp.Regexp) bool {
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
