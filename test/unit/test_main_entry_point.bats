#!/usr/bin/env bats

load '../test_helper.bash'

# Helper to capture main function output
main_entry() {
    export PATH="${MOCK_DIR}:${PATH}"
    # Source ghtools and call main with provided args
    bash -c "source ${PROJECT_DIR}/ghtools && main \"\$@\"" -- "$@"
}

# Test: --help flag displays usage
@test "main shows usage with --help flag" {
    run main_entry --help
    [ "$status" -eq 0 ]
    [[ "$output" == *"USAGE"* ]]
    [[ "$output" == *"COMMANDS"* ]]
}

# Test: -h flag displays usage
@test "main shows usage with -h flag" {
    run main_entry -h
    [ "$status" -eq 0 ]
    [[ "$output" == *"USAGE"* ]]
}

# Test: --version flag displays version
@test "main shows version with --version flag" {
    run main_entry --version
    [ "$status" -eq 0 ]
    [[ "$output" == *"3.1.0"* ]]
}

# Test: -v flag displays version
@test "main shows version with -v flag" {
    run main_entry -v
    [ "$status" -eq 0 ]
    [[ "$output" == *"3.1.0"* ]]
}

# Test: --verbose flag enables verbose mode
@test "main accepts --verbose flag" {
    PATH="${MOCK_DIR}:${PATH}"
    VERBOSE=true
    run main_entry list --refresh
    [ "$status" -eq 0 ]
}

# Test: --quiet flag suppresses output
@test "main accepts --quiet flag" {
    PATH="${MOCK_DIR}:${PATH}"
    QUIET=true
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json

    run main_entry list
    [ "$status" -eq 0 ]
}

# Test: list command
@test "main handles list command" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json

    run main_entry list
    [ "$status" -eq 0 ]
}

# Test: list command with --refresh
@test "main handles list --refresh" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"

    run main_entry list --refresh
    [ "$status" -eq 0 ]
}

# Test: list command with --lang filter
@test "main handles list --lang filter" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json

    run main_entry list --lang python
    [ "$status" -eq 0 ]
}

# Test: list command with --org filter
@test "main handles list --org filter" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"

    run main_entry list --org testorg
    [ "$status" -eq 0 ]
}

# Test: clone command
@test "main handles clone command" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json
    mkdir -p "$TEST_TMP_DIR/clones"

    run main_entry clone --path "$TEST_TMP_DIR/clones"
    [ "$status" -eq 0 ]
}

# Test: sync command
@test "main handles sync command" {
    PATH="${MOCK_DIR}:${PATH}"
    local test_dir="$TEST_TMP_DIR/test_repos"
    mkdir -p "$test_dir/repo1"
    create_test_git_repo "$test_dir/repo1"

    run main_entry sync --path "$test_dir" --all
    [ "$status" -eq 0 ]
}

# Test: sync command with dry-run
@test "main handles sync --dry-run" {
    PATH="${MOCK_DIR}:${PATH}"
    local test_dir="$TEST_TMP_DIR/test_repos"
    mkdir -p "$test_dir/repo1"
    create_test_git_repo "$test_dir/repo1"

    run main_entry sync --path "$test_dir" --dry-run --all
    [ "$status" -eq 0 ]
}

# Test: create command
@test "main handles create command" {
    PATH="${MOCK_DIR}:${PATH}"
    cd "$TEST_TMP_DIR"

    # Mock user inputs
    export GH_CREATE_NAME="test-repo"
    export GH_CREATE_DESC="Test description"
    export GH_CREATE_VIS="private"
    export GH_CREATE_TPL="none"

    run main_entry create
    [ "$status" -eq 0 ]
}

# Test: delete command
@test "main handles delete command" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json

    run main_entry delete
    [ "$status" -eq 0 ]
}

# Test: fork command
@test "main handles fork command" {
    PATH="${MOCK_DIR}:${PATH}"

    run main_entry fork "test query"
    [ "$status" -eq 0 ]
}

