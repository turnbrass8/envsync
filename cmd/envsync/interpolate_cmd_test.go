package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestRunInterpolate_ResolvesReferences(t *testing.T) {
	p := writeTempEnv(t, "BASE=http://localhost\nURL=${BASE}/api\n")
	out := filepath.Join(t.TempDir(), "out.env")
	err := runInterpolate([]string{"-env", p, "-out", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "URL=http://localhost/api") {
		t.Errorf("expected resolved URL, got:\n%s", data)
	}
}

func TestRunInterpolate_StrictFails(t *testing.T) {
	p := writeTempEnv(t, "URL=${MISSING}\n")
	err := runInterpolate([]string{"-env", p, "-strict=true"})
	if err == nil {
		t.Fatal("expected error for unresolved variable in strict mode")
	}
}

func TestRunInterpolate_NonStrictSucceeds(t *testing.T) {
	p := writeTempEnv(t, "URL=${MISSING}\n")
	err := runInterpolate([]string{"-env", p, "-strict=false"})
	if err != nil {
		t.Fatalf("unexpected error in non-strict mode: %v", err)
	}
}

func TestRunInterpolate_BadEnvFile(t *testing.T) {
	err := runInterpolate([]string{"-env", "/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing env file")
	}
}
