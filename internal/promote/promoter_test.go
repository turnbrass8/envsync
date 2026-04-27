package promote_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envsync/internal/envfile"
	"github.com	/yourorg/envsync/internal/manifest"
	"github.com/yourorg/envsync/internal/promote"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func buildManifest(keys ...string) *manifest.Manifest {
	m := &manifest.Manifest{}
	for _, k := range keys {
		m.Entries = append(m.Entries, manifest.Entry{Key: k})
	}
	return m
}

func TestPromote_CopiesKeysFromSrcToDst(t *testing.T) {
	src := writeTempEnv(t, "DB_HOST=prod.db\nDB_PORT=5432\n")
	dst := writeTempEnv(t, "APP_ENV=staging\n")

	mf := buildManifest("DB_HOST", "DB_PORT")
	results, err := promote.Promote(src, dst, mf, promote.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	promoted := 0
	for _, r := range results {
		if r.Promoted {
			promoted++
		}
	}
	if promoted != 2 {
		t.Errorf("expected 2 promoted, got %d", promoted)
	}

	dstEnv, _ := envfile.Parse(dst)
	if v, _ := dstEnv.Get("DB_HOST"); v != "prod.db" {
		t.Errorf("DB_HOST not written correctly, got %q", v)
	}
}

func TestPromote_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := writeTempEnv(t, "API_KEY=newsecret\n")
	dst := writeTempEnv(t, "API_KEY=oldsecret\n")

	mf := buildManifest("API_KEY")
	results, err := promote.Promote(src, dst, mf, promote.Options{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !results[0].Skipped {
		t.Error("expected key to be skipped")
	}

	dstEnv, _ := envfile.Parse(dst)
	if v, _ := dstEnv.Get("API_KEY"); v != "oldsecret" {
		t.Errorf("expected old value preserved, got %q", v)
	}
}

func TestPromote_OverwriteReplacesExisting(t *testing.T) {
	src := writeTempEnv(t, "API_KEY=newsecret\n")
	dst := writeTempEnv(t, "API_KEY=oldsecret\n")

	mf := buildManifest("API_KEY")
	_, err := promote.Promote(src, dst, mf, promote.Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dstEnv, _ := envfile.Parse(dst)
	if v, _ := dstEnv.Get("API_KEY"); v != "newsecret" {
		t.Errorf("expected new value, got %q", v)
	}
}

func TestPromote_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempEnv(t, "SECRET=abc123\n")
	dst := writeTempEnv(t, "OTHER=val\n")

	mf := buildManifest("SECRET")
	results, err := promote.Promote(src, dst, mf, promote.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !results[0].Promoted {
		t.Error("expected result to be marked promoted in dry run")
	}

	dstEnv, _ := envfile.Parse(dst)
	if _, ok := dstEnv.Get("SECRET"); ok {
		t.Error("expected SECRET not to be written in dry run")
	}
}

func TestPromote_MissingSourceKey_Skipped(t *testing.T) {
	src := writeTempEnv(t, "UNRELATED=x\n")
	dst := writeTempEnv(t, "")

	mf := buildManifest("MISSING_KEY")
	results, err := promote.Promote(src, dst, mf, promote.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !results[0].Skipped {
		t.Error("expected missing source key to be skipped")
	}
	_ = filepath.Join("") // keep import used
}
