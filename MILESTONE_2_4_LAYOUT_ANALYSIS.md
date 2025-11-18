# LUMOS Milestone 2.4: Layout Preservation - Technical Analysis

## Executive Summary

LUMOS currently extracts PDF text as a **flat, unstructured string** by concatenating individual text objects with spaces. The `ledongthuc/pdf` library provides **rich positional data** (x, y coordinates, font size, width) that is completely discarded. This analysis identifies what data is available, current limitations, and the minimum viable approach for layout preservation.

---

## 1. Current Text Extraction Approach

### How Text is Currently Extracted

**File**: `/home/user/LUMOS/pkg/pdf/document.go` (lines 98-103)

```go
// Extract plain text from page
content := ""
texts := page.Content().Text
for _, text := range texts {
    content += text.S + " "  // Just concatenate strings with spaces
}
```

### What We're Throwing Away

**File**: ledongthuc/pdf library - `Text` struct

The PDF library provides **per-text-element metadata**:

```go
type Text struct {
    Font     string  // Font name (e.g., "Helvetica")
    FontSize float64 // Font size in points
    X        float64 // X coordinate (points, increases left→right)
    Y        float64 // Y coordinate (points, increases bottom→top)
    W        float64 // Width of the text string (points)
    S        string  // Actual UTF-8 text content
}
```

### The Core Problem

1. **Positional data lost**: X, Y, W coordinates tell us exactly where text appears on the page
2. **Reading order corrupted**: PDF text objects aren't stored in reading order; we need X,Y to reconstruct it
3. **Layout signals ignored**: Gaps, alignments, columns are all encoded in X,Y spacing
4. **Font changes ignored**: Font size/name changes are lost

---

## 2. Available Positional Data

### Coordinate System (from ledongthuc/pdf)

- **X-axis**: Points from left edge, increasing left→right (0 = left margin)
- **Y-axis**: Points from bottom edge, increasing bottom→top (0 = bottom margin)
- **Units**: All in points (1 point = 1/72 inch)
- **Page dimensions**: Can be extracted from PDF metadata

### What Each Field Tells Us

