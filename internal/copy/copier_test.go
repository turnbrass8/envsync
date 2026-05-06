package copy

import (
	"os"
	"path/filepath"
	"testing"
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

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestCopy_AllKeys(t *testing.T) {
	src := writeTempEnv(t, "FOO=1\nBAR=2\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	res, err := Copy(src, dst, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Fatalf("expected 2 copied, got %d", len(res.Copied))
	}
}

func TestCopy_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := writeTempEnv(t, "FOO=new\n")
	dst := writeTempEnv(t, "FOO=old\n")

	res, err := Copy(src, dst, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if got := readFile(t, dst); got != "FOO=old\n" {
		t.Fatalf("dst modified unexpectedly: %q", got)
	}
}

func TestCopy_OverwriteReplacesExisting(t *testing.T) {
	src := writeTempEnv(t, "FOO=new\n")
	dst := writeTempEnv(t, "FOO=old\n")

	res, err := Copy(src, dst, Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 1 {
		t.Fatalf("expected 1 copied, got %d", len(res.Copied))
	}
}

func TestCopy_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempEnv(t, "FOO=1\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	res, err := Copy(src, dst, Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 1 {
		t.Fatalf("expected 1 reported, got %d", len(res.Copied))
	}
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		t.Fatal("dst should not exist after dry-run")
	}
}

func TestCopy_SpecificKeys(t *testing.T) {
	src := writeTempEnv(t, "FOO=1\nBAR=2\nBAZ=3\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	res, err := Copy(src, dst, Options{Keys: []string{"FOO", "BAZ"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Fatalf("expected 2 copied, got %d", len(res.Copied))
	}
}

func TestResult_Summary(t *testing.T) {
	r := Result{Copied: []string{"A", "B"}, Skipped: []string{"C"}}
	if s := r.Summary(); s != "2 copied, 1 skipped" {
		t.Fatalf("unexpected summary: %q", s)
	}
}
