# PHASE 1 MILESTONE 1.1: Build & Compile Foundation

**Quick Reference Guide**
**Duration**: 2-3 days
**Goal**: Get the project compiling and running

---

## Quick Start (5 minutes)

```bash
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMOS

# 1. Check Go version
go version              # Should be 1.21+

# 2. Download dependencies
go mod download

# 3. Build
make build

# 4. Test basic run
./build/lumos ~/Documents/test.pdf

# 5. Try keybindings
# Press ? for help
# j/k to scroll
# q to quit
```

---

## Detailed Checklist

### ✅ Task 1.1.1: Verify Go Environment

```bash
# Check Go version (need 1.21+)
go version

# Check GOPATH
echo $GOPATH

# Check workspace
pwd
# Should show: /Users/manu/Documents/LUXOR/PROJECTS/LUMOS
```

**Expected Output**:
```
go version go1.21.0 darwin/arm64
/Users/manu/go
/Users/manu/Documents/LUXOR/PROJECTS/LUMOS
```

**Troubleshooting**:
- If Go not installed: https://golang.org/doc/install
- If wrong version: Update Go to 1.21+

---

### ✅ Task 1.1.2: Download Dependencies

```bash
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMOS

# Download all dependencies (will create/update go.sum)
go mod download

# Verify dependencies
go mod verify

# List dependencies
go mod graph
```

**Expected**:
- `go.sum` is updated
- No errors from `go mod verify`

**Troubleshooting**:
```bash
# If issues, try:
go clean -modcache
go mod download

# Check network connection
ping github.com
```

---

### ✅ Task 1.1.3: Build the Project

```bash
# Option 1: Use Makefile (recommended)
make build
# Creates: ./build/lumos

# Option 2: Manual build
go build -o ./build/lumos ./cmd/lumos

# Verify binary was created
ls -lh ./build/lumos
```

**Expected**:
- Binary created at `./build/lumos`
- File size ~10-20MB
- No compilation errors

**If you get errors**:

