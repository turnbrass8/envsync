package strip_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/strip"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestStrip_ExactKey(t *testing.T) {
	src := writeTempEnv(t, "APP_NAME=myapp\nSECRET_KEY=abc123\nDEBUG=true\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	res, err := strip.Strip(src, dst, strip.Options{Keys: []string{"SECRET_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "SECRET_KEY" {
		t.Fatalf("expected [SECRET_KEY] removed, got %v", res.Removed)
	}

	env, _ := envfile.Parse(dst)
	if _, ok := env.Get("SECRET_KEY"); ok {
		t.Error("SECRET_KEY should have been removed")
	}
	if v, ok := env.Get("APP_NAME"); !ok || v != "myapp" {
		t.Errorf("APP_NAME should still be present, got %q %v", v, ok)
	}
}

func TestStrip_PatternMatch(t *testing.T) {
	src := writeTempEnv(t, "DB_HOST=localhost\nDB_PASS=secret\nAPP_ENV=prod\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	res, err := strip.Strip(src, dst, strip.Options{Patterns: []string{"^DB_"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d: %v", len(res.Removed), res.Removed)
	}

	env, _ := envfile.Parse(dst)
	if _, ok := env.Get("DB_HOST"); ok {
		t.Error("DB_HOST should have been removed")
	}
}

func TestStrip_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempEnv(t, "REMOVE_ME=yes\nKEEP=no\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	res, err := strip.Strip(src, dst, strip.Options{
		Keys:   []string{"REMOVE_ME"},
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 1 {
		t.Fatalf("expected 1 in dry-run result, got %d", len(res.Removed))
	}
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		t.Error("dst file should not exist in dry-run mode")
	}
}

func TestStrip_InvalidPattern_ReturnsError(t *testing.T) {
	src := writeTempEnv(t, "KEY=val\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	_, err := strip.Strip(src, dst, strip.Options{Patterns: []string{"[invalid"}})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestStrip_NoMatchLeavesFileUnchanged(t *testing.T) {
	src := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	res, err := strip.Strip(src, dst, strip.Options{Keys: []string{"NONEXISTENT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 0 {
		t.Fatalf("expected 0 removed, got %v", res.Removed)
	}
	env, _ := envfile.Parse(dst)
	if _, ok := env.Get("FOO"); !ok {
		t.Error("FOO should still be present")
	}
}
