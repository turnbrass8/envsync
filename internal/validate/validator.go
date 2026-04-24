package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yourusername/envsync/internal/manifest"
)

// Rule defines a validation rule for an environment variable.
type Rule struct {
	Key      string
	Pattern  string // optional regex pattern
	NonEmpty bool   // if true, value must not be empty
}

// Violation represents a single validation failure.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("key %q: %s", v.Key, v.Message)
}

// Result holds all violations found during validation.
type Result struct {
	Violations []Violation
}

func (r *Result) OK() bool {
	return len(r.Violations) == 0
}

func (r *Result) Summary() string {
	if r.OK() {
		return "all validations passed"
	}
	lines := make([]string, len(r.Violations))
	for i, v := range r.Violations {
		lines[i] = "  - " + v.Error()
	}
	return fmt.Sprintf("%d validation error(s):\n%s", len(r.Violations), strings.Join(lines, "\n"))
}

// Validate checks env values against manifest entries and optional rules.
func Validate(env map[string]string, entries []manifest.Entry, rules []Rule) *Result {
	ruleMap := make(map[string]Rule, len(rules))
	for _, r := range rules {
		ruleMap[r.Key] = r
	}

	result := &Result{}

	for _, entry := range entries {
		val, exists := env[entry.Key]

		if entry.Required && !exists {
			result.Violations = append(result.Violations, Violation{
				Key:     entry.Key,
				Message: "required key is missing",
			})
			continue
		}

		if !exists {
			continue
		}

		if rule, ok := ruleMap[entry.Key]; ok {
			if rule.NonEmpty && strings.TrimSpace(val) == "" {
				result.Violations = append(result.Violations, Violation{
					Key:     entry.Key,
					Message: "value must not be empty",
				})
			}
			if rule.Pattern != "" {
				matched, err := regexp.MatchString(rule.Pattern, val)
				if err != nil {
					result.Violations = append(result.Violations, Violation{
						Key:     entry.Key,
						Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
					})
				} else if !matched {
					result.Violations = append(result.Violations, Violation{
						Key:     entry.Key,
						Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern),
					})
				}
			}
		}
	}

	return result
}
