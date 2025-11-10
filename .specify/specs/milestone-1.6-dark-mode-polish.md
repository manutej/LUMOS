# Milestone 1.6: Dark Mode Polish & Phase 1 Completion

**Status**: Not Started
**Priority**: P0 (Phase 1 Completion)
**Estimated Duration**: 2-3 days
**Dependencies**: Milestone 1.5 (Vim Keybindings)

---

## Overview

Apply final polish to dark mode theme, optimize performance, achieve production-ready quality, and prepare for Phase 1 release. This milestone focuses on user experience refinement and quality assurance.

## User Stories

### US-1.6.1: As a user, I experience a beautiful dark mode by default
**Given** I launch LUMOS for the first time
**When** The interface renders
**Then** I see a carefully crafted dark theme that's easy on the eyes for long reading sessions

### US-1.6.2: As a user, I can switch to light mode when needed
**Given** I'm in dark mode
**When** I press `2`
**Then** The interface smoothly transitions to light mode with excellent readability

### US-1.6.3: As a user, I see consistent styling across all UI elements
**Given** I'm navigating the interface
**When** I switch panes, modes, or themes
**Then** All elements maintain visual consistency and professional appearance

### US-1.6.4: As a developer, I experience blazing fast performance
**Given** I'm using LUMOS with typical PDFs
**When** I navigate, search, and scroll
**Then** The interface is responsive with no lag or jank

### US-1.6.5: As a user, I have comprehensive documentation
**Given** I want to learn LUMOS
**When** I read the documentation
**Then** I find clear guides, examples, and troubleshooting help

---

## Acceptance Criteria

### AC-1.6.1: Dark Theme Refinement
```go
// pkg/config/theme.go - DarkTheme colors
Background:     "#1a1b26"  // Tokyo Night background
Foreground:     "#c0caf5"  // Light blue-white text
PaneBorder:     "#565f89"  // Subtle blue-gray
ActiveBorder:   "#7aa2f7"  // Bright blue for active
StatusBar:      "#1f2335"  // Slightly lighter than bg
StatusText:     "#9ece6a"  // Green for highlights
ErrorText:      "#f7768e"  // Soft red for errors
SearchMatch:    "#e0af68"  // Gold for search highlights
HelpOverlay:    "#24283b"  // Semi-transparent overlay
```

- [ ] All colors defined with hex codes
- [ ] Contrast ratio >7:1 for text (WCAG AAA)
- [ ] Colors tested on multiple terminals (iTerm2, Terminal.app, Alacritty)
- [ ] Consistent color usage across all panes

### AC-1.6.2: Light Theme Refinement
```go
// pkg/config/theme.go - LightTheme colors
Background:     "#ffffff"  // Pure white
Foreground:     "#383a42"  // Dark gray text
PaneBorder:     "#d0d0d0"  // Light gray
ActiveBorder:   "#4078f2"  // Bright blue
StatusBar:      "#f0f0f0"  // Off-white
StatusText:     "#50a14f"  // Green
ErrorText:      "#e45649"  // Red
SearchMatch:    "#c18401"  // Orange
HelpOverlay:    "#fafafa"  // Very light gray
```

- [ ] All colors defined
- [ ] Contrast ratio >7:1 for text
- [ ] No eye strain in bright environments
- [ ] Smooth theme transition (no flash)

### AC-1.6.3: Typography & Spacing
```go
// Lipgloss styles
TitleStyle:     Bold, Foreground(ActiveBorder)
BodyStyle:      Normal, Foreground(Foreground)
StatusStyle:    Background(StatusBar), Foreground(StatusText)
BorderStyle:    RoundedBorder, BorderForeground(PaneBorder)

// Spacing
PanePadding:    1 (all sides)
StatusPadding:  0 vertical, 2 horizontal
HelpPadding:    2 (all sides)
```

- [ ] Consistent padding across all panes
- [ ] Proper line height for readability
- [ ] Clear visual hierarchy
- [ ] No text cutoff or overflow

### AC-1.6.4: Visual Feedback
- [ ] Active pane clearly highlighted
- [ ] Page transitions have subtle animation (border flash)
- [ ] Loading states shown for operations >100ms
- [ ] Error messages clearly visible but not alarming
- [ ] Status bar updates immediately on state change
- [ ] Search matches highlighted prominently

### AC-1.6.5: Performance Optimization
```bash
# Startup Performance
time ./build/lumos test/fixtures/simple.pdf
# Target: <70ms

# Memory Profiling
go test -memprofile=mem.prof -bench=. ./pkg/...
# Target: 0 allocations in hot paths

# CPU Profiling
go test -cpuprofile=cpu.prof -bench=. ./pkg/...
# Target: No inefficient operations in critical paths
```