1. **Import errors** (can't find package):
   ```bash
   go mod tidy
   go mod download
   make build
   ```

2. **Compilation errors** (syntax issues):
   - Check error message carefully
   - Look at file mentioned in error
   - May need to fix code (document issue)

3. **Linker errors** (can't link):
   ```bash
   go clean
   make build
   ```

---

### ✅ Task 1.1.4: Basic Functionality Test

```bash
# Find a test PDF
find ~/ -name "*.pdf" -type f 2>/dev/null | head -5

# Or create a simple one
# Use this path if you have it: ~/Documents/

# Run with PDF
./build/lumos ~/Documents/your-file.pdf

# Or use an example from the system
./build/lumos /System/Library/CoreServices/System\ Information.app/Contents/Resources/English.lproj/system.pdf
```

**Expected**:
- LUMOS window opens
- Shows PDF with 3 panes
- Dark mode display
- No crashes

**Common Issues & Solutions**:

| Issue | Solution |
|-------|----------|
| PDF won't open | Check file exists: `ls -l /path/to/file.pdf` |
| Crash on load | Check file is valid PDF: `file /path/to/file.pdf` |
| No display | Try smaller terminal window or `reset` command |
| Garbled text | Terminal encoding issue - try: `export LANG=en_US.UTF-8` |

---

### ✅ Task 1.1.5: Test with Diverse PDFs

Create a test suite by gathering PDFs:

```bash
# Create test directory
mkdir -p test/fixtures

# Try different PDF types:
# 1. Simple text PDF (create manually if needed)
# 2. Multi-page PDF
# 3. PDF with images
# 4. PDF with tables
# 5. Large PDF (100+ pages)

# Test each one
for pdf in test/fixtures/*.pdf; do
  echo "Testing: $pdf"
  timeout 5 ./build/lumos "$pdf" < /dev/null
  echo "Status: $?"
done
```

**Document Your Findings**:
- Which PDFs load successfully
- Which PDFs cause issues
- Any crashes or errors
- File in: `PHASE_1_MILESTONE_1_1_TEST_RESULTS.md`

---

### ✅ Task 1.1.6: Verify All Keybindings Work

While the app is running, test each binding:

| Key | Expected Behavior | Works? |
|-----|-------------------|--------|
| `?` | Show help | [ ] |
| `j` | Scroll down | [ ] |
| `k` | Scroll up | [ ] |
| `d` | Half page down | [ ] |
| `u` | Half page up | [ ] |
| `g` | Start gg command | [ ] |
| `G` | End of document | [ ] |
| `Ctrl+N` | Next page | [ ] |
| `Ctrl+P` | Previous page | [ ] |
| `/` | Start search | [ ] |
| `Esc` | Exit search | [ ] |
| `Tab` | Switch pane | [ ] |
| `1` | Dark mode | [ ] |
| `2` | Light mode | [ ] |
| `q` | Quit | [ ] |

**Test Script**:
```bash
# Quick verification (exit after each test)
./build/lumos ~/Documents/test.pdf << 'EOF'
?
j
k
d
u
gg
G
q
EOF
```

---

### ✅ Task 1.1.7: CLI Argument Handling

Test command-line arguments:

```bash
# Help text
./build/lumos --help

# Version
./build/lumos --version

# Keybindings reference
./build/lumos --keys

# Invalid file (should show error)
./build/lumos /nonexistent/file.pdf

# Home directory expansion
./build/lumos ~/Documents/test.pdf

# With relative path
cd ~/Documents && /Users/manu/Documents/LUXOR/PROJECTS/LUMOS/build/lumos test.pdf
```

**Expected Results**:

| Command | Expected Output |
|---------|-----------------|
| `--help` | Usage information |
| `--version` | Version number (e.g., v0.1.0-alpha) |
| `--keys` | List of vim keybindings |
| `/nonexistent.pdf` | Error: "file not found" or similar |
| `~/path.pdf` | Opens the file |

---

## Success Criteria Checklist

After completing all tasks:

- [ ] Go version is 1.21+
- [ ] Dependencies downloaded successfully
- [ ] Binary builds without errors
- [ ] Binary is created at `./build/lumos`
- [ ] Can open and display a PDF
- [ ] No crashes during basic operations
- [ ] All keybindings respond (can test from help)
- [ ] `--help` flag works
- [ ] `--version` flag works
- [ ] `--keys` flag works
- [ ] CLI argument parsing works
- [ ] Home directory expansion works
- [ ] Error handling works (bad file path)

---

## Common Problems & Solutions

### Problem: `go mod download` fails

```bash
# Check internet connection
ping github.com

# Try clearing cache
go clean -modcache

# Set Go proxy (if behind proxy)
export GO111MODULE=on
go mod download

# Check git config (might need GitHub credentials)
git config --list | grep github
```

### Problem: `make build` fails - "command not found"

```bash
# Make sure you're in the right directory
pwd  # Should show: /Users/manu/Documents/LUXOR/PROJECTS/LUMOS

# Check Makefile exists
ls -la Makefile

# If Makefile doesn't exist, run manual build:
go build -o ./build/lumos ./cmd/lumos
```

### Problem: Binary runs but crashes immediately

```bash
# Get more info
./build/lumos 2>&1 | head -20

# Try with debug output
export GODEBUG=msan=1
./build/lumos ~/Documents/test.pdf

# Check for panic
./build/lumos ~/Documents/test.pdf > /tmp/lumos.log 2>&1
cat /tmp/lumos.log
```

### Problem: "Permission denied" when running

```bash
# Fix permissions
chmod +x ./build/lumos

# Verify
ls -la ./build/lumos  # Should show: -rwxr-xr-x
```

### Problem: Terminal display issues

```bash
# Reset terminal
reset

# Try different terminal
# Use Alacritty, iTerm2, or macOS Terminal

# Check terminal size
stty size

# Try in full terminal (not split pane)
```

---

## Creating Test PDF

If you don't have test PDFs, you can:

1. **Download sample PDFs**:
   ```bash
   # Example PDFs online:
   # https://www.w3.org/WAI/WCAG21/Techniques/pdf/img/table.pdf
   # https://www.w3.org/WAI/WCAG21/Techniques/pdf/PDFTableStructure.pdf

   curl -o test/fixtures/sample.pdf "https://www.w3.org/WAI/WCAG21/Techniques/pdf/PDFTableStructure.pdf"
   ```

2. **Use system PDFs**:
   ```bash
   # macOS has some sample PDFs
   find /System /Library /Applications -name "*.pdf" 2>/dev/null | head -10
   ```

3. **Use existing PDFs**:
   - Check your Documents folder
   - Check Downloads folder
   - Look in project documentation

---

## What to Document in Review

After completing all tasks, document these in `PHASE_1_MILESTONE_1_1_REVIEW.md`:

1. **What Compiled Successfully**
   - Go version used
   - Dependencies resolved
   - Build time

2. **PDFs Tested**
   - File names
   - File types
   - Results (success/failure)

3. **Keybindings Verified**
   - Which ones work
   - Which ones need fixes

4. **Issues Encountered**
   - Any errors during build
   - Any crashes during testing
   - Any UI issues

5. **CLI Testing Results**
   - Help flag: ✅/❌
   - Version flag: ✅/❌
   - Keys flag: ✅/❌
   - Argument parsing: ✅/❌

6. **Next Steps**
   - Any blockers
   - Recommended fixes
   - Ready for Milestone 1.2: ✅/❌

---

## Files to Create/Update

### Create These Files
- `PHASE_1_MILESTONE_1_1_REVIEW.md` - Milestone completion review
- `PHASE_1_MILESTONE_1_1_TEST_RESULTS.md` - Detailed test results
- `test/fixtures/` - Add test PDF files

### Update These Files
- `README.md` - Add actual performance metrics once measured
- `QUICKSTART.md` - Update with any issues found

---

## Next Steps After 1.1

Once Milestone 1.1 is complete:

1. **Create review document** - `PHASE_1_MILESTONE_1_1_REVIEW.md`
2. **Document any issues** - List bugs/problems found
3. **Pause for review** - Assess if ready for 1.2
4. **If blockers** - Fix issues before moving to 1.2
5. **If success** - Move to PHASE_1_MILESTONE_1_2_GUIDE.md

---

## Timeline Estimate

```
Day 1:
  - Task 1.1.1-1.1.3: 30 min (setup, download, build)
  - Task 1.1.4: 30 min (basic test)
  - Task 1.1.5: 1 hour (diverse PDFs)

Day 2:
  - Task 1.1.6: 1 hour (keybinding verification)
  - Task 1.1.7: 30 min (CLI args)
  - Documentation: 1 hour (write review)

Total: ~5 hours over 2 days
```

---

## Success Indicators

✅ **Milestone 1.1 is complete when**:
1. Binary builds without errors
2. Can open and display PDFs
3. No crashes on basic operations
4. All keybindings work
5. CLI arguments work
6. Review document is written
7. You feel confident to move to 1.2

---

**Start Time**: [Your date/time]
**Expected Completion**: [2-3 days from start]
**Status**: Ready to begin ✅

