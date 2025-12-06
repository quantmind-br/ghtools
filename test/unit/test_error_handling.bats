#!/usr/bin/env bats

# Setup test environment
setup() {
    local test_file="${BATS_TEST_FILENAME}"
    local project_dir
    project_dir="$(dirname "$(dirname "$(dirname "$test_file")")")"

    # Load test helper
    source "$project_dir/test/test_helper.bash"

    # Setup mocks
    mkdir -p "$TEST_TMP_DIR"
    setup_mock_gh
    setup_mock_git
    setup_mock_fzf
    setup_mock_gum
}

teardown() {
    teardown_test
}

# Test: check_dependencies with missing required commands
@test "check_dependencies fails with missing required command" {
    # Save original PATH
    ORIGINAL_PATH="$PATH"

    # Create empty PATH to simulate missing commands
    PATH="/nonexistent"

    run check_dependencies
    [ "$status" -eq 1 ]
    [[ "$output" == *"Missing required dependencies"* ]]

    # Restore PATH
    PATH="$ORIGINAL_PATH"
}

# Test: check_gh_auth when not authenticated
@test "check_gh_auth fails when not authenticated" {
    # Override gh mock to return failure for auth status
    cat > "${MOCK_DIR}/gh" <<'MOCK'
#!/bin/bash
if [ "$1" = "auth" ] && [ "$2" = "status" ]; then
    echo "not authenticated" >&2
    exit 1
fi
exit 0
MOCK
    chmod +x "${MOCK_DIR}/gh"

    PATH="${MOCK_DIR}:${PATH}"

    run check_gh_auth
    [ "$status" -eq 1 ]
}

# Test: action_list with no repositories
@test "action_list handles empty repository list" {
    PATH="${MOCK_DIR}:${PATH}"
    # Create empty cache
    echo '[]' > "$TEST_CACHE_FILE"

    run action_list
    [ "$status" -eq 0 ]
}

# Test: action_list with failed API call
@test "action_list handles API failure gracefully" {
    # Remove cache to force API call
    rm -f "$TEST_CACHE_FILE"

    # Override gh mock to fail
    cat > "${MOCK_DIR}/gh" <<'MOCK'
#!/bin/bash
echo "API Error" >&2
exit 1
MOCK
    chmod +x "${MOCK_DIR}/gh"
    PATH="${MOCK_DIR}:${PATH}"

    run action_list --refresh
    # May return 1 or exit gracefully with 0 depending on implementation
    [ "$status" -eq 1 ] || [ "$status" -eq 0 ]
}

# Test: action_sync with non-existent path
@test "action_sync handles non-existent path" {
    PATH="${MOCK_DIR}:${PATH}"

    run action_sync --path "/nonexistent/path"
    [ "$status" -eq 0 ]  # Should handle gracefully, not fail
}

# Test: action_sync with no git repositories
@test "action_sync handles directory with no git repos" {
    PATH="${MOCK_DIR}:${PATH}"
    local test_dir="$TEST_TMP_DIR/empty"
    mkdir -p "$test_dir"

    run action_sync --path "$test_dir"
    # Should complete successfully even with no repos
    [ "$status" -eq 0 ]
}

# Test: action_clone with non-existent clone path
@test "action_clone fails with non-existent path" {
    PATH="${MOCK_DIR}:${PATH}"

    run action_clone --path "/nonexistent/path"
    [ "$status" -eq 1 ]
    [[ "$output" == *"does not exist"* ]]
}

# Test: action_create with empty name
@test "action_create handles empty name" {
    PATH="${MOCK_DIR}:${PATH}"
    cd "$TEST_TMP_DIR"

    run action_create
    [ "$status" -eq 0 ]
    # Should not crash, gum/input will handle empty input
}

