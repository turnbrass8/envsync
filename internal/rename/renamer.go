// Package rename provides utilities for renaming keys in .env files.
package rename

import (
	"fmt"
	"os"

	"github.com/user/envsync/internal/envfile"
)

// Rule describes a single key rename operation.
type Rule struct {
	From string
	To   string
}

// Result holds the outcome of a rename operation.
type Result struct {
	Applied []Rule
	Skipped []Rule // From key not present in env
}

// Rename applies the given rules to the env file at path, writing the result
// back to the same file unless dryRun is true.
func Rename(path string, rules []Rule, dryRun bool) (Result, error) {
	env, err := envfile.Parse(path)
	if err != nil {
		return Result{}, fmt.Errorf("rename: parse %q: %w", path, err)
	}

	var result Result

	for _, rule := range rules {
		val, ok := env.Get(rule.From)
		if !ok {
			result.Skipped = append(result.Skipped, rule)
			continue
		}
		if _, exists := env.Get(rule.To); exists {
			return result, fmt.Errorf("rename: target key %q already exists", rule.To)
		}
		env.Set(rule.To, val)
		env.Delete(rule.From)
		result.Applied = append(result.Applied, rule)
	}

	if dryRun {
		return result, nil
	}

	f, err := os.Create(path)
	if err != nil {
		return result, fmt.Errorf("rename: open %q for writing: %w", path, err)
	}
	defer f.Close()

	if err := envfile.Write(f, env); err != nil {
		return result, fmt.Errorf("rename: write %q: %w", path, err)
	}

	return result, nil
}
