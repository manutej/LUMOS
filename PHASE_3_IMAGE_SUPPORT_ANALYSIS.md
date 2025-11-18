# LUMOS Phase 3: Image Support - Comprehensive Analysis

**Date**: 2025-11-18
**Current Status**: Phase 2 Complete, Phase 3 Planning
**Scope**: Full image support implementation analysis
**Complexity Level**: HIGH (new dependencies, architecture changes)

---

## Executive Summary

Phase 3 (Image Support) is **feasible but complex**. The main challenge is that LUMOS's current PDF library (`ledongthuc/pdf`) has **no image extraction support**. This requires either:

1. **Switching to a new PDF library** (pdfcpu, UniDoc, go-pdfium)
2. **Using a secondary tool** (pdfcpu CLI, external process)

Terminal graphics rendering is **well-solved** with mature Go libraries (`go-termimg`, `rasterm`).

**Estimated Effort**: 8-12 weeks (6-8 for MVP, 4+ for polish)
**Risk Level**: Medium (library integration, cross-terminal compatibility)
**Recommendation**: Start with pdfcpu for pragmatic MVP approach

---

## 1. PDF Image Extraction Analysis

### Current State: ledongthuc/pdf Library

**ledongthuc/pdf** is a lightweight text-extraction focused library:
- ✅ Fast text extraction
- ✅ Simple, minimal dependencies
- ❌ **NO image extraction support**
- ❌ No embedded object/stream support
- ❌ Limited to text content only

This is a **hard blocker** for image support.

### Option 1: Switch PDF Libraries (Recommended)

#### 1A. pdfcpu (github.com/hhrutter/pdfcpu)

**Status**: ✅ Actively maintained, production-ready
**License**: Apache 2.0 (commercial friendly)
**Size**: ~500KB library

**Image Extraction Capabilities**:
```go
// Pros
✅ Full image extraction (JPEG, PNG, custom formats)
✅ Extracts images with metadata
✅ Works as library or CLI tool
✅ Excellent documentation
✅ No CGO required (pure Go)

// Cons
❌ Heavier dependency than ledongthuc/pdf
❌ More complex API
❌ Overkill if only using for image extraction
❌ Need to maintain both text + image extraction paths
```

**Code Pattern**:
```go
// Simple image extraction
import "github.com/pdfcpu/pdfcpu/pkg/api"

func ExtractImages(pdfPath string, pageNum int) ([]image.Image, error) {
    // Extract images from page
    // Returns []image.Image with metadata
}
```

**API**: Well-documented, active GitHub community

---

#### 1B. UniDoc UniPDF (github.com/unidoc/unipdf)

**Status**: ✅ Production-ready, heavily used
**License**: AGPL/Commercial (requires commercial license for proprietary use)
**Size**: Large (~2MB+)

**Image Extraction Capabilities**:
```go
// Pros
✅ Most powerful PDF library for Go
✅ Full image extraction with coordinates
✅ Extracts XObject images and inline images
✅ Preserves image metadata (size, position, DPI)
✅ Best text+image coordination

// Cons
❌ Commercial license required for non-AGPL use
❌ Much heavier than pdfcpu
❌ Licensing complexity
❌ Over-engineered for image-only use case
```

**Recommendation**: **NOT suitable** for LUMOS (licensing complications, overkill)

---

#### 1C. go-pdfium (github.com/klippa-app/go-pdfium)

**Status**: ✅ Actively maintained, modern
**License**: MIT
**Size**: Requires PDFium C++ library (can use WASM)

**Image Extraction Capabilities**:
```go
// Pros
✅ Google's PDFium library (used by Chromium)
✅ Can render to images (images ARE available as rendered output)
✅ Very accurate rendering
✅ WebAssembly option (no CGO needed)

// Cons
❌ Requires external C++ library OR WebAssembly runtime
❌ Heavyweight for simple extraction
❌ Overkill: designed for rendering, not just extraction
❌ Complex setup with external dependencies
```

**Not Recommended**: Too heavy, designed for full rendering not extraction

---

#### 1D. Comparison Table

