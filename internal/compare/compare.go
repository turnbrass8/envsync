// Package compare provides utilities for comparing two .env files
// across environments and producing a structured diff report.
package compare

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Status represents the comparison state of a key between two env files.
type Status int

const (
	StatusMatch   Status = iota // key exists in both with equal values
	StatusDiffer                // key exists in both but values differ
	StatusOnlyLeft              // key exists only in the left (source) file
	StatusOnlyRight             // key exists only in the right (target) file
)

// Entry holds the comparison result for a single key.
type Entry struct {
	Key    string
	Left   string
	Right  string
	Status Status
}

// String returns a human-readable representation of the entry.
func (e Entry) String() string {
	switch e.Status {
	case StatusMatch:
		return fmt.Sprintf("  %s", e.Key)
	case StatusDiffer:
		return fmt.Sprintf("~ %s", e.Key)
	case StatusOnlyLeft:
		return fmt.Sprintf("< %s", e.Key)
	case StatusOnlyRight:
		return fmt.Sprintf("> %s", e.Key)
	}
	return e.Key
}

// Result holds all comparison entries between two env maps.
type Result struct {
	Entries []Entry
}

// HasDiff returns true if any entry is not a match.
func (r *Result) HasDiff() bool {
	for _, e := range r.Entries {
		if e.Status != StatusMatch {
			return true
		}
	}
	return false
}

// Print writes a formatted table of differences to w.
func (r *Result) Print(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "STATUS\tKEY\tLEFT\tRIGHT")
	for _, e := range r.Entries {
		if e.Status == StatusMatch {
			continue
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", e.String()[:1], e.Key, e.Left, e.Right)
	}
	tw.Flush()
}

// Compare compares two env maps (left=source, right=target) and returns a Result.
func Compare(left, right map[string]string) *Result {
	seen := make(map[string]bool)
	var entries []Entry

	for k, lv := range left {
		seen[k] = true
		rv, ok := right[k]
		switch {
		case !ok:
			entries = append(entries, Entry{Key: k, Left: lv, Right: "", Status: StatusOnlyLeft})
		case lv == rv:
			entries = append(entries, Entry{Key: k, Left: lv, Right: rv, Status: StatusMatch})
		default:
			entries = append(entries, Entry{Key: k, Left: lv, Right: rv, Status: StatusDiffer})
		}
	}

	for k, rv := range right {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Left: "", Right: rv, Status: StatusOnlyRight})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return &Result{Entries: entries}
}
