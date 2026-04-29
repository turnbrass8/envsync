package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFilterEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempFilterEnv: %v", err)
	}
	return p
}

func TestRunFilter_PrefixFiltersKeys(t *testing.T) {
	p := writeTempFilterEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\nDB_HOST=db\n")
	if err := runFilter([]string{"-env", p, "-prefix", "APP_"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunFilter_PatternFiltersKeys(t *testing.T) {
	p := writeTempFilterEnv(t, "APP_HOST=localhost\nDB_HOST=db\nLOG_LEVEL=info\n")
	if err := runFilter([]string{"-env", p, "-pattern", `^APP_`}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunFilter_InvalidPattern_ReturnsError(t *testing.T) {
	p := writeTempFilterEnv(t, "KEY=val\n")
	if err := runFilter([]string{"-env", p, "-pattern", `[bad`}); err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestRunFilter_StripPrefix(t *testing.T) {
	p := writeTempFilterEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\n")
	if err := runFilter([]string{"-env", p, "-prefix", "APP_", "-strip-prefix"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunFilter_JSONFormat(t *testing.T) {
	p := writeTempFilterEnv(t, "KEY=value\n")
	if err := runFilter([]string{"-env", p, "-format", "json"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunFilter_MissingEnvFile_ReturnsError(t *testing.T) {
	if err := runFilter([]string{"-env", "/no/such/file.env"}); err == nil {
		t.Fatal("expected error for missing file")
	}
}
