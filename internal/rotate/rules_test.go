package rotate_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envsync/internal/rotate"
)

func TestParseRules_BasicEntries(t *testing.T) {
	input := `
# rotation config
SECRET_KEY   random 32
API_TOKEN    uuid
SESSION_ID   alphanum 24
`
	rules, err := rotate.ParseRules(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
	if rules[0].Key != "SECRET_KEY" || rules[0].Strategy != rotate.StrategyRandom || rules[0].Length != 32 {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
	if rules[1].Strategy != rotate.StrategyUUID {
		t.Errorf("expected uuid strategy, got %q", rules[1].Strategy)
	}
	if rules[2].Length != 24 {
		t.Errorf("expected length 24, got %d", rules[2].Length)
	}
}

func TestParseRules_SkipsCommentsAndBlanks(t *testing.T) {
	input := "\n# skip me\n\nKEY random\n"
	rules, err := rotate.ParseRules(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules))
	}
}

func TestParseRules_MissingStrategy_ReturnsError(t *testing.T) {
	input := "ONLY_KEY\n"
	_, err := rotate.ParseRules(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for missing strategy")
	}
}

func TestParseRules_InvalidLength_ReturnsError(t *testing.T) {
	input := "KEY random notanumber\n"
	_, err := rotate.ParseRules(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid length")
	}
}
