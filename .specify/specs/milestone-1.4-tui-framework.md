# Milestone 1.4: Basic TUI Framework

**Status**: In Progress
**Priority**: P0 (Blocker for 1.5, 1.6)
**Estimated Duration**: 2-3 days
**Dependencies**: Milestone 1.3 (✅ Complete)

---

## Overview

Implement a functional Terminal User Interface using Bubble Tea framework with 3-pane layout, viewport rendering, and basic page navigation.

## User Stories

### US-1.4.1: As a user, I can launch LUMOS with a PDF file
**Given** I have a PDF file at `~/Documents/paper.pdf`
**When** I run `./build/lumos ~/Documents/paper.pdf`
**Then** The TUI launches in full-screen mode with the first page visible

### US-1.4.2: As a user, I can see a 3-pane layout
**Given** LUMOS is running
**When** The interface renders
**Then** I see three vertical panes: Metadata (20%), Viewer (60%), Search (20%)

### US-1.4.3: As a user, I can view PDF page content
**Given** A PDF is loaded
**When** I'm on page 1
**Then** The viewer pane shows the extracted text content

### US-1.4.4: As a user, I can see document metadata
**Given** A PDF is loaded
**When** The metadata pane renders
**Then** I see: filename, page count, current page, title, author

### US-1.4.5: As a user, I can quit the application
**Given** LUMOS is running
**When** I press `q` or `Ctrl+C`
**Then** The application exits cleanly to the terminal

---

## Acceptance Criteria

### AC-1.4.1: Bubble Tea Integration
- [ ] `tea.NewProgram()` initialized in `cmd/lumos/main.go`
- [ ] `tea.WithAltScreen()` used for full-screen mode
- [ ] Model implements `tea.Model` interface (Init, Update, View)
- [ ] Program runs without panics

### AC-1.4.2: Model State Management
```go
type Model struct {
    document      *pdf.Document     // ✅ Required
    cache         *pdf.LRUCache     // ✅ Required
    currentPage   int               // ✅ Required
    theme         config.Theme      // ✅ Required
    styles        config.Styles     // ✅ Required
    viewport      viewport.Model    // ✅ Required
    width, height int               // ✅ Required
}
```
- [ ] All required fields populated in `NewModel()`
- [ ] State is immutable (only changed via Update())

### AC-1.4.3: Three-Pane Layout
- [ ] Metadata pane renders at 20% width
- [ ] Viewer pane renders at 60% width
- [ ] Search pane renders at 20% width
- [ ] Panes adjust to terminal width
- [ ] Minimum width of 80 columns supported

### AC-1.4.4: Viewport Component
- [ ] Uses `bubbles/viewport` for scrollable content
- [ ] Viewport height = terminal height - 2 (status bar)
- [ ] Content scrollable with arrow keys
- [ ] Shows scroll position indicator

### AC-1.4.5: Status Bar
- [ ] Shows current page / total pages
- [ ] Shows filename
- [ ] Shows key hints: `[?] Help [q] Quit`
- [ ] Spans full terminal width
- [ ] Uses distinct styling from content

### AC-1.4.6: Page Loading
- [ ] Async page loading via `tea.Cmd`
- [ ] Loading indicator for operations >100ms
- [ ] Error messages shown in viewport
- [ ] First page loads on Init()

### AC-1.4.7: Window Resize Handling
- [ ] `tea.WindowSizeMsg` updates model dimensions
- [ ] Viewport resizes correctly
- [ ] Panes recalculate widths
- [ ] No visual artifacts on resize

### AC-1.4.8: Performance
- [ ] Initial render < 70ms (measured with time command)
- [ ] Viewport updates < 16ms (60 FPS)
- [ ] No memory leaks (run for 5 minutes, check memory)
- [ ] Clean shutdown (no goroutine leaks)

---

## Implementation Checklist

### Phase 1: Basic Structure (Day 1 Morning)
- [ ] Update `cmd/lumos/main.go`:
  - [ ] Import Bubble Tea packages
  - [ ] Initialize tea.Program with alt screen
  - [ ] Pass document to NewModel()
  - [ ] Run program and handle errors

