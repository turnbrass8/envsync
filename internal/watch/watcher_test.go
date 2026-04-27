package watch_test

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/yourorg/envsync/internal/watch"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestWatch_DetectsModification(t *testing.T) {
	path := writeTempEnv(t, "KEY=original\n")

	var mu sync.Mutex
	var events []watch.Event

	w := watch.New(20*time.Millisecond, func(e watch.Event) {
		mu.Lock()
		events = append(events, e)
		mu.Unlock()
	})

	errCh := make(chan error, 1)
	go func() { errCh <- w.Watch([]string{path}) }()

	time.Sleep(40 * time.Millisecond)

	// Modify the file.
	if err := os.WriteFile(path, []byte("KEY=updated\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	time.Sleep(60 * time.Millisecond)
	w.Stop()

	if err := <-errCh; err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(events) == 0 {
		t.Fatal("expected at least one event, got none")
	}
	if events[0].Path != path {
		t.Errorf("event path = %q, want %q", events[0].Path, path)
	}
}

func TestWatch_NoPaths_ReturnsError(t *testing.T) {
	w := watch.New(10*time.Millisecond, func(watch.Event) {})
	if err := w.Watch(nil); err == nil {
		t.Fatal("expected error for empty paths, got nil")
	}
}

func TestWatch_MissingFile_ReturnsError(t *testing.T) {
	w := watch.New(10*time.Millisecond, func(watch.Event) {})
	if err := w.Watch([]string{"/nonexistent/path/.env"}); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestWatch_StopWithoutChange_NoEvents(t *testing.T) {
	path := writeTempEnv(t, "KEY=value\n")

	var mu sync.Mutex
	var events []watch.Event

	w := watch.New(50*time.Millisecond, func(e watch.Event) {
		mu.Lock()
		events = append(events, e)
		mu.Unlock()
	})

	errCh := make(chan error, 1)
	go func() { errCh <- w.Watch([]string{path}) }()

	time.Sleep(20 * time.Millisecond)
	w.Stop()
	<-errCh

	mu.Lock()
	defer mu.Unlock()
	if len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}
