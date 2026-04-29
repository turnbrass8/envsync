package diff2_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/diff2"
)

func TestDiff_NoChanges(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}
	r := diff2.Diff(left, right)
	if r.HasChanges() {
		t.Fatal("expected no changes")
	}
	if r.Summary() != "no changes" {
		t.Fatalf("unexpected summary: %s", r.Summary())
	}
}

func TestDiff_AddedKey(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"A": "1", "B": "2"}
	r := diff2.Diff(left, right)
	if !r.HasChanges() {
		t.Fatal("expected changes")
	}
	if r.Summary() != "1 added" {
		t.Fatalf("unexpected summary: %s", r.Summary())
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1"}
	r := diff2.Diff(left, right)
	if r.Summary() != "1 removed" {
		t.Fatalf("unexpected summary: %s", r.Summary())
	}
}

func TestDiff_ChangedKey(t *testing.T) {
	left := map[string]string{"A": "old"}
	right := map[string]string{"A": "new"}
	r := diff2.Diff(left, right)
	if r.Summary() != "1 changed" {
		t.Fatalf("unexpected summary: %s", r.Summary())
	}
}

func TestDiff_MixedChanges(t *testing.T) {
	left := map[string]string{"A": "1", "B": "old", "C": "3"}
	right := map[string]string{"A": "1", "B": "new", "D": "4"}
	r := diff2.Diff(left, right)
	if r.Summary() != "1 added, 1 removed, 1 changed" {
		t.Fatalf("unexpected summary: %s", r.Summary())
	}
}

func TestLine_String_Added(t *testing.T) {
	l := diff2.Line{Key: "FOO", NewValue: "bar", Kind: diff2.KindAdded}
	if l.String() != "+ FOO=bar" {
		t.Fatalf("unexpected string: %s", l.String())
	}
}

func TestLine_String_Removed(t *testing.T) {
	l := diff2.Line{Key: "FOO", OldValue: "bar", Kind: diff2.KindRemoved}
	if l.String() != "- FOO=bar" {
		t.Fatalf("unexpected string: %s", l.String())
	}
}

func TestLine_String_Changed(t *testing.T) {
	l := diff2.Line{Key: "FOO", OldValue: "old", NewValue: "new", Kind: diff2.KindChanged}
	want := `~ FOO: "old" -> "new"`
	if l.String() != want {
		t.Fatalf("got %q, want %q", l.String(), want)
	}
}

func TestDiff_SortedOutput(t *testing.T) {
	left := map[string]string{"Z": "1", "A": "1", "M": "1"}
	right := map[string]string{"Z": "1", "A": "1", "M": "1"}
	r := diff2.Diff(left, right)
	keys := make([]string, len(r.Lines))
	for i, l := range r.Lines {
		keys[i] = l.Key
	}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Fatalf("output not sorted: %v", keys)
	}
}
