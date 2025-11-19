# LUMOS Dependencies Reference

**Project**: LUMOS - Dark Mode PDF Reader
**Language**: Go 1.24.1
**Last Updated**: 2025-11-18

---

## Overview

LUMOS is built with modern Go libraries for TUI (Terminal User Interface) development and PDF processing. This document maps dependencies to their purposes and references relevant skills from the global Claude Code configuration.

---

## Core Dependencies (Direct)

### 1. Bubble Tea Framework

**Package**: `github.com/charmbracelet/bubbletea v1.1.0`

**Purpose**: Core TUI framework implementing the Elm Architecture (Model-View-Update pattern)

**Usage in LUMOS**:
- `pkg/ui/model.go` - MVU pattern implementation
- `cmd/lumos/main.go` - Program initialization

**Key Concepts**:
```go
// Model-View-Update (MVU) pattern
type Model struct { ... }                    // State
func (m Model) Init() tea.Cmd { ... }       // Initialize
func (m Model) Update(msg) (Model, Cmd)     // Update state
func (m Model) View() string { ... }        // Render UI
```

**Relevant Global Skills**:
- **golang-backend-development** - Core Go patterns
- **frontend-architecture** - MVU pattern principles (applies to TUI)

**Documentation**:
- Official: https://github.com/charmbracelet/bubbletea
- Examples: https://github.com/charmbracelet/bubbletea/tree/master/examples

---

### 2. Bubbles Components

**Package**: `github.com/charmbracelet/bubbles v0.20.0`

**Purpose**: Reusable TUI components (viewport, textinput, spinner, etc.)

**Usage in LUMOS**:
- `pkg/ui/model.go:32-34` - Viewport for scrollable PDF content
- Three viewport instances: main viewer, metadata pane, search pane

**Key Component**:
```go
import "github.com/charmbracelet/bubbles/viewport"

vp := viewport.New(width, height)
vp.SetContent(content)
vp.LineDown(1)  // Scroll functionality
```

**Relevant Global Skills**:
- **frontend-architecture** - Component composition patterns
- **react-development** - Similar component model (conceptually)

**Documentation**:
- Official: https://github.com/charmbracelet/bubbles
- Viewport: https://github.com/charmbracelet/bubbles/tree/master/viewport

---

### 3. Lipgloss Styling

**Package**: `github.com/charmbracelet/lipgloss v0.13.1`

**Purpose**: Terminal styling and layout library (think CSS for TUI)

**Usage in LUMOS**:
- `pkg/config/theme.go` - Color themes and styles
- `pkg/ui/model.go:121-136` - Layout composition (JoinHorizontal, JoinVertical)

**Key Patterns**:
```go
import "github.com/charmbracelet/lipgloss"

// Style definition (like CSS)
style := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("62")).
    Padding(1).
    Width(width)

// Layout composition
content := lipgloss.JoinHorizontal(
    lipgloss.Top,
    metadataPane,
    viewerPane,
    searchPane,
)
```

**Relevant Global Skills**:
- **frontend-architecture** - Layout and styling patterns
- **mobile-design** - Responsive design principles (applies to terminal)

**Documentation**:
- Official: https://github.com/charmbracelet/lipgloss
- Examples: https://github.com/charmbracelet/lipgloss/tree/master/examples

---

### 4. PDF Library

**Package**: `github.com/ledongthuc/pdf v0.0.0-20250511090121-5959a4027728`

**Purpose**: Pure Go PDF parsing library (no CGO dependencies)

**Usage in LUMOS**:
- `pkg/pdf/document.go` - PDF loading and page extraction
- `pkg/pdf/search.go` - Text extraction for search

**Key Patterns**:
```go
import "github.com/ledongthuc/pdf"

// Open PDF
file, reader, err := pdf.Open(filepath)
defer file.Close()

// Extract text from page (1-indexed)
page := reader.Page(pageNum)
texts := page.Content().Text
for _, text := range texts {
    content += text.S + " "
}
```

**Known Limitations**:
- Must reopen file for each page extraction
- Text-only extraction (no images in this library)
- Character spacing quirks with some PDF generators

**Mitigation**: LRU cache minimizes file reopening impact

