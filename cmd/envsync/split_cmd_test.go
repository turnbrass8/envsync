package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempSplitEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestRunSplit_WritesOutputFiles(t *testing.T) {
	src := writeTempSplitEnv(t, "APP_HOST=localhost\nDB_HOST=pg\n")
	dir := t.TempDir()
	appOut := filepath.Join(dir, "app.env")
	dbOut := filepath.Join(dir, "db.env")

	err := runSplit([]string{
		"--src", src,
		"--map", "APP_:" + appOut + ",DB_:" + dbOut,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(appOut); os.IsNotExist(err) {
		t.Error("expected app.env to be created")
	}
	if _, err := os.Stat(dbOut); os.IsNotExist(err) {
		t.Error("expected db.env to be created")
	}
}

func TestRunSplit_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempSplitEnv(t, "APP_KEY=val\n")
	dir := t.TempDir()
	appOut := filepath.Join(dir, "app.env")

	err := runSplit([]string{
		"--src", src,
		"--map", "APP_:" + appOut,
		"--dry-run",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(appOut); !os.IsNotExist(err) {
		t.Error("expected app.env NOT to be created in dry-run")
	}
}

func TestRunSplit_MissingMapFlag_ReturnsError(t *testing.T) {
	src := writeTempSplitEnv(t, "KEY=val\n")
	err := runSplit([]string{"--src", src})
	if err == nil {
		t.Error("expected error when --map is missing")
	}
}

func TestParseSplitMappings_InvalidEntry_ReturnsError(t *testing.T) {
	_, err := parseSplitMappings("NOCOLON")
	if err == nil {
		t.Error("expected error for mapping without colon")
	}
}

func TestParseSplitMappings_ValidInput_ReturnsPrefixes(t *testing.T) {
	m, err := parseSplitMappings("APP_:app.env,DB_:db.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["APP_"] != "app.env" {
		t.Errorf("expected app.env for APP_, got %q", m["APP_"])
	}
	if !strings.Contains(m["DB_"], "db.env") {
		t.Errorf("expected db.env for DB_, got %q", m["DB_"])
	}
}
