package diff_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/diff"
)

func TestCompare_AllPresent(t *testing.T) {
	keys := []string{"DB_HOST", "DB_PORT"}
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}

	result := diff.Compare(keys, env)

	if result.HasMissing() {
		t.Error("expected no missing keys")
	}
	if len(result.Extra()) != 0 {
		t.Errorf("expected no extra keys, got %d", len(result.Extra()))
	}
	if len(result.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result.Entries))
	}
}

func TestCompare_MissingKeys(t *testing.T) {
	keys := []string{"DB_HOST", "DB_PORT", "SECRET_KEY"}
	env := map[string]string{"DB_HOST": "localhost"}

	result := diff.Compare(keys, env)

	if !result.HasMissing() {
		t.Error("expected missing keys")
	}
	missing := result.Missing()
	if len(missing) != 2 {
		t.Errorf("expected 2 missing keys, got %d", len(missing))
	}
	if missing[0].Key != "DB_PORT" {
		t.Errorf("expected DB_PORT missing, got %s", missing[0].Key)
	}
	if missing[1].Key != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY missing, got %s", missing[1].Key)
	}
}

func TestCompare_ExtraKeys(t *testing.T) {
	keys := []string{"DB_HOST"}
	env := map[string]string{"DB_HOST": "localhost", "EXTRA_VAR": "foo", "ANOTHER": "bar"}

	result := diff.Compare(keys, env)

	extra := result.Extra()
	if len(extra) != 2 {
		t.Errorf("expected 2 extra keys, got %d", len(extra))
	}
	// Extra keys are sorted
	if extra[0].Key != "ANOTHER" {
		t.Errorf("expected ANOTHER first extra key, got %s", extra[0].Key)
	}
	if extra[1].Key != "EXTRA_VAR" {
		t.Errorf("expected EXTRA_VAR second extra key, got %s", extra[1].Key)
	}
}

func TestCompare_EmptyEnv(t *testing.T) {
	keys := []string{"A", "B", "C"}
	env := map[string]string{}

	result := diff.Compare(keys, env)

	if len(result.Missing()) != 3 {
		t.Errorf("expected 3 missing keys, got %d", len(result.Missing()))
	}
}

func TestEntry_String(t *testing.T) {
	cases := []struct {
		entry    diff.Entry
		expected string
	}{
		{diff.Entry{Key: "FOO", Status: diff.StatusMissing}, "- FOO (missing)"},
		{diff.Entry{Key: "BAR", Status: diff.StatusExtra, Value: "baz"}, "+ BAR (extra)"},
		{diff.Entry{Key: "QUX", Status: diff.StatusPresent, Value: "1"}, "  QUX"},
	}
	for _, tc := range cases {
		got := tc.entry.String()
		if got != tc.expected {
			t.Errorf("String() = %q, want %q", got, tc.expected)
		}
	}
}
