package encrypt

import (
	"os"
	"path/filepath"
	"testing"
)

func tempKeystorePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "keystore.json")
}

func TestLoadKeystore_NonExistent_ReturnsEmpty(t *testing.T) {
	ks, err := LoadKeystore(tempKeystorePath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ks.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(ks.Entries))
	}
}

func TestAdd_And_Save_And_Reload(t *testing.T) {
	path := tempKeystorePath(t)
	ks, _ := LoadKeystore(path)

	if err := ks.Add("prod", "production key"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if err := ks.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	ks2, err := LoadKeystore(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if len(ks2.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(ks2.Entries))
	}
	if ks2.Entries[0].Alias != "prod" {
		t.Errorf("expected alias 'prod', got %q", ks2.Entries[0].Alias)
	}
	if ks2.Entries[0].Hint != "production key" {
		t.Errorf("unexpected hint: %q", ks2.Entries[0].Hint)
	}
}

func TestAdd_DuplicateAlias_ReturnsError(t *testing.T) {
	ks, _ := LoadKeystore(tempKeystorePath(t))
	_ = ks.Add("staging", "")
	if err := ks.Add("staging", "dup"); err == nil {
		t.Error("expected error for duplicate alias")
	}
}

func TestRemove_ExistingAlias(t *testing.T) {
	ks, _ := LoadKeystore(tempKeystorePath(t))
	_ = ks.Add("dev", "")
	if err := ks.Remove("dev"); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
	if len(ks.Entries) != 0 {
		t.Errorf("expected 0 entries after remove")
	}
}

func TestRemove_MissingAlias_ReturnsError(t *testing.T) {
	ks, _ := LoadKeystore(tempKeystorePath(t))
	if err := ks.Remove("ghost"); err == nil {
		t.Error("expected error removing missing alias")
	}
}

func TestSave_CreatesParentDirs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "dir", "ks.json")
	ks, _ := LoadKeystore(path)
	_ = ks.Add("ci", "ci hint")
	if err := ks.Save(); err != nil {
		t.Fatalf("Save with nested path failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
