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
	fmt.Fprintln(w, "   config\tSet preferences")
	fmt.Fprintln(w, "   new\tCreate new note")

	w.Flush()

	fmt.Println()
	fmt.Println(" FLAGS:")
	fmt.Fprintln(w, "   -h, --help\tShow this help")

	w.Flush()
}

func ConfigHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Println(" CONFIG:")
	fmt.Println("   dreadnotes config ― Set preferences")
	fmt.Println()
	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes config [OPTIONS]")
	fmt.Println()
	fmt.Println(" OPTIONS:")
	fmt.Fprintln(w, "   -h, --help\tShow this help")
	fmt.Fprintln(w, "   -p, --path '<PATH>'\tSet notes directory path (default: $HOME/Documents/dreadnotes)")
	fmt.Fprintln(w, "   -e, --editor '<NAME>'\tSet notes editor (default: nvim)")

	w.Flush()
}

func NewNoteHelp() {
	fmt.Println(" NEW:")
	fmt.Println("   dreadnotes new ― Create new note")
	fmt.Println()
	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes new \"<name>\"")
}

