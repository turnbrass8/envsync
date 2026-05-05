package patch_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/envsync/internal/patch"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestPatch_UpdatesExistingKey(t *testing.T) {
	p := writeTempEnv(t, "FOO=old\nBAR=keep\n")
	res, err := patch.Patch(p, []patch.Op{{Key: "FOO", Value: "new"}}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Updated) != 1 || res.Updated[0] != "FOO" {
		t.Fatalf("expected FOO updated, got %v", res.Updated)
	}
	if !strings.Contains(readFile(t, p), "FOO=new") {
		t.Error("file should contain FOO=new")
	}
}

func TestPatch_AddsNewKey(t *testing.T) {
	p := writeTempEnv(t, "EXISTING=1\n")
	res, err := patch.Patch(p, []patch.Op{{Key: "NEW_KEY", Value: "hello"}}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Added) != 1 || res.Added[0] != "NEW_KEY" {
		t.Fatalf("expected NEW_KEY added, got %v", res.Added)
	}
	if !strings.Contains(readFile(t, p), "NEW_KEY=hello") {
		t.Error("file should contain NEW_KEY=hello")
	}
}

func TestPatch_DeletesKey(t *testing.T) {
	p := writeTempEnv(t, "REMOVE=yes\nKEEP=no\n")
	res, err := patch.Patch(p, []patch.Op{{Key: "REMOVE", Delete: true}}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Deleted) != 1 {
		t.Fatalf("expected 1 deletion, got %v", res.Deleted)
	}
	if strings.Contains(readFile(t, p), "REMOVE") {
		t.Error("file should not contain REMOVE")
	}
}

func TestPatch_DryRunDoesNotWrite(t *testing.T) {
	original := "FOO=original\n"
	p := writeTempEnv(t, original)
	_, err := patch.Patch(p, []patch.Op{{Key: "FOO", Value: "changed"}}, true)
	if err != nil {
		t.Fatal(err)
	}
	if readFile(t, p) != original {
		t.Error("dry run should not modify file")
	}
}

func TestPatch_NoOps_ReturnsError(t *testing.T) {
	p := writeTempEnv(t, "A=1\n")
	_, err := patch.Patch(p, nil, false)
	if err == nil {
		t.Error("expected error for empty ops")
	}
}

func TestPatch_MissingFile_ReturnsError(t *testing.T) {
	_, err := patch.Patch("/nonexistent/.env", []patch.Op{{Key: "X", Value: "1"}}, false)
	if err == nil {
		t.Error("expected error for missing file")
	}
}
