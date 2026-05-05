package patch_test

import (
	"strings"
	"testing"

	"github.com/your-org/envsync/internal/patch"
)

func TestParseRules_SetAndDelete(t *testing.T) {
	input := `
# comment
FOO=bar
BAZ=qux
-OLD_KEY
`
	ops, err := patch.ParseRules(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 3 {
		t.Fatalf("expected 3 ops, got %d", len(ops))
	}
	if ops[0].Key != "FOO" || ops[0].Value != "bar" || ops[0].Delete {
		t.Errorf("unexpected op[0]: %+v", ops[0])
	}
	if ops[2].Key != "OLD_KEY" || !ops[2].Delete {
		t.Errorf("unexpected op[2]: %+v", ops[2])
	}
}

func TestParseRules_SkipsCommentsAndBlanks(t *testing.T) {
	input := "\n# skip me\n\nA=1\n"
	ops, err := patch.ParseRules(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 1 {
		t.Fatalf("expected 1 op, got %d", len(ops))
	}
}

func TestParseRules_InvalidLine_ReturnsError(t *testing.T) {
	_, err := patch.ParseRules(strings.NewReader("NODIVIDER\n"))
	if err == nil {
		t.Error("expected error for line without '='")
	}
}

func TestParseRules_EmptyDeleteKey_ReturnsError(t *testing.T) {
	_, err := patch.ParseRules(strings.NewReader("-\n"))
	if err == nil {
		t.Error("expected error for '-' with no key")
	}
}

func TestParseRules_EmptyInput_ReturnsNil(t *testing.T) {
	ops, err := patch.ParseRules(strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 0 {
		t.Errorf("expected empty ops, got %v", ops)
	}
}
