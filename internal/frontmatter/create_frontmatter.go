package frontmatter

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func CreateFrontmatter(path string, name string) {
	var sb strings.Builder

	tmpTime := time.Now()
	creationTime := tmpTime.Format("2006-01-02 15:04")

	fmt.Fprintf(&sb, "---\n")
	fmt.Fprintf(&sb, "title: %v\n", name)
	fmt.Fprintf(&sb, "created: %v\n", creationTime)
	fmt.Fprintf(&sb, "updated: %v\n", creationTime)
	fmt.Fprintf(&sb, "tags: \n")
	fmt.Fprintf(&sb, "---\n")
	fmt.Fprintf(&sb, "\n")

	content := sb.String()

	err := os.WriteFile(path, []byte(content), 0644)

	if err != nil {
		fmt.Println("Couldn't write frontmatter: ", err)
	}
}

