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

You can browse using titles and content filtered by creation and modification dates. You can also filter search by specific tags.

To move between found notes use Alt-j/k. It was made like this to avoid issues with using tmux.

To switch between creation/modification dates filter use Alt-d.

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

## Neovim tips

If you're using Neovim, I suggest using several functions to make the experience a bit more pleasant.

### Update time on write

```lua
vim.api.nvim_create_autocmd("BufWritePre", {
    pattern = "*.md",
    callback = function()
        local save_cursor = vim.fn.getpos(".")
        local n = math.min(10, vim.fn.line("$"))
        local lines = vim.api.nvim_buf_get_lines(0, 0, n, false)

        local now = os.time()
        local tz = os.date("%z", now)
        local tz_formatted = tz:sub(1, 3) .. ":" .. tz:sub(4, 5)
        local new_date = os.date("updated: %Y-%m-%d %H:%M")

        for i, line in ipairs(lines) do
            if line:match("^updated:") then
                vim.api.nvim_buf_set_lines(0, i-1, i, false, {new_date})
                break
            end
        end

        vim.fn.setpos(".", save_cursor)
    end,
})
```

This function updates the modified time of the note when you save the file. It will help **dreadnotes** know that the note is updated based on frontmatter field.

### Easy tagging

```lua
local function add_tags()
    vim.ui.input({ prompt = "Tags (separated by comma): " }, function(input)
        if not input or input == "" then
            return
        end

        local tags = {}
        for tag in input:gmatch("[^,]+") do
            tag = tag:match("^%s*(.-)%s*$")
            if tag ~= "" then
                table.insert(tags, tag)
            end
        end

        local tags_line = "tags: [" .. table.concat(tags, ", ") .. "]"

        local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
        local frontmatter_start = nil
        local frontmatter_end = nil
        local existing_tags_line = nil

        for i, line in ipairs(lines) do
            if line:match("^%-%-%-") then
                if not frontmatter_start then
                    frontmatter_start = i
                else
                    frontmatter_end = i
                    break
                end
            end
            if frontmatter_start and not frontmatter_end then
                if line:match("^tags:") then
                    existing_tags_line = i
                end
            end
        end

        if frontmatter_start and frontmatter_end then
            if existing_tags_line then
                local old = lines[existing_tags_line]:match("%[(.-)%]") or ""
                for tag in old:gmatch("[^,]+") do
                    tag = tag:match("^%s*(.-)%s*$")
                    if tag ~= "" then
                        local exists = false
                        for _, t in ipairs(tags) do
                            if t == tag then exists = true; break end
                        end
                        if not exists then
                            table.insert(tags, 1, tag)
                        end
                    end
                end
                tags_line = "tags: [" .. table.concat(tags, ", ") .. "]"
                vim.api.nvim_buf_set_lines(0, existing_tags_line - 1, existing_tags_line, false, { tags_line })
            end
        else
            local header = { "---", tags_line, "---", "" }
            vim.api.nvim_buf_set_lines(0, 0, 0, false, header)
        end
    end)
end

vim.keymap.set("n", "<leader>tt", add_tags, { desc = " Add tags to frontmatter" })
```

This function will let you add tags faster without having to move to tags frontmatter field. Be aware that it works only for tags placed in [].