| Field | Meaning | Example Use |
|-------|---------|-------------|
| `X` | Horizontal position | Detect columns (X < midpoint = left col, X >= midpoint = right col) |
| `Y` | Vertical position | Sort text top-to-bottom; detect line breaks (Y changes > threshold) |
| `W` | Width of text string | Detect spacing/gaps (if W < expected, there's a gap) |
| `FontSize` | Point size | Detect headings (FontSize > body text) |
| `Font` | Font name | Detect bold/italic (font name contains "Bold" or "Italic") |

### Example: Multi-Column Detection

```
Column 1 (X: 40-280pt)    Column 2 (X: 300-540pt)
┌──────────────┐          ┌──────────────┐
│ Text at      │          │ Text at      │
│ X=50, Y=700  │          │ X=310, Y=700 │
└──────────────┘          └──────────────┘

        Gap (X: 280-300)
```

---

## 3. Current Rendering Limitations

### Where Content Flows

1. **Extraction** (`pkg/pdf/document.go`): PDF → flat string
2. **Storage** (`pkg/pdf/cache.go`): Cached as plain string
3. **Rendering** (`pkg/ui/model.go`): 
   - Line 374: `m.viewport.SetContent(msg.Content)` 
   - Bubble Tea viewport wraps text at terminal width

### Actual Limitations

| Limitation | Impact | Example |
|------------|--------|---------|
| No line structure | Rewrapping loses original intent | "foo\nbar" becomes "foo bar" (space lost) |
| No column detection | Multi-column PDFs read wrong | Left col → right col → left col instead of left col fully |
| No spacing preservation | Dense text becomes dense mush | Proper spacing helps readability |
| No heading detection | All text treated equally | Section headers not visually distinct |
| Text concatenation order wrong | Reading order corrupted | Due to PDF not storing in reading order |

### Viewport Behavior (Bubble Tea)

**File**: `/home/user/LUMOS/pkg/ui/model.go` (line 374)

The viewport receives a plain string and:
1. Wraps at terminal width
2. Splits by newlines (which we've lost)
3. Renders line by line

**Result**: Multi-column PDFs appear as continuous text blocks.

---

## 4. Minimum Viable Approach for MVP

### Strategy: Progressive Enhancement (Not Perfect Layout)

Rather than perfect layout preservation (too complex), implement **smart line-breaking and column detection**:

#### Option A: Basic Line Preservation (Minimal effort, good ROI)

```go
// pkg/pdf/layout.go (NEW)

type TextElement struct {
    Text     string
    X, Y     float64
    FontSize float64
    Font     string
}

type LayoutAnalyzer struct {
    lineThreshold float64  // Y-distance before new line (e.g., 5 points)
    colThreshold  float64  // X-distance gap before new column (e.g., 20 points)
}

// Collect all text elements with coordinates
func (d *Document) ExtractRawElements(pageNum int) ([]TextElement, error) {
    // Same extraction, but preserve coordinates
    f, r, err := pdf.Open(d.filepath)
    defer f.Close()
    
    page := r.Page(pageNum)
    texts := page.Content().Text
    
    elements := make([]TextElement, len(texts))
    for i, text := range texts {
        elements[i] = TextElement{
            Text:     text.S,
            X:        text.X,
            Y:        text.Y,
            FontSize: text.FontSize,
            Font:     text.Font,
        }
    }
    return elements, nil
}
```

#### Option B: Column Detection

```go
// Detect column layout from X-coordinate distribution
func (la *LayoutAnalyzer) DetectColumns(elements []TextElement) []Column {
    // Group elements by X-range (clusters of X values)
    // Common case: X < 300 or X >= 300
    
    columns := make([]Column, 0)
    var currentCol Column
    
    // Sort by Y (top to bottom), then X (left to right)
    // Group consecutive Y-ranges by their X position
    
    return columns
}

type Column struct {
    Left, Right float64  // X boundaries
    Lines []string      // Text content
}
```

#### Option C: Smart Line Breaking

```go
// Reconstruct lines based on Y coordinates
func (la *LayoutAnalyzer) ExtractWithLineBreaks(elements []TextElement) string {
    // Sort by Y descending (top to bottom)
    sort.Slice(elements, func(i, j int) bool {
        return elements[i].Y > elements[j].Y
    })
    
    // Within each Y-range, sort by X (left to right)
    var result strings.Builder
    prevY := elements[0].Y
    
    for _, elem := range elements {
        // New line if Y dropped more than threshold
        if prevY - elem.Y > la.lineThreshold {
            result.WriteString("\n")
        }
        result.WriteString(elem.Text)
        result.WriteString(" ")
        prevY = elem.Y
    }
    
    return result.String()
}
```

---

## 5. Implementation Roadmap (Pragmatic MVP)

### Phase 1: Foundation (Low Risk, High Value)
**Effort**: 1-2 days

1. Create `pkg/pdf/layout.go` with `TextElement` struct
2. Modify `GetPage()` to optionally return structured elements
3. Add `ExtractRawElements()` method to Document
4. Write 5 tests for element extraction

**Benefit**: Preserve all coordinate data; no breaking changes

### Phase 2: Smart Line Breaking
**Effort**: 2-3 days

1. Implement `ExtractWithLineBreaks()` 
2. Make it configurable (threshold value)
3. Update `GetPage()` to use line-broken text
4. Write 10 tests with various PDF layouts

**Benefit**: Readable multi-paragraph text; preserve spacing intent

### Phase 3: Column Detection (Optional)
**Effort**: 3-5 days

1. Implement `DetectColumns()` - identify column boundaries
2. Extract per-column text
3. Render as interleaved text or labeled columns
4. 15+ tests for complex layouts

**Benefit**: Proper reading order for multi-column PDFs

---

## 6. Technical Decisions to Make

### Decision 1: Where to Store Coordinates?

**Option A**: In `PageInfo` struct
```go
type PageInfo struct {
    PageNum    int
    Text       string            // Current
    LineCount  int
    WordCount  int
    Elements   []TextElement     // NEW: Raw coordinates
    HasImages  bool
    HasTables  bool
}
```
**Pros**: Backward compatible, optional
**Cons**: Increases memory

**Option B**: Separate cache/method
```go
func (d *Document) GetPageElements(pageNum int) ([]TextElement, error) {}
```
**Pros**: Clean separation
**Cons**: Requires callers to know about elements

**Recommendation**: Option A (best for Phase 2 integration)

### Decision 2: Line Breaking Threshold

PDF coordinate system means Y-values are **precise but fragmented**:
- Same logical line might have Y-values like: 700.1, 699.8, 700.2 (font rendering)
- Different lines clearly separated: 700 vs 650

**Recommended threshold**: 
- Y-distance > `FontSize * 0.5` = new line
- Or absolute: Y-distance > 5 points

**Why**: Accounts for sub-pixel rendering, multiple fonts on same line

### Decision 3: Column Detection Method

**Option A**: X-coordinate clustering
```go
// Find X-value gaps > threshold → separate columns
// Simple, deterministic
```

**Option B**: Page width heuristic
```go
// If text at X < width/2 and text at X > width/2 → 2 columns
// Assumes balanced layout
```

**Recommendation**: Option A (more reliable)

---

## 7. Integration Points

### Where This Affects Current Code

#### Search (`pkg/pdf/search.go`)

**Current**: Works with flat text strings
```go
func (d *Document) Search(query string) ([]SearchResult, error) {
    // Line 139: page, err := d.GetPage(pageNum)
    // Works on string - fine as-is
}
```

**Needs**: Minor updates if using structured elements
- Can still work on final formatted string
- **No breaking changes needed**

#### Rendering (`pkg/ui/model.go`)

**Current**: Takes string, viewport wraps it
```go
func (m *Model) renderViewerPane(width, height int) string {
    // Line 374: m.viewport.SetContent(msg.Content)
}
```

**Improvement**: Pre-format text before passing to viewport
```go
// NEW: Use line breaks from analyzer instead of letting viewport wrap
formattedContent := layoutAnalyzer.ExtractWithLineBreaks(elements)
m.viewport.SetContent(formattedContent)
```

**Benefit**: Better preservation of original layout

#### Caching (`pkg/pdf/cache.go`)

**Current**: Caches raw string (5 pages, LRU)
**Change**: Could cache either:
1. Just the formatted text (minor memory change)
2. Elements + formatted text (more memory, recompute if needed)

**Recommendation**: Cache formatted text (less memory, same performance)

---

## 8. Testing Strategy

### Test Fixtures Needed

Use existing PDF fixtures but verify they test layout:

1. **simple.pdf**: 1 page, basic text (already exists)
2. **multipage.pdf**: 5 pages, paragraphs (already exists)
3. **search_test.pdf**: Multiple pages, various text (already exists)
4. **columns.pdf** (NEW): 2-column layout test
5. **headings.pdf** (NEW): Multiple font sizes

### Test Cases

```go
// Column detection
func TestDetectColumns_TwoColumn(t *testing.T) {
    // Left column X: 40-280, Right column X: 300-540
    // Should detect 2 columns
}

// Line preservation
func TestExtractWithLineBreaks(t *testing.T) {
    // Y values: 700, 695, 650, 645
    // Should produce 2 lines (700+695 grouped, 650+645 grouped)
}

// Font detection
func TestDetectHeadings(t *testing.T) {
    // Heading FontSize=24, Body FontSize=12
    // Should identify heading lines
}

// Integration
func TestLayoutPreservationWithSearch(t *testing.T) {
    // Extract with layout, search still works
}
```

---

## 9. Open Questions & Risks

### Question 1: Page Dimensions
How do we get page width/height?
```go
// From PDF library - check if available
page := r.Page(pageNum)
// Does page.V have width/height in its dict?
```
**Risk**: May need to parse PDF structure (Page MediaBox)

### Question 2: Rotated/Skewed PDFs
What if PDF page is rotated 90°?
**Risk**: Y-axis interpretation changes
**Mitigation**: Check page rotation metadata first

### Question 3: Performance Impact
Will storing elements in cache impact memory?
**Current**: ~8MB for 10MB PDF, 5-page cache
**Proposed**: +small amount (just coordinates, not huge)
**Mitigation**: Profile before/after

### Question 4: Complex Layouts
What about:
- Text in boxes/tables
- Vertical text (e.g., Japanese)
- Overlapping text

**Mitigation**: Start with simple heuristics, iterate

---

## 10. Success Criteria for MVP

### Minimum Success
- [x] Extract coordinates without breaking existing code
- [ ] Implement basic line breaking (Y-coordinate based)
- [ ] 90% of test PDFs preserve readability better than before
- [ ] No performance regression (same cache hits/misses)

### Nice-to-Have (Phase 2+)
- [ ] Column detection for 2-column layouts
- [ ] Heading detection (font size heuristics)
- [ ] Table detection (aligned columns)
- [ ] Spacing preservation (actual gap rendering)

---

## 11. Code Changes Checklist

### New Files
- [ ] `pkg/pdf/layout.go` - Layout analysis engine

### Modified Files
- [ ] `pkg/pdf/document.go` - Add `TextElement` to `PageInfo`, add `ExtractRawElements()`
- [ ] `pkg/pdf/document_test.go` - Add layout tests (10+ tests)

### Optional
- [ ] `pkg/ui/model.go` - Use formatted text if significant improvement

### No Changes Needed
- `pkg/ui/viewport.go` - Stays as-is (uses strings)
- `pkg/pdf/search.go` - Works with final text
- `pkg/pdf/cache.go` - Caches final text

---

## Summary Table

| Aspect | Current State | MVP Addition | Effort |
|--------|---------------|--------------|--------|
| **Data Available** | Text only | Text + Coordinates | 0 (just parse) |
| **Line Breaking** | Broken | Y-coord based | 1-2 days |
| **Column Detection** | None | X-coord clustering | 3-5 days (optional) |
| **Search Impact** | None | None (still works) | 0 |
| **Memory Impact** | None | +1-2% | Negligible |
| **Test Coverage** | 94.4% | +10-15 tests | 2-3 days |
| **API Changes** | None | 1 new method | Breaking? No |

---

## Recommended First Step

**Start with Phase 1**: Extract and store raw elements in `PageInfo`

```go
type PageInfo struct {
    PageNum   int
    Text      string
    Elements  []TextElement  // NEW
    // ...rest unchanged
}

// In GetPage(), after line 103:
d.cache.Put(pageNum, content)  // Current
// Add later: pageInfo.Elements = textsFromContent
```

This costs minimal effort but opens the door to all subsequent optimizations.

