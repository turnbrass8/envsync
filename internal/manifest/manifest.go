package manifest

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key entry in the manifest file.
type Entry struct {
	Key      string
	Required bool
	Default  string
	Comment  string
}

// Manifest holds all declared keys for an environment.
type Manifest struct {
	Entries []Entry
}

// Keys returns a slice of all key names in the manifest.
func (m *Manifest) Keys() []string {
	keys := make([]string, 0, len(m.Entries))
	for _, e := range m.Entries {
		keys = append(keys, e.Key)
	}
	return keys
}

// HasKey reports whether the manifest declares the given key.
func (m *Manifest) HasKey(key string) bool {
	for _, e := range m.Entries {
		if e.Key == key {
			return true
		}
	}
	return false
}

// Parse reads a manifest file from path.
// Manifest format per line:
//   KEY              # optional comment
//   KEY=default      # key with default value
//   KEY!             # required key (no default)
func Parse(path string) (*Manifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("manifest: open %q: %w", path, err)
	}
	defer f.Close()

	var m Manifest
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Strip inline comment
		comment := ""
		if idx := strings.Index(line, "#"); idx >= 0 {
			comment = strings.TrimSpace(line[idx+1:])
			line = line[:idx]
		}
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		entry := Entry{Comment: comment}

		if strings.HasSuffix(line, "!") {
			entry.Key = strings.TrimSuffix(line, "!")
			entry.Required = true
		} else if idx := strings.Index(line, "="); idx >= 0 {
			entry.Key = strings.TrimSpace(line[:idx])
			entry.Default = strings.TrimSpace(line[idx+1:])
		} else {
			entry.Key = line
		}

		if entry.Key == "" {
			return nil, fmt.Errorf("manifest: line %d: empty key", lineNum)
		}

		m.Entries = append(m.Entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("manifest: scan %q: %w", path, err)
	}

	return &m, nil
}
