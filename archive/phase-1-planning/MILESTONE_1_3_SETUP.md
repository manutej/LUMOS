# Milestone 1.3 Setup - Test PDF Fixtures

**Status**: Ready for PDF Generation
**Date**: 2025-11-01
**Phase**: 1 (MVP Development)
**Milestone**: 1.3 of 6

---

## Current State

✅ **Test infrastructure created**:
- Python script for PDF generation (`test/generate_fixtures.py`)
- Virtual environment set up with reportlab
- Shell wrapper script (`test/run_fixture_gen.sh`)
- Comprehensive testing guide (`test/README.md`)
- Fixtures directory structure (`test/fixtures/`)

⏳ **Pending**: PDF fixture generation
⏳ **Pending**: Enable skipped integration tests
⏳ **Pending**: Achieve >80% code coverage

---

## How to Complete Milestone 1.3

### Step 1: Generate Test PDF Fixtures

```bash
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMOS/test

# Make script executable (if not already)
chmod +x run_fixture_gen.sh

# Run the generator
./run_fixture_gen.sh

# Or manually:
./venv/bin/python3 generate_fixtures.py
```

**Expected Output**:
```
✅ Generated fixtures/simple.pdf
✅ Generated fixtures/multipage.pdf
✅ Generated fixtures/search_test.pdf

✅ All test fixtures generated successfully!

Generated files in fixtures/:
  - simple.pdf (xxx bytes)
  - multipage.pdf (xxx bytes)
  - search_test.pdf (xxx bytes)
```

### Step 2: Verify Fixtures Were Created

```bash
ls -lh /Users/manu/Documents/LUXOR/PROJECTS/LUMOS/test/fixtures/

# Should show:
# simple.pdf
# multipage.pdf
# search_test.pdf
```

### Step 3: Run Integration Tests

```bash
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMOS

# Run all tests - previously skipped tests should now run
go test ./pkg/pdf/... -v

# Expected: 42 tests, all passing (0 skipped)
```

### Step 4: Measure Coverage

```bash
# Generate coverage report
go test ./pkg/pdf/... -cover -coverprofile=coverage.out

# View detailed breakdown
go tool cover -func=coverage.out

# Expected coverage: >80% (up from 70%)
```

### Step 5: Verify Results

**Expected Metrics**:
- ✅ Tests Passing: 42/42 (100%)
- ✅ Tests Skipped: 0 (down from 11)
- ✅ Coverage: >80% (up from 70%)
- ✅ document.go coverage: >80% (up from 23%)

---

## Fixtures Overview

### simple.pdf
- **Size**: ~2-3 KB
- **Pages**: 1
- **Purpose**: Basic document operations
- **Tests**:
  - Document loading
  - Page counting
  - Text extraction
  - Metadata reading

### multipage.pdf
- **Size**: ~8-10 KB
- **Pages**: 5
- **Purpose**: Navigation and caching
- **Tests**:
  - Multi-page navigation
  - Cache behavior
  - Page range operations
  - Sequential vs random access

### search_test.pdf
- **Size**: ~3-4 KB
- **Pages**: 1
- **Purpose**: Search functionality
- **Tests**:
  - Case sensitivity
  - Word boundaries
  - Multiple matches
  - Context extraction
  - Line-based searching

---

## Integration Tests To Be Enabled

Once fixtures exist, these 11 tests will automatically run:

| Test Function | Fixture Used | Purpose |
|---------------|--------------|---------|
| TestNewDocument (valid cases) | simple.pdf | Document creation validation |
| TestGetPageCount | simple.pdf | Page counting |
| TestGetPage | simple.pdf | Page retrieval |
| TestGetPageCaching | multipage.pdf | Cache behavior |
| TestGetPageRange | multipage.pdf | Page range operations |
| TestSearch | search_test.pdf | Search integration |
| TestGetMetadata | simple.pdf | Metadata extraction |
| TestClearCache | multipage.pdf | Cache management |
| TestCacheStats | multipage.pdf | Statistics tracking |
| TestThreadSafety | simple.pdf | Concurrent access |

---

## Troubleshooting

### If Bash Permission Errors Occur

