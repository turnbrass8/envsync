package manifest

import (
	"os"
	"testing"
)

func writeTempManifest(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "manifest-*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParse_BasicKeys(t *testing.T) {
	path := writeTempManifest(t, "APP_ENV\nDB_HOST\nSECRET_KEY\n")
	m, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(m.Entries))
	}
	if m.Entries[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV, got %q", m.Entries[0].Key)
	}
}

func TestParse_RequiredKey(t *testing.T) {
	path := writeTempManifest(t, "SECRET_KEY!\n")
	m, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Entries[0].Required {
		t.Error("expected Required=true for SECRET_KEY!")
	}
	if m.Entries[0].Key != "SECRET_KEY" {
		t.Errorf("expected key SECRET_KEY, got %q", m.Entries[0].Key)
	}
}

func TestParse_DefaultValue(t *testing.T) {
	path := writeTempManifest(t, "APP_ENV=production\n")
	m, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Entries[0].Default != "production" {
		t.Errorf("expected default 'production', got %q", m.Entries[0].Default)
	}
}

func TestParse_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempManifest(t, "# header comment\n\nAPP_ENV # inline\n")
	m, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(m.Entries))
	}
	if m.Entries[0].Comment != "inline" {
		t.Errorf("expected comment 'inline', got %q", m.Entries[0].Comment)
	}
}

func TestManifest_HasKey(t *testing.T) {
	path := writeTempManifest(t, "APP_ENV\nDB_HOST\n")
	m, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.HasKey("APP_ENV") {
		t.Error("expected HasKey(APP_ENV) = true")
	}
	if m.HasKey("MISSING") {
		t.Error("expected HasKey(MISSING) = false")
	}
}

func TestManifest_Keys(t *testing.T) {
	path := writeTempManifest(t, "A\nB\nC\n")
	m, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := m.Keys()
	if len(keys) != 3 || keys[0] != "A" || keys[1] != "B" || keys[2] != "C" {
		t.Errorf("unexpected keys: %v", keys)
	}
}
