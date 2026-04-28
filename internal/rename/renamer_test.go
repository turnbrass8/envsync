package rename_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/rename"
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

func TestRename_AppliesRule(t *testing.T) {
	path := writeTempEnv(t, "OLD_KEY=hello\nOTHER=world\n")
	rules := []rename.Rule{{From: "OLD_KEY", To: "NEW_KEY"}}

	res, err := rename.Rename(path, rules, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 1 || res.Applied[0].From != "OLD_KEY" {
		t.Errorf("expected Applied=[{OLD_KEY NEW_KEY}], got %v", res.Applied)
	}

	env, _ := envfile.Parse(path)
	if _, ok := env.Get("OLD_KEY"); ok {
		t.Error("OLD_KEY should have been removed")
	}
	if v, ok := env.Get("NEW_KEY"); !ok || v != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q ok=%v", v, ok)
	}
}

func TestRename_SkipsMissingKey(t *testing.T) {
	path := writeTempEnv(t, "OTHER=world\n")
	rules := []rename.Rule{{From: "MISSING", To: "NEW_KEY"}}

	res, err := rename.Rename(path, rules, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0].From != "MISSING" {
		t.Errorf("expected Skipped=[{MISSING NEW_KEY}], got %v", res.Skipped)
	}
}

func TestRename_DryRunDoesNotWrite(t *testing.T) {
	original := "OLD_KEY=hello\n"
	path := writeTempEnv(t, original)
	rules := []rename.Rule{{From: "OLD_KEY", To: "NEW_KEY"}}

	_, err := rename.Rename(path, rules, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	if string(data) != original {
		t.Errorf("dry run should not modify file, got: %q", string(data))
	}
}

func TestRename_ErrorsOnConflict(t *testing.T) {
	path := writeTempEnv(t, "OLD_KEY=hello\nNEW_KEY=already\n")
	rules := []rename.Rule{{From: "OLD_KEY", To: "NEW_KEY"}}

	_, err := rename.Rename(path, rules, false)
	if err == nil {
		t.Error("expected error when target key already exists")
	}
}
