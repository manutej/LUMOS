# LUMOS Implementation Priorities

**Last Updated**: 2025-11-17
**Sprint Duration**: 9 days to MVP
**Critical Path**: TUI → Keybindings → Polish → Release

---

## 🔴 PRIORITY 0: CRITICAL BLOCKER - Complete PDF Reader Functionality

### THE MOST IMPORTANT THING
**LUMOS needs to be a working PDF reader FIRST**. Without the TUI, the excellent PDF engine (94.4% tested) is unusable.

---

## Priority Matrix

### 🔴 P0: URGENT - Must Complete Now (Days 1-3)
**Milestone 1.4: Basic TUI Framework**

These tasks block everything else:

1. **Initialize Bubble Tea Program** (2 hours)
   - Spec: `.specify/specs/milestone-1.4-tui-framework.md`
   - Modify `cmd/lumos/main.go`
   - Create tea.Program with alt screen
   - Wire up to existing PDF document

2. **Implement MVU Model** (3 hours)
   - Create Model struct with state
   - Implement Init() with page load command
   - Implement Update() with quit handling
   - Implement View() with basic rendering

3. **Create 3-Pane Layout** (4 hours)
   - 20% Metadata | 60% Viewer | 20% Search
   - Use lipgloss for borders and styling
   - Handle responsive width calculations
   - Join panes horizontally

4. **Add Viewport Component** (3 hours)
   - Integrate bubbles/viewport
   - Load PDF page content
   - Enable scrolling
   - Show page content

5. **Implement Status Bar** (2 hours)
   - Show: Page X/Y | filename.pdf | [?]Help [q]Quit
   - Position at bottom
   - Use distinct styling

6. **Handle Window Resize** (2 hours)
   - Listen for WindowSizeMsg
   - Recalculate pane widths
   - Update viewport dimensions
   - No visual artifacts

**Success Metric**: Can open PDF, see content in 3 panes, quit cleanly

---

### 🟠 P1: HIGH - Next Sprint (Days 4-6)
**Milestone 1.5: Vim Keybindings**

Cannot start until P0 complete:

1. **Basic Navigation** (3 hours)
   - j/k: Line down/up
   - d/u: Half page down/up
   - gg: First page
   - G: Last page

2. **Page Controls** (2 hours)
   - Ctrl+N: Next page
   - Ctrl+P: Previous page
   - Display page transition
   - Update status bar

3. **Search Mode** (4 hours)
   - /: Enter search mode
   - Type query in search pane
   - Enter: Execute search
   - n/N: Next/previous result
   - Highlight matches

4. **Help System** (2 hours)
   - ?: Show help overlay
   - List all keybindings
   - ESC: Close help
   - Markdown formatted

5. **Advanced Navigation** (3 hours)
   - Number+G: Go to page N
   - Tab: Cycle panes
   - Space/Shift+Space: Page down/up
   - Home/End: Document start/end

**Success Metric**: All vim users feel at home

---

### 🟡 P2: MEDIUM - Polish Sprint (Days 7-9)
**Milestone 1.6: Dark Mode Polish**

Quality and polish phase:

1. **Theme Refinement** (3 hours)
   - Perfect dark mode colors
   - Ensure >7:1 contrast ratios
   - Syntax highlighting for code blocks
   - Smooth theme transitions

2. **Performance Optimization** (4 hours)
   - Profile startup time (<100ms)
   - Optimize render pipeline
   - Minimize allocations
   - Cache styled strings

3. **Terminal Compatibility** (3 hours)
   - Test iTerm2, Terminal.app, Alacritty
   - Handle 256 color vs true color
   - Verify on Linux terminals
   - Document requirements

4. **Error Handling** (2 hours)
   - Graceful PDF load failures
   - Clear error messages
   - Recovery suggestions
   - No panics

5. **Documentation** (4 hours)
   - Update README with screenshots
   - Complete QUICKSTART guide
   - Document all keybindings
   - Create GIF demos

**Success Metric**: Professional, polished, production-ready

---

## 🔵 P3: LOW - Post-MVP Enhancements

These are explicitly NOT priorities for Phase 1:

