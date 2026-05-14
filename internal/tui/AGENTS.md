# TUI KNOWLEDGE BASE

**Generated:** 2025-05-14

## OVERVIEW

Bubble Tea + Lipgloss terminal UI components. All interactive prompts and styling live here.

## WHERE TO LOOK

| Task | File | Notes |
|------|------|-------|
| Styling/colors | `styles.go` | Color palette + pre-defined styles |
| Print helpers | `print.go` | Quiet-aware PrintError/Success/Info/Warning |
| Text input | `input.go` | Bubble Tea input model with header |
| Single select | `choose.go` | RunChoose / RunChooseWithTitle |
| Multi select | `multiselect.go` | RunMultiSelect with filter |
| Confirm prompt | `confirm.go` | Yes/no prompt |
| Spinner | `spinner.go` | RunWithSpinner for async tasks |
| Table printing | `table.go` | Dynamic column width calculation |

## CONVENTIONS

- All styles use `lipgloss.NewStyle()` - defined in `styles.go`
- Interactive components: struct model + Init/Update/View methods (Bubble Tea pattern)
- Export `Run*` functions, keep models unexported
- Check `Quiet` flag before non-error output in print helpers
- Colors: Primary=blue, Secondary=purple, Accent=green, Success=green, Warning=yellow, Error=red, Info=cyan, Muted=gray

## ANTI-PATTERNS

- Never use `fmt.Println` directly in cmd/ - use `tui.Print*` helpers
- Never create new styles outside `styles.go`
- Never skip `tui.Quiet` checks for informational output