- [ ] Update `pkg/ui/model.go`:
  - [ ] Add viewport field
  - [ ] Implement Init() - load first page
  - [ ] Implement basic Update() - handle quit
  - [ ] Implement basic View() - render viewport

- [ ] Create `pkg/ui/messages.go`:
  - [ ] Define PageLoadedMsg
  - [ ] Define WindowSizeMsg handling
  - [ ] Define QuitMsg

### Phase 2: Three-Pane Layout (Day 1 Afternoon)
- [ ] Create `pkg/ui/layout.go`:
  - [ ] `calculatePaneWidths(totalWidth) (int, int, int)`
  - [ ] `renderMetadataPane(width, height) string`
  - [ ] `renderViewerPane(width, height) string`
  - [ ] `renderSearchPane(width, height) string`

- [ ] Update View() in model.go:
  - [ ] Call layout functions
  - [ ] Use lipgloss.JoinHorizontal()
  - [ ] Add status bar at bottom

### Phase 3: Metadata & Styling (Day 2 Morning)
- [ ] Implement metadata pane:
  - [ ] Show document title
  - [ ] Show author (if available)
  - [ ] Show page count
  - [ ] Show current page highlight
  - [ ] Add border with lipgloss

- [ ] Implement viewer pane:
  - [ ] Render viewport content
  - [ ] Add title "📖 Viewer - Page X"
  - [ ] Add border with lipgloss
  - [ ] Handle empty content gracefully

- [ ] Implement search pane:
  - [ ] Placeholder text "🔍 Search"
  - [ ] Show "No results" when inactive
  - [ ] Add border with lipgloss

### Phase 4: Status Bar & Polish (Day 2 Afternoon)
- [ ] Create `pkg/ui/statusbar.go`:
  - [ ] Format: `Page 1/10 | paper.pdf | [?] Help [q] Quit`
  - [ ] Use theme colors
  - [ ] Center filename, align page left, hints right

- [ ] Add keyboard handling:
  - [ ] `q` or `Ctrl+C` → quit
  - [ ] `?` → show help (stub for now)
  - [ ] Arrow keys → scroll viewport

### Phase 5: Testing (Day 3)
- [ ] Create `pkg/ui/model_test.go`:
  - [ ] TestNewModel - verifies initialization
  - [ ] TestUpdate_Quit - verifies quit handling
  - [ ] TestView_Layout - verifies 3-pane structure
  - [ ] TestWindowResize - verifies resize handling

- [ ] Manual testing:
  - [ ] Test with simple.pdf (1 page)
  - [ ] Test with multipage.pdf (5 pages)
  - [ ] Test window resize
  - [ ] Test on small terminal (80x24)
  - [ ] Test on large terminal (200x60)

- [ ] Performance testing:
  - [ ] Measure startup time: `time ./build/lumos test/fixtures/simple.pdf`
  - [ ] Profile memory: Monitor with Activity Monitor
  - [ ] Test scroll smoothness

---

## Technical Design

### Message Flow
```
User Launches App
       ↓
tea.NewProgram(NewModel(document))
       ↓
Model.Init() → LoadPageCmd(1)
       ↓
[Async] Load page from PDF
       ↓
PageLoadedMsg → Update()
       ↓
viewport.SetContent(pageText)
       ↓
View() renders 3-pane layout
       ↓
Terminal displays UI
```

### Layout Structure
```
┌────────────────────────────────────────────────────────────┐
│ ┌──────┐ ┌─────────────────────────┐ ┌───────────────┐  │
│ │      │ │                         │ │               │  │
│ │ Meta │ │       Viewer            │ │    Search     │  │
│ │ 20%  │ │        60%              │ │      20%      │  │
│ │      │ │                         │ │               │  │
│ └──────┘ └─────────────────────────┘ └───────────────┘  │
│                                                           │
│ Page 1/10 | paper.pdf              [?] Help [q] Quit    │
└────────────────────────────────────────────────────────────┘
```