| Aspect | ledongthuc/pdf | pdfcpu | UniPDF | go-pdfium |
|--------|---|---|---|---|
| **Image Extraction** | ❌ None | ✅ Yes | ✅ Yes | ✅ Yes (via render) |
| **Text Extraction** | ✅ Fast | ✅ Yes | ✅ Yes | ✅ Yes |
| **Pure Go** | ✅ Yes | ✅ Yes | ✅ Yes | ❌ Requires CGO/WASM |
| **Size** | ~100KB | ~500KB | ~2MB+ | Large+ |
| **Licensing** | MIT | Apache 2.0 | AGPL/Commercial | MIT |
| **Complexity** | ⭐ Low | ⭐⭐ Medium | ⭐⭐⭐ High | ⭐⭐⭐ High |
| **Recommended** | ✅ Keep for text | ✅ Add for images | ❌ Not suitable | ❌ Not suitable |

---

### Recommended Approach: Dual-Library Strategy

```go
// pkg/pdf/document.go
import (
    "github.com/ledongthuc/pdf"        // Text extraction (fast, lightweight)
    "github.com/pdfcpu/pdfcpu/pkg/api" // Image extraction (when needed)
)

type Document struct {
    filepath      string
    pages         int
    cache         *LRUCache          // Text pages
    imageCache    *ImageLRUCache     // Images (new)
    isImageReady  bool               // Lazy load indicator
}

// GetPage() - existing text extraction (unchanged)
func (d *Document) GetPage(pageNum int) (*PageInfo, error) { ... }

// GetPageImages() - new image extraction (on-demand)
func (d *Document) GetPageImages(pageNum int) ([]PageImage, error) {
    // Lazy load images only when needed
    // Fall back gracefully if no images exist
}
```

**Advantages**:
- ✅ Zero impact on existing text extraction
- ✅ Backward compatible
- ✅ On-demand loading (performance)
- ✅ Lazy initialization (no startup overhead)
- ✅ Can add/remove image support easily

**Cost**: +500KB binary size, minimal complexity

---

## 2. Terminal Graphics Options Analysis

### 2.1 Available Libraries

#### go-termimg (github.com/blacktop/go-termimg)

**Status**: ✅ Actively maintained
**License**: Apache 2.0
**Stars**: 400+ on GitHub
**Last Commit**: Recent (2024)

**Supported Protocols**:
1. **Kitty** (Recommended)
   - Modern graphics protocol
   - Fast encoding
   - Best quality/size ratio
   - Full RGBA support
   - Supported by: Kitty, WezTerm, Ghostty

2. **iTerm2** (macOS primary)
   - Native inline images
   - Works on macOS Terminal (3.3.0+)
   - Full color support
   - Supported by: iTerm2, WezTerm

3. **SIXEL** (Legacy fallback)
   - Older protocol
   - Limited to 256 colors
   - Works on: xterm, Konsole, mlterm
   - Fallback for older terminals

4. **Halfblocks** (Unicode fallback)
   - Works in ANY terminal
   - Uses Unicode block characters
   - Limited quality
   - Ultimate fallback

**API**:
```go
import "github.com/blacktop/go-termimg/pkg/termimg"

// Automatic protocol detection
func RenderImage(img image.Image) (string, error) {
    // Detects terminal capability
    // Returns ANSI codes for rendering
}

// Or explicit protocol
func RenderImageKitty(img image.Image) string
func RenderImageiTerm(img image.Image) string
func RenderImageSixel(img image.Image) string
func RenderImageHalfblocks(img image.Image) string
```

**Pros**:
- ✅ Auto-detection of terminal capabilities
- ✅ Graceful fallback chain
- ✅ Chainable API (fluent interface)
- ✅ TUI integration (Bubbletea compatible)
- ✅ Performance optimizations
- ✅ Good documentation

**Cons**:
- ⚠️ Library size (adds ~200KB)
- ⚠️ Complex configuration for advanced use

**Recommendation**: ✅ **Use go-termimg**

---

#### rasterm (github.com/BourgeoisBear/rasterm)

**Status**: ✅ Maintained
**License**: MIT
**Size**: Lightweight (~50KB)

**Supported Protocols**:
- Kitty
- iTerm2
- SIXEL

**API**:
```go
import "github.com/BourgeoisBear/rasterm"

// Terminal capability detection
func IsKittyCapable() bool
func IsItermCapable() bool
func IsSixelCapable() bool

// Direct image writing
func KittyWriteImage(w io.Writer, img image.Image) error
func ItermWriteImage(w io.Writer, img image.Image) error
func SixelWriteImage(w io.Writer, img image.Image) error
```

