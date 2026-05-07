// Package group provides functionality to group env vars by prefix into named sections.
package group

import (
	"fmt"
	"sort"
	"strings"
)

// Group represents a named collection of env entries sharing a common prefix.
type Group struct {
	Name    string
	Prefix  string
	Entries map[string]string
}

// Options controls grouping behaviour.
type Options struct {
	// StripPrefix removes the prefix from keys in the resulting group.
	StripPrefix bool
	// IncludeUnmatched collects keys that don't match any prefix into a group named "_other".
	IncludeUnmatched bool
}

// GroupBy splits env into groups based on the provided prefix→name mapping.
// Prefixes are matched case-insensitively and must end without an underscore;
// the separator underscore is assumed (e.g. prefix "DB" matches "DB_HOST").
func GroupBy(env map[string]string, prefixes map[string]string, opts Options) ([]*Group, error) {
	if len(prefixes) == 0 {
		return nil, fmt.Errorf("group: at least one prefix mapping is required")
	}

	// Build ordered list of prefixes for deterministic output.
	type prefixEntry struct{ prefix, name string }
	var ordered []prefixEntry
	for p, n := range prefixes {
		ordered = append(ordered, prefixEntry{strings.ToUpper(p), n})
	}
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].name < ordered[j].name })

	groupMap := make(map[string]*Group, len(ordered))
	for _, pe := range ordered {
		groupMap[pe.prefix] = &Group{
			Name:    pe.name,
			Prefix:  pe.prefix,
			Entries: make(map[string]string),
		}
	}

	var other *Group
	if opts.IncludeUnmatched {
		other = &Group{Name: "_other", Entries: make(map[string]string)}
	}

	for k, v := range env {
		upper := strings.ToUpper(k)
		matched := false
		for _, pe := range ordered {
			if strings.HasPrefix(upper, pe.prefix+"_") || upper == pe.prefix {
				key := k
				if opts.StripPrefix && len(k) > len(pe.prefix) {
					key = k[len(pe.prefix)+1:]
				}
				groupMap[pe.prefix].Entries[key] = v
				matched = true
				break
			}
		}
		if !matched && other != nil {
			other.Entries[k] = v
		}
	}

	var result []*Group
	for _, pe := range ordered {
		result = append(result, groupMap[pe.prefix])
	}
	if other != nil && len(other.Entries) > 0 {
		result = append(result, other)
	}
	return result, nil
}
