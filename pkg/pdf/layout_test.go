package pdf

import (
	"strings"
	"testing"
)

// TestNewLayoutAnalyzer creates analyzer with defaults
func TestNewLayoutAnalyzer(t *testing.T) {
	analyzer := NewLayoutAnalyzer()

	if analyzer == nil {
		t.Error("Failed to create layout analyzer")
	}
	if analyzer.LineThreshold != 5.0 {
		t.Errorf("LineThreshold mismatch: got %f, want 5.0", analyzer.LineThreshold)
	}
	if analyzer.ColumnThreshold != 20.0 {
		t.Errorf("ColumnThreshold mismatch: got %f, want 20.0", analyzer.ColumnThreshold)
	}
	if !analyzer.UseRelativeThreshold {
		t.Error("UseRelativeThreshold should be true by default")
	}
}

// TestExtractWithLineBreaks preserves line structure
func TestExtractWithLineBreaks(t *testing.T) {
	analyzer := NewLayoutAnalyzer()
	analyzer.LineThreshold = 5.0
	analyzer.UseRelativeThreshold = false

	// Single line: elements all at same Y position
	elements := []TextElement{
		{Text: "Hello", X: 10, Y: 100, FontSize: 12, Font: "Helvetica", Width: 30},
		{Text: "World", X: 50, Y: 100, FontSize: 12, Font: "Helvetica", Width: 30},
	}

	result := analyzer.ExtractWithLineBreaks(elements)

	if !strings.Contains(result, "Hello") || !strings.Contains(result, "World") {
		t.Errorf("Result should contain both words: %s", result)
	}
	if strings.Count(result, "\n") > 0 {
		t.Error("Single line should not contain newlines")
	}
}

// TestExtractWithLineBreaks_MultiLine preserves multiple lines
func TestExtractWithLineBreaks_MultiLine(t *testing.T) {
	analyzer := NewLayoutAnalyzer()
	analyzer.LineThreshold = 5.0
	analyzer.UseRelativeThreshold = false

	// Multiple lines: different Y positions
	elements := []TextElement{
		// Line 1 (Y=100)
		{Text: "First", X: 10, Y: 100, FontSize: 12, Font: "Helvetica", Width: 30},
		{Text: "line", X: 50, Y: 100, FontSize: 12, Font: "Helvetica", Width: 30},
		// Line 2 (Y=80, more than 5 points below)
		{Text: "Second", X: 10, Y: 80, FontSize: 12, Font: "Helvetica", Width: 30},
		{Text: "line", X: 70, Y: 80, FontSize: 12, Font: "Helvetica", Width: 30},
	}

	result := analyzer.ExtractWithLineBreaks(elements)

	lines := strings.Split(result, "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d: %q", len(lines), lines)
	}

	if !strings.Contains(lines[0], "First") {
		t.Errorf("First line should contain 'First': %s", lines[0])
	}
	if !strings.Contains(lines[1], "Second") {
		t.Errorf("Second line should contain 'Second': %s", lines[1])
	}
}

// TestExtractWithLineBreaks_RelativeThreshold uses font size
func TestExtractWithLineBreaks_RelativeThreshold(t *testing.T) {
	analyzer := NewLayoutAnalyzer()
	analyzer.UseRelativeThreshold = true

	// Elements with different font sizes
	elements := []TextElement{
		// Large font (24pt), relative threshold = 24 * 0.5 = 12
		{Text: "Heading", X: 10, Y: 100, FontSize: 24, Font: "Helvetica-Bold", Width: 50},
		// Same Y within threshold
		{Text: "subheading", X: 80, Y: 100, FontSize: 24, Font: "Helvetica-Bold", Width: 60},
		// Small font (12pt), relative threshold = 12 * 0.5 = 6
		// Y=88 is 12 points below, more than 6, so new line
		{Text: "Body", X: 10, Y: 88, FontSize: 12, Font: "Helvetica", Width: 30},
	}

	result := analyzer.ExtractWithLineBreaks(elements)

	lines := strings.Split(result, "\n")
	if len(lines) < 2 {
		t.Errorf("Expected at least 2 lines with relative threshold, got %d", len(lines))
	}
}

