package diff3_test

import (
	"strings"
	"testing"

	"github.com/user/envsync/internal/diff3"
)

func TestDiff_NoChanges(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}

	r := diff3.Diff(base, left, right)
	if r.Conflicts != 0 {
		t.Fatalf("expected 0 conflicts, got %d", r.Conflicts)
	}
	for _, e := range r.Entries {
		if e.Kind != diff3.Unchanged {
			t.Errorf("expected Unchanged for %s, got %s", e.Key, e.Kind)
		}
	}
}

func TestDiff_AddedLeft(t *testing.T) {
	base := map[string]string{}
	left := map[string]string{"NEW": "val"}
	right := map[string]string{}

	r := diff3.Diff(base, left, right)
	if len(r.Entries) != 1 || r.Entries[0].Kind != diff3.AddedLeft {
		t.Fatalf("expected AddedLeft, got %+v", r.Entries)
	}
}

func TestDiff_AddedRight(t *testing.T) {
	base := map[string]string{}
	left := map[string]string{}
	right := map[string]string{"NEW": "val"}

	r := diff3.Diff(base, left, right)
	if len(r.Entries) != 1 || r.Entries[0].Kind != diff3.AddedRight {
		t.Fatalf("expected AddedRight, got %+v", r.Entries)
	}
}

func TestDiff_ModifiedLeft(t *testing.T) {
	base := map[string]string{"X": "old"}
	left := map[string]string{"X": "new"}
	right := map[string]string{"X": "old"}

	r := diff3.Diff(base, left, right)
	if r.Entries[0].Kind != diff3.ModifiedLeft {
		t.Errorf("expected ModifiedLeft, got %s", r.Entries[0].Kind)
	}
}

func TestDiff_ModifiedRight(t *testing.T) {
	base := map[string]string{"X": "old"}
	left := map[string]string{"X": "old"}
	right := map[string]string{"X": "new"}

	r := diff3.Diff(base, left, right)
	if r.Entries[0].Kind != diff3.ModifiedRight {
		t.Errorf("expected ModifiedRight, got %s", r.Entries[0].Kind)
	}
}

func TestDiff_Conflict(t *testing.T) {
	base := map[string]string{"X": "base"}
	left := map[string]string{"X": "left-val"}
	right := map[string]string{"X": "right-val"}

	r := diff3.Diff(base, left, right)
	if r.Conflicts != 1 {
		t.Fatalf("expected 1 conflict, got %d", r.Conflicts)
	}
	if r.Entries[0].Kind != diff3.Conflict {
		t.Errorf("expected Conflict kind, got %s", r.Entries[0].Kind)
	}
}

func TestDiff_AddedBothSameValue_NoConflict(t *testing.T) {
	base := map[string]string{}
	left := map[string]string{"K": "same"}
	right := map[string]string{"K": "same"}

	r := diff3.Diff(base, left, right)
	if r.Conflicts != 0 {
		t.Errorf("expected no conflict for identical additions, got %d", r.Conflicts)
	}
	if r.Entries[0].Kind != diff3.AddedBoth {
		t.Errorf("expected AddedBoth, got %s", r.Entries[0].Kind)
	}
}

func TestEntry_String_ContainsKey(t *testing.T) {
	e := diff3.Entry{Key: "MY_KEY", Left: "lv", Right: "rv", Base: "bv", Kind: diff3.Conflict}
	s := e.String()
	if !strings.Contains(s, "MY_KEY") {
		t.Errorf("expected key in string output, got: %s", s)
	}
	if !strings.Contains(s, "CONFLICT") {
		t.Errorf("expected CONFLICT in string output, got: %s", s)
	}
}