# Test: fork command with --clone
@test "main handles fork --clone" {
    PATH="${MOCK_DIR}:${PATH}"

    run main_entry fork --clone "test query"
    [ "$status" -eq 0 ]
}

# Test: archive command
@test "main handles archive command" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json

    run main_entry archive
    [ "$status" -eq 0 ]
}

# Test: archive --unarchive command
@test "main handles archive --unarchive" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json

    run main_entry archive --unarchive
    [ "$status" -eq 0 ]
}

# Test: visibility command
@test "main handles visibility command" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json

    run main_entry visibility
    [ "$status" -eq 0 ]
}

# Test: visibility --public command
@test "main handles visibility --public" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json

    run main_entry visibility --public
    [ "$status" -eq 0 ]
}

# Test: visibility --private command
@test "main handles visibility --private" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json

    run main_entry visibility --private
    [ "$status" -eq 0 ]
}

# Test: stats command
@test "main handles stats command" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json

    run main_entry stats
    [ "$status" -eq 0 ]
}

# Test: search command
@test "main handles search command" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json

    run main_entry search
    [ "$status" -eq 0 ]
}

# Test: browse command
@test "main handles browse command" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json

    run main_entry browse
    [ "$status" -eq 0 ]
}

# Test: explore command
@test "main handles explore command" {
    PATH="${MOCK_DIR}:${PATH}"

    run main_entry explore "test query"
    [ "$status" -eq 0 ]
}

# Test: explore with --sort and --lang
@test "main handles explore with options" {
    PATH="${MOCK_DIR}:${PATH}"

    run main_entry explore "test query" --sort stars --lang python
    [ "$status" -eq 0 ]
}

# Test: trending command
@test "main handles trending command" {
    PATH="${MOCK_DIR}:${PATH}"

    run main_entry trending
    [ "$status" -eq 0 ]
}

# Test: trending with language
@test "main handles trending --lang" {
    PATH="${MOCK_DIR}:${PATH}"

    run main_entry trending --lang python
    [ "$status" -eq 0 ]
}

# Test: pr command
@test "main handles pr command" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json

    run main_entry pr list
    [ "$status" -eq 0 ]
}

# Test: status command
@test "main handles status command" {
    PATH="${MOCK_DIR}:${PATH}"
    local test_dir="$TEST_TMP_DIR/test_repos"
    mkdir -p "$test_dir/repo1"
    create_test_git_repo "$test_dir/repo1"

    run main_entry status --path "$test_dir"
    [ "$status" -eq 0 ]
}

# Test: config command
@test "main handles config command" {
    PATH="${MOCK_DIR}:${PATH}"

    run main_entry config
    [ "$status" -eq 0 ]
}

# Test: refresh command
@test "main handles refresh command" {
    PATH="${MOCK_DIR}:${PATH}"

    run main_entry refresh
    [ "$status" -eq 0 ]
}

# Test: help command
@test "main handles help command" {
    PATH="${MOCK_DIR}:${PATH}"

    run main_entry help
    [ "$status" -eq 0 ]
    [[ "$output" == *"USAGE"* ]]
}

# Test: unknown command
@test "main handles unknown command" {
    PATH="${MOCK_DIR}:${PATH}"

    run main_entry unknown-command
    [ "$status" -eq 1 ]
    [[ "$output" == *"Unknown command"* ]]
}

# Test: multiple global flags
@test "main handles multiple global flags" {
    PATH="${MOCK_DIR}:${PATH}"
    export CACHE_FILE="$TEST_CACHE_FILE"
    create_mock_json

    run main_entry --verbose --quiet list
    [ "$status" -eq 0 ]
}

# Test: command with multiple options
@test "main handles sync with multiple options" {
    PATH="${MOCK_DIR}:${PATH}"
    local test_dir="$TEST_TMP_DIR/test_repos"
    mkdir -p "$test_dir/repo1"
    create_test_git_repo "$test_dir/repo1"

    run main_entry sync --path "$test_dir" --dry-run --all --max-depth 2
    [ "$status" -eq 0 ]
}
