package format

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestFormat_NoChanges_NotChanged(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	res, err := Format(p, Options{QuoteStyle: "none"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Changed {
		t.Errorf("expected unchanged, got changed; output: %q", res.Formatted)
	}
}

func TestFormat_DoubleQuoteStyle_AddsQuotes(t *testing.T) {
	p := writeTempEnv(t, "MSG=hello world\n")
	res, err := Format(p, Options{QuoteStyle: "double", DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Formatted, `"hello world"`) {
		t.Errorf("expected double-quoted value, got: %q", res.Formatted)
	}
}

func TestFormat_SingleQuoteStyle(t *testing.T) {
	p := writeTempEnv(t, "KEY=value\n")
	res, err := Format(p, Options{QuoteStyle: "single", DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Formatted, "'value'") {
		t.Errorf("expected single-quoted value, got: %q", res.Formatted)
	}
}

func TestFormat_SpaceAroundEquals(t *testing.T) {
	p := writeTempEnv(t, "KEY=value\n")
	res, err := Format(p, Options{QuoteStyle: "none", SpaceAroundEquals: true, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Formatted, "KEY = value") {
		t.Errorf("expected spaces around '=', got: %q", res.Formatted)
	}
}

func TestFormat_DryRun_DoesNotWriteFile(t *testing.T) {
	original := "KEY=value\n"
	p := writeTempEnv(t, original)
	_, err := Format(p, Options{QuoteStyle: "single", DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(p)
	if string(got) != original {
		t.Errorf("dry-run must not modify file; got %q", string(got))
	}
}

func TestFormat_WritesFile_WhenChanged(t *testing.T) {
	p := writeTempEnv(t, "MSG=hello world\n")
	_, err := Format(p, Options{QuoteStyle: "double"})
	if err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(p)
	if !strings.Contains(string(got), `"hello world"`) {
		t.Errorf("expected file to contain double-quoted value; got %q", string(got))
	}
}

func TestFormat_BadFile_ReturnsError(t *testing.T) {
	_, err := Format("/nonexistent/.env", Options{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
