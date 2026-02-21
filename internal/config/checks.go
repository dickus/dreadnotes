package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dickus/dreadnotes/internal/utils"
)

// getConfigPath safely retrieves the path to the configuration file.
// It prioritizes the DREADNOTES_CONFIG environment variable if set.
func getConfigPath() (string, error) {
	if envPath := os.Getenv("DREADNOTES_CONFIG"); envPath != "" {
		return utils.PathParse(envPath), nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "dreadnotes", "config.toml"), nil
}

// exists checks if the configuration file is present.
func exists() bool {
	path, err := getConfigPath()
	if err != nil {
		return false
	}

	_, err = os.Stat(path)

	return err == nil
}

// read fetches all lines from the configuration file.
func read() []string {
	path, err := getConfigPath()
	if err != nil {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening config: %v\n", err)
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

// split safely separates key and value.
func split(data string) (string, string) {
	parts := strings.SplitN(data, " = ", 2)
	if len(parts) != 2 {
		return "", ""
	}

	return parts[0], parts[1]
}

// validate verifies basic formatting rules.
func validate(key, value string) bool {
	switch key {
	case "notes_path", "editor", "templates_path":
		return strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")
	default:
		return false
	}
}
