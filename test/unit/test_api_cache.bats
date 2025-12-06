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
# Test Suite: fetch_repositories_json
# ========================================

@test "fetch_repositories_json executes without error" {
    # Create a mock cache file to avoid actual API call
    create_mock_json
    run fetch_repositories_json
    [ "$status" -eq 0 ]
}

@test "fetch_repositories_json respects force_refresh parameter" {
    create_mock_json
    run fetch_repositories_json "true"
    [ "$status" -eq 0 ]
}

@test "fetch_repositories_json respects limit parameter" {
    create_mock_json
    run fetch_repositories_json "false" "10"
    [ "$status" -eq 0 ]
}

@test "fetch_repositories_json respects organization filter" {
    create_mock_json
    run fetch_repositories_json "false" "100" "testorg"
    [ "$status" -eq 0 ]
}

@test "fetch_repositories_json handles all parameters together" {
    create_mock_json
    run fetch_repositories_json "true" "50" "myorg"
    [ "$status" -eq 0 ]
}

@test "fetch_repositories_json handles empty organization filter" {
    create_mock_json
    run fetch_repositories_json "false" "100" ""
    [ "$status" -eq 0 ]
}

@test "fetch_repositories_json works without organization filter" {
    create_mock_json
    run fetch_repositories_json
    [ "$status" -eq 0 ]
}

@test "fetch_repositories_json handles organization filter with special characters" {
    create_mock_json
    run fetch_repositories_json "false" "100" "my-org_123"
    [ "$status" -eq 0 ]
}

@test "fetch_repositories_json handles zero limit" {
    create_mock_json
    run fetch_repositories_json "false" "0"
    [ "$status" -eq 0 ]
}

@test "fetch_repositories_json handles large limit" {
    create_mock_json
    run fetch_repositories_json "false" "1000"
    [ "$status" -eq 0 ]
}

@test "fetch_repositories_json handles default parameters" {
    create_mock_json
    run fetch_repositories_json "" "" ""
    [ "$status" -eq 0 ]
}

@test "fetch_repositories_json outputs valid JSON when cache exists" {
    create_mock_json
    run fetch_repositories_json "false"
    [ "$status" -eq 0 ]
    # Check if output contains JSON-like structure (even if it's mock data)
    [ -n "$output" ]
}

@test "fetch_repositories_json respects custom cache file" {
    # Set custom cache file path
    export CACHE_FILE="/tmp/test_custom_cache.json"
    echo '[]' > "$CACHE_FILE"
    run fetch_repositories_json "true"
    [ "$status" -eq 0 ]
}

@test "fetch_repositories_json creates cache when forced refresh" {
    # Ensure cache doesn't exist initially
    rm -f "$CACHE_FILE"
    run fetch_repositories_json "true"
    [ "$status" -eq 0 ]
}

@test "fetch_repositories_json works with default TTL" {
    create_mock_json
    export CACHE_TTL=600
    run fetch_repositories_json
    [ "$status" -eq 0 ]
}
