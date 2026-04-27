// Package redact provides utilities for masking sensitive values in .env files
// before displaying or logging them.
package redact

import "strings"

// DefaultPatterns are key name substrings that trigger redaction.
var DefaultPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
}

const mask = "***REDACTED***"

// Options controls redaction behaviour.
type Options struct {
	// Patterns is the list of key substrings (case-insensitive) that trigger
	// redaction. If nil, DefaultPatterns is used.
	Patterns []string
	// Mask is the replacement string. Defaults to "***REDACTED***".
	Mask string
}

func (o Options) patterns() []string {
	if len(o.Patterns) > 0 {
		return o.Patterns
	}
	return DefaultPatterns
}

func (o Options) maskStr() string {
	if o.Mask != "" {
		return o.Mask
	}
	return mask
}

// IsSensitive reports whether key matches any of the given patterns
// (case-insensitive substring match).
func IsSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// Apply returns a copy of env where sensitive values are replaced with the
// configured mask string.
func Apply(env map[string]string, opts Options) map[string]string {
	patterns := opts.patterns()
	m := opts.maskStr()
	out := make(map[string]string, len(env))
	for k, v := range env {
		if IsSensitive(k, patterns) {
			out[k] = m
		} else {
			out[k] = v
		}
	}
	return out
}

// Value redacts a single value if the key is considered sensitive.
func Value(key, value string, opts Options) string {
	if IsSensitive(key, opts.patterns()) {
		return opts.maskStr()
	}
	return value
}
