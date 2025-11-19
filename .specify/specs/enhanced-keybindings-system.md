# Enhanced Keybindings System - Comprehensive Reference

**Status**: Specification
**Created**: 2025-11-18
**Target**: Phase 1 Enhancement (Post 1.5)
**Scope**: Replace basic vim keys with comprehensive, developer-friendly keybinding system

---

## Overview

Expand LUMOS keybindings from basic navigation to a **production-grade keybinding system** that includes:
- Full vim-style navigation with counts (5j, 10k, etc.)
- Word movement (w/b/e/W/B/E)
- Line commands (^/$, 0, g0, g$)
- Marks and jumps (ma, 'a, etc.)
- Vim-style page commands (:123 to jump to page)
- Mouse support (scroll, click for pane selection)
- Repeat functionality (. command)
- Customizable keybindings (via config file)

---

## Complete Keybindings Reference

### NAVIGATION - Character Level

#### Basic Scrolling (with repeat support)
```
j / Down        Scroll down 1 line (repeatable: 5j)
k / Up          Scroll up 1 line (repeatable: 5k)
h / Left        Scroll left 5 chars (for wide PDFs)
l / Right       Scroll right 5 chars (for wide PDFs)
```

#### Word Navigation
```
w               Jump to next word start
W               Jump to next WORD (whitespace delimited)
e               Jump to next word end
E               Jump to next WORD end
b               Jump to previous word start
B               Jump to previous WORD start
```

#### Line Navigation
```
^               Jump to first non-whitespace char on line
g^              Jump to first char on line
0               Jump to first column
g0              Jump to first visible column
$               Jump to end of line
g$              Jump to end of visible line
|               Jump to column N (pipe with number)
```

---

### NAVIGATION - Page Level

#### Page Scrolling (Vim-style, with repeat)
```
Ctrl+F          Full page down (forward page)
Ctrl+B          Full page up (backward page)
d               Half-page down (repeatable: 2d = full page)
u               Half-page up (repeatable: 2u = full page)
Ctrl+D          Half-page down (alternate)
Ctrl+U          Half-page up (alternate)
Ctrl+E          Scroll 3 lines down
Ctrl+Y          Scroll 3 lines up
```

#### Page Navigation (Document Level)
```
gg              Go to first page
G               Go to last page
Ctrl+N          Next page (repeatable: 5Ctrl+N)
Ctrl+P          Previous page (repeatable: 5Ctrl+P)
Ctrl+G          Show current page/total (status)
:N              Jump to specific page (e.g., :42 to go to page 42)
:N,M            Show range info (e.g., :1,10 shows pages 1-10)
```

---

### MARKS & JUMPS

#### Mark Position
```
ma              Set mark 'a' at current page/position
mA              Set uppercase mark 'A' (global)
```

#### Jump to Mark
```
'a              Jump to mark 'a'
'A              Jump to mark 'A' (global)
``              Jump to last mark
'.              Jump to last changed position
```

#### View Marks
```
:marks          Show all set marks
```

---

### SEARCH & REPLACE

#### Search
```
/pattern        Search forward for pattern (regex)
?pattern        Search backward for pattern
n               Next match (repeatable: 5n for 5 matches)
N               Previous match
*               Search for word under cursor (forward)
#               Search for word under cursor (backward)
gn              Select to next match (for highlight)
```

#### Search Options
```
/\c             Case-insensitive search
/\C             Case-sensitive search
:/pattern/i     Toggle ignore-case flag
```

---

### SELECTIONS & VISUAL

#### Visual Modes
```
v               Enter character-wise visual (highlight)
V               Enter line-wise visual
Ctrl+V          Enter block visual (for columns)
```

#### Selection Actions
```
y               Copy (yank) selection
p               Paste after cursor
P               Paste before cursor
```

---

### COMMANDS

#### Command Mode (Vim :commands)
```
:               Enter command mode
:q              Quit
:w              Save (if modified in future)
:wq             Save and quit
:!              Run shell command (future)
:set            Set options
:theme          Change theme (dark/light)
:help           Show help
:about          Show about info
:version        Show version
```

#### Theme Commands
```
:theme dark     Switch to dark theme
:theme light    Switch to light theme
:theme toggle   Toggle theme
:theme list     Show available themes
```

#### Search Commands
```
:/pattern       Jump to first match of pattern
:n              Jump to next match
:N              Jump to previous match
```

---

### UI CONTROLS

#### Pane Management
```
Tab             Cycle through panes (forward)
Shift+Tab       Cycle through panes (backward)
1               Focus metadata pane
2               Focus viewer pane
3               Focus search pane
Alt+N           Focus pane N
```

#### Help & Info
```
?               Toggle help screen
h               Show keybindings help
g?              Show keybindings menu (interactive)
Ctrl+H          Toggle help sidebar
```

#### Display Toggle
```
1               Switch to dark theme (alternative: :theme dark)
2               Switch to light theme (alternative: :theme light)
i               Toggle info panel
s               Toggle search panel
L               Toggle line numbers
w               Toggle word wrap
~               Toggle case display
```

---

### SPECIAL ACTIONS

#### Repeat
```
.               Repeat last command
@a              Replay macro 'a' (future)
q               Start recording macro (future)
q               Stop recording macro (future)
```

#### Undo/Redo (Future)
```
u               Undo last action
Ctrl+R          Redo
```

#### Marks Usage
```
m               Mark current position
'               Jump to mark
```

---

### MOUSE SUPPORT

#### Click Actions
```
Left Click      Focus clicked pane / move cursor to position
Double Click    Select word / toggle pane expand
Scroll Wheel    Scroll active pane
Scroll+Shift    Scroll horizontally
Scroll+Ctrl     Zoom (future)
Right Click     Context menu (future)
```

#### Drag Actions
```
Click+Drag      Select text (for copy)
```

---

### REPEAT & COUNTS

All repeatable commands support vim-style counts:

```
5j              Scroll down 5 lines
10k             Scroll up 10 lines
3Ctrl+N         Go to next page 3 times
42G             (reserved for page :42)
5n              Jump to next match 5 times
2d              Half-page down twice (= full page)
```

**Implementation**:
```
Buffer holds: "5"
User presses: "j"
Execute: scroll(5)
Buffer clears
```

**Timeout**: 1 second for count entry

---

## Implementation Details

### Keybinding Handler Architecture

```go
type KeybindingHandler struct {
    // Current state
    mode          KeyMode      // Normal, Search, Command, Visual
    countBuffer   string       // "5" in "5j"
    keyBuffer     string       // "g" in "gg"
    lastAction    Action       // for . repeat
    marks         map[rune]Position

    // Timing
    lastKeyTime   time.Time
    countTimeout  time.Duration  // 1 second
    keyTimeout    time.Duration  // 500ms

    // Settings
    caseSensitive bool
    useRegex      bool
    wrapSearch    bool

    // Callbacks
    onAction      func(Action)
}

type Action struct {
    Type     string      // "scroll", "navigate", "search", etc.
    Count    int         // From count buffer
    Argument string      // Pattern for search, theme name, etc.
    Keys     string      // Original keys pressed
}
```

### State Machine

```
┌─────────────────┐
│   NORMAL MODE   │◄──────────┐
│ (default)       │           │
└────────┬────────┘           │
         │                    │
    /    │ :  ?  v  V Ctrl+V  │ Esc
    │    │ │  │  │  │    │    │
    ▼    ▼ ▼  ▼  ▼  ▼    ▼    │
  ┌──┬──┬──┬──┬──┬──┬────┐    │
  │Se│Co│He│Vi│Vi│Bl│ ... │───┘
  │ar│mm│lp│su│lin│ock│    │
  │ch│nd│  │al│ual│    │    │
  └──┴──┴──┴──┴──┴──┴────┘
```

### Keybinding Processing Flow

```
Key Event
    │
    ▼
Check Mode
    │
    ├─ NORMAL: processNormalKey()
    │   ├─ Count Buffer Check (0-9)
    │   ├─ Key Buffer Check (multi-key sequences)
    │   ├─ Keybinding Lookup
    │   └─ Execute Action
    │
    ├─ SEARCH: processSearchKey()
    │   ├─ Typing mode (accumulate query)
    │   ├─ Navigation (n/N)
    │   └─ Exit (Esc)
    │
    ├─ COMMAND: processCommandKey()
    │   ├─ Typing mode (accumulate command)
    │   ├─ Execute (:command)
    │   └─ Exit (Esc)
    │
    └─ VISUAL: processVisualKey()
        ├─ Selection (y, p, d)
        ├─ Movement (j,k,w,b)
        └─ Exit (Esc, v)
```

---

## Keybinding Groups

### Group 1: Essential (Must Have)
- j/k/Up/Down - scroll
- Ctrl+N/P - next/previous page
- gg/G - first/last page
- / - search
- n/N - next/previous match
- q/Ctrl+C - quit
- ? - help

**Test Coverage**: 100%
**Documentation**: Primary help screen
**Keyboard Shortcuts Card**: Yes

### Group 2: Vim Standard (Should Have)
- w/b/e - word navigation
- ^/$ - line start/end
- d/u - half-page
- Ctrl+F/B - full-page
- Tab/Shift+Tab - pane cycling
- Marks (m/'') - position memory
- . - repeat last action

**Test Coverage**: 85%+
**Documentation**: Secondary help (g?)
**Keyboard Shortcuts Card**: Yes

### Group 3: Developer Convenience (Nice to Have)
- :N - jump to page number
- Ctrl+G - show position
- mouse scroll - accessibility
- 1/2 - theme toggle
- * # - word search
- h/l - horizontal scroll

**Test Coverage**: 70%+
**Documentation**: Help menu (g?)
**Keyboard Shortcuts Card**: Optional

### Group 4: Advanced (Future)
- @ - macro replay
- q - macro record
- u/Ctrl+R - undo/redo
- " - registers
- / flags - case/regex options
- :set - configuration

**Test Coverage**: -
**Documentation**: -
**Keyboard Shortcuts Card**: -

---

## Configuration System

### keybindings.toml (Future)

```toml
[bindings]
# Override defaults
scroll_down = ["j", "Down", "J"]          # j, Down arrow, Shift+J all do same thing
scroll_up = ["k", "Up", "K"]
word_next = ["w"]
word_prev = ["b"]

[theme]
dark = "Tokyo Night"
light = "Catppuccin Latte"

[search]
case_sensitive = false
use_regex = true
wrap = true
```

### Per-Mode Bindings
```toml
[normal]
"Ctrl+F" = "page_down"
"Ctrl+B" = "page_up"

[visual]
"y" = "yank"
"p" = "paste"

[command]
"Esc" = "exit_command"
"Ctrl+C" = "exit_command"
```

---

## Help System Design

### Primary Help (? key)
**Size**: 10-12 lines
**Content**: Most common keybindings (Group 1)
**Position**: Centered overlay
**Dismiss**: ? or Esc

```
┌─── LUMOS Keyboard Shortcuts ─────────────────┐
│                                               │
│ Navigation:  j/k scroll  d/u half-page       │
│              gg/G first/last page             │
│              Ctrl+N/P next/prev page          │
│                                               │
│ Search:      / search  n/N next/prev match    │
│                                               │
│ UI:          Tab cycle panes  ? help          │
│              1/2 dark/light theme             │
│                                               │
│ General:     q quit  : command mode           │
│                                               │
│ Press g? for more keybindings · Press ? to close
│
└───────────────────────────────────────────────┘
```

### Extended Help (g? key)
**Size**: Full screen
**Content**: All keybindings (Groups 1-3)
**Organization**: By category
**Features**: Searchable (future), clickable (future)

```
NAVIGATION
  Scrolling (Line)        Scrolling (Page)
  j/k/↑/↓    line         d/u       half-page
  h/l/←/→    horizontal   Ctrl+F/B  full-page
  w/b/e      word         Ctrl+N/P  next/prev
  ^/$/0      line pos     gg/G      first/last

SEARCH
  /pattern   search       */#       word search
  n/N        next/prev    gn        select match

MARKS
  ma/'a      set/jump     ''        last mark
  :marks     list all

VISUAL & COPY
  v/V/Ctrl+V visual       ./paste
  y          yank         p/P       paste

COMMANDS
  :N         page N       :theme    change theme
  :help      help         :q        quit

[Interactive keybindings guide · type /pattern to search · g? to close]
```

---

## Testing Strategy

### Unit Tests (40+)

```go
// Count handling
TestCount_SingleDigit_5j_ScrollsDown5Lines
TestCount_MultiDigit_42G_ReservedForPageJump
TestCount_Timeout_1Second_ClearsBuffer
TestCount_WithMultiKey_5gg_HandlesCorrectly

// Navigation keys
TestKey_w_JumpsToNextWord
TestKey_b_JumpsToPreviousWord
TestKey_e_JumpsToWordEnd
TestKey_HatDollar_LineNavigation
TestKey_Slash_StarHash_SearchVariants

// Marks
TestMark_ma_SetsMarkA
TestMark_QuoteA_JumpsToMarkA
TestMark_QuoteQuote_JumpsToLastMark
TestMark_GraveAccent_JumpsToLastChange

// Commands
TestCommand_ColonN_JumpsToPageN
TestCommand_ColonTheme_SwitchesTheme
TestCommand_ColonHelp_ShowsHelp
TestCommand_ColonMarks_ListsMarks

// Mouse
TestMouse_LeftClick_FocusesPane
TestMouse_Scroll_ScrollsActivePane
TestMouse_DoubleClick_SelectsWord

// Repeat
TestRepeat_Dot_RepeatsPreviousAction
TestRepeat_DotDot_RepeatsTwice
```

### Integration Tests (15+)

```go
// Full workflows
TestWorkflow_SearchWithCount_5nGoesToFifthMatch
TestWorkflow_MarkAndReturn_maJumpBackWithQuoteA
TestWorkflow_PageNavigateWithCount_5CtrlN_FivePages
TestWorkflow_VisualSelectAndCopy_vAwky_SelectAndCopy
TestWorkflow_CommandModeFlow_ColonThemeLightEnter

// Edge cases
TestEdgeCase_CountTimeoutAfter1Second
TestEdgeCase_MultiKeyTimeoutAfter500ms
TestEdgeCase_MarkOnLastPage
TestEdgeCase_SearchWrapAround
TestEdgeCase_CommandWithInvalidPageNumber
```

### Manual Test Scenarios

1. **Vim User Comfort Test**
   - Use only j/k/w/b/e/^/$ for movement
   - Search with / and n/N
   - Should feel like vim

2. **Developer Workflow Test**
   - Use :42 to jump to page 42
   - Use marks to jump between sections (ma, 'a)
   - Use . to repeat searches

3. **Discovery Test**
   - ? shows essential keys
   - g? shows all keys
   - :help explains features

4. **Accessibility Test**
   - All features work with mouse too
   - Scroll wheel works
   - Click to focus panes

---

## Success Criteria

### Completeness
- ✅ All Group 1 keybindings implemented (7 keys)
- ✅ All Group 2 keybindings implemented (15 keys)
- ✅ Most Group 3 keybindings implemented (8 keys)
- ✅ Count support (5j, 10k, 3n, etc.)
- ✅ Multi-key support (gg, '', etc.)
- ✅ Mouse support (scroll, click)
- ✅ Mark system (ma, 'a)
- ✅ Repeat functionality (.)

### Quality
- ✅ 40+ unit tests (all passing)
- ✅ 15+ integration tests (all passing)
- ✅ >90% coverage for keybindings
- ✅ No key conflicts or ambiguities
- ✅ Vim user validation (feels right)

### Performance
- ✅ Key response <16ms (60 FPS)
- ✅ Count timeout 1 second
- ✅ Multi-key timeout 500ms
- ✅ No input lag or stuttering

### Documentation
- ✅ Primary help (?) - 10 lines, essential keys
- ✅ Extended help (g?) - full screen, all keys
- ✅ README.md - keybinding reference
- ✅ Keybindings card - printable reference
- ✅ Inline help text (: for hints)

---

## Migration Path

### Phase 1 (Current - 1.4 Complete)
- Basic vim navigation (j/k/d/u/gg/G)
- Basic search (/)
- Basic UI (Tab, ?, q)
- Status: ✅ Done

### Phase 1.5 (Next)
- Multi-key sequences (gg working properly)
- Vim counts (5j, 3n, etc.)
- More search options (*, #)
- Extended help (g?)
- Target: 2-3 hours

### Phase 1.5+ (Post 1.6)
- Word navigation (w/b/e)
- Line navigation (^/$)
- Marks and jumps (ma, 'a)
- Mouse support
- Repeat functionality (.)
- Command mode enhancements (:theme, :marks, etc.)
- Target: 4-6 hours

### Phase 2 (Enhanced)
- Configurable keybindings (keybindings.toml)
- Macros (q, @)
- Undo/redo (u, Ctrl+R)
- Register system (", a-z)
- Visual modes (v, V, Ctrl+V)
- Target: Phase 2 (Dec 2025)

---

## Comparison: Basic vs Comprehensive

| Feature | Basic (Current) | Comprehensive (Enhanced) |
|---------|-----------------|--------------------------|
| Scroll keys | j/k/d/u | j/k/d/u/Ctrl+F/Ctrl+B/Ctrl+E/Ctrl+Y |
| Navigation | 4 basic keys | 15+ keys with word/line/mark support |
| Repeat support | None | Full count support (5j, 10k) |
| Search | / n N | / ? * # n N (forward/backward) |
| Commands | 6 basic | 20+ commands (:theme, :marks, :help) |
| Marks | None | Complete mark system (ma, 'a, '') |
| Repeat | None | . for repeat, macro support |
| Mouse | None | Full mouse support (scroll, click) |
| Help | Basic (?) | Primary (?) + Extended (g?) + Searchable |
| Configuration | None | Full keybindings.toml support |
| **Total Keys** | **~20** | **~50+** |

---

## Dependencies

- ✅ Milestone 1.4 (TUI framework)
- ✅ pkg/ui/keybindings.go (exists, needs enhancement)
- ✅ Bubble Tea viewport (scrolling)
- ✅ pkg/pdf/search.go (search patterns)

---

## Known Limitations

1. **Macro System (Future)**: q/@ require more state management
2. **Undo/Redo (Future)**: Requires action history tracking
3. **Regex Search (Future)**: More complex pattern matching
4. **Configuration (Future)**: Requires config file parsing
5. **Visual Modes (Future)**: Text selection requires viewport changes

---

## Timeline

**Quick Enhancement (Milestone 1.5 focus)**: 2-3 hours
- Counts (5j, 10k, 3n)
- Multi-key (gg)
- Extended help (g?)

**Full System (Post 1.6)**: 4-6 hours
- Word/line navigation
- Marks and jumps
- Mouse support
- Repeat functionality
- Command mode enhancements

**Advanced (Phase 2)**: 8-10 hours
- Configuration system
- Macros
- Undo/redo
- Visual modes
- Registers

---

## References

- Vim Documentation: https://vim.rtfm.io/
- Vim Cheat Sheet: https://vim.rtfm.io/cheatsheet/
- Bubble Tea Key Events: https://github.com/charmbracelet/bubbletea/wiki
- Terminal Input Handling: https://en.wikipedia.org/wiki/ANSI_escape_code

---

**Specification Version**: 1.0
**Last Updated**: 2025-11-18
**Status**: Ready for Implementation
**Priority**: P0 (Core Feature)
