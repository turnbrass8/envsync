// Package snapshot provides functionality to capture and compare
// point-in-time snapshots of .env files for drift detection.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envsync/internal/envfile"
)

// Snapshot represents a captured state of an env file at a point in time.
type Snapshot struct {
	CapturedAt time.Time         `json:"captured_at"`
	Source     string            `json:"source"`
	Values     map[string]string `json:"values"`
}

// Diff represents the difference between two snapshots.
type Diff struct {
	Added   map[string]string `json:"added"`
	Removed map[string]string `json:"removed"`
	Changed map[string][2]string `json:"changed"`
}

// HasChanges returns true if there are any differences.
func (d *Diff) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Capture reads an env file and returns a Snapshot.
func Capture(path string) (*Snapshot, error) {
	env, err := envfile.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: parse %q: %w", path, err)
	}
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		Source:     path,
		Values:     env,
	}, nil
}

// Save writes a snapshot to a JSON file.
func Save(snap *Snapshot, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("snapshot: create %q: %w", dest, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

// Load reads a snapshot from a JSON file.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open %q: %w", path, err)
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("snapshot: decode %q: %w", path, err)
	}
	return &snap, nil
}

// Compare returns a Diff between an older and newer snapshot.
func Compare(old, new *Snapshot) *Diff {
	d := &Diff{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}
	for k, v := range new.Values {
		if oldVal, ok := old.Values[k]; !ok {
			d.Added[k] = v
		} else if oldVal != v {
			d.Changed[k] = [2]string{oldVal, v}
		}
	}
	for k, v := range old.Values {
		if _, ok := new.Values[k]; !ok {
			d.Removed[k] = v
		}
	}
	return d
}
