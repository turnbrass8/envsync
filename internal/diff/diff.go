package diff

import (
	"fmt"
	"sort"
)

// Status represents the diff status of a key.
type Status int

const (
	StatusMissing  Status = iota // key exists in manifest but not in env file
	StatusExtra                  // key exists in env file but not in manifest
	StatusPresent                // key exists in both
)

// Entry represents a single diff result for a key.
type Entry struct {
	Key    string
	Status Status
	Value  string // populated when StatusPresent or StatusExtra
}

// String returns a human-readable representation of the diff entry.
func (e Entry) String() string {
	switch e.Status {
	case StatusMissing:
		return fmt.Sprintf("- %s (missing)", e.Key)
	case StatusExtra:
		return fmt.Sprintf("+ %s (extra)", e.Key)
	case StatusPresent:
		return fmt.Sprintf("  %s", e.Key)
	}
	return e.Key
}

// Result holds the full diff between a manifest key set and an env map.
type Result struct {
	Entries []Entry
}

// HasMissing returns true if any keys are missing from the env file.
func (r Result) HasMissing() bool {
	for _, e := range r.Entries {
		if e.Status == StatusMissing {
			return true
		}
	}
	return false
}

// Missing returns all entries with StatusMissing.
func (r Result) Missing() []Entry {
	var out []Entry
	for _, e := range r.Entries {
		if e.Status == StatusMissing {
			out = append(out, e)
		}
	}
	return out
}

// Extra returns all entries with StatusExtra.
func (r Result) Extra() []Entry {
	var out []Entry
	for _, e := range r.Entries {
		if e.Status == StatusExtra {
			out = append(out, e)
		}
	}
	return out
}

// Compare diffs manifestKeys (ordered list of expected keys) against envMap
// (the parsed env file). It returns a Result describing missing, extra, and
// present keys.
func Compare(manifestKeys []string, envMap map[string]string) Result {
	expected := make(map[string]struct{}, len(manifestKeys))
	for _, k := range manifestKeys {
		expected[k] = struct{}{}
	}

	var entries []Entry

	// Check each manifest key against the env map.
	for _, k := range manifestKeys {
		if v, ok := envMap[k]; ok {
			entries = append(entries, Entry{Key: k, Status: StatusPresent, Value: v})
		} else {
			entries = append(entries, Entry{Key: k, Status: StatusMissing})
		}
	}

	// Collect extra keys from the env map not in the manifest.
	extraKeys := make([]string, 0)
	for k, v := range envMap {
		if _, ok := expected[k]; !ok {
			_ = v
			extraKeys = append(extraKeys, k)
		}
	}
	sort.Strings(extraKeys)
	for _, k := range extraKeys {
		entries = append(entries, Entry{Key: k, Status: StatusExtra, Value: envMap[k]})
	}

	return Result{Entries: entries}
}
