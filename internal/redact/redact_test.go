package redact_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/redact"
)

func TestIsSensitive_MatchesDefaultPatterns(t *testing.T) {
	cases := []struct {
		key  string
		want bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"AUTH_TOKEN", true},
		{"PRIVATE_KEY", true},
		{"APP_SECRET", true},
		{"DATABASE_URL", false},
		{"PORT", false},
		{"LOG_LEVEL", false},
	}
	for _, tc := range cases {
		got := redact.IsSensitive(tc.key, redact.DefaultPatterns)
		if got != tc.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}

func TestApply_RedactsSensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
		"PORT":        "8080",
		"APP_NAME":    "envsync",
	}

	result := redact.Apply(env, redact.Options{})

	if result["DB_PASSWORD"] != "***REDACTED***" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != "***REDACTED***" {
		t.Errorf("expected API_KEY to be redacted, got %q", result["API_KEY"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected PORT to be unchanged, got %q", result["PORT"])
	}
	if result["APP_NAME"] != "envsync" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", result["APP_NAME"])
	}
}

func TestApply_CustomMask(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "topsecret", "HOST": "localhost"}
	opts := redact.Options{Mask: "<hidden>"}
	result := redact.Apply(env, opts)
	if result["SECRET_KEY"] != "<hidden>" {
		t.Errorf("expected custom mask, got %q", result["SECRET_KEY"])
	}
}

func TestApply_CustomPatterns(t *testing.T) {
	env := map[string]string{"INTERNAL_CODE": "xyz", "PORT": "9000"}
	opts := redact.Options{Patterns: []string{"INTERNAL"}}
	result := redact.Apply(env, opts)
	if result["INTERNAL_CODE"] != "***REDACTED***" {
		t.Errorf("expected INTERNAL_CODE to be redacted")
	}
	if result["PORT"] != "9000" {
		t.Errorf("expected PORT to be unchanged")
	}
}

func TestValue_SingleKey(t *testing.T) {
	opts := redact.Options{}
	if got := redact.Value("DB_TOKEN", "secret", opts); got != "***REDACTED***" {
		t.Errorf("expected redacted value, got %q", got)
	}
	if got := redact.Value("REGION", "us-east-1", opts); got != "us-east-1" {
		t.Errorf("expected plain value, got %q", got)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"PASSWORD": "original"}
	redact.Apply(env, redact.Options{})
	if env["PASSWORD"] != "original" {
		t.Error("Apply must not mutate the original map")
	}
}
