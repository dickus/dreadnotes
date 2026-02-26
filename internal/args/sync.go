package args

import (
	"flag"
	"fmt"
	"os"

	"github.com/dickus/dreadnotes/internal/config"
	"github.com/dickus/dreadnotes/internal/help"
	"github.com/dickus/dreadnotes/internal/sync"
)

func syncNotes() {
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
}
