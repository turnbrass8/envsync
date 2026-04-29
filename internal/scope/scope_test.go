package scope_test

import (
	"sort"
	"testing"

	"github.com/user/envsync/internal/scope"
)

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func TestExtract_ScopedKeysOnly(t *testing.T) {
	env := map[string]string{
		"PROD__DB_HOST": "db.prod.example.com",
		"PROD__DB_PORT": "5432",
		"STAGING__DB_HOST": "db.staging.example.com",
		"LOG_LEVEL": "info",
	}
	entries := scope.Extract(env, "prod", false)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Scope != "prod" {
			t.Errorf("expected scope=prod, got %q", e.Scope)
		}
	}
}

func TestExtract_IncludesGlobal(t *testing.T) {
	env := map[string]string{
		"PROD__DB_HOST": "db.prod.example.com",
		"LOG_LEVEL":     "info",
	}
	entries := scope.Extract(env, "prod", true)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestExtract_ExcludesOtherScopes(t *testing.T) {
	env := map[string]string{
		"STAGING__DB_HOST": "db.staging.example.com",
		"LOG_LEVEL":        "info",
	}
	entries := scope.Extract(env, "prod", false)
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestFlatten_ScopedTakesPrecedence(t *testing.T) {
	entries := []scope.Entry{
		{Key: "DB_HOST", Value: "global.db", Scope: ""},
		{Key: "DB_HOST", Value: "prod.db", Scope: "prod"},
	}
	out := scope.Flatten(entries)
	if out["DB_HOST"] != "prod.db" {
		t.Errorf("expected prod.db, got %q", out["DB_HOST"])
	}
}

func TestFlatten_GlobalUsedWhenNoScoped(t *testing.T) {
	entries := []scope.Entry{
		{Key: "LOG_LEVEL", Value: "info", Scope: ""},
	}
	out := scope.Flatten(entries)
	if out["LOG_LEVEL"] != "info" {
		t.Errorf("expected info, got %q", out["LOG_LEVEL"])
	}
}

func TestPrefix_FormatsCorrectly(t *testing.T) {
	got := scope.Prefix("staging", "DB_HOST")
	want := "STAGING__DB_HOST"
	if got != want {
		t.Errorf("Prefix: got %q, want %q", got, want)
	}
}

func TestFlatten_EmptyEntries_ReturnsEmptyMap(t *testing.T) {
	out := scope.Flatten(nil)
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", sortedKeys(out))
	}
}
