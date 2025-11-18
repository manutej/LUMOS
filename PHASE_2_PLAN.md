# LUMOS Phase 2: Enhanced Viewing Specification

**Status**: Planning
**Duration**: 3-4 weeks (4 milestones)
**Base**: Phase 1 complete, production-ready MVP
**Architecture**: Elm MVU pattern (proven in Phase 1)
**Priority**: Reuse LUMOS patterns, not external references

---

## Phase 2 Vision

Transform LUMOS from a basic reader into a **powerful document navigation and search tool** while maintaining the clean Elm architecture that makes Phase 1 solid.

### Target Capabilities
- Jump to specific pages via table of contents
- Regex-powered search with advanced filters
- Document bookmarks with persistence
- Better text layout preservation
- Configuration system for user preferences

---

## Milestone 2.1: Table of Contents Extraction (1 week)

### Goal
Extract PDF structure to create navigable table of contents and support `GoTo Page` command.

### User Stories
- Extract document outline/bookmarks from PDF metadata
- Generate simple TOC from text structure (headings)
- Jump to page via `:N` command (already implemented!)
- Display TOC in metadata pane

### Implementation

**pkg/pdf/toc.go** (New)
```go
// TableOfContents represents document structure
type TableOfContents struct {
    Entries []TOCEntry
}

type TOCEntry struct {
    Title    string
    Page     int
    Children []TOCEntry  // For nested structure
    Level    int         // 1=chapter, 2=section, etc.
}

// Extract from PDF metadata and structure
func (d *Document) ExtractTableOfContents() (*TableOfContents, error)
func (d *Document) FindHeadings(pageNum int) []TOCEntry
```

**pkg/ui/toc.go** (New)
```go
// TOCPane manages table of contents navigation
type TOCPane struct {
    toc          *pdf.TableOfContents
    selectedIdx  int
    viewport     viewport.Model
}

// Navigate via arrow keys, select with Enter
```

**Tests**: 15+ tests
- TOC extraction from sample PDFs
- Heading detection
- Jump-to-page functionality
- Edge cases (empty TOC, single page)

---

## Milestone 2.2: Advanced Search (1 week)

### Goal
Replace simple search with regex-powered search + filters.

### User Stories
- Search with regex patterns
- Case-sensitive toggle
- Whole word matching
- Match highlighting in context
- Search history

### Implementation

**pkg/pdf/search_advanced.go** (New)
```go
type SearchOptions struct {
    CaseSensitive bool
    UseRegex      bool
    WholeWord     bool
    MaxResults    int
}

type SearchEngine struct {
    doc     *Document
    options SearchOptions
}

func (s *SearchEngine) Search(pattern string) ([]SearchResult, error)
func (s *SearchEngine) SearchWithOptions(pattern string, opts SearchOptions) ([]SearchResult, error)
```

**pkg/ui/search.go** (Enhanced)
```go
type SearchPane struct {
    query       string
    results     []pdf.SearchResult
    currentIdx  int
    options     pdf.SearchOptions  // Toggle options
    history     []string           // Previous searches
}

// UI: Show search results in modal, navigate with n/N, toggle options with flags
```

**Tests**: 20+ tests
- Regex pattern matching
- Case sensitivity toggle
- Whole word matching
- Performance on large PDFs

---

## Milestone 2.3: Bookmarks & Persistent Configuration (1 week)

### Goal
Save user preferences and bookmarks persistently.

### User Stories
- Mark pages as bookmarks
- Jump between bookmarks
- Persist bookmarks to file
- Configuration file: `~/.config/lumos/config.toml`
- Save last read position

### Implementation

**pkg/config/config.go** (New/Enhanced)
```go
type Config struct {
    Theme           string              // Dark theme preference
    CacheSize       int                 // LRU cache size
    SearchOptions   pdf.SearchOptions   // Default search settings
    BookmarkedPages map[string][]int    // Per-document bookmarks
    LastPosition    map[string]LastPos  // Last page/scroll pos
}

type LastPos struct {
    Page      int
    ScrollPos int
    Timestamp time.Time
}

// Load/save from ~/.config/lumos/config.toml
func LoadConfig() (*Config, error)
func (c *Config) Save() error
```

**pkg/ui/bookmarks.go** (New)
```go
type BookmarkManager struct {
    bookmarks map[int]string    // pageNum -> note
    config    *config.Config
}

func (b *BookmarkManager) AddBookmark(page int, note string) error
func (b *BookmarkManager) RemoveBookmark(page int) error
func (b *BookmarkManager) ListBookmarks() []BookmarkEntry
func (b *BookmarkManager) Save() error
```

**Tests**: 15+ tests
- Config file reading/writing
- Bookmark persistence
- TOML parsing
- Edge cases

---

## Milestone 2.4: Layout & Text Preservation (1 week)

### Goal
Better preserve PDF structure during text extraction.

### User Stories
- Preserve line breaks and spacing
- Detect columns and multi-column layouts
- Preserve formatting (bold, italic hints)
- Better word wrapping in viewport

