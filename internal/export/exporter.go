package export

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Format represents an output format for exported env vars.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatShell  Format = "shell"
	FormatJSON   Format = "json"
)

// Options controls how the export is rendered.
type Options struct {
	Format  Format
	Sorted  bool
	Writer  io.Writer
}

// Export writes the given key=value map to the configured writer in the
// requested format. If opts.Writer is nil, os.Stdout is used.
func Export(env map[string]string, opts Options) error {
	w := opts.Writer
	if w == nil {
		w = os.Stdout
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if opts.Sorted {
		sort.Strings(keys)
	}

	switch opts.Format {
	case FormatShell:
		return exportShell(w, keys, env)
	case FormatJSON:
		return exportJSON(w, keys, env)
	default:
		return exportDotenv(w, keys, env)
	}
}

func exportDotenv(w io.Writer, keys []string, env map[string]string) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, quoteIfNeeded(env[k])); err != nil {
			return err
		}
	}
	return nil
}

func exportShell(w io.Writer, keys []string, env map[string]string) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "export %s=%s\n", k, quoteIfNeeded(env[k])); err != nil {
			return err
		}
	}
	return nil
}

func exportJSON(w io.Writer, keys []string, env map[string]string) error {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		sb.WriteString(fmt.Sprintf("  %q: %q", k, env[k]))
		if i < len(keys)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("}\n")
	_, err := fmt.Fprint(w, sb.String())
	return err
}

// quoteIfNeeded wraps the value in double quotes if it contains characters
// that may require quoting in shell or dotenv contexts (spaces, tabs, or #).
// Values that already contain a double quote are also quoted to avoid ambiguity.
func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t#\"") {
		return fmt.Sprintf("%q", v)
	}
	return v
}
