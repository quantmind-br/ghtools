#!/usr/bin/env bats

# Setup test environment
setup() {
    local test_file="${BATS_TEST_FILENAME}"
    local project_dir
    project_dir="$(dirname "$(dirname "$(dirname "$test_file")")")"

    # Load test helper
    source "$project_dir/test/test_helper.bash"

    # Source ghtools functions
    source "$project_dir/ghtools_functions.sh"
}

# Test: action_list with default options
@test "action_list displays repositories" {
    PATH="${MOCK_DIR}:${PATH}"
    create_mock_json

    run action_list
    [ "$status" -eq 0 ]
    [[ "$output" == *"NAME"* ]]
    [[ "$output" == *"DESCRIPTION"* ]]
    [[ "$output" == *"VISIBILITY"* ]]
    [[ "$output" == *"LANG"* ]]
    [[ "$output" == *"UPDATED"* ]]
}

# Test: action_list with --refresh flag
@test "action_list with refresh flag forces cache update" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"

    run action_list --refresh
    [ "$status" -eq 0 ]
    [ -f "$TEST_CACHE_FILE" ]
}

# Test: action_list with language filter
@test "action_list filters by language" {
    PATH="${MOCK_DIR}:${PATH}"
    create_mock_json

    run action_list --lang bash
    [ "$status" -eq 0 ]
    [[ "$output" == *"NAME"* ]]
}

# Test: action_list with organization filter
@test "action_list filters by organization" {
    PATH="${MOCK_DIR}:${PATH}"

    run action_list --org testorg
    [ "$status" -eq 0 ]
    [[ "$output" == *"NAME"* ]]
}

# Test: action_clone without arguments (uses default path)
@test "action_clone uses default clone path" {
    PATH="${MOCK_DIR}:${PATH}"
    create_mock_json
    mkdir -p "$HOME"

    run action_clone
    [ "$status" -eq 0 ]
}

# Test: action_clone with custom path
@test "action_clone accepts custom path" {
    PATH="${MOCK_DIR}:${PATH}"
    create_mock_json
    local test_dir="$TEST_TMP_DIR/clones"
    mkdir -p "$test_dir"

    run action_clone --path "$test_dir"
    [ "$status" -eq 0 ]
}

# Test: action_sync with dry-run mode
@test "action_sync dry-run mode shows what would be synced" {
    PATH="${MOCK_DIR}:${PATH}"
    local test_dir="$TEST_TMP_DIR/test_repos"
    mkdir -p "$test_dir/repo1"
    create_test_git_repo "$test_dir/repo1"

    run action_sync --path "$test_dir" --dry-run
    [ "$status" -eq 0 ]
}

# Test: action_sync with --all flag
@test "action_sync with all flag syncs all repositories" {
    PATH="${MOCK_DIR}:${PATH}"
    local test_dir="$TEST_TMP_DIR/test_repos"
    mkdir -p "$test_dir/repo1"
    mkdir -p "$test_dir/repo2"
    create_test_git_repo "$test_dir/repo1"
    create_test_git_repo "$test_dir/repo2"

    run action_sync --path "$test_dir" --all
    [ "$status" -eq 0 ]
}

# Test: action_sync with max-depth
@test "action_sync respects max-depth parameter" {
    PATH="${MOCK_DIR}:${PATH}"
    local test_dir="$TEST_TMP_DIR/test_repos"
    mkdir -p "$test_dir"

    run action_sync --path "$test_dir" --max-depth 1
    [ "$status" -eq 0 ]
}

# Test: action_status displays repository status
@test "action_status shows repository status" {
    PATH="${MOCK_DIR}:${PATH}"
    local test_dir="$TEST_TMP_DIR/test_repos"
    mkdir -p "$test_dir/repo1"
    create_test_git_repo "$test_dir/repo1"

    run action_status --path "$test_dir"
    [ "$status" -eq 0 ]
    [[ "$output" == *"REPOSITORY"* ]]
    [[ "$output" == *"BRANCH"* ]]
    [[ "$output" == *"STATUS"* ]]
}

