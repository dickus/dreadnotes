package utils

import (
	"flag"
	"fmt"
	"os"

	"github.com/dickus/dreadnotes/internal/help"
	"github.com/dickus/dreadnotes/internal/notes"
	"github.com/dickus/dreadnotes/internal/models"
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
	default:
		fmt.Println("Unknown argument: ", os.Args[1])

		os.Exit(1)
	}
}

