package envdiff_test

import (
	"strings"
	"testing"

	"github.com/user/envsync/internal/envdiff"
)

func TestDiff_NoChanges(t *testing.T) {
	base := map[string]string{"FOO": "bar", "BAZ": "qux"}
	target := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res := envdiff.Diff(base, target, false)
	if res.HasChanges() {
		t.Fatal("expected no changes")
	}
	if res.Summary() != "no changes" {
		t.Fatalf("unexpected summary: %s", res.Summary())
	}
}

func TestDiff_AddedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "NEW": "val"}
	res := envdiff.Diff(base, target, false)
	if !res.HasChanges() {
		t.Fatal("expected changes")
	}
	if len(res.Lines) != 1 || res.Lines[0].Kind != envdiff.Added {
		t.Fatalf("expected one Added line, got %+v", res.Lines)
	}
	if !strings.Contains(res.Summary(), "1 added") {
		t.Fatalf("unexpected summary: %s", res.Summary())
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD": "val"}
	target := map[string]string{"FOO": "bar"}
	res := envdiff.Diff(base, target, false)
	if len(res.Lines) != 1 || res.Lines[0].Kind != envdiff.Removed {
		t.Fatalf("expected one Removed line, got %+v", res.Lines)
	}
	if !strings.Contains(res.Summary(), "1 removed") {
		t.Fatalf("unexpected summary: %s", res.Summary())
	}
}

func TestDiff_ModifiedKey(t *testing.T) {
	base := map[string]string{"FOO": "old"}
	target := map[string]string{"FOO": "new"}
	res := envdiff.Diff(base, target, false)
	if len(res.Lines) != 1 || res.Lines[0].Kind != envdiff.Modified {
		t.Fatalf("expected one Modified line, got %+v", res.Lines)
	}
	if res.Lines[0].Old != "old" || res.Lines[0].New != "new" {
		t.Fatalf("unexpected old/new: %+v", res.Lines[0])
	}
	if !strings.Contains(res.Summary(), "1 modified") {
		t.Fatalf("unexpected summary: %s", res.Summary())
	}
}

func TestDiff_IncludeUnchanged(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	target := map[string]string{"A": "1", "B": "changed"}
	res := envdiff.Diff(base, target, true)
	if len(res.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(res.Lines))
	}
}

func TestLine_String_Formats(t *testing.T) {
	cases := []struct {
		line envdiff.Line
		want string
	}{
		{envdiff.Line{Key: "K", Kind: envdiff.Added, New: "v"}, "+ K=v"},
		{envdiff.Line{Key: "K", Kind: envdiff.Removed, Old: "v"}, "- K=v"},
		{envdiff.Line{Key: "K", Kind: envdiff.Modified, Old: "a", New: "b"}, "~ K: a -> b"},
		{envdiff.Line{Key: "K", Kind: envdiff.Unchanged, New: "v"}, "  K=v"},
	}
	for _, c := range cases {
		got := c.line.String()
		if got != c.want {
			t.Errorf("String() = %q, want %q", got, c.want)
		}
	}
}

func TestDiff_SortedOutput(t *testing.T) {
	base := map[string]string{"Z": "1", "A": "1", "M": "1"}
	target := map[string]string{"Z": "1", "A": "2", "M": "1"}
	res := envdiff.Diff(base, target, true)
	keys := make([]string, len(res.Lines))
	for i, l := range res.Lines {
		keys[i] = l.Key
	}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Fatalf("expected sorted keys A,M,Z got %v", keys)
	}
}