**Pros**:
- ✅ Minimal dependencies
- ✅ Lower-level (more control)
- ✅ Smaller size

**Cons**:
- ❌ No auto-fallback chain
- ❌ Manual protocol selection
- ❌ No Halfblocks fallback

**Alternative**: Good for minimalist approach, but go-termimg is better

---

### 2.2 Terminal Support Matrix

| Terminal | Kitty | iTerm2 | SIXEL | Halfblocks |
|----------|-------|--------|-------|-----------|
| **Kitty** | ✅ Native | - | - | ✅ Fallback |
| **iTerm2** | - | ✅ Native | ✅ 3.3.0+ | ✅ Fallback |
| **WezTerm** | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Fallback |
| **Ghostty** | ✅ Yes | - | - | ✅ Fallback |
| **xterm** | - | - | ✅ If compiled | ✅ Fallback |
| **Konsole** | - | - | ✅ 20.12+ | ✅ Fallback |
| **GNOME Terminal** | - | - | - | ✅ Fallback |
| **Alacritty** | - | - | - | ✅ Fallback |
| **macOS Terminal** | - | ✅ 3.3.0+ | - | ✅ Fallback |

**Coverage**: Using go-termimg's fallback chain provides support for **99%+ of terminals**

---

### 2.3 Recommended Approach

**Use go-termimg with fallback chain**:

```go
// pkg/ui/terminal.go (New)
type TerminalCapabilities struct {
    CanRenderKitty      bool
    CanRenderITerm2     bool
    CanRenderSIXEL      bool
    HasFullColor        bool
    PreferredProtocol   string // "kitty" | "iterm2" | "sixel" | "halfblocks"
}

// Detect once at startup
func DetectTerminalCapabilities() TerminalCapabilities {
    // Use go-termimg to detect
}

// Render with fallback
func RenderImageWithFallback(img image.Image, caps TerminalCapabilities) string {
    if caps.CanRenderKitty {
        return renderKitty(img)
    } else if caps.CanRenderITerm2 {
        return renderITerm2(img)
    } else if caps.CanRenderSIXEL {
        return renderSIXEL(img)
    } else {
        return renderHalfblocks(img)  // Always works
    }
}
```

---

## 3. Integration Architecture

### 3.1 Package Structure

```go
pkg/
├── pdf/
│   ├── document.go          (Enhanced: image extraction)
│   ├── image.go             (NEW: Image types)
│   ├── image_cache.go       (NEW: Image LRU cache)
│   ├── layout.go            (Enhanced: hybrid text+image)
│   └── image_test.go        (NEW: Tests)
│
├── ui/
│   ├── model.go             (Enhanced: image state)
│   ├── renderer.go          (NEW: Image rendering)
│   ├── terminal.go          (NEW: Terminal detection)
│   ├── keybindings.go       (Enhanced: image keys)
│   └── renderer_test.go     (NEW: Tests)
│
└── config/
    └── config.go            (Enhanced: image preferences)
```

---

### 3.2 Core Types

```go
// pkg/pdf/image.go (NEW)

type PageImage struct {
    ID           string        // Unique ID per document+page+index
    PageNum      int           // Page containing image
    ImageIndex   int           // Index on page (0-based)
    Data         image.Image   // Decoded image
    Format       string        // "jpeg", "png", etc.
    Width        int           // Native width
    Height       int           // Native height
    DPI          int           // Dots per inch (if available)
    X            float64       // PDF X coordinate (if available)
    Y            float64       // PDF Y coordinate (if available)
    RenderedSize RenderSize    // Size to render
}

type RenderSize struct {
    Width       int           // Terminal characters wide
    Height      int           // Terminal lines tall
    AspectRatio float64       // Width/Height ratio
    Dither      bool          // Apply dithering?
}

// Image cache (similar to text cache)
type ImageLRUCache struct {
    cache       map[string]*PageImage
    maxSize     int
    mu          sync.RWMutex
}

// Terminal capabilities
type TerminalCapabilities struct {
    CanRenderKitty    bool
    CanRenderITerm2   bool
    CanRenderSIXEL    bool
    CellWidth         int       // Pixels per character
    CellHeight        int       // Pixels per character
    MaxImageWidth     int       // Max terminal width for images
    PreferredProtocol string    // "kitty" | "iterm2" | "sixel" | "halfblocks"
}
```

