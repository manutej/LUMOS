# TUI Implementation Plan

**Milestone**: 1.4 - Basic TUI Framework
**Duration**: 2-3 days
**Status**: Ready to implement

---

## Day 1: Foundation (8 hours)

### Morning Session (4 hours)

#### Hour 1: Bubble Tea Setup
```go
// cmd/lumos/main.go modifications
import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/viewport"
)

func main() {
    // Existing PDF loading...

    // Initialize TUI
    model := ui.NewModel(document)
    p := tea.NewProgram(
        model,
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    )

    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

#### Hour 2: Model Structure
```go
// pkg/ui/model.go
type Model struct {
    // Document state
    document      *pdf.Document
    cache         *pdf.LRUCache
    currentPage   int
    totalPages    int
    pageContent   string

    // UI components
    viewport      viewport.Model
    metadataPane  string
    searchPane    string
    statusBar     string

    // Display state
    width, height int
    ready         bool
    loading       bool
    err           error

    // Theme
    theme         config.Theme
    styles        *Styles
}

func NewModel(doc *pdf.Document) Model {
    return Model{
        document:    doc,
        cache:       pdf.NewLRUCache(5),
        currentPage: 1,
        totalPages:  doc.PageCount(),
        theme:       config.DarkTheme,
        styles:      DefaultStyles(),
    }
}
```

#### Hour 3: Init and Basic Update
```go
func (m Model) Init() tea.Cmd {
    return tea.Batch(
        LoadPageCmd(m.currentPage),
        tea.EnterAltScreen,
    )
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "j":
            // Scroll down
            m.viewport.LineDown(1)
        case "k":
            // Scroll up
            m.viewport.LineUp(1)
        }

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.viewport.Width = m.calculateViewerWidth()
        m.viewport.Height = msg.Height - 2 // Status bar

    case PageLoadedMsg:
        m.pageContent = msg.Content
        m.viewport.SetContent(msg.Content)
        m.loading = false
    }

    // Update viewport
    m.viewport, cmd := m.viewport.Update(msg)
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}
```

#### Hour 4: Basic View
```go
func (m Model) View() string {
    if !m.ready {
        return "Initializing..."
    }

    if m.err != nil {
        return fmt.Sprintf("Error: %v", m.err)
    }

    // For now, just viewport
    return m.viewport.View() + "\n" + m.renderStatusBar()
}

func (m Model) renderStatusBar() string {
    status := fmt.Sprintf(
        " Page %d/%d | %s | [?]Help [q]Quit",
        m.currentPage,
        m.totalPages,
        m.document.Filename(),
    )
    return m.styles.StatusBar.Render(status)
}
```

### Afternoon Session (4 hours)

#### Hour 5: Messages and Commands
```go
// pkg/ui/messages.go
type PageLoadedMsg struct {
    PageNum int
    Content string
}

type ErrorMsg struct {
    Err error
}

func LoadPageCmd(pageNum int) tea.Cmd {
    return func() tea.Msg {
        // This runs async
        content, err := loadPageContent(pageNum)
        if err != nil {
            return ErrorMsg{Err: err}
        }
        return PageLoadedMsg{
            PageNum: pageNum,
            Content: content,
        }
    }
}
```

#### Hour 6: Three-Pane Layout Structure
```go
// pkg/ui/layout.go
func (m Model) calculatePaneWidths() (meta, viewer, search int) {
    total := m.width
    meta = total * 20 / 100
    search = total * 20 / 100
    viewer = total - meta - search - 4 // Account for borders
    return
}

