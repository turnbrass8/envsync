package normalize

import (
	"testing"
)

func TestNormalize_UppercaseKeys(t *testing.T) {
	env := map[string]string{"db_host": "localhost", "DB_PORT": "5432"}
	res, err := Normalize(env, Options{UppercaseKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be present")
	}
	if _, ok := res.Env["DB_PORT"]; !ok {
		t.Error("expected DB_PORT to be present")
	}
	if len(res.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(res.Changes))
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	env := map[string]string{"KEY": "  hello  "}
	res, err := Normalize(env, Options{TrimValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", res.Env["KEY"])
	}
	if len(res.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(res.Changes))
	}
}

func TestNormalize_RemoveEmptyValues(t *testing.T) {
	env := map[string]string{"PRESENT": "yes", "EMPTY": ""}
	res, err := Normalize(env, Options{RemoveEmptyValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["EMPTY"]; ok {
		t.Error("expected EMPTY to be removed")
	}
	if _, ok := res.Env["PRESENT"]; !ok {
		t.Error("expected PRESENT to remain")
	}
}

func TestNormalize_DryRunDoesNotModify(t *testing.T) {
	env := map[string]string{"lower_key": "  val  "}
	opts := Options{UppercaseKeys: true, TrimValues: true, DryRun: true}
	res, err := Normalize(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["lower_key"]; !ok {
		t.Error("dry-run should preserve original key")
	}
	if res.Env["lower_key"] != "  val  " {
		t.Error("dry-run should preserve original value")
	}
	if len(res.Changes) == 0 {
		t.Error("dry-run should still report changes")
	}
}

func TestNormalize_NilEnv_ReturnsError(t *testing.T) {
	_, err := Normalize(nil, DefaultOptions())
	if err == nil {
		t.Error("expected error for nil env")
	}
}

func TestChange_String(t *testing.T) {
	c := Change{Key: "FOO", Reason: "value trimmed"}
	if c.String() != "FOO: value trimmed" {
		t.Errorf("unexpected string: %s", c.String())
	}
}
