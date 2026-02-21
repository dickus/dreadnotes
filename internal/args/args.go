// Package args implements command-line arguments parsing and routing.
package args

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dickus/dreadnotes/internal/config"
	"github.com/dickus/dreadnotes/internal/help"
	"github.com/dickus/dreadnotes/internal/notes"
	"github.com/dickus/dreadnotes/internal/search"
	"github.com/dickus/dreadnotes/internal/sync"
	"github.com/dickus/dreadnotes/internal/ui"
	"github.com/dickus/dreadnotes/internal/utils"
)

// ArgsParser parses command-line arguments and routes execution to the appropriate subcommand.
func ArgsParser() {
	if len(os.Args) < 2 {
		help.Short()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "-h", "--help":
		help.Long()

		os.Exit(0)

	case "new":
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

	case "open":
		openCmd := flag.NewFlagSet("open", flag.ExitOnError)

		openCmd.Usage = func() {
			help.OpenNoteHelp()

			os.Exit(0)
		}

		tagMode := openCmd.Bool("t", false, "search by tag")
		openCmd.BoolVar(tagMode, "tag", false, "search by tag")

		openCmd.Parse(os.Args[2:])

		idx, err := search.BuildIndex(config.Cfg.NotesPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to build search index: %v\n", err)

			os.Exit(1)
		}
		defer idx.Close()

		if err := search.ReindexAll(idx, config.Cfg.NotesPath); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to reindex notes: %v\n", err)

			os.Exit(1)
		}

		p := tea.NewProgram(ui.NewSearchModel(idx, *tagMode))
		result, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Search UI error: %v\n", err)

			os.Exit(1)
		}

		if sm, ok := result.(ui.SearchModel); ok && sm.Chosen() != "" {
			if err := notes.OpenNote(sm.Chosen()); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to open note: %v\n", err)

				os.Exit(1)
			}
		}

	case "random":
		randomCmd := flag.NewFlagSet("random", flag.ExitOnError)

		randomCmd.Usage = func() {
			help.RandomNoteHelp()

			os.Exit(0)
		}

		randomCmd.Parse(os.Args[2:])

		path, err := notes.RandomNote(config.Cfg.NotesPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get random note: %v\n", err)

			os.Exit(1)
		}

		if err := notes.OpenNote(path); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open random note: %v\n", err)

			os.Exit(1)
		}

	case "sync":
		syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)

		syncCmd.Usage = func() {
			help.SyncHelp()

			os.Exit(0)
		}

		syncCmd.Parse(os.Args[2:])

		err := sync.Sync(config.Cfg.NotesPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Sync failed: %v\n", err)

			os.Exit(1)
		}

	default:
		help.Short()

		os.Exit(0)
	}
}
