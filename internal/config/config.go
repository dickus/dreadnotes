// Package config handles the application configuration management.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dickus/dreadnotes/internal/utils"
)

// Load sets default configuration values and overrides them with settings from the user's config file if it exists.
func Load() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't find home directory: %v\n", err)
	}

	conf, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't find config directory: %v\n", err)
	}

	// Set robust default absolute paths using filepath.Join
	Cfg.RepoPath = filepath.Join(home, "Documents", "dreadnotes")
	Cfg.NotesPath = filepath.Join(Cfg.RepoPath, "notes")
	Cfg.Editor = "nvim"
	Cfg.Templates = filepath.Join(conf, "dreadnotes", "templates")

	if !exists() {
		return
	}

	configStrings := read()
	var pathSeen, editorSeen, templateSeen bool

	for _, data := range configStrings {
		if !strings.Contains(data, "=") {
			fmt.Printf("Incorrect data format in '%s'.\n", data)
			continue
		}

		key, value := split(data)

		if !validate(key, value) {
			fmt.Printf("'%s' is not a valid data type.\n", value)
			continue
		}

		// Idiomatic way to remove surrounding quotes
		value = strings.Trim(value, "\"")

		switch key {
		case "notes_path":
			if !pathSeen {
				// Expand $HOME and update BOTH RepoPath and NotesPath correctly
				parsedPath := utils.PathParse(value)
				Cfg.RepoPath = parsedPath
				Cfg.NotesPath = filepath.Join(parsedPath, "notes")
				pathSeen = true
			} else {
				fmt.Printf("Duplicate '%s'. Using: %s\n", key, Cfg.NotesPath)
			}

		case "editor":
			if !editorSeen {
				Cfg.Editor = value
				editorSeen = true
			} else {
				fmt.Printf("Duplicate '%s'. Using: %s\n", key, Cfg.Editor)
			}

		case "templates_path":
			if !templateSeen {
				// We now store the TRUE absolute path, handling $HOME via PathParse
				Cfg.Templates = utils.PathParse(value)
				templateSeen = true
			} else {
				fmt.Printf("Duplicate '%s'. Using: %s\n", key, Cfg.Templates)
			}

		default:
			fmt.Printf("Key '%s' is unknown. Check config.toml.\n", key)
		}
	}
}
