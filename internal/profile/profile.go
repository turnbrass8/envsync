// Package profile provides support for named environment profiles,
// allowing users to switch between sets of .env overrides (e.g. dev, staging, prod).
package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Profile represents a named environment configuration profile.
type Profile struct {
	Name string
	Path string
}

// Manager manages multiple named profiles rooted at a base directory.
type Manager struct {
	BaseDir string
}

// NewManager creates a Manager that stores profiles under baseDir.
func NewManager(baseDir string) *Manager {
	return &Manager{BaseDir: baseDir}
}

// List returns all profiles found in the base directory.
// Profiles are .env files named as "<name>.env".
func (m *Manager) List() ([]Profile, error) {
	entries, err := os.ReadDir(m.BaseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("profile: reading directory %q: %w", m.BaseDir, err)
	}

	var profiles []Profile
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".env") {
			profileName := strings.TrimSuffix(name, ".env")
			profiles = append(profiles, Profile{
				Name: profileName,
				Path: filepath.Join(m.BaseDir, name),
			})
		}
	}
	return profiles, nil
}

// Resolve returns the filesystem path for a named profile.
func (m *Manager) Resolve(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("profile: name must not be empty")
	}
	path := filepath.Join(m.BaseDir, name+".env")
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("profile %q not found (expected %s)", name, path)
		}
		return "", fmt.Errorf("profile: stat %q: %w", path, err)
	}
	return path, nil
}

// Exists reports whether a profile with the given name exists.
func (m *Manager) Exists(name string) bool {
	_, err := m.Resolve(name)
	return err == nil
}
