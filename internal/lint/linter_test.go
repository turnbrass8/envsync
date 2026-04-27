package lint

import (
	"strings"
	"testing"
)

func TestLint_CleanEnv_NoFindings(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/mydb",
		"APP_PORT":     "8080",
	}
	findings := Lint(env)
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d: %v", len(findings), findings)
	}
}

func TestLint_LowercaseKey_Warning(t *testing.T) {
	env := map[string]string{"db_host": "localhost"}
	findings := Lint(env)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != SeverityWarning {
		t.Errorf("expected warning, got %s", findings[0].Severity)
	}
	if !strings.Contains(findings[0].Message, "UPPER_SNAKE_CASE") {
		t.Errorf("unexpected message: %s", findings[0].Message)
	}
}

func TestLint_EmptyValue_Warning(t *testing.T) {
	env := map[string]string{"APP_NAME": ""}
	findings := Lint(env)
	found := false
	for _, f := range findings {
		if strings.Contains(f.Message, "empty") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected empty-value warning")
	}
}

func TestLint_LeadingWhitespace_Warning(t *testing.T) {
	env := map[string]string{"APP_HOST": "  localhost"}
	findings := Lint(env)
	found := false
	for _, f := range findings {
		if strings.Contains(f.Message, "whitespace") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected whitespace warning")
	}
}

func TestLint_ShortSecret_Error(t *testing.T) {
	env := map[string]string{"APP_SECRET": "abc"}
	findings := Lint(env)
	found := false
	for _, f := range findings {
		if f.Severity == SeverityError && strings.Contains(f.Message, "short") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected short-secret error finding")
	}
}

func TestFinding_String_ContainsFields(t *testing.T) {
	f := Finding{Key: "MY_KEY", Message: "some issue", Severity: SeverityWarning}
	s := f.String()
	if !strings.Contains(s, "MY_KEY") || !strings.Contains(s, "warning") || !strings.Contains(s, "some issue") {
		t.Errorf("unexpected String() output: %s", s)
	}
}
