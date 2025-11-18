package pdf

import (
	"sort"
	"strings"
)

// TextElement represents a single text object with layout information
type TextElement struct {
	Text     string  // Actual UTF-8 text content
	X        float64 // Horizontal position in points (left edge)
	Y        float64 // Vertical position in points (baseline)
	FontSize float64 // Font size in points
	Font     string  // Font name (e.g., "Helvetica", "Times-Bold")
	Width    float64 // Width of text string in points
}

// LayoutAnalyzer provides layout-aware text extraction and formatting
type LayoutAnalyzer struct {
	// LineThreshold: Y-distance threshold for detecting new lines
	// If Y-distance between elements > threshold, treat as new line
	// Can be absolute (e.g., 5.0 points) or relative (e.g., FontSize * 0.5)
	LineThreshold float64

	// ColumnThreshold: X-distance threshold for detecting column breaks
	// If X-gap between elements > threshold, treat as column break
	ColumnThreshold float64

	// UseRelativeThreshold: if true, line threshold is relative to font size
	UseRelativeThreshold bool
}

// NewLayoutAnalyzer creates a layout analyzer with sensible defaults
func NewLayoutAnalyzer() *LayoutAnalyzer {
	return &LayoutAnalyzer{
		LineThreshold:        5.0, // 5 points for absolute threshold
		ColumnThreshold:      20.0,
		UseRelativeThreshold: true,
	}
}

// ExtractWithLineBreaks sorts text elements by position and preserves line structure
// Returns formatted text with proper line breaks instead of simple concatenation
func (la *LayoutAnalyzer) ExtractWithLineBreaks(elements []TextElement) string {
	if len(elements) == 0 {
		return ""
	}

	// Create working copy to avoid modifying original
	sorted := make([]TextElement, len(elements))
	copy(sorted, elements)

	// Sort by Y descending (top of page first, since PDF Y increases upward)
	// Then by X ascending (left to right)
	sort.Slice(sorted, func(i, j int) bool {
		// Y descending (larger Y first = higher on page)
		if sorted[i].Y != sorted[j].Y {
			return sorted[i].Y > sorted[j].Y
		}
		// X ascending (smaller X first = left to right)
		return sorted[i].X < sorted[j].X
	})

	var result strings.Builder
	var lastY float64 = -1.0
	lineThreshold := la.LineThreshold

	for i, elem := range sorted {
		// Determine if this is a new line based on Y-coordinate
		if i > 0 {
			yDistance := lastY - elem.Y
			// Use relative threshold if enabled
			if la.UseRelativeThreshold && elem.FontSize > 0 {
				lineThreshold = elem.FontSize * 0.5
			}

			// If Y distance exceeds threshold, start new line
			if yDistance > lineThreshold {
				result.WriteString("\n")
			} else if i > 0 && result.Len() > 0 {
				// Within same line, add space between elements
				result.WriteString(" ")
			}
		}

		result.WriteString(elem.Text)
		lastY = elem.Y
	}

	return result.String()
}

// ExtractWithColumns detects columns and reconstructs text in proper reading order
// Returns text formatted with column breaks (empty lines between columns)
func (la *LayoutAnalyzer) ExtractWithColumns(elements []TextElement) (string, []Column) {
	if len(elements) == 0 {
		return "", nil
	}

	// Group elements by Y-coordinate ranges (logical lines)
	lines := la.groupIntoLines(elements)

	// Detect column boundaries from X-coordinates across all lines
	columns := la.detectColumnBoundaries(lines)

	// Format text respecting column structure
	formatted := la.formatByColumns(lines, columns)

	return formatted, columns
}

// Column represents a vertical region of text
type Column struct {
	Left   float64 // Minimum X coordinate
	Right  float64 // Maximum X coordinate
	Width  float64 // Right - Left
	Lines  []string // Text content for each line in this column
}

