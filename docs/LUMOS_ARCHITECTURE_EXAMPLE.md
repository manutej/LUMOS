# LUMOS: Architecture & Code Examples

**Purpose**: Practical implementation guide with code examples
**Date**: 2025-10-21

---

## Project Structure

```
lumos/
├── cmd/
│   └── lumos/
│       └── main.go                 # Entry point
│
├── internal/
│   ├── document/
│   │   ├── manager.go             # PDF document manager
│   │   ├── cache.go               # LRU page cache
│   │   ├── metadata.go            # Metadata extraction
│   │   └── types.go               # Common types
│   │
│   ├── rendering/
│   │   ├── text.go                # Text extraction/formatting
│   │   ├── image.go               # Image rendering
│   │   ├── hybrid.go              # Mixed rendering
│   │   ├── ansi.go                # ANSI styling
│   │   └── detect.go              # Capability detection
│   │
│   ├── ui/
│   │   ├── model.go               # Bubble Tea model
│   │   ├── view.go                # View rendering
│   │   ├── update.go              # Update logic
│   │   ├── keybindings.go         # Keyboard handling
│   │   └── styles.go              # Lipgloss styles
│   │
│   └── search/
│       ├── index.go               # Text search index
│       └── query.go               # Search queries
│
├── pkg/
│   └── terminal/
│       ├── graphics.go            # Terminal graphics support
│       └── detect.go              # Terminal detection
│
├── go.mod
├── go.sum
└── README.md
```

---

## Core Types

### internal/document/types.go

```go
package document

import (
    "image"
)

// PageContent represents a rendered PDF page
type PageContent struct {
    PageNum     int
    Content     string        // Rendered text or image reference
    RenderMode  RenderMode
    HasImages   bool
    WordCount   int
    RenderTime  time.Duration
}

type RenderMode int

const (
    RenderModeText RenderMode = iota
    RenderModeHybrid
    RenderModeImage
)

// Metadata contains PDF document information
type Metadata struct {
    Title      string
    Author     string
    Subject    string
    Keywords   string
    PageCount  int
    FileSize   int64
    PDFVersion string
}

// CacheStats tracks cache performance
type CacheStats struct {
    Hits       int64
    Misses     int64
    Evictions  int64
    MemoryUsed int64
}
```

---

## Document Manager

### internal/document/manager.go

```go
package document

import (
    "fmt"
    "os"
    "sync"

    lru "github.com/hashicorp/golang-lru"
    "github.com/ledongthuc/pdf"
)

const (
    DefaultCacheSize = 5
    MinTextThreshold = 100
)

type Manager struct {
    path      string
    reader    *pdf.Reader
    file      *os.File
    metadata  Metadata
    cache     *lru.Cache
    stats     CacheStats
    mu        sync.RWMutex
}

func NewManager(path string) (*Manager, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open PDF: %w", err)
    }

    stat, err := f.Stat()
    if err != nil {
        f.Close()
        return nil, fmt.Errorf("failed to stat file: %w", err)
    }

    reader, err := pdf.NewReader(f, stat.Size())
    if err != nil {
        f.Close()
        return nil, fmt.Errorf("failed to create PDF reader: %w", err)
    }

    cache, err := lru.New(DefaultCacheSize)
    if err != nil {
        f.Close()
        return nil, fmt.Errorf("failed to create cache: %w", err)
    }

    m := &Manager{
        path:   path,
        reader: reader,
        file:   f,
        cache:  cache,
        metadata: Metadata{
            PageCount: reader.NumPage(),
            FileSize:  stat.Size(),
        },
    }

    // Extract metadata
    if err := m.extractMetadata(); err != nil {
        // Non-fatal, continue without metadata
        _ = err
    }

    return m, nil
}

func (m *Manager) Close() error {
    return m.file.Close()
}

func (m *Manager) GetPage(pageNum int) (*PageContent, error) {
    if pageNum < 1 || pageNum > m.metadata.PageCount {
        return nil, fmt.Errorf("page %d out of range (1-%d)", pageNum, m.metadata.PageCount)
    }

    // Check cache
    if cached, ok := m.cache.Get(pageNum); ok {
        m.mu.Lock()
        m.stats.Hits++
        m.mu.Unlock()
        return cached.(*PageContent), nil
    }

    m.mu.Lock()
    m.stats.Misses++
    m.mu.Unlock()

    // Render page
    start := time.Now()
    content, err := m.renderPage(pageNum)
    if err != nil {
        return nil, err
    }
    content.RenderTime = time.Since(start)

    // Cache result
    m.cache.Add(pageNum, content)

    return content, nil
}

func (m *Manager) renderPage(pageNum int) (*PageContent, error) {
    page := m.reader.Page(pageNum)
    if page.V.IsNull() {
        return nil, fmt.Errorf("page %d not found", pageNum)
    }

    // Try text extraction
    text, err := page.GetPlainText(nil)
    if err == nil && len(text) >= MinTextThreshold {
        return &PageContent{
            PageNum:    pageNum,
            Content:    text,
            RenderMode: RenderModeText,
            HasImages:  false,
            WordCount:  estimateWordCount(text),
        }, nil
    }

    // TODO: Fallback to image rendering
    return nil, fmt.Errorf("text extraction failed, image rendering not yet implemented")
}

func (m *Manager) GetMetadata() Metadata {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.metadata
}

func (m *Manager) GetStats() CacheStats {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.stats
}

func (m *Manager) extractMetadata() error {
    // PDF metadata extraction from Info dictionary
    // This is simplified - real implementation would parse Info dict
    m.metadata.Title = "Untitled"
    m.metadata.Author = "Unknown"
    return nil
}

func estimateWordCount(text string) int {
    // Rough estimate: count spaces + 1
    count := 1
    for _, r := range text {
        if r == ' ' || r == '\n' || r == '\t' {
            count++
        }
    }
    return count
}

// Preload pages in background
func (m *Manager) PreloadPages(pages []int) {
    go func() {
        for _, pageNum := range pages {
            if _, ok := m.cache.Get(pageNum); !ok {
                _, _ = m.GetPage(pageNum)
            }
        }
    }()
}
```

