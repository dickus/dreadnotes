package utils

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dickus/dreadnotes/internal/help"
	"github.com/dickus/dreadnotes/internal/models"
	"github.com/dickus/dreadnotes/internal/notes"
	"github.com/dickus/dreadnotes/internal/search"
	"github.com/dickus/dreadnotes/internal/sync"
	"github.com/dickus/dreadnotes/internal/ui"
)

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

		name := newCmd.Arg(0)

		var tmplPath string

		if *pick {
			picked, err := ui.RunTemplatePicker(models.Cfg.Templates)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)

				os.Exit(1)
			}

			tmplPath = picked
		} else if *tmpl != "" {
			tmplPath = filepath.Join(models.Cfg.Templates, *tmpl+".md")
			confDir, _ := os.UserConfigDir()
			tmplPath = confDir + "/" + tmplPath

			if _, err := os.Stat(tmplPath); err != nil {
				fmt.Fprintf(os.Stderr, "Template not found: %s\n", *tmpl)

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

		idx, err := search.BuildIndex(models.Cfg.NotesPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}
		defer idx.Close()

		if err := search.ReindexAll(idx, models.Cfg.NotesPath); err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}

		p := tea.NewProgram(ui.NewSearchModel(idx, *tagMode))
		result, err := p.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}

		if sm, ok := result.(ui.SearchModel); ok && sm.Chosen() != "" {
			if err := notes.OpenNote(sm.Chosen()); err != nil {
				fmt.Fprintln(os.Stderr, err)

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

		path, err := notes.RandomNote(models.Cfg.NotesPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}

		if err := notes.OpenNote(path); err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}

	case "sync":
		syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)
		syncCmd.Usage = func() {
			help.SyncHelp()

			os.Exit(0)
		}

		syncCmd.Parse(os.Args[2:])

		err := sync.Sync(models.Cfg.NotesPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}

	default:
		help.Short()

		os.Exit(0)
	}
}
