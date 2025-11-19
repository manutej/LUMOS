# LUMOS Testing Strategy

**Version**: 1.0.0
**Last Updated**: 2025-11-17
**Current Coverage**: 94.4%
**Target Coverage**: >80% (Exceeding)

---

## Testing Philosophy

**"Test behavior, not implementation"**

We test from the outside in:
1. What the user experiences (integration tests)
2. What the API promises (unit tests)
3. What performance we guarantee (benchmarks)

---

## Testing Pyramid

```
        /\           E2E Tests (5%)
       /  \          - Full user workflows
      /    \         - Real PDFs
     /      \
    /  Intg  \       Integration Tests (25%)
   /   Tests  \      - Component interactions
  /            \     - TUI behavior
 /              \    - File operations
/   Unit Tests   \   Unit Tests (70%)
──────────────────   - Pure functions
                     - State management
                     - Error paths
```

---

## Current Test Status

### Coverage by Package
```
pkg/pdf/document.go    94.2%  ✅
pkg/pdf/cache.go       96.2%  ✅
pkg/pdf/search.go      97.9%  ✅
pkg/ui/                [TODO] 🚧
pkg/config/theme.go    100%   ✅
────────────────────────────
Overall:               94.4%  ✅
```

### Test Inventory
- **Unit Tests**: 42 (all passing)
- **Benchmarks**: 9 (all meeting targets)
- **Integration**: 0 (pending TUI)
- **E2E**: 0 (pending full stack)

---

## Test Categories

### 1. Unit Tests (Pure Functions)

**Location**: `*_test.go` alongside source
**Naming**: `Test{Function}_{Scenario}_{Expected}`
**Tool**: Standard `testing` package

```go
func TestCacheGet_ExistingKey_ReturnsValue(t *testing.T) {
    // Given: Setup
    cache := NewLRUCache(5)
    cache.Put(1, "content")

    // When: Action
    result, found := cache.Get(1)

    // Then: Assert
    assert.True(t, found)
    assert.Equal(t, "content", result)
}
```

### 2. Table-Driven Tests (Multiple Scenarios)

**When**: Testing multiple inputs/outputs
**Pattern**: Single test, multiple cases

```go
func TestExtractContext(t *testing.T) {
    tests := []struct {
        name     string
        text     string
        pos      int
        before   int
        after    int
        expected string
    }{
        {
            name:     "middle_of_text",
            text:     "The quick brown fox jumps",
            pos:      10,
            before:   5,
            after:    5,
            expected: "quick brown fox",
        },
        // ... more cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ExtractContext(tt.text, tt.pos, tt.before, tt.after)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 3. Integration Tests (Component Interaction)

**Location**: `test/integration/`
**Focus**: Real component interaction
**Mock**: Only external systems

```go
func TestPDFLoadAndCache(t *testing.T) {
    // Real PDF, real cache, real search
    doc, err := pdf.LoadDocument("test/fixtures/multipage.pdf")
    require.NoError(t, err)

    // First load - cache miss
    content1, err := doc.GetPage(1)
    require.NoError(t, err)

    // Second load - cache hit (should be faster)
    start := time.Now()
    content2, err := doc.GetPage(1)
    duration := time.Since(start)

    assert.Equal(t, content1, content2)
    assert.Less(t, duration, 1*time.Millisecond) // Cache hit
}
```

### 4. TUI Tests (Bubble Tea)

**Challenge**: Testing interactive TUI
**Solution**: Message-based testing

```go
func TestTUINavigation(t *testing.T) {
    // Create model with test PDF
    doc := loadTestPDF(t)
    model := ui.NewModel(doc)

    // Simulate initialization
    cmd := model.Init()
    assert.NotNil(t, cmd)

    // Simulate key press
    model, cmd = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

    // Verify viewport scrolled
    assert.Greater(t, model.viewport.YOffset, 0)
}
```

### 5. Benchmark Tests (Performance)

**Location**: `*_test.go` with `Benchmark` prefix
**Run**: `make bench` or `go test -bench=.`

```go
func BenchmarkCacheGet(b *testing.B) {
    cache := NewLRUCache(100)
    // Setup cache with data
    for i := 0; i < 100; i++ {
        cache.Put(i, fmt.Sprintf("content-%d", i))
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Get(i % 100)
    }
}

