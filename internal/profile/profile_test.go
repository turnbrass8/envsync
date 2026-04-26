package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envsync/internal/profile"
)

func makeProfileDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("setup: write %s: %v", name, err)
		}
	}
	return dir
}

func TestList_ReturnsProfiles(t *testing.T) {
	dir := makeProfileDir(t, map[string]string{
		"dev.env":     "APP_ENV=dev\n",
		"staging.env": "APP_ENV=staging\n",
		"README.md":   "# not a profile\n",
	})

	m := profile.NewManager(dir)
	profiles, err := m.List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(profiles))
	}
	names := map[string]bool{}
	for _, p := range profiles {
		names[p.Name] = true
	}
	if !names["dev"] || !names["staging"] {
		t.Errorf("unexpected profile names: %v", names)
	}
}

func TestList_EmptyDir_ReturnsNil(t *testing.T) {
	dir := makeProfileDir(t, map[string]string{})
	m := profile.NewManager(dir)
	profiles, err := m.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("expected empty list, got %v", profiles)
	}
}

func TestList_NonExistentDir_ReturnsNil(t *testing.T) {
	m := profile.NewManager("/tmp/does-not-exist-envsync-profile-test")
	profiles, err := m.List()
	if err != nil {
		t.Fatalf("expected nil error for missing dir, got: %v", err)
	}
	if profiles != nil {
		t.Errorf("expected nil profiles, got %v", profiles)
	}
}

func TestResolve_ExistingProfile(t *testing.T) {
	dir := makeProfileDir(t, map[string]string{"prod.env": "APP_ENV=prod\n"})
	m := profile.NewManager(dir)
	path, err := m.Resolve("prod")
	if err != nil {
		t.Fatalf("Resolve() error: %v", err)
	}
	if path == "" {
		t.Error("expected non-empty path")
	}
}

func TestResolve_MissingProfile_ReturnsError(t *testing.T) {
	dir := makeProfileDir(t, map[string]string{})
	m := profile.NewManager(dir)
	_, err := m.Resolve("ghost")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestExists(t *testing.T) {
	dir := makeProfileDir(t, map[string]string{"local.env": "DEBUG=true\n"})
	m := profile.NewManager(dir)
	if !m.Exists("local") {
		t.Error("expected local profile to exist")
	}
	if m.Exists("nope") {
		t.Error("expected nope profile to not exist")
	}
}
