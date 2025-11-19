package pdf

import (
	"strings"
	"testing"
)

// TestTOCEntry_Basic verifies TOC entry structure
func TestTOCEntry_Basic(t *testing.T) {
	entry := TOCEntry{
		Title: "Chapter 1",
		Page:  5,
		Level: 1,
	}

	if entry.Title != "Chapter 1" {
		t.Errorf("Title mismatch: got %q, want %q", entry.Title, "Chapter 1")
	}
	if entry.Page != 5 {
		t.Errorf("Page mismatch: got %d, want %d", entry.Page, 5)
	}
	if entry.Level != 1 {
		t.Errorf("Level mismatch: got %d, want %d", entry.Level, 1)
	}
}

// TestTableOfContents_Empty verifies empty TOC
func TestTableOfContents_Empty(t *testing.T) {
	toc := &TableOfContents{
		Entries: []TOCEntry{},
		Source:  "none",
	}

	if len(toc.Entries) != 0 {
		t.Errorf("Expected empty TOC, got %d entries", len(toc.Entries))
	}
	if toc.Source != "none" {
		t.Errorf("Source mismatch: got %q, want %q", toc.Source, "none")
	}
}

// TestTableOfContents_SingleEntry verifies single entry TOC
func TestTableOfContents_SingleEntry(t *testing.T) {
	toc := &TableOfContents{
		Entries: []TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
		},
		Source: "metadata",
	}

	if len(toc.Entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(toc.Entries))
	}
	if toc.Entries[0].Title != "Chapter 1" {
		t.Errorf("Entry title mismatch: got %q", toc.Entries[0].Title)
	}
}

// TestTableOfContents_MultipleEntries verifies TOC with multiple entries
func TestTableOfContents_MultipleEntries(t *testing.T) {
	toc := &TableOfContents{
		Entries: []TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
			{Title: "Chapter 2", Page: 10, Level: 1},
			{Title: "Chapter 3", Page: 20, Level: 1},
		},
		Source: "metadata",
	}

	if len(toc.Entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(toc.Entries))
	}

	expectedPages := []int{1, 10, 20}
	for i, expectedPage := range expectedPages {
		if toc.Entries[i].Page != expectedPage {
			t.Errorf("Entry %d page mismatch: got %d, want %d", i, toc.Entries[i].Page, expectedPage)
		}
	}
}

// TestDetectHeadingLevel_VeryShort detects level 1 for very short text
func TestDetectHeadingLevel_VeryShort(t *testing.T) {
	level := detectHeadingLevel("Chapter")
	if level != 1 {
		t.Errorf("Expected level 1 for very short text, got %d", level)
	}

	level = detectHeadingLevel("Part Two")
	if level != 1 {
		t.Errorf("Expected level 1 for 2-word text, got %d", level)
	}
}

// TestDetectHeadingLevel_Medium detects level 2 for medium text
func TestDetectHeadingLevel_Medium(t *testing.T) {
	level := detectHeadingLevel("Chapter One Introduction")
	if level != 2 {
		t.Errorf("Expected level 2 for 3-word text, got %d", level)
	}

	level = detectHeadingLevel("Getting Started with Python")
	if level != 2 {
		t.Errorf("Expected level 2 for 4-word text, got %d", level)
	}
}

// TestDetectHeadingLevel_Long detects level 3 for long text
func TestDetectHeadingLevel_Long(t *testing.T) {
	level := detectHeadingLevel("Advanced Topics in Computer Science and Engineering")
	if level != 3 {
		t.Errorf("Expected level 3 for long text, got %d", level)
	}
}