// groupIntoLines groups text elements by Y-coordinate (logical lines)
func (la *LayoutAnalyzer) groupIntoLines(elements []TextElement) [][]TextElement {
	if len(elements) == 0 {
		return nil
	}

	// Sort by Y descending first
	sorted := make([]TextElement, len(elements))
	copy(sorted, elements)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Y != sorted[j].Y {
			return sorted[i].Y > sorted[j].Y
		}
		return sorted[i].X < sorted[j].X
	})

	var lines [][]TextElement
	var currentLine []TextElement
	var lastY float64 = -1.0
	lineThreshold := la.LineThreshold

	for _, elem := range sorted {
		if lastY < 0 {
			// First element
			currentLine = append(currentLine, elem)
			lastY = elem.Y
		} else {
			yDistance := lastY - elem.Y
			// Use relative threshold if enabled
			if la.UseRelativeThreshold && elem.FontSize > 0 {
				lineThreshold = elem.FontSize * 0.5
			}

			if yDistance > lineThreshold {
				// New line detected
				lines = append(lines, currentLine)
				currentLine = []TextElement{elem}
				lastY = elem.Y
			} else {
				// Same line
				currentLine = append(currentLine, elem)
			}
		}
	}

	// Add last line
	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

// detectColumnBoundaries analyzes X-coordinates to find column divisions
func (la *LayoutAnalyzer) detectColumnBoundaries(lines [][]TextElement) []Column {
	if len(lines) == 0 {
		return nil
	}

	// Collect all X values and gaps to detect columns
	var columns []Column
	xValues := make(map[float64]bool)

	// Get all unique X boundaries
	for _, line := range lines {
		for _, elem := range line {
			xValues[elem.X] = true
			xValues[elem.X + elem.Width] = true
		}
	}

	// Simple heuristic: if there's a significant gap, it's likely a column boundary
	var sortedX []float64
	for x := range xValues {
		sortedX = append(sortedX, x)
	}
	sort.Float64s(sortedX)

	// Find gaps larger than column threshold
	var columnStarts []float64
	columnStarts = append(columnStarts, sortedX[0])

	for i := 1; i < len(sortedX); i++ {
		gap := sortedX[i] - sortedX[i-1]
		if gap > la.ColumnThreshold {
			columnStarts = append(columnStarts, sortedX[i])
		}
	}

	// Create column objects
	for i, start := range columnStarts {
		col := Column{Left: start}
		if i+1 < len(columnStarts) {
			col.Right = columnStarts[i+1]
		} else {
			col.Right = sortedX[len(sortedX)-1]
		}
		col.Width = col.Right - col.Left
		columns = append(columns, col)
	}

	return columns
}

// formatByColumns constructs text respecting column structure
func (la *LayoutAnalyzer) formatByColumns(lines [][]TextElement, columns []Column) string {
	var result strings.Builder

	for lineIdx, line := range lines {
		// Process each column in this line
		for colIdx, col := range columns {
			var lineText strings.Builder
			first := true

			// Collect elements for this column in this line
			for _, elem := range line {
				if elem.X >= col.Left && elem.X+elem.Width <= col.Right+1 {
					if !first {
						lineText.WriteString(" ")
					}
					lineText.WriteString(elem.Text)
					first = false
				}
			}

			if lineText.Len() > 0 {
				result.WriteString(lineText.String())
				if colIdx < len(columns)-1 {
					result.WriteString(" | ") // Column separator
				}
			}
		}

		// Add line break between lines (but not after last line)
		if lineIdx < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// DetectHeadings identifies likely headings by font size
// Returns indices of elements that are likely headings
func (la *LayoutAnalyzer) DetectHeadings(elements []TextElement) []int {
	if len(elements) == 0 {
		return nil
	}

	// Find median and mean font sizes
	var sizes []float64
	for _, elem := range elements {
		if elem.FontSize > 0 {
			sizes = append(sizes, elem.FontSize)
		}
	}

	if len(sizes) == 0 {
		return nil
	}

	sort.Float64s(sizes)
	medianSize := sizes[len(sizes)/2]
	headingThreshold := medianSize * 1.2 // 20% larger than median

	var headings []int
	for i, elem := range elements {
		if elem.FontSize > headingThreshold {
			headings = append(headings, i)
		}
	}

	return headings
}
