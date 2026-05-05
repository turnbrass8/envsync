// Package resolve provides environment variable resolution across multiple
// sources with precedence ordering and fallback support.
package resolve

import (
	"fmt"
	"sort"
)

// Source represents a named environment variable source.
type Source struct {
	Name   string
	Values map[string]string
}

// Result holds a resolved value along with its origin source.
type Result struct {
	Key    string
	Value  string
	Source string
	Found  bool
}

// Options controls resolution behaviour.
type Options struct {
	// Strict causes Resolve to return an error if any key is unresolved.
	Strict bool
	// Fallback is returned when a key is not found and Strict is false.
	Fallback string
}

// Resolve looks up each key across sources in order, returning the first match.
// Sources are checked in the order provided (index 0 = highest precedence).
func Resolve(keys []string, sources []Source, opts Options) ([]Result, error) {
	results := make([]Result, 0, len(keys))

	for _, key := range keys {
		r := Result{Key: key}
		for _, src := range sources {
			if v, ok := src.Values[key]; ok {
				r.Value = v
				r.Source = src.Name
				r.Found = true
				break
			}
		}
		if !r.Found {
			if opts.Strict {
				return nil, fmt.Errorf("resolve: key %q not found in any source", key)
			}
			r.Value = opts.Fallback
			r.Source = ""
		}
		results = append(results, r)
	}
	return results, nil
}

// ResolveAll resolves every key present across all sources, merging them with
// source precedence (index 0 wins on conflict).
func ResolveAll(sources []Source, opts Options) []Result {
	seen := make(map[string]Result)
	order := []string{}

	// Iterate in reverse so higher-precedence sources overwrite.
	for i := len(sources) - 1; i >= 0; i-- {
		src := sources[i]
		for k, v := range src.Values {
			if _, exists := seen[k]; !exists {
				order = append(order, k)
			}
			seen[k] = Result{Key: k, Value: v, Source: src.Name, Found: true}
		}
	}

	sort.Strings(order)
	out := make([]Result, 0, len(order))
	for _, k := range order {
		out = append(out, seen[k])
	}
	return out
}
