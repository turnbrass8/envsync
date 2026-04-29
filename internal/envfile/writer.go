package envfile

import (
	"fmt"
	"io"
	"strings"
)

// Write serialises an EnvFile back to .env format.
func Write(w io.Writer, ef *EnvFile) error {
	for _, entry := range ef.Entries {
		val := quoteIfNeeded(entry.Value)
		line := fmt.Sprintf("%s=%s", entry.Key, val)
		if entry.Comment != "" {
			line += " " + entry.Comment
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return fmt.Errorf("writing entry %q: %w", entry.Key, err)
		}
	}
	return nil
}

// Merge merges src entries into dst, overwriting existing keys.
// Returns the list of keys that were added or updated.
func Merge(dst, src *EnvFile) (added, updated []string) {
	for _, entry := range src.Entries {
		if idx, exists := dst.Index[entry.Key]; exists {
			if dst.Entries[idx].Value != entry.Value {
				dst.Entries[idx].Value = entry.Value
				updated = append(updated, entry.Key)
			}
		} else {
			dst.Index[entry.Key] = len(dst.Entries)
			dst.Entries = append(dst.Entries, entry)
			added = append(added, entry.Key)
		}
	}
	return added, updated
}

// Delete removes a key from the EnvFile if it exists.
// Returns true if the key was found and removed, false otherwise.
func Delete(ef *EnvFile, key string) bool {
	idx, exists := ef.Index[key]
	if !exists {
		return false
	}
	ef.Entries = append(ef.Entries[:idx], ef.Entries[idx+1:]...)
	delete(ef.Index, key)
	// Shift all indexes that come after the removed entry.
	for k, i := range ef.Index {
		if i > idx {
			ef.Index[k] = i - 1
		}
	}
	return true
}

// quoteIfNeeded wraps value in double quotes if it contains spaces or special chars.
func quoteIfNeeded(val string) string {
	if strings.ContainsAny(val, " \t#") {
		return `"` + val + `"`
	}
	return val
}
