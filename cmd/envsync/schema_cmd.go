package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envsync/internal/envfile"
	"github.com/yourorg/envsync/internal/schema"
)

// schemaFieldDef is the JSON representation of a schema field.
type schemaFieldDef struct {
	Key      string `json:"key"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Pattern  string `json:"pattern,omitempty"`
}

func runSchema(args []string) error {
	fs := flag.NewFlagSet("schema", flag.ContinueOnError)
	envPath := fs.String("env", ".env", "path to .env file")
	schemaPath := fs.String("schema", "schema.json", "path to JSON schema definition")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Load schema definition
	schemaData, err := os.ReadFile(*schemaPath)
	if err != nil {
		return fmt.Errorf("reading schema file: %w", err)
	}

	var defs []schemaFieldDef
	if err := json.Unmarshal(schemaData, &defs); err != nil {
		return fmt.Errorf("parsing schema JSON: %w", err)
	}

	s := &schema.Schema{}
	for _, d := range defs {
		s.Fields = append(s.Fields, schema.Field{
			Key:      d.Key,
			Type:     schema.FieldType(d.Type),
			Required: d.Required,
			Pattern:  d.Pattern,
		})
	}

	// Load env file
	env, err := envfile.Parse(*envPath)
	if err != nil {
		return fmt.Errorf("parsing env file: %w", err)
	}

	violations := s.Validate(env)
	if len(violations) == 0 {
		fmt.Println("schema: all checks passed")
		return nil
	}

	fmt.Fprintf(os.Stderr, "schema violations (%d):\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(os.Stderr, "  [FAIL] %s\n", v.Error())
	}
	return fmt.Errorf("%d schema violation(s) found", len(violations))
}
