package sort_test

import (
	"os"
	"strings"
	"testing"

	"github.com/user/envsync/internal/sort"
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

func TestSort_AlphaOrder(t *testing.T) {
	path := writeTempEnv(t, "ZEBRA=1\nAPPLE=2\nMANGO=3\n")
	res, err := sort.Sort(path, sort.Options{Strategy: sort.StrategyAlpha})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Reordered == 0 {
		t.Error("expected some keys to be reordered")
	}
	data, _ := os.ReadFile(path)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if lines[0] != "APPLE=2" {
		t.Errorf("first key should be APPLE, got %q", lines[0])
	}
	if lines[2] != "ZEBRA=1" {
		t.Errorf("last key should be ZEBRA, got %q", lines[2])
	}
}

func TestSort_ReverseOrder(t *testing.T) {
	path := writeTempEnv(t, "APPLE=2\nMANGO=3\nZEBRA=1\n")
	res, err := sort.Sort(path, sort.Options{Strategy: sort.StrategyReverse})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = res
	data, _ := os.ReadFile(path)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if lines[0] != "ZEBRA=1" {
		t.Errorf("first key should be ZEBRA, got %q", lines[0])
	}
}

func TestSort_DryRunDoesNotWrite(t *testing.T) {
	original := "ZEBRA=1\nAPPLE=2\n"
	path := writeTempEnv(t, original)
	_, err := sort.Sort(path, sort.Options{Strategy: sort.StrategyAlpha, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	if string(data) != original {
		t.Error("dry run should not modify the file")
	}
}

func TestSort_AlreadySorted_ZeroReordered(t *testing.T) {
	path := writeTempEnv(t, "ALPHA=1\nBETA=2\nGAMMA=3\n")
	res, err := sort.Sort(path, sort.Options{Strategy: sort.StrategyAlpha})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Reordered != 0 {
		t.Errorf("expected 0 reordered, got %d", res.Reordered)
	}
}

func TestSort_UnknownStrategy_ReturnsError(t *testing.T) {
	path := writeTempEnv(t, "KEY=val\n")
	_, err := sort.Sort(path, sort.Options{Strategy: "bogus"})
	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestSort_MissingFile_ReturnsError(t *testing.T) {
	_, err := sort.Sort("/nonexistent/.env", sort.Options{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
