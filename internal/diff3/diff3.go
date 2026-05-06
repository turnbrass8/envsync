// Package diff3 provides three-way diff between a base, left, and right .env file.
// It identifies keys that were added, removed, modified, or have conflicting changes.
package diff3

import "fmt"

// ChangeKind describes how a key changed relative to the base.
type ChangeKind string

const (
	Unchanged  ChangeKind = "unchanged"
	AddedLeft  ChangeKind = "added_left"
	AddedRight ChangeKind = "added_right"
	AddedBoth  ChangeKind = "added_both"
	RemovedLeft  ChangeKind = "removed_left"
	RemovedRight ChangeKind = "removed_right"
	ModifiedLeft  ChangeKind = "modified_left"
	ModifiedRight ChangeKind = "modified_right"
	Conflict ChangeKind = "conflict"
)

// Entry represents a single key in the three-way diff result.
type Entry struct {
	Key        string
	Base       string
	Left       string
	Right      string
	Kind       ChangeKind
}

func (e Entry) String() string {
	switch e.Kind {
	case Conflict:
		return fmt.Sprintf("CONFLICT %s: left=%q right=%q (base=%q)", e.Key, e.Left, e.Right, e.Base)
	case AddedLeft:
		return fmt.Sprintf("+ (left)  %s=%q", e.Key, e.Left)
	case AddedRight:
		return fmt.Sprintf("+ (right) %s=%q", e.Key, e.Right)
	case AddedBoth:
		if e.Left == e.Right {
			return fmt.Sprintf("+ (both)  %s=%q", e.Key, e.Left)
		}
		return fmt.Sprintf("CONFLICT %s: left=%q right=%q (added in both)", e.Key, e.Left, e.Right)
	case RemovedLeft:
		return fmt.Sprintf("- (left)  %s", e.Key)
	case RemovedRight:
		return fmt.Sprintf("- (right) %s", e.Key)
	case ModifiedLeft:
		return fmt.Sprintf("~ (left)  %s=%q (was %q)", e.Key, e.Left, e.Base)
	case ModifiedRight:
		return fmt.Sprintf("~ (right) %s=%q (was %q)", e.Key, e.Right, e.Base)
	default:
		return fmt.Sprintf("  %s=%q", e.Key, e.Base)
	}
}

// Result holds the full three-way diff output.
type Result struct {
	Entries   []Entry
	Conflicts int
}

// Diff performs a three-way diff of base, left, and right env maps.
func Diff(base, left, right map[string]string) Result {
	keys := unionKeys(base, left, right)
	var entries []Entry
	conflicts := 0

	for _, k := range keys {
		bv, inBase := base[k]
		lv, inLeft := left[k]
		rv, inRight := right[k]

		var e Entry
		e.Key = k
		e.Base = bv
		e.Left = lv
		e.Right = rv

		switch {
		case !inBase && inLeft && !inRight:
			e.Kind = AddedLeft
		case !inBase && !inLeft && inRight:
			e.Kind = AddedRight
		case !inBase && inLeft && inRight:
			e.Kind = AddedBoth
			if lv != rv {
				conflicts++
			}
		case inBase && !inLeft && inRight:
			e.Kind = RemovedLeft
		case inBase && inLeft && !inRight:
			e.Kind = RemovedRight
		case inBase && !inLeft && !inRight:
			// removed in both — treat as unchanged/removed
			e.Kind = RemovedLeft
		case lv == bv && rv == bv:
			e.Kind = Unchanged
		case lv != bv && rv == bv:
			e.Kind = ModifiedLeft
		case lv == bv && rv != bv:
			e.Kind = ModifiedRight
		case lv == rv:
			e.Kind = ModifiedLeft // same change on both sides
		default:
			e.Kind = Conflict
			conflicts++
		}
		entries = append(entries, e)
	}
	return Result{Entries: entries, Conflicts: conflicts}
}

func unionKeys(maps ...map[string]string) []string {
	seen := map[string]struct{}{}
	for _, m := range maps {
		for k := range m {
			seen[k] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sortStrings(out)
	return out
}

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
