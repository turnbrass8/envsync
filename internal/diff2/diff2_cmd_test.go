package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempDiff2Env(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp env: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestRunDiff2_NoChanges(t *testing.T) {
	content := "APP_NAME=envsync\nDEBUG=false\n"
	left := writeTempDiff2Env(t, content)
	right := writeTempDiff2Env(t, content)

	var sb strings.Builder
	err := runDiff2(&sb, left, right, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "no differences") {
		t.Errorf("expected 'no differences' message, got: %q", out)
	}
}

func TestRunDiff2_DetectsAddedKey(t *testing.T) {
	left := writeTempDiff2Env(t, "APP_NAME=envsync\n")
	right := writeTempDiff2Env(t, "APP_NAME=envsync\nNEW_KEY=hello\n")

	var sb strings.Builder
	err := runDiff2(&sb, left, right, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("expected NEW_KEY in output, got: %q", out)
	}
}

func TestRunDiff2_DetectsRemovedKey(t *testing.T) {
	left := writeTempDiff2Env(t, "APP_NAME=envsync\nOLD_KEY=bye\n")
	right := writeTempDiff2Env(t, "APP_NAME=envsync\n")

	var sb strings.Builder
	err := runDiff2(&sb, left, right, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "OLD_KEY") {
		t.Errorf("expected OLD_KEY in output, got: %q", out)
	}
}

func TestRunDiff2_DetectsChangedValue(t *testing.T) {
	left := writeTempDiff2Env(t, "PORT=8080\n")
	right := writeTempDiff2Env(t, "PORT=9090\n")

	var sb strings.Builder
	err := runDiff2(&sb, left, right, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got: %q", out)
	}
	if !strings.Contains(out, "8080") || !strings.Contains(out, "9090") {
		t.Errorf("expected both old and new values in output, got: %q", out)
	}
}

func TestRunDiff2_QuietMode_NoOutput(t *testing.T) {
	left := writeTempDiff2Env(t, "APP=foo\n")
	right := writeTempDiff2Env(t, "APP=bar\n")

	var sb strings.Builder
	err := runDiff2(&sb, left, right, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sb.Len() != 0 {
		t.Errorf("expected no output in quiet mode, got: %q", sb.String())
	}
}

func TestRunDiff2_MissingLeftFile_ReturnsError(t *testing.T) {
	right := writeTempDiff2Env(t, "APP=foo\n")
	missing := filepath.Join(t.TempDir(), "nonexistent.env")

	var sb strings.Builder
	err := runDiff2(&sb, missing, right, false)
	if err == nil {
		t.Fatal("expected error for missing left file, got nil")
	}
}

func TestRunDiff2_MissingRightFile_ReturnsError(t *testing.T) {
	left := writeTempDiff2Env(t, "APP=foo\n")
	missing := filepath.Join(t.TempDir(), "nonexistent.env")

	var sb strings.Builder
	err := runDiff2(&sb, left, missing, false)
	if err == nil {
		t.Fatal("expected error for missing right file, got nil")
	}
}
