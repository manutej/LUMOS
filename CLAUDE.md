# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**LUMOS** is a dark mode PDF reader for developers built with Go and Bubble Tea TUI framework. It features vim-style keybindings, 3-pane layout, full-text search, and LRU page caching. A companion to LUMINA (markdown viewer).

**Current Phase**: Phase 1 MVP (Milestone 1.6 - Dark Mode Polish)
**Progress**: 83% complete (5/6 milestones done)
**Test Coverage**: 94.4% (PDF), 47.8% (UI)
**Status**: Functionally complete, pending test coverage improvement

## Quick References

- **Dependencies & Skills**: See `.claude/DEPENDENCIES.md` for:
  - Complete dependency breakdown
  - Relevant global skills (golang-backend-development, frontend-architecture, testing)
  - Development workflow with dependencies
  - Troubleshooting guide

- **Specifications**: See `.specify/` for:
  - Complete Phase 1 specification
  - Milestone requirements and acceptance criteria
  - Architecture decisions and testing strategy

## Common Commands

### Build & Run
```bash
# Build the binary
make build

# Build and run with test PDF
make run

# Install globally to ~/bin/
make install

# Show help/version/keybindings
./build/lumos --help
./build/lumos --version
./build/lumos --keys

# Open a PDF
./build/lumos ~/Documents/sample.pdf
```

### Testing
```bash
# Run all tests (42 tests, all passing)
make test

# Run tests with verbose output
make test-v

# Run tests with race detector
make test-race

# Generate coverage report (current: 94.4%)
make coverage
```

### Quality & Performance
```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet

# Run all CI checks
make ci-check

# Run benchmarks
make bench

# Profile CPU/memory
make profile-cpu
make profile-mem
```

### Development
```bash
# Setup development tools (air, golangci-lint)
make dev-setup

# Watch for changes and rebuild
make watch

# Clean build artifacts
make clean
```

## Architecture

### High-Level Structure

```
┌─────────────────┐
│  CLI Entry      │  cmd/lumos/main.go
│  (tea.Program)  │  - Parse flags
└────────┬────────┘  - Load PDF
         │           - Run TUI
         ▼
┌─────────────────┐
│  PDF Package    │  pkg/pdf/
│  (Document)     │  - document.go: PDF loading/parsing
│                 │  - search.go: Full-text search
│                 │  - cache.go: LRU page caching
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  UI Package     │  pkg/ui/
│  (Bubble Tea)   │  - model.go: MVU pattern state
│                 │  - keybindings.go: Vim keybindings
│                 │  - messages.go: Tea messages
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Config Package │  pkg/config/
│  (Themes)       │  - theme.go: Dark/light themes
└─────────────────┘
```

### Package Organization

