# LUMOS Phase 1 Plan: MVP - Basic PDF Reader

**Status**: Planning & Definition Complete
**Date Started**: 2025-10-21
**Target Completion**: ~3 weeks
**Milestone-Based Delivery**: Review point at end of each milestone

---

## Overview

Phase 1 transforms the foundation code into a working MVP with comprehensive testing, performance validation, and documentation. The phase is divided into **6 core milestones**, each with clear success criteria and review points.

### Key Goals
- ✅ Compile and run without errors
- ✅ Load and display PDFs correctly
- ✅ All vim keybindings functional
- ✅ Meet performance targets (<100ms startup, <50MB memory)
- ✅ 80%+ test coverage
- ✅ Production-ready documentation

---

## Phase 1 Milestones

### 🚀 MILESTONE 1.1: Build & Compile Foundation

**Duration**: 2-3 days
**Owner**: Developer
**Goal**: Get the project compiling and running with basic PDF loading

#### Tasks

- [ ] **1.1.1** - Verify Go environment (1.21+)
  - Check `go version`
  - Verify GOPATH is set
  - Confirm workspace structure

- [ ] **1.1.2** - Download dependencies
  - Run `go mod download`
  - Resolve any import issues
  - Document dependency versions

- [ ] **1.1.3** - Build the project
  - Run `make build`
  - Fix any compilation errors
  - Verify binary creation at `./build/lumos`

- [ ] **1.1.4** - Basic functionality test
  - Load a simple PDF: `./build/lumos ~/Documents/test.pdf`
  - Verify basic display works
  - Document any issues

- [ ] **1.1.5** - Test with diverse PDFs
  - Simple text PDF
  - Multi-page document
  - PDF with images
  - PDF with tables
  - Large file (100+ pages)

- [ ] **1.1.6** - Verify all keybindings work
  - Navigation (j/k, d/u, gg/G)
  - Page navigation (Ctrl+N/P)
  - Search (/)
  - Help (?)
  - Quit (q)

- [ ] **1.1.7** - CLI argument handling
  - Verify `./build/lumos --help` works
  - Test `./build/lumos --version`
  - Test `./build/lumos --keys`
  - Test with invalid file
  - Test with home directory expansion (~)

#### Success Criteria

- [x] Binary compiles without warnings
- [ ] Can open and display PDFs
- [ ] No crashes on basic operations
- [ ] All keybindings respond
- [ ] CLI arguments work correctly
- [ ] Help and version info display

#### Deliverables

- `./build/lumos` binary (working)
- Build troubleshooting guide in PHASE_1_MILESTONE_1_1.md
- List of tested PDF files

#### Review Point

**After completing 1.1, file: `PROJECTS/LUMOS/PHASE_1_MILESTONE_1_1_REVIEW.md`**
- Document what compiled successfully
- List any issues encountered
- Note any needed fixes
- Pause for review before proceeding to 1.2

---

### 📊 MILESTONE 1.2: Core Testing & Benchmarking

**Duration**: 3-4 days
**Owner**: Developer
**Goal**: Establish baseline performance and identify optimization opportunities

#### Tasks

- [ ] **1.2.1** - Create test PDF fixtures
  - Simple 1-page PDF
  - Multi-page PDF (10 pages)
  - Large PDF (100+ pages)
  - PDF with images
  - PDF with special characters
  - Place in `test/fixtures/`

- [ ] **1.2.2** - Write unit tests for `pkg/pdf/document.go`
  - `TestLoadDocument` - Load PDF file
  - `TestGetPage` - Extract single page
  - `TestGetPageRange` - Extract multiple pages
  - `TestGetMetadata` - Get PDF metadata
  - `TestDocumentSearch` - Full-text search
  - `TestThreadSafety` - Concurrent access
  - Target: 90%+ coverage for document.go

- [ ] **1.2.3** - Write unit tests for `pkg/pdf/search.go`
  - `TestCaseSensitiveMatch`
  - `TestCaseInsensitiveMatch`
  - `TestWordMatch`
  - `TestHighlightMatches`
  - `TestLineExtraction`
  - Target: 90%+ coverage for search.go

