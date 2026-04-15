# UI/UX Improvements Analysis Report

## Executive Summary

`ghtools` is a Go CLI with a Bubble Tea/Lip Gloss TUI layer, not a browser app. I analyzed the main user-facing flows in `cmd/` and the shared UI primitives in `internal/tui/`.

Biggest UX gaps sit in three places: long-running actions show weak progress feedback, selection/confirmation rely too heavily on color and implicit keyboard knowledge, and destructive flows do not provide enough structured reassurance before irreversible changes.

**Key statistics**

- Main UI surface: Charmbracelet Bubble Tea + Lip Gloss + Cobra CLI
- Total components analyzed: 22
- Total issues found: 9
- High priority: 4
- Medium priority: 3
- Low priority: 2

## Issues Found

### High Priority

#### UIUX-001: Long-running actions have weak progress feedback and inconsistent loading states

**Category:** performance

**Affected Components:**
- `internal/tui/spinner.go`
- `cmd/explore.go`
- `cmd/search.go`
- `cmd/clone.go`
- `cmd/sync_cmd.go`
- `cmd/pr.go`

**Current State:**
Network fetches and repo mutations mostly print one static info line, then block until work completes. `RunWithSpinner` already exists in [`internal/tui/spinner.go:60`](/home/diogo/dev/ghtools/internal/tui/spinner.go:60), but main flows do not use it. Clone/sync print raw `\r Progress: x/y` lines in [`cmd/clone.go:81`](/home/diogo/dev/ghtools/cmd/clone.go:81) and [`cmd/sync_cmd.go:103`](/home/diogo/dev/ghtools/cmd/sync_cmd.go:103), while fetch-heavy flows such as [`cmd/explore.go:53`](/home/diogo/dev/ghtools/cmd/explore.go:53) and [`cmd/pr.go:61`](/home/diogo/dev/ghtools/cmd/pr.go:61) have no visual waiting state beyond text.

**Proposed Change:**
Standardize on a single loading pattern:
- Use `tui.RunWithSpinner(...)` for fetch-only waits.
- Add a reusable batch-progress model for clone/sync/fork/delete/archive flows.
- Show current item, completed count, failed count, and cancel hint during long jobs.

**User Benefit:**
Users stop guessing whether app froze, especially on slow networks or large batch jobs.

**Code Example:**
```tsx
// Current
tui.PrintInfo(fmt.Sprintf("Searching GitHub for '%s' (sorted by %s)...", query, sort))
results, err := gh.SearchRepos(query, sort, lang, limit)

// Proposed
var results []types.SearchResult
err := tui.RunWithSpinner(
	 fmt.Sprintf("Searching GitHub for '%s'...", query),
	 func() error {
	 	 var err error
	 	 results, err = gh.SearchRepos(query, sort, lang, limit)
	 	 return err
	 },
)
```

**Estimated Effort:** small

---

#### UIUX-002: Selection and confirmation states rely too much on color, hurting accessibility in terminals

**Category:** accessibility

**Affected Components:**
- `internal/tui/styles.go`
- `internal/tui/choose.go`
- `internal/tui/multiselect.go`
- `internal/tui/confirm.go`

**Current State:**
Primary state changes depend heavily on purple/cyan/pink/gray foreground swaps in [`internal/tui/styles.go:6`](/home/diogo/dev/ghtools/internal/tui/styles.go:6). Active choice uses colored `> ` plus colored text in [`internal/tui/choose.go:94`](/home/diogo/dev/ghtools/internal/tui/choose.go:94), and confirm uses green/red emphasis in [`internal/tui/confirm.go:56`](/home/diogo/dev/ghtools/internal/tui/confirm.go:56). On low-color terminals, custom themes, or for color-blind users, state distinction drops fast.

**Proposed Change:**
- Add non-color cues: bold, underline, inverse background, explicit labels like `(selected)` or `[*]`.
- Detect `NO_COLOR` / low-color terminals and switch to high-contrast fallback tokens.
- Use one semantic helper for focus, selected, danger, muted instead of per-screen ad hoc styling.

**User Benefit:**
Critical actions stay understandable even without reliable color rendering.

**Code Example:**
```tsx
// Current
if i == m.cursor {
	cursor = StyleAccent.Render("> ")
	opt = StyleSecondary.Render(opt)
}

// Proposed
if i == m.cursor {
	cursor = "> "
	opt = FocusStyle.Render("[selected] " + opt)
}
```

**Estimated Effort:** small

---

#### UIUX-003: Main menu flow adds friction after every action and drops user context

**Category:** usability

**Affected Components:**
- `cmd/menu.go`
- `internal/tui/choose.go`
- `internal/tui/confirm.go`

**Current State:**
After each action, menu exits into a second confirmation step, [`tui.RunConfirm("Continue?", true)` in `cmd/menu.go:85`](/home/diogo/dev/ghtools/cmd/menu.go:85). This adds one extra keypress to every flow and discards menu position, so repeated admin work feels slower than necessary.