---

## Text Rendering

### internal/rendering/text.go

```go
package rendering

import (
    "fmt"
    "strings"

    "github.com/charmbracelet/lipgloss"
)

var (
    headerStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("12")).
        MarginBottom(1)

    paragraphStyle = lipgloss.NewStyle().
        MarginBottom(1)

    codeStyle = lipgloss.NewStyle().
        Background(lipgloss.Color("236")).
        Foreground(lipgloss.Color("220")).
        Padding(0, 1)

    listStyle = lipgloss.NewStyle().
        MarginLeft(2)
)

type TextRenderer struct {
    width  int
    height int
}

func NewTextRenderer(width, height int) *TextRenderer {
    return &TextRenderer{
        width:  width,
        height: height,
    }
}

func (r *TextRenderer) Render(text string) string {
    var rendered strings.Builder

    // Split into paragraphs
    paragraphs := strings.Split(text, "\n\n")

    for _, para := range paragraphs {
        para = strings.TrimSpace(para)
        if para == "" {
            continue
        }

        // Detect paragraph type and apply styling
        switch {
        case isHeader(para):
            rendered.WriteString(r.renderHeader(para))
        case isCodeBlock(para):
            rendered.WriteString(r.renderCodeBlock(para))
        case isList(para):
            rendered.WriteString(r.renderList(para))
        default:
            rendered.WriteString(r.renderParagraph(para))
        }

        rendered.WriteString("\n")
    }

    return rendered.String()
}

func (r *TextRenderer) renderHeader(text string) string {
    // Remove common header markers
    text = strings.TrimLeft(text, "#- ")
    text = strings.ToUpper(text)
    return headerStyle.Render(text) + "\n"
}

func (r *TextRenderer) renderParagraph(text string) string {
    // Wrap text to terminal width
    wrapped := wrapText(text, r.width-4)
    return paragraphStyle.Render(wrapped) + "\n"
}

func (r *TextRenderer) renderCodeBlock(text string) string {
    lines := strings.Split(text, "\n")
    var rendered strings.Builder

    for _, line := range lines {
        rendered.WriteString(codeStyle.Render(line))
        rendered.WriteString("\n")
    }

    return rendered.String()
}

func (r *TextRenderer) renderList(text string) string {
    lines := strings.Split(text, "\n")
    var rendered strings.Builder

    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }

        // Add bullet if not present
        if !strings.HasPrefix(line, "•") && !strings.HasPrefix(line, "-") {
            line = "• " + line
        }

        rendered.WriteString(listStyle.Render(line))
        rendered.WriteString("\n")
    }

    return rendered.String()
}

// Helper functions

func isHeader(text string) bool {
    // Simple heuristic: short, no punctuation at end, possibly uppercase
    if len(text) > 80 {
        return false
    }

    // Check for all caps or title case
    upper := strings.ToUpper(text)
    if text == upper {
        return true
    }

    // Check for markdown-style headers
    if strings.HasPrefix(text, "#") {
        return true
    }

    // No period at end, relatively short
    return !strings.HasSuffix(text, ".") && len(text) < 60
}

func isCodeBlock(text string) bool {
    // Heuristics for code:
    // - Contains programming symbols
    // - Indented
    // - Contains common keywords

    codeIndicators := []string{
        "func ", "def ", "class ", "import ", "return",
        "if (", "for (", "while (", "{", "}", "[]", "=>",
    }

    for _, indicator := range codeIndicators {
        if strings.Contains(text, indicator) {
            return true
        }
    }

    // Check for consistent indentation
    lines := strings.Split(text, "\n")
    indented := 0
    for _, line := range lines {
        if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
            indented++
        }
    }

    return indented > len(lines)/2
}

func isList(text string) bool {
    lines := strings.Split(text, "\n")
    bullets := 0

    for _, line := range lines {
        line = strings.TrimSpace(line)
        if strings.HasPrefix(line, "•") ||
            strings.HasPrefix(line, "-") ||
            strings.HasPrefix(line, "*") ||
            (len(line) > 2 && line[0] >= '0' && line[0] <= '9' && line[1] == '.') {
            bullets++
        }
    }

    return bullets > len(lines)/2
}

func wrapText(text string, width int) string {
    if width <= 0 {
        width = 80
    }

    var wrapped strings.Builder
    words := strings.Fields(text)
    lineLen := 0

    for i, word := range words {
        wordLen := len(word)

        if lineLen+wordLen+1 > width {
            wrapped.WriteString("\n")
            lineLen = 0
        }

        if i > 0 && lineLen > 0 {
            wrapped.WriteString(" ")
            lineLen++
        }

        wrapped.WriteString(word)
        lineLen += wordLen
    }

    return wrapped.String()
}
```

