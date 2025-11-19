# LUMOS Specification Framework

**Version**: 1.0.0
**Last Updated**: 2025-11-17
**Methodology**: GitHub spec-kit
**Philosophy**: Specifications drive implementation, not vice versa

---

## Executive Summary

LUMOS is a dark mode PDF reader for developers built with Go and Bubble Tea TUI framework. This specification framework provides the source of truth for all implementation decisions, architectural constraints, and feature requirements.

**Current Status**: Phase 1 MVP - 50% Complete (3/6 milestones)
**Critical Priority**: Complete TUI implementation for functional PDF reader

---

## Framework Structure

```
.specify/
├── SPECIFICATION_INDEX.md       # This file - master index
├── constitution.md              # ✅ Immutable architectural principles
├── HANDOFF.md                  # Quick start for new developers
├── PRIORITIES.md               # Ordered implementation priorities
├── ARCHITECTURE_DECISIONS.md   # Key decisions and constraints
├── TESTING_STRATEGY.md        # Testing approach and standards
├── specs/
│   ├── phase-1-mvp.md         # Phase 1 master specification
│   ├── milestone-1.4-tui-framework.md    # ✅ Current milestone
│   ├── milestone-1.5-vim-keybindings.md  # ✅ Next milestone
│   └── milestone-1.6-dark-mode-polish.md # ✅ Final MVP milestone
├── plans/
│   ├── tui-implementation.md   # Detailed TUI implementation plan
│   ├── keybinding-matrix.md    # Complete keybinding specification
│   └── performance-targets.md  # Performance budget and targets
└── contracts/
    ├── pdf-interface.md         # PDF package contract
    ├── ui-interface.md          # UI package contract
    └── config-interface.md      # Config package contract
```

---

## Specification Hierarchy

### 1. Constitutional Framework (Immutable)
- **Location**: `.specify/constitution.md`
- **Status**: ✅ Established
- **Purpose**: Defines non-negotiable quality standards and architectural principles
- **Enforcement**: All code must comply or document justified exceptions

### 2. Phase Specifications (Strategic)
- **Location**: `.specify/specs/phase-*.md`
- **Current**: Phase 1 MVP (50% complete)
- **Purpose**: High-level feature goals and success criteria
- **Next**: Phase 2 Enhanced Viewing (Q1 2026)

### 3. Milestone Specifications (Tactical)
- **Location**: `.specify/specs/milestone-*.md`
- **Current**: Milestone 1.4 - Basic TUI Framework
- **Purpose**: Detailed implementation requirements with acceptance criteria
- **Completed**: 1.1 Build, 1.2 Testing, 1.3 Fixtures (94.4% coverage)

### 4. Implementation Plans (Operational)
- **Location**: `.specify/plans/*.md`
- **Purpose**: Step-by-step implementation guides
- **Usage**: Developer reference during coding

### 5. Contract Specifications (Interfaces)
- **Location**: `.specify/contracts/*.md`
- **Purpose**: Define package boundaries and APIs
- **Enforcement**: Integration tests validate contracts

---

## Current Implementation Status

### Completed Components (50%)
```
✅ Core PDF Engine       - 94.4% test coverage
✅ LRU Cache System      - <100ns operations
✅ Search Implementation - <50μs performance
✅ Test Infrastructure   - 42 tests passing
✅ Build System         - 4.6MB binary
```

### In-Progress Components (Current Sprint)
```
🚧 TUI Framework (Milestone 1.4)
  └─ Bubble Tea integration
  └─ 3-pane layout system
  └─ Viewport rendering
  └─ Status bar implementation
```

### Pending Components (Next Sprints)
```
⏳ Vim Keybindings (Milestone 1.5)
⏳ Dark Mode Polish (Milestone 1.6)
⏳ Documentation Completion
⏳ Performance Optimization
```

---

## Critical Path to MVP

### Priority 1: Complete TUI (URGENT - 2-3 days)
**Specification**: `.specify/specs/milestone-1.4-tui-framework.md`
- Initialize Bubble Tea program
- Implement MVU pattern
- Create 3-pane layout
- Add viewport for PDF content
- Handle window resize

