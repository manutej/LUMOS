# LUMOS Constitutional Framework

**Version**: 1.0.0
**Last Updated**: 2025-11-10
**Status**: Active

This document establishes the immutable architectural principles and quality standards for LUMOS development. All implementations MUST comply with these constraints.

---

## I. Code Quality Principles

### 1.1 Go Idioms & Style
- **MUST** follow official Go style guidelines
- **MUST** pass `go fmt` without changes
- **MUST** pass `go vet` without warnings
- **MUST** pass `golangci-lint run` with zero issues
- **MUST** use clear, descriptive variable names (no single-letter except loop indices)

### 1.2 Code Organization
- **MUST** organize by package concern: `pdf/`, `ui/`, `config/`
- **MUST** keep functions under 50 lines (complex logic should be extracted)
- **MUST** avoid circular dependencies between packages
- **MUST** comment all exported functions, types, and constants

### 1.3 Error Handling
- **MUST** return errors, never panic in library code
- **MUST** wrap errors with context using `fmt.Errorf("context: %w", err)`
- **MUST** handle all error returns (no `_ = err`)
- **MUST** validate all inputs at package boundaries

---

## II. Testing Standards

### 2.1 Test Coverage
- **MUST** maintain minimum 80% overall test coverage
- **MUST** achieve 90%+ coverage for critical packages (`pdf/`, `ui/`)
- **MUST** test all exported functions
- **MUST** include edge cases and error paths

### 2.2 Test-Driven Development
- **MUST** write tests before implementation (TDD approach)
- **MUST** use table-driven tests for multiple scenarios
- **MUST** name tests clearly: `TestFunctionName_Scenario_ExpectedBehavior`
- **MUST** include both unit and integration tests

### 2.3 Test Structure
```go
// Example test structure
func TestCacheGet_ExistingKey_ReturnsValue(t *testing.T) {
    // Given: Setup test fixtures
    cache := NewLRUCache(5)
    cache.Put(1, "test content")

    // When: Perform action
    result, ok := cache.Get(1)

    // Then: Verify expectations
    if !ok {
        t.Error("Expected key to exist")
    }
    if result != "test content" {
        t.Errorf("Expected 'test content', got %q", result)
    }
}
```

### 2.4 Benchmarking
- **MUST** include benchmarks for performance-critical operations
- **MUST** run benchmarks in CI before merging
- **MUST** document performance regressions

---

## III. Performance Budgets

### 3.1 Startup Performance
- Cold start: **MUST** be < 100ms
- PDF open: **MUST** be < 50ms for typical files
- First page render: **MUST** be < 70ms total

### 3.2 Runtime Performance
- Page switch (cached): **MUST** be < 50ms
- Page switch (uncached): **MUST** be < 200ms
- Search operation: **MUST** be < 100ms per 100 pages
- Viewport scroll: **MUST** maintain 60 FPS

### 3.3 Memory Constraints
- Baseline memory: **MUST** be < 10MB without PDF loaded
- 10MB PDF loaded: **MUST** use < 50MB total
- 100MB PDF loaded: **MUST** use < 200MB total
- LRU cache: **MUST** respect configured limits

### 3.4 Validation
```bash
# Performance tests MUST pass
go test -bench=. ./pkg/...

# Benchmarks MUST show:
# - Cache operations: < 100ns/op
# - Search operations: < 50μs/op per KB
# - Memory allocations: minimize in hot paths
```

---

## IV. Architectural Constraints

### 4.1 MVU Pattern (Bubble Tea)
- **MUST** follow Model-View-Update pattern strictly
- **MUST** keep all state in Model struct
- **MUST** use messages for all state changes
- **MUST NOT** mutate state outside Update() function

### 4.2 Package Separation
```
cmd/lumos/          # CLI entry - thin wrapper
pkg/pdf/            # PDF operations - no UI dependencies
pkg/ui/             # TUI components - no PDF internals
pkg/config/         # Configuration - pure data
```

- **MUST NOT** import `pkg/ui` from `pkg/pdf`
- **MUST NOT** import `pkg/pdf` internals from `cmd/lumos`
- **MUST** use interfaces for testability

