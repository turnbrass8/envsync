package sync

import (
	"fmt"

	"github.com/yourorg/envsync/internal/diff"
	"github.com/yourorg/envsync/internal/envfile"
	"github.com/yourorg/envsync/internal/manifest"
)

// Result holds the outcome of a sync operation.
type Result struct {
	Applied []string
	Skipped []string
	Errors  []error
}

// Syncer applies manifest-driven changes to a target .env file.
type Syncer struct {
	DryRun bool
}

// New returns a new Syncer.
func New(dryRun bool) *Syncer {
	return &Syncer{DryRun: dryRun}
}

// Sync compares the manifest against the target env map and fills in missing
// or mismatched keys. It writes the result back to targetPath unless DryRun
// is true.
func (s *Syncer) Sync(man *manifest.Manifest, target envfile.Env, targetPath string) (*Result, error) {
	entries := diff.Compare(man, target)

	result := &Result{}
	updated := envfile.Env{}
	for k, v := range target {
		updated[k] = v
	}

	for _, e := range entries {
		switch e.Status {
		case diff.StatusMissing:
			if e.Default != "" {
				updated[e.Key] = e.Default
				result.Applied = append(result.Applied, fmt.Sprintf("%s (default: %q)", e.Key, e.Default))
			} else {
				result.Skipped = append(result.Skipped, e.Key)
				if e.Required {
					result.Errors = append(result.Errors, fmt.Errorf("required key %q is missing and has no default", e.Key))
				}
			}
		case diff.StatusOK, diff.StatusExtra:
			// nothing to do
		}
	}

	if !s.DryRun && len(result.Applied) > 0 {
		if err := envfile.Write(targetPath, updated); err != nil {
			return result, fmt.Errorf("writing target file: %w", err)
		}
	}

	return result, nil
}
