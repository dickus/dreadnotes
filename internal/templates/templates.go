package templates

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
	"time"
)

type TemplateData struct {
	Title 	string
	Date 	string
	Year 	string
}

func ApplyTemplate(tmplPath string, name string) ([]byte, error) {
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read template: %w", err)
	}

	now := time.Now()

	data := TemplateData{
		Title: 	name,
		Date: 	now.Format("2006-01-02 15:04"),
		Year: 	now.Format("2006"),
	}

	tmpl, err := template.New("note").Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("Couldn't execute template: %w", err)
	}

	return buf.Bytes(), nil
}

