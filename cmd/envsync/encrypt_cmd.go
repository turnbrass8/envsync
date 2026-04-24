package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envsync/internal/encrypt"
)

// runEncrypt handles the `envsync encrypt` and `envsync decrypt` sub-commands.
// Usage:
//
//	envsync encrypt -pass <passphrase> <value>
//	envsync decrypt -pass <passphrase> <ciphertext>
func runEncrypt(args []string, mode string) error {
	fs := flag.NewFlagSet("envsync "+mode, flag.ContinueOnError)
	pass := fs.String("pass", "", "Passphrase used for AES-GCM encryption (required)")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envsync %s -pass <passphrase> <value>\n", mode)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *pass == "" {
		fs.Usage()
		return fmt.Errorf("flag -pass is required")
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return fmt.Errorf("expected exactly one positional argument, got %d", fs.NArg())
	}
	input := fs.Arg(0)

	switch mode {
	case "encrypt":
		ct, err := encrypt.Encrypt(input, *pass)
		if err != nil {
			return fmt.Errorf("encrypt: %w", err)
		}
		fmt.Println(ct)
	case "decrypt":
		plain, err := encrypt.Decrypt(input, *pass)
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}
		fmt.Println(plain)
	default:
		return fmt.Errorf("unknown mode: %s", mode)
	}
	return nil
}
