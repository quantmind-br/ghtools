# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**ghtools** is a Go CLI tool that provides an interactive TUI for managing GitHub repositories. It wraps the `gh` CLI and provides menu-driven access to repository operations like listing, searching, cloning, forking, creating, deleting, and more.

## Build & Run Commands

```bash
# Build the binary
make build
# or: go build -o ghtools .

# Install to ~/.local/bin
make install

# Uninstall
make uninstall

# Clean build artifacts
make clean
```

**Requirements:**
- Go 1.25+
- `gh` CLI (GitHub CLI) - must be authenticated
- `git` CLI

## Architecture

```
ghtools/
├── main.go           # Entry point, calls cmd.Execute()
├── cmd/              # CLI commands (Cobra)
│   ├── root.go       # Root command, version, global flags
│   ├── menu.go       # Interactive menu (main navigation)
│   ├── list.go       # List repositories
│   ├── search.go     # Search repositories
│   ├── clone.go      # Clone repositories
│   ├── sync_cmd.go   # Sync local repos
│   ├── status.go     # Local repo status
│   ├── fork.go       # Fork repository
│   ├── create.go     # Create repository
│   ├── delete.go     # Delete repository
│   ├── archive.go    # Archive/Unarchive
│   ├── visibility.go # Change visibility
│   ├── browse.go     # Browse in browser
│   ├── pr.go         # Pull requests
│   ├── config_cmd.go # Config management
│   ├── refresh.go    # Refresh cache
│   ├── trending.go   # Trending repos
│   └── stats.go      # Statistics dashboard
│
├── internal/         # Core packages
│   ├── cache/        # Caching functionality
│   ├── config/       # Configuration (JSON file in ~/.config/ghtools/)
│   ├── gh/           # GitHub API wrapper (calls `gh` CLI)
│   ├── git/          # Git operations
│   ├── runner/       # Parallel execution
│   ├── template/     # Template handling
│   ├── types/        # Type definitions
│   └── tui/          # TUI components (Bubble Tea + Lipgloss)
│       ├── styles.go
│       ├── print.go
│       ├── input.go
│       ├── confirm.go
│       ├── choose.go
│       ├── multiselect.go
│       ├── spinner.go
│       └── table.go
│
└── _legacy/          # Original bash script (deprecated)
```

## Key Frameworks

- **Cobra** (v1.10.2) - CLI command framework
- **Bubble Tea** (v1.3.10) - TUI framework
- **Lipgloss** (v1.1.0) - Terminal styling

## Global Flags

- `-V, --verbose` - Enable verbose output
- `-q, --quiet` - Suppress non-error output
- `-y, --yes` - Non-interactive mode (auto-confirm defaults)

## TUI Navigation

The main menu provides keyboard-driven access:

| Key | Action |
|-----|--------|
| L | List Repositories |
| S | Search My Repos |
| E | Explore GitHub |
| T | Trending Repos |
| D | Statistics Dashboard |
| C | Clone Repositories |
| Y | Sync Local Repos |
| O | Local Repo Status |
| F | Fork Repository |
| R | Create Repository |
| X | Delete Repositories |
| A | Archive/Unarchive |
| V | Change Visibility |
| B | Browse in Browser |
| P | Pull Requests |
| G | Config |
| M | Refresh Cache |
| Q | Exit |

## Notes

- No test files exist in this project
- The TUI components are in `internal/tui/` - useful for adding new UI elements
- Configuration is stored in `~/.config/ghtools/config.json`
- Cache is stored in `~/.cache/ghtools/`
