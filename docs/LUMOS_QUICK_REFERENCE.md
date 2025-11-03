# LUMOS: Quick Reference & Decision Matrix

**Last Updated**: 2025-10-21

---

## TL;DR Recommendation

**For LUMOS MVP**: Use **ledongthuc/pdf** + **go-termimg** + **Bubble Tea**

**Architecture**: Text-first rendering with image fallback
**Performance**: <100ms cold start, <50ms page switch
**Memory**: <50MB for typical docs

---

## PDF Library Comparison

| Feature | pdfcpu | UniPDF | ledongthuc/pdf |
|---------|--------|--------|----------------|
| **License** | Apache 2.0 | Commercial | BSD-3 |
| **Cost** | Free | $$ | Free |
| **Text Extraction** | Basic | Excellent | Good |
| **Layout Preservation** | No | Yes | No |
| **Table Detection** | No | Yes | No |
| **Performance** | Excellent | Good | Good |
| **Memory** | Excellent | Good | Excellent |
| **Pure Go** | ✅ | ✅ | ✅ |
| **Active Development** | ✅ | ✅ | ✅ |
| **Best For** | Images, metadata | Professional text | MVP, OSS |
| **MVP Ready** | 🟡 | ✅ | ✅ |

### Verdict

- **MVP/OSS**: ledongthuc/pdf
- **Commercial**: UniPDF
- **Image ops**: pdfcpu

---

## Terminal Graphics Protocol Comparison

| Protocol | Speed | Quality | Compatibility | Recommendation |
|----------|-------|---------|---------------|----------------|
| **Kitty** | ⚡⚡⚡ | ⚡⚡⚡ | Modern terminals | Primary |
| **iTerm2** | ⚡⚡⚡ | ⚡⚡⚡ | macOS, some Linux | Secondary |
| **SIXEL** | 🟡 | 🟡 | Widest (legacy) | Fallback |

### Terminal Support Matrix

| Terminal | Kitty | iTerm2 | SIXEL | Recommendation |
|----------|-------|--------|-------|----------------|
| Kitty | ✅ | ✅ | ✅ | All protocols |
| iTerm2 | ❌ | ✅ | ✅ | iTerm2 or SIXEL |
| WezTerm | ✅ | ✅ | ✅ | All protocols |
| Ghostty | ✅ | ❌ | ❌ | Kitty only |
| Alacritty | ❌ | ❌ | ❌ | Text only |
| xterm | ❌ | ❌ | ✅ | SIXEL only |

---

## Go Libraries for Terminal Graphics

| Library | Protocols | Auto-Detect | API | Recommendation |
|---------|-----------|-------------|-----|----------------|
| **go-termimg** | Kitty, iTerm2 | ✅ | Simple | Primary |
| **rasterm** | Kitty, iTerm2, SIXEL | ❌ | Explicit | Secondary |

### Code Example Comparison

**go-termimg** (simpler):
```go
import "github.com/blacktop/go-termimg"

// Automatically detects supported protocol
termimg.Render(img, os.Stdout)
```

**rasterm** (more control):
```go
import "github.com/BourgeoisBear/rasterm"

// Explicit protocol selection
rasterm.Encode(img, rasterm.ProtocolKitty)
```

---

## Rendering Strategy Comparison

| Strategy | Speed | Memory | Quality | Complexity | Use Case |
|----------|-------|--------|---------|------------|----------|
| **Text-only** | ⚡⚡⚡ | ⚡⚡⚡ | 🟢 | 🟢 | 80% cases |
| **Hybrid** | ⚡⚡ | ⚡⚡ | ⚡⚡ | 🟡 | 15% cases |
| **Image-only** | 🟡 | 🟡 | ⚡⚡⚡ | 🟡 | 5% cases |

### Decision Flow

```
PDF Page
    │
    ├─ Extract text (ledongthuc/pdf)
    │       │
    │       ├─ Success + >100 chars? → TEXT RENDERING (90%)
    │       │
    │       └─ Failure / <100 chars
    │               │
    │               ├─ Has images? → HYBRID RENDERING (8%)
    │               │
    │               └─ Complex layout → IMAGE RENDERING (2%)
```

---

## Performance Targets

### Conservative (MVP)

| Operation | Target | Library Choice |
|-----------|--------|----------------|
| Cold start | <100ms | ledongthuc/pdf |
| Page switch (cached) | <50ms | LRU cache |
| Page switch (uncached) | <200ms | ledongthuc/pdf |
| Text search | <100ms | ledongthuc/pdf |
| Memory (50 pages) | <50MB | Cache 5 pages |

### Optimistic (Production)

| Operation | Target | Library Choice |
|-----------|--------|----------------|
| Cold start | <50ms | UniPDF |
| Page switch (cached) | <20ms | Optimized cache |
| Page switch (uncached) | <100ms | UniPDF |
| Text search | <50ms | Pre-indexed |
| Memory (50 pages) | <30MB | Smart caching |

---

## Technology Stack

### Recommended

```toml
[core]
pdf = "github.com/ledongthuc/pdf"               # Text extraction
image = "github.com/pdfcpu/pdfcpu"              # Image fallback

[rendering]
graphics = "github.com/blacktop/go-termimg"    # Auto-detect protocol
ansi = "github.com/charmbracelet/x/ansi"       # Text formatting

[tui]
framework = "github.com/charmbracelet/bubbletea"
components = "github.com/charmbracelet/bubbles"
styling = "github.com/charmbracelet/lipgloss"

[utilities]
cache = "github.com/hashicorp/golang-lru"      # Page cache
```

