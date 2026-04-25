package export

import (
	"strings"
	"testing"
)

func TestExport_DotenvFormat(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"PORT":    "8080",
	}
	var buf strings.Builder
	err := Export(env, Options{Format: FormatDotenv, Sorted: true, Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got:\n%s", got)
	}
	if !strings.Contains(got, "PORT=8080") {
		t.Errorf("expected PORT=8080 in output, got:\n%s", got)
	}
}

func TestExport_ShellFormat(t *testing.T) {
	env := map[string]string{"DB_URL": "postgres://localhost/db"}
	var buf strings.Builder
	err := Export(env, Options{Format: FormatShell, Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(buf.String(), "export DB_URL=") {
		t.Errorf("expected shell export prefix, got: %s", buf.String())
	}
}

func TestExport_JSONFormat(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	var buf strings.Builder
	err := Export(env, Options{Format: FormatJSON, Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, `"KEY"`) {
		t.Errorf("expected JSON key, got: %s", got)
	}
	if !strings.Contains(got, `"value"`) {
		t.Errorf("expected JSON value, got: %s", got)
	}
}

func TestExport_QuotesValueWithSpaces(t *testing.T) {
	env := map[string]string{"MSG": "hello world"}
	var buf strings.Builder
	_ = Export(env, Options{Format: FormatDotenv, Writer: &buf})
	if !strings.Contains(buf.String(), `"hello world"`) {
		t.Errorf("expected quoted value, got: %s", buf.String())
	}
}

func TestExport_SortedOutput(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	var buf strings.Builder
	_ = Export(env, Options{Format: FormatDotenv, Sorted: true, Writer: &buf})
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected A_KEY first, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected Z_KEY last, got: %s", lines[2])
	}
}

func TestExport_EmptyMap(t *testing.T) {
	env := map[string]string{}
	var buf strings.Builder
	err := Export(env, Options{Format: FormatDotenv, Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error for empty map: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "" {
		t.Errorf("expected empty output for empty map, got: %s", got)
	}
}
