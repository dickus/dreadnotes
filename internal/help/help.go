// Package help provides usage instructions and help messages for the CLI application.
// It utilizes text/tabwriter to align command descriptions and flags neatly.
package help

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// HelpData holds the information required to print a standardized help message.
type HelpData struct {
	Title       string
	Description string
	Usage       string
	Flags       [][2]string // Each element is {flag, description}
	Examples    []string
}

// printHelp formats and prints the help message using tabwriter.
func printHelp(data HelpData) {
	// Header and Description
	if data.Title != "" {
		fmt.Printf(" %s:\n", strings.ToUpper(data.Title))
	}
	if data.Description != "" {
		fmt.Printf("   dreadnotes %s ― %s\n", data.Title, data.Description)
		fmt.Println()
	}

	// Usage
	if data.Usage != "" {
		fmt.Println(" USAGE:")
		fmt.Printf("   %s\n", data.Usage)
		fmt.Println()
	}

	// Flags (aligned)
	if len(data.Flags) > 0 {
		fmt.Println(" FLAGS:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for _, f := range data.Flags {
			fmt.Fprintf(w, "   %s\t%s\n", f[0], f[1])
		}
		w.Flush()
		fmt.Println()
	}

	// Examples
	if len(data.Examples) > 0 {
		fmt.Println(" EXAMPLES:")
		for _, ex := range data.Examples {
			fmt.Printf("   %s\n", ex)
		}
	}
}

// Short prints a concise summary of available commands.
func Short() {
	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes <COMMAND> [FLAGS]")
	fmt.Println()
	fmt.Println(" COMMANDS:")
	fmt.Println("   new, open, random, sync, doctor")
	fmt.Println()
	fmt.Println(" Run 'dreadnotes --help' for detailed usage.")
}

// Long prints the full help message (root command).
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
	fmt.Fprintln(w, "   doctor\tCheck for problems")
	w.Flush()

	fmt.Println()
	fmt.Println(" FLAGS:")
	fmt.Fprintln(w, "   -h, --help\tShow this help")
	w.Flush()

	fmt.Println()
	fmt.Println(" ENVIRONMENT:")
	fmt.Fprintln(w, "   DREADNOTES_CONFIG\tOverride the default config file path")
	w.Flush()

	fmt.Println()
	fmt.Println(" Run 'dreadnotes <COMMAND> --help' for more information on a command.")
}

// NewNoteHelp displays usage for 'new' command.
func NewNoteHelp() {
	printHelp(HelpData{
		Title:       "new",
		Description: "Create new note",
		Usage:       "dreadnotes new [FLAGS] \"<NAME>\"",
		Flags: [][2]string{
			{"-h, --help", "Show this help"},
			{"-T <name>", "Use a specific template"},
			{"-i", "Pick a template interactively"},
		},
		Examples: []string{
			"dreadnotes new \"My Note\"",
			"dreadnotes new -T daily \"My Note\"",
			"dreadnotes new -i \"My Note\"",
		},
	})
}

// OpenNoteHelp displays usage for 'open' command.
func OpenNoteHelp() {
	printHelp(HelpData{
		Title:       "open",
		Description: "Search notes",
		Usage:       "dreadnotes open [FLAGS]",
		Flags: [][2]string{
			{"-h, --help", "Show this help"},
			{"-t, --tag", "Search by tag"},
		},
		Examples: []string{
			"dreadnotes open",
			"dreadnotes open -t work",
		},
	})
}

// RandomNoteHelp displays usage for 'random' command.
func RandomNoteHelp() {
	printHelp(HelpData{
		Title:       "random",
		Description: "Open random note",
		Usage:       "dreadnotes random [FLAGS]",
		Flags: [][2]string{
			{"-h, --help", "Show this help"},
		},
		Examples: []string{
			"dreadnotes random",
		},
	})
}

// SyncHelp displays usage for 'sync' command.
func SyncHelp() {
	printHelp(HelpData{
		Title:       "sync",
		Description: "Update local git repo. Also updates remote repo if it exists",
		Usage:       "dreadnotes sync [FLAGS]",
		Flags: [][2]string{
			{"-h, --help", "Show this help"},
		},
		Examples: []string{
			"dreadnotes sync",
		},
	})
}

// DoctorHelp displays usage for 'doctor' command.
func DoctorHelp() {
	printHelp(HelpData{
		Title:       "doctor",
		Description: "Check notes for duplicates, empty content and broken links",
		Usage:       "dreadnotes doctor [FLAGS]",
		Flags: [][2]string{
			{"-h, --help", "Show this help"},
		},
		Examples: []string{
			"dreadnotes doctor",
		},
	})
}
