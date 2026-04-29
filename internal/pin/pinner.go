// Package pin provides functionality to pin (lock) specific env keys to
// their current values, preventing accidental overwrites during sync.
package pin

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/example/envsync/internal/envfile"
)

// PinnedKey represents a single pinned key and its locked value.
type PinnedKey struct {
	Key   string
	Value string
}

// Result holds the outcome of a pin operation.
type Result struct {
	Pinned  []PinnedKey
	Skipped []string // keys not found in the env file
}

// Options controls Pin behaviour.
type Options struct {
	DryRun bool
}

// Pin locks the given keys in envFile to their current values by writing a
// pin-file (envFile + ".pinned"). Returns an error if the env file cannot be
// read or written.
func Pin(envFile string, keys []string, opts Options) (*Result, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("pin: no keys specified")
	}

	env, err := envfile.Parse(envFile)
	if err != nil {
		return nil, fmt.Errorf("pin: parse %s: %w", envFile, err)
	}

	result := &Result{}
	pinMap := make(map[string]string)

	for _, k := range keys {
		v, ok := env.Get(k)
		if !ok {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		result.Pinned = append(result.Pinned, PinnedKey{Key: k, Value: v})
		pinMap[k] = v
	}

	if opts.DryRun || len(pinMap) == 0 {
		return result, nil
	}

	return result, writePinFile(envFile+".pinned", pinMap)
}

// LoadPins reads the pin-file associated with envFile and returns the locked
// key→value map. Returns an empty map if no pin-file exists.
func LoadPins(envFile string) (map[string]string, error) {
	path := envFile + ".pinned"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	env, err := envfile.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("pin: load pins %s: %w", path, err)
	}
	out := make(map[string]string)
	for _, k := range env.Keys() {
		v, _ := env.Get(k)
		out[k] = v
	}
	return out, nil
}

func writePinFile(path string, pins map[string]string) error {
	keys := make([]string, 0, len(pins))
	for k := range pins {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString("# envsync pin file — do not edit manually\n")
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, pins[k])
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
