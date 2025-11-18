# Phase 3.2: Pdfcpu Integration & Image Extraction

**Status**: In Progress
**Timeline**: 2-3 weeks
**Goal**: Extract images from PDFs and prepare for terminal rendering

## Overview

Phase 3.2 focuses on integrating pdfcpu library for image extraction and connecting it to the existing image infrastructure from Phase 3.1.

## Implementation Strategy

### 1. Pdfcpu Integration

**Add Dependency**:
```bash
go get github.com/pdfcpu/pdfcpu@latest
```

**Key Imports**:
```go
import (
  "github.com/pdfcpu/pdfcpu/pkg/api"
  "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
  "image"
)
```

**Why pdfcpu over others**:
- ✅ Apache 2.0 license (permissive)
- ✅ Pure Go implementation (no C/C++ deps)
- ✅ Active maintenance
- ✅ Comprehensive PDF support
- ✅ Image extraction API available
- ✅ ~4.5KB per PDF file for extraction (acceptable overhead)

### 2. Dual-Library Architecture

**Text Extraction**: `ledongthuc/pdf` (proven, fast, unchanged)
**Image Extraction**: `pdfcpu` (new, comprehensive)

```
Document
├─ text via ledongthuc/pdf ← GetPage() ✅ Already working
├─ images via pdfcpu        ← GetPageImages() NEW
└─ metadata                  ← Both libraries support
```

**Benefits**:
- Zero impact on proven text extraction
- Lazy loading: images extracted on-demand only
- Easy to disable (optional feature)
- Clear separation of concerns

### 3. GetPageImages Implementation

**Location**: `pkg/pdf/document.go`

**Method Signature**:
```go
func (d *Document) GetPageImages(pageNum int, opts ImageExtractionOptions) ([]PageImage, error)
```

**Algorithm**:
1. Validate page number
2. Check image cache first
3. Open PDF with pdfcpu
4. Extract images from page
5. Filter by options (size, type, etc.)
6. Convert to PageImage structs
7. Cache and return

**Performance Targets**:
- First extraction: <100ms per page
- Cached retrieval: <1ms
- Memory: <2MB per 10 images

### 4. Image Filtering & Conversion

**ImageExtractionOptions**:
```go
type ImageExtractionOptions struct {
    MinWidth         float64 // Skip tiny images
    MinHeight        float64 // Skip tiny images
    OnlyInlineImages bool    // Skip background images
    PreserveDPI      bool    // Respect DPI info
    MaxImagesPerPage int     // Limit image count
}
```

**Filter Logic**:
1. Skip images smaller than MinWidth x MinHeight
2. Skip background images if OnlyInlineImages=true
3. Skip images beyond MaxImagesPerPage limit
4. Convert from pdfcpu Image → PageImage

**Supported Formats**:
- JPEG ✅ Native support in pdfcpu
- PNG ✅ Native support in pdfcpu
- TIFF ⚠️ May require conversion
- Others → Convert to PNG fallback

### 5. Model State Management

**Add to Model**:
```go
type Model struct {
    // ... existing fields ...

    // Image support (Phase 3.2)
    imageCache      *pdf.ImagePageCache
    showImages      bool              // Toggle with 'i' key
    imagesOnPage    []pdf.PageImage  // Current page images
    imageRenderCfg  ui.ImageRenderConfig
}
```

**Initialize in NewModel**:
```go
imageCache:     pdf.NewImagePageCache(10),  // 10-page cache
showImages:     true,                        // Default to enabled
imagesOnPage:   []pdf.PageImage{},
imageRenderCfg: ui.GetImageRenderConfig(width, height),
```

### 6. Message Handlers

**New Messages**:
```go
type ToggleImagesMsg struct{}
type LoadPageImagesMsg struct {
    PageNum int
}
type ImagesLoadedMsg struct {
    PageNum int
    Images  []pdf.PageImage
    Error   error
}
```

