// Package filter provides utilities for selecting a subset of env entries
// based on key prefixes, patterns, or explicit key lists.
package filter

import (
	"fmt"
	"regexp"
	"strings"
)

// Options controls how entries are filtered.
type Options struct {
	// Keys is an explicit list of keys to include. If non-empty, only these keys
	// are returned (unless also excluded by Exclude).
	Keys []string

	// Prefix restricts results to keys that start with the given string.
	Prefix string

	// Pattern is a regular expression that key names must match.
	Pattern string

	// Exclude is a list of exact key names to drop from the result.
	Exclude []string
}

// Filter returns a new map containing only the entries from src that satisfy
// the provided Options. An error is returned if Pattern is not a valid regex.
func Filter(src map[string]string, opts Options) (map[string]string, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, fmt.Errorf("filter: invalid pattern %q: %w", opts.Pattern, err)
		}
	}

	allowSet := buildSet(opts.Keys)
	excludeSet := buildSet(opts.Exclude)

	out := make(map[string]string, len(src))
	for k, v := range src {
		if len(allowSet) > 0 && !allowSet[k] {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		if re != nil && !re.MatchString(k) {
			continue
		}
		if excludeSet[k] {
			continue
		}
		out[k] = v
	}
	return out, nil
}

// StripPrefix removes the given prefix from all keys in src, returning a new
// map. Keys that do not carry the prefix are passed through unchanged.
func StripPrefix(src map[string]string, prefix string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[strings.TrimPrefix(k, prefix)] = v
	}
	return out
}

// Keys returns a sorted slice of all keys present in src that would be
// retained by Filter with the given opts. It is a convenience wrapper for
// callers that only need the key names rather than the full map.
func Keys(src map[string]string, opts Options) ([]string, error) {
	filtered, err := Filter(src, opts)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(filtered))
	for k := range filtered {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}

func buildSet(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
