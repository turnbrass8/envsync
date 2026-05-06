package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempCopyEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestRunCopy_CopiesAllKeys(t *testing.T) {
	src := writeTempCopyEnv(t, "FOO=1\nBAR=2\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	if err := runCopy([]string{src, dst}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := os.ReadFile(dst)
	if !strings.Contains(string(b), "FOO=1") {
		t.Errorf("expected FOO in dst, got: %s", string(b))
	}
}

func TestRunCopy_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempCopyEnv(t, "FOO=1\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	if err := runCopy([]string{"-dry-run", src, dst}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		t.Error("dst should not exist after dry-run")
	}
}

func TestRunCopy_SpecificKeys(t *testing.T) {
	src := writeTempCopyEnv(t, "FOO=1\nBAR=2\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	if err := runCopy([]string{"-keys", "FOO", src, dst}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := os.ReadFile(dst)
	if strings.Contains(string(b), "BAR") {
		t.Errorf("BAR should not be in dst, got: %s", string(b))
	}
}

func TestRunCopy_MissingArgs_ReturnsError(t *testing.T) {
	if err := runCopy([]string{"only-one-arg"}); err == nil {
		t.Error("expected error for missing dst argument")
	}
}
