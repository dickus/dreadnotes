package utils

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dickus/dreadnotes/internal/help"
)

func ArgsParser() {
	if len(os.Args) < 2 {
		help.Short()

		os.Exit(0)
	}

	//CONFIG
	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	var (
		configNotesPath string
		configEditor string
		configHelp bool
	)
	configCmd.StringVar(&configNotesPath, "path", "$HOME/Documents/dreadnotes", "config editor flag")
	configCmd.StringVar(&configNotesPath, "p", "$HOME/Documents/dreadnotes", "config editor flag")
	configCmd.StringVar(&configEditor, "editor", "nvim", "config editor flag")
	configCmd.StringVar(&configEditor, "e", "nvim", "config editor flag")
	configCmd.BoolVar(&configHelp, "help", false, "config flag")
	configCmd.BoolVar(&configHelp, "h", false, "config help flag")

	switch os.Args[1] {
	case "config":
		configCmd.Parse(os.Args[2:])

		if configHelp {
			help.ConfigHelp()

			os.Exit(0)
		}

		var sb strings.Builder

		fmt.Fprintf(&sb, "notes_path = \"%v\"\n", configNotesPath)
		fmt.Fprintf(&sb, "editor = \"%v\"\n", configEditor)

		finalString := sb.String()

		SaveConfig(finalString)
	case "-h":
		fallthrough
	case "--help":
		help.Long()

		os.Exit(0)
	default:
		fmt.Println("Unknown argument: ", os.Args[1])

		os.Exit(1)
	}
}

