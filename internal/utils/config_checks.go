package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	configDir, err = os.UserConfigDir()
	configPath = filepath.Join(configDir, "dreadnotes", "config.toml")
)

func ConfigInPlace() bool {
	if err != nil {
		fmt.Println("Couldn't find config directory: ", err)

		return false
	}

	if _, err := os.Stat(configPath); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return false
}

func ReadConfig() []string {
	file, _ := os.Open(configPath)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var configStrings []string

	for scanner.Scan() {
		configStrings = append(configStrings, scanner.Text())
	}

	return configStrings
}

func SplitConfig(data string) (string, string) {
	parts := strings.Split(data, " = ")

	return parts[0], parts[1]
}

func DataValidation(key string, value string) bool {
	if key == "notes_path" || key == "editor" || key == "templates_path" {
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			return true
		} else {
			return false
		}
	}

	return false
}

