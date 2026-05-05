package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempDefaultsEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunDefaults_AppliesMissingKeys(t *testing.T) {
	target := writeTempDefaultsEnv(t, "EXISTING=yes\n")
	defs := writeTempDefaultsEnv(t, "EXISTING=other\nNEW_KEY=hello\n")

	err := runDefaults([]string{"--defaults", defs, target})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(target)
	body := string(data)
	if !strings.Contains(body, "NEW_KEY=hello") {
		t.Errorf("expected NEW_KEY in output, got:\n%s", body)
	}
	if !strings.Contains(body, "EXISTING=yes") {
		t.Errorf("expected EXISTING unchanged, got:\n%s", body)
	}
}

func TestRunDefaults_DryRunDoesNotWrite(t *testing.T) {
	target := writeTempDefaultsEnv(t, "A=1\n")
	defs := writeTempDefaultsEnv(t, "B=2\n")

	err := runDefaults([]string{"--dry-run", "--defaults", defs, target})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(target)
	if strings.Contains(string(data), "B=2") {
		t.Error("dry-run should not write to file")
	}
}

func TestRunDefaults_MissingDefaultsFlag_ReturnsError(t *testing.T) {
	target := writeTempDefaultsEnv(t, "A=1\n")
	err := runDefaults([]string{target})
	if err == nil || !strings.Contains(err.Error(), "--defaults") {
		t.Errorf("expected --defaults error, got %v", err)
	}
}

func TestRunDefaults_BadTargetFile_ReturnsError(t *testing.T) {
	defs := writeTempDefaultsEnv(t, "KEY=val\n")
	err := runDefaults([]string{"--defaults", defs, "/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for bad target file")
	}
}
