package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envsync/internal/template"
)

func TestRender_BasicSubstitution(t *testing.T) {
	src := `APP_HOST={{ .HOST }}
APP_PORT={{ .PORT }}
`
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	out, err := template.Render(src, env, template.RenderOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "APP_HOST=localhost\nAPP_PORT=8080\n"
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestRender_DefaultFunc(t *testing.T) {
	src := `LEVEL={{ default "info" .LOG_LEVEL }}`
	env := map[string]string{}
	out, err := template.Render(src, env, template.RenderOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "LEVEL=info" {
		t.Errorf("got %q", out)
	}
}

func TestRender_UpperFunc(t *testing.T) {
	src := `ENV={{ upper .ENV }}`
	env := map[string]string{"ENV": "production"}
	out, err := template.Render(src, env, template.RenderOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "ENV=PRODUCTION" {
		t.Errorf("got %q", out)
	}
}

func TestRender_StrictMode_MissingKey_ReturnsError(t *testing.T) {
	src := `VALUE={{ .MISSING }}`
	env := map[string]string{}
	_, err := template.Render(src, env, template.RenderOptions{Strict: true})
	if err == nil {
		t.Fatal("expected error for missing key in strict mode, got nil")
	}
}

func TestRender_NonStrict_MissingKey_EmptyString(t *testing.T) {
	src := `VALUE={{ .MISSING }}`
	env := map[string]string{}
	out, err := template.Render(src, env, template.RenderOptions{Strict: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "VALUE=" {
		t.Errorf("got %q", out)
	}
}

func TestRenderFile_ReadsAndRenders(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "env.tmpl")
	content := `DB_URL={{ .DB_HOST }}:{{ .DB_PORT }}/{{ .DB_NAME }}`
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	env := map[string]string{"DB_HOST": "db", "DB_PORT": "5432", "DB_NAME": "mydb"}
	out, err := template.RenderFile(p, env, template.RenderOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "DB_URL=db:5432/mydb" {
		t.Errorf("got %q", out)
	}
}

func TestRenderFile_MissingFile_ReturnsError(t *testing.T) {
	_, err := template.RenderFile("/nonexistent/path.tmpl", nil, template.RenderOptions{})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
