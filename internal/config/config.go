package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dickus/dreadnotes/internal/utils"
	"github.com/dickus/dreadnotes/internal/models"
)

func LoadConfig() {
	home, _ := os.UserHomeDir()
	conf, _ := os.UserConfigDir()

	models.Cfg.NotesPath = filepath.Join(home, "Documents/dreadnotes")
	models.Cfg.Editor = "nvim"
	models.Cfg.Templates = filepath.Join(conf, "dreadnotes", "templates")

	if utils.ConfigInPlace() {
		configStrings := utils.ReadConfig()

		pathKey := 0
		editorKey := 0
		templateKey := 0

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

			switch key {
			case "notes_path":
				if pathKey == 0 {
					value = strings.ReplaceAll(value, "\"", "")

					models.Cfg.NotesPath = value

					pathKey++
				} else {
					fmt.Printf("'%s' has a duplicate. Check config.toml to resolve this issue. Path %s will be used now.\n", key, models.Cfg.NotesPath)
				}

			case "editor":
				if editorKey == 0 {
					value = strings.ReplaceAll(value, "\"", "")

					models.Cfg.Editor = value

					editorKey++
				} else {
					fmt.Printf("'%s' has a duplicate. Check config.toml to resolve this issue. Editor %s will be used now.\n", key, models.Cfg.Editor)
				}

			case "templates_path":
				if templateKey == 0 {
					value = strings.ReplaceAll(value, "\"", "")

					models.Cfg.Templates = value

					models.Cfg.Templates = strings.Replace(models.Cfg.Templates, "$HOME/.config/", "", 1)

					templateKey++
				} else {
					fmt.Printf("'%s' has a duplicate. Check config.toml to resolve this issue. Templates path %s will be used now.\n", key, models.Cfg.Templates)
				}

			default:
				fmt.Printf("Key '%s' is unknown. Check config.toml to resolve this issue.\n", key)
			}
		}
	}
}

