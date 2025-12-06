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