**cmd/lumos/**
- `main.go` - CLI entry point with flag parsing and tea.Program initialization

**pkg/pdf/**
- `document.go` - Core PDF document handling with ledongthuc/pdf library
- `search.go` - Full-text search with case-insensitive matching and context extraction
- `cache.go` - LRU cache for page content (5 pages default)
- Performance: Cache hits <100ns, search <50μs

**pkg/ui/**
- `model.go` - Bubble Tea MVU model (Model-View-Update pattern)
- `keybindings.go` - Vim-style key handling (j/k/d/u/gg/G/Ctrl+N/Ctrl+P)
- `messages.go` - Custom Bubble Tea messages
- 3-pane layout: Metadata (20%) | Viewer (60%) | Search (20%)

**pkg/config/**
- `theme.go` - Color themes (Dark default, Light optional)

### Data Flow

**Page Loading**:
```
User Action (Ctrl+N) → KeyPress → Update() → LoadPageCmd
  → Check LRU Cache → (Hit: <1ms) | (Miss: ~50ms PDF extraction)
  → PageLoadedMsg → viewport.SetContent() → View() re-renders
```

**Search**:
```
User types '/' → Search mode → Type query → Enter
  → Document.Search(query) → Extract all pages → Find matches
  → Store results → Jump to first match → Highlight in search pane
```

### State Management

The Model struct holds all application state:
- **document**: *pdf.Document reference
- **cache**: *pdf.LRUCache (5 pages)
- **currentPage**: int (1-indexed)
- **theme**: Dark/Light
- **searchActive**: bool
- **searchResults**: []pdf.SearchResult
- **viewport**: Bubble Tea viewport component

All state changes flow through `Update(msg tea.Msg)` - no shared state mutations.

## Testing Strategy

### Test Structure

```
test/
├── fixtures/           # Generated PDF test files
│   ├── simple.pdf      # 1 page, basic operations
│   ├── multipage.pdf   # 5 pages, navigation tests
│   └── search_test.pdf # Search pattern tests
├── generate_fixtures.py # Python script to generate PDFs
├── TESTING_GUIDE.md    # Comprehensive testing guide
└── README.md           # Test documentation
```

### Running Specific Tests

```bash
# Run only cache tests
go test -v ./pkg/pdf -run TestCache

# Run only search tests
go test -v ./pkg/pdf -run TestSearch

# Run specific test
go test -v ./pkg/pdf -run TestCacheHit

# Run with coverage for single package
go test -coverprofile=cover.out ./pkg/pdf
go tool cover -html=cover.out
```

### Regenerating Test Fixtures

```bash
cd test
python3 generate_fixtures.py
# Creates simple.pdf, multipage.pdf, search_test.pdf
```

## Key Implementation Details

### PDF Library (ledongthuc/pdf)

The project uses `github.com/ledongthuc/pdf` for parsing. Key patterns:

```go
// Open PDF (must close file handle)
f, r, err := pdf.Open(filepath)
if err != nil { return err }
defer f.Close()

// Get page count
pages := r.NumPage()

// Extract text from page (1-indexed)
page := r.Page(pageNum)
texts := page.Content().Text
for _, text := range texts {
    content += text.S + " "
}
```

**Important**: Each page extraction requires reopening the file (ledongthuc/pdf limitation).

### LRU Cache Implementation

5-page LRU cache with thread-safe access:
- Uses `sync.RWMutex` for concurrent reads
- Simple map-based implementation (no external dependencies)
- Cache hit: ~16ns/op, Cache miss: ~8ns/op

### Search Performance

Case-insensitive full-text search:
- `TextToLines()` - Split text into line objects
- `CaseInsensitiveMatch()` - Find all positions of query
- `ExtractContext()` - Get context before/after match
- Performance: ~10μs/op for 1KB text

### Vim Keybindings

Implemented in `pkg/ui/keybindings.go`:
- j/k: Scroll down/up one line
- d/u: Half page down/up
- gg/G: First/last page
- Ctrl+N/Ctrl+P: Next/previous page
- /: Search mode
- n/N: Next/previous match
- Tab: Cycle panes
- 1/2: Dark/light mode toggle
- q: Quit

## Development Workflow

### Adding New Features

1. Write tests first (TDD approach used throughout project)
2. Implement in appropriate package (pdf/ui/config)
3. Update documentation (PROGRESS.md, milestone reviews)
4. Run full test suite and verify coverage
5. Update README if user-facing changes

### Code Style

- Follow Go idioms and conventions
- Use clear variable names (no abbreviations)
- Comment public functions and types
- Organize by package concern (pdf, ui, config)
- Keep functions small and focused

### Performance Targets

| Metric | Target | Current |
|--------|--------|---------|
| Cold start | <100ms | ~70ms ✅ |
| Page switch (cached) | <50ms | <20ms ✅ |
| Page switch (uncached) | <200ms | ~65ms ✅ |
| Memory (10MB PDF) | <50MB | ~8MB ✅ |
| Text search | <100ms | <50μs ✅ |

## Current Development Status

**Phase 1 Progress**: 50% (3/6 milestones complete)

### Completed Milestones
- ✅ 1.1: Build & Compile - Clean build, dependencies resolved
- ✅ 1.2: Core Testing - 42 tests, 9 benchmarks, 70% coverage
- ✅ 1.3: Test Fixtures - 3 PDF fixtures, 94.4% coverage

### Next Milestone
- ⏳ 1.4: Basic TUI Framework - Implement Bubble Tea UI components

### Upcoming
- 1.5: Vim Keybindings - Complete vim-style navigation
- 1.6: Dark Mode Polish - Theme refinement and final touches

## Important Files

### Documentation
- `START_HERE.md` - Quick orientation guide
- `PROGRESS.md` - Current development status (update after each milestone)
- `PHASE_1_PLAN.md` - Complete Phase 1 breakdown
- `PHASE_1_MILESTONE_X_X_REVIEW.md` - Milestone completion reviews
- `test/TESTING_GUIDE.md` - Comprehensive testing guide
- `docs/ARCHITECTURE.md` - Detailed architecture document

### Planning Documents (in project root)
- Multiple PHASE_1_* files tracking development progress
- Review documents for completed milestones
- Guides for upcoming milestones

## Known Limitations

1. **Text-only extraction**: Complex PDF layouts may not render perfectly
2. **Character spacing**: reportlab-generated PDFs have spacing quirks (use real PDFs for production testing)
3. **No image support**: Phase 3 feature (go-termimg integration)
4. **No encrypted PDFs**: Phase 2+ feature
5. **File reopening**: ledongthuc/pdf requires reopening file for each page extraction

## Dependencies

Core dependencies (go.mod):
- `github.com/charmbracelet/bubbletea` - TUI framework (MVU pattern)
- `github.com/charmbracelet/bubbles` - UI components (viewport)
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/ledongthuc/pdf` - PDF parsing

## Phase Roadmap

- **Phase 0**: ✅ Design & Research (complete)
- **Phase 1**: 🚧 MVP - Basic PDF Reader (50% complete)
- **Phase 2**: 📅 Enhanced Viewing (TOC, bookmarks, better search)
- **Phase 3**: 📅 Image Support (go-termimg integration)
- **Phase 4**: 📅 AI Integration (Claude Agent SDK, NotebookLM-like features)

## Companion Project

**LUMINA** (`/Users/manu/Documents/LUXOR/PROJECTS/LUMINA`)
- Markdown viewer with similar TUI patterns
- Shares Bubble Tea, Glamour, Lipgloss stack
- Reference for UI component patterns

## Additional Context

When working on LUMOS:
1. Update `PROGRESS.md` after significant changes
2. Create milestone review documents upon completion
3. Maintain 80%+ test coverage (currently 94.4%)
4. Follow established patterns from pkg/pdf and pkg/ui
5. Performance benchmarks validate architecture decisions
6. Refer to `test/TESTING_GUIDE.md` for testing best practices
