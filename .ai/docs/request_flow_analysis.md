# Request Flow Analysis
The `ghtools` project is a command-line interface (CLI) tool implemented as a single Bash script. The "request flow" is defined by the execution path of command-line arguments, which are parsed and dispatched to specific handler functions. The system relies heavily on the `gh` (GitHub CLI), `jq`, `fzf`, and `gum` tools for data interaction, processing, and user interface.

## Entry Points Overview
The sole entry point for the application logic is the `main` function in the `./ghtools` script.

1.  **Direct Command Execution:** `ghtools <COMMAND> [ARGS]`
    *   The `main` function parses global options (`-V`, `-q`) and then uses a `case` statement to dispatch the primary command (`$1`) to the corresponding `action_*` function.
2.  **Interactive Menu (TUI):** `ghtools` (with no arguments)
    *   The `main` function enters a `while true` loop, repeatedly calling `show_menu`.
    *   `show_menu` uses `gum choose` (or `fzf` fallback) to present a list of actions, and the user's selection is then mapped to the appropriate `action_*` function.

## Request Routing Map
Routing is handled by the `case` statement within the `main` function, which maps the command-line verb to a dedicated Bash function.

| Command Category | Command | Handler Function | Core GitHub CLI Command |
| :--- | :--- | :--- | :--- |
| **Repository Management** | `list` | `action_list` | `gh repo list` |
| | `clone` | `action_clone` | `gh repo clone` |
| | `create` | `action_create` | `gh repo create` |
| | `delete` | `action_delete` | `gh repo delete` |
| | `fork` | `action_fork` | `gh repo fork` |
| | `archive` | `action_archive` | `gh repo archive`/`unarchive` |
| | `visibility` | `action_visibility` | `gh repo edit` |
| **Local Operations** | `sync` | `action_sync` | `git pull --ff-only` |
| | `status` | `action_status` | `git branch`, `git diff-index`, `git rev-list` |
| **Discovery** | `search` | `action_search` | (Uses cached data) |
| | `browse` | `action_browse` | `gh browse` |
| | `stats` | `action_stats` | (Uses cached data) |
| | `explore` | `action_explore` | `gh search repos` |
| | `trending` | `action_trending` | `gh search repos` |
| **Pull Requests** | `pr list` | `action_pr_list` | `gh pr list` |
| | `pr create` | `action_pr_create` | `gh pr create` |
| **Utilities** | `refresh` | (Inline logic) | `rm -f $CACHE_FILE` |
| | `config` | `init_config` | `mkdir -p`, `cat > $CONFIG_FILE` |

## Middleware Pipeline
The script executes several preprocessing steps before command dispatch, acting as a middleware chain:

1.  **Configuration Loading:** `load_config` reads `$CONFIG_FILE` to set global variables (`CACHE_TTL`, `MAX_JOBS`, `DEFAULT_ORG`, etc.). It includes a security check to validate the config file contents.
2.  **Global Option Parsing:** The `main` function parses `-V` (`VERBOSE`) and `-q` (`QUIET`), setting global flags that affect all subsequent output functions (`print_*`).
3.  **Dependency Check:** `check_dependencies` verifies the presence of required tools (`gh`, `fzf`, `git`, `jq`) and warns about optional ones (`gum`). The script exits if required dependencies are missing.
4.  **Authentication Check:** `check_gh_auth` verifies that the user is authenticated with GitHub CLI (`gh auth status`). The script exits if authentication fails.

## Controller/Handler Analysis
The `action_*` functions serve as controllers, orchestrating the data flow and user interaction.

