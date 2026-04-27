package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/envsync/internal/watch"
)

func runWatch(args []string) error {
	fs := flag.NewFlagSet("watch", flag.ContinueOnError)
	intervalMs := fs.Int("interval", 500, "polling interval in milliseconds")

	if err := fs.Parse(args); err != nil {
		return err
	}

	paths := fs.Args()
	if len(paths) == 0 {
		return fmt.Errorf("watch: at least one .env file path required")
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			return fmt.Errorf("watch: cannot access %q: %w", p, err)
		}
	}

	interval := time.Duration(*intervalMs) * time.Millisecond
	fmt.Fprintf(os.Stdout, "Watching %v (interval: %s) — press Ctrl+C to stop\n", paths, interval)

	w := watch.New(interval, func(e watch.Event) {
		fmt.Fprintf(os.Stdout, "[%s] changed: %s\n",
			e.ModTime.Format("15:04:05"), e.Path)
	})

	return w.Watch(paths)
}
