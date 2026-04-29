package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/example/envsync/internal/pin"
)

func runPin(args []string) error {
	fs := flag.NewFlagSet("pin", flag.ContinueOnError)
	envFile := fs.String("env", ".env", "path to the .env file")
	keys := fs.String("keys", "", "comma-separated list of keys to pin")
	dryRun := fs.Bool("dry-run", false, "preview without writing the pin file")
	enforce := fs.Bool("enforce", false, "check that pinned keys have not drifted")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *enforce {
		res, err := pin.Enforce(*envFile)
		if err != nil {
			return err
		}
		if res.OK {
			fmt.Println("pin: all pinned keys match — no drift detected")
			return nil
		}
		fmt.Fprintf(os.Stderr, "pin: %d violation(s) detected:\n", len(res.Violations))
		for _, v := range res.Violations {
			fmt.Fprintf(os.Stderr, "  %s\n", v)
		}
		return fmt.Errorf("pin: drift detected in %d key(s)", len(res.Violations))
	}

	if *keys == "" {
		return fmt.Errorf("pin: --keys is required")
	}
	keyList := strings.Split(*keys, ",")
	for i, k := range keyList {
		keyList[i] = strings.TrimSpace(k)
	}

	res, err := pin.Pin(*envFile, keyList, pin.Options{DryRun: *dryRun})
	if err != nil {
		return err
	}

	for _, p := range res.Pinned {
		fmt.Printf("pinned: %s\n", p.Key)
	}
	for _, s := range res.Skipped {
		fmt.Printf("skipped (not found): %s\n", s)
	}
	if *dryRun {
		fmt.Println("(dry-run: pin file not written)")
	}
	return nil
}
