package audit

import (
	"fmt"
	"io"
	"time"
)

// EventKind describes the type of audit event.
type EventKind string

const (
	EventApplied EventKind = "applied"
	EventSkipped EventKind = "skipped"
	EventMissing EventKind = "missing"
	EventInvalid EventKind = "invalid"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time
	Kind      EventKind
	Key       string
	Message   string
}

// String returns a human-readable representation of the event.
func (e Event) String() string {
	return fmt.Sprintf("%s [%s] key=%q %s",
		e.Timestamp.Format(time.RFC3339),
		e.Kind,
		e.Key,
		e.Message,
	)
}

// Logger records audit events and can write them to an io.Writer.
type Logger struct {
	events []Event
	out    io.Writer
}

// New creates a new Logger that writes to out.
func New(out io.Writer) *Logger {
	return &Logger{out: out}
}

// Record appends an event to the log and writes it to the underlying writer.
func (l *Logger) Record(kind EventKind, key, message string) {
	e := Event{
		Timestamp: time.Now(),
		Kind:      kind,
		Key:       key,
		Message:   message,
	}
	l.events = append(l.events, e)
	if l.out != nil {
		fmt.Fprintln(l.out, e.String())
	}
}

// Events returns all recorded events.
func (l *Logger) Events() []Event {
	return l.events
}

// Summary returns a count breakdown of events by kind.
func (l *Logger) Summary() map[EventKind]int {
	counts := make(map[EventKind]int)
	for _, e := range l.events {
		counts[e.Kind]++
	}
	return counts
}
