package args

import (
	"flag"
	"fmt"
	"os"

	"github.com/dickus/dreadnotes/internal/config"
	"github.com/dickus/dreadnotes/internal/doctor"
	"github.com/dickus/dreadnotes/internal/help"
)

func doctorNotes() {
	doctorCmd := flag.NewFlagSet("doctor", flag.ExitOnError)

	doctorCmd.Usage = func() {
		help.DoctorHelp()

		os.Exit(0)
	}

	doctorCmd.Parse(os.Args[2:])

	report, err := doctor.Run(config.Cfg.NotesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Doctor failed: %v\n", err)

		os.Exit(1)
	}

	doctor.PrintReport(report)
}
