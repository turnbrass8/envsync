package flatten_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/flatten"
)

func TestFlatten_SimpleObject(t *testing.T) {
	input := []byte(`{"host":"localhost","port":5432}`)
	got, err := flatten.Flatten(input, flatten.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"HOST=localhost", "PORT=5432"}
	assertEqual(t, want, got)
}

func TestFlatten_NestedObject(t *testing.T) {
	input := []byte(`{"db":{"host":"localhost","port":5432}}`)
	got, err := flatten.Flatten(input, flatten.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"DB_HOST=localhost", "DB_PORT=5432"}
	assertEqual(t, want, got)
}

func TestFlatten_WithPrefix(t *testing.T) {
	opts := flatten.DefaultOptions()
	opts.Prefix = "APP"
	input := []byte(`{"debug":true}`)
	got, err := flatten.Flatten(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"APP_DEBUG=true"}
	assertEqual(t, want, got)
}

func TestFlatten_ArrayValues(t *testing.T) {
	input := []byte(`{"hosts":["a","b"]}`)
	got, err := flatten.Flatten(input, flatten.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"HOSTS_0=a", "HOSTS_1=b"}
	assertEqual(t, want, got)
}

func TestFlatten_NullValue(t *testing.T) {
	input := []byte(`{"token":null}`)
	got, err := flatten.Flatten(input, flatten.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"TOKEN="}
	assertEqual(t, want, got)
}

func TestFlatten_CustomSeparator(t *testing.T) {
	opts := flatten.DefaultOptions()
	opts.Separator = "__"
	input := []byte(`{"a":{"b":"v"}}`)
	got, err := flatten.Flatten(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"A__B=v"}
	assertEqual(t, want, got)
}

func TestFlatten_InvalidJSON_ReturnsError(t *testing.T) {
	_, err := flatten.Flatten([]byte(`not json`), flatten.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func assertEqual(t *testing.T, want, got []string) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("length mismatch: want %d got %d\nwant: %v\ngot:  %v", len(want), len(got), want, got)
	}
	for i := range want {
		if want[i] != got[i] {
			t.Errorf("index %d: want %q got %q", i, want[i], got[i])
		}
	}
}
