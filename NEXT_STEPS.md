# LUMOS Next Steps

**Last Updated**: 2025-11-18
**Current Phase**: Phase 1 MVP (83% complete)
**Status**: Functionally complete, pending test coverage & polish

---

## 🎯 Immediate Priorities (In Order)

### 1. ⚠️ CRITICAL: Increase UI Test Coverage (4-6 hours)
**Current**: 47.8% | **Target**: 80%+ | **Gap**: ~15-20 tests needed

**Why Critical**: Quality gate blocker for Phase 1 completion

**Tasks**:
- [ ] Add rendering function tests
  - [ ] `TestRenderMetadataPane_DisplaysCorrectInfo`
  - [ ] `TestRenderViewerPane_ShowsPageNumber`
  - [ ] `TestRenderSearchPane_ShowsResultCount`
  - [ ] `TestRenderStatusBar_FormatsCorrectly`
  - [ ] `TestRenderHelp_ShowsKeybindings`

- [ ] Add layout calculation tests
  - [ ] `TestCalculatePaneWidths_80Columns`
  - [ ] `TestCalculatePaneWidths_200Columns`
  - [ ] `TestCalculatePaneWidths_MinimumWidth`
  - [ ] `TestHandleWindowResize_UpdatesViewport`

- [ ] Add search functionality tests
  - [ ] `TestExecuteSearch_FindsMatches`
  - [ ] `TestJumpToSearchResult_NavigatesToPage`
  - [ ] `TestSearchMode_EnterAndExit`
  - [ ] `TestNextPreviousMatch_Cycles`

- [ ] Add KeyHandler tests
  - [ ] `TestKeyHandler_ModeSwitch_NormalToSearch`
  - [ ] `TestKeyHandler_ModeSwitch_SearchToNormal`
  - [ ] `TestKeyHandler_HandleSearchInput`

**Files to Edit**:
- `pkg/ui/model_test.go` - Add new test cases
- `pkg/ui/layout_test.go` - Create new file for layout tests
- `pkg/ui/keybindings_test.go` - Create new file for keybinding tests

**Success Criteria**:
- ✅ `go test ./pkg/ui -cover` shows ≥80%
- ✅ All tests passing
- ✅ No race conditions (`go test -race`)

**Estimated Time**: 4-6 hours

---

### 2. 🔍 HIGH: Performance Verification (1 hour)
**Status**: Not yet measured | **Target**: Meet all targets

**Tasks**:
- [ ] Measure startup time
  ```bash
  time ~/bin/lumos test/fixtures/simple.pdf
  # Target: <70ms
  ```

- [ ] Measure render performance
  ```bash
  # Use profiling
  make profile-cpu
  # Verify: <16ms per frame (60 FPS)
  ```

- [ ] Verify memory usage
  ```bash
  # Run for 5 minutes, monitor memory
  # Target: <10MB baseline
  ```

- [ ] Document results
  - [ ] Add to PROGRESS.md
  - [ ] Create benchmark in `pkg/ui/model_bench_test.go`

**Files to Create**:
- `pkg/ui/model_bench_test.go` - Benchmark tests

**Success Criteria**:
- ✅ Startup: <70ms
- ✅ Render: <16ms
- ✅ Memory: <10MB
- ✅ Results documented

**Estimated Time**: 1 hour

---

### 3. 🎨 MEDIUM: Dark Mode Polish (2-3 hours)
**Status**: 50% complete | **Target**: Theme refinement

**Tasks**:
- [ ] Refine color contrast
  - [ ] Test readability on different backgrounds
  - [ ] Adjust border colors for clarity
  - [ ] Verify WCAG contrast ratios

- [ ] Optimize border styles
  - [ ] Test rounded vs square borders
  - [ ] Adjust padding for comfort
  - [ ] Ensure alignment across panes

- [ ] Terminal compatibility
  - [ ] Test on iTerm2 (primary)
  - [ ] Test on Terminal.app
  - [ ] Test on Alacritty
  - [ ] Document any quirks

- [ ] Light mode verification
  - [ ] Test light mode colors
  - [ ] Ensure readability
  - [ ] Document usage

**Files to Edit**:
- `pkg/config/theme.go` - Adjust colors and styles

**Success Criteria**:
- ✅ Readable on all 3 terminals
- ✅ Good contrast ratios
- ✅ Both themes usable

**Estimated Time**: 2-3 hours

---

## 📊 Phase 1 Completion Checklist

### Milestones Status
- [x] 1.1: Build & Compile (100%)
- [x] 1.2: Core Testing (100%)
- [x] 1.3: Test Fixtures (100%)
- [x] 1.4: TUI Framework (95% - functionally complete)
- [x] 1.5: Vim Keybindings (95% - functionally complete)
- [ ] 1.6: Dark Mode Polish (50% - in progress)

### Quality Gates
- [x] All builds pass
- [x] PDF package tests: 94.4%
- [ ] UI package tests: 47.8% → 80%+ (BLOCKER)
- [ ] Performance targets verified (BLOCKER)
- [ ] Terminal compatibility verified
- [x] Documentation complete
- [x] README updated
- [x] Installation working

### Final Tasks Before Release
- [ ] Create Phase 1 completion review
- [ ] Update README with screenshots (text art)
- [ ] Tag v0.1.0 release
- [ ] Push to remote repository
- [ ] Celebrate! 🎉

