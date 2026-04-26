// Package interpolate resolves variable references within env values.
// Supports ${VAR} and $VAR syntax, with optional default via ${VAR:-default}.
package interpolate

import (
	"fmt"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Resolve replaces variable references in value using the provided env map.
// Returns an error if a required reference cannot be resolved.
func Resolve(value string, env map[string]string) (string, error) {
	var resolveErr error
	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}
		var key, defaultVal string
		hasDefault := false

		if strings.HasPrefix(match, "${") {
			inner := match[2 : len(match)-1]
			if idx := strings.Index(inner, ":-"); idx >= 0 {
				key = inner[:idx]
				defaultVal = inner[idx+2:]
				hasDefault = true
			} else {
				key = inner
			}
		} else {
			key = match[1:]
		}

		if v, ok := env[key]; ok {
			return v
		}
		if hasDefault {
			return defaultVal
		}
		resolveErr = fmt.Errorf("interpolate: unresolved variable %q", key)
		return match
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}

// ResolveAll applies Resolve to every value in the map in-place.
// References may refer to keys within the same map.
func ResolveAll(env map[string]string) error {
	for k, v := range env {
		resolved, err := Resolve(v, env)
		if err != nil {
			return fmt.Errorf("key %q: %w", k, err)
		}
		env[k] = resolved
	}
	return nil
}