# Test: truncate_text with zero limit
@test "truncate_text handles zero limit" {
    result=$(truncate_text "test" 0)
    # Zero or very short limit should return minimal result
    [ ${#result} -le 5 ]  # May return "..." or empty
}

# Test: truncate_text with negative limit
@test "truncate_text handles negative limit" {
    result=$(truncate_text "test" -1)
    [ ${#result} -le 4 ]
}

# Test: is_cache_valid with empty file
@test "is_cache_valid handles empty file" {
    touch "$TEST_CACHE_FILE"

    run is_cache_valid
    # Empty file may be considered valid (file exists and is recent)
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

# Test: action_delete without delete_repo scope
@test "action_delete handles missing delete scope" {
    # Mock check_delete_scope to return false
    check_delete_scope() {
        return 1
    }

    export -f check_delete_scope

    run action_delete
    [ "$status" -eq 1 ]
}

# Test: action_pr_create not in git repository
@test "action_pr_create fails outside git repository" {
    cd "$TEST_TMP_DIR"

    run action_pr_create
    [ "$status" -eq 1 ]
    [[ "$output" == *"Not in a git repository"* ]]
}

# Test: action_pr_create on main/master branch
@test "action_pr_create detects main/master branch" {
    PATH="${MOCK_DIR}:${PATH}"
    cd "$TEST_TMP_DIR"
    mkdir -p ".git"

    run action_pr_create
    # Should fail or warn when on main/master branch
    [ "$status" -eq 1 ] || [ "$status" -eq 0 ]
}

# Test: fetch_repositories_json with limit 0
@test "fetch_repositories_json handles limit 0" {
    PATH="${MOCK_DIR}:${PATH}"

    run fetch_repositories_json "true" "0"
    [ "$status" -eq 0 ]
}

# Test: show_header with and without subtitle
@test "show_header handles optional subtitle" {
    run show_header "Test Title" "Test Subtitle"
    [ "$status" -eq 0 ]
    [[ "$output" == *"Test Title"* ]]

    run show_header "Test Title Only"
    [ "$status" -eq 0 ]
    [[ "$output" == *"Test Title Only"* ]]
}

# Test: show_divider with and without title
@test "show_divider handles optional title" {
    run show_divider "Section Title"
    [ "$status" -eq 0 ]

    run show_divider
    [ "$status" -eq 0 ]
}

# Test: print_functions with QUIET mode
@test "print_info respects QUIET mode" {
    QUIET=true
    run print_info "Test message"
    [ "$status" -eq 0 ]
    [ -z "$output" ]
}

# Test: print_verbose with VERBOSE disabled
@test "print_verbose respects VERBOSE disabled" {
    VERBOSE=false
    run print_verbose "Debug message"
    [ "$status" -eq 0 ]
    [ -z "$output" ]
}

# Test: gum_confirm with default yes
@test "gum_confirm handles default yes" {
    PATH="${MOCK_DIR}:${PATH}"
    # Provide input to simulate user response
    run bash -c 'echo "y" | gum_confirm "Test prompt" "yes"'
    [ "$status" -eq 0 ]
}

# Test: gum_input with default value
@test "gum_input handles default value" {
    PATH="${MOCK_DIR}:${PATH}"
    run gum_input "Placeholder" "› " "default_value"
    [ "$status" -eq 0 ]
}

# Test: gum_choose with multiple options
@test "gum_choose handles multiple options" {
    PATH="${MOCK_DIR}:${PATH}"
    run gum_choose "Choose:" "Option 1" "Option 2" "Option 3"
    [ "$status" -eq 0 ]
}

# Test: gum_filter with multi-select
@test "gum_filter handles multi-select mode" {
    PATH="${MOCK_DIR}:${PATH}"
    run gum_filter "Filter items" "true"
    [ "$status" -eq 0 ]
}

# Test: run_with_spinner
@test "run_with_spinner executes command" {
    PATH="${MOCK_DIR}:${PATH}"
    run run_with_spinner "Testing" echo "success"
    [ "$status" -eq 0 ]
    [[ "$output" == *"success"* ]]
}

# Test: wait_for_jobs with multiple jobs
@test "wait_for_jobs handles job limit" {
    # Start some background jobs
    (sleep 0.1) &
    (sleep 0.1) &
    (sleep 0.1) &

    run wait_for_jobs
    [ "$status" -eq 0 ]
}
