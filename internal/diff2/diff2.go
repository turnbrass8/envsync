// Package diff2 provides a line-by-line unified diff between two .env files.
package diff2

import (
	"fmt"
	"sort"
	"strings"
)

// LineKind describes the nature of a diff line.
type LineKind int

const (
	KindEqual   LineKind = iota
	KindAdded             // present in right, absent in left
	KindRemoved           // present in left, absent in right
	KindChanged           // key present in both but value differs
)

// Line represents a single diff entry.
type Line struct {
	Key      string
	OldValue string
	NewValue string
	Kind     LineKind
}

// String returns a human-readable representation of the diff line.
func (l Line) String() string {
	switch l.Kind {
	case KindAdded:
		return fmt.Sprintf("+ %s=%s", l.Key, l.NewValue)
	case KindRemoved:
		return fmt.Sprintf("- %s=%s", l.Key, l.OldValue)
	case KindChanged:
		return fmt.Sprintf("~ %s: %q -> %q", l.Key, l.OldValue, l.NewValue)
	default:
		return fmt.Sprintf("  %s=%s", l.Key, l.NewValue)
	}
}

// Result holds the full diff between two env maps.
type Result struct {
	Lines []Line
}

// HasChanges reports whether any non-equal lines exist.
func (r Result) HasChanges() bool {
	for _, l := range r.Lines {
		if l.Kind != KindEqual {
			return true
		}
	}
	return false
}

// Summary returns a one-line summary of the diff.
func (r Result) Summary() string {
	var added, removed, changed int
	for _, l := range r.Lines {
		switch l.Kind {
		case KindAdded:
			added++
		case KindRemoved:
			removed++
		case KindChanged:
			changed++
		}
	}
	if added == 0 && removed == 0 && changed == 0 {
		return "no changes"
	}
	parts := []string{}
	if added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", added))
	}
	if removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", removed))
	}
	if changed > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", changed))
	}
	return strings.Join(parts, ", ")
}

// Diff computes the unified diff between left and right env maps.
func Diff(left, right map[string]string) Result {
	keys := unionKeys(left, right)
	sort.Strings(keys)

	var lines []Line
	for _, k := range keys {
		lv, inLeft := left[k]
		rv, inRight := right[k]
		switch {
		case inLeft && inRight && lv == rv:
			lines = append(lines, Line{Key: k, OldValue: lv, NewValue: rv, Kind: KindEqual})
		case inLeft && inRight:
			lines = append(lines, Line{Key: k, OldValue: lv, NewValue: rv, Kind: KindChanged})
		case inLeft:
			lines = append(lines, Line{Key: k, OldValue: lv, Kind: KindRemoved})
		default:
			lines = append(lines, Line{Key: k, NewValue: rv, Kind: KindAdded})
		}
	}
	return Result{Lines: lines}
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}
