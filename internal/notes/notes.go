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

	"github.com/dickus/dreadnotes/internal/frontmatter"
	"github.com/dickus/dreadnotes/internal/models"
	"github.com/dickus/dreadnotes/internal/templates"
)

func NewNote(name string, tmplPath string) {
	notesDir := PathParse(models.Cfg.NotesPath)

	fmt.Println(notesDir)

	timestamp := time.Now().Unix()

	os.MkdirAll(notesDir, 0755)

	var filePath string

	if name == "" {
		filePath = notesDir + "/" + strconv.FormatInt(timestamp, 10) + ".md"
	} else {
		filePath = notesDir + "/" + strconv.FormatInt(timestamp, 10) + "_" + strings.ReplaceAll(name, " ", "_") + ".md"
	}

	file, err := os.Create(filePath)

	if err != nil {
		fmt.Println("Couldn't create note: ", err)
	}
	file.Close()

	if tmplPath != "" {
		content, err := templates.ApplyTemplate(tmplPath, name)
		if err != nil {
			fmt.Println("Couldn't read template: ", err)

			return
		}

		os.WriteFile(filePath, content, 0644)
	} else {
		frontmatter.CreateFrontmatter(filePath, name)
	}

	OpenNote(filePath)
}

func OpenNote(file string) error {
	cmd := exec.Command(models.Cfg.Editor, file)

	if models.Cfg.Editor == "nvim" {
		lineNumber, err := nvimFindContent(file)

		if err != nil { return err }

		cmd = exec.Command(models.Cfg.Editor, fmt.Sprintf("+%d", lineNumber), file)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}

func nvimFindContent(file string) (int, error) {
	content, err := os.ReadFile(file)

	if err != nil { return 1, err }

	lines := strings.Split(string(content), "\n")

	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" { return i + 2, nil }
		}
	}

	return 1, nil
}

func RandomNote(path string) (string, error) {
	notesDir := PathParse(path)

	entries, err := os.ReadDir(notesDir)
	if err != nil {
		return "", fmt.Errorf("Failed to read notes directory: %w", err)
	}

	var notes []string
	for _, e := range entries {
		if e.IsDir() { continue }

		if strings.HasSuffix(e.Name(), ".md") {
			notes = append(notes, filepath.Join(notesDir, e.Name()))
		}
	}

	if len(notes) == 0 {
		return "", fmt.Errorf("No notes found in %s", notesDir)
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(notes))))
	if err != nil {
		return "", fmt.Errorf("Failed to generate random index: %w", err)
	}

	return notes[n.Int64()], nil
}

func PathParse(path string) string {
	homeDir, _ := os.UserHomeDir()
	if strings.HasPrefix(path, "$HOME") {
		return strings.Replace(path, "$HOME", homeDir, 1)
	}

	return path
}

