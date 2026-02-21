// Package templates provides handling templates user can specify.
package templates

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
	"time"
)

// TemplateData holds the variables that will be injected into the markdown templates.
type TemplateData struct {
	Title string
	Date  string
}

// ApplyTemplate reads a markdown template from the given path, injects the note's name and current date, and returns the rendered bytes.
func ApplyTemplate(tmplPath string, name string) ([]byte, error) {
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read template: %w", err)
	}

	data := TemplateData{
		Title: name,
		Date:  time.Now().Format("2006-01-02 15:04"),
	}

	tmpl, err := template.New("note").Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("couldn't parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("couldn't execute template: %w", err)
	}

	return buf.Bytes(), nil
}
