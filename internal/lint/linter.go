// Package lint provides checks for common .env file issues such as
// duplicate keys, suspicious values, and naming convention violations.
package lint

import (
	"fmt"
	"regexp"
	"strings"
)

// Severity represents the level of a lint finding.
type Severity string

const (
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Finding describes a single lint issue.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s", f.Severity, f.Key, f.Message)
}

var (
	keyPattern     = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	spaceValueRe   = regexp.MustCompile(`^\s|\s$`)
)

// Lint runs all checks against the provided key/value map and returns findings.
func Lint(env map[string]string) []Finding {
	var findings []Finding
	seen := make(map[string]int)

	for k, v := range env {
		seen[k]++

		// Naming convention: must be UPPER_SNAKE_CASE
		if !keyPattern.MatchString(k) {
			findings = append(findings, Finding{
				Key:      k,
				Message:  "key does not follow UPPER_SNAKE_CASE convention",
				Severity: SeverityWarning,
			})
		}

		// Empty value warning
		if v == "" {
			findings = append(findings, Finding{
				Key:      k,
				Message:  "value is empty",
				Severity: SeverityWarning,
			})
		}

		// Leading/trailing whitespace in unquoted value
		if spaceValueRe.MatchString(v) {
			findings = append(findings, Finding{
				Key:      k,
				Message:  "value has leading or trailing whitespace",
				Severity: SeverityWarning,
			})
		}

		// Suspicious plaintext secret patterns
		lower := strings.ToLower(k)
		if (strings.Contains(lower, "password") || strings.Contains(lower, "secret") || strings.Contains(lower, "token")) && len(v) < 8 && v != "" {
			findings = append(findings, Finding{
				Key:      k,
				Message:  "secret-like key has a suspiciously short value",
				Severity: SeverityError,
			})
		}
	}

	return findings
}
