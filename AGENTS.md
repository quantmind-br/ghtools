# Repository Guidelines

## Project Structure & Architecture
The project is a **Monolithic Bash CLI script** (`ghtools`) acting as an **Orchestrator** for external tools. It uses a Command Dispatcher pattern to route CLI arguments to `action_*` functions. Core data flow involves: `gh` (API) -> `jq` (Transform) -> `fzf`/`gum` (UI/Select).

**Key Components:**
*   `ghtools`: Main executable, handles routing and core logic.
*   `fetch_repositories_json`: Manages the local repository cache (`$CACHE_FILE`).
*   `run_parallel_jobs`: Utility for concurrent execution (e.g., cloning, syncing).
*   `print_*`: Functions for consistent, styled terminal output.

## Build, Test, and Development Commands
The script requires no compilation. Dependencies: `gh`, `jq`, `fzf`, `git`.

| Task | Command | Notes |
| :--- | :--- | :--- |
| **Install** | `bash install.sh` | Copies `ghtools` to `~/scripts` and updates PATH. |
| **Run** | `./ghtools` | Executes the interactive menu. |
| **Lint** | `shellcheck ghtools` | **Required** before any commit. |
| **Test** | `bats test/` | Runs the BATS test suite. |
| **Dry-Run Sync** | `./ghtools sync --dry-run` | Use for safe validation of sync logic. |

## Coding Style & Naming Conventions
*   **Strict Mode**: All scripts must start with `set -euo pipefail`.
*   **Indentation**: Use **4 spaces**.
*   **Naming**: `snake_case` for functions (`action_clone`), `ALL_CAPS` for constants (`MAX_JOBS`, `CACHE_TTL`).
*   **Output**: Use `print_info`, `print_success`, `print_error` helpers for all user-facing messages.
*   **Data**: Use `jq` for all JSON parsing and manipulation of `gh` output.

## Testing Guidelines
Testing is primarily manual and via the BATS suite.
1.  Verify `gh auth status` is successful before testing.
2.  Always test destructive commands (`delete`, `archive`) in a sandbox environment and ensure the `check_delete_scope` and confirmation prompts are working.
3.  Validate parallel operations (`clone`, `sync`) respect the `MAX_JOBS` limit.
4.  Ensure the cache is correctly invalidated (`rm -f $CACHE_FILE`) after write operations.

## Git Workflows
*   **Commit Format**: `Type: concise summary` (e.g., `Feat: Add --lang filter to list command`). Keep subject under 72 chars.
*   **PR Process**: Squash cosmetic commits. Link related issues. Summarize verification steps, including `shellcheck` results and manual command matrix.

## Security & Configuration Notes
*   **Authentication**: Must be handled by `gh auth login`. Sensitive actions require refreshing scopes: `gh auth refresh -s delete_repo`.
*   **Configuration**: Settings are loaded from `~/.config/ghtools/config`. Variables like `CACHE_TTL` and `MAX_JOBS` control runtime behavior.
*   **Safety**: `action_sync` must use `git pull --ff-only` to prevent non-fast-forward merges.
*   **Cache**: The cache file is secured with `umask 077` (permissions 600).

## Issue Tracking with bd (beads)

**IMPORTANT**: This project uses **bd (beads)** for ALL issue tracking. Do NOT use markdown TODOs, task lists, or other tracking methods.

### Why bd?

- Dependency-aware: Track blockers and relationships between issues
- Git-friendly: Auto-syncs to JSONL for version control
- Agent-optimized: JSON output, ready work detection, discovered-from links
- Prevents duplicate tracking systems and confusion

### Quick Start

**Check for ready work:**
```bash
bd ready --json
```

**Create new issues:**
```bash
bd create "Issue title" -t bug|feature|task -p 0-4 --json
bd create "Issue title" -p 1 --deps discovered-from:bd-123 --json
bd create "Subtask" --parent <epic-id> --json  # Hierarchical subtask (gets ID like epic-id.1)
```

**Claim and update:**
```bash
bd update bd-42 --status in_progress --json
bd update bd-42 --priority 1 --json
```

**Complete work:**
```bash
bd close bd-42 --reason "Completed" --json
```

### Issue Types

- `bug` - Something broken
- `feature` - New functionality
- `task` - Work item (tests, docs, refactoring)
- `epic` - Large feature with subtasks
- `chore` - Maintenance (dependencies, tooling)

### Priorities

- `0` - Critical (security, data loss, broken builds)
- `1` - High (major features, important bugs)
- `2` - Medium (default, nice-to-have)
- `3` - Low (polish, optimization)
- `4` - Backlog (future ideas)

### Workflow for AI Agents

1. **Check ready work**: `bd ready` shows unblocked issues
2. **Claim your task**: `bd update <id> --status in_progress`
3. **Work on it**: Implement, test, document
4. **Discover new work?** Create linked issue:
   - `bd create "Found bug" -p 1 --deps discovered-from:<parent-id>`
5. **Complete**: `bd close <id> --reason "Done"`
6. **Commit together**: Always commit the `.beads/issues.jsonl` file together with the code changes so issue state stays in sync with code state

### Auto-Sync

bd automatically syncs with git:
- Exports to `.beads/issues.jsonl` after changes (5s debounce)
- Imports from JSONL when newer (e.g., after `git pull`)
- No manual export/import needed!

### GitHub Copilot Integration

If using GitHub Copilot, also create `.github/copilot-instructions.md` for automatic instruction loading.
Run `bd onboard` to get the content, or see step 2 of the onboard instructions.

### MCP Server (Recommended)

If using Claude or MCP-compatible clients, install the beads MCP server:

```bash
pip install beads-mcp
```

Add to MCP config (e.g., `~/.config/claude/config.json`):
```json
{
  "beads": {
    "command": "beads-mcp",
    "args": []
  }
}
```

Then use `mcp__beads__*` functions instead of CLI commands.

### Managing AI-Generated Planning Documents

AI assistants often create planning and design documents during development:
- PLAN.md, IMPLEMENTATION.md, ARCHITECTURE.md
- DESIGN.md, CODEBASE_SUMMARY.md, INTEGRATION_PLAN.md
- TESTING_GUIDE.md, TECHNICAL_DESIGN.md, and similar files

**Best Practice: Use a dedicated directory for these ephemeral files**

**Recommended approach:**
- Create a `history/` directory in the project root
- Store ALL AI-generated planning/design docs in `history/`
- Keep the repository root clean and focused on permanent project files
- Only access `history/` when explicitly asked to review past planning

**Example .gitignore entry (optional):**
```
# AI planning documents (ephemeral)
history/
```

**Benefits:**
- ✅ Clean repository root
- ✅ Clear separation between ephemeral and permanent documentation
- ✅ Easy to exclude from version control if desired
- ✅ Preserves planning history for archeological research
- ✅ Reduces noise when browsing the project

### CLI Help

Run `bd <command> --help` to see all available flags for any command.
For example: `bd create --help` shows `--parent`, `--deps`, `--assignee`, etc.

### Important Rules

- ✅ Use bd for ALL task tracking
- ✅ Always use `--json` flag for programmatic use
- ✅ Link discovered work with `discovered-from` dependencies
- ✅ Check `bd ready` before asking "what should I work on?"
- ✅ Store AI planning docs in `history/` directory
- ✅ Run `bd <cmd> --help` to discover available flags
- ❌ Do NOT create markdown TODO lists
- ❌ Do NOT use external issue trackers
- ❌ Do NOT duplicate tracking systems
- ❌ Do NOT clutter repo root with planning documents

For more details, see README.md and QUICKSTART.md.