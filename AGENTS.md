# PROJECT KNOWLEDGE BASE

**Generated:** 2025-05-14
**Commit:** 1415f7d
**Branch:** master

## OVERVIEW

**ghtools** is a Go CLI tool that provides an interactive TUI for managing GitHub repositories. It wraps the `gh` CLI and provides menu-driven access to repository operations.

- **Go Version**: 1.25+
- **Dependencies**: Cobra v1.10.2, Bubble Tea v1.3.10, Lipgloss v1.1.0

---

## STRUCTURE

```
.
├── cmd/              # CLI commands (Cobra)
├── internal/         # Core packages
│   ├── cache/        # Caching functionality
│   ├── config/       # Configuration management
│   ├── gh/           # GitHub API wrapper (gh CLI)
│   ├── git/          # Git operations
│   ├── runner/       # Parallel execution
│   ├── template/     # Template handling
│   ├── types/        # Type definitions
│   └── tui/          # TUI components (Bubble Tea + Lipgloss)
├── _legacy/          # Deprecated bash script
├── main.go           # Entry point
├── go.mod            # Module definition
└── Makefile          # Build automation
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add new CLI command | `cmd/` | Follow `verb.go` naming pattern |
| Add TUI component | `internal/tui/` | Use existing styles + Bubble Tea patterns |
| GitHub API call | `internal/gh/` | Wraps `gh` CLI, not direct API |
| Git operations | `internal/git/` | Shells out to `git` CLI |
| Config changes | `internal/config/` | JSON file in `~/.config/ghtools/` |
| Cache logic | `internal/cache/` | File-based cache in `~/.cache/ghtools/` |
| Parallel tasks | `internal/runner/` | Worker pool pattern |

## CODE MAP

| Symbol | Type | Location | Role |
|--------|------|----------|------|
| `Execute` | Function | `cmd/root.go:62` | Entry point from main.go |
| `runMenu` | Function | `cmd/menu.go:9` | Interactive menu loop |
| `FetchRepos` | Function | `internal/gh/repos.go:13` | List user repos via gh CLI |
| `CloneRepo` | Function | `internal/gh/repos.go:36` | Clone a repository |
| `RunWithSpinner` | Function | `internal/tui/spinner.go:60` | Async task with spinner UI |
| `RunChoose` | Function | `internal/tui/choose.go:185` | Single-select prompt |
| `RunMultiSelect` | Function | `internal/tui/multiselect.go:197` | Multi-select prompt |
| `ParallelRunner` | Struct | `internal/runner/parallel.go:14` | Worker pool for parallel ops |
| `Config` | Struct | `internal/config/config.go:10` | App configuration |
| `Repo` | Struct | `internal/types/types.go:5` | GitHub repo data model |

## Build Commands

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

# Run with verbose output
./ghtools -V

# Run in quiet mode
./ghtools -q

# Run in non-interactive mode (auto-confirm defaults)
./ghtools -y
```

**Testing**: There are currently no test files in this project (`go test` will pass with no test files).

---

## Code Style Guidelines

### Imports

Group imports in the following order (no blank lines between groups):

1. Standard library
2. External packages (github.com, etc.)

```go
import (
    "fmt"
    "os"

    "github.com/diogo/ghtools/internal/config"
    "github.com/spf13/cobra"
)
```

### Formatting

- Use `go fmt` before committing
- 4-space indentation (standard Go)
- Maximum line length: ~100 characters (soft limit)
- Use blank lines between logical sections in functions

### Naming Conventions

- **Variables**: camelCase (`cfg`, `verbose`, `quiet`)
- **Functions**: PascalCase / camelCase depending on export status
  - Exported: `RunChooseWithTitle`, `CheckInstalled`
  - Unexported: `runMenu`, `checkAuth`
- **Constants**: PascalCase for exported, camelCase for unexported
- **Package names**: short, lowercase, no underscores (`tui`, `gh`, `git`)
- **File names**: lowercase with underscores (`root.go`, `menu.go`)

### Error Handling

- Return errors explicitly with context:
  ```go
  if err != nil {
      return fmt.Errorf("description: %w", err)
  }
  ```
- For user-facing errors, use `tui.StyleError.Render()`:
  ```go
  fmt.Fprintln(os.Stderr, tui.StyleError.Render("ERROR")+" "+err.Error())
  ```
- Check prerequisites at the start of commands (e.g., `gh` CLI installed, authenticated)

### Type Usage

- Use explicit types for package-level variables
- Use `interface{}` sparingly; prefer concrete types
- For JSON parsing, use struct types with json tags

### TUI Patterns

When adding new TUI components:

1. Use `lipgloss` for styling (see `internal/tui/styles.go`)
2. Use `bubbletea` for interactive components
3. Use the existing helper functions in `internal/tui/`
4. Check `tui.Quiet` flag before printing informational messages

Example TUI style:
```go
var StyleMyComponent = lipgloss.NewStyle().
    Foreground(ColorPrimary).
    Bold(true).
    Padding(1, 2)
```

### Command Structure

When adding new Cobra commands:

1. Create new file in `cmd/` directory
2. Use the naming pattern: `verb.go` (e.g., `list.go`, `search.go`)
3. Register in `menu.go` switch statement
4. Follow existing pattern for `PersistentPreRunE` checks

---

## Architecture

```
ghtools/
├── main.go              # Entry point
├── cmd/                 # CLI commands (Cobra)
│   ├── root.go         # Root command, global flags
│   ├── menu.go         # Interactive menu
│   ├── list.go         # List repositories
│   └── ...
├── internal/
│   ├── cache/          # Caching functionality
│   ├── config/         # Configuration management
│   ├── gh/             # GitHub API wrapper (gh CLI)
│   ├── git/            # Git operations
│   ├── runner/         # Parallel execution
│   ├── template/       # Template handling
│   ├── types/          # Type definitions
│   └── tui/            # TUI components (Bubble Tea + Lipgloss)
└── _legacy/            # Deprecated bash script
```

---

## Common Operations

### Adding a New Command

1. Create `cmd/newcmd.go`
2. Define command using Cobra:
   ```go
   var newCmd = &cobra.Command{
       Use:   "newcmd",
       Short: "Description",
       RunE:  runNewCmd,
   }
   ```
3. Add flag handling in `init()`
4. Implement `runNewCmd` function
5. Register in `root.go`:
   ```go
   rootCmd.AddCommand(newCmd)
   ```
6. Add menu entry in `cmd/menu.go` switch

### Adding a New TUI Component

1. Add to appropriate file in `internal/tui/`
2. Export style variables in `styles.go`
3. Use existing patterns from similar components

---

## Configuration

- **Config file**: `~/.config/ghtools/config.json`
- **Cache directory**: `~/.cache/ghtools/`
- **Config is loaded via**: `config.Load()`

---

## Notes for Agents

- This project uses the `gh` CLI (must be authenticated)
- No test files exist - testing is manual
- No CI/CD currently configured
- Error messages should be user-friendly
- Prefer interactive TUI over flag-heavy CLI

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
