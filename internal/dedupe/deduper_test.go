package dedupe_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envsync/internal/dedupe"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readFile: %v", err)
	}
	return string(b)
}

func TestDedupe_NoDuplicates_NoChanges(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	res, err := dedupe.Dedupe(p, dedupe.Options{Strategy: dedupe.StrategyLast})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected no removals, got %v", res.Removed)
	}
}

func TestDedupe_StrategyLast_KeepsLastValue(t *testing.T) {
	p := writeTempEnv(t, "FOO=first\nBAR=keep\nFOO=last\n")
	res, err := dedupe.Dedupe(p, dedupe.Options{Strategy: dedupe.StrategyLast})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "FOO" {
		t.Errorf("expected FOO removed, got %v", res.Removed)
	}
	got := readFile(t, p)
	if !strings.Contains(got, "FOO=last") {
		t.Errorf("expected FOO=last in output, got:\n%s", got)
	}
	if strings.Contains(got, "FOO=first") {
		t.Errorf("FOO=first should have been removed, got:\n%s", got)
	}
}

func TestDedupe_StrategyFirst_KeepsFirstValue(t *testing.T) {
	p := writeTempEnv(t, "FOO=first\nBAR=keep\nFOO=last\n")
	res, err := dedupe.Dedupe(p, dedupe.Options{Strategy: dedupe.StrategyFirst})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 1 {
		t.Errorf("expected 1 removal, got %v", res.Removed)
	}
	got := readFile(t, p)
	if !strings.Contains(got, "FOO=first") {
		t.Errorf("expected FOO=first in output, got:\n%s", got)
	}
}

func TestDedupe_DryRun_DoesNotWrite(t *testing.T) {
	original := "FOO=first\nFOO=last\n"
	p := writeTempEnv(t, original)
	_, err := dedupe.Dedupe(p, dedupe.Options{Strategy: dedupe.StrategyLast, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := readFile(t, p); got != original {
		t.Errorf("dry run should not modify file; got:\n%s", got)
	}
}

func TestDedupe_BadFile_ReturnsError(t *testing.T) {
	_, err := dedupe.Dedupe("/nonexistent/.env", dedupe.Options{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
