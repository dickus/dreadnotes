package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// PathParse expands "~" and "$HOME" to the absolute path of the user's home directory.
// It returns the cleaned absolute path, or the original path if no expansion is needed.
func PathParse(path string) string {
	if path == "" {
		return path
	}

	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}

		return filepath.Join(homeDir, path[1:])
	}

	if strings.HasPrefix(path, "$HOME") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}

		return filepath.Join(homeDir, path[5:])
	}

	return filepath.Clean(path)
}