// TestExtractWithLineBreaks_Empty handles empty input
func TestExtractWithLineBreaks_Empty(t *testing.T) {
	analyzer := NewLayoutAnalyzer()

	result := analyzer.ExtractWithLineBreaks(nil)

	if result != "" {
		t.Errorf("Empty input should produce empty output, got %q", result)
	}
}

// TestExtractWithLineBreaks_SingleElement single element
func TestExtractWithLineBreaks_SingleElement(t *testing.T) {
	analyzer := NewLayoutAnalyzer()

	elements := []TextElement{
		{Text: "Only", X: 10, Y: 100, FontSize: 12, Font: "Helvetica", Width: 30},
	}

	result := analyzer.ExtractWithLineBreaks(elements)

	if result != "Only" {
		t.Errorf("Single element should be returned as-is, got %q", result)
	}
}

// TestGroupIntoLines groups elements by Y coordinate
func TestGroupIntoLines(t *testing.T) {
	analyzer := NewLayoutAnalyzer()
	analyzer.LineThreshold = 5.0
	analyzer.UseRelativeThreshold = false

	elements := []TextElement{
		{Text: "Line1A", X: 10, Y: 100, FontSize: 12, Font: "Helvetica", Width: 30},
		{Text: "Line1B", X: 50, Y: 100, FontSize: 12, Font: "Helvetica", Width: 30},
		{Text: "Line2A", X: 10, Y: 85, FontSize: 12, Font: "Helvetica", Width: 30},
		{Text: "Line2B", X: 50, Y: 85, FontSize: 12, Font: "Helvetica", Width: 30},
	}

	lines := analyzer.groupIntoLines(elements)

	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	if len(lines[0]) != 2 {
		t.Errorf("First line should have 2 elements, got %d", len(lines[0]))
	}
	if len(lines[1]) != 2 {
		t.Errorf("Second line should have 2 elements, got %d", len(lines[1]))
	}
}

// TestDetectHeadings identifies headings by font size
func TestDetectHeadings(t *testing.T) {
	analyzer := NewLayoutAnalyzer()

	elements := []TextElement{
		{Text: "Heading", FontSize: 24, Font: "Helvetica-Bold"},     // Larger
		{Text: "Body", FontSize: 12, Font: "Helvetica"},             // Normal
		{Text: "Body", FontSize: 12, Font: "Helvetica"},             // Normal
		{Text: "SubHeading", FontSize: 16, Font: "Helvetica-Bold"}, // Larger
		{Text: "Body", FontSize: 12, Font: "Helvetica"},             // Normal
	}

	headings := analyzer.DetectHeadings(elements)

	// Should detect at least the 24pt heading
	if len(headings) == 0 {
		t.Error("Should detect at least one heading")
	}

	// First heading should be at index 0
	if len(headings) > 0 && headings[0] != 0 {
		t.Errorf("First heading should be at index 0, got %d", headings[0])
	}
}

// TestDetectHeadings_NoHeadings handles text without headings
func TestDetectHeadings_NoHeadings(t *testing.T) {
	analyzer := NewLayoutAnalyzer()

	elements := []TextElement{
		{Text: "Body", FontSize: 12, Font: "Helvetica"},
		{Text: "Body", FontSize: 12, Font: "Helvetica"},
		{Text: "Body", FontSize: 12, Font: "Helvetica"},
	}

	headings := analyzer.DetectHeadings(elements)

	if len(headings) != 0 {
		t.Errorf("Uniform font should detect no headings, got %d", len(headings))
	}
}

