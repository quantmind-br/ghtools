# Code Structure Analysis
## Architectural Overview
The `ghtools` project is a command-line interface (CLI) tool implemented entirely in **Bash**. It follows a **Monolithic Script Architecture** where all logic, configuration, utility functions, and command handlers are contained within a single executable file, `ghtools`.

The architecture is structured around a **Command Dispatcher Pattern**, where the main script parses command-line arguments and delegates execution to specific functions (e.g., `action_list`, `action_clone`, `action_sync`).

Key architectural principles include:
1.  **Dependency on External Tools**: The tool heavily relies on external CLI utilities (`gh`, `gum`, `fzf`, `jq`, `git`) for core functionality (GitHub API interaction, TUI, fuzzy finding, JSON parsing, Git operations). This makes the script an **Orchestrator** of these tools.
2.  **Configuration and Caching Layer**: It includes a simple configuration loading mechanism (`load_config`, `init_config`) and a file-based caching system (`CACHE_FILE`, `CACHE_TTL`) to improve performance by reducing repeated calls to the GitHub API via `gh`.
3.  **Presentation Layer Abstraction**: It uses a set of `print_*` utility functions and `gum_style` helpers to abstract the terminal output, providing a modern, colored TUI when `gum` is available, and falling back to standard ANSI colors otherwise. This introduces a basic form of **Presentation Abstraction**.
4.  **Modular Actions**: The core capabilities (list, clone, sync, create, delete) are implemented as distinct `action_*` functions, promoting separation of concerns within the single script file.

## Core Components
| Component | File/Function | Purpose & Responsibility |
| :--- | :--- | :--- |
| **Main Entry Point** | `ghtools` (main execution block) | Parses command-line arguments, handles global flags (`-v`, `-q`), and dispatches control to the appropriate `action_*` function or the interactive menu (`show_interactive_menu`). |
| **Configuration Manager** | `load_config`, `init_config` | Manages application settings (e.g., `CACHE_TTL`, `MAX_JOBS`, `DEFAULT_ORG`). Loads settings from `$CONFIG_FILE` and creates a default file if none exists. |
| **Cache Manager** | `get_repos_from_cache`, `fetch_and_cache_repos` | Handles the repository list cache. Checks cache validity based on `CACHE_TTL` and orchestrates fetching fresh data from the GitHub API (`gh`) when the cache is stale or missing. |
| **TUI/Style Utilities** | `print_error`, `print_success`, `gum_style`, `use_gum` | Provides a consistent, styled output interface. Detects the presence of `gum` to enable a modern TUI, falling back to standard color codes otherwise. |
| **Action Dispatcher** | `main` function (implicit) | The central control flow that maps user commands (`list`, `clone`, `sync`, etc.) to their corresponding handler functions. |

## Service Definitions
The "services" in this Bash application are the high-level functions that encapsulate the application's core business logic, primarily by orchestrating external tools.

| Service Function | Responsibility | Dependencies |
| :--- | :--- | :--- |
| `action_list` | Fetches, filters, sorts, and formats the list of GitHub repositories based on user-provided options. Handles various export formats (table, csv, json). | `gh`, `jq`, `fzf`, `get_all_repos`, `format_output` |
| `action_clone` | Allows interactive or bulk cloning of selected repositories. Manages parallel execution of clone operations. | `gh`, `git`, `select_repos_fzf`, `run_parallel_jobs` |
| `action_sync` | Scans a local directory for Git repositories, checks their status against the remote, and performs safe `git pull --ff-only` operations. Manages parallel synchronization. | `git`, `find_local_repos`, `get_repo_status`, `run_parallel_jobs` |
| `action_create` | Guides the user through creating a new GitHub repository using `gh repo create`. | `gh`, `gum` (for prompts) |
| `action_delete` | Provides a secure, interactive way to select and delete one or more GitHub repositories, requiring confirmation. | `gh`, `select_repos_fzf` |
| `show_interactive_menu` | Presents the main TUI menu using `gum` or `fzf` and dispatches to the selected action. | `gum`, `fzf` |

## Interface Contracts
Since the codebase is a single Bash script, formal interfaces are not present. However, the script establishes several implicit contracts:

| Contract Type | Description | Implementation Details |
| :--- | :--- | :--- |
| **External Tool Interface** | The script assumes the presence and specific command-line behavior of its dependencies. | **`gh`**: Must be authenticated and available for API calls. **`jq`**: Must be available for JSON parsing and manipulation. **`gum`/`fzf`**: Must be available for interactive selection and TUI. |
| **Repository Data Contract** | The `get_all_repos` function is the single source of truth for repository data, which is expected to be a JSON array of repository objects, conforming to the structure returned by `gh repo list --json`. | The `jq` filters within `get_all_repos` define the exact fields extracted (e.g., `.name`, `.owner.login`, `.isArchived`, `.stargazerCount`). |
| **Output Contract** | All user-facing output must pass through the `print_*` utility functions to ensure consistent styling and adherence to the `QUIET` and `VERBOSE` flags. | Functions like `print_info`, `print_error`, `print_success` are the mandatory gateway for logging. |
| **Configuration Contract** | The `config` file must only contain assignments for a predefined set of variables (`CACHE_TTL`, `MAX_JOBS`, etc.) to prevent arbitrary code execution during `source`. | Enforced by the security check within `load_config`. |

