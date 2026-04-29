package pin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/envsync/internal/pin"
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

func TestPin_PinsExistingKeys(t *testing.T) {
	env := writeTempEnv(t, "DB_PASS=secret\nAPI_KEY=abc123\n")
	res, err := pin.Pin(env, []string{"DB_PASS", "API_KEY"}, pin.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) != 2 {
		t.Fatalf("expected 2 pinned, got %d", len(res.Pinned))
	}
	if len(res.Skipped) != 0 {
		t.Fatalf("expected 0 skipped, got %d", len(res.Skipped))
	}
}

func TestPin_SkipsMissingKey(t *testing.T) {
	env := writeTempEnv(t, "DB_PASS=secret\n")
	res, err := pin.Pin(env, []string{"DB_PASS", "MISSING"}, pin.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Fatalf("expected MISSING in skipped, got %v", res.Skipped)
	}
}

func TestPin_DryRunDoesNotWriteFile(t *testing.T) {
	env := writeTempEnv(t, "TOKEN=xyz\n")
	_, err := pin.Pin(env, []string{"TOKEN"}, pin.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(env + ".pinned"); !os.IsNotExist(statErr) {
		t.Fatal("expected pin file not to exist in dry-run mode")
	}
}

func TestPin_NoKeysReturnsError(t *testing.T) {
	env := writeTempEnv(t, "A=1\n")
	_, err := pin.Pin(env, []string{}, pin.Options{})
	if err == nil {
		t.Fatal("expected error for empty key list")
	}
}

func TestLoadPins_RoundTrip(t *testing.T) {
	env := writeTempEnv(t, "SECRET=hunter2\nTOKEN=tok\n")
	_, err := pin.Pin(env, []string{"SECRET", "TOKEN"}, pin.Options{})
	if err != nil {
		t.Fatalf("pin error: %v", err)
	}
	pins, err := pin.LoadPins(env)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if pins["SECRET"] != "hunter2" {
		t.Errorf("expected hunter2, got %q", pins["SECRET"])
	}
	if pins["TOKEN"] != "tok" {
		t.Errorf("expected tok, got %q", pins["TOKEN"])
	}
}

func TestLoadPins_NonExistentFile_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	pins, err := pin.LoadPins(filepath.Join(dir, "missing.env"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pins) != 0 {
		t.Fatalf("expected empty map, got %v", pins)
	}
}