- [ ] **1.2.4** - Write unit tests for `pkg/pdf/cache.go`
  - `TestCacheGet` - Retrieve from cache
  - `TestCachePut` - Add to cache
  - `TestCacheEviction` - LRU eviction
  - `TestCacheStats` - Statistics
  - `TestCacheHitRate` - Performance metric
  - `TestCacheThreadSafety` - Concurrent access
  - Target: 95%+ coverage for cache.go

- [ ] **1.2.5** - Write integration tests for `pkg/ui/`
  - `TestUIInitialize` - Model initialization
  - `TestNavigationKeys` - j/k/d/u/gg/G
  - `TestPageNavigation` - Ctrl+N/P
  - `TestSearchFlow` - Open search, enter query
  - `TestThemeSwitch` - Switch dark/light mode
  - `TestPaneNavigation` - Tab between panes
  - Target: 80%+ coverage for ui package

- [ ] **1.2.6** - Run benchmark suite
  - `BenchmarkDocumentLoad` - PDF loading time
  - `BenchmarkGetPage` - Page extraction time
  - `BenchmarkSearch` - Search operation time
  - `BenchmarkCacheHitRate` - Cache performance
  - `make bench` - Run all benchmarks
  - Document results in `BENCHMARKS_1_2.md`

- [ ] **1.2.7** - Profile performance
  - CPU profiling: `make profile-cpu`
  - Memory profiling: `make profile-mem`
  - Race condition detection: `make test-race`
  - Document findings

- [ ] **1.2.8** - Measure performance targets
  - Cold start time (first run)
  - Page load time (cached)
  - Page load time (uncached)
  - Memory usage (small PDF, large PDF)
  - Search time (10 pages, 100 pages)
  - Document in `PERFORMANCE_BASELINE_1_2.md`

#### Success Criteria

- [ ] 80%+ overall test coverage
- [ ] All unit tests passing
- [ ] All integration tests passing
- [ ] No race conditions detected
- [ ] Cold start <200ms (Phase 1 target: relaxed from 100ms)
- [ ] Page load <100ms (cached)
- [ ] Page load <300ms (uncached)
- [ ] Memory usage <100MB (Phase 1 target: relaxed from 50MB)
- [ ] Benchmarks established and documented

#### Deliverables

- Test fixtures in `test/fixtures/`
- Unit test files: `*_test.go` in each package
- Integration test file: `test/ui_integration_test.go`
- Benchmark results in `BENCHMARKS_1_2.md`
- Performance baseline in `PERFORMANCE_BASELINE_1_2.md`
- CPU/memory profiles in `test/profiles/`

#### Review Point

**After completing 1.2, file: `PROJECTS/LUMOS/PHASE_1_MILESTONE_1_2_REVIEW.md`**
- Test coverage summary
- Benchmark results vs targets
- Performance bottlenecks identified
- Optimization opportunities listed
- Pause for review before proceeding to 1.3

---

### 🐛 MILESTONE 1.3: Fix Bugs & Edge Cases

**Duration**: 3-4 days
**Owner**: Developer
**Goal**: Handle edge cases and fix all discovered bugs

#### Tasks

- [ ] **1.3.1** - Fix compilation warnings
  - Address any compiler warnings
  - Remove unused imports
  - Fix linting issues: `make lint`

- [ ] **1.3.2** - Fix test failures
  - Run all tests: `make test`
  - Identify failing tests
  - Debug and fix each failure
  - Achieve 100% pass rate

- [ ] **1.3.3** - Handle edge cases in PDF loading
  - Empty PDFs
  - Corrupted PDFs
  - Very large PDFs (>1GB)
  - PDFs with special encoding
  - PDFs with no text (images only)
  - Missing file
  - Permission denied

- [ ] **1.3.4** - Handle edge cases in text extraction
  - Empty pages
  - Pages with only images
  - Special UTF-8 characters
  - Right-to-left text
  - Very long lines
  - Tables and complex layouts

