// Package docs generates documentation for .env files by extracting
// keys, types, defaults, and descriptions from a manifest and env file.
package docs

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envsync/internal/manifest"
)

// Format controls the output format of generated documentation.
type Format string

const (
	FormatMarkdown Format = "markdown"
	FormatText     Format = "text"
	FormatCSV      Format = "csv"
)

// Options configures documentation generation.
type Options struct {
	Format      Format
	Title       string
	ShowDefaults bool
	ShowRequired bool
}

// DefaultOptions returns sensible defaults for documentation generation.
func DefaultOptions() Options {
	return Options{
		Format:       FormatMarkdown,
		Title:        "Environment Variables",
		ShowDefaults: true,
		ShowRequired: true,
	}
}

// Generate writes documentation derived from the manifest entries to w.
func Generate(entries []manifest.Entry, env map[string]string, opts Options, w io.Writer) error {
	sorted := make([]manifest.Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	switch opts.Format {
	case FormatMarkdown:
		return writeMarkdown(sorted, env, opts, w)
	case FormatText:
		return writeText(sorted, env, opts, w)
	case FormatCSV:
		return writeCSV(sorted, env, opts, w)
	default:
		return fmt.Errorf("unsupported format: %q", opts.Format)
	}
}

func writeMarkdown(entries []manifest.Entry, env map[string]string, opts Options, w io.Writer) error {
	fmt.Fprintf(w, "# %s\n\n", opts.Title)
	fmt.Fprintf(w, "| Key | Required | Default | Current |\n")
	fmt.Fprintf(w, "|-----|----------|---------|---------|\n")
	for _, e := range entries {
		required := "-"
		if opts.ShowRequired && e.Required {
			required = "✓"
		}
		def := "-"
		if opts.ShowDefaults && e.Default != "" {
			def = "`" + e.Default + "`"
		}
		current := "-"
		if v, ok := env[e.Key]; ok {
			current = "`" + v + "`"
		}
		fmt.Fprintf(w, "| `%s` | %s | %s | %s |\n", e.Key, required, def, current)
	}
	return nil
}

func writeText(entries []manifest.Entry, env map[string]string, opts Options, w io.Writer) error {
	fmt.Fprintf(w, "%s\n%s\n\n", opts.Title, strings.Repeat("=", len(opts.Title)))
	for _, e := range entries {
		fmt.Fprintf(w, "%s\n", e.Key)
		if opts.ShowRequired {
			if e.Required {
				fmt.Fprintf(w, "  required: yes\n")
			} else {
				fmt.Fprintf(w, "  required: no\n")
			}
		}
		if opts.ShowDefaults && e.Default != "" {
			fmt.Fprintf(w, "  default:  %s\n", e.Default)
		}
		if v, ok := env[e.Key]; ok {
			fmt.Fprintf(w, "  current:  %s\n", v)
		}
		fmt.Fprintln(w)
	}
	return nil
}

func writeCSV(entries []manifest.Entry, env map[string]string, opts Options, w io.Writer) error {
	fmt.Fprintln(w, "key,required,default,current")
	for _, e := range entries {
		required := "false"
		if e.Required {
			required = "true"
		}
		def := ""
		if opts.ShowDefaults {
			def = csvEscape(e.Default)
		}
		current := ""
		if v, ok := env[e.Key]; ok {
			current = csvEscape(v)
		}
		fmt.Fprintf(w, "%s,%s,%s,%s\n", e.Key, required, def, current)
	}
	return nil
}

// csvEscape wraps a value in quotes if it contains a comma, quote, or newline.
func csvEscape(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}
