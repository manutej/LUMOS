# LUMOS Development Roadmap

**Project**: LUMOS - Dark Mode PDF Reader for Developers
**Version**: 1.0.0
**Last Updated**: 2025-11-10
**Status**: Phase 1 - 50% Complete

---

## Table of Contents

1. [Vision & Goals](#vision--goals)
2. [Current Status](#current-status)
3. [Phase 1: MVP](#phase-1-mvp-production-ready-reader)
4. [Phase 2: Enhanced Viewing](#phase-2-enhanced-viewing)
5. [Phase 3: Image Support](#phase-3-image-support)
6. [Phase 4: AI Integration](#phase-4-ai-integration)
7. [Success Metrics](#success-metrics)
8. [Risk Management](#risk-management)

---

## Vision & Goals

### North Star

**"The best dark mode PDF reader for developers who live in the terminal"**

### Core Principles

1. **Performance First**: Startup <100ms, navigation <50ms
2. **Vim Native**: Muscle memory compatibility for vim users
3. **Beautiful Dark Mode**: Optimized for long reading sessions
4. **Developer Friendly**: Clean architecture, extensible design
5. **Production Ready**: 80%+ test coverage, zero crashes

### Target Audience

- Developers reading technical papers, documentation, specifications
- Terminal power users who prefer keyboard navigation
- vim/neovim users wanting consistent keybindings
- Anyone needing a fast, distraction-free PDF reader

---

## Current Status

**Date**: November 10, 2025
**Phase**: 1 (MVP)
**Progress**: 50% complete (3/6 milestones done)

### Completed ✅

- ✅ **Phase 0**: Research & Architecture (Oct 21)
- ✅ **Milestone 1.1**: Build & Compile (Nov 1)
  - Clean build with Go 1.21+
  - Dependencies resolved
  - 4.6MB binary created

- ✅ **Milestone 1.2**: Core Testing (Nov 1)
  - 42 unit tests, 9 benchmarks
  - 70% code coverage
  - 5 bugs fixed
  - Performance validated

- ✅ **Milestone 1.3**: Test Fixtures (Nov 1)
  - 3 PDF fixtures generated
  - 94.4% code coverage achieved
  - Integration tests enabled
  - Search implementation complete

### In Progress 🚧

- 🚧 **Milestone 1.4**: Basic TUI Framework (Nov 10-13)
  - Status: Starting
  - Deliverable: Functional 3-pane layout with Bubble Tea

### Upcoming ⏳

- ⏳ **Milestone 1.5**: Vim Keybindings (Nov 13-16)
- ⏳ **Milestone 1.6**: Dark Mode Polish (Nov 16-19)
- ⏳ **Phase 2**: Enhanced Features (Dec 2025)
- ⏳ **Phase 3**: Image Support (Q1 2026)
- ⏳ **Phase 4**: AI Integration (Q2 2026)

---

## Phase 1: MVP (Production-Ready Reader)

**Timeline**: October 21 - November 19, 2025 (4 weeks)
**Status**: 50% Complete (3/6 milestones)
**Goal**: Functional dark mode PDF reader with vim keybindings

### Milestone Breakdown

#### ✅ Milestone 1.1: Build & Compile (Nov 1)
**Duration**: 1 day
**Status**: COMPLETE

**Deliverables**:
- [x] Clean build with `make build`
- [x] All dependencies resolved
- [x] Binary runs without errors
- [x] Makefile with test/lint commands

**Success Criteria**:
- [x] `go build` succeeds
- [x] Binary size <10MB (actual: 4.6MB)
- [x] No compile warnings
- [x] Documentation updated

---

#### ✅ Milestone 1.2: Core Testing (Nov 1)
**Duration**: 1 day
**Status**: COMPLETE

**Deliverables**:
- [x] 42 unit tests across pkg/pdf
- [x] 9 performance benchmarks
- [x] 70% code coverage
- [x] All tests passing

**Success Criteria**:
- [x] Cache operations <100ns ✅ (16ns actual)
- [x] Search operations <50μs ✅ (10μs actual)
- [x] 0 race conditions
- [x] All bugs fixed

**Key Achievements**:
- Found and fixed 5 bugs during testing
- Performance 5-10x better than targets
- High confidence in core PDF package

---

#### ✅ Milestone 1.3: Test Fixtures (Nov 1)
**Duration**: 1 day
**Status**: COMPLETE

**Deliverables**:
- [x] simple.pdf (1 page)
- [x] multipage.pdf (5 pages)
- [x] search_test.pdf (search patterns)
- [x] Comprehensive testing guide
- [x] 94.4% code coverage

**Success Criteria**:
- [x] All integration tests enabled
- [x] >80% coverage achieved ✅ (94.4%)
- [x] Real PDF loading validated
- [x] Search functionality working

**Key Achievements**:
- Coverage jumped from 70% → 94.4%
- Integration tests revealed stub functions
- Testing guide for future development

---

#### 🚧 Milestone 1.4: Basic TUI Framework (Nov 10-13)
**Duration**: 2-3 days
**Status**: STARTING
**Specification**: `.specify/specs/milestone-1.4-tui-framework.md`

**Deliverables**:
- [ ] Bubble Tea TUI initialized
- [ ] 3-pane layout (Metadata 20% | Viewer 60% | Search 20%)
- [ ] Viewport component integrated
- [ ] Status bar with page numbers
- [ ] Basic keyboard navigation (q to quit)
- [ ] Window resize handling

**Success Criteria**:
- [ ] Startup time <70ms
- [ ] UI renders at 60 FPS
- [ ] 3 terminals tested (iTerm2, Terminal.app, Alacritty)
- [ ] 15 unit tests + 5 integration tests
- [ ] >80% coverage for pkg/ui
- [ ] Clean shutdown, no goroutine leaks

**Key Tasks**:
1. Initialize tea.Program in main.go
2. Create Model with MVU pattern
3. Implement 3-pane layout with lipgloss
4. Add viewport for scrollable content
5. Create status bar component
6. Handle window resize messages
7. Test on multiple terminals

**Dependencies**:
- Upstream: Milestone 1.3 ✅
- Downstream: Milestone 1.5 (needs working viewport)

---

#### ⏳ Milestone 1.5: Vim Keybindings (Nov 13-16)
**Duration**: 2-3 days
**Status**: NOT STARTED
**Specification**: `.specify/specs/milestone-1.5-vim-keybindings.md`

**Deliverables**:
- [ ] All vim navigation keys (j/k/d/u/gg/G/Ctrl+N/Ctrl+P)
- [ ] Multi-key sequences (gg for first page)
- [ ] Search mode (/ to search, n/N for next/prev)
- [ ] Help screen (? key)
- [ ] Pane cycling (Tab key)
- [ ] Theme toggle (1/2 keys)

**Success Criteria**:
- [ ] All keybindings work <16ms response
- [ ] vim users confirm "feels right"
- [ ] Search <100ms for typical PDFs
- [ ] Multi-key timeout exactly 500ms
- [ ] 25 unit tests + 10 integration tests
- [ ] >85% coverage for ui package

**Key Tasks**:
1. Create KeyHandler with mode tracking
2. Implement multi-key sequence buffer
3. Add all vim navigation keys
4. Implement search mode (/  to enter)
5. Create help screen overlay
6. Add pane cycling with Tab
7. Vim user testing session

**Dependencies**:
- Upstream: Milestone 1.4
- Downstream: Milestone 1.6 (needs full nav for testing)

---

#### ⏳ Milestone 1.6: Dark Mode Polish (Nov 16-19)
**Duration**: 2-3 days
**Status**: NOT STARTED
**Specification**: `.specify/specs/milestone-1.6-dark-mode-polish.md`

**Deliverables**:
- [ ] Refined dark theme (Tokyo Night inspired)
- [ ] Refined light theme (high contrast)
- [ ] Smooth theme transitions
- [ ] Visual feedback for all actions
- [ ] Performance optimization
- [ ] Complete documentation
- [ ] Phase 1 completion review

**Success Criteria**:
- [ ] Contrast ratio >7:1 (WCAG AAA)
- [ ] All performance targets met
- [ ] Test coverage >80% overall
- [ ] 0 linter warnings
- [ ] All documentation complete
- [ ] **Phase 1 is production-ready**

**Key Tasks**:
1. Refine color palettes for both themes
2. Test on 3 terminal emulators
3. Add visual feedback (animations, indicators)
4. Profile and optimize performance
5. Final code quality pass
6. Complete all documentation
7. Create demo materials
8. Draft Phase 2 plan

**Dependencies**:
- Upstream: Milestone 1.5
- Downstream: Phase 1 complete → Phase 2 start

---

### Phase 1 Success Metrics

#### Functional Requirements
- [x] Opens PDF files
- [x] Extracts text from all pages
- [ ] Displays in 3-pane TUI layout
- [ ] All vim keybindings work
- [ ] Search finds matches
- [ ] Dark mode by default
- [ ] Light mode available
- [ ] Help screen accessible
- [ ] Clean quit (q/Ctrl+C)

#### Performance Requirements
- [x] Cold start <100ms (actual: ~70ms)
- [x] Page switch (cached) <50ms (actual: <20ms)
- [x] Page switch (uncached) <200ms (actual: ~65ms)
- [x] Memory baseline <10MB
- [x] Memory with PDF <50MB
- [x] Search <100ms per 100 pages
- [ ] 60 FPS UI rendering

#### Quality Requirements
- [x] Test coverage >80% (actual: 94.4%)
- [x] All tests passing (42/42)
- [x] 0 linter warnings
- [ ] 60+ total tests
- [ ] 0 TODOs in production code
- [ ] All exported symbols documented

#### Documentation Requirements
- [x] README.md complete
- [x] ARCHITECTURE.md current
- [x] PROGRESS.md updated
- [x] Testing guide created
- [ ] All milestone reviews complete
- [ ] Phase 1 completion document
- [ ] Phase 2 plan drafted

---

## Phase 2: Enhanced Viewing

**Timeline**: December 2025 - January 2026 (4-6 weeks)
**Status**: PLANNED
**Goal**: Professional-grade PDF viewer with advanced features

### Planned Features

#### 2.1: Advanced Search
- Regex search patterns
- Whole-word matching
- Case-sensitive toggle
- Search within page range
- Search history
- **Success**: Complex searches <200ms

#### 2.2: Table of Contents
- Extract PDF TOC if available
- Navigate by TOC entries
- Bookmarks within viewer
- Quick jump to sections
- **Success**: TOC extraction <50ms

#### 2.3: Vim Marks & Jumps
- Mark positions with `m{a-z}`
- Jump to marks with `'{a-z}`
- Global marks across sessions
- Jump list (`Ctrl+O/Ctrl+I`)
- **Success**: vim users feel at home

#### 2.4: Better Text Layout
- Preserve paragraph structure
- Handle columns correctly
- Respect reading order
- Better whitespace handling
- **Success**: 90% of PDFs readable

#### 2.5: Configuration System
- Config file: `~/.config/lumos/config.toml`
- Theme customization
- Keybinding remapping
- Default settings
- **Success**: User preferences persist

### Phase 2 Success Metrics
- [ ] All Phase 1 features stable
- [ ] TOC navigation works
- [ ] Vim marks functional
- [ ] Search supports regex
- [ ] Config system implemented
- [ ] Test coverage >85%
- [ ] User feedback positive

---

## Phase 3: Image Support

**Timeline**: Q1 2026 (6-8 weeks)
**Status**: RESEARCH
**Goal**: Render PDF images in terminal

### Planned Features

#### 3.1: Terminal Graphics Detection
- Auto-detect Kitty protocol
- Auto-detect iTerm2 inline images
- Auto-detect SIXEL support
- Fallback to text-only
- **Success**: Works on 3 terminals

#### 3.2: Image Rendering
- Extract images from PDF
- Convert to terminal format
- Render inline with text
- Image caching (LRU)
- **Success**: Images display correctly

#### 3.3: Hybrid Rendering
- Mix text and images
- Handle image placement
- Scroll with images
- Performance optimization
- **Success**: <300ms image render

#### 3.4: Image Quality
- Resolution optimization
- Color depth handling
- Dithering for terminals
- Aspect ratio preservation
- **Success**: Images recognizable

### Phase 3 Success Metrics
- [ ] Images render in 3 terminals
- [ ] Hybrid text+image layout works
- [ ] Image cache effective (>80% hit rate)
- [ ] Performance acceptable (<300ms)
- [ ] No crashes with image-heavy PDFs
- [ ] Test coverage >80%

### Technologies
- `github.com/blacktop/go-termimg` - Terminal image rendering
- Kitty graphics protocol
- iTerm2 inline images
- SIXEL fallback

---

## Phase 4: AI Integration

**Timeline**: Q2 2026 (8-12 weeks)
**Status**: CONCEPT
**Goal**: NotebookLM-like AI features for PDFs

### Planned Features

#### 4.1: Claude Agent SDK Integration
- Integrate Claude SDK (Go or API)
- PDF content chunking
- Context management
- API key configuration
- **Success**: SDK integrated, authenticated

#### 4.2: PDF Q&A
- `/ask` command for questions
- Contextual answers from PDF
- Citation to page numbers
- Multi-turn conversation
- **Success**: Useful answers in <5s

#### 4.3: Audio Generation
- TTS from PDF text
- Page-by-page audio
- Adjustable speed/voice
- Export to audio file
- **Success**: Clear audio output

#### 4.4: Summary Generation
- Automatic PDF summaries
- Section summaries
- Key points extraction
- Export summaries to markdown
- **Success**: Accurate summaries

#### 4.5: Multi-Document Analysis
- Compare multiple PDFs
- Find relationships
- Cross-reference citations
- Knowledge graph visualization
- **Success**: Useful insights

### Phase 4 Success Metrics
- [ ] Claude SDK integrated
- [ ] Q&A provides useful answers
- [ ] Audio generation works
- [ ] Summaries are accurate
- [ ] Multi-doc analysis insightful
- [ ] Privacy controls in place
- [ ] Test coverage >75%

### Technologies
- Claude Agent SDK (Go wrapper or API)
- OpenAI Whisper (for audio)
- Markdown export for summaries

---

## Success Metrics

### Phase 1 (MVP) - Target: Nov 19, 2025

#### Critical Success Factors
1. **Functionality**: All core features work
   - Load PDFs ✅
   - Display text ✅
   - Navigate with vim keys (in progress)
   - Search content ✅
   - Dark mode (in progress)

2. **Performance**: Meets all targets
   - Startup <100ms ✅ (70ms)
   - Navigation <50ms ✅ (20ms cached)
   - Memory <50MB ✅ (8MB)
   - Search <100ms ✅ (10μs/KB)

3. **Quality**: Production-ready code
   - Test coverage >80% ✅ (94.4%)
   - All tests passing ✅ (42/42)
   - 0 linter warnings ✅
   - Clean architecture ✅

4. **Documentation**: Complete and current
   - README ✅
   - ARCHITECTURE ✅
   - PROGRESS ✅
   - Milestone reviews (3/6)

#### Key Performance Indicators (KPIs)

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| Milestones Complete | 6/6 | 3/6 | 🟡 50% |
| Test Coverage | >80% | 94.4% | 🟢 Exceeded |
| Tests Passing | 100% | 100% | 🟢 Perfect |
| Startup Time | <100ms | ~70ms | 🟢 30% better |
| Page Switch (cached) | <50ms | <20ms | 🟢 60% better |
| Memory Usage | <50MB | ~8MB | 🟢 84% better |
| Linter Warnings | 0 | 0 | 🟢 Clean |

### Long-Term Success (Phase 2-4)

#### User Adoption
- **Phase 1**: 100 users
- **Phase 2**: 1,000 users
- **Phase 3**: 5,000 users
- **Phase 4**: 10,000+ users

#### Community Engagement
- **Phase 1**: GitHub repo public
- **Phase 2**: 100 GitHub stars
- **Phase 3**: 500 GitHub stars
- **Phase 4**: 1,000+ GitHub stars

#### Quality Metrics
- **Crash rate**: <0.1%
- **Performance**: 95th percentile meets targets
- **Test coverage**: >80% all phases
- **User satisfaction**: >4.5/5.0

---

## Risk Management

### High-Priority Risks

#### R1: Performance Degradation
**Risk**: Large PDFs cause slowdowns
**Impact**: High - Affects core UX
**Mitigation**:
- Aggressive caching strategy (LRU)
- Lazy page loading
- Performance benchmarks in CI
- Profile before optimizing
**Status**: Mitigated (cache working well)

#### R2: Terminal Compatibility
**Risk**: UI breaks on some terminals
**Impact**: Medium - Limits user base
**Mitigation**:
- Test on iTerm2, Terminal.app, Alacritty
- Graceful degradation for unsupported features
- Clear system requirements
**Status**: Testing in progress

#### R3: PDF Parsing Limitations
**Risk**: Some PDFs don't parse correctly
**Impact**: Medium - User frustration
**Mitigation**:
- Use proven library (ledongthuc/pdf)
- Handle errors gracefully
- Show helpful error messages
- Test with diverse PDFs
**Status**: Monitoring

### Medium-Priority Risks

#### R4: Scope Creep
**Risk**: Adding features delays Phase 1
**Impact**: Medium - Timeline slip
**Mitigation**:
- Strict milestone adherence
- Feature freeze after 1.6
- Document deferred features
**Status**: Controlled

#### R5: Image Support Complexity
**Risk**: Phase 3 harder than expected
**Impact**: Medium - Phase delay
**Mitigation**:
- Extensive research (already done)
- go-termimg library simplifies
- Can defer if needed
**Status**: Planned

### Low-Priority Risks

#### R6: API Changes
**Risk**: Bubble Tea API breaks
**Impact**: Low - Fixable
**Mitigation**:
- Pin exact versions
- Regular dependency updates
- Comprehensive tests
**Status**: Monitored

---

## Timeline Summary

```
Phase 0: Research & Design         ✅ Oct 21
├─ PDF library research
├─ TUI framework selection
└─ Architecture design

Phase 1: MVP                        🚧 Oct 21 - Nov 19 (4 weeks)
├─ 1.1: Build & Compile            ✅ Nov 1
├─ 1.2: Core Testing               ✅ Nov 1
├─ 1.3: Test Fixtures              ✅ Nov 1
├─ 1.4: Basic TUI                  🚧 Nov 10-13
├─ 1.5: Vim Keybindings            ⏳ Nov 13-16
└─ 1.6: Dark Mode Polish           ⏳ Nov 16-19

Phase 2: Enhanced Viewing           ⏳ Dec 2025 - Jan 2026 (6 weeks)
├─ Advanced search (regex)
├─ Table of contents
├─ Vim marks & jumps
├─ Better text layout
└─ Configuration system

Phase 3: Image Support              ⏳ Q1 2026 (8 weeks)
├─ Terminal graphics detection
├─ Image rendering
├─ Hybrid text+image layout
└─ Performance optimization

Phase 4: AI Integration             ⏳ Q2 2026 (12 weeks)
├─ Claude SDK integration
├─ PDF Q&A
├─ Audio generation
├─ Summary generation
└─ Multi-document analysis
```

---

## Next Actions

### Immediate (This Week)
1. **Complete Milestone 1.4** (Nov 10-13)
   - Implement Bubble Tea TUI
   - Create 3-pane layout
   - Add viewport component
   - Test on 3 terminals

2. **Start Milestone 1.5** (Nov 13)
   - Design key handler
   - Plan vim keybindings
   - Prepare test cases

### Short-Term (Next 2 Weeks)
3. **Complete Milestone 1.5** (Nov 13-16)
   - Implement all vim keys
   - Add search mode
   - Create help screen

4. **Complete Milestone 1.6** (Nov 16-19)
   - Polish themes
   - Optimize performance
   - Finalize documentation
   - **Ship Phase 1!** 🚀

### Medium-Term (Next Month)
5. **Plan Phase 2** (Dec 2025)
   - Gather user feedback
   - Prioritize features
   - Design TOC system
   - Spec out config system

---

## Constitutional Compliance

All development MUST comply with `.specify/constitution.md`:

- ✅ Code quality: go fmt, go vet, golangci-lint clean
- ✅ Testing: >80% coverage, TDD approach
- ✅ Performance: All budgets met
- ✅ Architecture: MVU pattern, package separation
- ✅ Documentation: All symbols documented

Non-compliance blocks merge. No exceptions.

---

## Conclusion

LUMOS is on track to become the best dark mode PDF reader for developers. Phase 1 is 50% complete with strong fundamentals:

- **Solid Core**: PDF package battle-tested (94.4% coverage)
- **Performance**: All targets exceeded (30-84% better)
- **Quality**: Zero-warning, well-tested codebase
- **Clear Path**: Detailed specs for remaining milestones

**Next Milestone**: 1.4 - Basic TUI Framework (Starting Nov 10)

**Phase 1 Ship Date**: November 19, 2025 ✅

---

**Roadmap Version**: 1.0.0
**Last Updated**: 2025-11-10
**Next Review**: Upon Phase 1 Completion
**Maintained By**: LUMOS Development Team
