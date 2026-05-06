// Package archive provides functionality to snapshot and bundle .env files
// into a compressed archive for backup or transfer purposes.
package archive

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Metadata holds information written into the archive manifest.
type Metadata struct {
	CreatedAt time.Time         `json:"created_at"`
	Files     []string          `json:"files"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// Options configures an Archive operation.
type Options struct {
	// Labels are arbitrary key/value pairs stored in the archive manifest.
	Labels map[string]string
	// DryRun reports what would be archived without writing the file.
	DryRun bool
}

// Archive bundles the given env files into a zip archive at dest.
// A JSON manifest (envsync-manifest.json) is included in the archive.
func Archive(dest string, envFiles []string, opts Options) (*Metadata, error) {
	if len(envFiles) == 0 {
		return nil, fmt.Errorf("archive: no files provided")
	}

	// Validate all source files exist before writing anything.
	for _, f := range envFiles {
		if _, err := os.Stat(f); err != nil {
			return nil, fmt.Errorf("archive: source file %q not found: %w", f, err)
		}
	}

	meta := &Metadata{
		CreatedAt: time.Now().UTC(),
		Files:     envFiles,
		Labels:    opts.Labels,
	}

	if opts.DryRun {
		return meta, nil
	}

	out, err := os.Create(dest)
	if err != nil {
		return nil, fmt.Errorf("archive: create %q: %w", dest, err)
	}
	defer out.Close()

	zw := zip.NewWriter(out)
	defer zw.Close()

	for _, src := range envFiles {
		if err := addFile(zw, src); err != nil {
			return nil, err
		}
	}

	if err := addManifest(zw, meta); err != nil {
		return nil, err
	}

	return meta, nil
}

func addFile(zw *zip.Writer, src string) error {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("archive: open %q: %w", src, err)
	}
	defer f.Close()

	w, err := zw.Create(filepath.Base(src))
	if err != nil {
		return fmt.Errorf("archive: zip entry for %q: %w", src, err)
	}

	if _, err := io.Copy(w, f); err != nil {
		return fmt.Errorf("archive: write %q: %w", src, err)
	}
	return nil
}

func addManifest(zw *zip.Writer, meta *Metadata) error {
	w, err := zw.Create("envsync-manifest.json")
	if err != nil {
		return fmt.Errorf("archive: create manifest entry: %w", err)
	}
	return json.NewEncoder(w).Encode(meta)
}