// TestTextElement_Structure verifies TextElement fields
func TestTextElement_Structure(t *testing.T) {
	elem := TextElement{
		Text:     "test",
		X:        10.5,
		Y:        20.5,
		FontSize: 12.0,
		Font:     "Helvetica",
		Width:    30.0,
	}

	if elem.Text != "test" {
		t.Errorf("Text mismatch: got %s", elem.Text)
	}
	if elem.X != 10.5 {
		t.Errorf("X mismatch: got %f", elem.X)
	}
	if elem.Y != 20.5 {
		t.Errorf("Y mismatch: got %f", elem.Y)
	}
	if elem.FontSize != 12.0 {
		t.Errorf("FontSize mismatch: got %f", elem.FontSize)
	}
	if elem.Font != "Helvetica" {
		t.Errorf("Font mismatch: got %s", elem.Font)
	}
	if elem.Width != 30.0 {
		t.Errorf("Width mismatch: got %f", elem.Width)
	}
}

// BenchmarkExtractWithLineBreaks benchmarks line breaking
func BenchmarkExtractWithLineBreaks(b *testing.B) {
	analyzer := NewLayoutAnalyzer()

	// Create a moderate-sized document
	elements := make([]TextElement, 100)
	for i := 0; i < 100; i++ {
		y := float64(1000 - (i/10)*20) // 10 lines
		elements[i] = TextElement{
			Text:     "Word",
			X:        float64((i % 10) * 50),
			Y:        y,
			FontSize: 12.0,
			Font:     "Helvetica",
			Width:    40.0,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = analyzer.ExtractWithLineBreaks(elements)
	}
}

// BenchmarkDetectHeadings benchmarks heading detection
func BenchmarkDetectHeadings(b *testing.B) {
	analyzer := NewLayoutAnalyzer()

	// Create a document with mixed font sizes
	elements := make([]TextElement, 1000)
	for i := 0; i < 1000; i++ {
		fontSize := 12.0
		if i%50 == 0 {
			fontSize = 20.0 // Some headings
		}
		elements[i] = TextElement{
			Text:     "Word",
			FontSize: fontSize,
			Font:     "Helvetica",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = analyzer.DetectHeadings(elements)
	}
}

// TestExtractWithColumns detects column structure
func TestExtractWithColumns(t *testing.T) {
	analyzer := NewLayoutAnalyzer()
	analyzer.ColumnThreshold = 20.0

	// Two-column layout
	// Left column: X=10-40
	// Right column: X=100-130 (gap of 60 points)
	elements := []TextElement{
		// Left column, line 1
		{Text: "Left1", X: 10, Y: 100, FontSize: 12, Font: "Helvetica", Width: 30},
		// Right column, line 1 (should detect as column break)
		{Text: "Right1", X: 100, Y: 100, FontSize: 12, Font: "Helvetica", Width: 30},
		// Left column, line 2
		{Text: "Left2", X: 10, Y: 80, FontSize: 12, Font: "Helvetica", Width: 30},
		// Right column, line 2
		{Text: "Right2", X: 100, Y: 80, FontSize: 12, Font: "Helvetica", Width: 30},
	}

	formatted, cols := analyzer.ExtractWithColumns(elements)

	if len(cols) == 0 {
		t.Error("Should detect at least one column")
	}

	if formatted == "" {
		t.Error("Formatted output should not be empty")
	}

	// Check that both column contents are present
	if !strings.Contains(formatted, "Left") && !strings.Contains(formatted, "Right") {
		t.Errorf("Output should contain column content: %s", formatted)
	}
}

// TestExtractWithColumns_SingleColumn single column PDF
func TestExtractWithColumns_SingleColumn(t *testing.T) {
	analyzer := NewLayoutAnalyzer()

	elements := []TextElement{
		{Text: "Text1", X: 10, Y: 100, FontSize: 12, Font: "Helvetica", Width: 30},
		{Text: "Text2", X: 10, Y: 80, FontSize: 12, Font: "Helvetica", Width: 30},
	}

	formatted, _ := analyzer.ExtractWithColumns(elements)

	// Single column PDF should have minimal column detection
	if formatted == "" {
		t.Error("Formatted output should not be empty")
	}

	if !strings.Contains(formatted, "Text1") || !strings.Contains(formatted, "Text2") {
		t.Errorf("Output should contain text: %s", formatted)
	}
}
