// Package watch provides file watching for .env files, triggering
// a callback whenever the file is modified on disk.
package watch

import (
	"fmt"
	"os"
	"time"
)

// Event describes a change detected in a watched file.
type Event struct {
	Path    string
	ModTime time.Time
}

// Handler is called when a watched file changes.
type Handler func(Event)

// Watcher polls a set of files at a fixed interval and invokes Handler
// on any file whose modification time has advanced.
type Watcher struct {
	interval time.Duration
	handler  Handler
	stop     chan struct{}
}

// New creates a Watcher that polls at the given interval.
func New(interval time.Duration, handler Handler) *Watcher {
	return &Watcher{
		interval: interval,
		handler:  handler,
		stop:     make(chan struct{}),
	}
}

// Watch begins watching the provided paths. It blocks until Stop is called.
func (w *Watcher) Watch(paths []string) error {
	if len(paths) == 0 {
		return fmt.Errorf("watch: no paths provided")
	}

	modTimes := make(map[string]time.Time, len(paths))
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return fmt.Errorf("watch: cannot stat %q: %w", p, err)
		}
		modTimes[p] = info.ModTime()
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-w.stop:
			return nil
		case <-ticker.C:
			for _, p := range paths {
				info, err := os.Stat(p)
				if err != nil {
					continue
				}
				if info.ModTime().After(modTimes[p]) {
					modTimes[p] = info.ModTime()
					w.handler(Event{Path: p, ModTime: info.ModTime()})
				}
			}
		}
	}
}

// Stop signals the watcher to cease polling.
func (w *Watcher) Stop() {
	close(w.stop)
}
