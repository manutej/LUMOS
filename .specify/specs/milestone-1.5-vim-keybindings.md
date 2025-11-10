# Milestone 1.5: Vim Keybindings & Navigation

**Status**: Not Started
**Priority**: P0 (Blocker for 1.6)
**Estimated Duration**: 2-3 days
**Dependencies**: Milestone 1.4 (TUI Framework)

---

## Overview

Implement complete vim-style keyboard navigation including line scrolling, page navigation, search mode, and help system. All keybindings must feel natural to vim users.

## User Stories

### US-1.5.1: As a vim user, I can scroll content naturally
**Given** I'm viewing a PDF page
**When** I press `j`
**Then** Content scrolls down one line

**When** I press `k`
**Then** Content scrolls up one line

### US-1.5.2: As a vim user, I can navigate pages efficiently
**Given** I'm on page 3 of a 10-page PDF
**When** I press `Ctrl+N`
**Then** I advance to page 4

**When** I press `Ctrl+P`
**Then** I return to page 3

### US-1.5.3: As a vim user, I can jump to document boundaries
**Given** I'm on page 5
**When** I type `gg`
**Then** I jump to page 1

**When** I type `G` (shift+g)
**Then** I jump to the last page

### US-1.5.4: As a user, I can search for text
**Given** I'm viewing a PDF
**When** I press `/`
**Then** I enter search mode with a prompt

**When** I type "quantum" and press Enter
**Then** The first match is highlighted and visible

**When** I press `n`
**Then** I jump to the next match

### US-1.5.5: As a user, I can access help
**Given** LUMOS is running
**When** I press `?`
**Then** A help screen overlays the main view with all keybindings

**When** I press `?` again or `Esc`
**Then** I return to the main view

---

## Acceptance Criteria

### AC-1.5.1: Normal Mode Navigation
```go
// Required keybindings in Normal mode
"j", "down"    → viewport.LineDown(1)
"k", "up"      → viewport.LineUp(1)
"d"            → viewport.LineDown(viewport.Height/2)  // Half page
"u"            → viewport.LineUp(viewport.Height/2)    // Half page
"gg"           → goToFirstPage()
"G"            → goToLastPage()
"ctrl+n"       → goToNextPage()
"ctrl+p"       → goToPreviousPage()
"ctrl+f"       → viewport.LineDown(viewport.Height)    // Full page
"ctrl+b"       → viewport.LineUp(viewport.Height)      // Full page
```

