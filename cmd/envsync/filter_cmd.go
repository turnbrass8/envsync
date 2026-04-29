package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/yourorg/envsync/internal/envfile"
	"github.com/yourorg/envsync/internal/filter"
)

func runFilter(args []string) error {
	fs := flag.NewFlagSet("filter", flag.ContinueOnError)
	envPath := fs.String("env", ".env", "path to .env file")
	prefix := fs.String("prefix", "", "only include keys with this prefix")
	pattern := fs.String("pattern", "", "only include keys matching this regex")
	exclude := fs.String("exclude", "", "comma-separated keys to exclude")
	keys := fs.String("keys", "", "comma-separated explicit keys to include")
	stripPfx := fs.Bool("strip-prefix", false, "remove --prefix from output key names")
	format := fs.String("format", "dotenv", "output format: dotenv|json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	env, err := envfile.Parse(*envPath)
	if err != nil {
		return fmt.Errorf("filter: %w", err)
	}

	var excludeList, keyList []string
	if *exclude != "" {
		excludeList = strings.Split(*exclude, ",")
	}
	if *keys != "" {
		keyList = strings.Split(*keys, ",")
	}

	opts := filter.Options{
		Prefix:  *prefix,
		Pattern: *pattern,
		Exclude: excludeList,
		Keys:    keyList,
	}

	result, err := filter.Filter(env, opts)
	if err != nil {
		return err
	}

	if *stripPfx && *prefix != "" {
		result = filter.StripPrefix(result, *prefix)
	}

	switch *format {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(result)
	default:
		var sorted []string
		for k := range result {
			sorted = append(sorted, k)
		}
		sort.Strings(sorted)
		for _, k := range sorted {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, result[k])
		}
	}
	return nil
}