// TestBuildHierarchy_Flat converts flat list to hierarchy
func TestBuildHierarchy_Flat(t *testing.T) {
	toc := &TableOfContents{
		Entries: []TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1, Children: []TOCEntry{}},
			{Title: "Section 1.1", Page: 2, Level: 2, Children: []TOCEntry{}},
			{Title: "Section 1.2", Page: 5, Level: 2, Children: []TOCEntry{}},
			{Title: "Chapter 2", Page: 10, Level: 1, Children: []TOCEntry{}},
		},
		Source: "headings",
	}

	hierarchical := toc.BuildHierarchy()

	// Should have 2 root entries (2 chapters)
	if len(hierarchical.Entries) != 2 {
		t.Errorf("Expected 2 root entries, got %d", len(hierarchical.Entries))
	}

	// First chapter should have 2 children
	if len(hierarchical.Entries[0].Children) != 2 {
		t.Errorf("Expected 2 children for Chapter 1, got %d", len(hierarchical.Entries[0].Children))
	}

	// Second chapter should have 0 children
	if len(hierarchical.Entries[1].Children) != 0 {
		t.Errorf("Expected 0 children for Chapter 2, got %d", len(hierarchical.Entries[1].Children))
	}
}

// TestFindEntryByPage_Found returns entries for specific page
func TestFindEntryByPage_Found(t *testing.T) {
	toc := &TableOfContents{
		Entries: []TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
			{Title: "Chapter 2", Page: 10, Level: 1},
			{Title: "Chapter 3", Page: 20, Level: 1},
		},
		Source: "metadata",
	}

	results := toc.FindEntryByPage(10)
	if len(results) != 1 {
		t.Errorf("Expected 1 result for page 10, got %d", len(results))
	}
	if results[0].Title != "Chapter 2" {
		t.Errorf("Expected Chapter 2, got %s", results[0].Title)
	}
}

// TestFindEntryByPage_NotFound returns empty for non-existent page
func TestFindEntryByPage_NotFound(t *testing.T) {
	toc := &TableOfContents{
		Entries: []TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
			{Title: "Chapter 2", Page: 10, Level: 1},
		},
		Source: "metadata",
	}

	results := toc.FindEntryByPage(999)
	if len(results) != 0 {
		t.Errorf("Expected 0 results for non-existent page, got %d", len(results))
	}
}

// TestFindEntryByPage_InHierarchy finds entries in hierarchical TOC
func TestFindEntryByPage_InHierarchy(t *testing.T) {
	toc := &TableOfContents{
		Entries: []TOCEntry{
			{
				Title: "Chapter 1",
				Page:  1,
				Level: 1,
				Children: []TOCEntry{
					{Title: "Section 1.1", Page: 2, Level: 2},
					{Title: "Section 1.2", Page: 5, Level: 2},
				},
			},
		},
		Source: "headings",
	}

	// Find subsection
	results := toc.FindEntryByPage(5)
	if len(results) != 1 {
		t.Errorf("Expected 1 result for page 5, got %d", len(results))
	}
	if results[0].Title != "Section 1.2" {
		t.Errorf("Expected Section 1.2, got %s", results[0].Title)
	}
}

// TestGetPageForEntry returns correct page
func TestGetPageForEntry(t *testing.T) {
	entry := &TOCEntry{
		Title: "Chapter 5",
		Page:  42,
		Level: 1,
	}

	page := entry.GetPageForEntry()
	if page != 42 {
		t.Errorf("Expected page 42, got %d", page)
	}
}

// TestFormatTOC_Single formats single entry TOC
func TestFormatTOC_Single(t *testing.T) {
	toc := &TableOfContents{
		Entries: []TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
		},
		Source: "metadata",
	}

	formatted := toc.FormatTOC()
	if !strings.Contains(formatted, "Chapter 1") {
		t.Errorf("Formatted TOC missing chapter title")
	}
	if !strings.Contains(formatted, "page 1") {
		t.Errorf("Formatted TOC missing page number")
	}
	if !strings.Contains(formatted, "metadata") {
		t.Errorf("Formatted TOC missing source")
	}
}

