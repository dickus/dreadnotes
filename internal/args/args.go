// Package args implements command-line arguments parsing and routing.
package args

import (
	"os"

	"github.com/dickus/dreadnotes/internal/help"
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
		newNote()

	case "open":
		openNote()

	case "random":
		randomNote()

	case "sync":
		syncNotes()

	case "doctor":
		doctorNotes()

	default:
		help.Short()

		os.Exit(0)
	}
}