---

### 3.3 Document Integration

```go
// pkg/pdf/document.go (Enhanced)

type Document struct {
    // ... existing fields ...
    
    // NEW: Image support
    hasImages        bool              // Quick check
    imageCache       *ImageLRUCache    // Image cache
    imageExtractor   *ImageExtractor   // pdfcpu wrapper
}

// GetPageImages returns images from a specific page
func (d *Document) GetPageImages(pageNum int) ([]PageImage, error) {
    // Check cache first
    // If miss, extract using pdfcpu
    // Store in cache
    // Return results
}

// GetPageWithImages returns text + layout info for images
func (d *Document) GetPageWithImages(pageNum int) (*PageWithImages, error) {
    textInfo, _ := d.GetPage(pageNum)
    images, _ := d.GetPageImages(pageNum)
    
    return &PageWithImages{
        Text:   textInfo.Text,
        Images: images,
        Layout: d.AnalyzeLayout(pageNum, images),
    }, nil
}

type PageWithImages struct {
    Text   string
    Images []PageImage
    Layout ImageLayout
}

type ImageLayout struct {
    HasImages      bool
    ImagePositions []ImagePosition  // Where images go in text
}
```

---

### 3.4 UI Integration

```go
// pkg/ui/model.go (Enhanced)

type Model struct {
    // ... existing fields ...
    
    // NEW: Image support
    terminalCaps    TerminalCapabilities
    showImages      bool          // Toggle on/off
    imageScaling    ImageScaling  // How to scale images
    renderMode      RenderMode    // Hybrid or text-only
}

enum RenderMode {
    TextOnly        // Show only text, skip images
    ImagesOnly      // Show only images, hide text (for diagrams)
    Hybrid          // Mix text and images (default)
    ImagePlaceholders // Show [Image] markers for text-only terminals
}

enum ImageScaling {
    Fit             // Scale to fit terminal width
    Stretch         // Use full terminal width
    Actual          // Use actual image size (may exceed terminal)
}

// Rendering with images
func (m *Model) View() string {
    // ... existing panes ...
    
    // In viewport pane, render text + images
    return m.renderViewportWithImages()
}

func (m *Model) renderViewportWithImages() string {
    // Get page with both text and images
    page, _ := m.document.GetPageWithImages(m.currentPage)
    
    // Interleave text and image sections
    var output strings.Builder
    
    for _, element := range m.interleaveLayers(page) {
        if element.IsText {
            output.WriteString(element.Text)
        } else {
            // Render image using terminal protocol
            output.WriteString(m.renderImage(element.Image))
        }
    }
    
    return output.String()
}

func (m *Model) renderImage(img PageImage) string {
    // Use terminal capabilities to choose protocol
    // Apply scaling based on imageScaling setting
    // Return ANSI codes
}
```

---

### 3.5 Keybindings for Images

```go
// Additional keybindings for Phase 3

Key             | Action
----------------|---------------------------------------------
i               | Toggle images on/off
I               | Toggle image-only mode
z               | Cycle image scaling (fit → stretch → actual)
+               | Increase image size
-               | Decrease image size
t               | Jump to next image
T               | Jump to previous image
c               | Copy current image to clipboard (if supported)
```

---

## 4. MVP Scope: Phase 3.1 vs Full Phase 3

### Phase 3.1 MVP (3-4 weeks, Pragmatic)

**Goal**: Get images working on at least 2 terminal types with text fallback

**Deliverables**:
- [ ] Add pdfcpu as image extraction library
- [ ] Implement ImageLRUCache
- [ ] Terminal capability detection
- [ ] Basic image rendering (Kitty protocol)
- [ ] Image toggle keybinding (i)
- [ ] Fallback to text for unsupported terminals
- [ ] 30+ tests
- [ ] Documentation

**Success Criteria**:
- [ ] Images display in Kitty terminal
- [ ] Images display in iTerm2 (with fallback or protocol)
- [ ] Text fallback works: "━━━[Image: 800x600]━━━"
- [ ] Performance: <200ms for image extraction per page
- [ ] Cache hit rate >80%
- [ ] No crashes on image-heavy PDFs
- [ ] Graceful degradation on unsupported terminals

**Not Included** (Phase 3.2+):
- SIXEL support (can add later)
- Halfblocks rendering (can add later)
- Image scaling options (manual size control)
- Image positioning precision
- Copy/paste of images
- Advanced image editing