### Priority 2: Vim Keybindings (HIGH - 2-3 days)
**Specification**: `.specify/specs/milestone-1.5-vim-keybindings.md`
- Navigation: j/k/d/u/gg/G
- Page control: Ctrl+N/Ctrl+P
- Search mode: /
- Help system: ?
- Quit: q

### Priority 3: Dark Mode Polish (MEDIUM - 2-3 days)
**Specification**: `.specify/specs/milestone-1.6-dark-mode-polish.md`
- Theme refinement
- Color contrast optimization
- Terminal compatibility
- Performance tuning

---

## Key Architectural Decisions

### Decision 1: Bubble Tea for TUI
**Rationale**: MVU pattern provides predictable state management
**Trade-off**: Learning curve vs maintainability
**Alternative rejected**: Direct terminal manipulation (too complex)

### Decision 2: LRU Cache for Pages
**Rationale**: Balance memory usage with performance
**Configuration**: 5 pages cached by default
**Performance**: Cache hits <100ns

### Decision 3: Ledongthuc/pdf Library
**Rationale**: Pure Go, no CGO dependencies
**Limitation**: File reopening for each page
**Mitigation**: LRU cache minimizes impact

### Decision 4: 3-Pane Layout
**Rationale**: Familiar IDE-like interface
**Distribution**: 20% metadata | 60% viewer | 20% search
**Flexibility**: Responsive to terminal width

---

## Quality Gates

### Per-Milestone Gates
1. **Test Coverage**: Maintain >80% (currently 94.4%)
2. **Performance**: Meet targets in constitution
3. **Code Quality**: Pass fmt, vet, lint
4. **Documentation**: Update PROGRESS.md, create review doc

### Pre-Release Gates (Phase 1 Complete)
1. **Functionality**: All 6 milestones complete
2. **Stability**: 0 crashes on diverse PDFs
3. **Performance**: <100ms startup, <50ms navigation
4. **Documentation**: README, QUICKSTART, architecture complete

---

## Specification Evolution

### Update Triggers
- Milestone completion discoveries
- Performance bottleneck identification
- User feedback from testing
- Technical constraint changes

### Update Process
1. Document issue in milestone review
2. Update relevant specification
3. Validate constitutional compliance
4. Regenerate affected implementations
5. Update this index

---

## Navigation Guide

### For New Developers
1. Start with `.specify/HANDOFF.md`
2. Review `.specify/constitution.md`
3. Study `.specify/PRIORITIES.md`
4. Reference current milestone spec

### For Implementation
1. Read current milestone specification
2. Follow implementation plan in `.specify/plans/`
3. Validate against contracts in `.specify/contracts/`
4. Check constitutional compliance

### For Architecture Decisions
1. Review `.specify/ARCHITECTURE_DECISIONS.md`
2. Check constitutional principles
3. Document new decisions with ADRs

---

## Success Metrics

### Phase 1 MVP Success (Target: Nov 19, 2025)
- ✅ 6/6 milestones complete
- ✅ 80%+ test coverage maintained
- ✅ Performance targets met
- ✅ 0 crashes on test suite
- ✅ Documentation complete

### Current Progress Metrics
- **Milestones**: 3/6 (50%)
- **Test Coverage**: 94.4% ✅
- **Performance**: Exceeding targets ✅
- **Stability**: Production-ready core ✅
- **Documentation**: 70% complete

---

## Quick Links

### Specifications
- [Phase 1 MVP Specification](specs/phase-1-mvp.md)
- [Current Milestone 1.4](specs/milestone-1.4-tui-framework.md)
- [Next Milestone 1.5](specs/milestone-1.5-vim-keybindings.md)
- [Constitution](constitution.md)

### Implementation
- [TUI Implementation Plan](plans/tui-implementation.md)
- [Testing Strategy](TESTING_STRATEGY.md)
- [Performance Targets](plans/performance-targets.md)

### Documentation
- [Handoff Guide](HANDOFF.md)
- [Priorities](PRIORITIES.md)
- [Architecture Decisions](ARCHITECTURE_DECISIONS.md)

---

**Remember**: Specifications are the source of truth. Code serves specifications, not the other way around.