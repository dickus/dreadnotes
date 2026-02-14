package help

import (
	"fmt"
	"os"
	"text/tabwriter"
)


func Short() {
	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes <command> [flags] [<args>]")
	fmt.Println()
	fmt.Println(" Run 'dreadnotes --help' for detailed usage.")
}

func Long() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes [FLAGS] <COMMAND> [OPTIONS]")
	fmt.Println()
	fmt.Println(" COMMANDS:")
	fmt.Fprintln(w, "   new\tCreate new note")

	w.Flush()

	fmt.Println()
	fmt.Println(" FLAGS:")
	fmt.Fprintln(w, "   -h, --help\tShow this help")

	w.Flush()
}

func NewNoteHelp() {
	fmt.Println(" NEW:")
	fmt.Println("   dreadnotes new ― Create new note")
	fmt.Println()
	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes new \"<name>\"")
}

func OpenNoteHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Println(" OPEN:")
	fmt.Println("   dreadnotes open ― Search notes")
	fmt.Println()
	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes open [FLAGS]")
	fmt.Println()
	fmt.Println(" FLAGS:")
	fmt.Fprintln(w, "   -h, --help\tShow this help")
	fmt.Fprintln(w, "   -t, --tag\tSearch by tag")

	w.Flush()
}

