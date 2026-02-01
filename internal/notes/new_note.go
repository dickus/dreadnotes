package notes

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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
		filePath = notesDir(path) + "/" + strconv.FormatInt(timestamp, 10) + "_" + name + ".md"
	}

	file, err := os.Create(filePath)

	if err != nil {
		fmt.Println("Couldn't create note: ", err)
	}
	defer file.Close()
}

