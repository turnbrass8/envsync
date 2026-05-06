// Package envdiff provides a line-level unified diff between two .env files.
package envdiff

import (
	"fmt"
	"sort"
	"strings"
)

// ChangeKind describes the type of change for a key.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
	Unchanged ChangeKind = "unchanged"
)

// Line represents a single diff entry for one key.
type Line struct {
	Key  string
	Kind ChangeKind
	Old  string
	New  string
}

// String returns a human-readable representation of the diff line.
func (l Line) String() string {
	switch l.Kind {
	case Added:
		return fmt.Sprintf("+ %s=%s", l.Key, l.New)
	case Removed:
		return fmt.Sprintf("- %s=%s", l.Key, l.Old)
	case Modified:
		return fmt.Sprintf("~ %s: %s -> %s", l.Key, l.Old, l.New)
	default:
		return fmt.Sprintf("  %s=%s", l.Key, l.New)
	}
}

// Result holds all diff lines produced by Diff.
type Result struct {
	Lines []Line
}

// HasChanges returns true if any line is not Unchanged.
func (r Result) HasChanges() bool {
	for _, l := range r.Lines {
		if l.Kind != Unchanged {
			return true
		}
	}
	return false
}

// Summary returns a short human-readable summary.
func (r Result) Summary() string {
	var added, removed, modified int
	for _, l := range r.Lines {
		switch l.Kind {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	if added+removed+modified == 0 {
		return "no changes"
	}
	parts := []string{}
	if added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", added))
	}
	if removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", removed))
	}
	if modified > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", modified))
	}
	return strings.Join(parts, ", ")
}

// Diff computes the difference between two env maps (key -> value).
// All keys from both maps are included; ordering is alphabetical.
func Diff(base, target map[string]string, includeUnchanged bool) Result {
	keys := unionKeys(base, target)
	sort.Strings(keys)

	var lines []Line
	for _, k := range keys {
		bv, inBase := base[k]
		tv, inTarget := target[k]
		switch {
		case inBase && !inTarget:
			lines = append(lines, Line{Key: k, Kind: Removed, Old: bv})
		case !inBase && inTarget:
			lines = append(lines, Line{Key: k, Kind: Added, New: tv})
		case bv != tv:
			lines = append(lines, Line{Key: k, Kind: Modified, Old: bv, New: tv})
		default:
			if includeUnchanged {
				lines = append(lines, Line{Key: k, Kind: Unchanged, Old: bv, New: tv})
			}
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