- [ ] **1.3.5** - Handle edge cases in search
  - Empty search query
  - Search with special characters
  - Search across page boundaries
  - No matches found
  - Very large result sets

- [ ] **1.3.6** - Handle UI edge cases
  - Window resize during operation
  - Very small terminal window
  - Very large terminal window
  - Rapid key presses
  - Invalid page numbers

- [ ] **1.3.7** - Verify error messages
  - Clear, helpful error messages
  - Proper error handling
  - Graceful degradation

- [ ] **1.3.8** - Test with real PDFs
  - Academic papers
  - Technical specifications
  - Books
  - Scanned documents
  - Mixed content (text + images)

#### Success Criteria

- [ ] 100% test pass rate
- [ ] No unhandled panics
- [ ] All error cases handled gracefully
- [ ] Clear error messages for users
- [ ] No race conditions
- [ ] Application stable with diverse inputs

#### Deliverables

- Fixed source code in `pkg/`
- Bug fix log in `PHASE_1_MILESTONE_1_3_FIXES.md`
- Edge case test coverage in test files
- Error handling documentation

#### Review Point

**After completing 1.3, file: `PROJECTS/LUMOS/PHASE_1_MILESTONE_1_3_REVIEW.md`**
- All bugs fixed summary
- Edge cases handled
- Test pass rate: 100%
- Stability assessment
- Pause for review before proceeding to 1.4

---

### ⚡ MILESTONE 1.4: Performance Optimization

**Duration**: 2-3 days
**Owner**: Developer
**Goal**: Achieve performance targets through optimization

#### Tasks

- [ ] **1.4.1** - Analyze performance bottlenecks
  - Review profiling data from 1.2
  - Identify hot paths
  - Prioritize optimizations

- [ ] **1.4.2** - Optimize PDF loading
  - Cache parsed PDFs
  - Lazy-load pages
  - Reduce memory allocations
  - Benchmark improvements

- [ ] **1.4.3** - Optimize text extraction
  - Reduce string allocations
  - Use buffers efficiently
  - Cache extracted text
  - Benchmark improvements

- [ ] **1.4.4** - Optimize search performance
  - Index full-text search
  - Use efficient string matching
  - Cache search results
  - Benchmark improvements

- [ ] **1.4.5** - Optimize UI rendering
  - Reduce screen refreshes
  - Cache styled strings
  - Optimize pane rendering
  - Benchmark improvements

- [ ] **1.4.6** - Optimize memory usage
  - Profile memory allocations
  - Reduce allocations in hot paths
  - Use sync.Pool for temporary objects
  - Measure before/after

- [ ] **1.4.7** - Verify cache effectiveness
  - Measure cache hit rate
  - Verify LRU eviction working
  - Benchmark page load times
  - Document in `CACHE_ANALYSIS_1_4.md`

- [ ] **1.4.8** - Measure final performance
  - Cold start time
  - Page load time (cached/uncached)
  - Memory usage
  - Search time
  - Compare to Phase 1 targets
  - Document in `PERFORMANCE_FINAL_1_4.md`

#### Performance Targets (Phase 1 - Relaxed)

| Metric | Phase 0 Target | Phase 1 Target | Success? |
|--------|---|---|---|
| Cold start | <100ms | <200ms | |
| Page load (cached) | <50ms | <100ms | |
| Page load (uncached) | <200ms | <300ms | |
| Memory usage (10MB PDF) | <50MB | <100MB | |
| Search time (10 pages) | <100ms | <200ms | |

#### Success Criteria

- [ ] Cold start <200ms
- [ ] Page load (cached) <100ms
- [ ] Page load (uncached) <300ms
- [ ] Memory usage <100MB
- [ ] Search time <200ms
- [ ] Cache hit rate >80%
- [ ] No memory leaks detected
- [ ] No performance regressions from 1.3

#### Deliverables

- Optimized source code in `pkg/`
- Cache analysis in `CACHE_ANALYSIS_1_4.md`
- Final performance metrics in `PERFORMANCE_FINAL_1_4.md`
- Optimization techniques documented
- Before/after benchmarks

