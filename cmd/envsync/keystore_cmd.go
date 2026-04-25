package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envsync/internal/encrypt"
)

func runKeystore(args []string) error {
	fs := flag.NewFlagSet("keystore", flag.ContinueOnError)
	ksPath := fs.String("keystore", ".envsync-keys.json", "path to keystore file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	sub := fs.Arg(0)
	if sub == "" {
		return fmt.Errorf("usage: envsync keystore <list|add|remove> [flags]")
	}

	ks, err := encrypt.LoadKeystore(*ksPath)
	if err != nil {
		return fmt.Errorf("loading keystore: %w", err)
	}

	switch sub {
	case "list":
		if len(ks.Entries) == 0 {
			fmt.Println("no key aliases defined")
			return nil
		}
		fmt.Printf("%-20s  %-30s  %s\n", "ALIAS", "HINT", "CREATED")
		for _, e := range ks.Entries {
			fmt.Printf("%-20s  %-30s  %s\n", e.Alias, e.Hint, e.CreatedAt.Format("2006-01-02"))
		}

	case "add":
		alias := fs.Arg(1)
		hint := fs.Arg(2)
		if alias == "" {
			return fmt.Errorf("usage: envsync keystore add <alias> [hint]")
		}
		if err := ks.Add(alias, hint); err != nil {
			return err
		}
		if err := ks.Save(); err != nil {
			return fmt.Errorf("saving keystore: %w", err)
		}
		fmt.Fprintf(os.Stdout, "added alias %q to %s\n", alias, *ksPath)

	case "remove":
		alias := fs.Arg(1)
		if alias == "" {
			return fmt.Errorf("usage: envsync keystore remove <alias>")
		}
		if err := ks.Remove(alias); err != nil {
			return err
		}
		if err := ks.Save(); err != nil {
			return fmt.Errorf("saving keystore: %w", err)
		}
		fmt.Fprintf(os.Stdout, "removed alias %q from %s\n", alias, *ksPath)

	default:
		return fmt.Errorf("unknown keystore subcommand: %q", sub)
	}

	return nil
}
