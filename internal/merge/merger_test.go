package merge_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsync/internal/merge"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestMerge_StrategyLast_OverwritesConflicts(t *testing.T) {
	a := writeTempEnv(t, "FOO=alpha\nBAR=base\n")
	b := writeTempEnv(t, "FOO=beta\nBAZ=extra\n")

	res, err := merge.Merge([]string{a, b}, merge.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["FOO"] != "beta" {
		t.Errorf("expected FOO=beta, got %q", res.Env["FOO"])
	}
	if res.Env["BAR"] != "base" {
		t.Errorf("expected BAR=base, got %q", res.Env["BAR"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0].Key != "FOO" {
		t.Errorf("expected 1 conflict on FOO, got %+v", res.Conflicts)
	}
}

func TestMerge_StrategyFirst_KeepsOriginal(t *testing.T) {
	a := writeTempEnv(t, "FOO=first\n")
	b := writeTempEnv(t, "FOO=second\n")

	res, err := merge.Merge([]string{a, b}, merge.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["FOO"] != "first" {
		t.Errorf("expected FOO=first, got %q", res.Env["FOO"])
	}
}

func TestMerge_StrategyError_ReturnsErrorOnConflict(t *testing.T) {
	a := writeTempEnv(t, "KEY=one\n")
	b := writeTempEnv(t, "KEY=two\n")

	_, err := merge.Merge([]string{a, b}, merge.StrategyError)
	if err == nil {
		t.Fatal("expected error for conflict, got nil")
	}
}

func TestMerge_NoConflicts(t *testing.T) {
	a := writeTempEnv(t, "A=1\n")
	b := writeTempEnv(t, "B=2\n")

	res, err := merge.Merge([]string{a, b}, merge.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %+v", res.Conflicts)
	}
	if res.Env["A"] != "1" || res.Env["B"] != "2" {
		t.Errorf("unexpected env: %+v", res.Env)
	}
}

func TestMerge_NoPaths_ReturnsError(t *testing.T) {
	_, err := merge.Merge([]string{}, merge.StrategyFirst)
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestMerge_BadFile_ReturnsError(t *testing.T) {
	_, err := merge.Merge([]string{filepath.Join(t.TempDir(), "missing.env")}, merge.StrategyFirst)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
