// Package format provides utilities for normalising .env file formatting:
// consistent quoting style, spacing around '=', and trailing newlines.
package format

import (
	"fmt"
	"strings"

	"github.com/user/envsync/internal/envfile"
)

// Options controls how the formatter rewrites a .env file.
type Options struct {
	// QuoteStyle is "none", "double" (default), or "single".
	QuoteStyle string
	// SpaceAroundEquals adds spaces around the '=' sign when true.
	SpaceAroundEquals bool
	// DryRun returns the formatted content without writing to disk.
	DryRun bool
}

// Result holds the outcome of a Format call.
type Result struct {
	Path     string
	Changed  bool
	Formatted string
}

// Format reads the .env file at path, applies formatting rules defined by
// opts, and writes the result back unless DryRun is set.
func Format(path string, opts Options) (Result, error) {
	if opts.QuoteStyle == "" {
		opts.QuoteStyle = "double"
	}

	env, err := envfile.Parse(path)
	if err != nil {
		return Result{}, fmt.Errorf("format: parse %q: %w", path, err)
	}

	var sb strings.Builder
	for _, key := range env.Keys() {
		val, _ := env.Get(key)
		formatted := formatLine(key, val, opts)
		sb.WriteString(formatted)
		sb.WriteByte('\n')
	}

	output := sb.String()

	original, _ := readRaw(path)
	changed := output != original

	result := Result{Path: path, Changed: changed, Formatted: output}

	if !opts.DryRun && changed {
		if err := writeRaw(path, output); err != nil {
			return Result{}, fmt.Errorf("format: write %q: %w", path, err)
		}
	}

	return result, nil
}

func formatLine(key, val string, opts Options) string {
	quoted := applyQuoteStyle(val, opts.QuoteStyle)
	if opts.SpaceAroundEquals {
		return key + " = " + quoted
	}
	return key + "=" + quoted
}

func applyQuoteStyle(val, style string) string {
	switch style {
	case "single":
		return "'" + strings.ReplaceAll(val, "'", `\'`) + "'"
	case "none":
		return val
	default: // "double"
		if strings.ContainsAny(val, " \t#") {
			return `"` + strings.ReplaceAll(val, `"`, `\"`) + `"`
		}
		return val
	}
}
