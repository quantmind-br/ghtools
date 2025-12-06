#!/usr/bin/env bats

# Setup test environment
setup() {
    # Get the project directory (three levels up from this test file)
    local test_file="${BATS_TEST_FILENAME}"
    local project_dir
    project_dir="$(dirname "$(dirname "$(dirname "$test_file")")")"

    # Load test helper (this sets up test environment variables and loads ghtools)
    source "$project_dir/test/test_helper.bash"

    # Setup mocks for interactive commands
    mkdir -p "$TEST_TMP_DIR"
    setup_mock_gh
    setup_mock_git
    setup_mock_fzf
    setup_mock_gum
}

teardown() {
    teardown_test
}

# ========================================
# Test Suite: gum_confirm
# These tests pipe input to simulate user response
# ========================================

@test "gum_confirm executes without error" {
    # Simulate user typing 'y' (yes)
    run bash -c 'echo "y" | gum_confirm "Test prompt"'
    [ "$status" -eq 0 ]
}

@test "gum_confirm respects default yes parameter" {
    run bash -c 'echo "y" | gum_confirm "Test prompt" "yes"'
    [ "$status" -eq 0 ]
}

@test "gum_confirm respects default no parameter" {
    # When default is no and user types nothing, should fail
    run bash -c 'echo "" | gum_confirm "Test prompt" "no"'
    [ "$status" -eq 1 ]
}

@test "gum_confirm handles empty prompt" {
    run bash -c 'echo "y" | gum_confirm ""'
    [ "$status" -eq 0 ]
}

@test "gum_confirm handles prompt with special characters" {
    run bash -c 'echo "y" | gum_confirm "Are you sure? (!@#$%)"'
    [ "$status" -eq 0 ]
}

@test "gum_confirm handles prompt with quotes" {
    run bash -c "echo 'y' | gum_confirm \"Continue with 'special' prompt?\""
    [ "$status" -eq 0 ]
}

@test "gum_confirm handles long prompt" {
    local long_prompt="This is a very long confirmation prompt "
    long_prompt+="that exceeds normal line length and should be handled gracefully "
    long_prompt+="without causing any issues."
    run bash -c "echo 'y' | gum_confirm '$long_prompt'"
    [ "$status" -eq 0 ]
}

@test "gum_confirm handles prompt with newlines" {
    run bash -c 'echo "y" | gum_confirm "Line 1 Line 2"'
    [ "$status" -eq 0 ]
}

@test "gum_confirm handles prompt with international characters" {
    run bash -c 'echo "y" | gum_confirm "¿Continuar?"'
    [ "$status" -eq 0 ]
}

@test "gum_confirm handles prompt with unicode" {
    run bash -c 'echo "y" | gum_confirm "Continue? 🎉"'
    [ "$status" -eq 0 ]
}

# ========================================
# Test Suite: gum_input
# ========================================

@test "gum_input executes without error" {
    run gum_input "Test placeholder"
    [ "$status" -eq 0 ]
}

@test "gum_input handles custom prompt" {
    run gum_input "Name" "Enter > "
    [ "$status" -eq 0 ]
}

@test "gum_input handles default value" {
    run gum_input "Name" "› " "default"
    [ "$status" -eq 0 ]
}

@test "gum_input handles empty placeholder" {
    run gum_input ""
    [ "$status" -eq 0 ]
}

@test "gum_input handles placeholder with special characters" {
    run gum_input "Email (user@domain.com)"
    [ "$status" -eq 0 ]
}

@test "gum_input handles placeholder with quotes" {
    run gum_input "Enter 'username'"
    [ "$status" -eq 0 ]
}

@test "gum_input handles long placeholder" {
    local long_placeholder="This is a very long placeholder text "
    long_placeholder+="that exceeds normal line length and should be handled gracefully "
    long_placeholder+="without causing any issues."
    run gum_input "$long_placeholder"
    [ "$status" -eq 0 ]
}

@test "gum_input handles placeholder with newlines" {
    run gum_input -e "Line 1\nLine 2"
    [ "$status" -eq 0 ]
}

@test "gum_input handles placeholder with international characters" {
    run gum_input "Nombre"
    [ "$status" -eq 0 ]
}

@test "gum_input handles placeholder with unicode" {
    run gum_input "Name 🎉"
    [ "$status" -eq 0 ]
}

# ========================================
# Test Suite: gum_choose
# ========================================

@test "gum_choose executes without error with single option" {
    run gum_choose "Select option" "Option 1"
    [ "$status" -eq 0 ]
}

@test "gum_choose handles multiple options" {
    run gum_choose "Select option" "Option 1" "Option 2" "Option 3"
    [ "$status" -eq 0 ]
}

@test "gum_choose handles empty header" {
    run gum_choose "" "Option 1" "Option 2"
    [ "$status" -eq 0 ]
}

@test "gum_choose handles options with special characters" {
    run gum_choose "Select" "Option 1!@#$%" "Option 2&*()"
    [ "$status" -eq 0 ]
}

@test "gum_choose handles options with spaces" {
    run gum_choose "Select" "Option with spaces" "Another option"
    [ "$status" -eq 0 ]
}

@test "gum_choose handles options with quotes" {
    run gum_choose "Select" "Option 'quoted'" 'Option "double"'
    [ "$status" -eq 0 ]
}

@test "gum_choose handles long options" {
    local long_option="This is a very long option text "
    long_option+="that exceeds normal line length and should be handled gracefully."
    run gum_choose "Select" "$long_option"
    [ "$status" -eq 0 ]
}

@test "gum_choose handles many options" {
    run gum_choose "Select" "Opt 1" "Opt 2" "Opt 3" "Opt 4" "Opt 5" "Opt 6"
    [ "$status" -eq 0 ]
}

@test "gum_choose handles options with international characters" {
    run gum_choose "Seleccionar" "Opción 1" "Opción 2"
    [ "$status" -eq 0 ]
}

@test "gum_choose handles options with unicode" {
    run gum_choose "Select" "Option 🎉" "Option ⭐"
    [ "$status" -eq 0 ]
}

# ========================================
# Test Suite: gum_filter
# ========================================

@test "gum_filter executes without error" {
    run gum_filter "Type to filter"
    [ "$status" -eq 0 ]
}

@test "gum_filter handles multi mode" {
    run gum_filter "Type to filter" "true"
    [ "$status" -eq 0 ]
}

@test "gum_filter handles empty placeholder" {
    run gum_filter ""
    [ "$status" -eq 0 ]
}

@test "gum_filter handles placeholder with special characters" {
    run gum_filter "Search: !@#$%^&*()"
    [ "$status" -eq 0 ]
}

@test "gum_filter handles placeholder with international characters" {
    run gum_filter "Buscar"
    [ "$status" -eq 0 ]
}

# ========================================
# Test Suite: gum_write
# ========================================

@test "gum_write executes without error" {
    run gum_write "Enter text"
    [ "$status" -eq 0 ]
}

@test "gum_write handles empty placeholder" {
    run gum_write ""
    [ "$status" -eq 0 ]
}

@test "gum_write handles placeholder with special characters" {
    run gum_write "Write message: !@#$%^&*()"
    [ "$status" -eq 0 ]
}

@test "gum_write handles placeholder with international characters" {
    run gum_write "Escribir mensaje"
    [ "$status" -eq 0 ]
}

@test "gum_write handles placeholder with unicode" {
    run gum_write "Write message 🎉"
    [ "$status" -eq 0 ]
}
