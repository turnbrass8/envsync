package filter_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/filter"
)

var base = map[string]string{
	"APP_HOST":    "localhost",
	"APP_PORT":    "8080",
	"DB_HOST":     "db.local",
	"DB_PASSWORD": "secret",
	"LOG_LEVEL":   "info",
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	out, err := filter.Filter(base, filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(base) {
		t.Fatalf("expected %d entries, got %d", len(base), len(out))
	}
}

func TestFilter_ByPrefix(t *testing.T) {
	out, err := filter.Filter(base, filter.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
}

func TestFilter_ByPattern(t *testing.T) {
	out, err := filter.Filter(base, filter.Options{Pattern: `^DB_`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestFilter_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := filter.Filter(base, filter.Options{Pattern: `[invalid`})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestFilter_ByExplicitKeys(t *testing.T) {
	out, err := filter.Filter(base, filter.Options{Keys: []string{"LOG_LEVEL", "APP_PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestFilter_Exclude(t *testing.T) {
	out, err := filter.Filter(base, filter.Options{Prefix: "DB_", Exclude: []string{"DB_PASSWORD"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been excluded")
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("DB_HOST should be present")
	}
}

func TestStripPrefix(t *testing.T) {
	src := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	out := filter.StripPrefix(src, "APP_")
	if v, ok := out["HOST"]; !ok || v != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", v)
	}
	if v, ok := out["PORT"]; !ok || v != "8080" {
		t.Errorf("expected PORT=8080, got %q", v)
	}
}
