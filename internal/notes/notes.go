// Package notes handles the creation, opening, and retrieval of markdown notes.
package notes

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dickus/dreadnotes/internal/config"
	"github.com/dickus/dreadnotes/internal/frontmatter"
	"github.com/dickus/dreadnotes/internal/templates"
	"github.com/dickus/dreadnotes/internal/utils"
)

// NewNote creates a new note file and opens it in the configured editor.
// It returns an error if any step (creation, template application, or opening) fails.
func NewNote(name string, tmplPath string) error {
	notesDir := utils.PathParse(config.Cfg.NotesPath)

	// Ensure the notes directory exists
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return fmt.Errorf("failed to create notes directory: %w", err)
	}

	timestamp := time.Now().Unix()
	var filename string

	// Generate filename: timestamp.md or timestamp_Note_Name.md
	if name == "" {
		filename = strconv.FormatInt(timestamp, 10) + ".md"
	} else {
		// Replace spaces with underscores for a cleaner filename
		cleanName := strings.ReplaceAll(name, " ", "_")
		filename = strconv.FormatInt(timestamp, 10) + "_" + cleanName + ".md"
	}

	filePath := filepath.Join(notesDir, filename)

	// Create the note content based on whether a template is used
	if tmplPath != "" {
		content, err := templates.ApplyTemplate(tmplPath, name)
		if err != nil {
			return fmt.Errorf("failed to apply template: %w", err)
		}

		// WriteFile creates the file if it doesn't exist, or truncates it if it does.
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			return fmt.Errorf("failed to write note file: %w", err)
		}
	} else {
		// Generate default frontmatter if no template is provided
		if err := frontmatter.Create(filePath, name); err != nil {
			return fmt.Errorf("failed to create default note: %w", err)
		}
	}

	// Open the newly created note in the editor
	return OpenNote(filePath)
}

// OpenNote opens the specified file in the configured editor.
// It handles special logic for Neovim to jump to the content line.
func OpenNote(file string) error {
	editor := config.Cfg.Editor
	var args []string

	// Special handling for Neovim: try to position the cursor after the frontmatter
	if editor == "nvim" {
		lineNumber, err := nvimFindContent(file)
		// Only add the line number argument if we successfully found a valid line
		if err == nil && lineNumber > 1 {
			args = append(args, fmt.Sprintf("+%d", lineNumber))
		}
	}

	// Append the file path as the last argument
	args = append(args, file)

	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// nvimFindContent reads the file and returns the line number where the content starts (immediately after the second "---" separator of the YAML frontmatter).
func nvimFindContent(file string) (int, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return 1, err
	}

	lines := strings.Split(string(content), "\n")
	separatorCount := 0

	// Check if the file starts with frontmatter delimiter
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		for i, line := range lines {
			if strings.TrimSpace(line) == "---" {
				separatorCount++
				// If we found the second "---", the content starts on the next line.
				// i is 0-indexed, so line number is i+1. Next line is i+2.
				if separatorCount == 2 {
					return i + 2, nil
				}
			}
		}
	}

	// Default to the first line if no valid frontmatter is found
	return 1, nil
}

// RandomNote selects a random markdown file from the notes directory.
func RandomNote(path string) (string, error) {
	notesDir := utils.PathParse(path)

	entries, err := os.ReadDir(notesDir)
	if err != nil {
		return "", fmt.Errorf("failed to read notes directory: %w", err)
	}

	var notes []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			notes = append(notes, filepath.Join(notesDir, e.Name()))
		}
	}

	if len(notes) == 0 {
		return "", fmt.Errorf("no notes found in %s", notesDir)
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(notes))))
	if err != nil {
		return "", fmt.Errorf("failed to generate random index: %w", err)
	}

	return notes[n.Int64()], nil
}
