package format

import "os"

func readRaw(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func writeRaw(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}
