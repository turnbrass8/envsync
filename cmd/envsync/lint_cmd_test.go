package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempLintEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestRunLint_CleanFile_NoError(t *testing.T) {
	p := writeTempLintEnv(t, "APP_PORT=8080\nDATABASE_URL=postgres://localhost/db\n")
	if err := runLint([]string{p}); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestRunLint_ShortSecret_ReturnsError(t *testing.T) {
	p := writeTempLintEnv(t, "APP_SECRET=abc\n")
	if err := runLint([]string{p}); err == nil {
		t.Error("expected error for short secret, got nil")
	}
}

func TestRunLint_ErrorsOnlyFlag_SkipsWarnings(t *testing.T) {
	// lowercase key produces a warning; short secret produces an error
	p := writeTempLintEnv(t, "app_name=foo\nAPP_TOKEN=toolongvalue\n")
	// errors-only: only warnings present (app_name), no errors => no error returned
	if err := runLint([]string{"-errors-only", p}); err != nil {
		t.Errorf("expected no error with errors-only flag, got: %v", err)
	}
}

func TestRunLint_MissingFile_ReturnsError(t *testing.T) {
	if err := runLint([]string{"/nonexistent/.env"}); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRunLint_NoArgs_ReturnsError(t *testing.T) {
	if err := runLint([]string{}); err == nil {
		t.Error("expected error when no args provided")
	}
}
