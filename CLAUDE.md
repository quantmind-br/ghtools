# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development

- **Type**: Bash script (no compilation required)
- **Run**: `./ghtools` or `bash ghtools`
- **Install**: `./install.sh` (installs to `~/scripts` and updates PATH)
- **Dependencies**: `gh` (GitHub CLI), `fzf`, `jq`, `git`
- **Test**: Manual testing of commands (e.g., `./ghtools list --dry-run` where applicable, though most commands interact with live GitHub API)

## Architecture

- **Core Script**: `ghtools` is a monolithic bash script.
- **State Management**:
  - Caches GitHub repo data in `/tmp/ghtools_repos.json` (TTL: 10 mins).
  - Environment variables control behavior (`MAX_JOBS`, `CACHE_TTL`).
- **Key Functions**:
  - `main`: Entry point, handles argument parsing and routing.
  - `action_*`: Implementations for specific commands (list, clone, sync, create, delete).
  - `fetch_repositories_json`: Handles API calls to GitHub and caching.
  - `wait_for_jobs`: Manages parallelism for bulk operations (clone/sync).
- **UI/UX**:
  - Uses `fzf` for interactive selection menus.
  - Uses ANSI escape codes for colored output (`print_error`, `print_success`, etc.).
  - Uses `jq` for parsing and formatting JSON data from GitHub API.
- **Installation**: `install.sh` handles dependency checking, file copying, and shell configuration (zsh support).

## Code Style

- **Strict Mode**: Scripts use `set -uo pipefail`.
- **Formatting**: Indentation is 2 spaces.
- **Output**: Use helper functions `print_info`, `print_success`, `print_warning`, `print_error` for consistency.
- **Safety**:
  - Destructive actions (delete) require explicit confirmation and scope checks.
  - Sync uses `pull --ff-only` to prevent accidental merges.
