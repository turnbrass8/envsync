// Package patch applies a set of key=value edits to an existing .env file,
// preserving comments, blank lines, and key ordering.
package patch

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Op represents a single patch operation.
type Op struct {
	Key   string
	Value string
	Delete bool
}

// Result summarises what Patch did.
type Result struct {
	Updated []string
	Added   []string
	Deleted []string
}

// Summary returns a one-line description of the result.
func (r Result) Summary() string {
	return fmt.Sprintf("%d updated, %d added, %d deleted",
		len(r.Updated), len(r.Added), len(r.Deleted))
}

// Patch applies ops to the file at path. When dryRun is true the file is not
// written but the Result is still returned.
func Patch(path string, ops []Op, dryRun bool) (Result, error) {
	if len(ops) == 0 {
		return Result{}, errors.New("patch: no operations provided")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return Result{}, fmt.Errorf("patch: read %s: %w", path, err)
	}

	index := make(map[string]Op, len(ops))
	for _, op := range ops {
		if op.Key == "" {
			return Result{}, errors.New("patch: op with empty key")
		}
		index[op.Key] = op
	}

	lines := strings.Split(string(raw), "\n")
	out := make([]string, 0, len(lines))
	var res Result
	seen := make(map[string]bool)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			out = append(out, line)
			continue
		}
		idx := strings.IndexByte(trimmed, '=')
		if idx < 0 {
			out = append(out, line)
			continue
		}
		key := strings.TrimSpace(trimmed[:idx])
		op, ok := index[key]
		if !ok {
			out = append(out, line)
			continue
		}
		seen[key] = true
		if op.Delete {
			res.Deleted = append(res.Deleted, key)
			// omit line
		} else {
			out = append(out, key+"="+op.Value)
			res.Updated = append(res.Updated, key)
		}
	}

	// Append ops whose keys were not found (not deletes).
	for _, op := range ops {
		if !seen[op.Key] && !op.Delete {
			out = append(out, op.Key+"="+op.Value)
			res.Added = append(res.Added, op.Key)
		}
	}

	if dryRun {
		return res, nil
	}

	content := strings.Join(out, "\n")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return res, fmt.Errorf("patch: write %s: %w", path, err)
	}
	return res, nil
}
