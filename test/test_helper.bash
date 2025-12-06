#!/usr/bin/env bash

# Test helper for ghtools
# This file sets up the test environment

# Get the directory where this script is located
TEST_DIR="$(cd "$(dirname "${BATS_TEST_FILENAME}")" && pwd)"
PROJECT_DIR="$(cd "${TEST_DIR}/../.." && pwd)"

# Enable test mode to prevent auto-execution and strict mode issues
export GHTOOLS_TEST_MODE=1

# Source the ghtools script directly (it now supports being sourced)
source "${PROJECT_DIR}/ghtools"
export -f load_config 2>/dev/null || true
export -f init_config 2>/dev/null || true
export -f use_gum 2>/dev/null || true
export -f gum_style 2>/dev/null || true
export -f print_error 2>/dev/null || true
export -f print_success 2>/dev/null || true
export -f print_info 2>/dev/null || true
export -f print_warning 2>/dev/null || true
export -f print_verbose 2>/dev/null || true
export -f show_header 2>/dev/null || true
export -f show_divider 2>/dev/null || true
export -f run_with_spinner 2>/dev/null || true
export -f gum_confirm 2>/dev/null || true
export -f gum_input 2>/dev/null || true
export -f gum_choose 2>/dev/null || true
export -f gum_filter 2>/dev/null || true
export -f gum_write 2>/dev/null || true
export -f print_table_row 2>/dev/null || true
export -f show_usage 2>/dev/null || true
export -f check_dependencies 2>/dev/null || true
export -f check_gh_auth 2>/dev/null || true
export -f is_cache_valid 2>/dev/null || true
export -f fetch_repositories_json 2>/dev/null || true
export -f wait_for_jobs 2>/dev/null || true
export -f truncate_text 2>/dev/null || true
export -f action_list 2>/dev/null || true
export -f action_clone 2>/dev/null || true
export -f action_sync 2>/dev/null || true
export -f check_delete_scope 2>/dev/null || true
export -f action_delete 2>/dev/null || true
export -f action_create 2>/dev/null || true
export -f apply_template 2>/dev/null || true
export -f action_fork 2>/dev/null || true
export -f action_explore 2>/dev/null || true
export -f action_trending 2>/dev/null || true
export -f action_archive 2>/dev/null || true
export -f action_stats 2>/dev/null || true
export -f action_search 2>/dev/null || true
export -f action_browse 2>/dev/null || true
export -f action_visibility 2>/dev/null || true
export -f action_pr 2>/dev/null || true
export -f action_pr_list 2>/dev/null || true
export -f action_pr_create 2>/dev/null || true
export -f action_status 2>/dev/null || true
export -f show_menu 2>/dev/null || true

# Set up test environment variables
export TEST_MODE=1
export TEST_TMP_DIR="${BATS_TMPDIR:-/tmp}/ghtools_test_$$"
export TEST_CACHE_FILE="${TEST_TMP_DIR}/test_repos.json"
export TEST_CONFIG_DIR="${TEST_TMP_DIR}/config"
export TEST_CONFIG_FILE="${TEST_CONFIG_DIR}/config"

# Override ghtools global variables to use test paths
export CACHE_FILE="$TEST_CACHE_FILE"
export CONFIG_DIR="$TEST_CONFIG_DIR"
export CONFIG_FILE="$TEST_CONFIG_FILE"
export CACHE_TTL=600
export MAX_JOBS=2
export VERBOSE=false
export QUIET=true

# Create test directory
mkdir -p "$TEST_TMP_DIR"
mkdir -p "$TEST_CONFIG_DIR"

# Mock directory for simulating external commands
export MOCK_DIR="${TEST_DIR}/mocks"
mkdir -p "$MOCK_DIR"

