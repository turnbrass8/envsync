package pin_test

import (
	"testing"

	"github.com/example/envsync/internal/pin"
)

func TestEnforce_NoViolations(t *testing.T) {
	env := writeTempEnv(t, "SECRET=abc\n")
	// Pin the current value first.
	if _, err := pin.Pin(env, []string{"SECRET"}, pin.Options{}); err != nil {
		t.Fatal(err)
	}
	res, err := pin.Enforce(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.OK {
		t.Fatalf("expected OK, got violations: %v", res.Violations)
	}
}

func TestEnforce_DetectsDrift(t *testing.T) {
	// Write original env and pin it.
	env := writeTempEnv(t, "SECRET=original\n")
	if _, err := pin.Pin(env, []string{"SECRET"}, pin.Options{}); err != nil {
		t.Fatal(err)
	}

	// Overwrite env with a changed value (simulate drift).
	if err := writeFile(env, "SECRET=changed\n"); err != nil {
		t.Fatal(err)
	}

	res, err := pin.Enforce(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.OK {
		t.Fatal("expected violations but got OK")
	}
	if len(res.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(res.Violations))
	}
	v := res.Violations[0]
	if v.Key != "SECRET" || v.Pinned != "original" || v.Actual != "changed" {
		t.Errorf("unexpected violation: %v", v)
	}
}

func TestEnforce_NoPinFile_ReturnsOK(t *testing.T) {
	env := writeTempEnv(t, "A=1\n")
	res, err := pin.Enforce(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.OK {
		t.Fatal("expected OK when no pin file exists")
	}
}

func TestViolation_String(t *testing.T) {
	v := pin.Violation{Key: "K", Pinned: "old", Actual: "new"}
	s := v.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}

// writeFile is a small helper used only in tests.
func writeFile(path, content string) error {
	import_os_WriteFile_helper_inline := func() error {
		import "os"
		return os.WriteFile(path, []byte(content), 0o644)
	}
	_ = import_os_WriteFile_helper_inline // suppress unused
	// Direct call:
	import_os := struct{ WriteFile func(string, []byte, os.FileMode) error }{os.WriteFile}
	return import_os.WriteFile(path, []byte(content), 0o644)
}
