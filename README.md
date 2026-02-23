[![Go Reference](https://pkg.go.dev/badge/github.com/dickus/dreadnotes.svg)](https://pkg.go.dev/github.com/dickus/dreadnotes)
[![Go Report Card](https://goreportcard.com/badge/github.com/dickus/dreadnotes)](https://goreportcard.com/report/github.com/dickus/dreadnotes)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/dickus/dreadnotes)](https://github.com/dickus/dreadnotes/releases)

**Dreadnotes** is yet another simple and lightweight CLI tool for managing your notes/knowledge base.

## Installation

### Go-way

```bash
go install github.com/dickus/dreadnotes@latest
```

Make sure your $GOPATH is in your $PATH.

### Binary

1. Download the archive from [releases page](https://github.com/dickus/dreadnotes/releases).
2. Extract the file.
3. Put the file to your $PATH.

### Build from source

Clone the repo.

#### Locally

Run `make install PREFIX=$HOME/.local`.

Make sure that $HOME/.local/bin is in your $PATH.

#### System-wide

Run `sudo make install`.

## Configuration

**Dreadnotes** has minimal configuration setup. It looks into $HOME/.config/dreadnotes for config.toml by default.

### Default config

```TOML
notes_path = "$HOME/Documents/dreadnotes"
editor = "nvim"
templates_path = "$HOME/.config/dreadnotes/templates"
```

### Multiple "vaults"

If you wish to split your notes into several "vaults", you can use the `DREADNOTES_CONFIG` environment variable. It will work as a different storage, so different Git repo, different search index, etc.

1. Create a secondary config file (e.g., `$HOME/.config/dreadnotes/work.toml`) with different `notes_path`.
2. Add an alias to your `~/.bashrc` or `~/.zshrc`:

```bash
# 'dn' for default notes vault
alias dn="dreadnotes"

# 'dnw' for work notes
alias dnw="DREADNOTES_CONFIG=$HOME/.config/dreadnotes/work.toml dreadnotes"
```

`dnw open` will only search your work documents, etc.

## Usage
The general syntax for the CLI is:

```bash
dreadnotes <COMMAND> [FLAGS]
```

Run `dreadnotes --help` at any time to see the available commands.

### Create (`new`)

You can create a note with a default layout, specify a template, or choose one interactively.

Avoid using / (slash) character in template and note names.

**Usage:**
```bash
dreadnotes new [FLAGS] "<NAME>"
```

**Options:**
| Flag | Description |
| :--- | :--- |
| `-T <name>` | Use a specific template by name (e.g., `-T daily`) |
| `-i` | Pick a template via an interactive menu |
| `-h, --help` | Show help for this command |

**Examples:**
```bash
# Create a simple note
dreadnotes new "Project Idea"

# Create a note using the 'meeting' template
dreadnotes new -T meeting "Team Sync"

# Create a note using the 'meeting plan' template
dreadnotes new -T "meeting plan" tomorrow meeting plan

# Choose a template interactively
dreadnotes new -i "Refactoring Plan"
```

### Find (`open`)

You can browse using titles and content or filter by specific tags.

**Usage:**
```bash
dreadnotes open [FLAGS]
```

**Options:**
| Flag | Description |
| :--- | :--- |
| `-t, --tag` | Filter search results by a specific tag |
| `-h, --help` | Show help for this command |

**Examples:**
```bash
# Search notes by their titles and content
dreadnotes open

# Search notes by a tag specified in search field
dreadnotes open -t
```

### Rediscover (`random`)

Open a random note from your vault.

**Usage:**
```bash
dreadnotes random
```

### Sync (`sync`)

Update your local notes repository. This command handles local git operations and pushes changes to a remote repository if one is configured.

Before you can use this command, you need to set up local repo and link it to a remote one if you wish.

**Usage:**
```bash
dreadnotes sync
```

### Fix (`doctor`)

Check your notes for broken wikilinks, duplicate titles and empty content.

**Usage:**
```bash
dreadnotes doctor
```
