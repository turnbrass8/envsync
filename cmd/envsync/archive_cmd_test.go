package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempArchiveEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestRunArchive_CreatesZip(t *testing.T) {
	env := writeTempArchiveEnv(t, "FOO=bar\n")
	dest := filepath.Join(t.TempDir(), "result.zip")

	if err := runArchive([]string{"-out", dest, env}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(dest); err != nil {
		t.Errorf("expected zip file at %s: %v", dest, err)
	}

	zr, err := zip.OpenReader(dest)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer zr.Close()

	if len(zr.File) < 2 {
		t.Errorf("expected at least 2 entries in zip, got %d", len(zr.File))
	}
}

func TestRunArchive_DryRunDoesNotWrite(t *testing.T) {
	env := writeTempArchiveEnv(t, "A=1\n")
	dest := filepath.Join(t.TempDir(), "dry.zip")

	if err := runArchive([]string{"-out", dest, "-dry-run", env}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Error("expected no zip file in dry-run mode")
	}
}

func TestRunArchive_NoFiles_ReturnsError(t *testing.T) {
	err := runArchive([]string{"-out", "/tmp/x.zip"})
	if err == nil {
		t.Error("expected error when no files provided")
	}
}

func TestParseLabels_ValidInput(t *testing.T) {
	labels := parseLabels("env=prod,version=2")
	if labels["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", labels["env"])
	}
	if labels["version"] != "2" {
		t.Errorf("expected version=2, got %q", labels["version"])
	}
}

func TestParseLabels_Empty_ReturnsNil(t *testing.T) {
	if parseLabels("") != nil {
		t.Error("expected nil for empty input")
	}
}

func TestRunArchive_WithLabels_DoesNotError(t *testing.T) {
	env := writeTempArchiveEnv(t, "KEY=val\n")
	dest := filepath.Join(t.TempDir(), "labeled.zip")

	err := runArchive([]string{"-out", dest, "-labels", "env=staging", env})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = strings.Contains(dest, "labeled") // suppress unused warning
}
