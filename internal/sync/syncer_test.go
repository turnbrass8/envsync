package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envsync/internal/envfile"
	"github.com/yourorg/envsync/internal/manifest"
	"github.com/yourorg/envsync/internal/sync"
)

func buildManifest(keys []manifest.Key) *manifest.Manifest {
	return &manifest.Manifest{Keys: keys}
}

func TestSync_AppliesDefaults(t *testing.T) {
	man := buildManifest([]manifest.Key{
		{Name: "EXISTING", Required: false},
		{Name: "MISSING", Required: false, Default: "fallback"},
	})
	target := envfile.Env{"EXISTING": "yes"}

	tmp := filepath.Join(t.TempDir(), ".env")
	s := sync.New(false)
	res, err := s.Sync(man, target, tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 1 || res.Applied[0] != `MISSING (default: "fallback")` {
		t.Errorf("expected MISSING to be applied, got %v", res.Applied)
	}
}

func TestSync_DryRunDoesNotWrite(t *testing.T) {
	man := buildManifest([]manifest.Key{
		{Name: "KEY", Default: "val"},
	})
	target := envfile.Env{}
	tmp := filepath.Join(t.TempDir(), ".env")

	s := sync.New(true)
	_, err := s.Sync(man, target, tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(tmp); !os.IsNotExist(statErr) {
		t.Error("dry-run should not create the file")
	}
}

func TestSync_RequiredMissingReturnsError(t *testing.T) {
	man := buildManifest([]manifest.Key{
		{Name: "SECRET", Required: true},
	})
	target := envfile.Env{}
	tmp := filepath.Join(t.TempDir(), ".env")

	s := sync.New(true)
	res, _ := s.Sync(man, target, tmp)
	if len(res.Errors) == 0 {
		t.Error("expected an error for missing required key")
	}
}
