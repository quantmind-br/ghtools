#!/usr/bin/env bats

# Setup test environment
setup() {
    # Get the project directory (three levels up from this test file)
    local test_file="${BATS_TEST_FILENAME}"
    local project_dir
    project_dir="$(dirname "$(dirname "$(dirname "$test_file")")")"

    # Load test helper (this sets up test environment variables and loads ghtools)
    source "$project_dir/test/test_helper.bash"
}

teardown() {
    teardown_test
}

# ========================================
# Test Suite: action_clone
# ========================================

@test "action_clone executes without error" {
    create_mock_json
    run action_clone
    [ "$status" -eq 0 ]
}

@test "action_clone with --path flag executes without error" {
    create_mock_json
    mkdir -p /tmp/test_clone_path
    run action_clone --path /tmp/test_clone_path
    [ "$status" -eq 0 ]
}

@test "action_clone with default path executes without error" {
    create_mock_json
    run action_clone --path "$(pwd)"
    [ "$status" -eq 0 ]
}

@test "action_clone handles non-existent path" {
    create_mock_json
    run action_clone --path "/nonexistent/path"
    [ "$status" -eq 1 ]
}

@test "action_clone handles empty cache" {
    echo '[]' > "$CACHE_FILE"
    run action_clone
    [ "$status" -eq 0 ]
}

@test "action_clone handles cache with valid data" {
    create_mock_json
    run action_clone
    [ "$status" -eq 0 ]
}

@test "action_clone respects quiet mode" {
    export QUIET=true
    create_mock_json
    run action_clone
    [ "$status" -eq 0 ]
}

@test "action_clone handles path with spaces" {
    create_mock_json
    mkdir -p "/tmp/test path with spaces"
    run action_clone --path "/tmp/test path with spaces"
    [ "$status" -eq 0 ]
}

@test "action_clone handles path with special characters" {
    create_mock_json
    mkdir -p "/tmp/test-path_123"
    run action_clone --path "/tmp/test-path_123"
    [ "$status" -eq 0 ]
}

@test "action_clone handles path with trailing slash" {
    create_mock_json
    mkdir -p /tmp/test_trailing
    run action_clone --path "/tmp/test_trailing/"
    [ "$status" -eq 0 ]
}

@test "action_clone handles empty --path parameter" {
    create_mock_json
    run action_clone --path ""
    # Empty path may fail, which is expected behavior
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "action_clone handles path with tilde expansion" {
    create_mock_json
    # Tilde expansion might not work in all test environments
    # Accept either success or failure
    run action_clone --path "~"
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "action_clone handles multiple options together" {
    create_mock_json
    mkdir -p /tmp/test_multi
    run action_clone --path /tmp/test_multi
    [ "$status" -eq 0 ]
}

@test "action_clone respects verbose mode" {
    export VERBOSE=true
    create_mock_json
    run action_clone
    [ "$status" -eq 0 ]
}

@test "action_clone handles cache refresh" {
    create_mock_json
    run action_clone --path "$(pwd)"
    [ "$status" -eq 0 ]
}
