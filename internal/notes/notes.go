package notes

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/dickus/dreadnotes/internal/frontmatter"
	"github.com/dickus/dreadnotes/internal/models"
)

func NewNote(name string, path string) {
	homeDir, _ := os.UserHomeDir()
	notesDir := func(path string) string {
		if strings.HasPrefix(path, "$HOME") {
			return strings.Replace(path, "$HOME", homeDir, 1)
		}

		return path
	}

	timestamp := time.Now().Unix()

	os.MkdirAll(notesDir(path), 0755)

	var filePath string

	if name == "" {
		filePath = notesDir(path) + "/" + strconv.FormatInt(timestamp, 10) + ".md"
	} else {
		filePath = notesDir(path) + "/" + strconv.FormatInt(timestamp, 10) + "_" + strings.ReplaceAll(name, " ", "_") + ".md"
	}

	file, err := os.Create(filePath)

	if err != nil {
		fmt.Println("Couldn't create note: ", err)
	}
	file.Close()

	frontmatter.CreateFrontmatter(filePath)

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

