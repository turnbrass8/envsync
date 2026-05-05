package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envsync/internal/defaults"
	"github.com/yourorg/envsync/internal/envfile"
)

func runDefaults(args []string) error {
	fs := flag.NewFlagSet("defaults", flag.ContinueOnError)
	overwrite := fs.Bool("overwrite", false, "replace existing keys with default values")
	dryRun := fs.Bool("dry-run", false, "print what would change without writing")
	defsFile := fs.String("defaults", "", "path to .env file containing default values (required)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *defsFile == "" {
		return fmt.Errorf("--defaults flag is required")
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: envsync defaults --defaults <file> <target.env>")
	}
	targetPath := fs.Arg(0)

	targetEnv, err := envfile.Parse(targetPath)
	if err != nil {
		return fmt.Errorf("parsing target: %w", err)
	}

	defsEnv, err := envfile.Parse(*defsFile)
	if err != nil {
		return fmt.Errorf("parsing defaults file: %w", err)
	}

	env := make(map[string]string, len(targetEnv))
	for _, e := range targetEnv {
		env[e.Key] = e.Value
	}
	defsMap := make(map[string]string, len(defsEnv))
	for _, e := range defsEnv {
		defsMap[e.Key] = e.Value
	}

	res, err := defaults.Apply(env, defsMap, defaults.Options{
		Overwrite: *overwrite,
		DryRun:    *dryRun,
	})
	if err != nil {
		return err
	}

	if *dryRun {
		for _, k := range res.Applied {
			fmt.Fprintf(os.Stdout, "+ %s=%s\n", k, defsMap[k])
		}
		for _, k := range res.Overwritten {
			fmt.Fprintf(os.Stdout, "~ %s=%s\n", k, defsMap[k])
		}
		fmt.Fprintln(os.Stdout, res.Summary())
		return nil
	}

	if err := envfile.Write(targetPath, env); err != nil {
		return fmt.Errorf("writing target: %w", err)
	}
	fmt.Fprintln(os.Stdout, res.Summary())
	return nil
}