- ❌ Image rendering (Phase 3)
- ❌ Table of contents (Phase 2)
- ❌ Bookmarks (Phase 2)
- ❌ Annotations (Phase 2)
- ❌ Multiple tabs (Phase 2)
- ❌ Configuration file (Phase 2)
- ❌ Export features (Phase 3)
- ❌ AI integration (Phase 4)

**Focus**: MVP functionality only

---

## Daily Priority Checklist

### Day 1 (Monday)
- [ ] Morning: Initialize Bubble Tea program
- [ ] Afternoon: Implement basic MVU model
- [ ] Evening: Test basic launch and quit

### Day 2 (Tuesday)
- [ ] Morning: Create 3-pane layout structure
- [ ] Afternoon: Integrate viewport component
- [ ] Evening: Load and display first page

### Day 3 (Wednesday)
- [ ] Morning: Implement status bar
- [ ] Afternoon: Handle window resize
- [ ] Evening: Complete Milestone 1.4 testing

### Day 4 (Thursday)
- [ ] Morning: Add j/k/d/u navigation
- [ ] Afternoon: Add gg/G page jumps
- [ ] Evening: Test navigation flow

### Day 5 (Friday)
- [ ] Morning: Implement Ctrl+N/P paging
- [ ] Afternoon: Add search mode UI
- [ ] Evening: Search integration

### Day 6 (Saturday)
- [ ] Morning: Help system overlay
- [ ] Afternoon: Advanced keybindings
- [ ] Evening: Complete Milestone 1.5 testing

### Day 7 (Sunday)
- [ ] Morning: Theme refinement
- [ ] Afternoon: Performance profiling
- [ ] Evening: Optimization implementation

### Day 8 (Monday)
- [ ] Morning: Terminal compatibility testing
- [ ] Afternoon: Error handling improvements
- [ ] Evening: Bug fixes

### Day 9 (Tuesday)
- [ ] Morning: Documentation updates
- [ ] Afternoon: Final testing suite
- [ ] Evening: Phase 1 complete! 🎉

---

## Decision Framework

When prioritizing tasks, ask:

1. **Does this make the PDF reader usable?** → P0
2. **Does this improve user experience?** → P1
3. **Does this add polish?** → P2
4. **Is this a nice-to-have?** → P3 (skip for now)

---

## Resource Allocation

### Time Budget (9 days)
- P0 Tasks: 3 days (33%)
- P1 Tasks: 3 days (33%)
- P2 Tasks: 3 days (33%)
- Buffer: Included in estimates

### Testing Budget
- Unit tests: Write alongside code
- Integration tests: End of each milestone
- Manual testing: 2 hours per milestone
- Performance testing: Day 7

### Documentation Budget
- Code comments: As you write
- User docs: Day 9
- Architecture updates: Per milestone
- README screenshots: Day 9

---

## Risk Mitigation Priority

### High Risk → High Priority
1. **TUI doesn't render correctly** → P0
2. **Performance targets missed** → P1
3. **Keybindings conflict** → P1
4. **Terminal incompatibility** → P2

### Low Risk → Low Priority
1. **Theme colors not perfect** → P2
2. **Help text formatting** → P2
3. **Advanced features** → P3

---

## Communication Priorities

### Daily Updates in PROGRESS.md
- What was completed
- What's blocking
- What's next
- Updated percentage

### Per-Milestone Reviews
- Create review document
- Update specifications if needed
- Commit and tag

### Phase Completion
- Full documentation
- Performance report
- Test coverage report
- Phase 2 planning

---

## The One Thing

If you can only do ONE THING today:

**Get the TUI to launch and display a PDF page in a viewport.**

Everything else builds on this foundation.

---

## Success Definition

Phase 1 is successful when a developer can:

1. ✅ Run `lumos paper.pdf`
2. ✅ See the PDF content immediately
3. ✅ Navigate with vim keys naturally
4. ✅ Search for text quickly
5. ✅ Work for hours without eye strain

That's it. Everything else is Phase 2+.

---

**Remember**:
- **Specifications drive priorities**
- **Working software over perfect code**
- **User experience over features**
- **Done is better than perfect**

**Current Sprint**: P0 - Make it work!
**Next Sprint**: P1 - Make it good!
**Final Sprint**: P2 - Make it beautiful!