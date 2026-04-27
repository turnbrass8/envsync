package schema_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/schema"
)

func makeSchema(fields ...schema.Field) *schema.Schema {
	return &schema.Schema{Fields: fields}
}

func TestValidate_AllValid(t *testing.T) {
	s := makeSchema(
		schema.Field{Key: "PORT", Type: schema.TypeInt, Required: true},
		schema.Field{Key: "DEBUG", Type: schema.TypeBool},
		schema.Field{Key: "API_URL", Type: schema.TypeURL},
	)
	env := map[string]string{"PORT": "8080", "DEBUG": "true", "API_URL": "https://api.example.com"}
	violations := s.Validate(env)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestValidate_RequiredMissing(t *testing.T) {
	s := makeSchema(schema.Field{Key: "SECRET", Type: schema.TypeString, Required: true})
	violations := s.Validate(map[string]string{})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "SECRET" {
		t.Errorf("expected key SECRET, got %s", violations[0].Key)
	}
}

func TestValidate_InvalidInt(t *testing.T) {
	s := makeSchema(schema.Field{Key: "PORT", Type: schema.TypeInt})
	violations := s.Validate(map[string]string{"PORT": "not-a-number"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_InvalidBool(t *testing.T) {
	s := makeSchema(schema.Field{Key: "FLAG", Type: schema.TypeBool})
	violations := s.Validate(map[string]string{"FLAG": "yes"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_ValidBoolVariants(t *testing.T) {
	s := makeSchema(schema.Field{Key: "FLAG", Type: schema.TypeBool})
	for _, v := range []string{"true", "false", "1", "0", "TRUE", "FALSE"} {
		violations := s.Validate(map[string]string{"FLAG": v})
		if len(violations) != 0 {
			t.Errorf("expected no violations for %q, got %v", v, violations)
		}
	}
}

func TestValidate_InvalidURL(t *testing.T) {
	s := makeSchema(schema.Field{Key: "ENDPOINT", Type: schema.TypeURL})
	violations := s.Validate(map[string]string{"ENDPOINT": "not-a-url"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_InvalidEmail(t *testing.T) {
	s := makeSchema(schema.Field{Key: "EMAIL", Type: schema.TypeEmail})
	violations := s.Validate(map[string]string{"EMAIL": "notanemail"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	s := makeSchema(schema.Field{Key: "ENV", Type: schema.TypeString, Pattern: `^(dev|staging|prod)$`})
	violations := s.Validate(map[string]string{"ENV": "local"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	s := makeSchema(schema.Field{Key: "ENV", Type: schema.TypeString, Pattern: `^(dev|staging|prod)$`})
	violations := s.Validate(map[string]string{"ENV": "prod"})
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestViolation_Error(t *testing.T) {
	v := schema.Violation{Key: "FOO", Message: "something wrong"}
	if v.Error() != "FOO: something wrong" {
		t.Errorf("unexpected error string: %s", v.Error())
	}
}
