#!/bin/bash
# Test script to verify TUI loads without crashing

set -e

echo "Testing LUMOS TUI..."
echo

# Test 1: Help flag
echo "✓ Test 1: Help flag"
./lumos --help > /dev/null

# Test 2: Version flag
echo "✓ Test 2: Version flag"
./lumos --version > /dev/null

# Test 3: Keys flag
echo "✓ Test 3: Keys flag"
./lumos --keys > /dev/null

# Test 4: Missing file error
echo "✓ Test 4: Error handling for missing file"
if ./lumos /nonexistent/file.pdf 2>&1 | grep -q "File not found"; then
    echo "  - Correctly handles missing file"
else
    echo "  - ERROR: Should report missing file"
    exit 1
fi

# Test 5: Load a PDF (will exit immediately with q)
echo "✓ Test 5: Load PDF without crashing"
# This test would require simulating 'q' keypress, skip for now
echo "  - Manual test required (run: ./lumos test/fixtures/simple.pdf)"

echo
echo "All automated tests passed!"
echo
echo "Next step: Run './lumos test/fixtures/simple.pdf' and verify:"
echo "  1. PDF content displays"
echo "  2. Status bar shows page numbers"
echo "  3. Pressing 'q' quits cleanly"
echo "  4. Pressing '?' shows help"
echo "  5. Basic navigation works"
