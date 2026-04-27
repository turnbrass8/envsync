// Package mask provides utilities for masking sensitive environment variable
// values in output, logs, and diffs — replacing them with a fixed placeholder.
package mask

import (
	"strings"
)

const DefaultMask = "****"

// Rule defines a key pattern and the mask to apply.
type Rule struct {
	Key  string // exact key name (case-insensitive)
	Mask string // replacement string; defaults to DefaultMask
}

// Masker holds a set of rules and applies them to env maps.
type Masker struct {
	rules []Rule
}

// New creates a Masker with the provided rules. If a rule's Mask is empty,
// DefaultMask is used.
func New(rules []Rule) *Masker {
	normalized := make([]Rule, len(rules))
	for i, r := range rules {
		if r.Mask == "" {
			r.Mask = DefaultMask
		}
		r.Key = strings.ToUpper(r.Key)
		normalized[i] = r
	}
	return &Masker{rules: normalized}
}

// Apply returns a copy of env with sensitive values replaced according to rules.
func (m *Masker) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	for _, rule := range m.rules {
		for k := range out {
			if strings.ToUpper(k) == rule.Key {
				out[k] = rule.Mask
			}
		}
	}
	return out
}

// Value masks a single key/value pair. Returns the original value if no rule
// matches, or the mask string if a rule does.
func (m *Masker) Value(key, value string) string {
	upper := strings.ToUpper(key)
	for _, rule := range m.rules {
		if rule.Key == upper {
			return rule.Mask
		}
	}
	return value
}

// DefaultRules returns a sensible set of rules covering common secret key names.
func DefaultRules() []Rule {
	keys := []string{
		"PASSWORD", "PASSWD", "SECRET", "TOKEN", "API_KEY", "APIKEY",
		"PRIVATE_KEY", "ACCESS_KEY", "AUTH_TOKEN", "CLIENT_SECRET",
	}
	rules := make([]Rule, len(keys))
	for i, k := range keys {
		rules[i] = Rule{Key: k, Mask: DefaultMask}
	}
	return rules
}
