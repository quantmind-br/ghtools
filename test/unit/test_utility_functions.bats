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

# Test: use_gum function
# Note: use_gum checks for both 'gum' command AND terminal (-t 1)
# In test environment without terminal, this will fail even with mock
@test "use_gum returns false in non-terminal environment" {
    # In tests, we're not in a terminal, so use_gum should return 1
    run use_gum
    # This is expected to fail in test environment (no terminal)
    [ "$status" -eq 1 ]
}

# Test: truncate_text function
@test "truncate_text returns original text when within limit" {
    result=$(truncate_text "short text" 20)
    [ "$result" = "short text" ]
}

@test "truncate_text truncates text when exceeding limit" {
    result=$(truncate_text "this is a very long text that should be truncated" 20)
    [ ${#result} -le 20 ]
    [[ "$result" == *"..."* ]]
}

@test "truncate_text handles empty input" {
    result=$(truncate_text "" 10)
    [ "$result" = "" ]
}

# Test: truncate_text edge case exactly at limit
@test "truncate_text handles text exactly at limit" {
    result=$(truncate_text "exactly ten" 11)
    [ "$result" = "exactly ten" ]
}

# Test: print_table_row function
@test "print_table_row formats output correctly" {
    run print_table_row "Column1" "Column2"
    [ "$status" -eq 0 ]
    [[ "$output" == *"Column1"* ]]
    [[ "$output" == *"Column2"* ]]
}

# Test: wait_for_jobs function
@test "wait_for_jobs runs without error" {
    run wait_for_jobs
    [ "$status" -eq 0 ]
}

# Test: is_cache_valid function
@test "is_cache_valid returns false when cache file doesn't exist" {
    rm -f "$CACHE_FILE"
    run is_cache_valid
    [ "$status" -eq 1 ]
}

@test "is_cache_valid returns false when cache is expired" {
    # Create cache file with old timestamp (year 2020)
    echo '{}' > "$CACHE_FILE"
    touch -t 202001010000 "$CACHE_FILE"
    run is_cache_valid
    [ "$status" -eq 1 ]
}

@test "is_cache_valid returns true when cache is valid" {
    # Create cache file with current timestamp (just touched = fresh)
    echo '{}' > "$CACHE_FILE"
    # touch without -t uses current time, making it fresh
    touch "$CACHE_FILE"
    run is_cache_valid
    [ "$status" -eq 0 ]
}

# Test: check_dependencies function (mock dependencies)
@test "check_dependencies passes with mocked commands" {
    PATH="${MOCK_DIR}:${PATH}"
    run check_dependencies
    [ "$status" -eq 0 ]
}

# Test: check_gh_auth function (mocked)
@test "check_gh_auth passes with mocked gh auth" {
    PATH="${MOCK_DIR}:${PATH}"
    run check_gh_auth
    [ "$status" -eq 0 ]
}

# Test: load_config function
@test "load_config loads valid config without error" {
    mkdir -p "$CONFIG_DIR"
    # Use a valid config variable
    echo 'CACHE_TTL=300' > "$CONFIG_FILE"
    chmod 600 "$CONFIG_FILE"
    run load_config
    [ "$status" -eq 0 ]
}

@test "load_config handles missing config file" {
    rm -f "$CONFIG_FILE"
    run load_config
    [ "$status" -eq 0 ]
}

# Test: init_config function
@test "init_config creates config directory and file" {
    # Remove the config dir that test_helper creates
    rm -rf "$CONFIG_DIR"
    run init_config
    [ "$status" -eq 0 ]
    [ -d "$CONFIG_DIR" ]
    [ -f "$CONFIG_FILE" ]
}

@test "init_config doesn't overwrite existing config" {
    mkdir -p "$CONFIG_DIR"
    echo 'CACHE_TTL=999' > "$CONFIG_FILE"
    original_content=$(cat "$CONFIG_FILE")
    run init_config
    [ "$status" -eq 0 ]
    [ "$(cat "$CONFIG_FILE")" = "$original_content" ]
}

# Test: show_usage function
@test "show_usage outputs usage information" {
    run show_usage
    [ "$status" -eq 0 ]
    [[ "$output" == *"USAGE"* ]]
    [[ "$output" == *"COMMANDS"* ]]
    [[ "$output" == *"OPTIONS"* ]]
}
