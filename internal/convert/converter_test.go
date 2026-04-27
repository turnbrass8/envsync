package convert_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envsync/internal/convert"
)

var sampleEnv = map[string]string{
	"APP_HOST": "localhost",
	"APP_PORT": "8080",
	"APP_NAME": "my app",
}

func TestConvert_DotenvFormat(t *testing.T) {
	out, err := convert.Convert(sampleEnv, convert.Options{Format: convert.FormatDotenv, Sorted: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_HOST=localhost") {
		t.Errorf("expected APP_HOST=localhost in output, got:\n%s", out)
	}
	if !strings.Contains(out, `APP_NAME="my app"`) {
		t.Errorf("expected quoted APP_NAME in output, got:\n%s", out)
	}
}

func TestConvert_JSONFormat(t *testing.T) {
	out, err := convert.Convert(sampleEnv, convert.Options{Format: convert.FormatJSON, Sorted: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"APP_HOST": "localhost"`) {
		t.Errorf("expected JSON key/value, got:\n%s", out)
	}
	if !strings.HasPrefix(out, "{") {
		t.Errorf("expected JSON object, got:\n%s", out)
	}
}

func TestConvert_YAMLFormat(t *testing.T) {
	out, err := convert.Convert(sampleEnv, convert.Options{Format: convert.FormatYAML, Sorted: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_PORT:") {
		t.Errorf("expected YAML key, got:\n%s", out)
	}
}

func TestConvert_TOMLFormat(t *testing.T) {
	out, err := convert.Convert(sampleEnv, convert.Options{Format: convert.FormatTOML, Sorted: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_HOST = ") {
		t.Errorf("expected TOML assignment, got:\n%s", out)
	}
}

func TestConvert_UnsupportedFormat(t *testing.T) {
	_, err := convert.Convert(sampleEnv, convert.Options{Format: "xml"})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestConvert_SortedOutput(t *testing.T) {
	out, err := convert.Convert(sampleEnv, convert.Options{Format: convert.FormatDotenv, Sorted: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "APP_HOST") {
		t.Errorf("expected first line to be APP_HOST, got %s", lines[0])
	}
}

func TestConvert_StripPrefix(t *testing.T) {
	env := map[string]string{"PROD_DB_URL": "postgres://localhost/db"}
	out, err := convert.Convert(env, convert.Options{Format: convert.FormatDotenv, Prefix: "PROD_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_URL=") {
		t.Errorf("expected prefix stripped, got: %s", out)
	}
	if strings.Contains(out, "PROD_") {
		t.Errorf("expected prefix to be removed, got: %s", out)
	}
}
