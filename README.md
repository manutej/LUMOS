# LUMOS - PDF Dark Mode Reader

**Project Status**: 🚀 Phase 0 Design Complete | Phase 1 Implementation Starting
**Date Created**: 2025-10-21
**Companion To**: [LUMINA](../LUMINA/) - Markdown Editor & Claude Code TUI

---

## 🎯 Vision

**LUMOS** is a developer-friendly PDF reader with a dark mode interface, inspired by modern markdown editors. It provides:

- 🌙 Beautiful dark mode rendering optimized for reading
- 📖 Smooth PDF scrolling and navigation
- 🎨 3-pane layout (consistent with LUMINA's ccn)
- 📸 Automatic image rendering with terminal graphics
- ⌨️ Vim-style keybindings
- 🔍 Full text search and document indexing
- 🤖 **Claude Agent SDK wrapper** (Phase 2+) for NotebookLM-like features

### Key Difference from LUMINA

While **LUMINA** (ccn) is a markdown file navigator with Glamour rendering, **LUMOS** extends that vision to **PDFs** with:

1. **PDF-Specific Challenges**:
   - Binary parsing instead of text files
   - Page-based navigation vs file-based
   - Image rendering in terminal
   - Text extraction and layout preservation

2. **Unique Opportunities**:
   - AI-powered PDF analysis (Claude Agent SDK)
   - Audio generation from PDF content
   - Interactive annotations
   - Cross-document relationship mapping

---

## 📋 Project Phases

### ✅ Phase 0: Research & Design (COMPLETE)

**Deliverables**:
- PDF library research (3 comprehensive options)
- Terminal graphics protocol analysis
- Architecture design document
- Performance benchmarks
- Technology stack decision

**Files**:
- `LUMOS_PDF_LIBRARY_RESEARCH.md` (30KB)
- `LUMOS_QUICK_REFERENCE.md` (9KB)
- `LUMOS_ARCHITECTURE_EXAMPLE.md` (22KB)
- `LUMOS_RESEARCH_INDEX.md` (12KB)

**Decision**: ledongthuc/pdf + go-termimg + Bubble Tea

---

### 🔄 Phase 1: MVP - Basic PDF Reader (2-3 weeks)

**Goal**: Production-ready dark mode PDF viewer with navigation

**Tasks**:
- [ ] Initialize Go project with dependencies
- [ ] Create model structure for PDF document handling
- [ ] Implement basic PDF rendering (text extraction)
- [ ] Build 3-pane layout with Bubble Tea
  - Left: Document metadata & page list
  - Center: PDF content viewer
  - Right: Text search preview
- [ ] Implement vim keybindings
  - `j/k` - Scroll down/up
  - `d/u` - Half page down/up
  - `gg/G` - Top/bottom of document
  - `/` - Search
  - `n/N` - Next/previous match
- [ ] Add dark mode by default with color scheme
- [ ] Page navigation (Ctrl+P: previous, Ctrl+N: next)
- [ ] Test with diverse PDFs

**Success Criteria**:
- [ ] Cold start <100ms
- [ ] Page navigation <50ms
- [ ] Memory usage <50MB for typical PDFs
- [ ] Smooth scrolling at 60fps
- [ ] All vim keybindings working

**Deliverables**:
- Working CLI: `lumos /path/to/file.pdf`
- 3-pane layout rendering
- Basic search functionality
- Vim keybindings
- Dark mode UI

---

### 🔄 Phase 2: Enhanced Viewing (1-2 weeks)

**Goal**: Improve PDF viewing experience with better text handling

**Tasks**:
- [ ] Fuzzy search with ripgrep integration
- [ ] Text extraction improvements
- [ ] Better layout preservation
- [ ] Bookmark support (vim marks)
- [ ] PDF metadata display
- [ ] Table of contents extraction and navigation
- [ ] Advanced vim commands (marks, registers)

**Deliverables**:
- Fast full-text search
- TOC navigation
- Bookmarks/marks system
- Better text formatting

---

### 🔄 Phase 3: Image Support (1-2 weeks)

**Goal**: Render complex PDFs with images and diagrams

**Tasks**:
- [ ] Integrate image rendering (go-termimg)
- [ ] Terminal protocol auto-detection (Kitty/iTerm2/SIXEL)
- [ ] Hybrid rendering (text + images)
- [ ] Image caching (LRU)
- [ ] Performance optimization

**Deliverables**:
- Image rendering in terminal
- Automatic protocol detection
- Fast image-heavy PDF support
- <200ms page switch for images

---

### 🔄 Phase 4: AI Integration (2-3 weeks)

**Goal**: Add Claude Agent SDK wrapper for NotebookLM-like features

**Tasks**:
- [ ] Claude Agent SDK integration (Go or via API)
- [ ] PDF content extraction and chunking
- [ ] Audio generation from PDF text
- [ ] Interactive summary generation
- [ ] Multi-document analysis
- [ ] Chat interface for PDF Q&A

**Deliverables**:
- `/ask` command for PDF Q&A
- Audio generation from pages
- Summary generation
- Cross-document analysis

---

## 🛠️ Technology Stack

### Core Dependencies

```go
// TUI Framework & Components
github.com/charmbracelet/bubbletea       // MVU framework (26k⭐)
github.com/charmbracelet/bubbles         // UI components (5.5k⭐)
github.com/charmbracelet/lipgloss        // Terminal styling (8k⭐)
github.com/charmbracelet/glamour         // Markdown rendering (2.5k⭐)

// PDF Processing
github.com/ledongthuc/pdf                // PDF parsing (1.5k⭐)
github.com/pdfcpu/pdfcpu                 // PDF manipulation (2k⭐)

// Terminal Graphics
github.com/blacktop/go-termimg           // Auto-detect terminal protocols
```

### Architecture Comparison

```
LUMINA (ccn)              LUMOS
┌──────────────┐          ┌──────────────┐
│ Markdown     │          │ PDF          │
│ Navigator    │          │ Reader       │
├──────────────┤          ├──────────────┤
│ Text Files   │          │ Binary PDFs   │
│ Glamour      │          │ ledongthuc/pdf
│ Rendering    │          │ + go-termimg │
├──────────────┤          ├──────────────┤
│ Bubble Tea   │◄────────►│ Bubble Tea   │
│ 3-pane       │ Shared   │ 3-pane       │
│ Pattern      │ Patterns │ Pattern      │
└──────────────┘          └──────────────┘
```

### Performance Targets

| Metric | Target | Status |
|--------|--------|--------|
| Cold start | <100ms | 🎯 |
| Page switch (cached) | <50ms | 🎯 |
| Page switch (uncached) | <200ms | 🎯 |
| Memory (10MB PDF) | <50MB | 🎯 |
| Text search | <100ms | 🎯 |
| Image render | <300ms | 🎯 |

---

## 📁 Project Structure

```
LUMOS/
├── README.md                          # This file
├── go.mod                             # Go module definition
├── go.sum                             # Dependency lock
├── Makefile                           # Build & test automation
├── .gitignore                         # Git ignore rules
│
├── cmd/
│   └── lumos/
│       └── main.go                    # CLI entrypoint
│
├── pkg/
│   ├── pdf/
│   │   ├── document.go                # PDF document model
│   │   ├── renderer.go                # Text rendering
│   │   ├── cache.go                   # Page caching (LRU)
│   │   └── search.go                  # Full-text search
│   │
│   ├── ui/
│   │   ├── model.go                   # Bubble Tea model
│   │   ├── view.go                    # Rendering
│   │   ├── keybindings.go             # Vim keybindings
│   │   ├── panes.go                   # 3-pane layout
│   │   └── styles.go                  # Dark mode styles
│   │
│   ├── terminal/
│   │   ├── graphics.go                # go-termimg integration
│   │   ├── protocols.go               # Kitty/iTerm2/SIXEL
│   │   └── colors.go                  # Terminal color management
│   │
│   └── config/
│       ├── defaults.go                # Default configuration
│       ├── loader.go                  # Config file loading
│       └── theme.go                   # Theme management
│
├── test/
│   ├── fixtures/                      # Test PDF files
│   │   ├── simple.pdf
│   │   ├── images.pdf
│   │   ├── tables.pdf
│   │   └── large.pdf
│   │
│   └── benchmarks/
│       ├── rendering_bench_test.go
│       └── search_bench_test.go
│
├── docs/
│   ├── ARCHITECTURE.md                # Detailed architecture
│   ├── KEYBINDINGS.md                 # Vim keybinding reference
│   ├── DEVELOPMENT.md                 # Development guide
│   ├── TESTING.md                     # Testing strategy
│   └── PERFORMANCE.md                 # Performance notes
│
├── scripts/
│   ├── build.sh                       # Build script
│   ├── test.sh                        # Test runner
│   └── benchmark.sh                   # Benchmark runner
│
└── examples/
    ├── simple-reader.go               # Minimal example
    └── advanced-usage.go              # Full-featured example
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- macOS / Linux / WSL

### Build from Source

```bash
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMOS

# Install dependencies
go mod download

# Build
go build -o lumos ./cmd/lumos

# Run
./lumos /path/to/file.pdf
```

### Global Installation (Planned)

```bash
# Will be similar to LUMINA
ln -s ~/LUMOS/lumos ~/bin/lumos

# Then use anywhere
lumos ~/Documents/paper.pdf
```

---

## ⌨️ Keyboard Controls (Phase 1 MVP)

### General Navigation
- `q` or `Ctrl+C` - Quit
- `Tab` - Cycle panes (Metadata → Viewer → Search)
- `?` - Show help

### PDF Navigation
- `j` / `k` or `↓` / `↑` - Scroll down/up one line
- `d` / `u` - Half page down/up
- `gg` - Go to first page
- `G` - Go to last page
- `Ctrl+N` / `Ctrl+P` - Next/previous page
- `N` [page]` - Jump to page N

### Search
- `/` - Start search
- `n` - Next match
- `N` - Previous match
- `Esc` - Exit search

### UI
- `1` - Toggle dark mode (Phase 1)
- `2` - Toggle light mode (Phase 1)
- `:` - Command mode (Phase 2+)

---

## 🎨 Dark Mode Design

### Color Scheme (Default)

```
Background:     #1e1e1e (near black)
Text:           #e0e0e0 (light gray)
Accent:         #61afef (blue, active selection)
Warning:        #e06c75 (red, for errors)
Success:        #98c379 (green, for search matches)
```

### Reference
- Inspired by VSCode Dark+ theme
- Optimized for long reading sessions
- Adjustable in config (Phase 2+)

---

## 📊 Development Roadmap

```
Phase 0: Design ✅
    └─→ Research, architecture, tech stack

Phase 1: MVP 🚀 (current)
    └─→ Basic PDF reader, dark mode, vim keys

Phase 2: Enhanced 📖 (after Phase 1)
    └─→ Search, TOC, bookmarks

Phase 3: Images 🖼️ (after Phase 2)
    └─→ Terminal image rendering

Phase 4: AI 🤖 (after Phase 3)
    └─→ Claude Agent SDK integration
```

---

## 🔗 Integration with LUMINA

LUMOS and LUMINA (ccn) form a complementary pair:

| Feature | LUMINA | LUMOS |
|---------|--------|-------|
| Purpose | Markdown navigation | PDF reading |
| Input | Markdown files | PDF documents |
| Rendering | Glamour | PDF + Images |
| Use Case | Documentation | Reference materials |
| Keybindings | Vim-style | Vim-style |
| Framework | Bubble Tea | Bubble Tea |

**Unified Workflow**:
```
Workflow: Technical Reading + Documentation

1. Use LUMINA to navigate markdown docs
   → lumina ~/docs/guide.md

2. Reference PDF technical spec
   → lumos ~/docs/spec.pdf

3. Both tools use consistent vim keybindings
   → No context switching

4. Eventually: One unified interface?
   → Phase 4+ future: Integration layer
```

---

## 📚 Documentation

All planning documents are in parent `/LUXOR/` directory:

- **LUMOS_QUICK_REFERENCE.md** - Technology decisions and comparison tables
- **LUMOS_PDF_LIBRARY_RESEARCH.md** - Comprehensive library analysis
- **LUMOS_ARCHITECTURE_EXAMPLE.md** - Code examples and patterns
- **LUMOS_RESEARCH_INDEX.md** - Navigation guide

---

## 🧪 Testing Strategy

**Phase 1 Focus**: Basic functionality and correctness

```
Unit Tests:
  ├── PDF parsing
  ├── Text extraction
  ├── Page caching
  └── Keybinding routing

Integration Tests:
  ├── Layout rendering
  ├── Navigation flow
  └── Search functionality

Performance Tests:
  ├── <100ms cold start
  ├── <50ms page switch
  └── <50MB memory usage

Manual Testing:
  ├── Various PDF types
  ├── Large documents (100+ pages)
  └── Terminal compatibility
```

---

## 🎯 Success Criteria

### Phase 1 (MVP)
- [ ] Reads and displays PDF content
- [ ] Dark mode by default
- [ ] All vim keybindings working
- [ ] <100ms startup
- [ ] <50MB memory for typical PDFs
- [ ] Can scroll through 100-page documents smoothly
- [ ] Text search works
- [ ] 50+ test cases passing
- [ ] Documentation complete

### Overall (Future)
- [ ] Companion tool to LUMINA
- [ ] AI-powered PDF analysis
- [ ] Audio generation from PDFs
- [ ] 10k+ GitHub stars (aspirational)
- [ ] Used by developers daily

---

## 📝 Development Guidelines

### Code Style
- Follow Go idioms and best practices
- Use clear variable names
- Comment public functions and types
- Organize code by concern (pdf, ui, config, etc.)

### Testing
- Aim for 80%+ coverage
- Include integration tests
- Performance benchmarks for critical paths
- Test with diverse PDF types

### Performance
- Profile before optimizing
- LRU cache for pages
- Lazy load PDF content
- Stream text rendering

### Documentation
- Clear README for each package
- Code comments for complex logic
- Architecture docs in /docs/
- Examples in /examples/

---

## 🤝 Contributing

To extend LUMOS:

1. Follow the project structure
2. Write tests first (TDD)
3. Update documentation
4. Run benchmarks
5. Ensure <100ms startup target

---

## 📜 License

TBD - Following LUMINA (likely MIT or Apache 2.0)

---

## 📖 References

- **PDF Library**: [ledongthuc/pdf](https://github.com/ledongthuc/pdf)
- **Terminal Graphics**: [go-termimg](https://github.com/blacktop/go-termimg)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Similar Tools**:
  - [Glow](https://github.com/charmbracelet/glow) - Markdown viewer
  - [Pdfless](https://github.com/Aaaaatyle/pdfless) - PDF reader (Rust)

---

## 📍 Project Location

```
/Users/manu/Documents/LUXOR/PROJECTS/LUMOS
```

---

**Status**: Ready for Phase 1 Implementation
**Last Updated**: 2025-10-21
**Next**: Create Go project structure and begin Phase 1 development
