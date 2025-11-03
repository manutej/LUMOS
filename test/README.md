# LUMOS Test Suite

This directory contains the test suite for LUMOS, including unit tests, integration tests, and test fixtures.

## Directory Structure

```
test/
├── README.md                 # This file
├── generate_fixtures.py      # Python script to generate test PDFs
├── run_fixture_gen.sh        # Shell script wrapper
├── venv/                     # Python virtual environment
└── fixtures/                 # Test PDF files
    ├── simple.pdf           # Single-page basic PDF
    ├── multipage.pdf        # 5-page PDF for navigation testing
    └── search_test.pdf      # PDF with specific search patterns
```

## Quick Start

### 1. Generate Test Fixtures

The test suite requires PDF fixtures for integration testing. Generate them using the provided Python script:

```bash
cd test

# Generate PDFs (using virtual environment)
./run_fixture_gen.sh

# Or manually:
./venv/bin/python3 generate_fixtures.py
```

This will create three test PDFs in the `fixtures/` directory:

- **simple.pdf** (1 page) - Basic document for testing loading, page counting, and text extraction
- **multipage.pdf** (5 pages) - Multi-page document for testing navigation and caching
- **search_test.pdf** (1 page) - Document with specific patterns for search testing

### 2. Run Tests

```bash
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMOS

# Run all tests
go test ./pkg/pdf/... -v

# Run with coverage
go test ./pkg/pdf/... -cover -coverprofile=coverage.out

# View detailed coverage
go tool cover -func=coverage.out

# View coverage in browser
go tool cover -html=coverage.out

# Run benchmarks
go test ./pkg/pdf/... -bench=. -benchmem
```

### 3. Expected Results

After generating fixtures:

```
Tests:           42 total (100% passing)
Coverage:        >80% (up from 70% without fixtures)
Benchmarks:      9 total
Execution Time:  <1 second
```

## Test Files

### Unit Tests

1. **cache_test.go** - LRU cache operations
   - 12 test functions
   - 5 benchmark functions
   - Tests: creation, get/put, eviction, thread safety, statistics

2. **search_test.go** - Search functionality
   - 14 test functions
   - 4 benchmark functions
   - Tests: navigation, matching algorithms, context extraction, line splitting

3. **document_test.go** - Document operations
   - 16 test functions (11 require PDF fixtures)
   - Tests: document creation, page retrieval, caching, metadata

## Test Fixtures

### simple.pdf

**Purpose**: Basic PDF operations testing

**Contents**:
- 1 page
- Simple text content
- Metadata fields (title, author, subject, creator)
- ~50 words
- Multiple lines

**Tests**:
- Document loading
- Page counting
- Text extraction
- Metadata reading

### multipage.pdf

**Purpose**: Pagination and caching testing

**Contents**:
- 5 pages with different content per page
- Page numbers in text
- Clear page boundaries
- Sequential navigation markers

**Tests**:
- Multi-page navigation
- Cache behavior (LRU eviction with limited cache)
- Page range operations
- Sequential vs random access

### search_test.pdf

**Purpose**: Search functionality testing

**Contents**:
- Case sensitivity test words ("test", "Test", "TEST")
- Word boundary test words ("testing", "test", "retest")
- Multiple occurrences of "search"
- Special characters (email, hyphenated, underscored)
- Long text with "MATCH" for context extraction
- Line boundary keywords

**Tests**:
- Case-sensitive vs case-insensitive matching
- Word boundary detection
- Multiple match highlighting
- Context extraction
- Line-based searching

## Coverage Goals

| Module | Without Fixtures | With Fixtures | Target |
|--------|------------------|---------------|--------|
| cache.go | 96.2% | 96.2% | >90% |
| search.go | 97.9% | 97.9% | >90% |
| document.go | 23.1% | >80% | >80% |
| **Overall** | **70.0%** | **>80%** | **>80%** |

## Fixture Generation Details

### Requirements

- Python 3.x
- reportlab library
- pillow library (dependency of reportlab)

### Setup Virtual Environment (First Time Only)