**Relevant Global Skills**:
- **golang-backend-development** - File I/O patterns

**Documentation**:
- GitHub: https://github.com/ledongthuc/pdf

---

## Indirect Dependencies (Transitive)

### Terminal & System

**Packages**:
- `github.com/mattn/go-isatty v0.0.20` - TTY detection
- `github.com/mattn/go-localereader v0.0.1` - Locale handling
- `golang.org/x/sys v0.24.0` - OS-specific syscalls

**Purpose**: Cross-platform terminal interaction

---

### Text Rendering

**Packages**:
- `github.com/mattn/go-runewidth v0.0.16` - Unicode width calculation
- `github.com/rivo/uniseg v0.4.7` - Unicode grapheme segmentation
- `golang.org/x/text v0.17.0` - Text processing

**Purpose**: Correct text rendering in terminal (CJK support, emoji, etc.)

---

### Color & Styling

**Packages**:
- `github.com/lucasb-eyer/go-colorful v1.2.0` - Color manipulation
- `github.com/muesli/termenv v0.15.2` - Terminal environment detection
- `github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6` - ANSI escape codes
- `github.com/aymanbagabas/go-osc52/v2 v2.0.1` - OSC 52 clipboard support

**Purpose**: Rich terminal styling and color support

---

### Input Handling

**Packages**:
- `github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f` - Windows console input
- `github.com/muesli/cancelreader v0.2.2` - Cancellable reader

**Purpose**: Cross-platform keyboard/mouse input

---

### Concurrency

**Packages**:
- `golang.org/x/sync v0.8.0` - Synchronization primitives

**Purpose**: Thread-safe operations (used in LRU cache)

---

## Global Skills Integration

### Recommended Skills for LUMOS Development

From `~/.claude/skills/`, these are most relevant:

#### 1. **golang-backend-development** ✅
**Location**: `~/.claude/skills/golang-backend-development/`

**Covers**:
- Go idioms and conventions
- Concurrency patterns (goroutines, channels)
- Error handling
- Testing strategies
- Package organization

**LUMOS Usage**:
- All of `pkg/pdf/` - backend data processing
- Concurrency in LRU cache
- Error handling patterns throughout

**Invoke**: Automatically available, or explicitly reference with `@golang-backend-development`

---

#### 2. **frontend-architecture** (TUI = Frontend!) ✅
**Location**: `~/.claude/skills/frontend-architecture/`

**Covers**:
- Component architecture
- State management
- Event handling
- Layout patterns
- Separation of concerns

**LUMOS Usage**:
- `pkg/ui/` - TUI is frontend development
- MVU pattern = Component architecture
- Message passing = Event system
- 3-pane layout = Responsive design

**Invoke**: Automatically available

---

#### 3. **testing** ✅
**Location**: Referenced in global skills

**Covers**:
- Unit testing
- Table-driven tests
- Mocking and fixtures
- Benchmarking
- Coverage targets

**LUMOS Usage**:
- `pkg/pdf/*_test.go` - 42 tests, 94.4% coverage
- `pkg/ui/*_test.go` - 9 tests, 47.8% coverage (needs improvement)
- `test/fixtures/` - PDF test fixtures

**Invoke**: Automatically available

---

### Skills NOT Needed for LUMOS

These global skills are **not applicable** to LUMOS:

❌ **Web Framework Skills** (Express, FastAPI, etc.) - LUMOS is a CLI tool
❌ **Database Skills** (PostgreSQL, MongoDB) - No database in LUMOS
❌ **API Skills** (REST, GraphQL) - No API server
❌ **Container Skills** (Docker, Kubernetes) - Not containerized yet (Phase 2+)
❌ **Cloud Skills** (AWS, Azure) - Local-only application

---

## Development Workflow with Dependencies

### Installing Dependencies

```bash
# Download all dependencies
go mod download

# Verify dependencies
go mod verify

# Update dependencies (careful!)
go get -u ./...
go mod tidy
```

### Checking Dependency Graph

```bash
# See all dependencies
go list -m all

# See dependency tree
go mod graph

# Why is package X included?
go mod why -m github.com/charmbracelet/bubbletea
```

