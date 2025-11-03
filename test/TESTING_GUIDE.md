# LUMOS Testing Guide

**Last Updated**: 2025-11-01  
**Phase**: 1 (MVP Development)  
**Milestone**: 1.3 (Test PDF Fixtures & Integration Testing)

---

## Table of Contents

1. [Overview](#overview)
2. [Test Fixtures](#test-fixtures)
3. [Running Tests](#running-tests)
4. [Test Categories](#test-categories)
5. [Code Coverage](#code-coverage)
6. [Benchmark Testing](#benchmark-testing)
7. [Adding New Tests](#adding-new-tests)
8. [CI/CD Integration](#cicd-integration)
9. [Troubleshooting](#troubleshooting)

---

## Overview

LUMOS uses Go's built-in testing framework with comprehensive unit tests, integration tests, and benchmarks. The test suite covers:

- **PDF document operations** (loading, parsing, metadata extraction)
- **Page caching** (LRU cache with thread safety)
- **Search functionality** (case-sensitive, word boundaries, context extraction)
- **Text processing** (line splitting, highlighting, navigation)

### Test Statistics

```
Total Tests: 42
Unit Tests: 31
Integration Tests: 11
Benchmarks: 9
Coverage: 70%+ (target >80%)
```

---

## Test Fixtures

Test PDFs are located in `test/fixtures/` and are generated using the Python script `test/generate_fixtures.py`.

### Available Fixtures

| File | Pages | Size | Purpose |
|------|-------|------|---------|
| `simple.pdf` | 1 | 1.9KB | Basic document loading and text extraction |
| `multipage.pdf` | 5 | 4.3KB | Multi-page navigation and caching |
| `search_test.pdf` | 1 | 2.2KB | Search pattern matching and context extraction |

### Fixture Contents

#### simple.pdf
- **Title**: "LUMOS Test Document"
- **Author**: "LUMOS Test Suite"
- **Subject**: "Testing PDF Reading Functionality"
- **Content**: Basic text with metadata for testing document loading

#### multipage.pdf
- **Title**: "LUMOS Multi-Page Test"
- **Pages**: 5 distinct pages with unique content
- **Content**: Each page has identifiable text for testing:
  - Page navigation
  - Cache behavior
  - Sequential and random access patterns

#### search_test.pdf
- **Title**: "LUMOS Search Test"
- **Content**: Specific patterns for search testing:
  - Case sensitivity: "test", "Test", "TEST"
  - Word boundaries: "testing", "test", "retest"
  - Multiple occurrences: "search" appears 3 times
  - Special characters: email addresses, hyphens, underscores
  - Long text for context extraction

### Regenerating Fixtures

If you need to regenerate or modify fixtures:

```bash
cd test
source venv/bin/activate  # Or create: python3 -m venv venv
python3 generate_fixtures.py
```

The script will regenerate all PDFs in `test/fixtures/`.

---

## Running Tests

### Run All Tests

```bash
# From project root
go test ./...

# Verbose output
go test ./... -v

# Run only PDF package tests
go test ./pkg/pdf/... -v
```

### Run Specific Tests

```bash
# Run tests matching a pattern
go test ./pkg/pdf/... -run TestNewDocument

# Run cache tests only
go test ./pkg/pdf/... -run TestLRUCache

# Run search tests only
go test ./pkg/pdf/... -run TestTextSearch
```

### Run With Coverage

```bash
# Generate coverage report
go test ./pkg/pdf/... -cover

# Generate detailed coverage profile
go test ./pkg/pdf/... -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out

# Show function-level coverage
go tool cover -func=coverage.out
```

### Run Benchmarks

```bash
# Run all benchmarks
go test ./pkg/pdf/... -bench=.

# Run benchmarks with memory stats
go test ./pkg/pdf/... -bench=. -benchmem

# Run specific benchmark
go test ./pkg/pdf/... -bench=BenchmarkLRUCache

# Run benchmarks multiple times for accuracy
go test ./pkg/pdf/... -bench=. -benchtime=10s -count=5
```

---

## Test Categories

### Unit Tests

#### Cache Tests (`cache_test.go`)

Tests for LRU cache implementation:
- `TestNewLRUCache` - Cache creation with various sizes
- `TestLRUCache_PutAndGet` - Basic get/put operations
- `TestLRUCache_Eviction` - LRU eviction behavior
- `TestLRUCache_ThreadSafety` - Concurrent access
- `TestLRUCache_Stats` - Statistics tracking

**Coverage**: 96.2%

#### Search Tests (`search_test.go`)

Tests for search functionality:
- `TestCaseSensitiveMatch` - Case-sensitive pattern matching
- `TestCaseInsensitiveMatch` - Case-insensitive matching
- `TestWordMatch` - Word boundary detection
- `TestExtractContext` - Context extraction around matches
- `TestTextToLines` - Text line splitting
- `TestHighlightMatches` - Match highlighting

**Coverage**: 97.9%

#### Document Tests (`document_test.go`)

Tests for PDF document operations:
- `TestNewDocument` - Document creation
- `TestGetPageCount` - Page counting
- `TestGetPage` - Page retrieval
- `TestGetPageCaching` - Cache behavior
- `TestSearch` - Search integration
- `TestGetMetadata` - Metadata extraction

**Coverage**: 23.1% (unit tests) → 80%+ (with integration tests)

### Integration Tests

Integration tests use real PDF fixtures from `test/fixtures/`:

```go
func TestNewDocument(t *testing.T) {
    doc, err := NewDocument("../../test/fixtures/simple.pdf", 5)
    if err != nil {
        t.Skipf("Test PDF fixture not found: %v", err)
    }
    // Test with real PDF...
}
```

These tests validate:
- Real PDF parsing with ledongthuc/pdf library
- Metadata extraction from actual PDF files
- Text extraction accuracy
- Multi-page document handling
- Cache behavior with real data

---

## Code Coverage

### Current Coverage

```
Overall: 70.0%

By File:
- cache.go:    96.2%  (excellent)
- search.go:   97.9%  (excellent)
- document.go: 23.1%  (unit tests only)
```

### Coverage Goals

- **Target**: >80% overall coverage
- **Critical paths**: >90% coverage
- **Edge cases**: All known edge cases covered

### Improving Coverage

To increase `document.go` coverage:

1. **Enable integration tests** - Use real PDF fixtures
2. **Test error paths** - Invalid PDFs, missing files, corrupt data
3. **Test metadata extraction** - Once implemented in Phase 1.3
4. **Test concurrent document access** - Thread safety validation

### Coverage Report

```bash
# Generate detailed report
go test ./pkg/pdf/... -coverprofile=coverage.out

# View in terminal
go tool cover -func=coverage.out

# View in browser (interactive)
go tool cover -html=coverage.out
```

The HTML report shows:
- Green: Covered code
- Red: Uncovered code
- Grey: Not executable (comments, declarations)

---

## Benchmark Testing

### Benchmark Results

```
BenchmarkLRUCache_Put-12            16,764,399 ops    61.50 ns/op    24 B/op    2 allocs/op
BenchmarkLRUCache_Get_Hit-12        76,452,598 ops    16.05 ns/op     0 B/op    0 allocs/op
BenchmarkLRUCache_Get_Miss-12      147,783,690 ops     8.18 ns/op     0 B/op    0 allocs/op
BenchmarkCaseSensitiveMatch-12        113,108 ops 10,612.00 ns/op 25,208 B/op   12 allocs/op
BenchmarkHighlightMatches-12        1,000,000 ops  1,072.00 ns/op 12,840 B/op   10 allocs/op
```

### Performance Targets

| Operation | Target | Actual | Status |
|-----------|--------|--------|--------|
| Cache Hit | <100ns | 16ns | ✅ 6x better |
| Search (1KB) | <50μs | 10.6μs | ✅ 5x better |
| Highlight | <10μs | 1μs | ✅ 10x better |

### Writing Benchmarks

Example benchmark structure:

```go
func BenchmarkOperation(b *testing.B) {
    // Setup
    data := setupTestData()
    
    // Reset timer after setup
    b.ResetTimer()
    
    // Run operation b.N times
    for i := 0; i < b.N; i++ {
        performOperation(data)
    }
}
```

Best practices:
- Reset timer after setup
- Avoid allocations in the hot path
- Use realistic data sizes
- Run with `-benchmem` to track allocations
- Run multiple times for consistency

---

## Adding New Tests

### Test Structure

Follow Go testing conventions:

```go
func TestFeatureName(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:    "descriptive test case name",
            input:   testInput,
            want:    expectedOutput,
            wantErr: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionUnderTest(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionUnderTest() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("FunctionUnderTest() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Test Naming Conventions

- **Test functions**: `TestFunctionName` or `TestFeature_Scenario`
- **Benchmark functions**: `BenchmarkFunctionName`
- **Sub-tests**: Descriptive names with spaces
- **Test files**: `*_test.go` in same package

### Test Organization

```
pkg/pdf/
├── cache.go           # Implementation
├── cache_test.go      # Tests
├── search.go          # Implementation
├── search_test.go     # Tests
└── document.go        # Implementation
    └── document_test.go  # Tests
```

---

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Generate test fixtures
        run: |
          cd test
          python3 -m venv venv
          source venv/bin/activate
          pip install reportlab
          python3 generate_fixtures.py
      
      - name: Run tests
        run: go test ./... -v -race -coverprofile=coverage.out
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

### Pre-commit Hook

```bash
#!/bin/sh
# .git/hooks/pre-commit

echo "Running tests..."
go test ./... -short

if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi

echo "Tests passed!"
```

---

## Troubleshooting

### Common Issues

#### 1. Test PDFs Not Found

**Error**: `Test PDF fixture not found: ../../test/fixtures/simple.pdf`

**Solution**:
```bash
cd test
python3 -m venv venv
source venv/bin/activate
pip install reportlab
python3 generate_fixtures.py
```

#### 2. Import Errors

**Error**: `package github.com/luxor/lumos/pkg/pdf is not in GOPATH`

**Solution**:
```bash
# Ensure you're in the project root
go mod download
go mod tidy
```

#### 3. Race Conditions

**Error**: Test failures only when run with `-race`

**Solution**:
- Check for unprotected shared state
- Add proper mutex locking
- Use `go test -race` during development

#### 4. Flaky Tests

**Symptoms**: Tests pass sometimes, fail other times

**Common causes**:
- Time-dependent logic
- Race conditions
- Reliance on external state
- Order-dependent tests

**Solution**:
- Make tests deterministic
- Use fixed random seeds for testing
- Isolate test state
- Use `t.Parallel()` to detect race conditions

#### 5. Slow Tests

**Solution**:
```bash
# Skip slow tests in development
go test -short ./...

# Mark slow tests
func TestSlowOperation(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping slow test")
    }
    // Test code...
}
```

---

## Quick Reference

### Most Common Commands

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run benchmarks
go test ./... -bench=. -benchmem

# Run specific package
go test ./pkg/pdf/...

# Verbose output
go test ./... -v

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Selection

```bash
# Run only cache tests
go test ./... -run Cache

# Run only benchmark tests
go test ./... -bench=. -run=^$

# Run one specific test
go test ./... -run TestNewDocument/valid_PDF
```

### Debugging Tests

```bash
# Print all test output
go test ./... -v

# Run with race detector
go test ./... -race

# Run with CPU profiling
go test ./... -cpuprofile=cpu.prof

# Run with memory profiling
go test ./... -memprofile=mem.prof
```

---

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table-Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Advanced Testing](https://about.sourcegraph.com/blog/go/advanced-testing-in-go)
- [Benchmarking](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)

---

**Next Steps**: After completing Milestone 1.3, proceed to Milestone 1.4 (Basic TUI Framework)
