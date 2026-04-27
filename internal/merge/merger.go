// Package merge provides functionality to merge multiple .env files
// with configurable precedence and conflict resolution strategies.
package merge

import (
	"fmt"
	"sort"

	"github.com/user/envsync/internal/envfile"
)

// Strategy defines how conflicts are resolved when the same key
// appears in multiple source files.
type Strategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last file that defines the key.
	StrategyLast
	// StrategyError returns an error when a key conflict is detected.
	StrategyError
)

// Result holds the merged environment map and metadata about the merge.
type Result struct {
	Env      map[string]string
	Sources  map[string]string // key -> source file path
	Conflicts []Conflict
}

// Conflict describes a key that appeared in more than one source file.
type Conflict struct {
	Key    string
	Files  []string
	Chosen string
}

// Merge combines the given env files using the provided strategy.
// Files are processed in order; earlier files have lower index.
func Merge(paths []string, strategy Strategy) (*Result, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("merge: no source files provided")
	}

	result := &Result{
		Env:     make(map[string]string),
		Sources: make(map[string]string),
	}

	// Track which files each key was seen in for conflict detection.
	seen := make(map[string][]string)

	for _, path := range paths {
		env, err := envfile.Parse(path)
		if err != nil {
			return nil, fmt.Errorf("merge: reading %s: %w", path, err)
		}

		keys := make([]string, 0, len(env))
		for k := range env {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := env[k]
			seen[k] = append(seen[k], path)

			_, exists := result.Env[k]
			if exists {
				switch strategy {
				case StrategyError:
					return nil, fmt.Errorf("merge: conflict on key %q in %s and %s", k, result.Sources[k], path)
				case StrategyFirst:
					// keep existing value, record conflict
				case StrategyLast:
					result.Env[k] = v
					result.Sources[k] = path
				}
			} else {
				result.Env[k] = v
				result.Sources[k] = path
			}
		}
	}

	for k, files := range seen {
		if len(files) > 1 {
			result.Conflicts = append(result.Conflicts, Conflict{
				Key:    k,
				Files:  files,
				Chosen: result.Sources[k],
			})
		}
	}
	sort.Slice(result.Conflicts, func(i, j int) bool {
		return result.Conflicts[i].Key < result.Conflicts[j].Key
	})

	return result, nil
}
