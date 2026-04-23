package envfile

import (
	"strings"
	"testing"
)

func TestParse_BasicKeyValue(t *testing.T) {
	input := `APP_ENV=production
DB_HOST=localhost
DB_PORT=5432
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(ef.Entries))
	}

	cases := []struct{ key, want string }{
		{"APP_ENV", "production"},
		{"DB_HOST", "localhost"},
		{"DB_PORT", "5432"},
	}
	for _, c := range cases {
		val, ok := ef.Get(c.key)
		if !ok {
			t.Errorf("key %q not found", c.key)
		}
		if val != c.want {
			t.Errorf("key %q: got %q, want %q", c.key, val, c.want)
		}
	}
}

func TestParse_StripsQuotes(t *testing.T) {
	input := `SECRET="my secret value"
TOKEN='abc123'
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := ef.Get("SECRET"); v != "my secret value" {
		t.Errorf("expected unquoted value, got %q", v)
	}
	if v, _ := ef.Get("TOKEN"); v != "abc123" {
		t.Errorf("expected unquoted value, got %q", v)
	}
}

func TestParse_SkipsCommentsAndBlanks(t *testing.T) {
	input := `
# This is a comment
FOO=bar

# Another comment
BAZ=qux
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(ef.Entries))
	}
}

func TestParse_InvalidLine(t *testing.T) {
	input := `INVALID_LINE_NO_EQUALS
`
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestGet_MissingKey(t *testing.T) {
	ef := &EnvFile{Index: make(map[string]int)}
	_, ok := ef.Get("NONEXISTENT")
	if ok {
		t.Error("expected false for missing key")
	}
}
