// Package schema provides validation of .env files against a JSON Schema-like
// definition, ensuring key types, formats, and constraints are enforced.
package schema

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// FieldType represents the expected type of an env variable's value.
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInt     FieldType = "int"
	TypeBool    FieldType = "bool"
	TypeURL     FieldType = "url"
	TypeEmail   FieldType = "email"
)

// Field describes the schema definition for a single env key.
type Field struct {
	Key      string
	Type     FieldType
	Required bool
	Pattern  string // optional regex pattern
}

// Schema holds a collection of field definitions.
type Schema struct {
	Fields []Field
}

// Violation describes a single schema violation.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

var (
	urlPattern   = regexp.MustCompile(`^https?://[^\s]+$`)
	emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

// Validate checks the provided env map against the schema and returns any violations.
func (s *Schema) Validate(env map[string]string) []Violation {
	var violations []Violation

	for _, field := range s.Fields {
		val, ok := env[field.Key]
		if !ok || val == "" {
			if field.Required {
				violations = append(violations, Violation{Key: field.Key, Message: "required key is missing or empty"})
			}
			continue
		}

		if v := validateType(field.Key, val, field.Type); v != nil {
			violations = append(violations, *v)
			continue
		}

		if field.Pattern != "" {
			re, err := regexp.Compile(field.Pattern)
			if err != nil {
				violations = append(violations, Violation{Key: field.Key, Message: fmt.Sprintf("invalid pattern %q: %v", field.Pattern, err)})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, Violation{Key: field.Key, Message: fmt.Sprintf("value %q does not match pattern %q", val, field.Pattern)})
			}
		}
	}

	return violations
}

func validateType(key, val string, t FieldType) *Violation {
	switch t {
	case TypeInt:
		if _, err := strconv.Atoi(strings.TrimSpace(val)); err != nil {
			return &Violation{Key: key, Message: fmt.Sprintf("expected int, got %q", val)}
		}
	case TypeBool:
		lower := strings.ToLower(strings.TrimSpace(val))
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			return &Violation{Key: key, Message: fmt.Sprintf("expected bool, got %q", val)}
		}
	case TypeURL:
		if !urlPattern.MatchString(val) {
			return &Violation{Key: key, Message: fmt.Sprintf("expected URL, got %q", val)}
		}
	case TypeEmail:
		if !emailPattern.MatchString(val) {
			return &Violation{Key: key, Message: fmt.Sprintf("expected email, got %q", val)}
		}
	}
	return nil
}
