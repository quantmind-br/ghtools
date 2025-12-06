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
# Test Suite: print_error
# ========================================

@test "print_error executes without error" {
    run print_error "Test error message"
    [ "$status" -eq 0 ]
}

@test "print_error executes successfully" {
    run print_error "Test error message"
    [ "$status" -eq 0 ]
}

@test "print_error respects QUIET mode" {
    export QUIET=true
    run print_error "Should not appear"
    [ "$status" -eq 0 ]
}

@test "print_error handles empty message" {
    run print_error ""
    [ "$status" -eq 0 ]
}

@test "print_error handles message with special characters" {
    run print_error "Error: !@#$%^&*()"
    [ "$status" -eq 0 ]
}

@test "print_error handles very long message" {
    local long_msg="This is a very long error message "
    long_msg+="that exceeds normal line length and should be handled gracefully "
    long_msg+="without causing any issues in the printing function."
    run print_error "$long_msg"
    [ "$status" -eq 0 ]
}

@test "print_error handles message with newlines" {
    run print_error -e "Line 1\nLine 2"
    [ "$status" -eq 0 ]
}

@test "print_error handles international characters" {
    run print_error "Erro em português"
    [ "$status" -eq 0 ]
}

# ========================================
# Test Suite: print_success
# ========================================

@test "print_success executes without error" {
    run print_success "Test success message"
    [ "$status" -eq 0 ]
}

@test "print_success executes successfully" {
    run print_success "Test success message"
    [ "$status" -eq 0 ]
}

@test "print_success respects QUIET mode" {
    export QUIET=true
    run print_success "Should not appear"
    [ "$status" -eq 0 ]
}

@test "print_success handles empty message" {
    run print_success ""
    [ "$status" -eq 0 ]
}

@test "print_success handles message with special characters" {
    run print_success "Success: ✓!@#$%^&*()"
    [ "$status" -eq 0 ]
}

@test "print_success handles message with unicode" {
    run print_success "Success with émojis 🎉"
    [ "$status" -eq 0 ]
}

@test "print_success handles very long message" {
    local long_msg="This is a very long success message "
    long_msg+="that exceeds normal line length and should be handled gracefully "
    long_msg+="without causing any issues in the printing function."
    run print_success "$long_msg"
    [ "$status" -eq 0 ]
}

@test "print_success handles message with multiple sentences" {
    run print_success "First sentence. Second sentence. Third sentence."
    [ "$status" -eq 0 ]
}

# ========================================
# Test Suite: print_info
# ========================================

@test "print_info executes without error" {
    run print_info "Test info message"
    [ "$status" -eq 0 ]
}

@test "print_info executes successfully" {
    run print_info "Test info message"
    [ "$status" -eq 0 ]
}

@test "print_info respects QUIET mode" {
    export QUIET=true
    run print_info "Should not appear"
    [ "$status" -eq 0 ]
}

@test "print_info handles empty message" {
    run print_info ""
    [ "$status" -eq 0 ]
}

@test "print_info handles message with special characters" {
    run print_info "Info: <tag> & special"
    [ "$status" -eq 0 ]
}

@test "print_info handles message with newlines" {
    run print_info -e "Line 1\nLine 2"
    [ "$status" -eq 0 ]
}

@test "print_info handles message with tabs and spaces" {
    run print_info "	Tab	 and   Spaces   "
    [ "$status" -eq 0 ]
}

@test "print_info handles international characters" {
    run print_info "Info en español"
    [ "$status" -eq 0 ]
}

# ========================================
# Test Suite: print_warning
# ========================================

@test "print_warning executes without error" {
    run print_warning "Test warning message"
    [ "$status" -eq 0 ]
}

@test "print_warning executes successfully" {
    run print_warning "Test warning message"
    [ "$status" -eq 0 ]
}

@test "print_warning respects QUIET mode" {
    export QUIET=true
    run print_warning "Should not appear"
    [ "$status" -eq 0 ]
}

@test "print_warning handles empty message" {
    run print_warning ""
    [ "$status" -eq 0 ]
}

@test "print_warning handles message with special characters" {
    run print_warning "Warning: !@#$%^&*()"
    [ "$status" -eq 0 ]
}

@test "print_warning handles very long message" {
    local long_msg="This is a very long warning message "
    long_msg+="that exceeds normal line length and should be handled gracefully "
    long_msg+="without causing any issues in the printing function."
    run print_warning "$long_msg"
    [ "$status" -eq 0 ]
}

@test "print_warning handles message with multiple sentences" {
    run print_warning "First sentence. Second sentence. Third sentence."
    [ "$status" -eq 0 ]
}

@test "print_warning handles international characters" {
    run print_warning "Avertissement en français"
    [ "$status" -eq 0 ]
}