### 4.3 Concurrency
- **MUST** use `sync.RWMutex` for shared data structures
- **MUST** document all goroutine lifecycle
- **MUST** avoid goroutine leaks (proper cleanup)
- **MUST** pass race detector: `go test -race ./...`

### 4.4 Dependencies
- **MUST** use only approved dependencies:
  - `github.com/charmbracelet/bubbletea` (TUI framework)
  - `github.com/charmbracelet/bubbles` (UI components)
  - `github.com/charmbracelet/lipgloss` (styling)
  - `github.com/ledongthuc/pdf` (PDF parsing)
- **MUST** justify any new dependency with benchmarks
- **MUST** vendor dependencies for reproducible builds

---

## V. User Experience Standards

### 5.1 Vim Keybindings
- **MUST** support standard vim navigation: j/k/d/u/gg/G
- **MUST** provide visual feedback for all actions
- **MUST** show help with `?` key
- **MUST** allow quit with `q` or Ctrl+C

### 5.2 Dark Mode Default
- **MUST** enable dark mode by default
- **MUST** provide high contrast text (>7:1 ratio)
- **MUST** support light mode toggle
- **MUST** persist theme preference (Phase 2+)

### 5.3 Responsiveness
- **MUST** respond to keypresses within 16ms (60 FPS)
- **MUST** show loading indicators for operations >100ms
- **MUST** handle terminal resize gracefully
- **MUST** never block UI thread

### 5.4 Error Messages
- **MUST** show user-friendly error messages
- **MUST** suggest corrective actions
- **MUST** log technical details for debugging
- **SHOULD** recover from non-fatal errors

---

## VI. Documentation Requirements

### 6.1 Code Documentation
- **MUST** document all exported symbols with godoc comments
- **MUST** include examples for complex functions
- **MUST** explain non-obvious algorithms
- **MUST** document performance characteristics

### 6.2 User Documentation
- **MUST** update README.md for user-facing changes
- **MUST** maintain accurate PROGRESS.md
- **MUST** create milestone review documents
- **MUST** update ROADMAP.md with completion status

### 6.3 Development Documentation
- **MUST** update ARCHITECTURE.md for structural changes
- **MUST** document breaking changes in CHANGELOG.md
- **MUST** provide migration guides for API changes

---

## VII. Git & Release Standards

### 7.1 Commit Messages
```
<type>: <subject>

<body>

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `chore`

### 7.2 Branch Strategy
- `master` - production-ready code
- `develop` - integration branch (Phase 2+)
- `feature/*` - feature branches (Phase 2+)

### 7.3 Pre-merge Checklist
- **MUST** pass all tests: `make test`
- **MUST** pass CI checks: `make ci-check`
- **MUST** maintain coverage: `make coverage`
- **MUST** update documentation

---

## VIII. Security & Safety

### 8.1 Input Validation
- **MUST** validate all file paths
- **MUST** check file existence before opening
- **MUST** handle corrupted PDFs gracefully
- **MUST** limit memory usage (prevent DoS)

### 8.2 File Handling
- **MUST** close all file handles (use defer)
- **MUST** handle read/write errors
- **MUST** never expose file system internals to UI

### 8.3 Safe Defaults
- **MUST** use safe defaults for all configuration
- **MUST** fail securely on errors
- **MUST** avoid undefined behavior

---

## IX. Validation & Enforcement

### 9.1 Automated Checks
```bash
# Pre-commit validation
make ci-check    # Runs: fmt, vet, lint, test

# Coverage validation
make coverage    # Must show >80% coverage

# Performance validation
make bench       # Benchmarks must meet targets
```

### 9.2 Manual Review
- Code review checklist includes constitutional compliance
- Milestone reviews verify adherence
- Quarterly constitutional review for updates

### 9.3 Exemptions
Exemptions require:
1. Clear justification in commit message
2. Documentation in ADR (Architecture Decision Record)
3. Plan to remove exemption in future

---

## X. Amendment Process

**This constitution is immutable for Phase 1.**

Post-Phase 1 amendments require:
1. Proposal in GitHub issue
2. Impact analysis
3. Team consensus
4. Version bump in this document

---

**Constitutional Compliance**: All code merged into LUMOS MUST comply with this framework. Non-compliant code will be rejected regardless of functionality.

**Last Review**: 2025-11-10
**Next Review**: Upon Phase 2 start
