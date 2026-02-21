// Dreadnotes is a lightweight command-line interface (CLI) tool for managing your personal knowledge base.
//
// It streamlines note creation with templates, offers quick search capabilities, and handles simple Git synchronization.
package main

import (
	"github.com/dickus/dreadnotes/internal/args"
	"github.com/dickus/dreadnotes/internal/config"
)

func main() {
	config.Load()
	args.ArgsParser()
}
