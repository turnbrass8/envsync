package clone_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsync/internal/clone"
	"github.com/user/envsync/internal/envfile"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestClone_CopiesKeys(t *testing.T) {
	src := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	dst := writeTempEnv(t, "EXISTING=yes\n")

	res, err := clone.Clone(src, dst, clone.Options{
		Rules: []clone.Rule{{SrcKey: "FOO"}, {SrcKey: "BAZ"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Fatalf("expected 2 copied, got %d", len(res.Copied))
	}

	env, _ := envfile.Parse(dst)
	if env.Values["FOO"] != "bar" {
		t.Errorf("FOO: got %q, want %q", env.Values["FOO"], "bar")
	}
	if env.Values["EXISTING"] != "yes" {
		t.Error("existing key should be preserved")
	}
}

func TestClone_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := writeTempEnv(t, "FOO=new\n")
	dst := writeTempEnv(t, "FOO=old\n")

	res, err := clone.Clone(src, dst, clone.Options{
		Rules:     []clone.Rule{{SrcKey: "FOO"}},
		Overwrite: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(res.Skipped))
	}

	env, _ := envfile.Parse(dst)
	if env.Values["FOO"] != "old" {
		t.Error("value should not have been overwritten")
	}
}

func TestClone_RenamesKey(t *testing.T) {
	src := writeTempEnv(t, "DB_PASS=secret\n")
	dst := writeTempEnv(t, "")

	_, err := clone.Clone(src, dst, clone.Options{
		Rules: []clone.Rule{{SrcKey: "DB_PASS", DstKey: "DATABASE_PASSWORD"}},
	})
	if err != nil {
		t.Fatal(err)
	}

	env, _ := envfile.Parse(dst)
	if env.Values["DATABASE_PASSWORD"] != "secret" {
		t.Errorf("renamed key not found or wrong value: %v", env.Values)
	}
}

func TestClone_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempEnv(t, "FOO=bar\n")
	dstDir := t.TempDir()
	dst := filepath.Join(dstDir, "out.env")

	_, err := clone.Clone(src, dst, clone.Options{
		Rules:  []clone.Rule{{SrcKey: "FOO"}},
		DryRun: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		t.Error("dry-run should not create the destination file")
	}
}

func TestClone_MissingSourceKey_ReturnsError(t *testing.T) {
	src := writeTempEnv(t, "FOO=bar\n")
	dst := writeTempEnv(t, "")

	_, err := clone.Clone(src, dst, clone.Options{
		Rules: []clone.Rule{{SrcKey: "MISSING"}},
	})
	if err == nil {
		t.Fatal("expected error for missing source key")
	}
}
