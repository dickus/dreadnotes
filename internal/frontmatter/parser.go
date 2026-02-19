package frontmatter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func SplitFile(note string) (path string, frontmatter []byte, content []byte) {
	homeDir, _ := os.UserHomeDir()
	notesDir := func(path string) string {
		if strings.HasPrefix(path, "$HOME") {
			return strings.Replace(path, "$HOME", homeDir, 1)
		}

		return path
	}

	file, err := os.Open(notesDir(note))

	if err != nil {
		fmt.Println("Note not found.")
	}
	defer file.Close()

	var frontBuilder strings.Builder
	var contentBuilder strings.Builder

	scanner := bufio.NewScanner(file)

	front := true
	separator := 0

	for scanner.Scan() {
		line := scanner.Text()

		if front {
			if line == "---" {
				separator++

				continue
			}

			if separator > 1 {
				front = false

				continue
			}

			frontBuilder.WriteString(line)
			frontBuilder.WriteString("\n")
		} else {
			contentBuilder.WriteString(line)
			contentBuilder.WriteString("\n")
		}
	}

	if err = scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return notesDir(note), []byte(frontBuilder.String()), []byte(contentBuilder.String())
}

func Parser(path string, frontmatter []byte, content []byte) Document {
	reader := strings.NewReader(string(frontmatter))
	scanner := bufio.NewScanner(reader)

	var (
		title   string
		created string
		updated string
		tags    []string
	)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "title:") {
			title = strings.Replace(line, "title: ", "", 1)

			continue
		}

		if strings.Contains(line, "created:") {
			created = strings.Replace(line, "created: ", "", 1)

			continue
		}

		if strings.Contains(line, "updated:") {
			updated = strings.Replace(line, "updated: ", "", 1)

			continue
		}

		if strings.Contains(line, "tags:") {
			prePreTags := strings.Replace(line, "tags: [", "", 1)
			preTags := strings.Replace(prePreTags, "]", "", 1)
			tags = strings.Split(preTags, ", ")

			continue
		}
	}

	return Document{
		Meta: Frontmatter{
			Title:   title,
			Created: created,
			Updated: updated,
			Tags:    tags,
		},
		Content: content,
		Path:    path,
	}
}
