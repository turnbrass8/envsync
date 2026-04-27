package rotate

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParseRules reads rotation rules from a simple line-based format:
//
//	# comment
//	KEY strategy [length]
//
// Example:
//	SECRET_KEY   random 32
//	API_TOKEN    uuid
//	SESSION_ID   alphanum 24
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
		if len(parts) < 2 {
			return nil, fmt.Errorf("rotate rules line %d: expected KEY STRATEGY [LENGTH], got %q", lineNum, line)
		}

		rule := Rule{
			Key:      parts[0],
			Strategy: Strategy(strings.ToLower(parts[1])),
		}

		if len(parts) >= 3 {
			n, err := strconv.Atoi(parts[2])
			if err != nil {
				return nil, fmt.Errorf("rotate rules line %d: invalid length %q: %w", lineNum, parts[2], err)
			}
			rule.Length = n
		}

		rules = append(rules, rule)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("rotate rules: scanning: %w", err)
	}

	return rules, nil
}
