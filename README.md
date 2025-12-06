# ghtools - Unified GitHub Repository Management Tool

## Project Overview

The `ghtools` project is a powerful, single-file **Bash script** designed to simplify and accelerate the management of GitHub repositories directly from the command line. It acts as an interactive, feature-rich wrapper and orchestrator for essential command-line tools like the GitHub CLI (`gh`), `git`, `fzf`, and `jq`.

### Purpose and Main Functionality

The primary purpose of `ghtools` is to provide a unified, interactive terminal user interface (TUI) for common and bulk GitHub operations, abstracting away complex API calls and repetitive command sequences.

### Key Features and Capabilities

*   **Interactive TUI:** Utilizes `gum` or `fzf` for fuzzy-finding, selection, and interactive prompts across all major commands.
*   **Bulk Operations:** Supports parallel execution for cloning (`clone`) and synchronizing (`sync`) multiple repositories, limited by a configurable `MAX_JOBS` setting.
*   **Performance Caching:** Implements a time-to-live (TTL) file-based cache for the repository list, significantly reducing repeated calls to the GitHub API.
*   **Comprehensive Actions:** Provides dedicated commands for listing, cloning, syncing, creating, deleting, archiving, and changing the visibility of repositories.
*   **Local Status Checks:** The `sync` and `status` commands provide detailed information on the local state of repositories (e.g., dirty working directory, ahead/behind remote).

### Likely Intended Use Cases

*   **Developer Onboarding:** Quickly clone a large set of required repositories for a new project or team.
*   **Daily Workflow:** Rapidly check the status and pull updates for all local repositories with a single command (`ghtools sync`).
*   **Repository Administration:** Securely and interactively delete or archive multiple repositories.
*   **Discovery:** Search, explore, and view trending repositories directly in the terminal.

## Table of Contents

