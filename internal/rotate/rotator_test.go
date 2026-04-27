package rotate_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envsync/internal/rotate"
)

func TestRotate_GeneratesNewValues(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "old"}
	rules := []rotate.Rule{{Key: "SECRET_KEY", Strategy: rotate.StrategyRandom, Length: 16}}

	out, results, err := rotate.Rotate(env, rules, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].OldValue != "old" {
		t.Errorf("expected old value 'old', got %q", results[0].OldValue)
	}
	if out["SECRET_KEY"] == "old" {
		t.Error("expected rotated value to differ from old value")
	}
	if len(out["SECRET_KEY"]) != 16 {
		t.Errorf("expected length 16, got %d", len(out["SECRET_KEY"]))
	}
}

func TestRotate_DryRunDoesNotChange(t *testing.T) {
	env := map[string]string{"API_SECRET": "original"}
	rules := []rotate.Rule{{Key: "API_SECRET", Strategy: rotate.StrategyAlphaNum, Length: 20}}

	out, results, err := rotate.Rotate(env, rules, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_SECRET"] != "original" {
		t.Errorf("dry run should not change value, got %q", out["API_SECRET"])
	}
	if results[0].Rotated {
		t.Error("dry run result should have Rotated=false")
	}
}

func TestRotate_MissingKeyReturnsError(t *testing.T) {
	env := map[string]string{}
	rules := []rotate.Rule{{Key: "MISSING", Strategy: rotate.StrategyRandom}}

	_, _, err := rotate.Rotate(env, rules, false)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRotate_UUIDStrategy(t *testing.T) {
	env := map[string]string{"TOKEN": "old-token"}
	rules := []rotate.Rule{{Key: "TOKEN", Strategy: rotate.StrategyUUID}}

	out, _, err := rotate.Rotate(env, rules, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	parts := strings.Split(out["TOKEN"], "-")
	if len(parts) != 5 {
		t.Errorf("expected UUID with 5 parts, got %q", out["TOKEN"])
	}
}

func TestRotate_UnknownStrategyReturnsError(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	rules := []rotate.Rule{{Key: "KEY", Strategy: "bogus"}}

	_, _, err := rotate.Rotate(env, rules, false)
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestRotate_PreservesUnrelatedKeys(t *testing.T) {
	env := map[string]string{"SECRET": "s", "KEEP": "unchanged"}
	rules := []rotate.Rule{{Key: "SECRET", Strategy: rotate.StrategyAlphaNum, Length: 8}}

	out, _, err := rotate.Rotate(env, rules, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEEP"] != "unchanged" {
		t.Errorf("expected KEEP to be unchanged, got %q", out["KEEP"])
	}
}
