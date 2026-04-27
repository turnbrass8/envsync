package mask_test

import (
	"testing"

	"github.com/user/envsync/internal/mask"
)

func TestApply_MasksMatchingKeys(t *testing.T) {
	m := mask.New([]mask.Rule{
		{Key: "PASSWORD"},
		{Key: "TOKEN", Mask: "[REDACTED]"},
	})
	env := map[string]string{
		"PASSWORD": "s3cr3t",
		"TOKEN":    "abc123",
		"HOST":     "localhost",
	}
	out := m.Apply(env)
	if out["PASSWORD"] != mask.DefaultMask {
		t.Errorf("expected PASSWORD to be masked, got %q", out["PASSWORD"])
	}
	if out["TOKEN"] != "[REDACTED]" {
		t.Errorf("expected TOKEN to be [REDACTED], got %q", out["TOKEN"])
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST to remain unchanged, got %q", out["HOST"])
	}
}

func TestApply_CaseInsensitiveKey(t *testing.T) {
	m := mask.New([]mask.Rule{{Key: "api_key"}})
	env := map[string]string{"API_KEY": "supersecret"}
	out := m.Apply(env)
	if out["API_KEY"] != mask.DefaultMask {
		t.Errorf("expected API_KEY masked, got %q", out["API_KEY"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	m := mask.New([]mask.Rule{{Key: "SECRET"}})
	env := map[string]string{"SECRET": "original"}
	_ = m.Apply(env)
	if env["SECRET"] != "original" {
		t.Error("Apply must not mutate the original map")
	}
}

func TestValue_MasksKnownKey(t *testing.T) {
	m := mask.New([]mask.Rule{{Key: "PASSWORD"}})
	if got := m.Value("PASSWORD", "hunter2"); got != mask.DefaultMask {
		t.Errorf("expected default mask, got %q", got)
	}
}

func TestValue_PassesThroughUnknownKey(t *testing.T) {
	m := mask.New([]mask.Rule{{Key: "PASSWORD"}})
	if got := m.Value("USERNAME", "alice"); got != "alice" {
		t.Errorf("expected original value, got %q", got)
	}
}

func TestDefaultRules_CoverCommonSecrets(t *testing.T) {
	rules := mask.DefaultRules()
	if len(rules) == 0 {
		t.Fatal("expected at least one default rule")
	}
	m := mask.New(rules)
	env := map[string]string{
		"API_KEY":  "key123",
		"PASSWORD": "pass",
		"HOST":     "example.com",
	}
	out := m.Apply(env)
	if out["API_KEY"] != mask.DefaultMask {
		t.Errorf("API_KEY should be masked")
	}
	if out["PASSWORD"] != mask.DefaultMask {
		t.Errorf("PASSWORD should be masked")
	}
	if out["HOST"] != "example.com" {
		t.Errorf("HOST should not be masked")
	}
}
