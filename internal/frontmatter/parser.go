package frontmatter

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/dickus/dreadnotes/internal/utils"
	"gopkg.in/yaml.v3"
)

// ParseFile reads a markdown note, extracts the YAML metadata, parses it into a Frontmatter struct, and returns the full Document.
func ParseFile(filePath string) (Document, error) {
	resolvedPath := utils.PathParse(filePath)

	file, err := os.Open(resolvedPath)
	if err != nil {
		return Document{}, fmt.Errorf("failed to open note: %w", err)
	}
	defer file.Close()

	var frontBuffer bytes.Buffer
	var contentBuffer bytes.Buffer

	scanner := bufio.NewScanner(file)

	isFrontmatter := false
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if lineCount == 1 && strings.TrimSpace(line) == "---" {
			isFrontmatter = true

			continue
		}

		if isFrontmatter && strings.TrimSpace(line) == "---" {
			isFrontmatter = false

			continue
		}

		if isFrontmatter {
			frontBuffer.WriteString(line + "\n")
		} else {
			contentBuffer.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return Document{}, fmt.Errorf("error reading file %s: %w", resolvedPath, err)
	}

	var meta Frontmatter
	frontBytes := frontBuffer.Bytes()

	if len(frontBytes) > 0 {
		if err := yaml.Unmarshal(frontBytes, &meta); err != nil {
			return Document{}, fmt.Errorf("failed to parse YAML in %s: %w", resolvedPath, err)
		}
	}

	return Document{
		Meta:    meta,
		Content: contentBuffer.Bytes(),
		Path:    resolvedPath,
	}, nil
}