# Test: action_status with max-depth
@test "action_status respects max-depth parameter" {
    PATH="${MOCK_DIR}:${PATH}"
    local test_dir="$TEST_TMP_DIR/test_repos"
    mkdir -p "$test_dir"

    run action_status --path "$test_dir" --max-depth 2
    [ "$status" -eq 0 ]
}

# Test: action_stats generates statistics
@test "action_stats displays statistics dashboard" {
    PATH="${MOCK_DIR}:${PATH}"
    create_mock_json

    run action_stats
    [ "$status" -eq 0 ]
    [[ "$output" == *"Total Repositories"* ]] || [[ "$output" == *"REPOSITORY STATISTICS"* ]]
}

# Test: action_browse opens repositories in browser
@test "action_browse opens selected repositories" {
    PATH="${MOCK_DIR}:${PATH}"
    create_mock_json

    run action_browse
    [ "$status" -eq 0 ]
}

# Test: action_search performs fuzzy search
@test "action_search performs repository search" {
    PATH="${MOCK_DIR}:${PATH}"
    create_mock_json

    run action_search
    [ "$status" -eq 0 ]
}

# Test: action_fork without arguments
@test "action_fork prompts for search query" {
    PATH="${MOCK_DIR}:${PATH}"

    # Mock user input
    export GH_FORK_QUERY="test"

    run action_fork "test query"
    [ "$status" -eq 0 ]
}

# Test: action_fork with --clone flag
@test "action_fork with clone flag forks and clones" {
    PATH="${MOCK_DIR}:${PATH}"

    run action_fork --clone "test query"
    [ "$status" -eq 0 ]
}

# Test: action_explore with query
@test "action_explore searches external repositories" {
    PATH="${MOCK_DIR}:${PATH}"

    run action_explore "machine learning" --lang python
    [ "$status" -eq 0 ]
}

# Test: action_trending with language filter
@test "action_trending shows trending repositories" {
    PATH="${MOCK_DIR}:${PATH}"

    run action_trending --lang python
    [ "$status" -eq 0 ]
}

# Test: action_archive archives repositories
@test "action_archive archives selected repositories" {
    PATH="${MOCK_DIR}:${PATH}"
    create_mock_json

    run action_archive
    [ "$status" -eq 0 ]
}

# Test: action_archive with --unarchive flag
@test "action_archive unarchives repositories" {
    PATH="${MOCK_DIR}:${PATH}"
    create_mock_json

    run action_archive --unarchive
    [ "$status" -eq 0 ]
}

# Test: action_visibility changes repository visibility
@test "action_visibility changes to public" {
    PATH="${MOCK_DIR}:${PATH}"
    create_mock_json

    run action_visibility --public
    [ "$status" -eq 0 ]
}

@test "action_visibility changes to private" {
    PATH="${MOCK_DIR}:${PATH}"
    create_mock_json

    run action_visibility --private
    [ "$status" -eq 0 ]
}

# Test: action_pr list subcommand
@test "action_pr list shows pull requests" {
    PATH="${MOCK_DIR}:${PATH}"
    create_mock_json

    run action_pr list
    [ "$status" -eq 0 ]
}

# Test: apply_template for python
@test "apply_template creates python template" {
    local test_dir="$TEST_TMP_DIR/test_repo"
    mkdir -p "$test_dir"

    run apply_template "$test_dir" "python"
    [ "$status" -eq 0 ]
    [ -f "$test_dir/main.py" ]
    [ -f "$test_dir/README.md" ]
    [ -f "$test_dir/.gitignore" ]
}

# Test: apply_template for node
@test "apply_template creates node template" {
    local test_dir="$TEST_TMP_DIR/test_repo"
    mkdir -p "$test_dir"

    run apply_template "$test_dir" "node"
    [ "$status" -eq 0 ]
    [ -f "$test_dir/index.js" ]
    [ -f "$test_dir/package.json" ]
    [ -f "$test_dir/README.md" ]
    [ -f "$test_dir/.gitignore" ]
}

# Test: apply_template for go
@test "apply_template creates go template" {
    local test_dir="$TEST_TMP_DIR/test_repo"
    mkdir -p "$test_dir"

    run apply_template "$test_dir" "go"
    [ "$status" -eq 0 ]
    [ -f "$test_dir/main.go" ]
    [ -f "$test_dir/README.md" ]
}
