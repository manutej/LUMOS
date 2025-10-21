# LUMOS Quick Start Guide

**Version**: 0.1.0
**Date**: 2025-10-21
**Status**: Ready for Phase 1 Development

---

## 5-Minute Setup

### 1. Build the Project

```bash
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMOS

# Download dependencies
go mod download

# Build binary
make build

# Binary is now at: ./build/lumos
```

### 2. Get a Test PDF

```bash
# Option 1: Use an existing PDF
lumos ~/Documents/your-file.pdf

# Option 2: Create a simple test PDF (using pdfcpu - Phase 2)
# For now, find any PDF on your system
```

### 3. Navigate with Vim Keys

```
j/k   - Scroll up/down
d/u   - Half page up/down
gg/G  - Top/bottom of document
q     - Quit
?     - Show help
```

---

## Project Structure Overview

```
LUMOS/
├── cmd/lumos/
│   └── main.go              # CLI entry point
│
├── pkg/
│   ├── pdf/
│   │   ├── document.go      # PDF loading & page extraction
│   │   ├── search.go        # Text search
│   │   └── cache.go         # LRU cache
│   │
│   ├── ui/
│   │   ├── model.go         # Bubble Tea model (main logic)
│   │   └── keybindings.go   # Vim keybindings
│   │
│   ├── config/
│   │   └── theme.go         # Dark/light mode themes
│   │
│   └── terminal/            # (Phase 3+: image support)
│
├── test/fixtures/           # Test PDFs
├── docs/                    # Documentation
│
├── README.md                # Project overview
├── go.mod                   # Dependencies
├── Makefile                 # Build commands
└── QUICKSTART.md            # This file
```

---

## Common Commands

### Building

```bash
make build              # Build binary
make build-all         # Build for macOS/Linux
make install           # Install to ~/bin/
```

### Running & Testing

```bash
make run               # Build and run
make test              # Run all tests
make test-v            # Verbose test output
make coverage          # Generate coverage report
```

### Performance

```bash
make bench             # Run benchmarks
make profile-cpu       # CPU profiling
make profile-mem       # Memory profiling
```

### Code Quality

```bash
make fmt               # Format code
make vet               # Lint with go vet
make lint              # Full linter check (requires golangci-lint)
```

### Cleanup

```bash
make clean             # Remove build artifacts
make clean-all         # Clean everything including cache
```

---

## Development Workflow

### 1. Create Feature Branch

```bash
git checkout -b feature/my-feature
```

### 2. Make Changes

Edit files in `pkg/pdf`, `pkg/ui`, or `cmd/`

### 3. Test Your Changes

```bash
# Build
make build

# Run
./build/lumos ~/Documents/test.pdf

# Run tests
make test
```

### 4. Commit & Push

```bash
git add .
git commit -m "feat: description of changes"
git push origin feature/my-feature
```

---

## Understanding the Code

### Entry Point: `cmd/lumos/main.go`
- Parses command line arguments
- Loads PDF file
- Creates Bubble Tea UI model
- Runs the application

### Core Model: `pkg/pdf/document.go`
- Loads PDF using `ledongthuc/pdf`
- Extracts text from pages
- Implements LRU caching
- Provides search API

### UI Layer: `pkg/ui/model.go`
- Bubble Tea MVU pattern implementation
- Handles keyboard input
- Manages 3-pane layout
- Processes navigation and search

### Styling: `pkg/config/theme.go`
- Dark theme (default, VSCode Dark+)
- Light theme (alternative)
- All colors defined here
- Easy to add new themes

---

## Testing Your Changes

### Unit Tests

```bash
# Run all tests
go test ./...

# Test specific package
go test ./pkg/pdf/

# Verbose output
go test -v ./...

# With race detector
go test -race ./...
```

### Manual Testing

```bash
# Build current code
make build

# Test with various PDFs
./build/lumos ~/Documents/simple.pdf
./build/lumos ~/Documents/complex.pdf
./build/lumos ~/Documents/large.pdf
```

### Performance Testing

```bash
# Check startup time
time ./build/lumos ~/Documents/test.pdf

# Profile
make profile-cpu
make profile-mem
```

---

## Architecture Overview

```
User Input (Keyboard)
        ↓
KeyHandler (pkg/ui/keybindings.go)
        ↓
Model.Update() (pkg/ui/model.go)
        ↓
Async Commands (e.g., LoadPageCmd)
        ↓
PDF Processing (pkg/pdf/document.go)
        ↓
LRU Cache (pkg/pdf/cache.go)
        ↓
Model receives PageLoadedMsg
        ↓
Model.View() renders TUI
        ↓
Terminal displays 3-pane layout
```

---

## Key Files to Understand

### Start Here (Reading Order)

1. **README.md** (5 min)
   - Project overview
   - Quick start
   - Feature list

2. **cmd/lumos/main.go** (5 min)
   - CLI argument parsing
   - PDF loading
   - TUI launch

3. **pkg/pdf/document.go** (10 min)
   - PDF model
   - Page extraction
   - Caching logic

4. **pkg/ui/model.go** (10 min)
   - Bubble Tea pattern
   - State management
   - Rendering

5. **docs/ARCHITECTURE.md** (15 min)
   - System design
   - Data flow
   - Performance characteristics

