package snapshot_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsync/internal/snapshot"
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

func TestCapture_ReadsValues(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	snap, err := snapshot.Capture(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Values["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", snap.Values["FOO"])
	}
	if snap.Source != path {
		t.Errorf("expected source %q, got %q", path, snap.Source)
	}
}

func TestCapture_BadFile_ReturnsError(t *testing.T) {
	_, err := snapshot.Capture("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := writeTempEnv(t, "KEY=value\n")
	snap, err := snapshot.Capture(path)
	if err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshot.Save(snap, dest); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	loaded, err := snapshot.Load(dest)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if loaded.Values["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", loaded.Values["KEY"])
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "*.json")
	f.WriteString("not json")
	f.Close()
	_, err := snapshot.Load(f.Name())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestCompare_DetectsChanges(t *testing.T) {
	old := &snapshot.Snapshot{Values: map[string]string{"A": "1", "B": "2", "C": "3"}}
	new := &snapshot.Snapshot{Values: map[string]string{"A": "1", "B": "changed", "D": "4"}}
	diff := snapshot.Compare(old, new)
	if !diff.HasChanges() {
		t.Fatal("expected changes")
	}
	if diff.Added["D"] != "4" {
		t.Errorf("expected D added")
	}
	if diff.Removed["C"] != "3" {
		t.Errorf("expected C removed")
	}
	if diff.Changed["B"] != ([2]string{"2", "changed"}) {
		t.Errorf("expected B changed")
	}
	_ = json.Marshal // ensure encoding import used in main file
}

func TestCompare_NoChanges(t *testing.T) {
	old := &snapshot.Snapshot{Values: map[string]string{"X": "y"}}
	new := &snapshot.Snapshot{Values: map[string]string{"X": "y"}}
	if snapshot.Compare(old, new).HasChanges() {
		t.Fatal("expected no changes")
	}
}
