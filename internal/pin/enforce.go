package pin

import (
	"fmt"

	"github.com/example/envsync/internal/envfile"
)

// Violation describes a pinned key whose value has drifted.
type Violation struct {
	Key      string
	Pinned   string
	Actual   string
}

func (v Violation) String() string {
	return fmt.Sprintf("key %q: pinned=%q actual=%q", v.Key, v.Pinned, v.Actual)
}

// EnforceResult is returned by Enforce.
type EnforceResult struct {
	Violations []Violation
	OK         bool
}

// Enforce checks that every pinned key in the pin-file still matches the
// value in envFile. Returns violations for any drifted keys.
func Enforce(envFile string) (*EnforceResult, error) {
	pins, err := LoadPins(envFile)
	if err != nil {
		return nil, err
	}
	if len(pins) == 0 {
		return &EnforceResult{OK: true}, nil
	}

	env, err := envfile.Parse(envFile)
	if err != nil {
		return nil, fmt.Errorf("enforce: parse %s: %w", envFile, err)
	}

	result := &EnforceResult{OK: true}
	for k, pinned := range pins {
		actual, ok := env.Get(k)
		if !ok || actual != pinned {
			result.Violations = append(result.Violations, Violation{
				Key:    k,
				Pinned: pinned,
				Actual: actual,
			})
			result.OK = false
		}
	}
	return result, nil
}