#### Review Point

**After completing 1.4, file: `PROJECTS/LUMOS/PHASE_1_MILESTONE_1_4_REVIEW.md`**
- Performance metrics vs targets
- Optimization techniques applied
- Memory usage verified
- Cache effectiveness measured
- Pause for review before proceeding to 1.5

---

### 📈 MILESTONE 1.5: Test Coverage & Quality

**Duration**: 2-3 days
**Owner**: Developer
**Goal**: Achieve 80%+ test coverage with quality focus

#### Tasks

- [ ] **1.5.1** - Measure current coverage
  - Run `make coverage`
  - Generate coverage report
  - Identify untested code
  - Document in `COVERAGE_1_5.md`

- [ ] **1.5.2** - Fill coverage gaps
  - Add tests for uncovered code
  - Focus on public APIs first
  - Then internal functions
  - Target 80%+ coverage

- [ ] **1.5.3** - Write negative test cases
  - Invalid inputs
  - Boundary conditions
  - Error paths
  - Resource exhaustion

- [ ] **1.5.4** - Add stress tests
  - Load 1000-page PDF
  - Search with 1000 results
  - Rapid page navigation
  - Concurrent operations

- [ ] **1.5.5** - Code quality checks
  - Run `make fmt` - Format code
  - Run `make vet` - Vet code
  - Run `make lint` - Lint code
  - Fix all issues

- [ ] **1.5.6** - Document test strategy
  - Test coverage by package
  - Test types (unit/integration/stress)
  - How to run tests
  - How to add new tests
  - File: `docs/TESTING.md` (update)

- [ ] **1.5.7** - Verify reproducibility
  - Run tests multiple times
  - Run tests in different order
  - Run with `-race` flag
  - No flaky tests

- [ ] **1.5.8** - Create test documentation
  - How to run tests: `make test`
  - How to run with coverage: `make coverage`
  - How to run specific tests
  - How to add new tests
  - File: `TEST_GUIDE.md`

#### Success Criteria

- [ ] 80%+ overall code coverage
- [ ] 90%+ coverage for `pkg/pdf/`
- [ ] 80%+ coverage for `pkg/ui/`
- [ ] 90%+ coverage for `pkg/config/`
- [ ] 100% pass rate on all tests
- [ ] No flaky tests
- [ ] All code formatted, vetted, linted
- [ ] Test documentation complete

#### Deliverables

- Additional test files increasing coverage
- Coverage report in `COVERAGE_1_5.md`
- Updated test documentation
- Test guide in `TEST_GUIDE.md`
- Stress test results

#### Review Point

**After completing 1.5, file: `PROJECTS/LUMOS/PHASE_1_MILESTONE_1_5_REVIEW.md`**
- Code coverage metrics
- Test quality assessment
- Code quality checks passed
- Stress test results
- Pause for review before proceeding to 1.6

---

### 📚 MILESTONE 1.6: Documentation & Release

**Duration**: 3-4 days
**Owner**: Developer
**Goal**: Complete documentation and prepare for Phase 2

#### Tasks

- [ ] **1.6.1** - Update README.md
  - Installation instructions
  - Quick start guide
  - Feature list
  - Performance metrics
  - Contributing guidelines

- [ ] **1.6.2** - Complete QUICKSTART.md
  - 5-minute setup
  - Common commands
  - Understanding the code
  - Common development tasks
  - Debugging tips

- [ ] **1.6.3** - Update docs/ARCHITECTURE.md
  - System design review
  - Data flow diagrams
  - State management explanation
  - Performance characteristics
  - Design decisions rationale

- [ ] **1.6.4** - Update docs/DEVELOPMENT.md
  - Development setup
  - Code style guide
  - Common tasks with examples
  - Performance optimization guide
  - Debugging guide

- [ ] **1.6.5** - Create keybinding reference
  - File: `docs/KEYBINDINGS.md`
  - All keybindings listed
  - Vim commands explained
  - Examples of common workflows