---

## Bubble Tea UI

### internal/ui/model.go

```go
package ui

import (
    "fmt"
    "time"

    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"

    "lumos/internal/document"
    "lumos/internal/rendering"
)

type Model struct {
    docManager  *document.Manager
    renderer    *rendering.TextRenderer
    viewport    viewport.Model
    ready       bool

    currentPage int
    totalPages  int
    metadata    document.Metadata

    loading     bool
    err         error

    width       int
    height      int
}

type pageLoadedMsg struct {
    content *document.PageContent
    err     error
}

func NewModel(path string) (*Model, error) {
    mgr, err := document.NewManager(path)
    if err != nil {
        return nil, err
    }

    meta := mgr.GetMetadata()

    return &Model{
        docManager:  mgr,
        currentPage: 1,
        totalPages:  meta.PageCount,
        metadata:    meta,
    }, nil
}

func (m Model) Init() tea.Cmd {
    return m.loadPage(m.currentPage)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit

        case "j", "down":
            if !m.loading && m.currentPage < m.totalPages {
                m.currentPage++
                return m, m.loadPage(m.currentPage)
            }

        case "k", "up":
            if !m.loading && m.currentPage > 1 {
                m.currentPage--
                return m, m.loadPage(m.currentPage)
            }

        case "g", "home":
            if !m.loading && m.currentPage != 1 {
                m.currentPage = 1
                return m, m.loadPage(1)
            }

        case "G", "end":
            if !m.loading && m.currentPage != m.totalPages {
                m.currentPage = m.totalPages
                return m, m.loadPage(m.totalPages)
            }

        case " ", "pgdown":
            m.viewport, cmd = m.viewport.Update(msg)
            cmds = append(cmds, cmd)
        }

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

        if !m.ready {
            m.viewport = viewport.New(msg.Width, msg.Height-3)
            m.viewport.YPosition = 2
            m.renderer = rendering.NewTextRenderer(msg.Width, msg.Height-3)
            m.ready = true
            cmds = append(cmds, m.loadPage(m.currentPage))
        } else {
            m.viewport.Width = msg.Width
            m.viewport.Height = msg.Height - 3
            m.renderer = rendering.NewTextRenderer(msg.Width, msg.Height-3)
        }

    case pageLoadedMsg:
        m.loading = false
        if msg.err != nil {
            m.err = msg.err
        } else {
            content := m.renderer.Render(msg.content.Content)
            m.viewport.SetContent(content)
            m.viewport.GotoTop()
        }
    }

    m.viewport, cmd = m.viewport.Update(msg)
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}

func (m Model) View() string {
    if !m.ready {
        return "Loading..."
    }

    // Header
    header := m.renderHeader()

    // Footer
    footer := m.renderFooter()

    // Main content
    return lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        m.viewport.View(),
        footer,
    )
}

func (m Model) renderHeader() string {
    style := lipgloss.NewStyle().
        Bold(true).
        Background(lipgloss.Color("62")).
        Foreground(lipgloss.Color("230")).
        Padding(0, 1).
        Width(m.width)

    title := m.metadata.Title
    if title == "" {
        title = "PDF Viewer"
    }

    return style.Render(title)
}

func (m Model) renderFooter() string {
    leftStyle := lipgloss.NewStyle().
        Background(lipgloss.Color("236")).
        Foreground(lipgloss.Color("250")).
        Padding(0, 1)

    rightStyle := lipgloss.NewStyle().
        Background(lipgloss.Color("236")).
        Foreground(lipgloss.Color("250")).
        Padding(0, 1)

    left := fmt.Sprintf("Page %d/%d", m.currentPage, m.totalPages)
    right := "q: quit | j/k: next/prev page | g/G: first/last"

    leftRendered := leftStyle.Render(left)
    rightRendered := rightStyle.Render(right)

    gap := m.width - lipgloss.Width(leftRendered) - lipgloss.Width(rightRendered)
    if gap < 0 {
        gap = 0
    }

    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        leftRendered,
        strings.Repeat(" ", gap),
        rightRendered,
    )
}

func (m *Model) loadPage(pageNum int) tea.Cmd {
    m.loading = true
    return func() tea.Msg {
        content, err := m.docManager.GetPage(pageNum)
        return pageLoadedMsg{content: content, err: err}
    }
}

func (m Model) Close() {
    if m.docManager != nil {
        m.docManager.Close()
    }
}
```

