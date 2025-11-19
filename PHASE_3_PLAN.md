# Phase 3: Image Support - Complete Plan

**Status**: 30% COMPLETE (Phase 3.1-3.2 done, Phase 3.3 in progress)
**Timeline**: 6-8 weeks total for full implementation
**Goal**: Render PDF images in terminal with intelligent fallbacks

---

## Phase 3 Breakdown

### Phase 3.1: Foundation ✅ COMPLETE
- Image data structures (PageImage, RenderSize, etc.)
- Terminal capability detection (Kitty, iTerm2, SIXEL, Halfblock)
- Image renderer with placeholder stubs
- Image cache infrastructure
- 4 new modules, 495 lines

### Phase 3.2: Extraction Infrastructure ✅ COMPLETE
- Image extraction preparation
- LRU cache system for images
- Document integration (imageCache)
- 14 comprehensive tests
- Pdfcpu integration pattern ready

### Phase 3.3: UI Integration 🚧 IN PROGRESS
**Timeline**: 1 week
**Goal**: Connect image infrastructure to Model and keybindings

#### What we're doing:
1. Add image state to Model
2. Create message types for image operations
3. Add 'i' keybinding to toggle images
4. Implement message handlers
5. Write tests for state management

#### What we're NOT doing yet:
- Actual pdfcpu integration (Phase 3.4)
- Image rendering (Phase 3.5)
- Advanced features (Phase 3.6)

### Phase 3.4: Pdfcpu Integration ⏳ PLANNED
**Timeline**: 2-3 weeks
**Goal**: Extract actual images from PDFs

#### Tasks:
1. Add pdfcpu dependency
2. Implement GetPageImages() with real extraction
3. Handle image format conversion
4. Error handling and fallbacks
5. Comprehensive extraction tests

### Phase 3.5: Terminal Rendering ⏳ PLANNED
**Timeline**: 2-3 weeks
**Goal**: Render images in terminal

#### Tasks:
1. Implement Kitty graphics protocol
2. Implement iTerm2 inline images
3. Implement SIXEL rendering
4. Implement halfblock fallback
5. Test on multiple terminals

### Phase 3.6: Polish & Optimization ⏳ PLANNED
**Timeline**: 1-2 weeks
**Goal**: Performance, edge cases, documentation

#### Tasks:
1. Performance profiling
2. Memory optimization
3. Edge case handling
4. Comprehensive documentation
5. User testing and feedback

---

## Phase 3.3: UI Integration (Current Focus)

### Model Changes

**Add to Model struct**:
```go
imageCache      *pdf.ImagePageCache
showImages      bool                  // Toggle with 'i' key
imagesOnPage    []pdf.PageImage      // Current page images
imageRenderCfg  ui.ImageRenderConfig // Terminal config
imageLoading    bool                  // Loading state
```

**Message Types**:
```go
type ToggleImagesMsg struct{}
type LoadPageImagesMsg struct{ PageNum int }
type ImagesLoadedMsg struct {
    PageNum int
    Images  []pdf.PageImage
    Error   error
}
```

**Message Handlers in Update()**:
- ToggleImagesMsg: Toggle showImages flag
- LoadPageImagesMsg: Async load images for page
- ImagesLoadedMsg: Store loaded images in model

### Keybinding Changes

**Add to keybindings.go**:
```go
case "i":
    return ToggleImages

var ToggleImages = func() tea.Msg {
    return ToggleImagesMsg{}
}
```

**Update VimKeybindingReference**:
```go
"i": "Toggle images on/off",
```

### View Changes (Future - Phase 3.5)

Once we have image rendering:
```go
if m.showImages && len(m.imagesOnPage) > 0 {
    images := m.renderImages(m.imagesOnPage)
    content = content + "\n\n" + images
}
```

For now: Just add the state, tests confirm it works.

---

## Testing Strategy for Phase 3.3

### Model Tests
```go
TestModel_ToggleImages          // 'i' key toggles flag
TestModel_ImagesState           // Initial state correct
TestModel_LoadImagesMsg         // Message handler works
TestModel_ClearsImagesOnPageChange // State management
```

### Integration Tests
```go
TestKeyBinding_ImageToggle      // 'i' key recognized
TestKeyBinding_ReferenceUpdated // Help shows 'i'
```

---

## Pragmatic Design Decisions

1. **Decouple Extraction from Rendering**: Don't wait for pdfcpu, build UI layer now
2. **Use Stubs**: GetPageImages() returns empty for now, tests prove framework works
3. **Lazy Load Images**: Only load if showImages=true
4. **Cache Everything**: Use existing LRU cache pattern
5. **Fallback Chain**: Text placeholder works on any terminal
6. **Keep It Simple**: One keybinding, one toggle, one state

---

## Success Criteria for Phase 3.3

✅ Model has image state
✅ 'i' keybinding toggles images
✅ Message handlers process image commands
✅ Tests prove state management works
✅ All 203 existing tests still pass
✅ 5-10 new tests for image state
✅ Zero breaking changes
✅ Ready for Phase 3.4 (pdfcpu integration)

---

## Architecture Diagram

```
User Input ('i' key)
    ↓
KeyHandler
    ↓
ToggleImagesMsg
    ↓
Model.Update()
    ↓
showImages = !showImages
    ↓
Model.View() (no change yet, but ready for 3.5)
    ↓
Viewport displays (text placeholder for now)
```

Once Phase 3.4 (pdfcpu):
```
User navigates to page
    ↓
LoadPageImagesMsg
    ↓
Document.GetPageImages() [NOW USES PDFCPU]
    ↓
ImagesLoadedMsg
    ↓
Model.imagesOnPage = images
    ↓
Model.View() calls ImageRenderer
    ↓
Viewport displays [REAL IMAGES OR TEXT FALLBACK]
```

---

## Implementation Checklist

### Phase 3.3 (This Week)
- [ ] Add image state to Model struct
- [ ] Create ToggleImagesMsg message type
- [ ] Add message handler in Update()
- [ ] Add 'i' keybinding
- [ ] Write 5-8 tests
- [ ] Commit and push
- [ ] All tests passing (208+ total)

### Phase 3.4 (Next 2-3 Weeks)
- [ ] Add pdfcpu library
- [ ] Implement actual GetPageImages()
- [ ] Handle image format conversion
- [ ] Error handling
- [ ] Write extraction tests
- [ ] Test with real PDFs

### Phase 3.5 (2-3 Weeks After)
- [ ] Implement Kitty protocol
- [ ] Implement iTerm2 support
- [ ] Implement SIXEL support
- [ ] Test on multiple terminals
- [ ] Optimize performance

### Phase 3.6 (Final Polish)
- [ ] Performance profiling
- [ ] Memory optimization
- [ ] Edge case handling
- [ ] Comprehensive docs
- [ ] Phase 3 completion review

---

## Risk Mitigation

| Risk | Mitigation |
|------|-----------|
| Pdfcpu API complexity | Extensive documentation, gradual integration |
| Terminal incompatibility | Complete fallback chain (5 levels) |
| Memory leaks | Bounded LRU cache, tests for leaks |
| Performance regression | Lazy loading, aggressive caching, benchmarks |
| Breaking changes | Stub-based approach, backward compatible |

---

## Next: Phase 3.3 Implementation

Ready to:
1. Add image state to Model
2. Create message handlers
3. Add keybindings
4. Write tests
5. Commit (Phase 3.3 - UI Integration)

This keeps things lean and pragmatic: complete the UI framework now, plug in the extraction layer later when pdfcpu is ready.