- [ ] Cold start <70ms
- [ ] Page switch (cached) <20ms
- [ ] Page switch (uncached) <65ms
- [ ] Search <100ms for 100-page PDF
- [ ] Memory <10MB baseline, <50MB with PDF
- [ ] No memory leaks during extended use

### AC-1.6.6: Code Quality Final Pass
```bash
# All quality checks must pass
make fmt          # No changes needed
make vet          # 0 warnings
make lint         # 0 issues
make test         # 100% pass rate
make coverage     # >80% overall
make ci-check     # All checks pass
```

- [ ] All code formatted with gofmt
- [ ] All linter warnings resolved
- [ ] No TODO or FIXME comments in production code
- [ ] All exported symbols documented
- [ ] Consistent error handling patterns

### AC-1.6.7: Documentation Completeness
**README.md**:
- [ ] Clear project description
- [ ] Installation instructions
- [ ] Usage examples with screenshots (text art)
- [ ] Keybindings reference
- [ ] Troubleshooting section

**ARCHITECTURE.md**:
- [ ] Updated with final architecture
- [ ] All diagrams current
- [ ] Performance characteristics documented

**PROGRESS.md**:
- [ ] All milestones marked complete
- [ ] Final metrics documented
- [ ] Lessons learned captured

**Phase 1 Completion**:
- [ ] Create PHASE_1_COMPLETE.md
- [ ] Document achievements
- [ ] Note deviations from plan
- [ ] Outline Phase 2 preparation

---

## Implementation Checklist

### Phase 1: Theme Refinement (Day 1 Morning)
- [ ] Review and optimize dark theme colors:
  - [ ] Test on iTerm2
  - [ ] Test on Terminal.app
  - [ ] Test on Alacritty
  - [ ] Adjust for consistency

- [ ] Review and optimize light theme colors:
  - [ ] Ensure readability
  - [ ] Test in bright environments
  - [ ] Verify contrast ratios

- [ ] Create theme preview documentation:
  - [ ] Screenshots (text art) of both themes
  - [ ] Color palette reference
  - [ ] Usage guidelines

### Phase 2: Visual Polish (Day 1 Afternoon)
- [ ] Enhance status bar:
  - [ ] Better layout
  - [ ] Clearer information hierarchy
  - [ ] Add helpful hints based on context

- [ ] Improve pane borders:
  - [ ] Active pane clearly highlighted
  - [ ] Inactive panes subtle but visible
  - [ ] Consistent border styling

- [ ] Add visual feedback:
  - [ ] Page transition animation
  - [ ] Loading indicators
  - [ ] Error state styling

- [ ] Polish help screen:
  - [ ] Better layout
  - [ ] Clearer sections
  - [ ] Examples for complex keys

### Phase 3: Performance Optimization (Day 2 Morning)
- [ ] Profile startup time:
  - [ ] Identify bottlenecks
  - [ ] Optimize initialization
  - [ ] Lazy load where possible

- [ ] Profile memory usage:
  - [ ] Check for leaks
  - [ ] Optimize cache usage
  - [ ] Minimize allocations

- [ ] Profile rendering:
  - [ ] Optimize View() function
  - [ ] Minimize lipgloss calls
  - [ ] Cache computed styles

- [ ] Run benchmarks:
  - [ ] Document all results
  - [ ] Compare against targets
  - [ ] Optimize if needed

### Phase 4: Quality Assurance (Day 2 Afternoon)
- [ ] Code review:
  - [ ] Check all packages
  - [ ] Remove dead code
  - [ ] Improve comments
  - [ ] Ensure consistency

- [ ] Testing:
  - [ ] Run full test suite
  - [ ] Check coverage gaps
  - [ ] Add missing tests
  - [ ] Verify edge cases

- [ ] Linting:
  - [ ] Fix all warnings
  - [ ] Apply suggestions
  - [ ] Format all files

### Phase 5: Documentation & Release Prep (Day 3)
- [ ] Update README.md:
  - [ ] Add usage examples
  - [ ] Create text art screenshots
  - [ ] Document installation
  - [ ] Add troubleshooting

- [ ] Finalize technical docs:
  - [ ] Update ARCHITECTURE.md
  - [ ] Complete all milestone reviews
  - [ ] Update PROGRESS.md
  - [ ] Create PHASE_1_COMPLETE.md

- [ ] Create Phase 2 plan:
  - [ ] Outline features
  - [ ] Set timeline
  - [ ] Identify dependencies
  - [ ] Create PHASE_2_PLAN.md

- [ ] Prepare demo:
  - [ ] Create demo PDF
  - [ ] Record usage session (script)
  - [ ] Prepare presentation notes

---

## Test Requirements

