This project, `ghtools`, is a **Command-Line Interface (CLI) tool** written in **Bash**. It does not expose a traditional HTTP API (like REST or GraphQL). Instead, its "API" is the set of commands it provides to the user, and it acts as a **client** to the **GitHub API** via the `gh` (GitHub CLI) tool.

The documentation below treats the CLI commands as the "APIs Served" and the underlying GitHub interactions as "External API Dependencies."

# API Documentation

## APIs Served by This Project

The service exposes a unified command-line interface for managing GitHub repositories. All interactions are performed via the `ghtools` executable.

### Endpoints (CLI Commands)

| Command | Method and Path | Description |
| :--- | :--- | :--- |
| `list` | `ghtools list` | Lists repositories, optionally filtered by language or organization. Uses a local cache. |
| `clone` | `ghtools clone` | Interactively selects and clones repositories. Supports parallel cloning. |
| `create` | `ghtools create` | Creates a new GitHub repository interactively. |
| `delete` | `ghtools delete` | Interactively selects and deletes repositories. |
| `fork` | `ghtools fork <query>` | Forks an external repository found via search. |
| `archive` | `ghtools archive` | Archives or unarchives selected repositories. |
| `visibility` | `ghtools visibility` | Changes the visibility (public/private) of selected repositories. |
| `sync` | `ghtools sync` | Synchronizes local repositories with their remotes (pulls changes). Supports parallelism. |
| `status` | `ghtools status` | Shows the git status of local repositories. |
| `search` | `ghtools search` | Interactive fuzzy search across local and remote repositories with quick actions. |
| `browse` | `ghtools browse` | Opens selected repositories in the web browser. |
| `stats` | `ghtools stats` | Displays repository statistics dashboard. |
| `explore` | `ghtools explore <query>` | Searches external GitHub repositories. |
| `trending` | `ghtools trending` | Shows trending GitHub repositories. |
| `pr list` | `ghtools pr list` | Lists Pull Requests for a repository. |
| `pr create` | `ghtools pr create` | Creates a Pull Request from the current branch. |
| `refresh` | `ghtools refresh` | Clears the local repository cache. |
| `config` | `ghtools config` | Initializes or shows the location of the configuration file. |

#### Endpoint Details: `list`

*   **Method and Path:** `ghtools list [--refresh] [--lang <lang>] [--org <org>]`
*   **Description:** Fetches and displays a list of repositories accessible to the authenticated user. Results are cached for performance.
*   **Request (Parameters):**
    *   `--refresh` (Flag): Forces a refresh of the repository cache, ignoring the `CACHE_TTL`.
    *   `--lang <lang>` (String): Filters the list to repositories matching the specified primary language.
    *   `--org <org>` (String): Filters the list to repositories belonging to the specified organization. Overrides `DEFAULT_ORG` if set.
*   **Response (Success Format):** A formatted table printed to standard output, including `NAME`, `DESCRIPTION`, `VISIBILITY`, `LANG`, and `UPDATED` date.
*   **Response (Error Format):** Prints an error message to `stderr` if the GitHub API call fails or if authentication is missing.
*   **Authentication:** Relies on `gh CLI` authentication (see below).
*   **Examples:**
    ```bash
    ghtools list --lang python
    ghtools list --refresh --org my-company
    ```

#### Endpoint Details: `clone`

*   **Method and Path:** `ghtools clone [--path <dir>]`
*   **Description:** Presents an interactive fuzzy-search list of remote repositories for selection and cloning.
*   **Request (Parameters):**
    *   `--path <dir>` (String): Specifies the target directory for cloning. Overrides `DEFAULT_CLONE_PATH`.
*   **Authentication:** Relies on `gh CLI` authentication.
*   **Resilience:** Uses `MAX_JOBS` (default 5) to limit the number of parallel clone operations.

### Authentication & Security

The `ghtools` application does not manage its own authentication. It relies entirely on the **GitHub CLI (`gh`)** for all API interactions.

*   **Mechanism:** The script calls `gh auth status` to verify the user is logged in.
*   **Requirement:** A user must run `gh auth login` and successfully authenticate with GitHub (typically using a Personal Access Token or web flow) before using `ghtools`.
*   **Security:** The script uses `umask 077` when writing the cache file (`CACHE_FILE`) to ensure it has secure permissions (600), preventing other users from reading the potentially sensitive repository list.