// Results:
// BenchmarkCacheGet-8  75000000  16.0 ns/op  0 B/op  0 allocs/op
```

### 6. Property-Based Tests (Edge Cases)

**When**: Testing invariants
**Tool**: `testing/quick` or similar

```go
func TestSearchInvariant(t *testing.T) {
    // Property: Search results are always within text bounds
    f := func(text string, query string) bool {
        if query == "" {
            return true // Skip empty
        }

        results := Search(text, query)
        for _, r := range results {
            if r.Start < 0 || r.End > len(text) {
                return false
            }
        }
        return true
    }

    if err := quick.Check(f, nil); err != nil {
        t.Error(err)
    }
}
```

---

## Test Fixtures

### PDF Test Files
```
test/fixtures/
├── simple.pdf       # 1 page, basic text
├── multipage.pdf    # 5 pages, navigation testing
├── search_test.pdf  # Known search patterns
├── large.pdf        # 100+ pages (Phase 2)
├── complex.pdf      # Tables, images (Phase 3)
└── encrypted.pdf    # Error handling
```

### Generating Fixtures
```python
# test/generate_fixtures.py
from reportlab.pdfgen import canvas

def create_simple_pdf():
    c = canvas.Canvas("fixtures/simple.pdf")
    c.drawString(100, 750, "Hello, LUMOS!")
    c.save()
```

---

## Testing Commands

### Quick Test (Unit Only)
```bash
make test
# or
go test ./pkg/...
```

### Full Test Suite
```bash
make test-all
# Runs: unit, integration, benchmarks
```

### Coverage Report
```bash
make coverage
# Opens HTML report in browser
```

### Specific Package
```bash
go test -v ./pkg/pdf/
go test -v ./pkg/ui/
```

### Race Detection
```bash
make test-race
# or
go test -race ./...
```

### Benchmarks
```bash
make bench
# or
go test -bench=. -benchmem ./...
```

### Single Test
```bash
go test -v -run TestCacheGet ./pkg/pdf/
```

---

## Test Writing Guidelines

### 1. Test Naming Convention
```go
// Pattern: Test{Function}_{Scenario}_{Expected}
TestLoadDocument_ValidPDF_ReturnsDocument
TestLoadDocument_MissingFile_ReturnsError
TestLoadDocument_CorruptPDF_ReturnsError
```

### 2. Test Structure (AAA)
```go
func TestExample(t *testing.T) {
    // Arrange (Given)
    setup := createTestData()

    // Act (When)
    result := FunctionUnderTest(setup)

    // Assert (Then)
    assert.Equal(t, expected, result)
}
```

### 3. Use Subtests for Groups
```go
func TestNavigation(t *testing.T) {
    t.Run("move_down", func(t *testing.T) {
        // Test j key
    })

    t.Run("move_up", func(t *testing.T) {
        // Test k key
    })
}
```

### 4. Cleanup Resources
```go
func TestWithFile(t *testing.T) {
    file := createTempFile(t)
    defer os.Remove(file.Name()) // Always cleanup

    // Test with file
}
```

### 5. Parallel Tests (When Safe)
```go
func TestIndependent(t *testing.T) {
    t.Parallel() // Run concurrently
    // Test code
}
```

---

## Coverage Requirements

### Constitutional Minimums
- **Overall**: 80% minimum
- **Critical packages** (`pdf/`, `ui/`): 90% target
- **New code**: Must include tests

### Current Status (94.4% ✅)
Exceeding requirements significantly!

### Coverage Gaps to Address
1. **pkg/ui/**: Pending TUI implementation
2. **cmd/lumos/**: Minimal coverage needed (thin layer)
3. **Error paths**: Some error branches untested

### Measuring Coverage
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View in terminal
go tool cover -func=coverage.out

# View in browser (recommended)
go tool cover -html=coverage.out

# Check specific package
go test -cover ./pkg/pdf/
```

