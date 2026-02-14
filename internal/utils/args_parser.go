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

	//NOTES
	newNoteCmd := flag.NewFlagSet("new", flag.ExitOnError)

	switch os.Args[1] {
	case "-h":
		fallthrough

	case "--help":
		help.Long()

		os.Exit(0)

	case "new":

		newNoteCmd.Parse(os.Args[2:])

		if newNoteCmd.NArg() == 0 {
			notes.NewNote("", models.Cfg.NotesPath)
		} else {
			notes.NewNote(newNoteCmd.Arg(0), models.Cfg.NotesPath)
		}

	case "open":
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

		p := tea.NewProgram(ui.NewSearchModel(idx))
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

	default:
		fmt.Println("Unknown argument: ", os.Args[1])

		os.Exit(1)
	}
}

