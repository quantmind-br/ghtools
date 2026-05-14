# CMD KNOWLEDGE BASE

**Generated:** 2025-05-14

## OVERVIEW

Cobra CLI commands - each file implements one user-facing command.

## WHERE TO LOOK

| Task | File | Notes |
|------|------|-------|
| List repos | `list.go` | Table output with language filter |
| Search repos | `search.go` | GitHub search with table output |
| Clone repo | `clone.go` | Clones to configured path |
| Sync repos | `sync_cmd.go` | Parallel sync with spinner |
| Local status | `status.go` | Git status across repos |
| Fork repo | `fork.go` | Forks to user or org |
| Create repo | `create.go` | Creates with optional template |
| Delete repo | `delete.go` | Requires delete scope, confirms |
| Archive/Unarchive | `archive.go` | Toggle archive status |
| Change visibility | `visibility.go` | public/private/internal |
| Browse repo | `browse.go` | Opens browser |
| PR operations | `pr.go` | List or create PRs |
| Config management | `config_cmd.go` | View/edit settings |
| Cache refresh | `refresh.go` | Clears and rebuilds cache |
| Trending repos | `trending.go` | GitHub trending by language |
| Statistics | `stats.go` | Dashboard summary |
| Explore repo | `explore.go` | Interactive repo actions |
| Menu loop | `menu.go` | Main interactive navigation |
| Root command | `root.go` | Global flags, PersistentPreRunE |
| Helpers | `helpers.go` | reposToItems, searchResultsToItems |

## CONVENTIONS

- Each command: `var <name>Cmd` variable + `init()` + `run<Name>()` function
- Register in `root.go`: `rootCmd.AddCommand(<name>Cmd)`
- Add menu entry in `menu.go` switch statement
- Use `PersistentPreRunE` for prerequisite checks (gh auth, etc.)
- Check `tui.Quiet` before printing informational messages
- User-facing errors: `tui.StyleError.Render("ERROR") + err.Error()`

## ANTI-PATTERNS

- Never add business logic in cmd/ - delegate to internal/gh/ or internal/git/
- Never bypass the menu.go switch for new commands
- Never skip PersistentPreRunE checks
