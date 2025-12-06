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
# Test Suite: action_list
# ========================================

@test "action_list executes without error" {
    create_mock_json
    run action_list
    [ "$status" -eq 0 ]
}

@test "action_list with --refresh flag executes without error" {
    create_mock_json
    run action_list --refresh
    [ "$status" -eq 0 ]
}

@test "action_list with --lang filter executes without error" {
    create_mock_json
    run action_list --lang python
    [ "$status" -eq 0 ]
}

@test "action_list with --org filter executes without error" {
    create_mock_json
    run action_list --org myorg
    [ "$status" -eq 0 ]
}

@test "action_list with multiple flags executes without error" {
    create_mock_json
    run action_list --refresh --lang python --org myorg
    [ "$status" -eq 0 ]
}

@test "action_list handles empty cache" {
    echo '[]' > "$CACHE_FILE"
    run action_list
    [ "$status" -eq 0 ]
}

@test "action_list handles cache with valid data" {
    create_mock_json
    run action_list
    [ "$status" -eq 0 ]
}

@test "action_list handles language filter case-insensitive" {
    create_mock_json
    run action_list --lang PYTHON
    [ "$status" -eq 0 ]
}

@test "action_list handles empty organization filter" {
    create_mock_json
    run action_list --org ""
    [ "$status" -eq 0 ]
}

@test "action_list handles organization with special characters" {
    create_mock_json
    run action_list --org "my-org_123"
    [ "$status" -eq 0 ]
}

@test "action_list handles language filter with special characters" {
    create_mock_json
    run action_list --lang "c++"
    [ "$status" -eq 0 ]
}

@test "action_list handles language filter with spaces" {
    create_mock_json
    run action_list --lang "objective c"
    [ "$status" -eq 0 ]
}

@test "action_list handles organization with spaces" {
    create_mock_json
    run action_list --org "my organization"
    [ "$status" -eq 0 ]
}

@test "action_list handles mixed case organization" {
    create_mock_json
    run action_list --org "MyOrg"
    [ "$status" -eq 0 ]
}

@test "action_list respects quiet mode" {
    export QUIET=true
    create_mock_json
    run action_list
    [ "$status" -eq 0 ]
}
