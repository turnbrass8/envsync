package validate

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ParseRules reads a simple rules file where each line is:
//
//	KEY [nonempty] [pattern=<regex>]
//
// Lines starting with '#' or blank lines are ignored.
func ParseRules(r io.Reader) ([]Rule, error) {
	var rules []Rule
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		rule := Rule{Key: parts[0]}

		for _, token := range parts[1:] {
			switch {
			case token == "nonempty":
				rule.NonEmpty = true
			case strings.HasPrefix(token, "pattern="):
				rule.Pattern = strings.TrimPrefix(token, "pattern=")
			default:
				return nil, fmt.Errorf("line %d: unknown token %q for key %q", lineNum, token, rule.Key)
			}
		}

		rules = append(rules, rule)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading rules: %w", err)
	}

	return rules, nil
}
