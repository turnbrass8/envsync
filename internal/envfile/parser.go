package envfile

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry represents a single key-value pair in an .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
}

// EnvFile holds all entries parsed from an .env file.
type EnvFile struct {
	Entries []Entry
	Index   map[string]int // key -> index in Entries
}

// Parse reads an .env file from the given reader and returns an EnvFile.
func Parse(r io.Reader) (*EnvFile, error) {
	ef := &EnvFile{
		Index: make(map[string]int),
	}

	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format %q", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNum)
		}

		// Strip inline comments
		comment := ""
		if idx := strings.Index(val, " #"); idx != -1 {
			comment = strings.TrimSpace(val[idx+1:])
			val = strings.TrimSpace(val[:idx])
		}

		// Strip surrounding quotes
		val = stripQuotes(val)

		entry := Entry{Key: key, Value: val, Comment: comment}
		ef.Index[key] = len(ef.Entries)
		ef.Entries = append(ef.Entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning error: %w", err)
	}

	return ef, nil
}

// Get returns the value for a given key, and whether it exists.
func (ef *EnvFile) Get(key string) (string, bool) {
	if idx, ok := ef.Index[key]; ok {
		return ef.Entries[idx].Value, true
	}
	return "", false
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