---

## Main Entry Point

### cmd/lumos/main.go

```go
package main

import (
    "fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"

    "lumos/internal/ui"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: lumos <pdf-file>")
        os.Exit(1)
    }

    pdfPath := os.Args[1]

    // Verify file exists
    if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
        fmt.Printf("Error: file not found: %s\n", pdfPath)
        os.Exit(1)
    }

    // Create model
    model, err := ui.NewModel(pdfPath)
    if err != nil {
        fmt.Printf("Error: failed to load PDF: %v\n", err)
        os.Exit(1)
    }

    // Run Bubble Tea program
    p := tea.NewProgram(
        model,
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    )

    finalModel, err := p.Run()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    // Cleanup
    if m, ok := finalModel.(ui.Model); ok {
        m.Close()
    }
}
```

---

## go.mod

```go
module lumos

go 1.21

require (
    github.com/charmbracelet/bubbles v0.18.0
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/hashicorp/golang-lru v1.0.2
    github.com/ledongthuc/pdf v0.0.0-20240511090121-5959a4027728
)
```

---

## Build & Run

### Development

```bash
# Clone or create project
mkdir -p lumos
cd lumos

# Initialize module
go mod init lumos

# Install dependencies
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/lipgloss
go get github.com/hashicorp/golang-lru
go get github.com/ledongthuc/pdf

# Build
go build -o lumos ./cmd/lumos

# Run
./lumos example.pdf
```

### Production Build

```bash
# Build optimized binary
go build -ldflags="-s -w" -o lumos ./cmd/lumos

# Optional: compress with UPX
upx --best --lzma lumos
```

---

## Testing

### Unit Tests

