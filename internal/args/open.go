package args

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dickus/dreadnotes/internal/config"
	"github.com/dickus/dreadnotes/internal/help"
	"github.com/dickus/dreadnotes/internal/notes"
	"github.com/dickus/dreadnotes/internal/search"
	"github.com/dickus/dreadnotes/internal/ui"
)

func openNote() {
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
}
