# LUMOS Architecture Documentation

**Version**: 0.1.0
**Date**: 2025-10-21
**Status**: Phase 1 MVP Design

---

## System Architecture

### High-Level Design

```
┌─────────────────────────────────────────────────────────────┐
│                    CLI Entry Point                          │
│                   (cmd/lumos/main.go)                      │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│              PDF Document Handler                           │
│            (pkg/pdf/document.go)                           │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ • Load & parse PDF files                            │   │
│  │ • Extract text content per page                     │   │
│  │ • LRU cache (5 pages default)                       │   │
│  │ • Full-text search support                          │   │
│  └─────────────────────────────────────────────────────┘   │
└──────────────────────────┬──────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│   Search    │  │    Cache     │  │   Render    │
│  (search.go) │  │  (cache.go)  │  │ (view TBD)  │
└──────────────┘  └──────────────┘  └──────────────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                 UI Layer (Bubble Tea)                       │
│                  (pkg/ui/model.go)                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Model: Manages application state                    │   │
│  │ • Document reference                                │   │
│  │ • Current page                                       │   │
│  │ • Theme (dark/light)                                │   │
│  │ • Search state                                       │   │
│  │ • Active pane                                        │   │
│  │                                                      │   │
│  │ Update: Handles messages and state changes          │   │
│  │ View: Renders 3-pane layout                         │   │
│  └─────────────────────────────────────────────────────┘   │
└──────────────────────────┬──────────────────────────────────┘
        │                  │                  │
        ▼                  ▼                  ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ Keybindings  │  │   Styling    │  │   Config    │
│(keybindings.)│  │ (theme.go)   │  │(theme.go)   │
└──────────────┘  └──────────────┘  └──────────────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                   Terminal Output                           │
│              (Glamour + Lip Gloss Rendering)              │
└─────────────────────────────────────────────────────────────┘
```

---

## Package Organization

### pkg/pdf/

**Purpose**: PDF document handling and manipulation

**Files**:
- `document.go` - Main PDF document type and page extraction
- `search.go` - Full-text search functionality
- `cache.go` - LRU caching for pages

**Key Types**:
```go
Document       // Represents a loaded PDF
PageInfo       // Information about a single page
SearchResult   // Result from text search
LRUCache       // Page caching mechanism
```

**Performance**:
- Page extraction: ~50ms (depends on page complexity)
- Cache hit: <1ms
- Cache miss: ~50ms
- Text search: <100ms for typical documents

### pkg/ui/

**Purpose**: Terminal UI and user interaction

**Files**:
- `model.go` - Main Bubble Tea model (MVU pattern)
- `keybindings.go` - Vim keybinding handlers
- `panes.go` (planned) - Individual pane rendering
- `styles.go` (planned) - Shared style utilities

**Key Types**:
```go
Model       // Main application state
KeyHandler  // Processes keyboard input
Styles      // UI styling configuration
```

**Rendering Layers**:
1. Model: Manages application state
2. Update: Processes messages and state changes
3. View: Renders the 3-pane layout

### pkg/config/

**Purpose**: Configuration and theming

**Files**:
- `theme.go` - Color themes and styles

**Themes**:
- Dark (default, VSCode Dark+)
- Light (alternative)

### pkg/terminal/ (Phase 3+)

**Purpose**: Terminal graphics and image rendering

**Files**:
- `graphics.go` - go-termimg integration
- `protocols.go` - Kitty/iTerm2/SIXEL support
- `colors.go` - Color management

---

## Data Flow

### Page Loading

```
User Action: Page Navigation
    │
    ▼
KeyPress (Ctrl+N for next page)
    │
    ▼
Model.Update() receives NavigateMsg
    │
    ▼
goToNextPage() called
    │
    ▼
LoadPageCmd() returns async command
    │
    ▼
PDF Extraction Happens Asynchronously
    │
    ▼
Check LRU Cache
    ├─ Hit: Return cached text (~1ms)
    └─ Miss: Extract from PDF (~50ms)
    │
    ▼
PageLoadedMsg sent back to Update()
    │
    ▼
viewport.SetContent() called
    │
    ▼
View() re-renders with new content
    │
    ▼
Terminal Updated
```

### Search Flow

```
User Action: Type /
    │
    ▼
KeyMode set to Search
    │
    ▼
User types query
    │
    ▼
User presses Enter
    │
    ▼
Document.Search(query)
    │
    ▼
Iterate all pages (extract if needed)
    │
    ▼
Find matches (case-insensitive)
    │
    ▼
Store SearchResult objects
    │
    ▼
Jump to first match page
    │
    ▼
Highlight matches in search pane
```

---

## State Management

### Model State

```go
type Model struct {
    // Document state
    document   *pdf.Document
    cache      *pdf.LRUCache
    currentPage int

    // UI state
    theme      config.Theme
    styles     config.Styles
    showHelp   bool

    // Search state
    searchActive   bool
    searchQuery    string
    searchResults  []pdf.SearchResult
    currentMatch   int

    // Viewport state
    viewport   viewport.Model

    // Dimensions
    width, height int
}
```

### Message Types

