package encrypt

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// KeyEntry represents a named passphrase alias stored in the keystore.
type KeyEntry struct {
	Alias     string    `json:"alias"`
	Hint      string    `json:"hint,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Keystore holds a collection of key aliases and metadata (not the raw passphrases).
type Keystore struct {
	Entries []KeyEntry `json:"entries"`
	path    string
}

// LoadKeystore reads a keystore file from disk, or returns an empty one if the file does not exist.
func LoadKeystore(path string) (*Keystore, error) {
	ks := &Keystore{path: path}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return ks, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, ks); err != nil {
		return nil, err
	}
	return ks, nil
}

// Add inserts a new alias entry. Returns an error if the alias already exists.
func (ks *Keystore) Add(alias, hint string) error {
	for _, e := range ks.Entries {
		if e.Alias == alias {
			return errors.New("alias already exists: " + alias)
		}
	}
	ks.Entries = append(ks.Entries, KeyEntry{
		Alias:     alias,
		Hint:      hint,
		CreatedAt: time.Now().UTC(),
	})
	return nil
}

// Remove deletes an alias entry by name. Returns an error if not found.
func (ks *Keystore) Remove(alias string) error {
	for i, e := range ks.Entries {
		if e.Alias == alias {
			ks.Entries = append(ks.Entries[:i], ks.Entries[i+1:]...)
			return nil
		}
	}
	return errors.New("alias not found: " + alias)
}

// Save writes the keystore to disk, creating parent directories as needed.
func (ks *Keystore) Save() error {
	if err := os.MkdirAll(filepath.Dir(ks.path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(ks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ks.path, data, 0o600)
}
