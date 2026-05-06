package reorder_test

import (
	"os"
	"strings"
	"testing"

	"github.com/user/envsync/internal/reorder"
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

func TestReorder_AppliesOrder(t *testing.T) {
	path := writeTempEnv(t, "C=3\nA=1\nB=2\n")
	res, err := reorder.Reorder(path, reorder.Options{
		Order:  []string{"A", "B", "C"},
		Append: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Reordered != 3 {
		t.Errorf("expected 3 reordered, got %d", res.Reordered)
	}
	content := readFile(t, path)
	if idx := strings.Index(content, "A="); idx == -1 {
		t.Error("expected A in output")
	}
}

func TestReorder_DryRunDoesNotWrite(t *testing.T) {
	original := "C=3\nA=1\nB=2\n"
	path := writeTempEnv(t, original)
	_, err := reorder.Reorder(path, reorder.Options{
		Order:  []string{"A", "B", "C"},
		DryRun: true,
		Append: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != original {
		t.Errorf("dry run modified file: got %q", got)
	}
}

func TestReorder_DropsUnlistedKeys(t *testing.T) {
	path := writeTempEnv(t, "A=1\nB=2\nC=3\n")
	res, err := reorder.Reorder(path, reorder.Options{
		Order:  []string{"A", "B"},
		Append: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Dropped != 1 {
		t.Errorf("expected 1 dropped, got %d", res.Dropped)
	}
	content := readFile(t, path)
	if strings.Contains(content, "C=") {
		t.Error("expected C to be dropped")
	}
}

func TestReorder_AppendsUnlistedKeys(t *testing.T) {
	path := writeTempEnv(t, "A=1\nB=2\nC=3\n")
	res, err := reorder.Reorder(path, reorder.Options{
		Order:  []string{"A"},
		Append: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Appended != 2 {
		t.Errorf("expected 2 appended, got %d", res.Appended)
	}
}

func TestReorder_BadFile_ReturnsError(t *testing.T) {
	_, err := reorder.Reorder("/no/such/file.env", reorder.Options{Order: []string{"A"}})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestResult_Summary(t *testing.T) {
	r := reorder.Result{Reordered: 3, Appended: 1, Dropped: 2}
	s := r.Summary()
	if !strings.Contains(s, "reordered=3") {
		t.Errorf("unexpected summary: %s", s)
	}
}
