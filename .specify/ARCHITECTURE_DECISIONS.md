# LUMOS Architecture Decisions

**Version**: 1.0.0
**Last Updated**: 2025-11-17
**Format**: Architecture Decision Records (ADRs)

---

## ADR-001: Bubble Tea for Terminal User Interface

### Status
**Accepted** (2025-10-21)

### Context
LUMOS needs a robust TUI framework that can handle complex interactions, maintain state predictably, and provide a smooth user experience.

### Decision
Use Bubble Tea (github.com/charmbracelet/bubbletea) as the TUI framework.

### Rationale
- **MVU Pattern**: Model-View-Update provides predictable state management
- **Tea.Cmd**: Elegant handling of async operations
- **Ecosystem**: Rich component library (bubbles) available
- **Community**: Active development and support
- **Go Native**: Pure Go, no CGO requirements

### Consequences
- **Positive**: Clean architecture, testable code, maintainable state
- **Negative**: Learning curve for MVU pattern
- **Mitigation**: Study LUMINA's implementation for patterns

### Alternatives Rejected
- **tcell/tview**: More imperative, harder state management
- **termui**: Less active development
- **Raw terminal**: Too much complexity for timeline

---

## ADR-002: LRU Cache for PDF Pages

### Status
**Implemented** (2025-11-01)

### Context
PDF files can be large, and extracting text from pages is expensive due to the ledongthuc/pdf library requiring file reopening for each page.

### Decision
Implement a 5-page LRU (Least Recently Used) cache with thread-safe access.

### Rationale
- **Performance**: Cache hits <100ns vs ~50ms for extraction
- **Memory**: 5 pages ≈ 50KB typical, acceptable overhead
- **Simplicity**: Map-based implementation, no external dependencies
- **Configurable**: Can adjust size based on memory constraints

### Implementation
```go
type LRUCache struct {
    capacity int
    cache    map[int]string
    order    []int
    mu       sync.RWMutex
}
```

### Consequences
- **Positive**: 1000x performance improvement for cached pages
- **Positive**: Smooth navigation experience
- **Negative**: Additional memory usage
- **Mitigation**: Monitor memory, make size configurable

### Performance Metrics
- Cache hit: 16ns/op
- Cache miss: 8ns/op
- Put operation: 61ns/op

---

## ADR-003: Ledongthuc/pdf for PDF Parsing

### Status
**Accepted** (2025-10-21)

### Context
Need a PDF parsing library that is pure Go, well-maintained, and can extract text reliably.

### Decision
Use github.com/ledongthuc/pdf for PDF parsing.

### Rationale
- **Pure Go**: No CGO dependencies, easy cross-compilation
- **Maintained**: Active development, responsive maintainer
- **Sufficient**: Handles text extraction well
- **Simple API**: Clean, idiomatic Go interface

### Constraints
- **File Reopening**: Must reopen file for each page (library limitation)
- **Text Only**: No image extraction in current version
- **Layout**: Limited support for complex layouts

### Mitigations
- **Cache**: LRU cache eliminates reopening impact
- **Phase 3**: Plan image support with go-termimg
- **Validation**: Test with diverse PDFs

### Alternatives Rejected
- **pdfcpu**: More complex, overkill for text extraction
- **UniPDF**: Commercial license required
- **CGO libs**: Deployment complexity

---

## ADR-004: Three-Pane Layout Design

### Status
**Accepted** (2025-10-21)

### Context
Users need to see document metadata, content, and search results simultaneously without mode switching.

### Decision
Implement a three-pane vertical layout: Metadata (20%) | Viewer (60%) | Search (20%).

### Rationale
- **Familiar**: IDE-like layout developers know
- **Efficient**: All information visible at once
- **Responsive**: Percentages adapt to terminal width
- **Balanced**: 60% for content is optimal for reading