func (m Model) renderMetadataPane() string {
    width, _, _ := m.calculatePaneWidths()

    content := fmt.Sprintf(
        "📄 %s\n\n" +
        "Page: %d/%d\n" +
        "Author: %s\n" +
        "Title: %s",
        m.document.Filename(),
        m.currentPage,
        m.totalPages,
        m.document.Author(),
        m.document.Title(),
    )

    return m.styles.MetadataPane.
        Width(width).
        Height(m.height - 2).
        Render(content)
}
```

#### Hour 7: Viewport Integration
```go
func (m Model) renderViewerPane() string {
    _, width, _ := m.calculatePaneWidths()

    // Create viewport if not exists
    if m.viewport.Width == 0 {
        m.viewport = viewport.New(width, m.height-2)
        m.viewport.Style = m.styles.ViewerPane
    }

    title := fmt.Sprintf("📖 Page %d", m.currentPage)
    pane := m.styles.ViewerPane.
        Width(width).
        Render(title + "\n" + m.viewport.View())

    return pane
}
```

#### Hour 8: Join Panes
```go
func (m Model) View() string {
    if !m.ready {
        return m.styles.Loading.Render("Loading...")
    }

    // Render three panes
    metadata := m.renderMetadataPane()
    viewer := m.renderViewerPane()
    search := m.renderSearchPane()

    // Join horizontally
    content := lipgloss.JoinHorizontal(
        lipgloss.Top,
        metadata,
        viewer,
        search,
    )

    // Add status bar
    return lipgloss.JoinVertical(
        lipgloss.Left,
        content,
        m.renderStatusBar(),
    )
}
```

---

## Day 2: Refinement (8 hours)

### Morning Session (4 hours)

#### Hour 9: Styles and Themes
```go
// pkg/ui/styles.go
type Styles struct {
    MetadataPane lipgloss.Style
    ViewerPane   lipgloss.Style
    SearchPane   lipgloss.Style
    StatusBar    lipgloss.Style
    Loading      lipgloss.Style
    Error        lipgloss.Style
}

func DefaultStyles() *Styles {
    return &Styles{
        MetadataPane: lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("62")).
            Padding(1),

        ViewerPane: lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("63")).
            Padding(1),

        SearchPane: lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("64")).
            Padding(1),

        StatusBar: lipgloss.NewStyle().
            Background(lipgloss.Color("235")).
            Foreground(lipgloss.Color("252")),
    }
}
```

#### Hour 10: Page Navigation
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+n":
            if m.currentPage < m.totalPages {
                m.currentPage++
                return m, LoadPageCmd(m.currentPage)
            }
        case "ctrl+p":
            if m.currentPage > 1 {
                m.currentPage--
                return m, LoadPageCmd(m.currentPage)
            }
        case "g":
            if m.lastKey == "g" {
                // gg - go to first page
                m.currentPage = 1
                return m, LoadPageCmd(m.currentPage)
            }
            m.lastKey = "g"
        case "G":
            // Go to last page
            m.currentPage = m.totalPages
            return m, LoadPageCmd(m.currentPage)
        }
    }
}
```

#### Hour 11: Search Pane Placeholder
```go
func (m Model) renderSearchPane() string {
    _, _, width := m.calculatePaneWidths()

    content := "🔍 Search\n\n"
    if m.searchActive {
        content += fmt.Sprintf("Query: %s\n", m.searchQuery)
        content += fmt.Sprintf("Results: %d\n", len(m.searchResults))
    } else {
        content += "Press / to search"
    }

    return m.styles.SearchPane.
        Width(width).
        Height(m.height - 2).
        Render(content)
}
```

#### Hour 12: Window Resize Handling
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

        // Recalculate pane widths
        _, viewerWidth, _ := m.calculatePaneWidths()

        // Update viewport dimensions
        m.viewport.Width = viewerWidth - 2  // Account for borders
        m.viewport.Height = m.height - 3    // Status bar + borders

        // Force re-render
        m.ready = true

        return m, nil
    }
}
```

### Afternoon Session (4 hours)

#### Hour 13: Loading States
```go
func (m Model) View() string {
    if !m.ready {
        return m.centerMessage("Initializing LUMOS...")
    }

    if m.loading {
        return m.centerMessage(fmt.Sprintf(
            "Loading page %d/%d...",
            m.currentPage,
            m.totalPages,
        ))
    }

    if m.err != nil {
        return m.centerMessage(fmt.Sprintf(
            "Error: %v\n\nPress q to quit",
            m.err,
        ))
    }

    // Normal view...
}

