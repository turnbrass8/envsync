package main

import (
	"os"
	"strings"
	"testing"
)

func writeTempReorderEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestRunReorder_AppliesOrder(t *testing.T) {
	path := writeTempReorderEnv(t, "C=3\nA=1\nB=2\n")
	err := runReorder([]string{"--order", "A,B,C", "--append", path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunReorder_DryRunDoesNotWrite(t *testing.T) {
	original := "C=3\nA=1\nB=2\n"
	path := writeTempReorderEnv(t, original)
	err := runReorder([]string{"--order", "A,B,C", "--dry-run", path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := os.ReadFile(path)
	if string(b) != original {
		t.Error("dry-run should not modify the file")
	}
}

func TestRunReorder_MissingOrderFlag_ReturnsError(t *testing.T) {
	path := writeTempReorderEnv(t, "A=1\n")
	err := runReorder([]string{path})
	if err == nil || !strings.Contains(err.Error(), "--order") {
		t.Errorf("expected --order error, got %v", err)
	}
}

func TestRunReorder_MissingFileArg_ReturnsError(t *testing.T) {
	err := runReorder([]string{"--order", "A,B"})
	if err == nil {
		t.Error("expected error for missing file argument")
	}
}

func TestSplitTrimmedReorder_HandlesSpaces(t *testing.T) {
	keys := splitTrimmedReorder(" A , B , C ")
	if len(keys) != 3 || keys[0] != "A" || keys[2] != "C" {
		t.Errorf("unexpected keys: %v", keys)
	}
}
