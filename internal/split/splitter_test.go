package split_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsync/internal/split"
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

func TestSplit_ByPrefix_WritesCorrectFiles(t *testing.T) {
	src := writeTempEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\nDB_HOST=pg\nDB_PORT=5432\n")
	dir := t.TempDir()
	appOut := filepath.Join(dir, "app.env")
	dbOut := filepath.Join(dir, "db.env")

	res, err := split.Split(src, split.Options{
		Prefixes: map[string]string{"APP_": appOut, "DB_": dbOut},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Written[appOut] != 2 {
		t.Errorf("expected 2 app keys, got %d", res.Written[appOut])
	}
	if res.Written[dbOut] != 2 {
		t.Errorf("expected 2 db keys, got %d", res.Written[dbOut])
	}
}

func TestSplit_StripPrefix_RemovesPrefixFromKeys(t *testing.T) {
	src := writeTempEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\n")
	dir := t.TempDir()
	appOut := filepath.Join(dir, "app.env")

	_, err := split.Split(src, split.Options{
		Prefixes:    map[string]string{"APP_": appOut},
		StripPrefix: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(appOut)
	content := string(data)
	if !containsKey(content, "HOST") {
		t.Errorf("expected stripped key HOST in output, got:\n%s", content)
	}
}

func TestSplit_DryRun_DoesNotWriteFiles(t *testing.T) {
	src := writeTempEnv(t, "APP_HOST=localhost\n")
	dir := t.TempDir()
	appOut := filepath.Join(dir, "app.env")

	res, err := split.Split(src, split.Options{
		Prefixes: map[string]string{"APP_": appOut},
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Written[appOut] != 1 {
		t.Errorf("expected dry-run count 1, got %d", res.Written[appOut])
	}
	if _, err := os.Stat(appOut); !os.IsNotExist(err) {
		t.Error("expected file not to be created in dry-run mode")
	}
}

func TestSplit_UnmatchedKeys_ReportedInResult(t *testing.T) {
	src := writeTempEnv(t, "APP_HOST=localhost\nSECRET=abc\n")
	dir := t.TempDir()
	appOut := filepath.Join(dir, "app.env")

	res, err := split.Split(src, split.Options{
		Prefixes: map[string]string{"APP_": appOut},
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Unmatched) != 1 || res.Unmatched[0] != "SECRET" {
		t.Errorf("expected [SECRET] unmatched, got %v", res.Unmatched)
	}
}

func TestSplit_NoPrefixes_ReturnsError(t *testing.T) {
	src := writeTempEnv(t, "KEY=val\n")
	_, err := split.Split(src, split.Options{})
	if err == nil {
		t.Error("expected error when no prefixes provided")
	}
}

func containsKey(content, key string) bool {
	for _, line := range splitLines(content) {
		if len(line) >= len(key) && line[:len(key)] == key {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
