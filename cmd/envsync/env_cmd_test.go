package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunEnv_PrintsKeyValues(t *testing.T) {
	p := writeTempEnvFile(t, "APP=hello\nDEBUG=true\n")
	var buf bytes.Buffer
	if err := runEnv([]string{"-f", p, "-no-os"}, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP=hello") {
		t.Errorf("expected APP=hello in output, got: %s", out)
	}
	if !strings.Contains(out, "DEBUG=true") {
		t.Errorf("expected DEBUG=true in output, got: %s", out)
	}
}

func TestRunEnv_FilterKey_ReturnsValue(t *testing.T) {
	p := writeTempEnvFile(t, "APP=hello\nDEBUG=true\n")
	var buf bytes.Buffer
	if err := runEnv([]string{"-f", p, "-no-os", "-key", "APP"}, &buf); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(buf.String()) != "hello" {
		t.Errorf("expected 'hello', got %q", buf.String())
	}
}

func TestRunEnv_FilterKey_Missing_ReturnsError(t *testing.T) {
	p := writeTempEnvFile(t, "APP=hello\n")
	var buf bytes.Buffer
	err := runEnv([]string{"-f", p, "-no-os", "-key", "MISSING"}, &buf)
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestRunEnv_ShowOrigin_IncludesSource(t *testing.T) {
	p := writeTempEnvFile(t, "MYKEY=myval\n")
	var buf bytes.Buffer
	if err := runEnv([]string{"-f", p, "-no-os", "-origin"}, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "[file]") {
		t.Errorf("expected origin [file] in output, got: %s", out)
	}
}

func TestRunEnv_MissingFile_StillRunsWithOSEnv(t *testing.T) {
	var buf bytes.Buffer
	// Should not error — missing file is silently skipped, OS env used
	if err := runEnv([]string{"-f", "/nonexistent/.env"}, &buf); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