### Layout Calculation
```go
func calculatePaneWidths(totalWidth int) (meta, viewer, search int) {
    meta = totalWidth * 20 / 100
    search = totalWidth * 20 / 100
    viewer = totalWidth - meta - search
    return
}
```

### Consequences
- **Positive**: Information density, no context switching
- **Positive**: Predictable layout users can learn
- **Negative**: Minimum 80 columns required
- **Mitigation**: Graceful degradation for narrow terminals

### Future Enhancement
- Phase 2: Resizable panes
- Phase 2: Collapsible panels
- Phase 3: Customizable layouts

---

## ADR-005: Vim Keybindings as Primary Interface

### Status
**Accepted** (2025-10-21)

### Context
Target audience is developers who likely use vim/neovim and expect familiar navigation.

### Decision
Implement vim keybindings as the primary (and only) navigation method.

### Keybinding Matrix
```
Navigation: j/k (line), d/u (half page), gg/G (document)
Paging: Ctrl+N/Ctrl+P
Search: / (open), n/N (next/prev)
UI: Tab (cycle panes), ? (help), q (quit)
```

### Rationale
- **Muscle Memory**: Developers already know these
- **Efficient**: Minimal hand movement
- **Consistent**: Same as vim, less, man pages
- **Complete**: Covers all navigation needs

### Consequences
- **Positive**: Immediate productivity for vim users
- **Positive**: No mouse requirement
- **Negative**: Learning curve for non-vim users
- **Mitigation**: Help overlay with all bindings

### Future Considerations
- Phase 2: Customizable keybindings
- Phase 2: Emacs mode option
- Phase 3: Mouse support addition

---

## ADR-006: Dark Mode by Default

### Status
**Accepted** (2025-10-21)

### Context
Developers spend long hours reading documentation and papers. Eye strain is a real concern.

### Decision
Dark mode is the default theme with high contrast (>7:1) ratios.

### Color Palette
```go
Background: #1a1b26 (Tokyo Night)
Foreground: #c0caf5
Accent:     #7aa2f7
Border:     #3b4261
Error:      #f7768e
```

### Rationale
- **Health**: Reduces eye strain in long sessions
- **Modern**: Matches VS Code, terminal themes
- **Performance**: Less light emission on OLED
- **Professional**: Clean, focused appearance

### Consequences
- **Positive**: Better for extended use
- **Positive**: Matches developer environments
- **Negative**: Some prefer light mode
- **Mitigation**: Toggle with key binding (1/2)

---

## ADR-007: Test-Driven Development

### Status
**Enforced** (2025-11-01)

### Context
Complex state management and user interactions require robust testing.

### Decision
Follow strict TDD: Write tests first, implementation second.

### Testing Strategy
```
1. Unit tests for all packages (target: 90%)
2. Integration tests for user flows
3. Benchmark tests for performance
4. Property tests for edge cases
```

### Current Metrics
- **Coverage**: 94.4% (exceeds 80% requirement)
- **Tests**: 42 passing
- **Benchmarks**: 9 performance tests

### Consequences
- **Positive**: High confidence in code
- **Positive**: Refactoring safety
- **Positive**: Living documentation
- **Negative**: Slower initial development
- **Mitigation**: Time saved in debugging

---

## ADR-008: Single Binary Distribution

### Status
**Planned** (Phase 1)

### Context
Users need simple installation without dependency management.

### Decision
Distribute LUMOS as a single static binary.

### Build Configuration
```makefile
CGO_ENABLED=0 go build -ldflags="-s -w"
```

### Rationale
- **Simple**: Copy binary, done
- **Portable**: Works anywhere Go runs
- **Small**: ~5MB compressed
- **Secure**: No dynamic dependencies

### Distribution Channels
- GitHub Releases (primary)
- Homebrew tap (Phase 2)
- AUR package (Phase 2)

---

## ADR-009: Configuration via Code

### Status
**Accepted** (2025-10-21)

### Context
Phase 1 needs to be simple. Configuration files add complexity.

