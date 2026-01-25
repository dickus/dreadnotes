package utils

import (
	"fmt"
	"os"
	"errors"
	"path/filepath"
)

func ConfigChecks() (string, string) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		fmt.Println("Couldn't find config directory: ", err)
	}

	configPath := filepath.Join(configDir, "dreadnotes", "config.toml")

	if _, err := os.Stat(configPath); err == nil {
		fmt.Println("Config is in place.")

		return "", ""
	} else if errors.Is(err, os.ErrNotExist) {
		return "$HOME/Documents/notes", "nvim"
	}

	return "", ""
}