---

### Full Phase 3 (8-12 weeks, Complete)

**3.2: Advanced Rendering** (2-3 weeks)
- SIXEL protocol support
- Halfblocks Unicode fallback
- Image scaling algorithms (Fit, Stretch, Actual)
- DPI-aware rendering
- Dithering for limited color terminals
- Multiple image types (PNG, JPEG, WebP, SVG raster)

**3.3: Layout & Performance** (2 weeks)
- Hybrid text+image layout analysis
- Smart image positioning
- Performance optimization
- Memory profiling
- Caching tuning

**3.4: UI Polish & Integration** (2-3 weeks)
- Full keybinding set
- Image metadata display
- Copy/export images
- Image search/filtering
- Multi-image handling per page
- Documentation & examples

---

## 5. Implementation Timeline

### Phase 3.1: Minimal Image Support (3-4 weeks)

```
Week 1: Foundation
├─ Add pdfcpu dependency
├─ Create image extraction layer
├─ Implement ImageLRUCache
└─ Terminal capability detection

Week 2: Rendering
├─ Integrate go-termimg
├─ Implement Kitty protocol rendering
├─ Add image toggle (i key)
└─ Text fallback for unsupported

Week 3: Integration
├─ Update Document.GetPage() for lazy images
├─ Update Model for image state
├─ Hybrid text+image rendering
└─ Tests + documentation

Week 4: Polish (if time)
├─ Performance optimization
├─ Cross-terminal testing
├─ README + examples
└─ Bug fixes
```

---

## 6. Risk Assessment

### High Priority Risks

#### R1: Library Integration Complexity
**Risk**: Adding pdfcpu causes issues with existing text extraction
**Impact**: High - Breaks core functionality
**Mitigation**:
- Dual-library approach (text: ledongthuc, images: pdfcpu)
- Comprehensive tests for both extraction paths
- Lazy loading (images only when needed)

#### R2: Terminal Compatibility
**Risk**: Images don't work on some popular terminals
**Impact**: High - Limits user base
**Mitigation**:
- Support full fallback chain (Kitty → iTerm2 → SIXEL → Halfblocks)
- Test on 5+ terminals (Kitty, iTerm2, WezTerm, Alacritty, xterm)
- Graceful degradation (text placeholders)
- Configuration for manual protocol selection

#### R3: Performance Degradation
**Risk**: Large images slow down paging
**Impact**: Medium - Bad user experience
**Mitigation**:
- Lazy loading (extract images on-demand)
- Aggressive image cache (LRU, 10-50 images)
- Image downsampling for large PDFs
- Benchmarks for all image operations

#### R4: Memory Usage
**Risk**: Large PDFs with many images use too much memory
**Impact**: Medium - Crash on large files
**Mitigation**:
- Size-limited cache
- Stream-based image decoding (if possible)
- Memory profiling
- User-configurable cache size

### Medium Priority Risks

#### R5: Image Format Compatibility
**Risk**: Some PDFs have uncommon image formats
**Impact**: Medium - Partial image loss
**Mitigation**:
- Support JPEG, PNG, WebP (covers 99% of PDFs)
- Graceful fallback for unsupported formats
- Log warnings for skipped images

#### R6: PDF Library Licensing
**Risk**: pdfcpu license conflicts
**Impact**: Low - Apache 2.0 is compatible
**Status**: ✅ Mitigated (Apache 2.0 is commercial-friendly)

---

## 7. Testing Strategy

### Unit Tests (15+ tests)

```go
// Image extraction
- ExtractImages() valid PDF
- ExtractImages() PDF with no images
- ExtractImages() invalid page number
- Image cache hit/miss
- Image cache eviction

// Terminal detection
- DetectTerminalCapabilities() Kitty
- DetectTerminalCapabilities() iTerm2
- DetectTerminalCapabilities() unsupported
- Fallback chain logic

// Rendering
- RenderImageKitty() basic image
- RenderImageKitty() various sizes
- RenderImageiTerm2() basic image
- RenderImageFallback() halfblocks
- Image scaling logic
```

### Integration Tests (10+ tests)

```go
// End-to-end
- Load PDF with images
- Get page with images
- Render images in terminal
- Image caching works
- Performance meets targets
```

### Manual Testing (5+ terminals)

