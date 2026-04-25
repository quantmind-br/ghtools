# AGENTS.md - Agentic Coding Guidelines for ghtools

## Project Overview

**ghtools** is a Go CLI tool that provides an interactive TUI for managing GitHub repositories. It wraps the `gh` CLI and provides menu-driven access to repository operations.

- **Go Version**: 1.25+
- **Dependencies**: Cobra v1.10.2, Bubble Tea v1.3.10, Lipgloss v1.1.0

---

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