# Make mocks executable (if any exist)
if [ -d "$MOCK_DIR" ] && [ "$(ls -A "$MOCK_DIR" 2>/dev/null)" ]; then
    chmod +x "$MOCK_DIR"/*
fi

# Function to create a temporary git repository for testing
create_test_git_repo() {
    local repo_dir="$1"
    mkdir -p "$repo_dir"
    cd "$repo_dir"
    git init -q
    echo "test" > test.txt
    git add test.txt
    git commit -q -m "Initial commit"
}

# Function to create mock JSON response
create_mock_json() {
    cat > "$TEST_CACHE_FILE" <<'EOF'
[
  {
    "name": "test-repo",
    "nameWithOwner": "user/test-repo",
    "description": "A test repository",
    "visibility": "PUBLIC",
    "primaryLanguage": {
      "name": "bash"
    },
    "stargazerCount": 10,
    "forkCount": 2,
    "diskUsage": 100,
    "updatedAt": "2024-01-01T00:00:00Z",
    "createdAt": "2024-01-01T00:00:00Z",
    "isArchived": false,
    "url": "https://github.com/user/test-repo",
    "sshUrl": "git@github.com:user/test-repo.git"
  },
  {
    "name": "private-repo",
    "nameWithOwner": "user/private-repo",
    "description": "A private test repository",
    "visibility": "PRIVATE",
    "primaryLanguage": {
      "name": "python"
    },
    "stargazerCount": 5,
    "forkCount": 1,
    "diskUsage": 200,
    "updatedAt": "2024-01-02T00:00:00Z",
    "createdAt": "2024-01-01T00:00:00Z",
    "isArchived": false,
    "url": "https://github.com/user/private-repo",
    "sshUrl": "git@github.com:user/private-repo.git"
  }
]
EOF
}

# Function to setup mock gh command
setup_mock_gh() {
    cat > "${MOCK_DIR}/gh" <<'MOCK_SCRIPT'
#!/bin/bash
case "$1" in
    "auth")
        if [ "$2" = "status" ]; then
            echo "github.com
  Logged in to github.com as user (oauth_token)
  Git operations done over the http on api.github.com
  Active user: user
  Token scopes: delete_repo, repo, read:org, gist"
            exit 0
        fi
        if [ "$2" = "refresh" ]; then
            # Mock auth refresh
            echo "Auth refreshed"
            exit 0
        fi
        ;;
    "repo")
        case "$2" in
            "list")
                # Check if --json flag is present
                if [[ "$*" == *"--json"* ]]; then
                    # Output mock JSON data
                    cat <<'JSON_DATA'
[
  {
    "name": "test-repo",
    "nameWithOwner": "user/test-repo",
    "description": "A test repository",
    "visibility": "PUBLIC",
    "primaryLanguage": {"name": "bash"},
    "stargazerCount": 10,
    "forkCount": 2,
    "diskUsage": 100,
    "updatedAt": "2024-01-01T00:00:00Z",
    "createdAt": "2024-01-01T00:00:00Z",
    "isArchived": false,
    "url": "https://github.com/user/test-repo",
    "sshUrl": "git@github.com:user/test-repo.git"
  },
  {
    "name": "private-repo",
    "nameWithOwner": "user/private-repo",
    "description": "A private test repository",
    "visibility": "PRIVATE",
    "primaryLanguage": {"name": "python"},
    "stargazerCount": 5,
    "forkCount": 1,
    "diskUsage": 200,
    "updatedAt": "2024-01-02T00:00:00Z",
    "createdAt": "2024-01-01T00:00:00Z",
    "isArchived": false,
    "url": "https://github.com/user/private-repo",
    "sshUrl": "git@github.com:user/private-repo.git"
  }
]
JSON_DATA
                    exit 0
                elif [ "$3" = "--limit" ]; then
                    # Output mock repository data (non-JSON)
                    echo "user/test-repo"
                    echo "user/private-repo"
                fi
                ;;
            "clone")
                echo "Mock clone of $3"
                exit 0
                ;;
            "create")
                echo "Mock create of $3"
                exit 0
                ;;
            "delete")
                echo "Mock delete of $3"
                exit 0
                ;;
            "edit")
                echo "Mock edit of $3"
                exit 0
                ;;
        esac
        ;;
    "search")
        if [ "$1" = "search" ] && [ "$2" = "repos" ]; then
            # Output mock search results
            echo "example/repo1"
            echo "example/repo2"
        fi
        ;;
    "api")
        # Mock API calls
        if [[ "$*" == *"user/starred"* ]]; then
            echo "Mock starred"
            exit 0
        fi
        ;;
esac
exit 0
MOCK_SCRIPT
    chmod +x "${MOCK_DIR}/gh"
}

# Function to setup mock jq command
setup_mock_jq() {
    cat > "${MOCK_DIR}/jq" <<'MOCK_SCRIPT'
#!/bin/bash
# Simple jq mock - just pass through for basic operations
cat
MOCK_SCRIPT
    chmod +x "${MOCK_DIR}/jq"
}

# Function to setup mock git command
setup_mock_git() {
    cat > "${MOCK_DIR}/git" <<'MOCK_SCRIPT'
#!/bin/bash
case "$1" in
    "init")
        # Handle git init with or without -q flag
        # git init -q (init in current dir, quiet)
        # git init <dir> (init in specified dir)
        # git init -q <dir> (init in specified dir, quiet)
        local dir=""
        shift  # skip "init"
        while [[ $# -gt 0 ]]; do
            case "$1" in
                -q|--quiet) shift ;;  # ignore quiet flag
                *) dir="$1"; shift ;;
            esac
        done
        # If no dir specified, use current directory
        if [[ -z "$dir" ]]; then
            mkdir -p ".git"
            echo "Initialized empty Git repository in .git/"
        else
            mkdir -p "$dir/.git"
            echo "Initialized empty Git repository in $dir/.git/"
        fi
        exit 0
        ;;
    "add")
        exit 0
        ;;
    "commit")
        echo "[master (root-commit) 1234567] Test commit"
        exit 0
        ;;
    "branch")
        if [ "$1" = "branch" ] && [ "$2" = "--show-current" ]; then
            echo "main"
            exit 0
        fi
        if [ "$1" = "branch" ]; then
            echo "* main"
            exit 0
        fi
        ;;
    "diff-index")
        # Pretend working tree is clean
        exit 0
        ;;
    "fetch")
        exit 0
        ;;
    "pull")
        echo "Already up to date."
        exit 0
        ;;
    "push")
        echo "Everything up-to-date"
        exit 0
        ;;
    "ls-remote")
        # Pretend remote branch exists
        exit 0
        ;;
    "rev-list")
        echo "0"
        exit 0
        ;;
esac
# Default: succeed
exit 0
MOCK_SCRIPT
    chmod +x "${MOCK_DIR}/git"
}

# Function to setup mock fzf command
setup_mock_fzf() {
    cat > "${MOCK_DIR}/fzf" <<'MOCK_SCRIPT'
#!/bin/bash
# Mock fzf - just select first line or echo input
if [ "$*" == *"--multi"* ]; then
    # Multi-select mode - return all lines
    cat
else
    # Single select - return first line
    head -n 1
fi
exit 0
MOCK_SCRIPT
    chmod +x "${MOCK_DIR}/fzf"
}

# Function to setup mock gum command
setup_mock_gum() {
    cat > "${MOCK_DIR}/gum" <<'MOCK_SCRIPT'
#!/bin/bash
# Mock gum - just echo the text or return first line for choose
case "$1" in
    "style")
        # Just output the text
        shift
        echo "$@"
        ;;
    "choose")
        # Return first option
        head -n 1
        ;;
    "input")
        # Return default or echo input
        if [[ "$*" == *"--value"* ]]; then
            # Get default value
            echo "default"
        else
            read -r input
            echo "${input:-default}"
        fi
        ;;
    "confirm")
        # Always return true
        exit 0
        ;;
    "filter")
        # Return input
        head -n 1
        ;;
    "write")
        # Return default text
        echo "Test commit message"
        ;;
    "spin")
        # Execute command and return
        shift
        "$@"
        ;;
esac
exit 0
MOCK_SCRIPT
    chmod +x "${MOCK_DIR}/gum"
}

# Teardown function to clean up after each test
teardown_test() {
    # Clean up test directory
    rm -rf "$TEST_TMP_DIR"
}

# Add MOCK_DIR to PATH for tests
export PATH="${MOCK_DIR}:${PATH}"

# Setup mocks before each test
setup() {
    # Create fresh test environment
    mkdir -p "$TEST_TMP_DIR"
    setup_mock_gh
    setup_mock_jq
    setup_mock_git
    setup_mock_fzf
    setup_mock_gum
}

# Teardown after each test
teardown() {
    teardown_test
}
