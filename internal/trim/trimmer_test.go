package trim_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsync/internal/manifest"
	"github.com/user/envsync/internal/trim"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func buildManifest(keys ...string) *manifest.Manifest {
	mf := &manifest.Manifest{}
	for _, k := range keys {
		mf.Keys = append(mf.Keys, manifest.Key{Name: k})
	}
	return mf
}

func TestTrim_RemovesStaleKeys(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nSTALE_KEY=old\nDEBUG=true\n")
	mf := buildManifest("APP_NAME", "DEBUG")

	res, err := trim.Trim(path, mf, trim.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "STALE_KEY" {
		t.Errorf("expected [STALE_KEY] removed, got %v", res.Removed)
	}
	if _, ok := res.Retained["APP_NAME"]; !ok {
		t.Error("APP_NAME should be retained")
	}
}

func TestTrim_NoStaleKeys(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nDEBUG=true\n")
	mf := buildManifest("APP_NAME", "DEBUG")

	res, err := trim.Trim(path, mf, trim.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected no removals, got %v", res.Removed)
	}
	if res.Summary() != "no stale keys found" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestTrim_DryRunDoesNotWrite(t *testing.T) {
	original := "APP_NAME=myapp\nSTALE=bad\n"
	path := writeTempEnv(t, original)
	mf := buildManifest("APP_NAME")

	res, err := trim.Trim(path, mf, trim.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 1 {
		t.Fatalf("expected 1 removal, got %d", len(res.Removed))
	}

	// File must be unchanged.
	got, _ := os.ReadFile(path)
	if string(got) != original {
		t.Errorf("dry run mutated the file")
	}
}

func TestTrim_BadEnvFile_ReturnsError(t *testing.T) {
	_, err := trim.Trim("/nonexistent/.env", buildManifest("KEY"), trim.Options{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