### Decision
All configuration is in code for Phase 1. No config files.

### Defaults
```go
const (
    CacheSize = 5
    Theme = "dark"
    ViewportBuffer = 100
)
```

### Rationale
- **Simple**: One less thing to parse/validate
- **Fast**: No startup file I/O
- **Predictable**: Same behavior everywhere
- **Sufficient**: Phase 1 doesn't need customization

### Future (Phase 2)
- XDG config support
- TOML/YAML configuration
- Per-directory settings

---

## ADR-010: Performance Budgets

### Status
**Enforced** (Constitutional)

### Context
Performance is a feature. Slow readers won't be used.

### Decision
Strict performance budgets enforced by benchmarks.

### Budgets
| Operation | Budget | Current | Status |
|-----------|--------|---------|--------|
| Startup | <100ms | ~70ms | ✅ |
| Page Load (cached) | <50ms | <20ms | ✅ |
| Page Load (uncached) | <200ms | ~65ms | ✅ |
| Search | <100ms | <50μs | ✅ |
| Memory (10MB PDF) | <50MB | ~8MB | ✅ |

### Enforcement
```bash
make bench  # Must pass before merge
```

### Consequences
- **Positive**: Consistently fast experience
- **Positive**: Performance regressions caught
- **Negative**: Optimization complexity
- **Mitigation**: Profile-guided optimization

---

## ADR-011: Package Boundaries

### Status
**Enforced** (Constitutional)

### Context
Clear separation of concerns prevents coupling and aids testing.

### Decision
Strict package boundaries with no circular dependencies.

### Package Responsibilities
```
cmd/lumos/   → Entry point only
pkg/pdf/     → PDF operations (no UI knowledge)
pkg/ui/      → TUI components (no PDF internals)
pkg/config/  → Pure data structures
```

### Enforcement
- Import restrictions in code review
- `go mod graph` for dependency analysis
- Test isolation verification

### Consequences
- **Positive**: Clear responsibilities
- **Positive**: Easy testing
- **Positive**: Parallel development
- **Negative**: Some code duplication
- **Mitigation**: Shared interfaces

---

## ADR-012: Error Handling Philosophy

### Status
**Accepted** (2025-11-01)

### Context
PDF parsing can fail in many ways. Users need clear feedback.

### Decision
Never panic in library code. Return errors with context.

### Pattern
```go
if err != nil {
    return fmt.Errorf("loading page %d: %w", pageNum, err)
}
```

### User-Facing Errors
```
"Could not open PDF: file not found"
"This PDF appears to be encrypted"
"Page 5 could not be rendered"
```

### Consequences
- **Positive**: Graceful degradation
- **Positive**: Debugging context
- **Positive**: User understanding
- **Negative**: Verbose error handling
- **Mitigation**: Helper functions

---

## Decision Log

| Date | Decision | Impact |
|------|----------|--------|
| 2025-10-21 | Choose Bubble Tea | High - Entire UI architecture |
| 2025-10-21 | 3-pane layout | Medium - User experience |
| 2025-10-21 | Vim keybindings | High - User interaction model |
| 2025-11-01 | LRU cache implementation | High - Performance achieved |
| 2025-11-01 | TDD approach | High - 94.4% coverage |
| 2025-11-17 | Spec-driven development | High - Clear roadmap |

---

## Future Decision Points

### Phase 2 Decisions
- Configuration file format (TOML vs YAML)
- Plugin architecture (Lua vs WASM)
- Bookmark storage (SQLite vs JSON)
- TOC extraction strategy

### Phase 3 Decisions
- Image rendering library (go-termimg vs sixel)
- Caching strategy for images
- OCR integration approach

### Phase 4 Decisions
- AI provider (Claude vs OpenAI vs Local)
- Context window management
- Privacy considerations

---

**Note**: All decisions must comply with `.specify/constitution.md` or document justified exceptions.