### Viewport Configuration
```go
vp := viewport.New(viewerWidth, height-2)
vp.Style = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("62"))
```

---

## Test Requirements

### Unit Tests (Target: 15 tests)
```go
// pkg/ui/model_test.go
TestNewModel_InitializesAllFields
TestInit_ReturnsLoadPageCommand
TestUpdate_QuitKey_ExitsProgram
TestUpdate_PageLoadedMsg_UpdatesViewport
TestUpdate_WindowSizeMsg_ResizesPanes
TestView_EmptyDocument_ShowsPlaceholder
TestView_LoadedDocument_Shows3Panes
TestCalculatePaneWidths_80Columns_Returns16_48_16
TestCalculatePaneWidths_200Columns_Returns40_120_40
TestRenderStatusBar_FormatsCorrectly

// pkg/ui/layout_test.go
TestRenderMetadataPane_ShowsPageCount
TestRenderViewerPane_ShowsPageNumber
TestRenderSearchPane_ShowsPlaceholder
TestJoinPanes_CreatesCorrectLayout
```

### Integration Tests (Target: 5 tests)
```go
// test/test_tui.go
TestTUI_Launch_FirstPageDisplays
TestTUI_Quit_ExitsCleanly
TestTUI_Resize_AdjustsLayout
TestTUI_ScrollViewport_UpdatesPosition
TestTUI_LoadError_ShowsErrorMessage
```

### Performance Tests
```bash
# Startup benchmark
time ./build/lumos test/fixtures/simple.pdf
# Must show: real < 0.070s

# Memory test
# 1. Launch with large PDF
# 2. Monitor memory for 5 minutes
# 3. Quit and verify no leaks
```

---

## Success Criteria

### Functional
- [x] Application launches without errors
- [x] 3-pane layout renders correctly
- [x] First page content displays in viewer
- [x] Metadata pane shows document info
- [x] Status bar shows page number and hints
- [x] Quit command (q/Ctrl+C) works
- [x] Window resize handled gracefully

### Quality
- [x] All tests pass (minimum 15 unit + 5 integration)
- [x] Test coverage >80% for ui package
- [x] No linter warnings
- [x] No race conditions (go test -race passes)

### Performance
- [x] Startup time <70ms
- [x] UI responsive at 60 FPS
- [x] Memory usage <10MB baseline
- [x] Clean shutdown, no goroutine leaks

### Documentation
- [x] Update PROGRESS.md with milestone status
- [x] Create PHASE_1_MILESTONE_1_4_REVIEW.md
- [x] Update README.md with TUI screenshot (text art)
- [x] Document any architectural decisions

---

## Dependencies

### Upstream
- ✅ Milestone 1.3 complete (test fixtures available)
- ✅ pkg/pdf package complete (document loading works)
- ✅ pkg/config package complete (themes available)

### Downstream
- Milestone 1.5 requires: Working TUI framework with viewport
- Milestone 1.6 requires: Layout and styling in place

---

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Bubble Tea API changes | High | Use exact version in go.mod |
| Terminal compatibility issues | Medium | Test on iTerm2, Terminal.app, Alacritty |
| Performance on large PDFs | Medium | Use fixtures, defer large PDF testing to 1.5 |
| Layout calculations edge cases | Low | Extensive test coverage for widths |

---

## Definition of Done

Milestone 1.4 is **DONE** when:

1. ✅ All acceptance criteria met
2. ✅ All tests passing (>80% coverage)
3. ✅ Performance targets achieved
4. ✅ Manual testing complete on 3 terminals
5. ✅ Code reviewed and merged
6. ✅ Documentation updated
7. ✅ Milestone review document created
8. ✅ Ready to start Milestone 1.5

**Deliverable**: Working TUI that displays PDF content with 3-pane layout and basic navigation.

---

**Created**: 2025-11-10
**Target Completion**: 2025-11-13
**Assignee**: Claude Code + Human Reviewer