### Implementation Details

- `pkg/pdf/search.go` - Text search implementation
- `pkg/pdf/cache.go` - LRU cache with thread safety
- `pkg/ui/keybindings.go` - Vim keybinding handlers
- `pkg/config/theme.go` - Color schemes

---

## Common Development Tasks

### Add a New Keybinding

1. Edit `pkg/ui/keybindings.go`
2. Add case to `handleNormalKey()`
3. Define message type
4. Handle in `Model.Update()`

Example: Adding `H` for go to start of line

```go
// 1. In keybindings.go handleNormalKey
case "h":
    return GoStartOfLine

// 2. Define command
var GoStartOfLine = func() tea.Msg {
    return GoStartMsg{}
}

// 3. In model.go
case GoStartMsg:
    m.viewport.LineUp(m.viewport.Height)
    return nil
```

### Add a New Theme

1. Edit `pkg/config/theme.go`
2. Add new `Theme` variable
3. Update `GetTheme()` function

### Improve Performance

1. Profile with `make profile-cpu`
2. Identify bottleneck
3. Optimize and benchmark
4. Commit with benchmark results

---

## Debugging Tips

### Enable Logging

```go
import "log"

log.Printf("Debug: %v", value)
```

### Use Delve Debugger

```bash
dlv debug ./cmd/lumos -- ~/Documents/test.pdf

(dlv) help
(dlv) break main.main
(dlv) continue
(dlv) next
(dlv) print variable_name
```

### Check Memory Usage

```bash
# macOS
top -pid $(pgrep -f "./build/lumos")

# Linux
top -p $(pgrep -f "./build/lumos")
```

### Inspect PDF Structure

```go
// Add to main.go temporarily
doc, _ := pdf.NewDocument("test.pdf", 5)
log.Printf("Pages: %d\n", doc.GetPageCount())
meta := doc.GetMetadata()
log.Printf("Metadata: %+v\n", meta)
```

---

## Next Steps After Setup

### Phase 1 (MVP)

- [ ] Build and run locally
- [ ] Test with various PDFs
- [ ] Verify all vim keybindings work
- [ ] Check dark mode colors
- [ ] Profile startup time and memory
- [ ] Write unit tests
- [ ] Document findings

### Phase 2 (Enhancement)

- [ ] Add table of contents extraction
- [ ] Implement fuzzy search
- [ ] Add bookmark support
- [ ] Improve text layout

### Phase 3 (Images)

- [ ] Integrate go-termimg
- [ ] Add terminal protocol detection
- [ ] Render images in PDF
- [ ] Implement hybrid rendering

### Phase 4 (AI)

- [ ] Integrate Claude Agent SDK
- [ ] Add PDF Q&A
- [ ] Generate audio from PDF
- [ ] Multi-document analysis

---

## Troubleshooting

### Build Fails

```bash
# Clean and rebuild
make clean
make build

# Check Go version
go version  # Should be 1.21+

# Verify dependencies
go mod tidy
go mod verify
```

### Can't Load PDF

```bash
# Check file exists
ls -lh ~/Documents/your-file.pdf

# Check if valid PDF
file ~/Documents/your-file.pdf

# Run with error output
./build/lumos ~/Documents/your-file.pdf 2>&1
```

### Slow Performance

```bash
# Profile CPU
make profile-cpu

# Check cache hit rate
# (Add logging to cache.go)

# Profile memory
make profile-mem
```

---

## Resources

### Documentation
- `README.md` - Project overview
- `docs/ARCHITECTURE.md` - System design
- `docs/DEVELOPMENT.md` - Development guide
- `docs/KEYBINDINGS.md` - Keyboard reference (planned)

### External Resources
- [Go Tour](https://tour.golang.org/) - Learn Go
- [Bubble Tea Guide](https://github.com/charmbracelet/bubbletea/wiki)
- [ledongthuc/pdf](https://github.com/ledongthuc/pdf) - PDF library

### Related Projects
- [LUMINA CCN](../LUMINA/ccn) - Markdown viewer
- [Glow](https://github.com/charmbracelet/glow) - Markdown rendering
- [Helix](https://helix-editor.com/) - Modern editor (inspiration)

---

## Quick Reference: Vim Keybindings

```
Navigation:
  j/k, ↓/↑    - Scroll down/up
  d/u         - Half page down/up
  gg/G        - Top/bottom
  Ctrl+N/P    - Next/prev page

Search:
  /           - Start search
  n/N         - Next/prev match

UI:
  Tab         - Cycle panes
  1/2         - Dark/light mode
  ?           - Help
  q           - Quit
```

---

## Getting Started Right Now

```bash
# 1. Navigate to project
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMOS

# 2. Build
make build

# 3. Find a PDF
ls ~/Documents/*.pdf

# 4. Run LUMOS
./build/lumos ~/Documents/YOUR_PDF.pdf

# 5. Try keybindings
# Press: j k d u gg G q / n N ? Tab 1 2

# 6. Read help
./build/lumos --help
./build/lumos --keys
```

That's it! You're now running LUMOS.

---

**Happy developing! 🚀**

For questions, see:
- `docs/DEVELOPMENT.md` - Development guide
- `docs/ARCHITECTURE.md` - System design
- `README.md` - Project overview
