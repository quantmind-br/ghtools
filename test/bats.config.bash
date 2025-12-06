#!/usr/bin/env bash

# Bats configuration file
# This file is automatically loaded by bats

# Set timeout for individual tests
export BATS_TEST_TIMEOUT=60

# Enable parallel execution if supported
if command -v parallel &>/dev/null; then
    export BATS_NO_PARALLELIZE_FOREGROUND=
fi

# Load test helper
setup() {
    # This runs before each test file
    load 'test_helper.bash'
}

# Teardown after each test file
teardown() {
    # Clean up any remaining test artifacts
    :
}
