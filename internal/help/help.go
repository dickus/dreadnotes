package help

import "fmt"

func Short() {
	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes <command> [flags] [<args>]")
	fmt.Println()
	fmt.Println(" Run 'dreadnotes --help' for detailed usage.")
}

func Long() {
	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes [FLAGS] <COMMAND> [OPTIONS]")
	fmt.Println()
	fmt.Println(" COMMANDS:")
	fmt.Println("   config\tSet preferences")
	fmt.Println()
	fmt.Println(" FLAGS:")
	fmt.Println("   -h, --help\tShow this help")
}

func ConfigHelp() {
	fmt.Println(" CONFIG:")
	fmt.Println("   dreadnotes config ― Set preferences")
	fmt.Println()
	fmt.Println(" USAGE:")
	fmt.Println("   dreadnotes config [OPTIONS]")
	fmt.Println()
	fmt.Println(" OPTIONS:")
	fmt.Println("   -h, --help\t\t\tShow this help")
	fmt.Println("   -p, --path '<PATH>'\t\tSet notes directory path (default: $HOME/Documents/dreadnotes)")
	fmt.Println("   -e, --editor '<NAME>'\tSet notes editor (default: nvim)")
}