```go
// Navigation
NavigateMsg{Type: "first_page|last_page|next_page|prev_page"}

// Scrolling
ScrollMsg{Direction: "up|down", Amount: int}

// Search
SearchMsg{Direction: "next|prev", Query: string}

// UI
ThemeChangeMsg{Theme: "dark|light"}
ModeChangeMsg{Mode: KeyMode}
PaneChangeMsg{Direction: "next|prev"}
```

---

## Concurrency Model

### Thread Safety

1. **LRU Cache**:
   - Uses sync.RWMutex for thread-safe operations
   - Multiple goroutines can read cached pages
   - Write operations are serialized

2. **PDF Document**:
   - Uses sync.RWMutex for cache access
   - Page extraction can happen in background
   - ledongthuc/pdf handles internal concurrency

3. **UI Model**:
   - Single-threaded (Bubble Tea handles updates)
   - All state changes go through Update()
   - No race conditions expected

### Async Operations

- Page loading uses Bubble Tea async patterns
- Commands execute asynchronously
- Results come back as messages

---

## Performance Characteristics

### Cold Start

```
Time Breakdown:
1. Parse flags:        ~1ms
2. Open PDF:           ~5ms
3. Create model:       ~2ms
4. Initialize viewport: ~2ms
5. Load first page:    ~50ms (first time)
6. Render TUI:         ~10ms
────────────────────────────
TOTAL:                 ~70ms ✅ (target: <100ms)
```

### Page Navigation

**Cached (within LRU 5-page window)**:
- LRU lookup: <1ms
- Set viewport: <5ms
- Render: <10ms
- **Total**: <20ms ✅ (well under 50ms target)

**Uncached**:
- Load page: ~50ms
- Set viewport: ~5ms
- Render: ~10ms
- **Total**: ~65ms ✅ (under 200ms target)

### Memory Usage

For a typical 100-page PDF (~10MB file):
- Document structure: ~5MB
- LRU cache (5 pages): ~2MB
- UI components: ~1MB
- **Total**: ~8MB ✅ (target: <50MB)

---

## Error Handling

### PDF Loading Errors

```
File not found → User message + exit
Corrupted PDF → Extract as much as possible + warning
No pages → Error message + exit
Unreadable page → Skip + continue
```

### Runtime Errors

```
Out of memory → Graceful shutdown + message
Viewport resize → Adjust panes dynamically
Keyboard interrupt → Save state + exit
```

---

## Testing Strategy

### Unit Tests

**PDF Package**:
```
TestDocumentOpen
TestPageExtraction
TestCaching
TestSearch
```

**UI Package**:
```
TestKeyHandling
TestModelUpdates
TestNavigation
TestThemeChange
```

### Integration Tests

```
TestFullWorkflow (load → navigate → search)
TestLargeDocuments (100+ pages)
TestMemoryLeaks (extended use)
```

### Performance Tests

```
BenchmarkPageLoad
BenchmarkSearch
BenchmarkCacheHitRate
BenchmarkMemoryUsage
```

---

## Future Enhancements

### Phase 2 (Enhanced Viewing)
- Better text layout preservation
- Table of Contents extraction
- Bookmark support (vim marks)
- Advanced search (regex, whole-word)

### Phase 3 (Image Support)
- Terminal image rendering (go-termimg)
- Auto-detect terminal protocol
- Image caching (separate from text)
- Hybrid rendering (text + images)

### Phase 4 (AI Integration)
- Claude Agent SDK integration
- Audio generation from PDF
- Interactive Q&A
- Multi-document analysis

---

## Design Decisions

### Why Text-First Rendering?

1. **Performance**: Text extraction and display is very fast
2. **Accessibility**: Text can be searched, copied, manipulated
3. **Memory**: Text uses far less memory than images
4. **Terminal**: Text rendering works in all terminal emulators

### Why LRU Cache?

1. **Locality**: Users typically read sequentially (nearby pages)
2. **Memory**: 5 pages is a good balance (enough for context, not too much memory)
3. **Simplicity**: Easy to implement and understand

### Why Bubble Tea + Glamour?

1. **Consistency**: Same stack as LUMINA (ccn)
2. **Maturity**: 26k stars, proven in production
3. **Features**: Great components (viewport, list)
4. **Styling**: Lip Gloss makes theming trivial

---

## Limitations & Trade-offs

### Current Limitations

1. **Text-only extraction**: Complex layouts may not render perfectly
2. **No TLS support**: Can't handle encrypted PDFs (Phase 2+)
3. **Memory usage**: Large PDFs with many high-resolution images problematic
4. **Search speed**: Full-text search goes through all pages

### Accepted Trade-offs

1. **Simplicity > Features**: Minimal viable product before adding complexity
2. **Performance > Perfection**: 80% correctness with great performance better than 99% with poor UX
3. **Standard > Custom**: Use proven libraries over building from scratch

---

## References

- **ledongthuc/pdf**: https://github.com/ledongthuc/pdf
- **Bubble Tea**: https://github.com/charmbracelet/bubbletea
- **go-termimg**: https://github.com/blacktop/go-termimg (Phase 3+)
- **LUMINA CCN**: /Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn

---

**Last Updated**: 2025-10-21