### Rate Limiting & Constraints

The tool's rate limiting is inherited from the GitHub API limits applied to the authenticated user's token.

*   **Caching:** A local cache mechanism is implemented to reduce API calls for the repository list:
    *   **Cache File:** `/tmp/ghtools_repos_$(id -u).json`
    *   **Cache TTL:** Configurable via `CACHE_TTL` (default: 600 seconds / 10 minutes). The `list` command will only call the GitHub API if the cache is expired or if the `--refresh` flag is used.
*   **Parallelism:** Operations like `clone` and `sync` are parallelized using background jobs, limited by the configurable variable `MAX_JOBS` (default: 5).

## External API Dependencies

The project's core functionality is built upon consuming the GitHub API via the `gh` CLI tool.

### Services Consumed

#### 1. GitHub API (via `gh` CLI)

*   **Service Name & Purpose:** GitHub API. Used for fetching repository metadata, performing repository management actions (create, delete, archive, visibility), and managing Pull Requests.
*   **Base URL/Configuration:** The base URL is implicitly managed by the `gh` CLI, typically `https://api.github.com`.
*   **Endpoints Used (Mapped from `gh` commands):**

| `ghtools` Command | `gh` CLI Command | Underlying API Operation |
| :--- | :--- | :--- |
| `list`, `clone`, `search` | `gh repo list` | Fetching repository metadata (likely using GraphQL for efficiency). |
| `create` | `gh repo create` | Creating a new repository (REST POST /user/repos or /orgs/{org}/repos). |
| `delete` | `gh repo delete` | Deleting a repository (REST DELETE /repos/{owner}/{repo}). |
| `fork` | `gh repo fork` | Forking a repository (REST POST /repos/{owner}/{repo}/forks). |
| `archive` | `gh api repos/{owner}/{repo} -X PATCH` | Updating repository properties (e.g., `is_archived`). |
| `visibility` | `gh api repos/{owner}/{repo} -X PATCH` | Updating repository properties (e.g., `visibility`). |
| `explore`, `trending` | `gh search repos`, `gh api` | Searching and fetching trending data. |
| `pr list`, `pr create` | `gh pr list`, `gh pr create` | Pull Request management. |

*   **Authentication Method:** Inherited from `gh CLI` (Personal Access Token).
*   **Error Handling:**
    *   The `fetch_repositories_json` function uses `mktemp` to capture `stderr` from the `gh` command.
    *   If the `gh` command fails (e.g., due to network issues, rate limiting, or invalid credentials), the script prints a detailed error message to `stderr` and exits with a non-zero status code.
*   **Retry/Circuit Breaker Configuration:** None explicitly implemented in the Bash script. It relies on the underlying resilience of the `gh` CLI tool and the user's ability to manually re-run the command.

### Integration Patterns

*   **Data Transformation:** The script uses `jq` extensively to parse the JSON output from the `gh` CLI and transform it into a user-friendly, tabular format for display.
*   **Interactive UI:** The script uses `gum` (or falls back to `fzf`) for interactive selection, filtering, and input, providing a modern, user-friendly interface over the raw CLI commands.
*   **Idempotency:** The `clone` and `sync` operations are designed to be idempotent, ensuring that re-running them does not cause issues (e.g., `git clone` will fail if the directory exists, which is handled by the script's logic).

## Available Documentation

The project is primarily documented through its internal usage message and external markdown files.

| Path | Description | Quality Evaluation |
| :--- | :--- | :--- |
| `./ghtools` (Internal) | Contains the `show_usage` function, which serves as the primary command reference and help documentation. | **High.** Comprehensive list of commands, options, and examples. |
| `./README.md` | General project overview and setup instructions. | **Assumed High.** Standard entry point for new users. |
| `./.ai/docs/` | Directory containing additional documentation for AI agents. | **Assumed High.** Provides deeper context for codebase analysis. |
| `./AGENTS.md` | Documentation related to AI agents used in the project. | **Assumed High.** Relevant for understanding the development workflow. |