---

## Performance Testing

### Benchmark Requirements
Each performance-critical function needs a benchmark:

```go
BenchmarkCacheGet       # Target: <100ns
BenchmarkCachePut       # Target: <100ns
BenchmarkSearch         # Target: <50μs/KB
BenchmarkPageLoad       # Target: <50ms
BenchmarkTUIRender      # Target: <16ms (60fps)
```

### Performance Regression Detection
```bash
# Save baseline
go test -bench=. ./... > baseline.txt

# After changes, compare
go test -bench=. ./... > current.txt
benchcmp baseline.txt current.txt
```

### Profiling During Tests
```go
func TestExpensiveOperation(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping expensive test")
    }

    // CPU profiling
    defer profile.Start().Stop()

    // Run expensive operation
    result := ExpensiveOperation()
    assert.NotNil(t, result)
}
```

---

## TUI Testing Strategy

### Challenge
Bubble Tea TUIs are event-driven and async.

### Solution: Message Testing
```go
type mockCmd struct {
    msg tea.Msg
}

func (m mockCmd) Run() tea.Msg {
    return m.msg
}

func TestPageLoad(t *testing.T) {
    model := NewModel()

    // Simulate page loaded message
    msg := PageLoadedMsg{Content: "Page content"}
    newModel, cmd := model.Update(msg)

    // Verify state updated
    assert.Equal(t, "Page content", newModel.viewport.Content)
}
```

### Visual Regression Testing
```go
func TestViewOutput(t *testing.T) {
    model := setupTestModel(t)
    view := model.View()

    // Compare with golden file
    golden := readGoldenFile(t, "testdata/view_output.golden")
    assert.Equal(t, golden, view)
}
```

---

## Manual Testing Checklist

### Before Each Milestone
- [ ] All automated tests pass
- [ ] No race conditions detected
- [ ] Coverage meets requirements
- [ ] Benchmarks meet targets

### TUI Manual Tests
- [ ] Launch with various PDFs
- [ ] Test all keybindings
- [ ] Resize terminal window
- [ ] Test on min terminal (80x24)
- [ ] Test on large terminal (200x60)
- [ ] Verify no visual artifacts
- [ ] Check memory usage over time

### Terminal Compatibility
- [ ] macOS Terminal.app
- [ ] iTerm2
- [ ] Alacritty
- [ ] Linux gnome-terminal
- [ ] tmux/screen

---

## Test Maintenance

### Keep Tests Fast
- Target: Full suite <10 seconds
- Use test fixtures, not real files
- Parallelize where possible
- Skip expensive tests with -short

### Keep Tests Reliable
- No flaky tests allowed
- Deterministic outcomes only
- Mock time, randomness, I/O
- Clean up resources

### Keep Tests Readable
- Clear test names
- Minimal setup
- One assertion per test (ideal)
- Use helpers for common setup

---

## CI/CD Integration

### Pre-commit Checks
```bash
make ci-check
# Runs: fmt, vet, lint, test
```

### GitHub Actions (Future)
```yaml
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
      - run: make test
      - run: make coverage
```

---

## Testing Milestones

### Milestone 1.4 (TUI Framework)
- [ ] 15 unit tests for UI components
- [ ] 5 integration tests for TUI
- [ ] Benchmarks for render pipeline
- [ ] Manual testing on 3 terminals

### Milestone 1.5 (Keybindings)
- [ ] Test each keybinding
- [ ] Test key combinations
- [ ] Test help overlay
- [ ] Verify no conflicts

### Milestone 1.6 (Polish)
- [ ] Performance benchmarks
- [ ] Theme testing
- [ ] Error handling tests
- [ ] Final coverage check

---

## Success Metrics

### Quantitative
- ✅ Coverage >80% (Currently 94.4%)
- ✅ All tests passing (42/42)
- ✅ Benchmarks meeting targets
- ✅ No race conditions

### Qualitative
- ✅ Tests are documentation
- ✅ Tests catch real bugs
- ✅ Tests enable refactoring
- ✅ Tests build confidence

---

**Testing Mantra**: If it's not tested, it's broken.