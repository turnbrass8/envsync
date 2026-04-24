package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestRunAudit_AllPresent(t *testing.T) {
	dir := t.TempDir()
	env := writeTempFile(t, dir, ".env", "DB_HOST=localhost\nDB_PORT=5432\n")
	mf := writeTempFile(t, dir, "manifest.env", "DB_HOST\nDB_PORT\n")

	if err := runAudit(env, mf); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestRunAudit_MissingWithDefault(t *testing.T) {
	dir := t.TempDir()
	env := writeTempFile(t, dir, ".env", "DB_HOST=localhost\n")
	mf := writeTempFile(t, dir, "manifest.env", "DB_HOST\nDB_PORT default=5432\n")

	// Missing key has a default, so runAudit should succeed.
	if err := runAudit(env, mf); err != nil {
		t.Errorf("expected no error when default covers missing key, got: %v", err)
	}
}

func TestRunAudit_MissingRequired_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	env := writeTempFile(t, dir, ".env", "DB_HOST=localhost\n")
	mf := writeTempFile(t, dir, "manifest.env", "DB_HOST\nSECRET_KEY required\n")

	err := runAudit(env, mf)
	if err == nil {
		t.Error("expected error for missing required key, got nil")
	}
}

func TestRunAudit_BadEnvFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	mf := writeTempFile(t, dir, "manifest.env", "KEY\n")

	err := runAudit(filepath.Join(dir, "nonexistent.env"), mf)
	if err == nil {
		t.Error("expected error for missing env file")
	}
}