### Alternative (Commercial)

```toml
[core]
pdf = "github.com/unidoc/unipdf/v3"            # Better text extraction
# Rest same as recommended
```

---

## Implementation Checklist

### Phase 1: MVP (1-2 weeks)
- [ ] Set up Bubble Tea scaffolding
- [ ] Integrate ledongthuc/pdf
- [ ] Basic navigation (arrows, PgUp/PgDn)
- [ ] Display current page number
- [ ] Simple ANSI formatting (bold headers)

### Phase 2: Enhanced Text (1 week)
- [ ] Better text formatting (code blocks, lists)
- [ ] Extract table of contents
- [ ] Display metadata (title, author, pages)
- [ ] Basic search functionality
- [ ] Clipboard integration

### Phase 3: Image Support (1-2 weeks)
- [ ] Detect text extraction failures
- [ ] Integrate pdfcpu for page→image
- [ ] Integrate go-termimg for rendering
- [ ] Terminal protocol detection
- [ ] Fallback chain (Kitty→iTerm2→text)

### Phase 4: Performance (1 week)
- [ ] Implement LRU cache (5 pages)
- [ ] Async page pre-loading
- [ ] Progress indicators
- [ ] Configuration file
- [ ] Custom keybindings

---

## Risk Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| Poor text extraction | High | Tier 1→2→3 fallback strategy |
| Terminal incompatibility | Medium | Protocol detection + text fallback |
| Memory usage | Medium | LRU cache + lazy loading |
| Large PDF performance | Low | Async rendering + progress |
| Licensing (UniPDF) | Low | Start with ledongthuc, upgrade later |

---

## Code Templates

### Document Manager (Simplified)

```go
type PDFManager struct {
    path      string
    reader    *pdf.Reader
    pageCount int
    cache     *lru.Cache
}

func (m *PDFManager) GetPage(num int) (string, error) {
    // Check cache
    if cached, ok := m.cache.Get(num); ok {
        return cached.(string), nil
    }

    // Render page
    page := m.reader.Page(num)
    text, err := page.GetPlainText(nil)
    if err != nil {
        return "", err
    }

    // Format and cache
    formatted := formatForTerminal(text)
    m.cache.Add(num, formatted)
    return formatted, nil
}
```

### Bubble Tea Model

```go
type Model struct {
    pdfManager *PDFManager
    viewport   viewport.Model
    currentPage int
    totalPages  int
    err        error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "j", "down":
            m.currentPage = min(m.currentPage+1, m.totalPages)
            return m, m.loadPage(m.currentPage)
        case "k", "up":
            m.currentPage = max(m.currentPage-1, 1)
            return m, m.loadPage(m.currentPage)
        }
    }
    return m, nil
}
```

---

## Comparison with Existing Tools

| Feature | LUMOS | pdfless | zathura | termpdf | Advantage |
|---------|-------|---------|---------|---------|-----------|
| **Pure Terminal** | ✅ | ✅ | ❌ | ✅ | No X11 needed |
| **Text Search** | ✅ | ✅ | ✅ | ✅ | Standard feature |
| **Image Support** | ✅ | ❌ | ✅ | ✅ | Modern terminals |
| **Pure Go** | ✅ | ❌ | ❌ | ❌ | Easy deployment |
| **Bubble Tea** | ✅ | ❌ | ❌ | ❌ | Better UX |
| **LUMINA Integration** | ✅ | ❌ | ❌ | ❌ | Workflow fit |

---

## Success Metrics

### Functional
- ✅ Render 80% of PDFs with text extraction
- ✅ Graceful image fallback for remaining 20%
- ✅ Search works on text-rendered pages
- ✅ Table of contents extraction (when available)

### Performance
- ✅ <100ms cold start
- ✅ <50ms cached page switch
- ✅ <200ms uncached page switch
- ✅ <50MB memory for 50-page doc

### UX
- ✅ Familiar vim-style keybindings
- ✅ Page number indicator
- ✅ Progress for slow operations
- ✅ Help screen (?)

### Integration
- ✅ Works on same terminals as LUMINA ccn
- ✅ Shared keybinding patterns
- ✅ Consistent CLI experience

---

## FAQ

### Why not use pdftotext CLI?
- External dependency
- No control over formatting
- Harder to integrate search/navigation

### Why not pure image rendering?
- Slower (500ms vs 50ms)
- Higher memory usage
- No text search
- Not accessible

### Why Bubble Tea instead of Ratatui?
- LUMINA uses Go + Bubble Tea
- Consistent ecosystem
- Easier code sharing

### Why not UniPDF from start?
- Licensing complexity
- ledongthuc sufficient for MVP
- Can upgrade later if needed

### What about encrypted PDFs?
- ledongthuc has weak encryption support
- UniPDF handles better
- Consider separate decrypt step

---

## Next Actions

1. **Immediate** (Today):
   - Clone ledongthuc/pdf examples
   - Test with sample PDFs
   - Measure parse times

2. **This Week**:
   - Bubble Tea scaffolding
   - Basic page navigation
   - Simple text rendering

3. **Next Week**:
   - ANSI formatting
   - Search implementation
   - Table of contents

4. **Week 3-4**:
   - Image fallback
   - LRU cache
   - Performance tuning

---

**Status**: Ready to implement
**Confidence**: High (95%)
**Risk Level**: Low
