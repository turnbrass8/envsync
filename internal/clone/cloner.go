// Package clone copies a subset of keys from one env file to another,
// optionally renaming them via a mapping.
package clone

import (
	"fmt"

	"github.com/user/envsync/internal/envfile"
)

// Rule describes a single key to clone and an optional rename target.
type Rule struct {
	SrcKey string
	DstKey string // if empty, same as SrcKey
}

// Options controls Clone behaviour.
type Options struct {
	Rules     []Rule
	Overwrite bool
	DryRun    bool
}

// Result summarises what Clone did.
type Result struct {
	Copied  []string
	Skipped []string
}

// Clone reads keys defined in opts.Rules from src and writes them into dst.
// If a key already exists in dst and Overwrite is false it is skipped.
// When DryRun is true the destination file is never written.
func Clone(srcPath, dstPath string, opts Options) (Result, error) {
	srcEnv, err := envfile.Parse(srcPath)
	if err != nil {
		return Result{}, fmt.Errorf("clone: parse src: %w", err)
	}

	dstEnv, err := envfile.Parse(dstPath)
	if err != nil {
		// destination may not exist yet — start with an empty map
		dstEnv = &envfile.Env{Values: make(map[string]string)}
	}

	var res Result

	for _, rule := range opts.Rules {
		val, ok := srcEnv.Values[rule.SrcKey]
		if !ok {
			return Result{}, fmt.Errorf("clone: key %q not found in source", rule.SrcKey)
		}

		dst := rule.DstKey
		if dst == "" {
			dst = rule.SrcKey
		}

		if _, exists := dstEnv.Values[dst]; exists && !opts.Overwrite {
			res.Skipped = append(res.Skipped, dst)
			continue
		}

		dstEnv.Values[dst] = val
		res.Copied = append(res.Copied, dst)
	}

	if !opts.DryRun {
		if err := envfile.Write(dstPath, dstEnv); err != nil {
			return Result{}, fmt.Errorf("clone: write dst: %w", err)
		}
	}

	return res, nil
}
