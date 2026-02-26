package args

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dickus/dreadnotes/internal/config"
	"github.com/dickus/dreadnotes/internal/help"
	"github.com/dickus/dreadnotes/internal/notes"
	"github.com/dickus/dreadnotes/internal/ui"
	"github.com/dickus/dreadnotes/internal/utils"
)

func newNote() {
	newCmd := flag.NewFlagSet("new", flag.ExitOnError)

	newCmd.Usage = func() {
		help.NewNoteHelp()

		os.Exit(0)
	}

	tmpl := newCmd.String("T", "", "template name")
	pick := newCmd.Bool("i", false, "interactive template pick")

	newCmd.Parse(os.Args[2:])

	name := strings.Join(newCmd.Args(), " ")

	var tmplPath string

	tmplDir := utils.PathParse(config.Cfg.Templates)

	if !filepath.IsAbs(tmplDir) {
		confDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get user config directory: %v\n", err)

			os.Exit(1)
		}
		tmplDir = filepath.Join(confDir, tmplDir)
	}

	if *pick {
		picked, err := ui.RunTemplatePicker(tmplDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Template picker error: %v\n", err)

			os.Exit(1)
		}
		tmplPath = picked

	} else if *tmpl != "" {
		tmplPath = filepath.Join(tmplDir, *tmpl+".md")

		if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Template not found: %s\n", tmplPath)

			os.Exit(1)
		}
	}

	notes.NewNote(name, tmplPath)
}