```go
// internal/document/manager_test.go
package document

import (
    "testing"
)

func TestManager_GetPage(t *testing.T) {
    mgr, err := NewManager("testdata/sample.pdf")
    if err != nil {
        t.Fatalf("failed to create manager: %v", err)
    }
    defer mgr.Close()

    // Test valid page
    content, err := mgr.GetPage(1)
    if err != nil {
        t.Errorf("failed to get page 1: %v", err)
    }

    if content == nil {
        t.Error("expected content, got nil")
    }

    // Test cache hit
    content2, err := mgr.GetPage(1)
    if err != nil {
        t.Errorf("failed to get cached page: %v", err)
    }

    stats := mgr.GetStats()
    if stats.Hits < 1 {
        t.Errorf("expected cache hit, got %d hits", stats.Hits)
    }

    // Test invalid page
    _, err = mgr.GetPage(9999)
    if err == nil {
        t.Error("expected error for invalid page, got nil")
    }
}
```

### Integration Tests

```go
// internal/ui/model_test.go
package ui

import (
    "testing"

    tea "github.com/charmbracelet/bubbletea"
)

func TestModel_Navigation(t *testing.T) {
    model, err := NewModel("testdata/sample.pdf")
    if err != nil {
        t.Fatalf("failed to create model: %v", err)
    }

    // Simulate key press
    model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

    if model.currentPage != 2 {
        t.Errorf("expected page 2, got %d", model.currentPage)
    }
}
```

---

## Performance Profiling

### CPU Profile

```bash
# Build with profiling
go build -o lumos ./cmd/lumos

# Run with CPU profiling
./lumos -cpuprofile=cpu.prof example.pdf

# Analyze
go tool pprof cpu.prof
```

### Memory Profile

```bash
# Run with memory profiling
./lumos -memprofile=mem.prof example.pdf

# Analyze
go tool pprof mem.prof
```

### Benchmarks

```go
// internal/rendering/text_bench_test.go
package rendering

import (
    "strings"
    "testing"
)

func BenchmarkTextRenderer_Render(b *testing.B) {
    renderer := NewTextRenderer(80, 24)
    text := strings.Repeat("Lorem ipsum dolor sit amet. ", 1000)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = renderer.Render(text)
    }
}
```

---

## Configuration File (Future)

### ~/.config/lumos/config.yaml

```yaml
# Cache settings
cache:
  size: 5
  preload_adjacent: true

# Rendering
rendering:
  default_mode: text
  text_threshold: 100
  max_image_width: 1200

# Terminal
terminal:
  auto_detect: true
  preferred_protocol: kitty
  fallback_protocol: iterm2

# Keybindings
keybindings:
  quit: ["q", "ctrl+c"]
  next_page: ["j", "down", "space"]
  prev_page: ["k", "up"]
  first_page: ["g", "home"]
  last_page: ["G", "end"]
  search: ["/"]
```

---

## Future Enhancements

### Phase 2: Search

```go
// internal/search/index.go
type SearchIndex struct {
    pages map[int][]string  // page -> words
}

func (s *SearchIndex) Search(query string) []SearchResult {
    // Full-text search across all pages
}
```

### Phase 3: Table of Contents

```go
// internal/document/toc.go
type TOCEntry struct {
    Title    string
    PageNum  int
    Level    int
    Children []*TOCEntry
}

func (m *Manager) ExtractTOC() ([]*TOCEntry, error) {
    // Extract from PDF outline/bookmarks
}
```

### Phase 4: Image Rendering

```go
// internal/rendering/image.go
func (r *ImageRenderer) RenderPageAsImage(page int) (string, error) {
    // Convert PDF page to image
    // Detect terminal protocol
    // Encode and return
}
```

---

## Status

- ✅ Architecture defined
- ✅ Core types designed
- ✅ Document manager implemented
- ✅ Text rendering implemented
- ✅ Bubble Tea UI implemented
- ✅ Build system ready
- 🟡 Testing infrastructure (partial)
- ⏳ Image rendering (future)
- ⏳ Search (future)
- ⏳ TOC extraction (future)

---

**Next Steps**:
1. Create project structure
2. Implement core components
3. Test with sample PDFs
4. Iterate on rendering quality
5. Add image fallback
6. Performance optimization

**Estimated Timeline**: 2-3 weeks to working MVP