### Implementation

**pkg/pdf/layout.go** (New)
```go
type LayoutAnalyzer struct {
    pageWidth  float64
    pageHeight float64
}

type Column struct {
    Left   float64
    Right  float64
    Text   []string
}

// Analyze PDF structure to detect columns
func (l *LayoutAnalyzer) DetectColumns(page *pdfPage) []Column
func (l *LayoutAnalyzer) ExtractWithLayout(page *pdfPage) string
```

**Tests**: 15+ tests
- Column detection
- Line preservation
- Text wrapping
- Sample PDFs with various layouts

---

## Architecture Patterns (Proven in Phase 1)

### Model-View-Update (Elm Architecture)
- Same MVU pattern works for Phase 2 features
- Extend Model struct with new fields:
  ```go
  type Model struct {
      // ... Phase 1 fields

      // Phase 2
      tocPane          *TOCPane
      searchEngine     *pdf.SearchEngine
      bookmarks        *BookmarkManager
      config           *config.Config
  }
  ```

### Message Types (Extend as needed)
```go
// Add to messages.go
type TOCSelectedMsg struct {
    Entry *pdf.TOCEntry
}

type SearchOptionsChangedMsg struct {
    Options pdf.SearchOptions
}

type BookmarkAddedMsg struct {
    Page int
    Note string
}

type ConfigReloadedMsg struct {
    Config *config.Config
}
```

### Update Flow
```
Key Press
    ↓
keyHandler.HandleKey()
    ↓
Generate Message (TOCSelectedMsg, etc.)
    ↓
Model.Update()
    ↓
Update relevant sub-component (TOCPane, SearchEngine, etc.)
    ↓
View() re-renders
```

---

## UI Layout Enhancement

### Current (Phase 1)
```
┌─────────────┬──────────────────┬──────────────┐
│  Metadata   │    Viewer        │   Search     │
│  (20%)      │    (60%)         │   (20%)      │
└─────────────┴──────────────────┴──────────────┘
```

### Phase 2 Option A: Tab-based Panes
```
[Metadata] [Viewer] [Search] [TOC] [Bookmarks]

Active pane expands to fill space
```

### Phase 2 Option B: Modal for TOC
```
┌─────────────┬──────────────────┬──────────────┐
│             │    TOC Modal     │              │
│  Metadata   │  (overlay)       │   Search     │
│             │                  │              │
└─────────────┴──────────────────┴──────────────┘
```

---

## Testing Strategy

**Unit Tests**: 65+ tests across 4 milestones
- TOC extraction: 15 tests
- Advanced search: 20 tests
- Bookmarks/config: 15 tests
- Layout analysis: 15 tests

**Integration Tests**: 15+ tests
- Full search workflow
- TOC jump-to-page
- Config persistence
- Bookmark save/load

**Manual Testing**:
- Real PDF files with various structures
- Large documents (1000+ pages)
- Nested bookmarks
- Complex layouts

**Target Coverage**: >90% (maintain Phase 1's 94.4%)

---

## Known Dependencies

✅ Phase 1 Complete:
- Elm architecture working
- LRUCache for page caching
- Search engine (basic)
- Vim keybindings
- Dark theme system
- Test framework

---

## Phase 2 Timeline

**Week 1**: Milestone 2.1 (TOC extraction)
**Week 2**: Milestone 2.2 (Advanced search)
**Week 3**: Milestone 2.3 (Bookmarks & config)
**Week 4**: Milestone 2.4 (Layout preservation)

**Target Completion**: ~4 weeks from start

---

## Post-Phase 2

### Phase 3 (Q1 2026)
- Image rendering (go-termimg)
- Encrypted PDF support
- Better metadata extraction

### Phase 4 (Q2 2026)
- Claude Agent SDK integration
- AI-powered features
- PDF Q&A

---

## Success Criteria

### Functional
- [ ] Extract TOC from real PDFs
- [ ] Regex search with filters
- [ ] Bookmarks persist to disk
- [ ] Configuration file works
- [ ] Layout preserved in output

### Quality
- [ ] 65+ unit tests passing
- [ ] 15+ integration tests
- [ ] >90% coverage maintained
- [ ] Zero performance regressions
- [ ] All Phase 1 features still work

### Documentation
- [ ] Update PROGRESS.md
- [ ] Create Phase 2 milestone reviews
- [ ] Document new keybindings
- [ ] Update README

---

## Why This Approach

1. **Proven Pattern**: Elm MVU architecture already works perfectly
2. **Incremental**: Each milestone is independent, can ship separately
3. **User-Facing**: Every milestone adds visible value
4. **Low Risk**: Build on solid Phase 1 foundation
5. **Testable**: Clear interfaces, easy to test
6. **Maintainable**: Same code style and patterns as Phase 1

---

**Status**: Ready to implement
**Architecture**: Elm MVU (proven)
**Pattern**: Extend Phase 1 cleanly
**Quality**: Maintain 90%+ test coverage