- [ ] **1.6.6** - Create PHASE_1_SUMMARY.md
  - Milestones completed
  - Features implemented
  - Performance achieved
  - Test coverage
  - Known limitations

- [ ] **1.6.7** - Create Phase 2 roadmap
  - File: `PHASE_2_PLAN.md`
  - Features planned
  - Estimated timeline
  - Dependencies

- [ ] **1.6.8** - Code documentation
  - Comment all public functions
  - Document complex algorithms
  - Add examples in code
  - Generate godoc: `godoc -http=:6060`

- [ ] **1.6.9** - User-facing documentation
  - Help text in CLI
  - Error messages are clear
  - Man page (optional)
  - Troubleshooting guide

- [ ] **1.6.10** - Version and release
  - Tag version (v1.0.0-alpha or v0.1.0)
  - Create release notes
  - Update version in code
  - Prepare for distribution

#### Success Criteria

- [ ] README complete and accurate
- [ ] All docs up-to-date
- [ ] Code is well-commented
- [ ] Help text is clear
- [ ] User documentation complete
- [ ] Phase 2 roadmap clear
- [ ] All links work
- [ ] No TODOs left in docs

#### Deliverables

- Updated README.md
- Complete QUICKSTART.md
- Updated architecture documentation
- Updated development guide
- Keybindings reference (docs/KEYBINDINGS.md)
- Phase 1 summary (PHASE_1_SUMMARY.md)
- Phase 2 plan (PHASE_2_PLAN.md)
- Version tag in git
- Release notes

#### Review Point & Phase 1 Completion

**After completing 1.6, file: `PROJECTS/LUMOS/PHASE_1_COMPLETE.md`**
- Phase 1 deliverables checklist
- All milestones completed
- Final metrics and statistics
- Phase 2 readiness assessment
- **PHASE 1 COMPLETE** ✅

---

## Milestone Review Process

After **each milestone completion**, create a review document:

### Review Document Template

```markdown
# PHASE 1 MILESTONE X.Y REVIEW

**Date**: [Date completed]
**Status**: [Completed / Blocked / Needs revision]
**Duration**: [Actual time vs estimated]

## What Was Accomplished
- Bullet list of completed tasks

## Metrics & Results
- Test results
- Performance metrics
- Coverage
- Other relevant data

## Issues Encountered
- Any blockers
- Difficult problems solved
- Lessons learned

## Changes Made to Plan
- Any modifications to original plan
- Scope adjustments
- Timeline adjustments

## Quality Assessment
- Code quality (formatting, linting)
- Test quality
- Documentation quality

## Next Steps
- What to focus on in next milestone
- Known issues to address
- Recommendations

## Sign-off
Ready to proceed to next milestone: [YES / NO / NEEDS FIXES]
```

---

## Weekly Checkpoints

During Phase 1, use these weekly checkpoints to track progress:

### Week 1
- [ ] Milestone 1.1 Complete (Build & Compile)
- [ ] Milestone 1.2 Started (Testing & Benchmarking)

### Week 2
- [ ] Milestone 1.2 Complete (Testing & Benchmarking)
- [ ] Milestone 1.3 Started (Bug Fixes)
- [ ] Milestone 1.3 Complete or In Progress

### Week 3
- [ ] Milestone 1.3 Complete (Bug Fixes)
- [ ] Milestone 1.4 Complete (Performance)
- [ ] Milestone 1.5 In Progress (Test Coverage)

### Week 4
- [ ] Milestone 1.5 Complete (Test Coverage)
- [ ] Milestone 1.6 Complete (Documentation)
- [ ] **PHASE 1 COMPLETE** ✅

---

## Success Metrics

### Phase 1 Success Criteria (Final)

**Functionality**
- [x] Code is well-organized and readable
- [ ] Compiles cleanly without warnings
- [ ] Loads PDFs without crashing
- [ ] All vim keybindings work smoothly
- [ ] Dark mode works and looks good
- [ ] 3-pane layout displays correctly
- [ ] Search functionality works

