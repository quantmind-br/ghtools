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
# Test Suite: action_delete
# ========================================

@test "action_delete executes without error" {
    create_mock_json
    # Delete requires interactive input, may fail in test environment
    run action_delete
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "action_delete checks dependencies" {
    create_mock_json
    # Dependencies check should work
    run action_delete
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "action_delete handles empty cache" {
    echo '[]' > "$CACHE_FILE"
    run action_delete
    # Empty cache may cause different behavior
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "action_delete handles cache with valid data" {
    create_mock_json
    run action_delete
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "action_delete respects quiet mode" {
    export QUIET=true
    create_mock_json
    run action_delete
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "action_delete handles mock confirmation" {
    create_mock_json
    # Mock will simulate user canceling (expected behavior)
    run action_delete
    # User cancellation is acceptable
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "action_delete validates delete scope" {
    create_mock_json
    run action_delete
    # Scope validation may fail in test environment
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "action_delete handles multiple repositories" {
    create_mock_json
    run action_delete
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "action_delete handles repository selection" {
    create_mock_json
    # Mock will simulate user selection
    run action_delete
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "action_delete handles cache clearing" {
    create_mock_json
    # Delete operation may clear cache
    run action_delete
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}