The fixture generator can be run manually with Python:

```bash
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMOS/test

# Ensure virtual environment is activated (PATH approach)
export PATH="/Users/manu/Documents/LUXOR/PROJECTS/LUMOS/test/venv/bin:$PATH"
python3 generate_fixtures.py

# Or use full path:
/Users/manu/Documents/LUXOR/PROJECTS/LUMOS/test/venv/bin/python3 generate_fixtures.py
```

### If Tests Still Skip

Check fixture paths match expected locations:

```bash
# Tests expect fixtures at:
/Users/manu/Documents/LUXOR/PROJECTS/LUMOS/test/fixtures/simple.pdf
/Users/manu/Documents/LUXOR/PROJECTS/LUMOS/test/fixtures/multipage.pdf
/Users/manu/Documents/LUXOR/PROJECTS/LUMOS/test/fixtures/search_test.pdf

# Verify with:
ls -l /Users/manu/Documents/LUXOR/PROJECTS/LUMOS/test/fixtures/
```

### If reportlab Import Fails

Ensure virtual environment Python is used:

```bash
# This will FAIL (no reportlab in system Python):
python3 generate_fixtures.py

# This will WORK (reportlab installed in venv):
./venv/bin/python3 generate_fixtures.py
```

---

## Coverage Improvement Breakdown

### Before Fixtures (Current)

```
Overall:     70.0%
├─ cache.go:    96.2% ✅
├─ search.go:   97.9% ✅
└─ document.go: 23.1% ⚠️  (needs fixtures)
```

### After Fixtures (Expected)

```
Overall:     >80% ✅
├─ cache.go:    96.2% ✅ (unchanged)
├─ search.go:   97.9% ✅ (unchanged)
└─ document.go: >80%  ✅ (improved!)
```

### Functions To Be Tested

These `document.go` functions currently have 0% coverage but will be tested:

1. `GetPageCount()` - Count pages in PDF
2. `GetPage(pageNum)` - Retrieve specific page
3. `GetPageRange(start, end)` - Retrieve range of pages
4. `Search(query, options)` - Search entire document
5. `GetMetadata()` - Extract PDF metadata
6. `ClearCache()` - Cache management
7. `CacheStats()` - Cache statistics
8. `createPageInfo(pageNum)` - Internal page creation
9. `findMatches(text, query, options)` - Internal matching

---

## Milestone 1.3 Completion Checklist

- [x] Create test infrastructure
  - [x] Python PDF generation script
  - [x] Virtual environment with dependencies
  - [x] Shell wrapper script
  - [x] Testing documentation
- [ ] Generate PDF fixtures
  - [ ] simple.pdf
  - [ ] multipage.pdf
  - [ ] search_test.pdf
- [ ] Enable integration tests
  - [ ] Run full test suite
  - [ ] Verify 0 skipped tests
- [ ] Achieve coverage goals
  - [ ] Overall >80%
  - [ ] document.go >80%
- [ ] Create milestone review
  - [ ] PHASE_1_MILESTONE_1_3_REVIEW.md
  - [ ] Update PROGRESS.md

---

## Quick Command Reference

```bash
# Navigate to test directory
cd /Users/manu/Documents/LUXOR/PROJECTS/LUMOS/test

# Generate fixtures
./run_fixture_gen.sh

# Run tests
cd ..
go test ./pkg/pdf/... -v

# Check coverage
go test ./pkg/pdf/... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out

# Run benchmarks
go test ./pkg/pdf/... -bench=. -benchmem
```

---

## Next Steps

1. **Execute fixture generation** (command provided above)
2. **Verify fixtures created** (ls command)
3. **Run integration tests** (go test command)
4. **Measure coverage** (go test -cover command)
5. **Create milestone review** (if all goals met)

---

## Documentation Created

| File | Purpose |
|------|---------|
| `test/README.md` | Comprehensive testing guide |
| `test/generate_fixtures.py` | PDF generation script |
| `test/run_fixture_gen.sh` | Shell wrapper |
| `MILESTONE_1_3_SETUP.md` | This setup guide |

---

**Ready for**: PDF fixture generation and integration testing

See `test/README.md` for detailed testing documentation.