**Proposed Change:**
- Keep menu open by default after command completion.
- Preserve last cursor position in `RunChooseWithTitle`.
- Replace generic `Continue?` with footer help inside menu: `Enter=open • q=quit`.

**User Benefit:**
High-frequency workflows become faster and feel like one tool, not many disconnected prompts.

**Code Example:**
```tsx
// Current
cont, err := tui.RunConfirm("Continue?", true)
if err != nil || !cont {
	return nil
}

// Proposed
lastCursor = menuCursor
continue // return directly to menu with preserved selection
```

**Estimated Effort:** small

---

#### UIUX-004: Destructive flows do not give enough structured risk review before execution

**Category:** usability

**Affected Components:**
- `cmd/delete.go`
- `cmd/search.go`
- `cmd/archive.go`
- `cmd/visibility.go`

**Current State:**
Delete/archive/visibility flows print plain text summaries, then ask for `Continue?` or typed `DELETE` in a generic input box. Examples: [`cmd/delete.go:47`](/home/diogo/dev/ghtools/cmd/delete.go:47), [`cmd/search.go:73`](/home/diogo/dev/ghtools/cmd/search.go:73), [`cmd/archive.go:76`](/home/diogo/dev/ghtools/cmd/archive.go:76), [`cmd/visibility.go:104`](/home/diogo/dev/ghtools/cmd/visibility.go:104). There is no dedicated danger screen, no highlighted consequence summary, no clear default-safe choice beyond convention.

**Proposed Change:**
- Build one reusable `RunDangerConfirm` primitive.
- Show action, repo count, sample repo names, reversibility, and cache side effects.
- Default focus to cancel.
- Use typed confirmation only for delete; keep archive/visibility on explicit yes/no.

**User Benefit:**
Reduces accidental destructive changes and increases user trust in admin operations.

**Code Example:**
```tsx
// Current
confirm, err := tui.RunInput("Type 'DELETE' to confirm, or anything else for Dry Run", "DELETE", "")

// Proposed
confirm, err := tui.RunDangerConfirm(tui.DangerConfirm{
	Title: "Delete repositories",
	Count: len(selected),
	Items: previewNames(selected, 5),
	Consequence: "Permanent GitHub deletion. Cache will be cleared.",
	RequirePhrase: "DELETE",
	DefaultAction: "cancel",
})
```

**Estimated Effort:** medium

---

### Medium Priority

#### UIUX-005: Multi-select screen hides important interaction rules and filtered context

**Category:** interaction

**Affected Components:**
- `internal/tui/multiselect.go`
- `cmd/search.go`
- `cmd/explore.go`
- `cmd/clone.go`

**Current State:**
Multi-select shows basic shortcut text in [`internal/tui/multiselect.go:145`](/home/diogo/dev/ghtools/internal/tui/multiselect.go:145), but it does not show `Esc` to cancel, does not show filtered count vs total visible count, and auto-advances cursor after toggle in [`internal/tui/multiselect.go:79`](/home/diogo/dev/ghtools/internal/tui/multiselect.go:79), which can feel jumpy during careful selection.

**Proposed Change:**
- Footer should show full help: `↑↓ move • space toggle • enter confirm • esc cancel • ctrl+a all • ctrl+d clear`.
- Show `visible/total` count and active filter query.
- Make cursor auto-advance optional or remove it.

**User Benefit:**
Selection becomes easier to learn and less error-prone, especially for first-time users.

**Code Example:**
```tsx
// Current
b.WriteString(countStyle.Render(fmt.Sprintf("  %d/%d selected  Tab:toggle  Ctrl+A:all  Enter:confirm", selectedCount, len(m.items))) + "\n")

// Proposed
b.WriteString(countStyle.Render(
	fmt.Sprintf("  %d selected • %d visible / %d total • Space: toggle • Enter: confirm • Esc: cancel",
		selectedCount, len(m.filtered), len(m.items)),
) + "\n")
```

**Estimated Effort:** small

---

#### UIUX-006: Fixed-width tables and byte-based truncation reduce readability on narrow terminals and Unicode repo names

**Category:** accessibility

**Affected Components:**
- `internal/tui/table.go`
- `cmd/list.go`
- `cmd/status.go`
- `cmd/helpers.go`

**Current State:**
Tables use hardcoded widths in [`cmd/list.go:51`](/home/diogo/dev/ghtools/cmd/list.go:51) and [`cmd/status.go:43`](/home/diogo/dev/ghtools/cmd/status.go:43). `Truncate` slices raw bytes in [`internal/tui/table.go:8`](/home/diogo/dev/ghtools/internal/tui/table.go:8), which risks broken glyphs and poor rendering for multi-byte text. Long repo names and descriptions become hard to scan on small terminals.

