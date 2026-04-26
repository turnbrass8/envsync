package interpolate_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/interpolate"
)

func TestResolve_SimpleDollarBrace(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	got, err := interpolate.Resolve("http://${HOST}:8080", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "http://localhost:8080" {
		t.Errorf("got %q, want %q", got, "http://localhost:8080")
	}
}

func TestResolve_BareVariable(t *testing.T) {
	env := map[string]string{"PORT": "5432"}
	got, err := interpolate.Resolve("$PORT", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "5432" {
		t.Errorf("got %q, want %q", got, "5432")
	}
}

func TestResolve_DefaultValue(t *testing.T) {
	env := map[string]string{}
	got, err := interpolate.Resolve("${MISSING:-fallback}", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "fallback" {
		t.Errorf("got %q, want %q", got, "fallback")
	}
}

func TestResolve_UnresolvedReturnsError(t *testing.T) {
	env := map[string]string{}
	_, err := interpolate.Resolve("${MISSING}", env)
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestResolve_NoReferences(t *testing.T) {
	env := map[string]string{}
	got, err := interpolate.Resolve("plain-value", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "plain-value" {
		t.Errorf("got %q, want %q", got, "plain-value")
	}
}

func TestResolveAll_ResolvesMap(t *testing.T) {
	env := map[string]string{
		"BASE": "http://localhost",
		"URL":  "${BASE}/api",
	}
	if err := interpolate.ResolveAll(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["URL"] != "http://localhost/api" {
		t.Errorf("got %q, want %q", env["URL"], "http://localhost/api")
	}
}

func TestResolveAll_ErrorPropagates(t *testing.T) {
	env := map[string]string{"A": "${UNDEFINED}"}
	if err := interpolate.ResolveAll(env); err == nil {
		t.Fatal("expected error for unresolved variable in ResolveAll")
	}
}
