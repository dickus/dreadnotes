package help

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func Short() {
	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes <COMMAND> [FLAGS]")
	fmt.Println()
	fmt.Println(" COMMANDS:")
	fmt.Println("   new, open, random, sync")
	fmt.Println()
	fmt.Println(" Run 'dreadnotes --help' for detailed usage.")
}

func Long() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes <COMMAND> [FLAGS]")
	fmt.Println()
	fmt.Println(" COMMANDS:")
	fmt.Fprintln(w, "   new\tCreate new note")
	fmt.Fprintln(w, "   open\tSearch notes")
	fmt.Fprintln(w, "   random\tOpen random note")
	fmt.Fprintln(w, "   sync\tUpdate git repository")

	w.Flush()

	fmt.Println()
	fmt.Println(" FLAGS:")
	fmt.Fprintln(w, "   -h, --help\tShow this help")

	w.Flush()

	fmt.Println()
	fmt.Println(" Run 'dreadnotes <COMMAND> --help' for more information on a command.")
}

func NewNoteHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Println(" NEW:")
	fmt.Println("   dreadnotes new ― Create new note")
	fmt.Println()
	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes new [FLAGS] \"<NAME>\"")
	fmt.Println()
	fmt.Println(" FLAGS:")
	fmt.Fprintln(w, "   -h, --help\tShow this help")
	fmt.Fprintln(w, "   -T <name>\tUse a specific template")
	fmt.Fprintln(w, "   -i\tPick a template interactively")

	w.Flush()

	fmt.Println()
	fmt.Println(" EXAMPLES:")
	fmt.Println("   dreadnotes new \"My Note\"")
	fmt.Println("   dreadnotes new -T daily \"My Note\"")
	fmt.Println("   dreadnotes new -i \"My Note\"")
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

	fmt.Println()
	fmt.Println(" EXAMPLES:")
	fmt.Println("   dreadnotes open")
	fmt.Println("   dreadnotes open -t work")
}

func RandomNoteHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Println(" RANDOM:")
	fmt.Println("   dreadnotes random ― Open random note")
	fmt.Println()
	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes random [FLAGS]")
	fmt.Println()
	fmt.Println(" FLAGS:")
	fmt.Fprintln(w, "   -h, --help\tShow this help")

	w.Flush()

	fmt.Println()
	fmt.Println(" EXAMPLES:")
	fmt.Println("   dreadnotes random")
}

func SyncHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Println(" SYNC:")
	fmt.Println("   dreadnotes sync ― Update local git repo. Also updates remote repo if it exists")
	fmt.Println()
	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes sync [FLAGS]")
	fmt.Println()
	fmt.Println(" FLAGS:")
	fmt.Fprintln(w, "   -h, --help\tShow this help")

	w.Flush()

	fmt.Println()
	fmt.Println(" EXAMPLES:")
	fmt.Println("   dreadnotes sync")
}
