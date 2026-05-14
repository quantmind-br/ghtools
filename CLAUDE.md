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

<!-- gitnexus:start -->
# GitNexus — Code Intelligence

This project is indexed by GitNexus as **ghtools** (732 symbols, 2181 relationships, 56 execution flows). Use the GitNexus MCP tools to understand code, assess impact, and navigate safely.

> If any GitNexus tool warns the index is stale, run `npx gitnexus analyze` in terminal first.

## Always Do

- **MUST run impact analysis before editing any symbol.** Before modifying a function, class, or method, run `gitnexus_impact({target: "symbolName", direction: "upstream"})` and report the blast radius (direct callers, affected processes, risk level) to the user.
- **MUST run `gitnexus_detect_changes()` before committing** to verify your changes only affect expected symbols and execution flows.
- **MUST warn the user** if impact analysis returns HIGH or CRITICAL risk before proceeding with edits.
- When exploring unfamiliar code, use `gitnexus_query({query: "concept"})` to find execution flows instead of grepping. It returns process-grouped results ranked by relevance.
- When you need full context on a specific symbol — callers, callees, which execution flows it participates in — use `gitnexus_context({name: "symbolName"})`.

## Never Do

- NEVER edit a function, class, or method without first running `gitnexus_impact` on it.
- NEVER ignore HIGH or CRITICAL risk warnings from impact analysis.
- NEVER rename symbols with find-and-replace — use `gitnexus_rename` which understands the call graph.
- NEVER commit changes without running `gitnexus_detect_changes()` to check affected scope.

## Resources

| Resource | Use for |
|----------|---------|
| `gitnexus://repo/ghtools/context` | Codebase overview, check index freshness |
| `gitnexus://repo/ghtools/clusters` | All functional areas |
| `gitnexus://repo/ghtools/processes` | All execution flows |
| `gitnexus://repo/ghtools/process/{name}` | Step-by-step execution trace |

## CLI

| Task | Read this skill file |
|------|---------------------|
| Understand architecture / "How does X work?" | `.claude/skills/gitnexus/gitnexus-exploring/SKILL.md` |
| Blast radius / "What breaks if I change X?" | `.claude/skills/gitnexus/gitnexus-impact-analysis/SKILL.md` |
| Trace bugs / "Why is X failing?" | `.claude/skills/gitnexus/gitnexus-debugging/SKILL.md` |
| Rename / extract / split / refactor | `.claude/skills/gitnexus/gitnexus-refactoring/SKILL.md` |
| Tools, resources, schema reference | `.claude/skills/gitnexus/gitnexus-guide/SKILL.md` |
| Index, status, clean, wiki CLI commands | `.claude/skills/gitnexus/gitnexus-cli/SKILL.md` |

<!-- gitnexus:end -->
