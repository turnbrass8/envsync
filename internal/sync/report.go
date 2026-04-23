package sync

import (
	"fmt"
	"io"
	"strings"
)

// Print writes a human-readable sync report to w.
func (r *Result) Print(w io.Writer) {
	if len(r.Applied) == 0 && len(r.Skipped) == 0 && len(r.Errors) == 0 {
		fmt.Fprintln(w, "✔  env file is already in sync")
		return
	}

	if len(r.Applied) > 0 {
		fmt.Fprintf(w, "Applied (%d):\n", len(r.Applied))
		for _, a := range r.Applied {
			fmt.Fprintf(w, "  + %s\n", a)
		}
	}

	if len(r.Skipped) > 0 {
		fmt.Fprintf(w, "Skipped — no default (%d):\n", len(r.Skipped))
		for _, sk := range r.Skipped {
			fmt.Fprintf(w, "  ~ %s\n", sk)
		}
	}

	if len(r.Errors) > 0 {
		fmt.Fprintf(w, "Errors (%d):\n", len(r.Errors))
		for _, e := range r.Errors {
			fmt.Fprintf(w, "  ✗ %s\n", e)
		}
	}
}

// Summary returns a one-line summary string.
func (r *Result) Summary() string {
	parts := []string{}
	if n := len(r.Applied); n > 0 {
		parts = append(parts, fmt.Sprintf("%d applied", n))
	}
	if n := len(r.Skipped); n > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", n))
	}
	if n := len(r.Errors); n > 0 {
		parts = append(parts, fmt.Sprintf("%d error(s)", n))
	}
	if len(parts) == 0 {
		return "in sync"
	}
	return strings.Join(parts, ", ")
}
