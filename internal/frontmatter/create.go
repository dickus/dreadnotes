// Package frontmatter provides functionality to generate YAML frontmatter headers for Markdown files.
package frontmatter

import (
	"fmt"
	"os"
	"time"
)

// Create generates a YAML frontmatter block with title, creation/update timestamps, and writes it to the specified file path.
func Create(path string, name string) error {
	timestamp := time.Now().Format("2006-01-02 15:04")

	content := fmt.Sprintf(`---
title: %q
created: %s
updated: %s
tags: []
---

`, name, timestamp, timestamp)

	return os.WriteFile(path, []byte(content), 0644)
}
