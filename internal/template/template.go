// Package template provides functionality to render .env files from
// Go text/template sources, injecting values from an existing env map.
package template

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	gotemplate "text/template"
)

// RenderOptions controls how a template is rendered.
type RenderOptions struct {
	// Strict causes Render to return an error if any template variable is missing
	// from the provided env map.
	Strict bool
}

// Render parses src as a Go text/template and executes it with the provided
// env map as the data context. Variables are referenced as {{ .KEY }}.
// If opts.Strict is true, missing keys produce an error instead of an empty
// string.
func Render(src string, env map[string]string, opts RenderOptions) (string, error) {
	missingOption := "zero"
	if opts.Strict {
		missingOption = "error"
	}

	tmpl, err := gotemplate.New("env").
		Option("missingkey=" + missingOption).
		Funcs(gotemplate.FuncMap{
			"default": func(def, val string) string {
				if val == "" {
					return def
				}
				return val
			},
			"upper": strings.ToUpper,
			"lower": strings.ToLower,
		}).
		Parse(src)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, env); err != nil {
		return "", fmt.Errorf("template execute error: %w", err)
	}
	return buf.String(), nil
}

// RenderFile reads a template from path and renders it using Render.
func RenderFile(path string, env map[string]string, opts RenderOptions) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read template file: %w", err)
	}
	return Render(string(data), env, opts)
}
