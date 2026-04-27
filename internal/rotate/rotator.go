// Package rotate provides utilities for rotating secrets in .env files,
// generating new values and optionally backing up the originals.
package rotate

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Strategy defines how a new value is generated during rotation.
type Strategy string

const (
	StrategyRandom     Strategy = "random"
	StrategyUUID       Strategy = "uuid"
	StrategyAlphaNum   Strategy = "alphanum"
)

// Rule describes how a specific key should be rotated.
type Rule struct {
	Key      string
	Strategy Strategy
	Length   int
}

// Result holds the outcome of rotating a single key.
type Result struct {
	Key      string
	OldValue string
	NewValue string
	Rotated  bool
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

const alphanumChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const hexChars = "0123456789abcdef"

// Rotate applies the given rules to the env map, returning results for each key.
// It does not modify the original map; it returns a new map with rotated values.
func Rotate(env map[string]string, rules []Rule, dryRun bool) (map[string]string, []Result, error) {
	output := make(map[string]string, len(env))
	for k, v := range env {
		output[k] = v
	}

	var results []Result
	for _, rule := range rules {
		old, exists := env[rule.Key]
		if !exists {
			return nil, nil, fmt.Errorf("rotate: key %q not found in env", rule.Key)
		}

		newVal, err := generate(rule)
		if err != nil {
			return nil, nil, fmt.Errorf("rotate: generating value for %q: %w", rule.Key, err)
		}

		results = append(results, Result{
			Key:      rule.Key,
			OldValue: old,
			NewValue: newVal,
			Rotated:  !dryRun,
		})

		if !dryRun {
			output[rule.Key] = newVal
		}
	}

	return output, results, nil
}

func generate(rule Rule) (string, error) {
	length := rule.Length
	if length <= 0 {
		length = 32
	}

	switch rule.Strategy {
	case StrategyRandom, "":
		return randomHex(length), nil
	case StrategyAlphaNum:
		return randomAlphaNum(length), nil
	case StrategyUUID:
		return generateUUID(), nil
	default:
		return "", fmt.Errorf("unknown strategy %q", rule.Strategy)
	}
}

func randomHex(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(hexChars[rng.Intn(len(hexChars))])
	}
	return b.String()
}

func randomAlphaNum(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(alphanumChars[rng.Intn(len(alphanumChars))])
	}
	return b.String()
}

func generateUUID() string {
	var buf [16]byte
	_, _ = rng.Read(buf[:])
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:16])
}