1.  [Project Overview](#project-overview)
2.  [Architecture](#architecture)
3.  [C4 Model Architecture](#c4-model-architecture)
4.  [Repository Structure](#repository-structure)
5.  [Dependencies and Integration](#dependencies-and-integration)
6.  [API Documentation](#api-documentation)
7.  [Development Notes](#development-notes)
8.  [Known Issues and Limitations](#known-issues-and-limitations)
9.  [Additional Documentation](#additional-documentation)

## Architecture

### High-level Architecture Overview

The `ghtools` project employs a **Monolithic Script Architecture** written entirely in Bash. It functions as an **Orchestrator** or **Facade** layer, providing a simplified, unified interface over several powerful external command-line tools.

The core design follows a **Command Dispatcher Pattern**, where the main script parses the user's command and delegates execution to a specific, modular `action_*` function.

### Technology Stack and Frameworks

| Component | Technology | Role |
| :--- | :--- | :--- |
| **Core Logic** | Bash | Scripting language for control flow, configuration, and orchestration. |
| **GitHub Interaction** | GitHub CLI (`gh`) | Primary tool for all GitHub API calls (data retrieval, creation, deletion). |
| **Data Processing** | `jq` | JSON parser used for filtering and transforming data from `gh` output. |
| **User Interface** | `gum` (Optional) / `fzf` (Required) | Provides styled TUI components, interactive prompts, and fuzzy-finding selection. |
| **Local Git Operations** | `git` | Standard tool for cloning, pulling, and checking local repository status. |

### Component Relationships (with mermaid diagrams)

The script's internal structure is highly coupled, with core actions relying on a centralized set of data retrieval and UI utilities.

```mermaid
graph TD
    A[Main Entry Point] --> B(Configuration Loading);
    B --> C(Dependency/Auth Checks);
    C --> D{Command Dispatcher};

    D --> E[Action Functions (clone, sync, delete)];
    D --> F[Interactive Menu];

    E --> G[get_repo_list (Cache Manager)];
    E --> H[run_parallel_jobs (Template Method)];
    E --> I[UI/Printing Functions];

    G --> J[gh CLI];
    G --> K[jq CLI];

    E --> L[git CLI];
    F --> M[fzf/gum CLI];
    I --> M;

    style J fill:#f9f,stroke:#333,stroke-width:2px
    style K fill:#f9f,stroke:#333,stroke-width:2px
    style L fill:#f9f,stroke:#333,stroke-width:2px
    style M fill:#f9f,stroke:#333,stroke-width:2px
```

### Key Design Patterns

| Pattern | Description | Context |
| :--- | :--- | :--- |
| **Facade** | The entire script provides a simplified, unified interface for complex operations involving multiple underlying tools (`gh`, `git`, `fzf`). | The user only interacts with `ghtools`. |
| **Command Dispatcher** | Uses a `case` statement to route the command-line verb to the appropriate `action_*` handler function. | Main execution block of the script. |
| **Template Method** | The `run_parallel_jobs` function provides a reusable structure for concurrent execution, used by `action_clone` and `action_sync`. | Bulk operations. |
| **Cache-Aside** | The `fetch_repositories_json` function manages a file-based cache, checking its validity before calling the external GitHub API. | Repository data retrieval. |

## C4 Model Architecture

### <details>
<summary>Context Diagram</summary>

```mermaid
C4Context
    title Context Diagram for ghtools
    Person(user, "Developer", "Manages GitHub repositories via the command line.")
    System(ghtools, "ghtools CLI", "A Bash script that orchestrates repository management.")
    System_Ext(github, "GitHub API", "Provides repository data and handles all persistence (create, delete, update).")
    System_Ext(cli_tools, "External CLI Tools", "Required tools: gh, git, jq, fzf. Optional: gum.")

    user -- Executes commands --> ghtools
    ghtools -- Calls --> github : Fetches data, performs actions (via gh CLI)
    ghtools -- Orchestrates --> cli_tools : For data processing, TUI, and Git operations
```
</details>

### <details>
<summary>Container Diagram</summary>

```mermaid
C4Container
    title Container Diagram for ghtools
    System_Boundary(ghtools_system, "ghtools Repository Management Tool")
        Container(script, "ghtools Bash Script", "Bash", "The monolithic script containing all application logic, command dispatch, and orchestration.")
        Container(config, "Configuration File", "Plain Text (~/.config/ghtools/config)", "Stores user-defined settings like CACHE_TTL, MAX_JOBS, and DEFAULT_ORG.")
        Container(cache, "Repository Cache", "JSON File (/tmp/ghtools_repos_UID.json)", "Stores a time-to-live (TTL) copy of the GitHub repository list for performance.")
    System_Ext(gh_cli, "GitHub CLI (gh)", "Go Executable", "Primary interface for all GitHub API interactions.")
    System_Ext(jq, "jq CLI", "C Executable", "Used for parsing, filtering, and transforming JSON output.")
    System_Ext(fzf_gum, "TUI Tools (fzf/gum)", "Go Executable", "Provides interactive fuzzy-finding and styled terminal UI.")
    System_Ext(git, "Git CLI", "C Executable", "Used for local repository operations (clone, pull, status checks).")

    script --> config : Reads configuration on startup
    script --> cache : Reads/Writes cached repository data (securely)
    script --> gh_cli : Executes commands for API interaction
    script --> jq : Pipes gh output for JSON processing
    script --> fzf_gum : Pipes data for interactive selection/display
    script --> git : Executes local repository operations (sync, clone)
```
</details>

## Repository Structure

The project is characterized by its simplicity, being primarily a single executable file.

| Directory/File | Purpose |
| :--- | :--- |
| `ghtools` | The single, monolithic Bash script containing all application logic and functions. |
| `install.sh` | Script for checking dependencies and installing the `ghtools` executable. |
| `test/` | Contains the BATS (Bash Automated Testing System) test suite for unit and integration testing. |
| `~/.config/ghtools/` | User configuration directory (created on first run). |

## Dependencies and Integration

The application's core functionality is achieved by integrating with external services and relying on a set of internal utility functions.

### Internal Dependencies

The single-file structure results in **High Coupling** between functions. Core action functions (`action_clone`, `action_sync`, etc.) are tightly dependent on:

*   **Data Retrieval:** `get_repo_list` (which manages caching and calls `gh`/`jq`).
*   **Parallelism:** `run_parallel_jobs` (which manages concurrent execution).
*   **UI/Logging:** `print_error`, `print_success`, `print_info` (which manage styled output based on global flags).

### External Service Dependencies

The tool's primary integration is with the GitHub platform, mediated entirely through the `gh` CLI.

| Service | Integration Point | Protocol/Tool | Details |
| :--- | :--- | :--- | :--- |
| **GitHub** | Repository Management | `gh` CLI | All core actions (list, create, delete, archive) are executed by wrapping the `gh` command. Requires the user to be authenticated via `gh auth login`. |
| **Local Filesystem** | Configuration & Cache | Bash Utilities | Used to read configuration from `$CONFIG_FILE` and manage the temporary, secure repository cache at `$CACHE_FILE`. |

## API Documentation

The `ghtools` application does not expose an HTTP API; its "API" is the set of command-line interface (CLI) commands it provides.

### API Endpoints (CLI Commands)

The tool provides a comprehensive set of commands for repository management:

| Command | Description | Core Functionality |
| :--- | :--- | :--- |
| `list` | Lists repositories, supporting filtering by language or organization. | Data Retrieval, Filtering, Caching |
| `clone` | Interactively selects and clones repositories in parallel. | Interactive Selection, Parallel Execution |
| `sync` | Synchronizes local repositories (performs `git pull --ff-only`). | Local Operations, Parallel Execution |
| `create` | Interactively creates a new GitHub repository. | Remote Write Operation |
| `delete` | Securely and interactively deletes selected repositories. | Remote Write Operation, Authorization Check |
| `status` | Shows the git status (dirty, ahead/behind) of local repositories. | Local Operations, Status Reporting |
| `search` | Interactive fuzzy search across local and remote repositories. | Interactive Selection, Data Filtering |
| `refresh` | Clears the local repository cache, forcing a fresh API call. | Cache Invalidation |

### Request/Response Formats: `ghtools list`

The `list` command is the primary data retrieval endpoint, demonstrating the tool's caching and transformation capabilities.

| Detail | Description |
| :--- | :--- |
| **Method/Path** | `ghtools list [--refresh] [--lang <lang>] [--org <org>]` |
| **Request Parameters** | `--refresh` (Flag): Bypasses the `CACHE_TTL`. `--lang <lang>` (String): Filters by primary language. `--org <org>` (String): Filters by organization. |
| **Authentication** | Inherited from `gh CLI` (requires `gh auth login`). |
| **Response (Success)** | A formatted table printed to standard output, including columns like `NAME`, `DESCRIPTION`, `VISIBILITY`, `LANG`, and `UPDATED` date. |
| **Response (Error)** | An error message printed to `stderr` if the GitHub API call fails (e.g., rate limit exceeded, authentication failure). |
| **Data Flow** | `gh repo list --json` -> `jq` transformation -> Cache File -> `jq` filtering -> Console Output. |

## Development Notes

### Project-specific Conventions

*   **Strict Shell Mode:** The script uses `set -euo pipefail` to ensure robust error handling, immediately exiting on command failure or use of an unset variable.
*   **Function Naming:** Core commands are prefixed with `action_` (e.g., `action_clone`), while utility functions are descriptive (e.g., `print_error`, `check_dependencies`).
*   **Output Abstraction:** All user-facing output must pass through the `print_*` utility functions to respect the global `VERBOSE` and `QUIET` flags and ensure consistent styling.

### Testing Requirements

*   The project uses **BATS (Bash Automated Testing System)** for testing.
*   Tests cover both unit-level logic (e.g., configuration parsing, cache validation) and integration tests that verify the orchestration of external tools (`gh`, `git`).

### Performance Considerations

*   **Caching:** The `CACHE_TTL` (default 600 seconds) minimizes API calls for the repository list, which is the most expensive read operation.
*   **Parallelism:** The `action_clone` and `action_sync` handlers utilize `xargs -P $MAX_JOBS` to execute tasks concurrently, significantly improving performance for bulk operations.
*   **Data Transformation:** Extensive use of `jq` ensures efficient, in-memory JSON processing, avoiding slow shell string manipulation where possible.

## Known Issues and Limitations

| Issue/Limitation | Description | Technical Debt/Risk |
| :--- | :--- | :--- |
| **Fragility to External CLI Changes** | The script is tightly coupled to the specific command-line arguments and JSON output format of `gh` and `jq`. Changes to these external tools could break core functionality. | High Risk |
| **Monolithic Structure** | The single-file architecture and heavy reliance on global variables lead to high coupling, making it difficult to isolate functions for unit testing or refactoring. | Technical Debt |
| **Configuration Security** | The `load_config` function uses `source`, which is inherently risky. While a `grep` check is implemented to validate contents, a sophisticated bypass could lead to arbitrary code execution. | Security Risk |
| **UX Degradation** | The enhanced Terminal UI (TUI) is dependent on the optional `gum` tool. Users without `gum` will experience a functional but less visually appealing interface. | Limitation |
| **No Explicit Retry Logic** | The script relies on the underlying resilience of the `gh` CLI. There is no explicit retry or circuit breaker logic implemented in Bash for handling transient API failures. | Limitation |

## Additional Documentation

*   [Testing Documentation](./test/README.md): Details on the BATS test suite setup and execution.
*   [AI Agent Documentation](./AGENTS.md): Documentation related to the AI agents used in the development workflow.
*   [Claude AI Documentation](./CLAUDE.md): Specific documentation regarding the use of the Claude AI model.
*   [AI Command Definitions](./.claude/commands/): A collection of markdown files defining specific commands/prompts for the AI assistant. (Note: This directory is for development context.)