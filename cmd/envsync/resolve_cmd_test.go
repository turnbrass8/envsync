package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempResolveEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunResolve_BasicResolution(t *testing.T) {
	f1 := writeTempResolveEnv(t, "DB_HOST=prod\nAPI_KEY=secret\n")
	var buf bytes.Buffer
	err := runResolve([]string{"--files", f1, "--keys", "DB_HOST"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "DB_HOST=prod") {
		t.Errorf("expected DB_HOST=prod in output, got: %s", buf.String())
	}
}

func TestRunResolve_PrecedenceFirstWins(t *testing.T) {
	f1 := writeTempResolveEnv(t, "DB_HOST=primary\n")
	f2 := writeTempResolveEnv(t, "DB_HOST=secondary\nEXTRA=yes\n")
	var buf bytes.Buffer
	err := runResolve([]string{"--files", f1 + "," + f2, "--keys", "DB_HOST"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "DB_HOST=primary") {
		t.Errorf("expected primary to win, got: %s", buf.String())
	}
}

func TestRunResolve_StrictMissingReturnsError(t *testing.T) {
	f1 := writeTempResolveEnv(t, "PRESENT=yes\n")
	var buf bytes.Buffer
	err := runResolve([]string{"--files", f1, "--keys", "MISSING", "--strict"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing key in strict mode")
	}
}

func TestRunResolve_OriginFlag(t *testing.T) {
	f1 := writeTempResolveEnv(t, "TOKEN=abc\n")
	var buf bytes.Buffer
	err := runResolve([]string{"--files", f1, "--keys", "TOKEN", "--origin"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "#") {
		t.Errorf("expected origin comment in output, got: %s", buf.String())
	}
}

func TestRunResolve_MissingFilesFlag_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	err := runResolve([]string{"--keys", "X"}, &buf)
	if err == nil {
		t.Fatal("expected error when --files is absent")
	}
}