**Handle in Update**:
```go
case ToggleImagesMsg:
    m.showImages = !m.showImages
    return m, nil

case LoadPageImagesMsg:
    // Async load images for current page
    return m, tea.Batch(
        m.loadPageImages(msg.PageNum),
        // ... other commands ...
    )

case ImagesLoadedMsg:
    if msg.Error == nil {
        m.imagesOnPage = msg.Images
    }
    return m, nil
```

### 7. Keybinding Integration

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

### 8. Viewport Integration

**In View() method**:
```go
// Get page content as before
content := pageInfo.Text

// If images enabled, append image placeholders
if m.showImages && len(m.imagesOnPage) > 0 {
    images := m.renderImagesAsText(m.imagesOnPage)
    content = content + "\n\n" + images
}

m.viewport.SetContent(content)
```

**Image Rendering in View**:
```go
func (m *Model) renderImagesAsText(images []pdf.PageImage) string {
    var sb strings.Builder

    sb.WriteString("━━━ IMAGES ━━━\n\n")

    renderer := ui.NewImageRenderer(m.imageRenderCfg)

    for i, img := range images {
        // Use text fallback by default
        rendered := renderer.RenderImage(img.Data, img.Title)
        sb.WriteString(rendered)
        if i < len(images)-1 {
            sb.WriteString("\n\n")
        }
    }

    return sb.String()
}
```

## Testing Strategy

### Unit Tests (pkg/pdf/)

**Image Extraction Tests**:
```go
TestGetPageImages_EmptyPage        // Page with no images
TestGetPageImages_MultipleImages   // Page with 3+ images
TestGetPageImages_CachedRetrieval  // Cache hit/miss
TestGetPageImages_FilterBySize     // MinWidth/MinHeight
TestGetPageImages_FilterBackgroundImages
TestGetPageImages_LimitPerPage     // MaxImagesPerPage
```

### Integration Tests

**Model Tests**:
```go
TestModel_ToggleImages             // 'i' key toggles
TestModel_LoadImagesAsync          // Async loading
TestModel_ImagesInViewport         // Render in view
TestModel_ImageCachingPerformance  // <100ms first load
```

### Manual Testing

**Test PDFs**:
- Simple PDF with 1 image
- Multi-page with images on specific pages
- Image-heavy PDF (10+ images per page)
- Real-world academic paper
- Mixed text + images

**Terminals**:
- Current terminal (halfblock fallback)
- Kitty (if available for testing)
- iTerm2 (if on macOS)
- xterm (SIXEL support)

## Success Criteria

✅ Images extract from PDFs without errors
✅ Image cache working (>80% hit rate for accessed pages)
✅ Toggle works with 'i' key (on/off)
✅ Text fallback shows on all terminals
✅ Performance <100ms first extraction
✅ All 189 existing tests still pass
✅ 20+ new tests for image features
✅ No memory leaks (image cache size bounded)
✅ Backward compatible (text-only PDFs work unchanged)

## Risk Mitigation

| Risk | Mitigation |
|------|-----------|
| Pdfcpu library bugs | Fallback to text, log errors |
| Large images crash | Size limits, memory protection |
| Performance regression | Lazy loading, aggressive caching |
| Library incompatibility | Pin version, thorough testing |

## Next Steps (Phase 3.3+)

- Implement actual Kitty graphics protocol rendering
- Add iTerm2 inline image support
- SIXEL rendering for xterm
- Advanced image scaling and positioning
- Image caching optimization
- Performance profiling

## Timeline

**Week 1**:
- [ ] Add pdfcpu dependency
- [ ] Implement GetPageImages()
- [ ] Write extraction tests
- [ ] Basic image caching

**Week 2**:
- [ ] Integrate with Model
- [ ] Add keybindings
- [ ] Render in viewport
- [ ] Write integration tests

**Week 3** (if needed):
- [ ] Performance optimization
- [ ] Manual testing across terminals
- [ ] Documentation
- [ ] Polish and debugging

---

**Status**: Ready for implementation
**Estimated Effort**: 2-3 weeks for MVP
**Complexity**: Medium (new library integration, async loading)
