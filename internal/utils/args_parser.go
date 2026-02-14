package utils

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dickus/dreadnotes/internal/help"
	"github.com/dickus/dreadnotes/internal/models"
	"github.com/dickus/dreadnotes/internal/notes"
	"github.com/dickus/dreadnotes/internal/search"
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

		newCmd.Parse(os.Args[2:])

		name := newCmd.Arg(0)
		notes.NewNote(name, models.Cfg.NotesPath)

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

	default:
		help.Short()

		os.Exit(0)
	}
}

