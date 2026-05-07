package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/group"
)

func runGroup(args []string) error {
	fs := flag.NewFlagSet("group", flag.ContinueOnError)
	mappings := fs.String("map", "", "prefix:name pairs, comma-separated (e.g. DB:database,REDIS:cache)")
	strip := fs.Bool("strip-prefix", false, "strip matched prefix from output keys")
	includeOther := fs.Bool("include-other", false, "collect unmatched keys into _other group")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("group: usage: envsync group [flags] <envfile>")
	}
	if *mappings == "" {
		return fmt.Errorf("group: --map flag is required")
	}

	prefixes, err := parsePrefixMappings(*mappings)
	if err != nil {
		return err
	}

	env, err := envfile.Parse(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("group: %w", err)
	}

	groups, err := group.GroupBy(env, prefixes, group.Options{
		StripPrefix:      *strip,
		IncludeUnmatched: *includeOther,
	})
	if err != nil {
		return fmt.Errorf("group: %w", err)
	}

	w := os.Stdout
	for _, g := range groups {
		fmt.Fprintf(w, "[%s]\n", g.Name)
		keys := make([]string, 0, len(g.Entries))
		for k := range g.Entries {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "  %s=%s\n", k, g.Entries[k])
		}
		fmt.Fprintln(w)
	}
	return nil
}

func parsePrefixMappings(raw string) (map[string]string, error) {
	result := make(map[string]string)
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		idx := strings.IndexByte(part, ':')
		if idx < 1 {
			return nil, fmt.Errorf("group: invalid mapping %q, expected prefix:name", part)
		}
		result[part[:idx]] = part[idx+1:]
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("group: no valid prefix mappings found in %q", raw)
	}
	return result, nil
}
