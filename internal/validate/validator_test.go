package validate_test

import (
	"testing"

	"github.com/yourusername/envsync/internal/manifest"
	"github.com/yourusername/envsync/internal/validate"
)

func entries(keys ...string) []manifest.Entry {
	out := make([]manifest.Entry, len(keys))
	for i, k := range keys {
		out[i] = manifest.Entry{Key: k}
	}
	return out
}

func requiredEntry(key string) manifest.Entry {
	return manifest.Entry{Key: key, Required: true}
}

func TestValidate_AllPresent_NoRules(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := validate.Validate(env, entries("FOO", "BAZ"), nil)
	if !result.OK() {
		t.Fatalf("expected no violations, got: %s", result.Summary())
	}
}

func TestValidate_RequiredMissing(t *testing.T) {
	env := map[string]string{}
	result := validate.Validate(env, []manifest.Entry{requiredEntry("DB_URL")}, nil)
	if result.OK() {
		t.Fatal("expected violation for missing required key")
	}
	if len(result.Violations) != 1 || result.Violations[0].Key != "DB_URL" {
		t.Fatalf("unexpected violations: %+v", result.Violations)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	env := map[string]string{"PORT": "abc"}
	rules := []validate.Rule{{Key: "PORT", Pattern: `^\d+$`}}
	result := validate.Validate(env, entries("PORT"), rules)
	if result.OK() {
		t.Fatal("expected violation for pattern mismatch")
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	rules := []validate.Rule{{Key: "PORT", Pattern: `^\d+$`}}
	result := validate.Validate(env, entries("PORT"), rules)
	if !result.OK() {
		t.Fatalf("expected no violations, got: %s", result.Summary())
	}
}

func TestValidate_NonEmptyViolation(t *testing.T) {
	env := map[string]string{"SECRET": "   "}
	rules := []validate.Rule{{Key: "SECRET", NonEmpty: true}}
	result := validate.Validate(env, entries("SECRET"), rules)
	if result.OK() {
		t.Fatal("expected violation for empty value")
	}
}

func TestValidate_Summary_OK(t *testing.T) {
	result := &validate.Result{}
	if result.Summary() != "all validations passed" {
		t.Fatalf("unexpected summary: %s", result.Summary())
	}
}

func TestValidate_Summary_WithViolations(t *testing.T) {
	result := &validate.Result{
		Violations: []validate.Violation{
			{Key: "FOO", Message: "required key is missing"},
		},
	}
	if result.OK() {
		t.Fatal("expected not OK")
	}
	if result.Summary() == "" {
		t.Fatal("expected non-empty summary")
	}
}
