package archive_test

import (
	"archive/zip"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsync/internal/archive"
)

func writeTempEnv(t *testing.T, name, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestArchive_CreatesZipWithFiles(t *testing.T) {
	a := writeTempEnv(t, "a.env", "FOO=bar\n")
	b := writeTempEnv(t, "b.env", "BAZ=qux\n")
	dest := filepath.Join(t.TempDir(), "out.zip")

	meta, err := archive.Archive(dest, []string{a, b}, archive.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(meta.Files) != 2 {
		t.Fatalf("expected 2 files in metadata, got %d", len(meta.Files))
	}

	zr, err := zip.OpenReader(dest)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer zr.Close()

	names := map[string]bool{}
	for _, f := range zr.File {
		names[f.Name] = true
	}
	for _, want := range []string{filepath.Base(a), filepath.Base(b), "envsync-manifest.json"} {
		if !names[want] {
			t.Errorf("missing entry %q in zip", want)
		}
	}
}

func TestArchive_ManifestContainsLabels(t *testing.T) {
	a := writeTempEnv(t, "a.env", "KEY=val\n")
	dest := filepath.Join(t.TempDir(), "out.zip")

	_, err := archive.Archive(dest, []string{a}, archive.Options{
		Labels: map[string]string{"env": "staging"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	zr, err := zip.OpenReader(dest)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer zr.Close()

	for _, f := range zr.File {
		if f.Name != "envsync-manifest.json" {
			continue
		}
		rc, _ := f.Open()
		var meta archive.Metadata
		if err := json.NewDecoder(rc).Decode(&meta); err != nil {
			t.Fatalf("decode manifest: %v", err)
		}
		rc.Close()
		if meta.Labels["env"] != "staging" {
			t.Errorf("expected label env=staging, got %v", meta.Labels)
		}
		return
	}
	t.Error("manifest not found in zip")
}

func TestArchive_DryRunDoesNotCreateFile(t *testing.T) {
	a := writeTempEnv(t, "a.env", "X=1\n")
	dest := filepath.Join(t.TempDir(), "out.zip")

	_, err := archive.Archive(dest, []string{a}, archive.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Error("expected zip file to not exist in dry-run mode")
	}
}

func TestArchive_NoFiles_ReturnsError(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "out.zip")
	_, err := archive.Archive(dest, nil, archive.Options{})
	if err == nil {
		t.Error("expected error for empty file list")
	}
}

func TestArchive_MissingSourceFile_ReturnsError(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "out.zip")
	_, err := archive.Archive(dest, []string{"/nonexistent/path.env"}, archive.Options{})
	if err == nil {
		t.Error("expected error for missing source file")
	}
}