### Data Flow and Transformation
*   **Data Source:** The primary data source is the GitHub API, accessed via `gh repo list`.
*   **Caching:** `fetch_repositories_json` implements a time-to-live (TTL) cache (`$CACHE_FILE`, default 10 minutes). It handles fetching, secure file creation (`umask 077`), and cache validation.
*   **Transformation:** `jq` is the core transformation engine, used extensively within handlers (`action_list`, `action_stats`, `action_search`) to filter, project, and aggregate JSON data fetched from GitHub.
*   **User Interface:** `gum` (or `fzf` fallback) is used for interactive selection, filtering, and input, transforming raw data into user-selectable lists and capturing user choices.

### Parallelism
*   `action_clone` and `action_sync` implement parallelism using Bash job control (`&` for background execution and `wait_for_jobs`).
*   `wait_for_jobs` limits the number of concurrent operations to `$MAX_JOBS` (default 5), ensuring resource control during bulk operations.

### Request Validation
*   Handlers perform basic input validation, such as checking if a clone path exists (`action_clone`) or if the user is on a feature branch before creating a PR (`action_pr_create`).

## Authentication & Authorization Flow
Authentication is handled at two levels:

1.  **Global Authentication:** `check_gh_auth` ensures a valid `gh` session exists upon script startup.
2.  **Authorization Scope Check (Sensitive Operations):**
    *   The `action_delete` handler includes a specific authorization checkpoint: `check_delete_scope`.
    *   This function checks if the user's GitHub token has the necessary `delete_repo` scope by inspecting the output of `gh auth status`.
    *   If the scope is missing, the user is prompted to refresh their authentication using `gh auth refresh -s delete_repo`, ensuring the sensitive operation is only attempted with the correct permissions.

## Error Handling Pathways
The script employs a combination of global and localized error handling:

*   **Global Exit Strategy:** `set -euo pipefail` is set at the top of the script, ensuring that any command failure, use of an unset variable, or pipeline failure immediately terminates the script, preventing unexpected behavior.
*   **Output Reporting:** The `print_error` function is the standardized way to report failures to the user, often used after checking the exit status of a `gh` or `git` command (e.g., `if ! gh repo clone...; then print_error...`).
*   **API Failure Handling:** `fetch_repositories_json` explicitly captures and reports errors from the `gh repo list` command before exiting the script.
*   **Interactive Confirmation:** Sensitive actions like `action_delete` and `action_archive` require explicit user confirmation (typing 'DELETE' or 'y/N') before proceeding, mitigating accidental data loss.

## Request Lifecycle Diagram

```mermaid
graph TD
    A[Start: ghtools <CMD> [ARGS]] --> B{Global Options Parsing};
    B --> C[Load Config];
    C --> D[Check Dependencies];
    D --> E[Check gh Auth];
    E -- Success --> F{Command Provided?};
    F -- No --> G[Interactive Menu (show_menu)];
    F -- Yes --> H[Command Router (main case)];

    G --> I[User Selection];
    I --> H;

    H -- Dispatch --> J(Action Handler: action_*);

    J --> K{Needs Repo Data?};
    K -- Yes --> L[fetch_repositories_json];
    L --> M{Cache Valid?};
    M -- Yes --> N[Use Cache];
    M -- No --> O[gh repo list (API Call)];
    O -- Success --> P[Update Cache];
    O -- Failure --> Q[Print Error & Exit];
    N --> J;
    P --> J;

    J --> R{Sensitive Action?};
    R -- Yes (e.g., delete) --> S[Check Auth Scope (check_delete_scope)];
    S -- Fail --> Q;
    S -- Success --> T[User Confirmation];

    J --> U[Execute Core gh/git Command];
    T --> U;

    U -- Success --> V[print_success];
    U -- Failure --> W[print_error];

    V --> X[End/Return to Menu];
    W --> X;
    Q --> Z[Exit 1];

    style A fill:#f9f,stroke:#333
    style E fill:#ccf,stroke:#333
    style H fill:#bbf,stroke:#333
    style L fill:#ddf,stroke:#333
    style S fill:#fcc,stroke:#333
    style U fill:#afa,stroke:#333
    style Q fill:#f00,stroke:#333
```