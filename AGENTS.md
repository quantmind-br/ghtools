# Repository Guidelines

## Project Structure & Module Organization
The CLI entry point lives in `ghtools`, a Bash script that wires the interactive menu, direct subcommands, and all helper functions. Installation automation and PATH hygiene are handled in `install.sh`. Supporting docs and product specs sit under `PRPs/` and `SPEC_PRP/`; keep internal planning material there rather than alongside the executable script. There is no compiled output or build artifact committed to the repo—new assets should follow the same convention.

## Build, Test, and Development Commands
Run `./install.sh` to copy the CLI into `~/scripts` and ensure PATH hooks stay deduplicated. Use `./ghtools help` or `./ghtools list --help` while developing to verify argument parsing. Lint changes locally with `shellcheck ghtools` and `shellcheck install.sh`; resolve warnings before opening a PR. When iterating on command flows, call `./ghtools sync --dry-run --path <dir>` or `./ghtools delete` inside a sandbox repo to validate prompts without mutating production repositories.

## Coding Style & Naming Conventions
Scripts are POSIX-friendly but assume Bash; keep `#!/bin/bash` and `set -euo pipefail` at the top of new modules. Indent with four spaces, prefer snake_case for functions (`check_dependencies`) and ALL_CAPS for constants (color codes, paths). Echo user-facing messages through the provided `print_*` helpers so color formatting stays consistent. When adding flags, mirror the existing long-option style (`--dry-run`, `--max-depth`) and document them in both `show_usage` and the README.

## Testing Guidelines
There is no automated test harness yet, so exercise new logic manually with `gh auth status` confirmed and a GitHub test account when possible. Validate list/clone/sync flows against repositories with varied visibility to catch edge cases, and always include a `--dry-run` pathway for destructive commands. If a feature depends on GitHub scopes (e.g., `delete_repo`), add an explicit check similar to `check_delete_scope` and describe the requirement in the README update.

## Commit & Pull Request Guidelines
Commit messages follow a `Type: concise summary` pattern (see `Refactor: Replace ghclone and ghdelete with unified ghtools`). Keep subject lines under 72 characters and wrap additional context in the body when needed. Squash cosmetic commits before pushing. Pull requests should link related issues, summarize verification steps (`shellcheck`, manual command matrix), and note any auth scopes or environment impacts. Include screenshots or terminal captures only when they clarify interactive changes; otherwise reference exact commands exercised.

## Security & Configuration Notes
Authenticate early with `gh auth login` and refresh scopes via `gh auth refresh -s delete_repo` before testing deletion. Avoid hardcoding tokens or organization names; pass them through existing prompts or environment variables. When adding filesystem writes, respect the current `$HOME/scripts` install target and warn users before creating or overwriting files outside that directory.
