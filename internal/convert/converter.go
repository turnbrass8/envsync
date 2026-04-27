// Package convert provides utilities for converting .env files between formats.
package convert

import (
	"fmt"
	"sort"
	"strings"
)

// Format represents a supported conversion target format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatYAML   Format = "yaml"
	FormatTOML   Format = "toml"
)

// Options controls conversion behaviour.
type Options struct {
	Format  Format
	Sorted  bool
	Prefix  string // optional key prefix to strip before output
}

// Convert transforms a map of env key/value pairs into the target format string.
func Convert(env map[string]string, opts Options) (string, error) {
	keys := sortedKeys(env, opts.Sorted)

	switch opts.Format {
	case FormatDotenv:
		return toDotenv(env, keys, opts.Prefix), nil
	case FormatJSON:
		return toJSON(env, keys, opts.Prefix), nil
	case FormatYAML:
		return toYAML(env, keys, opts.Prefix), nil
	case FormatTOML:
		return toTOML(env, keys, opts.Prefix), nil
	default:
		return "", fmt.Errorf("unsupported format: %q", opts.Format)
	}
}

func stripPrefix(key, prefix string) string {
	if prefix != "" {
		return strings.TrimPrefix(key, prefix)
	}
	return key
}

func sortedKeys(env map[string]string, sorted bool) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if sorted {
		sort.Strings(keys)
	}
	return keys
}

func toDotenv(env map[string]string, keys []string, prefix string) string {
	var sb strings.Builder
	for _, k := range keys {
		out := stripPrefix(k, prefix)
		fmt.Fprintf(&sb, "%s=%s\n", out, quoteIfNeeded(env[k]))
	}
	return sb.String()
}

func toJSON(env map[string]string, keys []string, prefix string) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		out := stripPrefix(k, prefix)
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", out, env[k], comma)
	}
	sb.WriteString("}\n")
	return sb.String()
}

func toYAML(env map[string]string, keys []string, prefix string) string {
	var sb strings.Builder
	for _, k := range keys {
		out := stripPrefix(k, prefix)
		fmt.Fprintf(&sb, "%s: %q\n", out, env[k])
	}
	return sb.String()
}

func toTOML(env map[string]string, keys []string, prefix string) string {
	var sb strings.Builder
	for _, k := range keys {
		out := stripPrefix(k, prefix)
		fmt.Fprintf(&sb, "%s = %q\n", out, env[k])
	}
	return sb.String()
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t#") {
		return fmt.Sprintf("%q", v)
	}
	return v
}
