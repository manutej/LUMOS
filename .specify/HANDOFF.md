# LUMOS Developer Handoff Guide

**Welcome!** This guide gets you productive with LUMOS in 15 minutes.

---

## Quick Context

**What**: Dark mode PDF reader for terminal-loving developers
**Tech**: Go + Bubble Tea TUI + Vim keybindings
**Status**: 50% complete (core engine done, TUI needed)
**Your Mission**: Complete the TUI to make it usable

---

## 5-Minute Setup

```bash
# 1. Clone and enter project
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMOS

# 2. Verify Go environment
go version  # Need 1.21+

# 3. Download dependencies
go mod download

# 4. Build the project
make build

# 5. Run tests (should see 42 passing)
make test

# 6. Check coverage (should see 94.4%)
make coverage

# 7. Try loading a PDF (currently CLI only)
./build/lumos test/fixtures/simple.pdf
```

**Current State**: The PDF engine works perfectly. The TUI is the missing piece.

---

## Understanding the Codebase

### Architecture Overview
```
cmd/lumos/          # Entry point (needs TUI integration)
  └─ main.go        # Currently CLI, needs tea.Program

pkg/pdf/            # ✅ COMPLETE - Don't modify
  ├─ document.go    # PDF loading and parsing
  ├─ search.go      # Full-text search (<50μs)
  └─ cache.go       # LRU cache (<100ns ops)

pkg/ui/             # 🚧 IN PROGRESS - Your focus
  ├─ model.go       # Bubble Tea MVU model
  ├─ keybindings.go # Vim key handling
  └─ messages.go    # Tea message types

pkg/config/         # ✅ COMPLETE
  └─ theme.go       # Dark/light themes
```

### Key Files to Edit
1. **cmd/lumos/main.go** - Add Bubble Tea initialization
2. **pkg/ui/model.go** - Implement Model, Update, View
3. **pkg/ui/layout.go** - Create 3-pane layout system
4. **pkg/ui/statusbar.go** - Add status bar

---

## Your First Task: Get TUI Running

### Step 1: Initialize Bubble Tea (30 min)

Edit `cmd/lumos/main.go`:
```go
import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/lumos/pkg/ui"
)

func main() {
    // ... existing PDF loading code ...

    // Initialize TUI
    model := ui.NewModel(document)
    p := tea.NewProgram(model, tea.WithAltScreen())

    if err := p.Start(); err != nil {
        log.Fatal(err)
    }
}
```

### Step 2: Implement Basic Model (45 min)

Edit `pkg/ui/model.go`:
```go
type Model struct {
    document    *pdf.Document
    currentPage int
    viewport    viewport.Model
    ready       bool
}

func (m Model) Init() tea.Cmd {
    return LoadPageCmd(1)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" {
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m Model) View() string {
    if !m.ready {
        return "Loading..."
    }
    return m.viewport.View()
}
```

### Step 3: Test It Works (15 min)
```bash
make build
./build/lumos test/fixtures/simple.pdf
# Should see TUI launch, press 'q' to quit
```

---

## Implementation Priorities

### 🔴 URGENT - TUI Framework (Milestone 1.4)
**Deadline**: 2-3 days
**Spec**: `.specify/specs/milestone-1.4-tui-framework.md`

Must complete:
1. Bubble Tea integration ← You Are Here
2. 3-pane layout (metadata | viewer | search)
3. Viewport for scrolling
4. Status bar
5. Window resize handling

### 🟠 HIGH - Vim Keybindings (Milestone 1.5)
**Timeline**: Days 4-6
**Spec**: `.specify/specs/milestone-1.5-vim-keybindings.md`

Implement:
- j/k - scroll down/up
- d/u - half page down/up
- gg/G - first/last page
- Ctrl+N/P - next/previous page
- / - search mode
- q - quit

### 🟡 MEDIUM - Polish (Milestone 1.6)
**Timeline**: Days 7-9
**Spec**: `.specify/specs/milestone-1.6-dark-mode-polish.md`

Polish:
- Theme refinement
- Performance optimization
- Terminal compatibility
- Final testing

---

## Critical Implementation Details

### The 3-Pane Layout
```
┌─────────┬──────────────────┬──────────┐
│Metadata │     Viewer       │  Search  │
│  20%    │      60%         │   20%    │
├─────────┴──────────────────┴──────────┤
│ Page 1/5 | doc.pdf | [?]Help [q]Quit  │
└────────────────────────────────────────┘
```

### State Management Rules
1. **All state in Model struct** - no globals
2. **State changes via Update()** - never direct mutation
3. **Messages for async ops** - use tea.Cmd for loading
4. **Immutable updates** - return new model