// TestFormatTOC_Hierarchical formats hierarchical TOC with indentation
func TestFormatTOC_Hierarchical(t *testing.T) {
	toc := &TableOfContents{
		Entries: []TOCEntry{
			{
				Title: "Chapter 1",
				Page:  1,
				Level: 1,
				Children: []TOCEntry{
					{Title: "Section 1.1", Page: 2, Level: 2},
					{Title: "Section 1.2", Page: 5, Level: 2},
				},
			},
		},
		Source: "headings",
	}

	formatted := toc.FormatTOC()
	if !strings.Contains(formatted, "Chapter 1") {
		t.Errorf("Formatted TOC missing chapter")
	}
	if !strings.Contains(formatted, "Section 1.1") {
		t.Errorf("Formatted TOC missing section")
	}

	// Check for indentation (children should be indented)
	lines := strings.Split(formatted, "\n")
	var chapterLine, sectionLine int
	for i, line := range lines {
		if strings.Contains(line, "Chapter 1") {
			chapterLine = i
		}
		if strings.Contains(line, "Section 1.1") {
			sectionLine = i
		}
	}

	// Section should be indented more than chapter
	chapterIndent := len(lines[chapterLine]) - len(strings.TrimLeft(lines[chapterLine], " "))
	sectionIndent := len(lines[sectionLine]) - len(strings.TrimLeft(lines[sectionLine], " "))
	if sectionIndent <= chapterIndent {
		t.Errorf("Section should be indented more than chapter")
	}
}

// TestRegexHeadingDetector_Create verifies detector creation
func TestRegexHeadingDetector_Create(t *testing.T) {
	detector := NewRegexHeadingDetector()
	if detector == nil {
		t.Error("Failed to create detector")
	}
	if len(detector.patterns) == 0 {
		t.Error("Detector has no patterns")
	}
}

// TestRegexHeadingDetector_MatchesPattern detects chapter headings
func TestRegexHeadingDetector_MatchesPattern(t *testing.T) {
	detector := NewRegexHeadingDetector()

	tests := []struct {
		text      string
		shouldMatch bool
	}{
		{"Chapter 1", true},
		{"Chapter 42", true},
		{"Hello World", true},
		{"INTRODUCTION", true},
		{"1. Getting Started", true},
		{"this is body text", false},
		{"", false},
	}

	for _, tt := range tests {
		result := detector.MatchesPattern(tt.text)
		if result != tt.shouldMatch {
			t.Errorf("MatchesPattern(%q) = %v, want %v", tt.text, result, tt.shouldMatch)
		}
	}
}

// TestMinInt helper function
func TestMinInt(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{5, 10, 5},
		{10, 5, 5},
		{5, 5, 5},
		{0, 10, 0},
		{-5, 10, -5},
	}

	for _, tt := range tests {
		got := minInt(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("minInt(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

// BenchmarkBuildHierarchy benchmarks hierarchy building
func BenchmarkBuildHierarchy(b *testing.B) {
	toc := &TableOfContents{
		Entries: []TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1, Children: []TOCEntry{}},
			{Title: "Section 1.1", Page: 2, Level: 2, Children: []TOCEntry{}},
			{Title: "Section 1.2", Page: 5, Level: 2, Children: []TOCEntry{}},
			{Title: "Chapter 2", Page: 10, Level: 1, Children: []TOCEntry{}},
			{Title: "Section 2.1", Page: 11, Level: 2, Children: []TOCEntry{}},
		},
		Source: "headings",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = toc.BuildHierarchy()
	}
}

// BenchmarkFormatTOC benchmarks TOC formatting
func BenchmarkFormatTOC(b *testing.B) {
	toc := &TableOfContents{
		Entries: []TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
			{Title: "Chapter 2", Page: 10, Level: 1},
			{Title: "Chapter 3", Page: 20, Level: 1},
		},
		Source: "metadata",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = toc.FormatTOC()
	}
}

// BenchmarkDetectHeadingLevel benchmarks heading level detection
func BenchmarkDetectHeadingLevel(b *testing.B) {
	text := "Chapter One Advanced Topics in Computer Science and Engineering Principles"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = detectHeadingLevel(text)
	}
}
