package args

import (
	"flag"
	"fmt"
	"os"

	"github.com/dickus/dreadnotes/internal/config"
	"github.com/dickus/dreadnotes/internal/help"
	"github.com/dickus/dreadnotes/internal/notes"
)

func randomNote() {
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
}
