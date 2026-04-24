package audit_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envsync/internal/audit"
)

func TestRecord_AppendsEvent(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	l.Record(audit.EventApplied, "DB_HOST", "default applied")

	events := l.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Kind != audit.EventApplied {
		t.Errorf("expected kind %q, got %q", audit.EventApplied, events[0].Kind)
	}
	if events[0].Key != "DB_HOST" {
		t.Errorf("expected key %q, got %q", "DB_HOST", events[0].Key)
	}
}

func TestRecord_WritesToOutput(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	l.Record(audit.EventMissing, "SECRET", "required key absent")

	output := buf.String()
	if !strings.Contains(output, "missing") {
		t.Errorf("expected output to contain 'missing', got: %s", output)
	}
	if !strings.Contains(output, "SECRET") {
		t.Errorf("expected output to contain key name, got: %s", output)
	}
}

func TestSummary_CountsByKind(t *testing.T) {
	l := audit.New(nil)

	l.Record(audit.EventApplied, "A", "")
	l.Record(audit.EventApplied, "B", "")
	l.Record(audit.EventSkipped, "C", "")
	l.Record(audit.EventMissing, "D", "")

	summary := l.Summary()

	if summary[audit.EventApplied] != 2 {
		t.Errorf("expected 2 applied, got %d", summary[audit.EventApplied])
	}
	if summary[audit.EventSkipped] != 1 {
		t.Errorf("expected 1 skipped, got %d", summary[audit.EventSkipped])
	}
	if summary[audit.EventMissing] != 1 {
		t.Errorf("expected 1 missing, got %d", summary[audit.EventMissing])
	}
}

func TestEvent_String_ContainsFields(t *testing.T) {
	l := audit.New(nil)
	l.Record(audit.EventInvalid, "PORT", "must be numeric")

	events := l.Events()
	s := events[0].String()

	for _, want := range []string{"invalid", "PORT", "must be numeric"} {
		if !strings.Contains(s, want) {
			t.Errorf("expected string to contain %q, got: %s", want, s)
		}
	}
}

func TestLogger_NilWriter_DoesNotPanic(t *testing.T) {
	l := audit.New(nil)
	l.Record(audit.EventApplied, "KEY", "no output writer")
	if len(l.Events()) != 1 {
		t.Error("expected event to be recorded even with nil writer")
	}
}
