// Package models provides structure for config.
package config

// Config holds the global application configuration settings.
type Config struct {
	RepoPath  string // Absolute path to the git repository root
	NotesPath string // Absolute path to the directory containing notes
	Editor    string // Command to launch the preferred text editor (e.g., "vim", "code")
	Templates string // Path to the directory containing note templates
}

// Cfg is the global configuration instance used throughout the application.
// It is initialized once during startup.
var Cfg Config