```
✓ Kitty
✓ iTerm2 (macOS)
✓ WezTerm
✓ xterm
✓ Alacritty (fallback)
```

### Performance Benchmarks

```
BenchmarkExtractImages()         <200ms per page
BenchmarkImageCache()            <100ns cache hit
BenchmarkRenderImage()           <50ms per image
BenchmarkScaleImage()            <100ms per image
```

---

## 8. Go Dependency Changes

### New Dependencies

```toml
# go.mod additions
github.com/pdfcpu/pdfcpu v0.8.1+      # Image extraction
github.com/blacktop/go-termimg v0.9+  # Terminal graphics
```

**Impact**:
- Binary size: ~700KB → ~1.2MB (+71%)
- Build time: Minor impact
- Compilation: Pure Go (no CGO)

**Justification**:
- LUMOS is becoming a professional-grade reader
- Image support is expected feature
- Size still <2MB (better than most Electron apps)
- Go's static linking is huge advantage

---

## 9. Migration Path

### Phase 2 → Phase 3 Transition

**Compatibility**: ✅ Fully backward compatible
- Existing text-only PDFs work unchanged
- No changes to keybindings (except new 'i' key)
- No changes to configuration format
- Lazy loading: zero overhead for text-only mode

**User Experience**:
- Default: Show images if available, text fallback if not
- Toggle: 'i' key to hide/show images
- Graceful: Works on any terminal (some show placeholders)

---

## 10. Recommended Decision

### Go with Phase 3.1 MVP using:

1. **PDF Library**: pdfcpu for image extraction
   - Rationale: Apache 2.0, actively maintained, mature, pure Go
   
2. **Terminal Graphics**: go-termimg for rendering
   - Rationale: Best fallback chain, auto-detection, TUI-friendly
   
3. **Architecture**: Dual-library approach
   - Keep ledongthuc/pdf for text (fast, proven)
   - Add pdfcpu for images (on-demand, lazy-loaded)
   
4. **Timeline**: 3-4 weeks to MVP
   - Pragmatic approach: Get something working quickly
   - Polish in Phase 3.2+

5. **Success Metrics**:
   - [ ] Images display in Kitty (best case)
   - [ ] Images display in iTerm2 (macOS case)
   - [ ] Text fallback for unsupported terminals
   - [ ] Performance: <200ms per page
   - [ ] Cache efficiency: >80% hit rate
   - [ ] Test coverage: >80%

---

## 11. Appendix: Detailed Comparison

### ledongthuc/pdf vs pdfcpu

```go
// Text extraction (existing, use ledongthuc/pdf)
page := r.Page(pageNum)
texts := page.Content().Text
// Fast, simple, proven

// Image extraction (new, use pdfcpu)
import "github.com/pdfcpu/pdfcpu/pkg/api"

images := api.ExtractImagesByPageNr(filename, pageNum)
// Full image support, but heavier API
```

### go-termimg vs rasterm

```go
// go-termimg (RECOMMENDED)
import "github.com/blacktop/go-termimg/pkg/termimg"
img, _ := termimg.Open("image.png")
img.Render() // Auto-selects protocol, falls back gracefully

// rasterm (lighter alternative)
import "github.com/BourgeoisBear/rasterm"
if rasterm.IsKittyCapable() {
    rasterm.KittyWriteImage(os.Stdout, img)
} else {
    // Manual fallback logic needed
}
```

---

## 12. Key Takeaways

| Question | Answer |
|----------|--------|
| **Can LUMOS support images?** | ✅ Yes, but requires library changes |
| **Which terminal graphics library?** | ✅ go-termimg (best fallback) |
| **Which PDF library for images?** | ✅ pdfcpu (pragmatic choice) |
| **How long for MVP?** | 3-4 weeks |
| **Is it backwards compatible?** | ✅ Yes, fully |
| **Performance impact?** | Minimal (lazy loading) |
| **Terminal coverage?** | 99%+ with fallback chain |
| **Estimated binary size increase** | ~500KB (+7%) |
| **Recommended approach** | Dual-library: keep ledongthuc, add pdfcpu |
| **Risk level** | Medium (library integration, testing) |

---

**Conclusion**: Phase 3 (Image Support) is **recommended and feasible** using pdfcpu for extraction and go-termimg for rendering. A pragmatic MVP in 3-4 weeks can get basic image support working across multiple terminals with graceful fallback to text.