### Final Test Suite (Target: 60+ tests total)
```bash
# Package breakdown
pkg/pdf/         # 26 tests (existing)
pkg/ui/          # 30 tests (15 existing + 15 new)
pkg/config/      # 5 tests (new)
test/            # 5 integration tests (expand existing)

# Total: 66 tests
```

### New Tests for 1.6
```go
// pkg/config/theme_test.go
TestDarkTheme_ContrastRatio_ExceedsWCAG
TestLightTheme_ContrastRatio_ExceedsWCAG
TestThemeTransition_NoFlash
TestThemeColors_ValidHexCodes
TestThemeApplication_AllPanes

// pkg/ui/styling_test.go
TestRenderStatusBar_DarkTheme_CorrectColors
TestRenderStatusBar_LightTheme_CorrectColors
TestRenderPanes_ActiveHighlight
TestRenderPanes_InactiveSubtle
TestRenderHelp_ReadableLayout

// test/test_performance.go
TestStartup_UnderTarget
TestPageSwitch_UnderTarget
TestMemory_UnderTarget
TestSearch_UnderTarget
TestExtendedUse_NoMemoryLeak
```

### Manual Testing Checklist
- [ ] Dark mode on iTerm2 - all features work
- [ ] Dark mode on Terminal.app - all features work
- [ ] Dark mode on Alacritty - all features work
- [ ] Light mode on all three terminals
- [ ] Extended use test (30 minutes)
- [ ] Large PDF test (100+ pages)
- [ ] Search with 100+ matches
- [ ] Theme switching 20+ times (no issues)

---

## Success Criteria

### Functional
- [x] Dark theme is default and beautiful
- [x] Light theme available and readable
- [x] Theme switching works flawlessly
- [x] All UI elements visually consistent
- [x] Status bar informative and clear
- [x] Help screen comprehensive and readable

### Quality
- [x] Test coverage >80% (target: 85%)
- [x] All tests passing (66+ tests)
- [x] Zero linter warnings
- [x] Zero TODOs in production code
- [x] All exported symbols documented
- [x] Code review complete

### Performance
- [x] Startup <70ms ✅
- [x] Page switch cached <20ms ✅
- [x] Page switch uncached <65ms ✅
- [x] Memory baseline <10MB ✅
- [x] Memory with PDF <50MB ✅
- [x] No memory leaks ✅
- [x] Search <100ms ✅

### Documentation
- [x] README.md complete
- [x] ARCHITECTURE.md updated
- [x] PROGRESS.md finalized
- [x] PHASE_1_COMPLETE.md created
- [x] PHASE_2_PLAN.md drafted
- [x] All milestone reviews complete

---

## Phase 1 Completion Checklist

This milestone marks the end of Phase 1. Complete checklist:

### Functionality
- [x] PDF files load successfully
- [x] Text extraction works
- [x] 3-pane layout renders
- [x] All vim keybindings functional
- [x] Search finds matches
- [x] Page navigation works
- [x] Dark mode enabled
- [x] Light mode available
- [x] Help screen complete
- [x] Clean quit works

### Technical Excellence
- [x] 66+ tests passing
- [x] >80% code coverage
- [x] All linters clean
- [x] No memory leaks
- [x] Performance targets met
- [x] Clean architecture
- [x] Well documented

### Deliverables
- [x] Binary builds successfully
- [x] Installation script works
- [x] README.md complete
- [x] Documentation current
- [x] Demo materials ready
- [x] Phase 2 plan drafted

---

## Definition of Done

Milestone 1.6 is **DONE** when:

1. ✅ All acceptance criteria met
2. ✅ All tests passing (>80% coverage achieved)
3. ✅ All performance targets met
4. ✅ All documentation complete
5. ✅ Code quality perfect (0 warnings)
6. ✅ Demo materials prepared
7. ✅ Phase 1 complete review created
8. ✅ **Phase 1 is production-ready**

**Deliverable**: Production-ready LUMOS v0.1.0 - Dark Mode PDF Reader with vim keybindings.

---

## Post-Milestone Activities

### Release Preparation
- [ ] Tag version: `git tag -a v0.1.0 -m "Phase 1 Complete - MVP Release"`
- [ ] Create GitHub release notes
- [ ] Build binaries for all platforms
- [ ] Update project status everywhere

### Communication
- [ ] Demo to stakeholders
- [ ] Gather feedback
- [ ] Create user survey
- [ ] Plan Phase 2 kickoff

### Celebration
- [ ] Document wins and learnings
- [ ] Acknowledge contributors
- [ ] Take a break before Phase 2!

---

**Created**: 2025-11-10
**Target Completion**: 2025-11-19
**Assignee**: Claude Code + Quality Reviewer

**Phase 1 Timeline**: October 21 - November 19, 2025 (4 weeks)