```bash
cd test

# Create virtual environment
python3 -m venv venv

# Activate virtual environment
source venv/bin/activate  # macOS/Linux
# OR
venv\Scripts\activate     # Windows

# Install dependencies
pip install reportlab

# Deactivate when done
deactivate
```

### Manual PDF Generation

If you need to regenerate specific fixtures:

```python
# Edit generate_fixtures.py and run:
./venv/bin/python3 generate_fixtures.py

# Or generate specific PDFs:
./venv/bin/python3 -c "
from generate_fixtures import generate_simple_pdf
generate_simple_pdf('fixtures/simple.pdf')
"
```

## Test Organization

### Naming Conventions

- Test files: `*_test.go`
- Test functions: `Test<FunctionName>`
- Benchmark functions: `Benchmark<Operation>`
- Sub-tests: Use `t.Run("description", func(t *testing.T) {...})`

### Test Structure

All tests follow this pattern:

```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
    }{
        {"description", input1, output1},
        {"description", input2, output2},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Benchmark Structure

```go
func BenchmarkOperation(b *testing.B) {
    setup()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        operation()
    }
}
```

## Skipped Tests

The following tests are skipped until fixtures are generated:

1. `TestNewDocument` - valid PDF cases
2. `TestGetPageCount`
3. `TestGetPage`
4. `TestGetPageCaching`
5. `TestGetPageRange`
6. `TestSearch`
7. `TestGetMetadata`
8. `TestClearCache`
9. `TestCacheStats`
10. `TestThreadSafety`

All will be enabled once `test/fixtures/simple.pdf` exists.

## Troubleshooting

### Fixtures Not Found

**Error**: `Test PDF fixture not found: ../../test/fixtures/simple.pdf`

**Solution**: Generate fixtures first:
```bash
cd test
./run_fixture_gen.sh
```

### Permission Denied

**Error**: `permission denied: ./run_fixture_gen.sh`

**Solution**: Make script executable:
```bash
chmod +x run_fixture_gen.sh
```

### reportlab Not Found

**Error**: `ModuleNotFoundError: No module named 'reportlab'`

**Solution**: Ensure you're using the venv Python:
```bash
./venv/bin/python3 generate_fixtures.py
# NOT: python3 generate_fixtures.py
```

### Tests Still Skipped

**Cause**: Fixtures in wrong location

**Solution**: Verify fixtures location:
```bash
ls -l test/fixtures/
# Should show: simple.pdf, multipage.pdf, search_test.pdf
```

## Continuous Integration

When running tests in CI:

```yaml
# Example GitHub Actions workflow
- name: Generate Test Fixtures
  run: |
    cd test
    python3 -m venv venv
    ./venv/bin/pip install reportlab
    ./venv/bin/python3 generate_fixtures.py

- name: Run Tests
  run: go test ./pkg/pdf/... -v -cover

- name: Run Benchmarks
  run: go test ./pkg/pdf/... -bench=. -benchmem
```

## Performance Benchmarks

Expected benchmark results (with fixtures):

```
BenchmarkLRUCache_Put          ~60 ns/op
BenchmarkLRUCache_Get_Hit      ~16 ns/op
BenchmarkCaseSensitiveMatch    ~10 μs/op
BenchmarkHighlightMatches      ~1 μs/op
```

## Adding New Tests

### 1. Create Test Function

```go
func TestNewFeature(t *testing.T) {
    // Arrange
    input := setupInput()

    // Act
    result := NewFeature(input)

    // Assert
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

### 2. Add Fixture If Needed

Edit `generate_fixtures.py` to add new PDF patterns.

### 3. Update Coverage Target

Update this README with new coverage expectations.

## Resources

- Go Testing: https://golang.org/pkg/testing/
- Table-Driven Tests: https://github.com/golang/go/wiki/TableDrivenTests
- reportlab Documentation: https://www.reportlab.com/docs/reportlab-userguide.pdf

## Milestone Integration

This test suite is part of:
- **Milestone 1.2**: Core Testing & Benchmarking ✅
- **Milestone 1.3**: Test PDF Fixtures & Integration Testing (current)

See `../PHASE_1_PLAN.md` for full milestone details.