### Vendoring (Optional)

```bash
# Create vendor/ directory
go mod vendor

# Build using vendor/
go build -mod=vendor ./...
```

---

## Dependency Version Management

### Current Strategy

**Pinned versions** in go.mod for stability:
- ✅ Bubble Tea: v1.1.0 (stable)
- ✅ Bubbles: v0.20.0 (stable)
- ✅ Lipgloss: v0.13.1 (stable)
- ✅ PDF library: Pinned commit (no semver)

### Update Policy

**DO UPDATE**:
- Security patches
- Bug fixes in PDF library
- New Bubble Tea features (after testing)

**DON'T UPDATE** without testing:
- Major version bumps
- API-breaking changes
- Experimental features

### Testing Updates

```bash
# Test with new version
go get github.com/charmbracelet/bubbletea@latest
go mod tidy
make test
make build

# If tests pass, commit go.mod and go.sum
git add go.mod go.sum
git commit -m "deps: update bubbletea to vX.Y.Z"
```

---

## Troubleshooting Dependencies

### "Missing sum" errors

```bash
go mod tidy
go mod verify
```

### "Ambiguous import" errors

Check for duplicate dependencies:
```bash
go list -m all | grep <package-name>
```

### "Module not found" errors

```bash
# Clear cache and re-download
go clean -modcache
go mod download
```

### CGO errors

LUMOS uses **pure Go dependencies only** - no CGO required!

If you see CGO errors, check:
```bash
go env CGO_ENABLED  # Should work with 0 or 1
```

---

## Adding New Dependencies

### Before Adding

**Ask**:
1. Is it pure Go? (No CGO preferred)
2. Is it actively maintained?
3. Does it have good test coverage?
4. Is there a lighter alternative?
5. What's the license? (MIT/Apache preferred)

### Adding Process

```bash
# 1. Add dependency
go get github.com/author/package@v1.2.3

# 2. Import in code
# (import in your .go files)

# 3. Tidy up
go mod tidy

# 4. Test
make test

# 5. Document here
# (add to this file)

# 6. Commit
git add go.mod go.sum
git commit -m "deps: add github.com/author/package for X feature"
```

---

## Performance Considerations

### Dependency Impact

**Bubble Tea Ecosystem** (bubbletea + bubbles + lipgloss):
- Binary size: ~2.5MB
- Memory overhead: ~1MB
- Performance: Excellent (60 FPS rendering)

**PDF Library** (ledongthuc/pdf):
- Binary size: ~1.5MB
- Memory overhead: ~2MB per open document
- Performance: 50ms per page extraction (mitigated by cache)

**Total Binary**: 4.6MB (very good for Go TUI app)

### Optimization Tips

1. **Use LRU cache** - Already implemented for PDF pages
2. **Limit viewport content** - Don't render off-screen content
3. **Batch updates** - Use tea.Batch for multiple commands
4. **Profile before optimizing**:
   ```bash
   make profile-cpu
   make profile-mem
   ```

---

## Security Considerations

### Dependency Scanning

```bash
# Check for known vulnerabilities
go list -json -m all | nancy sleuth

# Or use govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

### Supply Chain Security

**Checksum verification**: go.sum ensures integrity
**Vendor audit**: All dependencies are open-source and auditable
**No network calls**: LUMOS is 100% offline after build

---

## Resources

### Official Documentation
- **Go Modules**: https://go.dev/ref/mod
- **Bubble Tea**: https://github.com/charmbracelet/bubbletea
- **Bubbles**: https://github.com/charmbracelet/bubbles
- **Lipgloss**: https://github.com/charmbracelet/lipgloss

### Learning Resources
- **Charm Tutorials**: https://charm.sh/
- **TUI Development Guide**: https://github.com/charmbracelet/bubbletea/tree/master/tutorials
- **Go Testing**: https://go.dev/doc/tutorial/add-a-test

### Community
- **Charm Slack**: https://charm.sh/slack
- **r/golang**: https://reddit.com/r/golang
- **Gophers Slack**: https://gophers.slack.com

---

**Last Updated**: 2025-11-18
**Maintained By**: Claude Code
**Version**: Phase 1 MVP (v0.1.0)