### Performance Requirements
- Startup: <100ms
- Page switch (cached): <50ms
- Scroll: 60 FPS (16ms updates)
- Memory: <50MB for 10MB PDF

---

## Testing Your Changes

### Run Existing Tests
```bash
make test          # All 42 tests should pass
make test-race     # Check for race conditions
make coverage      # Maintain >80% coverage
```

### Add New Tests
For every feature, add tests in `*_test.go`:
```go
func TestTUI_Launch_ShowsFirstPage(t *testing.T) {
    // Given
    doc := loadTestPDF(t, "simple.pdf")
    model := NewModel(doc)

    // When
    cmd := model.Init()

    // Then
    assert.NotNil(t, cmd)
    // ... verify first page loads
}
```

### Manual Testing Checklist
- [ ] Launch with simple.pdf (1 page)
- [ ] Launch with multipage.pdf (5 pages)
- [ ] Resize terminal window
- [ ] Test on 80x24 terminal (minimum)
- [ ] Test on 200x60 terminal (large)
- [ ] Verify no memory leaks (run 5 min)

---

## Common Gotchas & Solutions

### Gotcha 1: PDF Library Reopens Files
**Issue**: ledongthuc/pdf reopens file for each page
**Solution**: Already cached in LRU, just use it

### Gotcha 2: Terminal Compatibility
**Issue**: Different terminals render differently
**Solution**: Test on iTerm2, Terminal.app, Alacritty

### Gotcha 3: Viewport Not Updating
**Issue**: Forgetting to forward messages
**Solution**: Always forward WindowSizeMsg to viewport

### Gotcha 4: Memory Leaks
**Issue**: Goroutines not cleaned up
**Solution**: Use tea.Cmd, avoid manual goroutines

---

## Development Workflow

### Daily Routine
```bash
# Morning: Pull latest, check status
git pull
make test

# Start feature branch
git checkout -b feature/tui-implementation

# Development cycle
make watch     # Auto-rebuild on changes
# Edit code...
make test      # Verify nothing broke

# Before commit
make ci-check  # fmt, vet, lint, test

# Commit with spec reference
git commit -m "feat: implement 3-pane layout per spec 1.4

Implements milestone-1.4-tui-framework.md requirements:
- 20/60/20 split for panes
- Responsive to terminal width
- Minimum 80 column support"
```

### When Stuck
1. Check specification: `.specify/specs/milestone-*.md`
2. Review tests for examples: `pkg/pdf/*_test.go`
3. Check LUMINA for patterns: `../LUMINA/pkg/ui/`
4. Document blocker in PROGRESS.md

---

## Resources & References

### Project Documentation
- **Specifications**: `.specify/` - Source of truth
- **Constitution**: `.specify/constitution.md` - Quality standards
- **Progress**: `PROGRESS.md` - Current status
- **Architecture**: `docs/ARCHITECTURE.md` - System design

### External Resources
- [Bubble Tea Docs](https://github.com/charmbracelet/bubbletea)
- [Bubble Tea Examples](https://github.com/charmbracelet/bubbletea/tree/master/examples)
- [Lipgloss Styling](https://github.com/charmbracelet/lipgloss)
- [Viewport Component](https://github.com/charmbracelet/bubbles/tree/master/viewport)

### Similar Projects
- **LUMINA**: `../LUMINA/` - Markdown viewer, same stack
- Study their TUI implementation for patterns

---

## Success Checklist

You're successful when:
- [ ] TUI launches without crashes
- [ ] 3-pane layout renders correctly
- [ ] Can view PDF content in viewport
- [ ] Can quit with 'q' key
- [ ] Window resize works
- [ ] Tests maintain >80% coverage
- [ ] Performance meets targets

---

## Questions?

**Common Q&A**:

**Q: Why Bubble Tea instead of alternatives?**
A: MVU pattern provides predictable state management, essential for complex TUI.

**Q: Can I modify the PDF package?**
A: No, it's complete with 94.4% coverage. Focus on UI only.

**Q: What about images in PDFs?**
A: Phase 3 feature. Text-only for MVP.

**Q: How strict are the performance targets?**
A: Constitution defines them as MUST requirements. Use profiling if needed.

---

## Your Next 3 Hours

1. **Hour 1**: Get basic TUI launching (Step 1-3 above)
2. **Hour 2**: Add 3-pane layout structure
3. **Hour 3**: Implement viewport and test with PDFs

By end of Hour 3, you should have a working (if basic) PDF viewer!

---

**Remember**:
- Specifications drive code, not vice versa
- Test first, implement second
- When in doubt, check the constitution
- The core engine is solid - focus on UI

**Good luck! The PDF engine is waiting for its UI!** 🚀