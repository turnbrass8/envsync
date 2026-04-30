// Package sanitize provides utilities for sanitizing .env file values,
// trimming whitespace, normalizing line endings, and removing control characters.
package sanitize

import (
	"strings"
	"unicode"
)

// Result holds the sanitized key-value pairs and a list of changes made.
type Result struct {
	Env     map[string]string
	Changes []Change
}

// Change records a single sanitization action applied to a key.
type Change struct {
	Key    string
	Before string
	After  string
	Reason string
}

// Options controls which sanitization passes are applied.
type Options struct {
	TrimSpace       bool
	RemoveControl   bool
	NormalizeNewlines bool
	UppercaseKeys   bool
}

// DefaultOptions returns a sensible default set of sanitization options.
func DefaultOptions() Options {
	return Options{
		TrimSpace:         true,
		RemoveControl:     true,
		NormalizeNewlines: true,
		UppercaseKeys:     false,
	}
}

// Sanitize applies the given options to each key-value pair in env and
// returns a Result containing the cleaned map and a log of all changes.
func Sanitize(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))
	var changes []Change

	for k, v := range env {
		origKey := k
		origVal := v

		if opts.UppercaseKeys {
			k = strings.ToUpper(k)
		}

		if opts.NormalizeNewlines {
			v = strings.ReplaceAll(v, "\r\n", "\n")
			v = strings.ReplaceAll(v, "\r", "\n")
		}

		if opts.RemoveControl {
			v = strings.Map(func(r rune) rune {
				if unicode.IsControl(r) && r != '\n' && r != '\t' {
					return -1
				}
				return r
			}, v)
		}

		if opts.TrimSpace {
			v = strings.TrimSpace(v)
		}

		out[k] = v

		if k != origKey || v != origVal {
			reason := buildReason(origKey, k, origVal, v)
			changes = append(changes, Change{
				Key:    origKey,
				Before: origVal,
				After:  v,
				Reason: reason,
			})
		}
	}

	return Result{Env: out, Changes: changes}
}

func buildReason(origKey, newKey, origVal, newVal string) string {
	var parts []string
	if origKey != newKey {
		parts = append(parts, "key uppercased")
	}
	if origVal != newVal {
		parts = append(parts, "value sanitized")
	}
	return strings.Join(parts, ", ")
}
