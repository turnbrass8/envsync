package compare

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompare_AllMatch(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}

	r := Compare(left, right)
	if r.HasDiff() {
		t.Fatal("expected no diff")
	}
	if len(r.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Entries))
	}
}

func TestCompare_ValuesDiffer(t *testing.T) {
	left := map[string]string{"HOST": "localhost"}
	right := map[string]string{"HOST": "production.example.com"}

	r := Compare(left, right)
	if !r.HasDiff() {
		t.Fatal("expected diff")
	}
	if r.Entries[0].Status != StatusDiffer {
		t.Errorf("expected StatusDiffer, got %v", r.Entries[0].Status)
	}
}

func TestCompare_OnlyLeft(t *testing.T) {
	left := map[string]string{"ONLY_LEFT": "val"}
	right := map[string]string{}

	r := Compare(left, right)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	if r.Entries[0].Status != StatusOnlyLeft {
		t.Errorf("expected StatusOnlyLeft, got %v", r.Entries[0].Status)
	}
}

func TestCompare_OnlyRight(t *testing.T) {
	left := map[string]string{}
	right := map[string]string{"ONLY_RIGHT": "val"}

	r := Compare(left, right)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	if r.Entries[0].Status != StatusOnlyRight {
		t.Errorf("expected StatusOnlyRight, got %v", r.Entries[0].Status)
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	left := map[string]string{"Z": "1", "A": "2", "M": "3"}
	right := map[string]string{"Z": "1", "A": "2", "M": "3"}

	r := Compare(left, right)
	keys := make([]string, len(r.Entries))
	for i, e := range r.Entries {
		keys[i] = e.Key
	}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("expected sorted keys, got %v", keys)
	}
}

func TestEntry_String(t *testing.T) {
	tests := []struct {
		status Status
		prefix string
	}{
		{StatusMatch, "  "},
		{StatusDiffer, "~ "},
		{StatusOnlyLeft, "< "},
		{StatusOnlyRight, "> "},
	}
	for _, tt := range tests {
		e := Entry{Key: "FOO", Status: tt.status}
		if !strings.HasPrefix(e.String(), tt.prefix) {
			t.Errorf("status %v: expected prefix %q, got %q", tt.status, tt.prefix, e.String())
		}
	}
}

func TestResult_Print_SkipsMatches(t *testing.T) {
	left := map[string]string{"A": "same", "B": "left-only"}
	right := map[string]string{"A": "same"}

	r := Compare(left, right)
	var buf bytes.Buffer
	r.Print(&buf)

	output := buf.String()
	if strings.Contains(output, "same") {
		t.Error("Print should skip matching entries")
	}
	if !strings.Contains(output, "B") {
		t.Error("Print should include differing entries")
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	left := map[string]string{}
	right := map[string]string{}

	r := Compare(left, right)
	if r.HasDiff() {
		t.Fatal("expected no diff for two empty maps")
	}
	if len(r.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(r.Entries))
	}
}