## Design Patterns Identified
| Pattern | Description | Location/Context |
| :--- | :--- | :--- |
| **Command Dispatcher** | The main execution logic uses a `case` statement on the first argument to route execution to the appropriate action function. | The primary `if/case` block at the end of the `ghtools` script. |
| **Facade** | The `ghtools` script acts as a facade, providing a simplified, unified interface for complex operations that are internally handled by orchestrating multiple underlying tools (`gh`, `git`, `fzf`, `jq`). | The entire `ghtools` script. |
| **Template Method** | The `run_parallel_jobs` function provides a template for executing a list of tasks concurrently, abstracting the parallel execution logic (using `xargs -P`) while allowing different actions (clone, sync) to be plugged in. | Used by `action_clone` and `action_sync`. |
| **Decorator (Implicit)** | The `gum_style` and `print_*` functions decorate raw output with color, formatting, and status icons, enhancing the user experience without changing the core message content. | All output functions. |

## Component Relationships
1.  **Main Script (`ghtools`)** -> **Configuration Manager**: Loads configuration at startup to set global variables.
2.  **Action Functions (`action_*`)** -> **Cache Manager**: Actions that require the repository list (e.g., `list`, `clone`, `delete`) call `get_all_repos` which interacts with the cache.
3.  **Action Functions** -> **External Tools (`gh`, `git`)**: Actions directly invoke external tools to perform their core tasks (e.g., `action_clone` calls `git clone`).
4.  **Action Functions** -> **TUI/Style Utilities**: All actions use `print_*` functions for logging and status updates.
5.  **Interactive Menu (`show_interactive_menu`)** -> **Action Functions**: The menu selects and executes the chosen action function.
6.  **Parallel Execution (`run_parallel_jobs`)** -> **`action_clone` / `action_sync`**: These actions provide the list of tasks to be executed in parallel by the utility function.

## Key Methods & Functions
| Function Name | Purpose | Application Capability |
| :--- | :--- | :--- |
| `get_all_repos` | Fetches the list of all repositories from GitHub (via `gh`), caches the result, and returns the raw JSON data. | Core data retrieval and caching. |
| `select_repos_fzf` | Presents a list of repositories to the user using `fzf` (or `gum filter`) for interactive, multi-select filtering. | Interactive selection for `clone`, `sync`, and `delete` actions. |
| `run_parallel_jobs` | Executes a list of commands in parallel, respecting the `MAX_JOBS` limit. | Performance optimization for bulk operations like `clone` and `sync`. |
| `get_repo_status` | For a local Git repository, checks its status against the remote (ahead, behind, dirty). | Core logic for the `sync` action's status reporting. |
| `format_output` | Takes the filtered repository data and formats it into the requested output format (`table`, `csv`, or `json`). | Data presentation and export capability. |
| `show_interactive_menu` | The primary user interface for the tool, guiding the user through available commands. | User experience and command discovery. |

## Available Documentation
| Document Path | Purpose | Quality Evaluation |
| :--- | :--- | :--- |
| `/README.md` | High-level overview, functionalities, installation instructions, dependencies, and detailed usage examples for all main commands (`list`, `sync`, `clone`, etc.). | **Excellent**. Provides a comprehensive guide for end-users, covering both interactive and direct command usage, and detailing options for complex commands like `list` and `sync`. It clearly defines the tool's scope and requirements. |
| `/AGENTS.md` | Likely documentation for AI agents used in the development workflow. | **N/A (Development/Meta)**. Not directly related to the `ghtools` application structure but provides context on the development environment's automation tools. |
| `/CLAUDE.md` | Likely documentation related to the use of the Claude AI model in the development process. | **N/A (Development/Meta)**. Similar to `AGENTS.md`, it describes development tooling rather than the application itself. |
| `/.claude/commands/...` | A large collection of markdown files defining specific commands/prompts for an AI assistant (Claude). | **N/A (Development/Meta)**. These files define the *behavior* of the AI assistant for tasks like PR creation, conflict resolution, and code review, indicating a highly automated development workflow. |
| `/test/README.md` | Documentation for the BATS (Bash Automated Testing System) test suite. | **Good**. Explains the testing setup, which is crucial for a robust Bash application. The presence of unit and integration tests indicates a focus on code quality. |