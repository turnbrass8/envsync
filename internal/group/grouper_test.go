package group_test

import (
	"testing"

	"github.com/user/envsync/internal/group"
)

func TestGroupBy_BasicPrefixes(t *testing.T) {
	env := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"REDIS_HOST":  "redis",
		"APP_NAME":    "myapp",
	}
	prefixes := map[string]string{
		"DB":    "database",
		"REDIS": "cache",
		"APP":   "app",
	}
	groups, err := group.GroupBy(env, prefixes, group.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
	for _, g := range groups {
		switch g.Name {
		case "database":
			if g.Entries["DB_HOST"] != "localhost" {
				t.Errorf("expected DB_HOST=localhost, got %q", g.Entries["DB_HOST"])
			}
		case "cache":
			if g.Entries["REDIS_HOST"] != "redis" {
				t.Errorf("expected REDIS_HOST=redis")
			}
		case "app":
			if g.Entries["APP_NAME"] != "myapp" {
				t.Errorf("expected APP_NAME=myapp")
			}
		}
	}
}

func TestGroupBy_StripPrefix(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	prefixes := map[string]string{"DB": "database"}
	groups, err := group.GroupBy(env, prefixes, group.Options{StripPrefix: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := groups[0].Entries["HOST"]; !ok {
		t.Errorf("expected stripped key HOST, got %v", groups[0].Entries)
	}
}

func TestGroupBy_IncludeUnmatched(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"UNKNOWN": "value",
	}
	prefixes := map[string]string{"DB": "database"}
	groups, err := group.GroupBy(env, prefixes, group.Options{IncludeUnmatched: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var other *group.Group
	for _, g := range groups {
		if g.Name == "_other" {
			other = g
		}
	}
	if other == nil {
		t.Fatal("expected _other group")
	}
	if other.Entries["UNKNOWN"] != "value" {
		t.Errorf("expected UNKNOWN in _other group")
	}
}

func TestGroupBy_NoPrefixes_ReturnsError(t *testing.T) {
	_, err := group.GroupBy(map[string]string{"A": "1"}, nil, group.Options{})
	if err == nil {
		t.Fatal("expected error for empty prefixes")
	}
}

func TestGroupBy_EmptyEnv_ReturnsEmptyGroups(t *testing.T) {
	groups, err := group.GroupBy(map[string]string{}, map[string]string{"DB": "database"}, group.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups[0].Entries) != 0 {
		t.Errorf("expected empty entries, got %v", groups[0].Entries)
	}
}
