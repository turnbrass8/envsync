// Package env provides multi-source environment variable loading.
//
// It supports merging variables from OS environment, .env files, and
// in-memory maps with configurable priority ordering. Later sources
// override earlier ones, mirroring standard 12-factor app conventions.
//
// Example usage:
//
//	fileSrc, err := env.FileSource("dotenv", ".env")
//	if err != nil { ... }
//
//	loader := env.NewLoader(fileSrc, env.OSSource())
//	resolved := loader.Resolve()
//	// OS env takes precedence over .env file
//
// Origin tracking lets callers inspect which source last defined a key,
// useful for debugging environment configuration issues.
package env
