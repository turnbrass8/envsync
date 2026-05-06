// Package flatten provides utilities for flattening nested env-like structures
// (e.g. JSON objects) into dot-notation KEY=value pairs suitable for .env files.
package flatten

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Options controls flattening behaviour.
type Options struct {
	// Separator is the string placed between nested keys (default "_").
	Separator string
	// Prefix is prepended to every key.
	Prefix string
	// UpperCase converts all keys to upper-case.
	UpperCase bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Separator: "_",
		UpperCase: true,
	}
}

// Flatten parses raw JSON and returns a sorted slice of "KEY=value" strings.
func Flatten(data []byte, opts Options) ([]string, error) {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	var root interface{}
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("flatten: invalid JSON: %w", err)
	}

	pairs := make(map[string]string)
	flattenValue(pairs, opts.Prefix, root, opts)

	keys := make([]string, 0, len(pairs))
	for k := range pairs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]string, 0, len(pairs))
	for _, k := range keys {
		out = append(out, k+"="+pairs[k])
	}
	return out, nil
}

// FlattenMap parses raw JSON and returns a map of KEY -> value pairs.
// It is useful when callers need to look up individual keys without
// iterating over the full KEY=value slice returned by Flatten.
func FlattenMap(data []byte, opts Options) (map[string]string, error) {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	var root interface{}
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("flatten: invalid JSON: %w", err)
	}

	pairs := make(map[string]string)
	flattenValue(pairs, opts.Prefix, root, opts)
	return pairs, nil
}

func flattenValue(out map[string]string, prefix string, v interface{}, opts Options) {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, child := range val {
			newKey := join(prefix, k, opts)
			flattenValue(out, newKey, child, opts)
		}
	case []interface{}:
		for i, child := range val {
			newKey := join(prefix, fmt.Sprintf("%d", i), opts)
			flattenValue(out, newKey, child, opts)
		}
	case nil:
		out[prefix] = ""
	default:
		out[prefix] = fmt.Sprintf("%v", val)
	}
}

func join(prefix, key string, opts Options) string {
	if opts.UpperCase {
		key = strings.ToUpper(key)
	}
	if prefix == "" {
		return key
	}
	return prefix + opts.Separator + key
}