func (m Model) centerMessage(msg string) string {
    return lipgloss.Place(
        m.width, m.height,
        lipgloss.Center, lipgloss.Center,
        m.styles.Loading.Render(msg),
    )
}
```

#### Hour 14: Error Handling
```go
func LoadPageCmd(doc *pdf.Document, pageNum int) tea.Cmd {
    return func() tea.Msg {
        // Check bounds
        if pageNum < 1 || pageNum > doc.PageCount() {
            return ErrorMsg{
                Err: fmt.Errorf("page %d out of range (1-%d)",
                    pageNum, doc.PageCount()),
            }
        }

        // Try cache first
        if content, ok := cache.Get(pageNum); ok {
            return PageLoadedMsg{
                PageNum: pageNum,
                Content: content,
            }
        }

        // Load from PDF
        content, err := doc.GetPage(pageNum)
        if err != nil {
            return ErrorMsg{Err: err}
        }

        // Cache it
        cache.Put(pageNum, content)

        return PageLoadedMsg{
            PageNum: pageNum,
            Content: content,
        }
    }
}
```

#### Hour 15: Status Bar Polish
```go
func (m Model) renderStatusBar() string {
    // Left section: page info
    left := fmt.Sprintf(" Page %d/%d", m.currentPage, m.totalPages)

    // Center section: filename
    center := filepath.Base(m.document.Filename())

    // Right section: help
    right := "[?]Help [q]Quit "

    // Calculate spacing
    leftWidth := len(left)
    rightWidth := len(right)
    centerWidth := m.width - leftWidth - rightWidth

    if centerWidth < len(center) {
        center = center[:centerWidth-3] + "..."
    }

    status := left +
        lipgloss.PlaceHorizontal(centerWidth, lipgloss.Center, center) +
        right

    return m.styles.StatusBar.Width(m.width).Render(status)
}
```

#### Hour 16: Integration Testing
```go
// Manual test checklist:
// 1. Launch with simple.pdf
// 2. Verify 3-pane layout
// 3. Press j/k to scroll
// 4. Press Ctrl+N/P to change pages
// 5. Resize terminal
// 6. Press q to quit

// Automated test:
func TestTUILaunch(t *testing.T) {
    doc := loadTestPDF(t)
    model := NewModel(doc)

    // Test initialization
    cmd := model.Init()
    assert.NotNil(t, cmd)

    // Test view renders
    view := model.View()
    assert.Contains(t, view, "Page 1")

    // Test quit command
    model, cmd = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
    assert.IsType(t, tea.Quit(), cmd)
}
```

---

## Day 3: Polish & Testing (4 hours)

### Morning Session (4 hours)

#### Hour 17: Performance Optimization
- Profile render performance
- Cache styled strings
- Minimize allocations in View()
- Verify 60 FPS scrolling

#### Hour 18: Edge Cases
- Very small terminals (80x24)
- Very large terminals (300x100)
- Single page PDFs
- PDFs with no metadata
- Long filenames

#### Hour 19: Unit Tests
Write tests for:
- Model initialization
- Pane width calculations
- Message handling
- State transitions
- Error cases

#### Hour 20: Documentation
- Update PROGRESS.md
- Create milestone review
- Update README with screenshots
- Document architecture decisions

---

## Testing Checklist

### Automated Tests
- [ ] Model initialization
- [ ] Window size handling
- [ ] Page navigation
- [ ] Error handling
- [ ] View output format

### Manual Tests
- [ ] Launch with each test PDF
- [ ] All navigation keys work
- [ ] Resize at various sizes
- [ ] Test on 3 terminals
- [ ] Memory usage stable

### Performance Tests
- [ ] Startup time <70ms
- [ ] Render time <16ms
- [ ] Memory <10MB baseline
- [ ] No goroutine leaks

---

## Common Issues & Solutions

### Issue: Viewport not updating
```go
// Always forward update to viewport
m.viewport, cmd := m.viewport.Update(msg)
```

### Issue: Panes misaligned
```go
// Account for borders in width calculation
viewerWidth := totalWidth - metaWidth - searchWidth - 4
```

### Issue: Window resize artifacts
```go
// Clear and redraw on resize
case tea.WindowSizeMsg:
    m.ready = false  // Force full redraw
    // ... handle resize
    m.ready = true
```

---

## Success Criteria

The TUI implementation is complete when:
1. ✅ Application launches without crashes
2. ✅ 3-pane layout renders correctly
3. ✅ PDF content displays in viewport
4. ✅ Basic navigation works (j/k, Ctrl+N/P)
5. ✅ Window resize handled gracefully
6. ✅ Status bar shows correct info
7. ✅ Clean shutdown with q
8. ✅ All tests passing
9. ✅ Performance targets met
10. ✅ Ready for keybindings milestone

---

**Estimated Total Time**: 20 hours (2.5 days)
**Buffer**: 4 hours for unexpected issues
**Deadline**: November 13, 2025