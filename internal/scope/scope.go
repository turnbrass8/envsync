// Package scope provides utilities for scoping .env keys by environment
// prefix (e.g. "prod", "staging"), allowing multi-environment values to
// coexist in a single file.
package scope

import (
	"fmt"
	"strings"
)

// Entry holds a key/value pair together with its resolved scope.
type Entry struct {
	Key   string
	Value string
	Scope string
}

// Extract separates entries that belong to the given scope from a flat
// map of key/value pairs.  Keys are expected to follow the convention
//
//	<SCOPE>__<KEY>   (two underscores as separator)
//
// Keys that do not carry any scope prefix are included when
// includeGlobal is true.
func Extract(env map[string]string, scope string, includeGlobal bool) []Entry {
	prefix := strings.ToUpper(scope) + "__"
	var out []Entry
	for k, v := range env {
		if strings.HasPrefix(k, prefix) {
			bare := strings.TrimPrefix(k, prefix)
			out = append(out, Entry{Key: bare, Value: v, Scope: scope})
			continue
		}
		if includeGlobal && !strings.Contains(k, "__") {
			out = append(out, Entry{Key: k, Value: v, Scope: ""})
		}
	}
	return out
}

// Flatten converts a slice of Entries back to a plain map.  Scoped
// entries take precedence over global ones when keys collide.
func Flatten(entries []Entry) map[string]string {
	global := make(map[string]string)
	scoped := make(map[string]string)
	for _, e := range entries {
		if e.Scope == "" {
			global[e.Key] = e.Value
		} else {
			scoped[e.Key] = e.Value
		}
	}
	out := make(map[string]string, len(global)+len(scoped))
	for k, v := range global {
		out[k] = v
	}
	for k, v := range scoped {
		out[k] = v
	}
	return out
}

// Prefix returns the scoped form of a bare key.
func Prefix(scope, key string) string {
	return fmt.Sprintf("%s__%s", strings.ToUpper(scope), key)
}
