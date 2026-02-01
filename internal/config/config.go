package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dickus/dreadnotes/internal/utils"
)

type Config struct {
	NotesPath string
	Editor string
}

var cfg Config

func LoadConfig() {
	home, _ := os.UserHomeDir()

	cfg.NotesPath = filepath.Join(home, "Documents/dreadnotes")
	cfg.Editor = "nvim"

	if utils.ConfigInPlace() {
		configStrings := utils.ReadConfig()

		pathKey := 0
		editorKey := 0

		for _, data := range configStrings {
			if !strings.Contains(data, "=") {
				fmt.Printf("Incorrect data format in '%s'.\n", data)

				continue
			}

			key, value := utils.SplitConfig(data)

			if !utils.DataValidation(key, value) {
				fmt.Printf("'%s' is not a valid data type.\n", value)

				continue
			}

			if key == "notes_path" {
				if pathKey == 0 {
					value = strings.ReplaceAll(value, "\"", "")

					cfg.NotesPath = value

					pathKey++
				} else {
					fmt.Printf("'%s' has a duplicate. Check config.toml to resolve this issue. Path %s will be used now.\n", key, cfg.NotesPath)
				}
			} else if key == "editor" {
				if editorKey == 0 {
					value = strings.ReplaceAll(value, "\"", "")

					cfg.Editor = value

					editorKey++
				} else {
					fmt.Printf("'%s' has a duplicate. Check config.toml to resolve this issue. Editor %s will be used now.\n", key, cfg.Editor)
				}
			} else {
				fmt.Printf("Key '%s' is unknown. Check config.toml to resolve this issue.\n", key)
			}
		}
	}
}

