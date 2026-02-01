package utils

import (
	"os"
	"fmt"
	"path/filepath"
)

func SaveConfig(content string) {
	configDir, _ := os.UserConfigDir()
	configPath := filepath.Join(configDir, "dreadnotes")

	err := os.MkdirAll(configPath, 0755)

	if err != nil {
		fmt.Printf("Couldn't create config directory: %v\n", err)

		return
	}

	fullPath := filepath.Join(configPath, "config.toml")

	err = os.WriteFile(fullPath, []byte(content), 0644)
}
