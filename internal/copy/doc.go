// Package copy provides key-selective copying between .env files.
//
// Usage:
//
//	res, err := copy.Copy("staging.env", "production.env", copy.Options{
//		Keys:      []string{"DB_HOST", "DB_PORT"},
//		Overwrite: false,
//		DryRun:    false,
//	})
//
// When Keys is empty all keys from the source file are candidates.
// If Overwrite is false, keys that already exist in the destination are
// reported in Result.Skipped and left unchanged.
// DryRun reports what would be copied without modifying the destination file.
package copy
