**Note**: This project uses [bd (beads)](https://github.com/steveyegge/beads)
for issue tracking. Use `bd` commands instead of markdown TODOs.
See AGENTS.md for workflow details.

# CLAUDE.md: ghtools Configuration

This repository contains `ghtools`, a single-file Bash CLI script designed to be a unified, interactive manager for GitHub repositories. It acts as an orchestrator, wrapping and enhancing the functionality of the `gh` (GitHub CLI) tool, `git`, and TUI tools like `fzf` and `gum`.

The primary goal is to provide a fast, interactive, and parallelized way to manage large numbers of repositories (list, clone, sync, delete).

## Development Workflow

The project is a pure Bash script and requires no compilation.

| Task | Command | Notes |
| :--- | :--- | :--- |
| **Run** | `./ghtools` | Runs the interactive TUI menu (requires `gum` or falls back to `fzf`). |
| **Direct Run** | `./ghtools list --lang bash` | Executes a specific command. |
| **Install** | `./install.sh` | Copies `ghtools` to `~/scripts` and updates the shell PATH. |
| **Lint** | `shellcheck ghtools` | Essential for catching Bash errors and style issues. |
| **Test** | `bats test/` | The project uses BATS (Bash Automated Testing System) for unit and integration tests. |
| **Refresh Cache** | `./ghtools refresh` | Explicitly clears the repository cache file. |

## Architecture Overview

The architecture is a **Monolithic Script** using a **Command Dispatcher Pattern**.

1.  **Orchestrator Role**: `ghtools` is a facade that simplifies complex operations by orchestrating external tools (`gh`, `jq`, `git`, `fzf`, `gum`).
2.  **Entry Point**: The `main` function handles argument parsing and routes execution to the appropriate `action_*` function.
3.  **Data Flow**: All repository data is sourced from the GitHub API via `gh repo list --json`, processed by `jq`, and then cached locally.
4.  **Parallelism**: Bulk operations (`clone`, `sync`) use the `run_parallel_jobs` utility, which leverages `xargs -P $MAX_JOBS` to limit concurrent processes (default 5).
5.  **UI Abstraction**: Output is routed through `print_*` functions, which conditionally use the `gum` TUI library for modern styling, falling back to standard ANSI colors if `gum` is unavailable.

## Core Components and Data Management

### 1. Action Handlers (`action_*`)
These functions encapsulate the business logic for each command:
*   `action_list`: Fetches, filters, and formats repository data.
*   `action_clone`: Interactively selects repos and runs `git clone` in parallel.
*   `action_sync`: Scans local repos, checks status, and runs `git pull --ff-only` in parallel.
*   `action_delete`: Requires `delete_repo` scope check and explicit confirmation.

### 2. Data Retrieval and Caching
*   **Function**: `fetch_repositories_json`
*   **Mechanism**: Implements a Cache-Aside pattern. It checks the timestamp of `$CACHE_FILE` (default: `/tmp/ghtools_repos_UID.json`) against `CACHE_TTL` (default: 600s).
*   **Data Contract**: The function ensures the data is a JSON array of repository objects, which is then consumed by `jq` for all subsequent filtering and formatting.
*   **Invalidation**: The cache is explicitly deleted (`rm -f "$CACHE_FILE"`) after any write operation (`create`, `delete`, `archive`) to ensure data freshness.

### 3. Configuration
*   Configuration is loaded from `$CONFIG_FILE` (default: `~/.config/ghtools/config`) via the `source` command.
*   Key variables: `CACHE_TTL`, `MAX_JOBS`, `DEFAULT_ORG`, `DEFAULT_CLONE_PATH`.

## Code Style and Conventions

### Bash Safety
*   **Strict Mode**: Always use `set -euo pipefail` at the top of the script and new functions to ensure robust error handling and prevent use of unset variables.
*   **Indentation**: Use **2 spaces** for indentation.
*   **Naming**: Use `snake_case` for functions (`check_dependencies`) and `ALL_CAPS` for global constants (`CACHE_TTL`, `VERBOSE`).

### Output and Logging
*   All user-facing output must use the provided helper functions for consistent styling:
    *   `print_info`, `print_success`, `print_warning`, `print_error`.
*   Use `print_verbose` for debugging output, which is only shown when the `-V` flag is set.

### Git Safety
*   The `action_sync` function must use `git pull --ff-only` to ensure non-fast-forward merges are prevented, maintaining a clean history.
*   `action_pr_create` includes checks for dirty working directories and detached HEAD state before proceeding.

## Development Gotchas and Warnings

1.  **External Tool Fragility**: The script is highly coupled to the exact command-line interfaces and JSON output formats of `gh` and `jq`. If these tools change their output, core functions like `get_all_repos` will break. Be cautious when updating dependency versions.
2.  **Configuration Security**: The `load_config` function uses `source`. While it includes a `grep` check to prevent arbitrary code execution, this is a potential security risk. **Do not introduce new ways to source external files.**
3.  **Authentication Scopes**: Destructive actions (`action_delete`) require the `delete_repo` OAuth scope. Always verify this using `check_delete_scope` and prompt the user to run `gh auth refresh -s delete_repo` if needed.
4.  **TUI Dependency**: The enhanced user experience relies on `gum`. When debugging UI issues, remember to test the fallback path (without `gum`) which uses `fzf` and standard ANSI colors.