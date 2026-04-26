// Package interpolate provides variable reference resolution for .env files.
//
// Supported syntax:
//
//	${VAR}          — replaced with the value of VAR; error if absent
//	${VAR:-default} — replaced with the value of VAR, or "default" if absent
//	$VAR            — shorthand for ${VAR} (no default support)
//
// Example:
//
//	env := map[string]string{
//		"HOST": "db.internal",
//		"DSN":  "postgres://${HOST}:5432/mydb",
//	}
//	if err := interpolate.ResolveAll(env); err != nil {
//		log.Fatal(err)
//	}
//	// env["DSN"] == "postgres://db.internal:5432/mydb"
package interpolate
