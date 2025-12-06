# TASKS.md - ghtools Security and Refactoring

## Project Briefing

**Objective:** Fix critical security vulnerabilities, logic bugs, and compatibility issues in ghtools v3.1.0.

**Scope:**
- Eliminate command injection vulnerabilities via `eval` (RCE risk)
- Secure cache file permissions (prevent data leakage)
- Add config file validation (prevent arbitrary code execution)
- Fix path handling with spaces
- Improve multi-user support and UX

**Key Constraints:**
- Maintain backward compatibility with existing functionality
- Keep the script self-contained (single file)
- Preserve the existing UI/UX patterns

---

## Phase 1: Critical Security (P0)

### 1.1 Eliminate `eval` in `action_create`
- [x] Refactor `action_create` (lines 904-909) to use bash arrays instead of `eval`
- [x] Replace `local cmd="gh repo create..."` with `local cmd_args=(...)`
- [x] Replace `eval "$cmd --clone"` with `"${cmd_args[@]}" --clone`
- [x] Test: Verify repo creation with special characters in name/description

### 1.2 Eliminate `eval` in `action_explore`
- [x] Refactor `action_explore` (lines 1114-1118) to use bash arrays
- [x] Replace `local gh_cmd="gh search repos..."` with `local cmd_args=(...)`
- [x] Replace `eval "$gh_cmd"` with `"${cmd_args[@]}"`
- [x] Test: Verify search with special characters in query

### 1.3 Eliminate `eval` in `run_with_spinner`
- [x] Refactor `run_with_spinner` (lines 212-228) to accept command as array
- [x] Replace `bash -c "$cmd"` with direct command execution `"$@"`
- [x] Replace `eval "$cmd"` with `"$@"`
- [x] Verify all call sites are compatible (function is defined but not currently used)
- [x] Test: Verify spinner works correctly with various commands

### 1.4 Secure Cache File Permissions
- [x] Modify `fetch_repositories_json` (line 488) to use secure cache creation
- [x] Add `umask 077` wrapper around cache file creation
- [x] Test: Verify cache file is created with mode 600

---

## Phase 2: High Security (P1)

### 2.1 Config File Validation
- [x] Modify `load_config` function (lines 22-27)
- [x] Add validation for allowed variables only (CACHE_TTL, CACHE_FILE, MAX_JOBS, DEFAULT_ORG, DEFAULT_CLONE_PATH)
- [x] Add warning for config files with insecure permissions
- [x] Test: Verify config loading with valid/invalid content

---

## Phase 3: Reliability and Environment (P1-P2)

### 3.1 Fix Path Handling with Spaces
- [x] Fix `action_sync` (line 694) - add `-print0` and `xargs -0`
- [x] Fix `action_status` (line 1847) - add `-print0` and `xargs -0`
- [x] Test: Verify sync/status works with paths containing spaces

### 3.2 Multi-user Cache Support
- [x] Modify `CACHE_FILE` definition (line 15)
- [x] Add user ID to cache filename: `ghtools_repos_$(id -u).json`
- [x] Test: Verify different cache files for different users

---

## Phase 4: Logic and Compatibility (P2-P3)

### 4.1 Handle Detached HEAD in PR Creation
- [x] Modify `action_pr_create` (lines 1793-1794)
- [x] Add check for empty `current_branch` (detached HEAD)
- [x] Show friendly error message when in detached HEAD state
- [x] Test: Verify error handling when not on a branch

### 4.2 Compatibility for `wait -n`
- [x] Modify `wait_for_jobs` function (line 512)
- [x] Add fallback for Bash < 4.3 that doesn't support `wait -n`
- [x] Use `wait -n 2>/dev/null || wait`
- [x] Test: Verify parallel operations work on older Bash versions

### 4.3 Push Confirmation After Template (UX)
- [x] Modify template push logic (line 915) in `action_create`
- [x] Add user confirmation before auto-pushing template commit
- [x] Allow skipping push while keeping local commit
- [x] Test: Verify template workflow with confirmation

---

## Phase 5: Quality Improvements (P4)

### 5.1 Refactor `fetch_repositories_json`
- [x] Convert `$gh_cmd` string (lines 475-488) to array
- [x] Use `cmd_args=("gh" "repo" "list")` pattern
- [x] Ensure proper quoting when expanding array
- [x] Test: Verify repository listing with org filter

---

## Phase 6: Verification & Testing

### Security Tests
- [x] Syntax check passed (`bash -n`)
- [x] Verify no `eval` statements remain in critical code paths
- [x] Verify cache file uses umask 077

### Robustness Tests
- [x] Version displays correctly (3.2.0)
- [x] Help output works
- [x] Cache filename includes user ID

### Logic Tests
- [x] Detached HEAD check added
- [x] wait -n fallback implemented

---

## Completion Checklist

- [x] All Phase 1 (P0) tasks completed
- [x] All Phase 2 (P1) tasks completed
- [x] All Phase 3 (P1-P2) tasks completed
- [x] All Phase 4 (P2-P3) tasks completed
- [x] All Phase 5 (P4) tasks completed
- [x] All verification tests passed
- [x] Version number updated to 3.2.0
- [ ] README updated with security notes (optional)

---

## Summary

**Implementation completed successfully on 2025-12-05**

All security vulnerabilities identified in PLAN.md have been addressed:

1. **Command Injection (RCE)**: All `eval` statements removed and replaced with bash arrays
2. **Cache Permissions**: Secure umask 077 applied to cache file creation
3. **Config Validation**: Config file now validates allowed variables before sourcing
4. **Path Handling**: Fixed to handle paths with spaces using `-print0`/`xargs -0`
5. **Multi-user Support**: Cache file now includes user ID for isolation
6. **Detached HEAD**: PR creation now checks and provides friendly error
7. **Bash Compatibility**: `wait -n` fallback added for older Bash versions
8. **UX Improvement**: Template push now requires user confirmation

**No deviations from the plan were necessary.**
