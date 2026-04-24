package main

import (
	"fmt"
	"os"

	"github.com/yourorg/envsync/internal/audit"
	"github.com/yourorg/envsync/internal/diff"
	"github.com/yourorg/envsync/internal/envfile"
	"github.com/yourorg/envsync/internal/manifest"
)

// runAudit performs a diff between the manifest and the env file,
// records each key's status via the audit logger, and prints a summary.
func runAudit(envPath, manifestPath string) error {
	env, err := envfile.Parse(envPath)
	if err != nil {
		return fmt.Errorf("parsing env file: %w", err)
	}

	mf, err := manifest.Parse(manifestPath)
	if err != nil {
		return fmt.Errorf("parsing manifest: %w", err)
	}

	entries := diff.Compare(mf, env)
	log := audit.New(os.Stdout)

	for _, e := range entries {
		switch e.Status {
		case diff.StatusPresent:
			log.Record(audit.EventSkipped, e.Key, "key present, no action needed")
		case diff.StatusMissing:
			if e.Default != "" {
				log.Record(audit.EventApplied, e.Key, fmt.Sprintf("default=%q", e.Default))
			} else {
				log.Record(audit.EventMissing, e.Key, "no value and no default")
			}
		case diff.StatusExtra:
			log.Record(audit.EventSkipped, e.Key, "extra key not in manifest")
		}
	}

	summary := log.Summary()
	fmt.Fprintf(os.Stdout, "\nAudit summary: applied=%d skipped=%d missing=%d invalid=%d\n",
		summary[audit.EventApplied],
		summary[audit.EventSkipped],
		summary[audit.EventMissing],
		summary[audit.EventInvalid],
	)

	if summary[audit.EventMissing] > 0 {
		return fmt.Errorf("%d required key(s) missing", summary[audit.EventMissing])
	}
	return nil
}
