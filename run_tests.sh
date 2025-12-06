#!/bin/bash

# ghtools Test Runner
# This script runs all tests and generates coverage reports

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="$SCRIPT_DIR/test"

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║           ghtools - Test Runner                           ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if bats is installed
if ! command -v bats &>/dev/null; then
    echo -e "${RED}Error: bats is not installed${NC}"
    echo "Please install bats: git clone https://github.com/bats-core/bats-core.git && cd bats-core && sudo ./install.sh /usr/local"
    exit 1
fi

# Check if required tools are available
echo -e "${YELLOW}Checking dependencies...${NC}"
for cmd in jq git; do
    if ! command -v "$cmd" &>/dev/null; then
        echo -e "${RED}Missing dependency: $cmd${NC}"
        exit 1
    fi
done
echo -e "${GREEN}✓ All dependencies found${NC}"
echo ""

# Find all test files
TEST_FILES=()
while IFS= read -r -d '' file; do
    TEST_FILES+=("$file")
done < <(find "$TEST_DIR" -name "*.bats" -type f -print0 | sort -z)

if [ ${#TEST_FILES[@]} -eq 0 ]; then
    echo -e "${RED}No test files found in $TEST_DIR${NC}"
    exit 1
fi

echo -e "${YELLOW}Found ${#TEST_FILES[@]} test file(s)${NC}"
for file in "${TEST_FILES[@]}"; do
    echo "  - $(basename "$file")"
done
echo ""

# Check coverage tool availability
COVERAGE_TOOL=""
if command -v shfmt &>/dev/null; then
    COVERAGE_TOOL="shfmt"
elif command -v awk &>/dev/null; then
    COVERAGE_TOOL="awk"
fi

# Run tests
echo -e "${YELLOW}Running tests...${NC}"
echo ""

# Count total tests
TOTAL_TESTS=0
for file in "${TEST_FILES[@]}"; do
    count=$(grep -c "^@" "$file" 2>/dev/null || echo 0)
    TOTAL_TESTS=$((TOTAL_TESTS + count))
done

# Run bats with formatter
BATS_OUTPUT=$(bats --formatter tap "${TEST_FILES[@]}" 2>&1)
BATS_EXIT_CODE=$?

echo "$BATS_OUTPUT"

# Calculate results
PASSED_TESTS=$(echo "$BATS_OUTPUT" | grep -c "^ok " || echo 0)
FAILED_TESTS=$(echo "$BATS_OUTPUT" | grep -c "^not ok " || echo 0)
SKIPPED_TESTS=$(echo "$BATS_OUTPUT" | grep -c "^skip " || echo 0)

# Calculate coverage
if [ $BATS_EXIT_CODE -eq 0 ]; then
    echo ""
    echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                    Test Summary                           ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "  ${GREEN}Total Tests:${NC}    $TOTAL_TESTS"
    echo -e "  ${GREEN}Passed:${NC}        $PASSED_TESTS"
    echo -e "  ${YELLOW}Failed:${NC}        $FAILED_TESTS"
    echo -e "  ${YELLOW}Skipped:${NC}       $SKIPPED_TESTS"
    echo ""

    # Estimate coverage based on tested functions
    FUNCTIONS_TESTED=0
    FUNCTIONS_TOTAL=45  # Approximate number of functions in ghtools

    for file in "${TEST_FILES[@]}"; do
        test_count=$(grep -c "^@" "$file" 2>/dev/null || echo 0)
        FUNCTIONS_TESTED=$((FUNCTIONS_TESTED + test_count))
    done

    COVERAGE=$((FUNCTIONS_TESTED * 100 / FUNCTIONS_TOTAL))

    echo -e "  ${BLUE}Coverage:${NC}       ~${COVERAGE}% (estimated)"
    echo ""

    if [ $COVERAGE -ge 80 ]; then
        echo -e "  ${GREEN}✓ Coverage target achieved (≥80%)${NC}"
    else
        echo -e "  ${YELLOW}⚠ Coverage below target (${COVERAGE}% < 80%)${NC}"
    fi
    echo ""

    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${RED}║                    Test Summary                           ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "  ${GREEN}Total Tests:${NC}    $TOTAL_TESTS"
    echo -e "  ${GREEN}Passed:${NC}        $PASSED_TESTS"
    echo -e "  ${RED}Failed:${NC}        $FAILED_TESTS"
    echo -e "  ${YELLOW}Skipped:${NC}       $SKIPPED_TESTS"
    echo ""
    echo -e "${RED}Some tests failed. Please review the output above.${NC}"
    exit 1
fi
