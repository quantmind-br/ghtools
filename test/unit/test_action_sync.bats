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
# Test Suite: action_sync
# ========================================

@test "action_sync executes without error" {
    create_test_git_repo "/tmp/test_sync_repo"
    export DEFAULT_CLONE_PATH="/tmp"
    run action_sync --path /tmp/test_sync_repo --max-depth 1
    [ "$status" -eq 0 ]
}

@test "action_sync with --dry-run flag executes without error" {
    create_test_git_repo "/tmp/test_sync_dryrun"
    run action_sync --path /tmp/test_sync_dryrun --dry-run --max-depth 1
    [ "$status" -eq 0 ]
}

@test "action_sync with --all flag executes without error" {
    create_test_git_repo "/tmp/test_sync_all"
    run action_sync --path /tmp/test_sync_all --all --max-depth 1
    [ "$status" -eq 0 ]
}

@test "action_sync with --max-depth flag executes without error" {
    create_test_git_repo "/tmp/test_sync_depth"
    run action_sync --path /tmp/test_sync_depth --max-depth 2
    [ "$status" -eq 0 ]
}

@test "action_sync with custom path executes without error" {
    create_test_git_repo "/tmp/custom_sync_path"
    run action_sync --path /tmp/custom_sync_path --max-depth 1
    [ "$status" -eq 0 ]
}

@test "action_sync handles non-existent path" {
    run action_sync --path "/nonexistent/path"
    [ "$status" -eq 0 ]
}

@test "action_sync handles directory with no git repos" {
    mkdir -p "/tmp/no_git_dir"
    run action_sync --path /tmp/no_git_dir --max-depth 1
    [ "$status" -eq 0 ]
}

@test "action_sync handles multiple options together" {
    create_test_git_repo "/tmp/test_sync_multi"
    run action_sync --path /tmp/test_sync_multi --dry-run --all --max-depth 2
    [ "$status" -eq 0 ]
}

@test "action_sync respects quiet mode" {
    export QUIET=true
    create_test_git_repo "/tmp/test_sync_quiet"
    run action_sync --path /tmp/test_sync_quiet --max-depth 1
    [ "$status" -eq 0 ]
}

@test "action_sync respects verbose mode" {
    export VERBOSE=true
    create_test_git_repo "/tmp/test_sync_verbose"
    run action_sync --path /tmp/test_sync_verbose --max-depth 1
    [ "$status" -eq 0 ]
}

@test "action_sync handles path with spaces" {
    create_test_git_repo "/tmp/sync path test"
    run action_sync --path "/tmp/sync path test" --max-depth 1
    [ "$status" -eq 0 ]
}

@test "action_sync handles path with special characters" {
    create_test_git_repo "/tmp/sync-path_123"
    run action_sync --path "/tmp/sync-path_123" --max-depth 1
    [ "$status" -eq 0 ]
}

@test "action_sync handles zero max-depth" {
    create_test_git_repo "/tmp/sync_zero_depth"
    run action_sync --path /tmp/sync_zero_depth --max-depth 0
    [ "$status" -eq 0 ]
}

@test "action_sync handles large max-depth" {
    create_test_git_repo "/tmp/sync_large_depth"
    run action_sync --path /tmp/sync_large_depth --max-depth 10
    [ "$status" -eq 0 ]
}

@test "action_sync handles empty --path parameter" {
    # Empty path may fail or use default
    run action_sync --path ""
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]
}

@test "action_sync handles default path (no --path flag)" {
    create_test_git_repo "/tmp/test_sync_default"
    export DEFAULT_CLONE_PATH="/tmp"
    run action_sync --max-depth 1
    [ "$status" -eq 0 ]
}

@test "action_sync handles multiple git repos" {
    # Clean up any previous test artifacts
    rm -rf "/tmp/sync_multi_1" "/tmp/sync_multi_2" "/tmp/sync_multi_dir"
    create_test_git_repo "/tmp/sync_multi_1"
    create_test_git_repo "/tmp/sync_multi_2"
    mkdir -p "/tmp/sync_multi_dir"
    mv /tmp/sync_multi_1 /tmp/sync_multi_dir/
    mv /tmp/sync_multi_2 /tmp/sync_multi_dir/
    run action_sync --path /tmp/sync_multi_dir --max-depth 2
    [ "$status" -eq 0 ]
}

@test "action_sync respects max_jobs setting" {
    create_test_git_repo "/tmp/test_sync_jobs"
    export MAX_JOBS=2
    run action_sync --path /tmp/test_sync_jobs --max-depth 1
    [ "$status" -eq 0 ]
}

@test "action_sync handles sync_all with multiple repos" {
    mkdir -p "/tmp/sync_all_test"
    create_test_git_repo "/tmp/sync_all_test/repo1"
    create_test_git_repo "/tmp/sync_all_test/repo2"
    run action_sync --path /tmp/sync_all_test --all --max-depth 2
    [ "$status" -eq 0 ]
}

@test "action_sync handles dry-run with no changes" {
    create_test_git_repo "/tmp/sync_dryrun_nochanges"
    run action_sync --path /tmp/sync_dryrun_nochanges --dry-run --max-depth 1
    [ "$status" -eq 0 ]
}
