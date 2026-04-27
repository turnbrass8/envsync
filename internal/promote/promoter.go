// Package promote handles copying env values from one environment file to another
// based on a set of keys defined in a manifest.
package promote

import (
	"fmt"

	"github.com/yourorg/envsync/internal/envfile"
	"github.com/yourorg/envsync/internal/manifest"
)

// Result describes the outcome of a single key promotion.
type Result struct {
	Key       string
	Promoted  bool
	Skipped   bool
	Reason    string
}

// Options controls promotion behaviour.
type Options struct {
	DryRun    bool
	Overwrite bool
}

// Promote copies keys defined in mf from src env to dst env file.
// Keys already present in dst are skipped unless Overwrite is set.
func Promote(srcPath, dstPath string, mf *manifest.Manifest, opts Options) ([]Result, error) {
	srcEnv, err := envfile.Parse(srcPath)
	if err != nil {
		return nil, fmt.Errorf("promote: parse source %q: %w", srcPath, err)
	}

	dstEnv, err := envfile.Parse(dstPath)
	if err != nil {
		return nil, fmt.Errorf("promote: parse destination %q: %w", dstPath, err)
	}

	var results []Result

	for _, entry := range mf.Entries {
		val, srcOK := srcEnv.Get(entry.Key)
		if !srcOK {
			results = append(results, Result{
				Key:    entry.Key,
				Skipped: true,
				Reason: "not present in source",
			})
			continue
		}

		_, dstOK := dstEnv.Get(entry.Key)
		if dstOK && !opts.Overwrite {
			results = append(results, Result{
				Key:    entry.Key,
				Skipped: true,
				Reason: "already present in destination (use --overwrite to replace)",
			})
			continue
		}

		if !opts.DryRun {
			dstEnv.Set(entry.Key, val)
		}

		results = append(results, Result{
			Key:      entry.Key,
			Promoted: true,
		})
	}

	if !opts.DryRun {
		if err := envfile.Write(dstPath, dstEnv); err != nil {
			return nil, fmt.Errorf("promote: write destination %q: %w", dstPath, err)
		}
	}

	return results, nil
}
