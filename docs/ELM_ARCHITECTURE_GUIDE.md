# Elm Architecture for LUMOS TUI

**Last Updated**: 2025-11-01  
**Context**: Milestone 1.4 - Basic TUI Framework

---

## Table of Contents

1. [What is the Elm Architecture?](#what-is-the-elm-architecture)
2. [Core Concepts](#core-concepts)
3. [Bubble Tea Implementation](#bubble-tea-implementation)
4. [LUMOS Application Design](#lumos-application-design)
5. [Message Flow](#message-flow)
6. [State Management](#state-management)
7. [Testing Strategy](#testing-strategy)
8. [Best Practices](#best-practices)

---

## What is the Elm Architecture?

The Elm Architecture (TEA) is a pattern for building interactive applications with three core components:

1. **Model** - Application state
2. **Update** - State transformation logic
3. **View** - State rendering

### Key Principles

- ✅ **Unidirectional data flow**: Data flows in one direction (Model → View → Message → Update → Model)
- ✅ **Pure functions**: Update logic is deterministic and side-effect free
- ✅ **Explicit state**: All state is in the Model, no hidden state
- ✅ **Message-driven**: All changes happen through explicit messages

### The Update Loop

```
   ┌─────────────────────────────────────┐
   │                                     │
   │  ┌─────────┐    ┌────────┐    ┌────▼────┐
   │  │  Model  │───▶│  View  │───▶│ Message │
   │  └─────────┘    └────────┘    └────┬────┘
   │                                     │
   │  ┌─────────┐                        │
   └──│ Update  │◀───────────────────────┘
      └─────────┘
```

**Flow**:
1. **Model** holds current state
2. **View** renders state as HTML/TUI
3. User interaction generates **Message**
4. **Update** transforms Model based on Message
5. New Model triggers View re-render
6. Loop continues

---

## Core Concepts

### 1. Model (State)

The Model represents **all** application state. It should be:

- **Complete**: Contains everything needed to render the view
- **Serializable**: Can be saved/loaded if needed
- **Immutable**: Updated through pure functions, not mutated

**Example - Simple Counter**:
```go
type Model struct {
    count int
}

func initialModel() Model {
    return Model{count: 0}
}
```

**Example - Text Input**:
```go
type Model struct {
    content string
}

func initialModel() Model {
    return Model{content: ""}
}
```

### 2. Messages (Events)

Messages represent **all possible events** that can happen in the application.

- **Explicit**: Every user action is a message type
- **Type-safe**: Using Go types ensures correctness
- **Self-documenting**: Message names describe intent

**Example Messages**:
```go
type Msg interface{}

type IncrementMsg struct{}
type DecrementMsg struct{}
type ChangeMsg struct {
    newContent string
}
```

### 3. Update (State Transitions)

The Update function transforms state based on messages.

**Signature**: `Update(msg Msg, model Model) (Model, tea.Cmd)`

- **Input**: Current message and model
- **Output**: New model and optional command
- **Pure**: Same inputs always produce same outputs
- **Exhaustive**: Handles all message types

**Example**:
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case IncrementMsg:
        m.count++
        return m, nil
    
    case DecrementMsg:
        m.count--
        return m, nil
    
    case ChangeMsg:
        m.content = msg.newContent
        return m, nil
    }
    
    return m, nil
}
```

### 4. View (Rendering)

The View function renders the current state.

**Signature**: `View(model Model) string`

- **Input**: Current model
- **Output**: String representation (TUI)
- **Pure**: Same model always produces same view
- **No side effects**: Only renders, doesn't change state

**Example**:
```go
func (m Model) View() string {
    return fmt.Sprintf("Count: %d\n", m.count)
}
```

---

## Bubble Tea Implementation

Bubble Tea is Go's implementation of the Elm Architecture for terminal UIs.

### The `tea.Model` Interface

```go
type Model interface {
    Init() Cmd
    Update(Msg) (Model, Cmd)
    View() string
}
```

### Key Differences from Elm

| Concept | Elm | Bubble Tea |
|---------|-----|------------|
| Model | `type alias Model = ...` | `type Model struct { ... }` |
| Messages | `type Msg = ...` | `type Msg interface{}` with variants |
| Update | Returns `(Model, Cmd Msg)` | Returns `(tea.Model, tea.Cmd)` |
| View | Returns `Html Msg` | Returns `string` |
| Commands | `Cmd Msg` | `tea.Cmd` |
| Subscriptions | Built-in | Via `tea.Sub` |

### Commands (`tea.Cmd`)

Commands are **side effects** that return messages.

**Common commands**:
- `tea.Tick` - Timer events
- `tea.Batch` - Multiple commands
- Custom commands - HTTP requests, file I/O, etc.

**Example**:
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" {
            return m, tea.Quit
        }
    }
    return m, nil
}
```

---

## LUMOS Application Design

### Model Structure

```go
type Model struct {
    // Document state
    doc         *pdf.Document
    currentPage int
    totalPages  int
    
    // Display state
    viewport    viewport.Model
    statusBar   string
    
    // Search state
    searchMode  bool
    searchQuery string
    searchResults []pdf.SearchResult
    
    // UI state
    width       int
    height      int
    ready       bool
    
    // Error state
    err         error
}
```

**Design principles**:
- ✅ All state in one place
- ✅ Grouped by concern (document, display, search, UI, errors)
- ✅ Sufficient to render any view
- ✅ Can be initialized with defaults

### Message Types

```go
// Document messages
type LoadDocumentMsg struct {
    path string
}

type DocumentLoadedMsg struct {
    doc *pdf.Document
}

type PageChangedMsg struct {
    newPage int
}

// Navigation messages
type ScrollUpMsg struct{}
type ScrollDownMsg struct{}
type NextPageMsg struct{}
type PreviousPageMsg struct{}
type FirstPageMsg struct{}
type LastPageMsg struct{}

// Search messages
type EnterSearchModeMsg struct{}
type ExitSearchModeMsg struct{}
type SearchQueryChangedMsg struct {
    query string
}
type SearchExecutedMsg struct {
    results []pdf.SearchResult
}
type NextSearchResultMsg struct{}
type PreviousSearchResultMsg struct{}

// UI messages
type WindowSizeMsg tea.WindowSizeMsg
type QuitMsg struct{}
type ErrorMsg struct {
    err error
}
```

### Update Logic

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    
    // Keyboard input
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    
    // Window resize
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.viewport = viewport.New(msg.Width, msg.Height-2)
        return m, nil
    
    // Document loaded
    case DocumentLoadedMsg:
        m.doc = msg.doc
        m.totalPages = msg.doc.GetPageCount()
        m.currentPage = 1
        return m, m.loadPage(1)
    
    // Page changed
    case PageChangedMsg:
        m.currentPage = msg.newPage
        return m, m.loadPage(msg.newPage)
    
    // Search executed
    case SearchExecutedMsg:
        m.searchResults = msg.results
        return m, nil
    
    // Error occurred
    case ErrorMsg:
        m.err = msg.err
        return m, nil
    }
    
    return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // Search mode takes precedence
    if m.searchMode {
        return m.handleSearchKeyPress(msg)
    }
    
    // Normal mode navigation
    switch msg.String() {
    case "q", "ctrl+c":
        return m, tea.Quit
    
    case "j", "down":
        return m, ScrollUpMsg{}
    
    case "k", "up":
        return m, ScrollDownMsg{}
    
    case "n":
        if m.currentPage < m.totalPages {
            return m, func() tea.Msg {
                return PageChangedMsg{newPage: m.currentPage + 1}
            }
        }
    
    case "p":
        if m.currentPage > 1 {
            return m, func() tea.Msg {
                return PageChangedMsg{newPage: m.currentPage - 1}
            }
        }
    
    case "g":
        return m, func() tea.Msg {
            return PageChangedMsg{newPage: 1}
        }
    
    case "G":
        return m, func() tea.Msg {
            return PageChangedMsg{newPage: m.totalPages}
        }
    
    case "/":
        m.searchMode = true
        m.searchQuery = ""
        return m, nil
    }
    
    return m, nil
}
```

### View Rendering

```go
func (m Model) View() string {
    if !m.ready {
        return "Loading..."
    }
    
    if m.err != nil {
        return fmt.Sprintf("Error: %v\n", m.err)
    }
    
    var b strings.Builder
    
    // Render viewport (main content area)
    b.WriteString(m.viewport.View())
    b.WriteString("\n")
    
    // Render status bar
    b.WriteString(m.renderStatusBar())
    
    // Render search bar if in search mode
    if m.searchMode {
        b.WriteString("\n")
        b.WriteString(m.renderSearchBar())
    }
    
    return b.String()
}

func (m Model) renderStatusBar() string {
    left := fmt.Sprintf(" Page %d/%d", m.currentPage, m.totalPages)
    right := " q: quit | /: search | j/k: scroll | n/p: next/prev page "
    
    padding := m.width - len(left) - len(right)
    if padding < 0 {
        padding = 0
    }
    
    return lipgloss.NewStyle().
        Background(lipgloss.Color("240")).
        Foreground(lipgloss.Color("255")).
        Render(left + strings.Repeat(" ", padding) + right)
}
```

---

## Message Flow

### Example: User Presses 'n' Key

1. **Input**: User presses 'n'
2. **Message**: Bubble Tea sends `tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}`
3. **Update**: 
   ```go
   case "n":
       if m.currentPage < m.totalPages {
           return m, func() tea.Msg {
               return PageChangedMsg{newPage: m.currentPage + 1}
           }
       }
   ```
4. **Command**: Returns command that generates `PageChangedMsg`
5. **Update (again)**:
   ```go
   case PageChangedMsg:
       m.currentPage = msg.newPage
       return m, m.loadPage(msg.newPage)
   ```
6. **Command**: `loadPage` command fetches page content
7. **Message**: Returns `PageLoadedMsg` when complete
8. **Update (final)**:
   ```go
   case PageLoadedMsg:
       m.viewport.SetContent(msg.content)
       return m, nil
   ```
9. **View**: Re-renders with new page content

### Command Pattern

Commands are functions that return messages asynchronously:

```go
func (m Model) loadPage(pageNum int) tea.Cmd {
    return func() tea.Msg {
        page, err := m.doc.GetPage(pageNum)
        if err != nil {
            return ErrorMsg{err: err}
        }
        return PageLoadedMsg{
            pageNum: pageNum,
            content: page.Text,
        }
    }
}
```

**Key insights**:
- Commands are **side effects**
- They return **messages**
- Messages feed back into **Update**
- This keeps Update pure while allowing I/O

---

## State Management

### Derived State

Some state can be **computed** from other state. Don't store it!

**Bad** ❌:
```go
type Model struct {
    currentPage int
    totalPages  int
    isFirstPage bool  // Derived!
    isLastPage  bool  // Derived!
}
```

**Good** ✅:
```go
type Model struct {
    currentPage int
    totalPages  int
}

func (m Model) isFirstPage() bool {
    return m.currentPage == 1
}

func (m Model) isLastPage() bool {
    return m.currentPage == m.totalPages
}
```

### State Initialization

```go
func initialModel(pdfPath string) Model {
    return Model{
        currentPage: 1,
        searchMode:  false,
        searchQuery: "",
        ready:       false,
        viewport:    viewport.New(80, 24),
    }
}

func (m Model) Init() tea.Cmd {
    return func() tea.Msg {
        return LoadDocumentMsg{path: m.pdfPath}
    }
}
```

### State Validation

Ensure state is always valid:

```go
func (m Model) setPage(page int) Model {
    if page < 1 {
        page = 1
    }
    if page > m.totalPages {
        page = m.totalPages
    }
    m.currentPage = page
    return m
}
```

---

## Testing Strategy

### Testing Update Logic

Update functions are **pure**, making them easy to test:

```go
func TestUpdate_NextPage(t *testing.T) {
    model := Model{
        currentPage: 1,
        totalPages:  5,
    }
    
    newModel, _ := model.Update(tea.KeyMsg{
        Type:  tea.KeyRunes,
        Runes: []rune{'n'},
    })
    
    m := newModel.(Model)
    
    if m.currentPage != 2 {
        t.Errorf("Expected page 2, got %d", m.currentPage)
    }
}
```

### Testing View Rendering

Views are also **pure**:

```go
func TestView_StatusBar(t *testing.T) {
    model := Model{
        currentPage: 3,
        totalPages:  10,
        width:       80,
        ready:       true,
    }
    
    output := model.View()
    
    if !strings.Contains(output, "Page 3/10") {
        t.Error("Status bar missing page info")
    }
}
```

### Testing Commands

Commands return messages, which can be tested:

```go
func TestLoadPage_Success(t *testing.T) {
    doc := createTestDocument(t)
    model := Model{doc: doc}
    
    cmd := model.loadPage(1)
    msg := cmd()  // Execute command
    
    switch msg := msg.(type) {
    case PageLoadedMsg:
        if msg.pageNum != 1 {
            t.Errorf("Expected page 1, got %d", msg.pageNum)
        }
    case ErrorMsg:
        t.Errorf("Unexpected error: %v", msg.err)
    default:
        t.Errorf("Unexpected message type: %T", msg)
    }
}
```

---

## Best Practices

### 1. Keep Update Pure

**Bad** ❌:
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case LoadDocumentMsg:
        // Side effect in Update!
        doc, err := pdf.NewDocument(msg.path, 10)
        if err != nil {
            m.err = err
        }
        m.doc = doc
        return m, nil
}
```

**Good** ✅:
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case LoadDocumentMsg:
        // Return command for side effect
        return m, m.loadDocument(msg.path)
}

func (m Model) loadDocument(path string) tea.Cmd {
    return func() tea.Msg {
        doc, err := pdf.NewDocument(path, 10)
        if err != nil {
            return ErrorMsg{err: err}
        }
        return DocumentLoadedMsg{doc: doc}
    }
}
```

### 2. Use Message Types for Everything

**Bad** ❌:
```go
// Using string constants
const (
    MSG_NEXT_PAGE = "next"
    MSG_PREV_PAGE = "prev"
)
```

**Good** ✅:
```go
// Using typed messages
type NextPageMsg struct{}
type PreviousPageMsg struct{}
```

### 3. Model Contains All State

**Bad** ❌:
```go
var globalSearchResults []SearchResult  // Global state!

type Model struct {
    currentPage int
}
```

**Good** ✅:
```go
type Model struct {
    currentPage   int
    searchResults []SearchResult  // In model
}
```

### 4. Commands for Async Operations

**Bad** ❌:
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case SearchMsg:
        // Blocking operation!
        results := m.doc.Search(msg.query)
        m.searchResults = results
        return m, nil
}
```

**Good** ✅:
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case SearchMsg:
        return m, m.executeSearch(msg.query)
}

func (m Model) executeSearch(query string) tea.Cmd {
    return func() tea.Msg {
        results := m.doc.Search(query)
        return SearchResultsMsg{results: results}
    }
}
```

### 5. Handle All Messages

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case WindowSizeMsg:
        return m.handleResize(msg)
    // ... handle all message types
    default:
        // Always have a default case
        return m, nil
    }
}
```

---

## Summary

The Elm Architecture provides:

✅ **Predictable state management** - All state in Model  
✅ **Easy testing** - Pure functions  
✅ **Clear data flow** - Unidirectional  
✅ **Type safety** - Compile-time guarantees  
✅ **Debuggability** - All state changes are explicit messages  

For LUMOS, this means:
- PDF document state is always clear
- UI interactions are explicit messages
- Testing is straightforward
- Adding features is systematic
- Debugging is easier

**Next Steps**: Implement Milestone 1.4 using these patterns!

---

## References

- [Elm Architecture Guide](https://guide.elm-lang.org/architecture/)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Bubble Tea Examples](https://github.com/charmbracelet/bubbletea/tree/master/examples)