**Performance**
- [ ] Cold start <200ms (Phase 1 relaxed target)
- [ ] Page load (cached) <100ms
- [ ] Page load (uncached) <300ms
- [ ] Memory usage <100MB (Phase 1 relaxed target)
- [ ] Search time <200ms
- [ ] No memory leaks

**Quality**
- [ ] 80%+ test coverage
- [ ] 100% test pass rate
- [ ] No race conditions
- [ ] Code formatted, vetted, linted
- [ ] No flaky tests

**Documentation**
- [ ] README complete
- [ ] Quick start guide complete
- [ ] Architecture documented
- [ ] Development guide complete
- [ ] API documented
- [ ] Keybindings documented

**Stability**
- [ ] Handles 20 diverse PDFs
- [ ] Graceful error handling
- [ ] No unhandled panics
- [ ] Tested on macOS/Linux

---

## Phase 1 Timeline

```
Week 1:
  Mon-Wed: Milestone 1.1 (Build & Compile)
  Thu-Fri: Milestone 1.2 (Testing)

Week 2:
  Mon-Tue: Milestone 1.2 (Testing continued)
  Wed-Fri: Milestone 1.3 (Bug Fixes)

Week 3:
  Mon-Wed: Milestone 1.4 (Performance)
  Thu-Fri: Milestone 1.5 (Test Coverage)

Week 4:
  Mon-Fri: Milestone 1.6 (Documentation)

Week 5 (if needed):
  Final polish and release preparation
```

**Total Duration**: 3-5 weeks (depending on blockers)

---

## How to Use This Plan

1. **Start with Milestone 1.1**
   - Read all tasks in 1.1
   - Work through them systematically
   - Update progress in todo list

2. **Complete Milestone 1.1**
   - Create `PHASE_1_MILESTONE_1_1_REVIEW.md`
   - Review what was accomplished
   - Note any issues

3. **Continue to Next Milestone**
   - Start with Milestone 1.2
   - Repeat the process
   - Pause at each review point for assessment

4. **Track Progress**
   - Update todo list as you work
   - Note actual time vs estimates
   - Document lessons learned

5. **Complete Phase 1**
   - All 6 milestones complete
   - Create `PHASE_1_COMPLETE.md`
   - Ready to start Phase 2

---

## Key Files to Create During Phase 1

### Review & Status Documents
- `PHASE_1_MILESTONE_1_1_REVIEW.md` - After milestone 1.1
- `PHASE_1_MILESTONE_1_2_REVIEW.md` - After milestone 1.2
- `PHASE_1_MILESTONE_1_3_REVIEW.md` - After milestone 1.3
- `PHASE_1_MILESTONE_1_4_REVIEW.md` - After milestone 1.4
- `PHASE_1_MILESTONE_1_5_REVIEW.md` - After milestone 1.5
- `PHASE_1_COMPLETE.md` - After all milestones

### Analysis & Documentation
- `BENCHMARKS_1_2.md` - Benchmark results
- `PERFORMANCE_BASELINE_1_2.md` - Performance baseline
- `PHASE_1_MILESTONE_1_3_FIXES.md` - Bug fixes summary
- `CACHE_ANALYSIS_1_4.md` - Cache performance analysis
- `PERFORMANCE_FINAL_1_4.md` - Final performance metrics
- `COVERAGE_1_5.md` - Test coverage report
- `TEST_GUIDE.md` - How to run tests
- `PHASE_1_SUMMARY.md` - Phase 1 overall summary
- `PHASE_2_PLAN.md` - Phase 2 roadmap

---

## Notes

- **Flexibility**: This is a guide. Adjust as needed based on what you discover.
- **Review Points**: Each milestone ends with a review point. Pause and assess before moving forward.
- **Documentation**: Create review files at each milestone for future reference.
- **Questions**: If you get stuck, document the issue and review it at the milestone review point.

---

**Plan Created**: 2025-10-21
**Status**: Ready to begin Milestone 1.1
**Next Action**: Start Milestone 1.1 - Build & Compile Foundation