- [ ] All keybindings implemented
- [ ] Keybindings work at viewport boundaries (don't panic)
- [ ] Visual feedback for all navigation actions
- [ ] Smooth scrolling (no flicker)

### AC-1.5.2: Search Mode
```go
type KeyMode int
const (
    KeyModeNormal KeyMode = iota
    KeyModeSearch
    KeyModeCommand  // Phase 2
)
```

- [ ] `/` enters search mode
- [ ] Search prompt shows at bottom: `Search: _`
- [ ] Typing updates search query display
- [ ] `Enter` executes search
- [ ] `Esc` exits search mode without searching
- [ ] `n` jumps to next match
- [ ] `N` jumps to previous match
- [ ] Matches are highlighted in search pane
- [ ] Current match shown in status bar: `Match 3/17`

### AC-1.5.3: Key Handler Implementation
```go
// pkg/ui/keybindings.go
type KeyHandler struct {
    Mode         KeyMode
    Buffer       string  // For multi-key sequences like "gg"
    LastKeyTime  time.Time
}

func (k *KeyHandler) HandleKey(msg tea.KeyMsg) KeyAction
func (k *KeyHandler) IsSequenceComplete(keys string) bool
func (k *KeyHandler) ClearBuffer()
```

- [ ] KeyHandler tracks current mode
- [ ] Buffer handles multi-key sequences (gg, G, etc.)
- [ ] Timeout for incomplete sequences (500ms)
- [ ] Thread-safe key handling

### AC-1.5.4: Page Navigation Commands
```go
func (m *Model) goToNextPage() tea.Cmd
func (m *Model) goToPreviousPage() tea.Cmd
func (m *Model) goToFirstPage() tea.Cmd
func (m *Model) goToLastPage() tea.Cmd
func (m *Model) goToPage(pageNum int) tea.Cmd
```

- [ ] Page transitions are smooth
- [ ] Loading indicator shown for uncached pages
- [ ] Boundary checking (don't go beyond first/last)
- [ ] Page number updates in status bar
- [ ] Viewport resets to top on page change

### AC-1.5.5: Help Screen
```go
// pkg/ui/help.go
func renderHelpScreen() string
```

Content:
```
LUMOS - Keyboard Shortcuts

NAVIGATION
  j/k or ↑/↓     Scroll line up/down
  d/u            Scroll half page
  Ctrl+F/B       Scroll full page
  gg / G         First / Last page
  Ctrl+N/P       Next / Previous page

SEARCH
  /              Start search
  n / N          Next / Previous match
  Esc            Exit search

UI
  Tab            Cycle panes
  1 / 2          Dark / Light mode
  ?              Toggle help

GENERAL
  q or Ctrl+C    Quit

Press ? or Esc to close
```

- [ ] Help text formatted with lipgloss
- [ ] Centered on screen
- [ ] Semi-transparent background
- [ ] Easy to read contrast
- [ ] All current keybindings documented

### AC-1.5.6: Pane Cycling
- [ ] `Tab` cycles through panes: Metadata → Viewer → Search
- [ ] Active pane highlighted with border color
- [ ] Active pane shown in status bar
- [ ] Arrow keys scroll only active pane

### AC-1.5.7: Theme Toggle
- [ ] `1` switches to dark theme
- [ ] `2` switches to light theme
- [ ] Theme applies to all panes immediately
- [ ] Status bar updates theme name
- [ ] Vim keybindings remain active

---

## Implementation Checklist

### Phase 1: Key Handler Foundation (Day 1 Morning)
- [ ] Create `pkg/ui/keybindings.go`:
  - [ ] Define KeyMode enum
  - [ ] Define KeyHandler struct
  - [ ] Implement HandleKey() method
  - [ ] Add multi-key sequence support

- [ ] Update Model struct:
  - [ ] Add keyHandler field
  - [ ] Add searchQuery field
  - [ ] Add searchResults field
  - [ ] Add currentMatch field

### Phase 2: Normal Mode Navigation (Day 1 Afternoon)
- [ ] Implement basic navigation:
  - [ ] j/k for line scrolling
  - [ ] d/u for half-page scrolling
  - [ ] Ctrl+F/B for full-page scrolling

- [ ] Implement page navigation:
  - [ ] Ctrl+N/P for next/previous page
  - [ ] gg for first page (handle double-key)
  - [ ] G for last page

- [ ] Add visual feedback:
  - [ ] Flash border on page change
  - [ ] Update status bar immediately
  - [ ] Show loading indicator

### Phase 3: Search Mode (Day 2 Morning)
- [ ] Create `pkg/ui/search.go`:
  - [ ] `enterSearchMode()` - switch to search mode
  - [ ] `exitSearchMode()` - switch back to normal
  - [ ] `executeSearch()` - run Document.Search()
  - [ ] `jumpToMatch(index)` - navigate to result

- [ ] Implement search UI:
  - [ ] Search prompt at bottom
  - [ ] Show query as typed
  - [ ] Display match count
  - [ ] Highlight current match

- [ ] Implement n/N navigation:
  - [ ] n → currentMatch++, wrap to 0
  - [ ] N → currentMatch--, wrap to len-1
  - [ ] Jump to result page
  - [ ] Scroll to result position

### Phase 4: Help & Polish (Day 2 Afternoon)
- [ ] Create `pkg/ui/help.go`:
  - [ ] Design help layout
  - [ ] Format with lipgloss
  - [ ] Add all keybindings
  - [ ] Add examples

- [ ] Implement pane cycling:
  - [ ] Track activePaneIdx
  - [ ] Highlight active pane border
  - [ ] Route keys to active pane

- [ ] Theme toggle:
  - [ ] Load dark theme
  - [ ] Load light theme
  - [ ] Apply to all styles
  - [ ] Persist in Model

### Phase 5: Testing (Day 3)
- [ ] Create `pkg/ui/keybindings_test.go`:
  - [ ] TestHandleKey_NormalMode
  - [ ] TestHandleKey_MultiKeySequence
  - [ ] TestHandleKey_SearchMode
  - [ ] TestHandleKey_Timeout

- [ ] Create `pkg/ui/navigation_test.go`:
  - [ ] TestGoToNextPage
  - [ ] TestGoToPreviousPage
  - [ ] TestGoToFirstPage
  - [ ] TestGoToLastPage
  - [ ] TestBoundaryChecking

- [ ] Manual testing:
  - [ ] Test all vim keybindings
  - [ ] Test search with multiple matches
  - [ ] Test edge cases (empty query, no matches)
  - [ ] Test help screen toggle
  - [ ] Verify vim muscle memory feels right

---

## Technical Design

### Key Handling State Machine
```
┌─────────────┐
│  Normal     │◄───────────────┐
│  Mode       │                │
└──────┬──────┘                │
       │                       │
       │ Press /               │ Esc
       ▼                       │
┌─────────────┐                │
│  Search     │                │
│  Mode       │────────────────┘
└─────────────┘     Enter
       │
       │ Executes Search
       │
       ▼
  [Results stored]
  [Jump to first match]
  [Back to Normal mode]
```

### Multi-Key Sequence Handling
```go
// Example: "gg" sequence
keyHandler.Buffer = ""
User presses 'g' → Buffer = "g", wait 500ms
User presses 'g' → Buffer = "gg", isComplete("gg") → Execute goToFirstPage()
                   Buffer = ""

// Timeout
User presses 'g' → Buffer = "g"
500ms pass       → Buffer = "", no action
```

### Search Flow
```
User: /
   → Mode = Search
   → Prompt shown: "Search: _"

User: types "quantum"
   → searchQuery = "quantum"
   → Display updates in real-time

User: Enter
   → results = document.Search("quantum")
   → currentMatch = 0
   → Jump to results[0].PageNum
   → Mode = Normal

User: n
   → currentMatch++
   → Jump to results[currentMatch].PageNum
```

---

## Test Requirements

### Unit Tests (Target: 25 tests)
```go
// pkg/ui/keybindings_test.go
TestKeyHandler_NormalMode_ScrollDown
TestKeyHandler_NormalMode_ScrollUp
TestKeyHandler_NormalMode_HalfPageDown
TestKeyHandler_NormalMode_HalfPageUp
TestKeyHandler_MultiKey_GG_GoToFirst
TestKeyHandler_MultiKey_ShiftG_GoToLast
TestKeyHandler_MultiKey_Timeout_ClearBuffer
TestKeyHandler_SearchMode_EnterOnSlash
TestKeyHandler_SearchMode_TypeCharacters
TestKeyHandler_SearchMode_ExitOnEscape
TestKeyHandler_SearchMode_ExecuteOnEnter

// pkg/ui/navigation_test.go
TestGoToNextPage_ValidPage_IncrementsCurrent
TestGoToNextPage_LastPage_StaysOnLast
TestGoToPreviousPage_ValidPage_DecrementsCurrent
TestGoToPreviousPage_FirstPage_StaysOnFirst
TestGoToFirstPage_SetsPageToOne
TestGoToLastPage_SetsPageToMax
TestGoToPage_ValidNumber_LoadsPage
TestGoToPage_InvalidNumber_ReturnsError

// pkg/ui/search_test.go
TestSearch_ValidQuery_ReturnsResults
TestSearch_EmptyQuery_ReturnsEmpty
TestSearch_NoMatches_ReturnsEmpty
TestSearch_NextMatch_IncrementsIndex
TestSearch_NextMatch_Wraps_ToZero
TestSearch_PreviousMatch_DecrementsIndex
TestSearch_PreviousMatch_Wraps_ToLast
```

### Integration Tests (Target: 10 tests)
```go
// test/test_navigation.go (already exists, expand)
TestVimKeys_JK_ScrollsViewport
TestVimKeys_DU_ScrollsHalfPage
TestVimKeys_GG_JumpsToFirstPage
TestVimKeys_ShiftG_JumpsToLastPage
TestVimKeys_CtrlNP_ChangesPages
TestSearch_FindMatches_HighlightsResults
TestSearch_NextPrevious_NavigatesMatches
TestHelp_QuestionMark_TogglesHelp
TestPaneCycle_Tab_CyclesPanes
TestThemeToggle_12_SwitchesThemes
```

### Manual Test Scenarios
1. **Vim Muscle Memory Test**:
   - Navigate using only j/k/d/u
   - Should feel identical to vim

2. **Search Workflow**:
   - Search for common word
   - Navigate through all matches with n
   - Verify wrap-around behavior

3. **Edge Cases**:
   - Search with no results
   - Navigate at document boundaries
   - Multi-key sequence interruption

---

## Success Criteria

### Functional
- [x] All vim navigation keys work (j/k/d/u/gg/G/Ctrl+N/Ctrl+P)
- [x] Multi-key sequences (gg, G) work correctly
- [x] Search mode accessible via `/`
- [x] Search results navigable with n/N
- [x] Help screen shows all keybindings
- [x] Pane cycling works with Tab
- [x] Theme toggle works with 1/2
- [x] Quit works with q/Ctrl+C

### Quality
- [x] All tests pass (minimum 25 unit + 10 integration)
- [x] Test coverage >85% for ui package
- [x] No key handling bugs (buffer leaks, race conditions)
- [x] Vim users confirm "feels right"

### Performance
- [x] Key response time <16ms (60 FPS maintained)
- [x] Search executes in <100ms for typical PDFs
- [x] Page navigation <50ms (cached)
- [x] Multi-key timeout exactly 500ms

### Documentation
- [x] Update PROGRESS.md with milestone status
- [x] Create PHASE_1_MILESTONE_1_5_REVIEW.md
- [x] Document all keybindings in README.md
- [x] Create keybindings reference card

---

## Dependencies

### Upstream
- ✅ Milestone 1.4 complete (TUI framework working)
- ✅ Viewport component available
- ✅ Document search implemented (pkg/pdf/search.go)

### Downstream
- Milestone 1.6 requires: Full navigation for polish testing

---

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Multi-key timing issues | Medium | Extensive testing, adjustable timeout |
| Search performance on large PDFs | Medium | Add progress indicator, async search |
| Vim keybindings conflicts | Low | Follow vim conventions exactly |
| Help screen obscures content | Low | Clear visual design, easy dismissal |

---

## Definition of Done

Milestone 1.5 is **DONE** when:

1. ✅ All acceptance criteria met
2. ✅ All tests passing (>85% coverage)
3. ✅ Performance targets achieved
4. ✅ Vim user validation complete
5. ✅ All keybindings documented
6. ✅ Code reviewed and merged
7. ✅ Milestone review document created
8. ✅ Ready to start Milestone 1.6

**Deliverable**: Fully functional vim-style keyboard navigation with search and help system.

---

**Created**: 2025-11-10
**Target Completion**: 2025-11-16
**Assignee**: Claude Code + Vim User Reviewer
