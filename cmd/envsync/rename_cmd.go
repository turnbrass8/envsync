package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsync/internal/rename"
)

func runRename(args []string) error {
	fs := flag.NewFlagSet("rename", flag.ContinueOnError)
	envFile := fs.String("env", ".env", "path to the .env file")
	dryRun := fs.Bool("dry-run", false, "print changes without writing")

	if err := fs.Parse(args); err != nil {
		return err
	}

	pairs := fs.Args()
	if len(pairs) == 0 || len(pairs)%2 != 0 {
		return fmt.Errorf("rename: provide pairs of OLD_KEY NEW_KEY (got %d args)", len(pairs))
	}

	var rules []rename.Rule
	for i := 0; i < len(pairs); i += 2 {
		rules = append(rules, rename.Rule{From: pairs[i], To: pairs[i+1]})
	}

	res, err := rename.Rename(*envFile, rules, *dryRun)
	if err != nil {
		return err
	}

	for _, r := range res.Applied {
		if *dryRun {
			fmt.Fprintf(os.Stdout, "[dry-run] would rename %s -> %s\n", r.From, r.To)
		} else {
			fmt.Fprintf(os.Stdout, "renamed %s -> %s\n", r.From, r.To)
		}
	}

	if len(res.Skipped) > 0 {
		skipped := make([]string, len(res.Skipped))
		for i, r := range res.Skipped {
			skipped[i] = r.From
		}
		fmt.Fprintf(os.Stdout, "skipped (not found): %s\n", strings.Join(skipped, ", "))
	}

	return nil
}