**Proposed Change:**
- Switch truncation to display-width aware helpers (`lipgloss.Width`, runewidth, or Bubble Tea width helpers).
- Collapse secondary columns on narrow terminals.
- Prefer stacked mobile-style rows when terminal width is below threshold.

**User Benefit:**
Improves readability for small terminal windows and non-ASCII repo names.

**Code Example:**
```tsx
// Current
return s[:max-3] + "..."

// Proposed
return truncateDisplayWidth(s, max) // width-aware, rune-safe
```

**Estimated Effort:** medium

---

#### UIUX-007: Empty states tell users nothing about recovery or next step

**Category:** state

**Affected Components:**
- `cmd/list.go`
- `cmd/explore.go`
- `cmd/fork.go`
- `cmd/pr.go`
- `cmd/visibility.go`
- `cmd/archive.go`

**Current State:**
Most empty states stop at messages like `No repositories found`, `No open PRs found`, or `No repositories match the criteria` in [`cmd/list.go:47`](/home/diogo/dev/ghtools/cmd/list.go:47), [`cmd/explore.go:57`](/home/diogo/dev/ghtools/cmd/explore.go:57), [`cmd/pr.go:69`](/home/diogo/dev/ghtools/cmd/pr.go:69), and [`cmd/visibility.go:61`](/home/diogo/dev/ghtools/cmd/visibility.go:61). They do not suggest filter changes, refresh actions, or alternate commands.

**Proposed Change:**
- Standardize empty states with reason + next step.
- Example: suggest `--refresh`, `config`, different org, different query, or creating first repo.
- Add `EmptyState(title, hint, action)` helper for consistency.

**User Benefit:**
Users recover faster instead of guessing whether tool failed or data truly empty.

**Code Example:**
```tsx
// Current
tui.PrintWarning("No repositories found")

// Proposed
tui.PrintEmpty(
	"No repositories found",
	"Try --refresh, change --org, or create a new repository first.",
)
```

**Estimated Effort:** small

---

### Low Priority

#### UIUX-008: Statistics dashboard has weak hierarchy and low scan efficiency

**Category:** visual

**Affected Components:**
- `cmd/stats.go`
- `internal/tui/styles.go`

**Current State:**
Stats view mixes boxed cards with plain text section headers like `--- Languages Breakdown ---` in [`cmd/stats.go:83`](/home/diogo/dev/ghtools/cmd/stats.go:83). Important metrics are not visually prioritized beyond box border color, and sections do not align into a clear layout.

**Proposed Change:**
- Convert dashboard into consistent cards.
- Use shared spacing tokens from `internal/tui/styles.go`.
- Add bars or ranked pills for top languages/repos.
- Group headline metrics into a single responsive top row.

**User Benefit:**
Dashboard becomes faster to scan and feels more polished.

**Code Example:**
```tsx
// Current
fmt.Println(tui.StyleMuted.Render("--- Top Repositories (by Stars) ---"))

// Proposed
fmt.Println(tui.RenderSection("Top Repositories", renderTopRepoList(repos)))
```

**Estimated Effort:** medium

---

#### UIUX-009: Batch operations flood line-by-line logs but do not end with a strong summary or retry guidance

**Category:** usability

**Affected Components:**
- `cmd/clone.go`
- `cmd/sync_cmd.go`
- `cmd/search.go`
- `cmd/explore.go`
- `cmd/archive.go`
- `cmd/visibility.go`

**Current State:**
After batch work, app prints one result per repo, then one generic success line such as [`cmd/clone.go:94`](/home/diogo/dev/ghtools/cmd/clone.go:94) or [`cmd/sync_cmd.go:116`](/home/diogo/dev/ghtools/cmd/sync_cmd.go:116). Failed items are not regrouped into a final summary, so larger jobs are harder to review.

**Proposed Change:**
- End each batch with summary counts: succeeded, skipped, failed.
- If failures exist, print retry-ready repo list or suggested follow-up command.
- Reserve green `OK` closeout only when zero failures remain.

**User Benefit:**
Users can act on results immediately without rescanning long terminal output.

**Code Example:**
```tsx
// Proposed
tui.PrintBatchSummary(tui.BatchSummary{
	Succeeded: okRepos,
	Skipped: skippedRepos,
	Failed: failedRepos,
	RetryHint: "Retry failed items with: ghtools clone --path ...",
})
```

**Estimated Effort:** small

---

## Summary

| Category | Count |
|----------|-------|
| Usability | 3 |
| Accessibility | 2 |
| Performance Perception | 1 |
| Visual Polish | 1 |
| Interaction | 1 |
| State Handling | 1 |

**Total Components Analyzed:** 22
**Total Issues Found:** 9

## Notes On Scope

- Repo has no web frontend, HTML, CSS, or browser accessibility surface.
- Analysis focused on terminal UX: discoverability, keyboard flow, color dependency, feedback, readability, and safety of destructive actions.