---

## 🚀 After Phase 1 (Future Work)

### Phase 2: Enhanced Viewing (Q1 2026)
- Table of contents navigation
- Bookmarks and annotations
- Better search (regex, whole word)
- Page thumbnails sidebar
- Multi-document tabs

### Phase 3: Image Support
- Terminal image rendering with go-termimg
- Image extraction and display
- Fallback for non-graphical terminals

### Phase 4: AI Integration
- Claude Agent SDK integration
- NotebookLM-like features
- Document Q&A
- Summarization

---

## 📝 Session-by-Session Plan

### Session 1 (Today/Next): UI Tests
**Duration**: 4-6 hours
**Goal**: Reach 80%+ UI test coverage

**Workflow**:
1. Start with rendering function tests (2 hours)
2. Add layout calculation tests (1 hour)
3. Add search functionality tests (1-2 hours)
4. Add KeyHandler tests (1 hour)
5. Verify coverage: `go test ./pkg/ui -cover`
6. Fix any issues found
7. Update PROGRESS.md

### Session 2: Performance & Polish
**Duration**: 3-4 hours
**Goal**: Verify performance, polish themes

**Workflow**:
1. Measure startup, render, memory (1 hour)
2. Refine dark mode colors (1 hour)
3. Test on 3 terminals (1 hour)
4. Document results (30 min)
5. Update PROGRESS.md

### Session 3: Phase 1 Completion
**Duration**: 1-2 hours
**Goal**: Final review and release

**Workflow**:
1. Create Phase 1 completion review
2. Update all documentation
3. Tag v0.1.0 release
4. Push to remote
5. Plan Phase 2 kickoff

---

## 🛠️ Quick Commands Reference

### Development
```bash
# Run UI tests with coverage
go test -v -cover ./pkg/ui

# Run tests with race detector
go test -race ./...

# Profile CPU
make profile-cpu

# Profile memory
make profile-mem

# Build and install
make install

# Test TUI (after install)
lumos test/fixtures/multipage.pdf
```

### Git Workflow
```bash
# Check status
git status

# Run tests before commit
make test

# Commit with spec reference
git commit -m "test: add UI rendering tests

Increases pkg/ui coverage from 47.8% to 65%
Adds 10 tests for rendering functions

Per spec .specify/TESTING_STRATEGY.md"

# Push when ready
git push origin master
```

---

## 📍 Current Blockers

### BLOCKER #1: UI Test Coverage
**Severity**: High
**Impact**: Can't complete Phase 1 without 80%+ coverage
**Owner**: Next developer session
**ETA**: 4-6 hours

### BLOCKER #2: Performance Verification
**Severity**: Medium
**Impact**: Need to verify targets met
**Owner**: After test coverage complete
**ETA**: 1 hour

---

## 💡 Tips for Next Session

### Starting Fresh
1. **Read this file first** - Start here every session
2. **Check PROGRESS.md** - See detailed history
3. **Review `.specify/HANDOFF.md`** - 15-minute orientation
4. **Run tests** - `make test` to verify clean state

### Working on UI Tests
- **Reference existing tests** in `pkg/pdf/*_test.go` for patterns
- **Use table-driven tests** for multiple scenarios
- **Mock dependencies** when needed (document, viewport)
- **Test both success and error cases**
- **Run tests frequently** - `go test -v ./pkg/ui`

### Staying Organized
- **Update NEXT_STEPS.md** after each session
- **Update PROGRESS.md** with detailed notes
- **Commit frequently** with clear messages
- **Reference specs** in commits (`.specify/`)

---

## 📞 Getting Help

### Documentation
- **Quick Start**: `.specify/HANDOFF.md`
- **Dependencies**: `.claude/DEPENDENCIES.md`
- **Architecture**: `.specify/ARCHITECTURE_DECISIONS.md`
- **Testing**: `.specify/TESTING_STRATEGY.md`
- **Complete Specs**: `.specify/SPECIFICATION_INDEX.md`

### Code References
- **TUI Examples**: `pkg/ui/model.go` - Well-documented MVU pattern
- **Test Examples**: `pkg/pdf/cache_test.go` - Comprehensive test patterns
- **Keybindings**: `pkg/ui/keybindings.go` - Modal editing system

### External Resources
- **Bubble Tea**: https://github.com/charmbracelet/bubbletea
- **Testing Guide**: https://go.dev/doc/tutorial/add-a-test
- **Coverage Tools**: https://go.dev/blog/cover

---

## ✅ Definition of Done

Phase 1 is **DONE** when:
1. ✅ All 6 milestones complete
2. ✅ Test coverage >80% for all packages
3. ✅ Performance targets met and documented
4. ✅ Terminal compatibility verified
5. ✅ All documentation updated
6. ✅ Installation working globally
7. ✅ Tagged v0.1.0 release
8. ✅ Pushed to remote

**Current Status**: 6/8 criteria met, 2 in progress

**Estimated Completion**: 8-11 hours (2-3 focused sessions)

---

**Remember**: The TUI is already working! We're in polish and quality assurance mode, not build mode. This is the exciting final stretch! 🎯

---

**Last Updated**: 2025-11-18
**Next Review**: After UI test coverage session
**Maintained By**: Claude Code + Human Developer
