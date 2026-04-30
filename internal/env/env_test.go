package env_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsync/internal/env"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestResolve_MergesSourcesInOrder(t *testing.T) {
	a := env.Source{Name: "a", Values: map[string]string{"FOO": "from_a", "BAR": "bar_a"}}
	b := env.Source{Name: "b", Values: map[string]string{"FOO": "from_b", "BAZ": "baz_b"}}
	l := env.NewLoader(a, b)
	result := l.Resolve()
	if result["FOO"] != "from_b" {
		t.Errorf("expected FOO=from_b, got %s", result["FOO"])
	}
	if result["BAR"] != "bar_a" {
		t.Errorf("expected BAR=bar_a, got %s", result["BAR"])
	}
	if result["BAZ"] != "baz_b" {
		t.Errorf("expected BAZ=baz_b, got %s", result["BAZ"])
	}
}

func TestKeys_SortedAndDeduped(t *testing.T) {
	a := env.Source{Name: "a", Values: map[string]string{"FOO": "1", "BAR": "2"}}
	b := env.Source{Name: "b", Values: map[string]string{"FOO": "3", "ZAP": "4"}}
	l := env.NewLoader(a, b)
	keys := l.Keys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "BAR" || keys[1] != "FOO" || keys[2] != "ZAP" {
		t.Errorf("unexpected key order: %v", keys)
	}
}

func TestOrigin_ReturnsLastSource(t *testing.T) {
	a := env.Source{Name: "file", Values: map[string]string{"KEY": "v1"}}
	b := env.Source{Name: "os", Values: map[string]string{"KEY": "v2"}}
	l := env.NewLoader(a, b)
	if got := l.Origin("KEY"); got != "os" {
		t.Errorf("expected origin=os, got %s", got)
	}
}

func TestOrigin_MissingKey_ReturnsEmpty(t *testing.T) {
	a := env.Source{Name: "file", Values: map[string]string{"FOO": "bar"}}
	l := env.NewLoader(a)
	if got := l.Origin("MISSING"); got != "" {
		t.Errorf("expected empty origin, got %s", got)
	}
}

func TestFileSource_ParsesEnvFile(t *testing.T) {
	p := writeTempEnv(t, "APP_ENV=production\nDEBUG=false\n")
	src, err := env.FileSource("dotenv", p)
	if err != nil {
		t.Fatal(err)
	}
	if src.Values["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %s", src.Values["APP_ENV"])
	}
}

func TestFileSource_MissingFile_ReturnsError(t *testing.T) {
	_, err := env.FileSource("dotenv", "/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
