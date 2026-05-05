package resolve_test

import (
	"testing"

	"github.com/user/envsync/internal/resolve"
)

func src(name string, kv map[string]string) resolve.Source {
	return resolve.Source{Name: name, Values: kv}
}

func TestResolve_FirstSourceWins(t *testing.T) {
	sources := []resolve.Source{
		src("prod", map[string]string{"DB_HOST": "prod-db"}),
		src("dev", map[string]string{"DB_HOST": "dev-db", "API_KEY": "abc"}),
	}
	results, err := resolve.Resolve([]string{"DB_HOST", "API_KEY"}, sources, resolve.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Value != "prod-db" || results[0].Source != "prod" {
		t.Errorf("expected prod-db from prod, got %q from %q", results[0].Value, results[0].Source)
	}
	if results[1].Value != "abc" || results[1].Source != "dev" {
		t.Errorf("expected abc from dev, got %q from %q", results[1].Value, results[1].Source)
	}
}

func TestResolve_StrictMissingReturnsError(t *testing.T) {
	sources := []resolve.Source{
		src("env", map[string]string{"PRESENT": "yes"}),
	}
	_, err := resolve.Resolve([]string{"PRESENT", "MISSING"}, sources, resolve.Options{Strict: true})
	if err == nil {
		t.Fatal("expected error for missing key in strict mode")
	}
}

func TestResolve_NonStrictMissingUsesFallback(t *testing.T) {
	sources := []resolve.Source{}
	results, err := resolve.Resolve([]string{"GHOST"}, sources, resolve.Options{Fallback: "default"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Value != "default" {
		t.Errorf("expected fallback 'default', got %q", results[0].Value)
	}
	if results[0].Found {
		t.Error("expected Found=false for missing key")
	}
}

func TestResolveAll_MergesAndSorts(t *testing.T) {
	sources := []resolve.Source{
		src("a", map[string]string{"Z_KEY": "z", "A_KEY": "a-from-a"}),
		src("b", map[string]string{"A_KEY": "a-from-b", "M_KEY": "m"}),
	}
	results := resolve.ResolveAll(sources, resolve.Options{})
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	// Sorted: A_KEY, M_KEY, Z_KEY
	if results[0].Key != "A_KEY" || results[0].Source != "a" {
		t.Errorf("expected A_KEY from a, got %q from %q", results[0].Key, results[0].Source)
	}
	if results[1].Key != "M_KEY" {
		t.Errorf("expected M_KEY second, got %q", results[1].Key)
	}
}

func TestResolve_EmptySources(t *testing.T) {
	results, err := resolve.Resolve([]string{"X"}, nil, resolve.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Found {
		t.Error("expected Found=false with no sources")
	}
}
