package patch

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ParseRules reads patch operations from a reader.
// Each non-blank, non-comment line must be one of:
//
//	KEY=VALUE   – set/update
//	-KEY        – delete
func ParseRules(r io.Reader) ([]Op, error) {
	var ops []Op
	scanner := bufio.NewScanner(r)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "-") {
			key := strings.TrimSpace(line[1:])
			if key == "" {
				return nil, fmt.Errorf("patch rules line %d: empty key after '-'", lineNo)
			}
			ops = append(ops, Op{Key: key, Delete: true})
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("patch rules line %d: expected KEY=VALUE or -KEY, got %q", lineNo, line)
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		if key == "" {
			return nil, fmt.Errorf("patch rules line %d: empty key", lineNo)
		}
		ops = append(ops, Op{Key: key, Value: value})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("patch rules: scan error: %w", err)
	}
	return ops, nil
}
