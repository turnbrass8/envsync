// Package env provides utilities for loading, resolving, and inspecting
// environment variables from multiple sources (OS env, .env files, defaults).
package env

import (
	"os"
	"sort"
	"strings"

	"github.com/user/envsync/internal/envfile"
)

// Source represents a named origin of environment variables.
type Source struct {
	Name   string
	Values map[string]string
}

// Loader merges environment variables from multiple sources in priority order.
// Later sources override earlier ones.
type Loader struct {
	sources []Source
}

// NewLoader creates a Loader with the given sources (lowest to highest priority).
func NewLoader(sources ...Source) *Loader {
	return &Loader{sources: sources}
}

// OSSource returns a Source populated from the current process environment.
func OSSource() Source {
	env := os.Environ()
	m := make(map[string]string, len(env))
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return Source{Name: "os", Values: m}
}

// FileSource returns a Source populated by parsing a .env file.
func FileSource(name, path string) (Source, error) {
	env, err := envfile.Parse(path)
	if err != nil {
		return Source{}, err
	}
	return Source{Name: name, Values: env.All()}, nil
}

// Resolve merges all sources and returns the combined map.
func (l *Loader) Resolve() map[string]string {
	result := make(map[string]string)
	for _, src := range l.sources {
		for k, v := range src.Values {
			result[k] = v
		}
	}
	return result
}

// Keys returns all keys across all sources, sorted and deduplicated.
func (l *Loader) Keys() []string {
	seen := make(map[string]struct{})
	for _, src := range l.sources {
		for k := range src.Values {
			seen[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Origin returns the name of the source that last defined a key.
func (l *Loader) Origin(key string) string {
	origin := ""
	for _, src := range l.sources {
		if _, ok := src.Values[key]; ok {
			origin = src.Name
		}
	}
	return origin
}
