#!/usr/bin/env bats

load '../test_helper.bash'

# Test: fetch_repositories_json with refresh
@test "fetch_repositories_json with force refresh creates cache" {
    PATH="${MOCK_DIR}:${PATH}"
    rm -f "$TEST_CACHE_FILE"
    export CACHE_TTL=600
    export DEFAULT_ORG=""

    run fetch_repositories_json "true"
    [ "$status" -eq 0 ]
    [ -f "$TEST_CACHE_FILE" ]
}

# Test: fetch_repositories_json uses cache when valid
@test "fetch_repositories_json uses cache when valid" {
    PATH="${MOCK_DIR}:${PATH}"
    # Create a valid cache file
    create_mock_json
    touch -t 202412010000 "$TEST_CACHE_FILE"

    run fetch_repositories_json "false"
    [ "$status" -eq 0 ]
}

# Test: fetch_repositories_json with organization filter
@test "fetch_repositories_json respects organization filter" {
    PATH="${MOCK_DIR}:${PATH}"
    export DEFAULT_ORG="testorg"

    run fetch_repositories_json "true" "100" "testorg"
    [ "$status" -eq 0 ]
    # Check that cache was created
    [ -f "$TEST_CACHE_FILE" ]
}

# Test: fetch_repositories_json with invalid org (should still create cache)
@test "fetch_repositories_json handles org filter gracefully" {
    PATH="${MOCK_DIR}:${PATH}"
    rm -f "$TEST_CACHE_FILE"

    run fetch_repositories_json "true" "100" "nonexistent"
    [ "$status" -eq 0 ]
    [ -f "$TEST_CACHE_FILE" ]
}

# Test: Configuration loading
@test "load_config respects XDG_CONFIG_HOME" {
    export XDG_CONFIG_HOME="$TEST_TMP_DIR/custom_config"
    mkdir -p "$XDG_CONFIG_HOME/ghtools"
    # Use a valid config variable
    echo 'MAX_JOBS=8' > "$XDG_CONFIG_HOME/ghtools/config"

    # Reload to pick up new config dir
    CONFIG_DIR="${XDG_CONFIG_HOME}/ghtools"
    CONFIG_FILE="$CONFIG_DIR/config"

    run load_config
    [ "$status" -eq 0 ]
}

# Test: Default configuration values
@test "init_config creates default config with commented options" {
    rm -rf "$TEST_CONFIG_DIR"
    run init_config
    [ "$status" -eq 0 ]
    [ -f "$TEST_CONFIG_FILE" ]

    content=$(cat "$TEST_CONFIG_FILE")
    [[ "$content" == *"CACHE_TTL"* ]]
    [[ "$content" == *"MAX_JOBS"* ]]
    [[ "$content" == *"DEFAULT_ORG"* ]]
    [[ "$content" == *"DEFAULT_CLONE_PATH"* ]]
}

# Test: Configuration overrides
@test "load_config applies configuration overrides" {
    mkdir -p "$TEST_CONFIG_DIR"
    cat > "$TEST_CONFIG_FILE" <<EOF
CACHE_TTL=300
MAX_JOBS=10
DEFAULT_CLONE_PATH="/custom/path"
EOF

    load_config
    [ "$CACHE_TTL" = "300" ]
    [ "$MAX_JOBS" = "10" ]
    [ "$DEFAULT_CLONE_PATH" = "/custom/path" ]
}

# Test: Cache TTL validation
@test "is_cache_valid respects CACHE_TTL setting" {
    export CACHE_TTL=60  # 1 minute
    echo '{}' > "$TEST_CACHE_FILE"
    # File is 2 minutes old, should be invalid
    touch -d "2 minutes ago" "$TEST_CACHE_FILE"

    run is_cache_valid
    [ "$status" -eq 1 ]
}

@test "is_cache_valid accepts fresh cache" {
    export CACHE_TTL=3600  # 1 hour
    echo '{}' > "$TEST_CACHE_FILE"
    # File is fresh (just created), should be valid
    touch "$TEST_CACHE_FILE"

    run is_cache_valid
    [ "$status" -eq 0 ]
}
