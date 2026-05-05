package defaults_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/defaults"
)

func TestApply_MissingKeysAreAdded(t *testing.T) {
	env := map[string]string{"EXISTING": "yes"}
	defs := map[string]string{"NEW_KEY": "default_val", "EXISTING": "other"}

	res, err := defaults.Apply(env, defs, defaults.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 1 || res.Applied[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY applied, got %v", res.Applied)
	}
	if env["NEW_KEY"] != "default_val" {
		t.Errorf("expected default_val, got %q", env["NEW_KEY"])
	}
	if env["EXISTING"] != "yes" {
		t.Errorf("existing key should not be overwritten")
	}
}

func TestApply_OverwriteReplacesExisting(t *testing.T) {
	env := map[string]string{"KEY": "old"}
	defs := map[string]string{"KEY": "new"}

	res, err := defaults.Apply(env, defs, defaults.Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwritten) != 1 {
		t.Errorf("expected 1 overwritten, got %v", res.Overwritten)
	}
	if env["KEY"] != "new" {
		t.Errorf("expected new, got %q", env["KEY"])
	}
}

func TestApply_DryRunDoesNotModify(t *testing.T) {
	env := map[string]string{"A": "1"}
	defs := map[string]string{"B": "2"}

	res, err := defaults.Apply(env, defs, defaults.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied in dry-run report")
	}
	if _, ok := env["B"]; ok {
		t.Errorf("dry-run should not write to env map")
	}
}

func TestApply_NilEnvReturnsError(t *testing.T) {
	_, err := defaults.Apply(nil, map[string]string{"K": "v"}, defaults.Options{})
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestResult_Summary(t *testing.T) {
	r := defaults.Result{
		Applied:     []string{"A", "B"},
		Skipped:     []string{"C"},
		Overwritten: []string{},
	}
	s := r.Summary()
	if s != "applied=2 skipped=1 overwritten=0" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestApply_SameValueNoOverwrite(t *testing.T) {
	env := map[string]string{"KEY": "same"}
	defs := map[string]string{"KEY": "same"}

	res, _ := defaults.Apply(env, defs, defaults.Options{Overwrite: true})
	if len(res.Overwritten) != 0 {
		t.Errorf("identical value should not count as overwritten")
	}
	if len(res.Skipped) != 1 {
		t.Errorf("identical value should be skipped")
	}
}
