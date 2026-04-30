package sanitize_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/sanitize"
)

func TestSanitize_TrimSpace(t *testing.T) {
	env := map[string]string{
		"KEY": "  hello world  ",
	}
	opts := sanitize.DefaultOptions()
	r := sanitize.Sanitize(env, opts)

	if got := r.Env["KEY"]; got != "hello world" {
		t.Errorf("expected trimmed value, got %q", got)
	}
	if len(r.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(r.Changes))
	}
}

func TestSanitize_RemoveControlChars(t *testing.T) {
	env := map[string]string{
		"KEY": "val\x01ue\x07",
	}
	opts := sanitize.DefaultOptions()
	r := sanitize.Sanitize(env, opts)

	if got := r.Env["KEY"]; got != "value" {
		t.Errorf("expected control chars removed, got %q", got)
	}
}

func TestSanitize_NormalizeNewlines(t *testing.T) {
	env := map[string]string{
		"KEY": "line1\r\nline2\rline3",
	}
	opts := sanitize.DefaultOptions()
	opts.TrimSpace = false
	r := sanitize.Sanitize(env, opts)

	expected := "line1\nline2\nline3"
	if got := r.Env["KEY"]; got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSanitize_UppercaseKeys(t *testing.T) {
	env := map[string]string{
		"db_host": "localhost",
	}
	opts := sanitize.DefaultOptions()
	opts.UppercaseKeys = true
	r := sanitize.Sanitize(env, opts)

	if _, ok := r.Env["DB_HOST"]; !ok {
		t.Error("expected key to be uppercased to DB_HOST")
	}
	if len(r.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(r.Changes))
	}
}

func TestSanitize_NoChanges(t *testing.T) {
	env := map[string]string{
		"KEY": "clean_value",
	}
	opts := sanitize.DefaultOptions()
	r := sanitize.Sanitize(env, opts)

	if len(r.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(r.Changes))
	}
}

func TestSanitize_PreservesTabsAndNewlines(t *testing.T) {
	env := map[string]string{
		"KEY": "col1\tcol2",
	}
	opts := sanitize.DefaultOptions()
	r := sanitize.Sanitize(env, opts)

	if got := r.Env["KEY"]; got != "col1\tcol2" {
		t.Errorf("expected tab preserved, got %q", got)
	}
}

func TestSanitize_EmptyEnv(t *testing.T) {
	r := sanitize.Sanitize(map[string]string{}, sanitize.DefaultOptions())
	if len(r.Env) != 0 {
		t.Error("expected empty result")
	}
	if len(r.Changes) != 0 {
		t.Error("expected no changes")
	}
}
