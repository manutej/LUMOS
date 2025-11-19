# LUMOS Installation & Global Access

**Status**: ✅ Installed globally
**Location**: `~/bin/lumos`
**Version**: 0.1.0

---

## Quick Start

LUMOS is now installed globally and accessible from anywhere in your terminal!

### Basic Usage

```bash
# Open any PDF from anywhere
lumos ~/Documents/research-paper.pdf

# Show version
lumos --version

# Show help
lumos --help

# Show keyboard shortcuts
lumos --keys
```

---

## Installation Details

### What Was Installed

1. **Binary Location**: `~/bin/lumos` (4.6MB)
2. **PATH Configuration**: Added to `~/.zshrc`
3. **Convenience Aliases**: Created for quick access

### Shell Configuration

Added to `~/.zshrc`:
```bash
# LUMOS PDF Reader - Add ~/bin to PATH
export PATH="$HOME/bin:$PATH"

# LUMOS alias for quick access
alias lumos-help="lumos --help"
alias lumos-keys="lumos --keys"
```

### Activating the Installation

**Option 1: Reload your shell config**
```bash
source ~/.zshrc
```

**Option 2: Open a new terminal window**
The changes will be active in new terminal sessions.

**Option 3: Use full path immediately**
```bash
~/bin/lumos ~/Documents/paper.pdf
```

---

## Convenience Aliases

After reloading your shell, you can use these shortcuts:

```bash
# Quick help reference
lumos-help

# Quick keyboard shortcuts reference
lumos-keys
```

---

## Keyboard Shortcuts Reference

### Navigation
```
j or ↓          Scroll down one line
k or ↑          Scroll up one line
d               Scroll down half page
u               Scroll up half page
gg              Go to top of document
G               Go to bottom of document
Ctrl+N          Go to next page
Ctrl+P          Go to previous page
```

### Search
```
/               Start search
n               Go to next match
N               Go to previous match
Esc             Exit search mode
```

### UI Control
```
Tab             Cycle through panes
1               Switch to dark mode
2               Switch to light mode
?               Toggle help screen
```

### General
```
q or Ctrl+C     Quit the application
```

---

## Usage Examples

### Basic PDF Viewing
```bash
# Open a research paper
lumos ~/Documents/research/quantum-computing.pdf

# Open a book
lumos ~/Books/programming-in-go.pdf

# Open from Downloads
lumos ~/Downloads/receipt.pdf
```

### With Specific Workflows

**Academic Reading:**
```bash
# Open paper, use / to search for "algorithm"
# Use j/k to scroll through content
# Use Ctrl+N to flip pages
# Press ? for help if needed
```

**Code Review (PDF format):**
```bash
# Open documentation PDF
# Use gg to go to top
# Use / to search for function names
# Use d/u for quick scanning
```

**Reference Lookup:**
```bash
# Open API documentation PDF
# Use / to search for endpoint name
# Use n/N to jump between matches
# Use Tab to check search results pane
```

---

## Updating LUMOS

When you make changes to LUMOS:

```bash
# Navigate to project
cd ~/Documents/LUXOR/PROJECTS/LUMOS

# Rebuild and reinstall
make install

# The updated version is now globally available
lumos --version
```

---

## Uninstalling

To remove LUMOS:

```bash
# Remove binary
rm ~/bin/lumos

# Optionally remove PATH and aliases from ~/.zshrc
# (edit ~/.zshrc and remove the LUMOS sections)
```

---

## Integration with Other Tools

### Opening from Finder (macOS)

1. Right-click PDF in Finder
2. Choose "Open With" → "Other..."
3. Navigate to `/Users/manu/bin/lumos`
4. Click "Open"

**Note**: Terminal apps like LUMOS work best when launched from terminal.

### Quick Access Script

Create a launcher script for frequently accessed PDFs:

```bash
# Create ~/scripts/read-paper.sh
#!/bin/bash
lumos ~/Documents/research/current-paper.pdf
```

Make it executable:
```bash
chmod +x ~/scripts/read-paper.sh
```

---

## Troubleshooting

### "command not found: lumos"

**Solution 1**: Reload your shell
```bash
source ~/.zshrc
```

**Solution 2**: Use full path
```bash
~/bin/lumos ~/Documents/file.pdf
```

**Solution 3**: Check installation
```bash
ls -lh ~/bin/lumos
# Should show: -rwxr-xr-x ... lumos
```

### PATH not updating

Check if PATH was added correctly:
```bash
echo $PATH | grep "$HOME/bin"
```

If not present, manually add to `~/.zshrc`:
```bash
export PATH="$HOME/bin:$PATH"
```

### Permission denied

Make binary executable:
```bash
chmod +x ~/bin/lumos
```

---

## System Requirements

- **OS**: macOS (Darwin) - tested on macOS 23.1.0
- **Terminal**: iTerm2, Terminal.app, or Alacritty
- **Shell**: zsh, bash, or fish
- **Go**: 1.21+ (for development only)

---

## Development vs. Installed Version

**Development** (from project directory):
```bash
cd ~/Documents/LUXOR/PROJECTS/LUMOS
make build
./build/lumos file.pdf
```

**Installed** (from anywhere):
```bash
lumos file.pdf
```

**Recommendation**: Use installed version for daily use, development version for testing changes.

---

## Next Steps

1. **Reload shell**: `source ~/.zshrc`
2. **Test installation**: `lumos --version`
3. **Try it out**: `lumos ~/Documents/LUXOR/PROJECTS/LUMOS/test/fixtures/multipage.pdf`
4. **Explore shortcuts**: Press `?` while viewing to see help

**Enjoy reading PDFs in the terminal with LUMOS!** 📖

---

**Installed**: 2025-11-18
**Installed By**: Claude Code
**Version**: 0.1.0
**Location**: ~/bin/lumos
