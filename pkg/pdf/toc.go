package pdf

import (
	"fmt"
	"regexp"
	"strings"
)

// TableOfContents represents the hierarchical structure of a PDF document
type TableOfContents struct {
	Entries []TOCEntry
	Source  string // "metadata", "headings", or "none"
}

// TOCEntry represents a single entry in the table of contents
type TOCEntry struct {
	Title    string
	Page     int
	Children []TOCEntry
	Level    int // 1=chapter, 2=section, 3=subsection, etc.
}

// ExtractTableOfContents attempts to extract TOC from PDF metadata or generates one from content
func (d *Document) ExtractTableOfContents() (*TableOfContents, error) {
	if d == nil {
		return nil, fmt.Errorf("document is nil")
	}

	// Try to extract from PDF metadata first (future enhancement)
	// For now, generate from detected headings in content
	toc := &TableOfContents{
		Entries: []TOCEntry{},
		Source:  "none",
	}

	// Scan first 10 pages for heading patterns
	headingEntries := make(map[int][]TOCEntry)
	foundHeadings := false

	for pageNum := 1; pageNum <= minInt(10, d.GetPageCount()); pageNum++ {
		headings := d.FindHeadings(pageNum)
		if len(headings) > 0 {
			foundHeadings = true
			headingEntries[pageNum] = headings
		}
	}

	if foundHeadings {
		// Build hierarchical structure from detected headings
		for pageNum := 1; pageNum <= d.GetPageCount(); pageNum++ {
			if headings, ok := headingEntries[pageNum]; ok {
				for _, heading := range headings {
					heading.Page = pageNum
					toc.Entries = append(toc.Entries, heading)
				}
			}
		}
		toc.Source = "headings"
	}

	return toc, nil
}

// FindHeadings analyzes a page and detects potential headings
// Uses heuristics: line starting with capital letters, short lines, isolated text
func (d *Document) FindHeadings(pageNum int) []TOCEntry {
	page, err := d.GetPage(pageNum)
	if err != nil {
		return []TOCEntry{}
	}

	var headings []TOCEntry
	lines := TextToLines(page.Text)

	// Patterns for detecting headings
	// Capital letters at start, typically short
	for i, line := range lines {
		if len(line.Text) == 0 {
			continue
		}

		text := strings.TrimSpace(line.Text)
		if len(text) == 0 {
			continue
		}

		// Skip if line is too long (likely body text)
		if len(text) > 100 {
			continue
		}

		// Check if starts with capital letter
		if len(text) > 0 && text[0] >= 'A' && text[0] <= 'Z' {
			// Detect level based on indentation or length heuristics
			level := detectHeadingLevel(text)

			// Skip lines that are just single words (likely not headings)
			wordCount := len(strings.Fields(text))
			if wordCount < 2 && level == 1 {
				continue // Single word at level 1 is probably not a chapter title
			}

			// Create TOC entry
			entry := TOCEntry{
				Title:    text,
				Page:     pageNum,
				Level:    level,
				Children: []TOCEntry{},
			}

			// Add to headings if not a duplicate of previous
			if len(headings) == 0 || headings[len(headings)-1].Title != text {
				headings = append(headings, entry)
			}
		}

		// Limit to reasonable number of headings per page
		if len(headings) >= 5 {
			break
		}

		// Stop after scanning first 50 lines
		if i > 50 {
			break
		}
	}

	return headings
}

// detectHeadingLevel determines the heading level based on text characteristics
func detectHeadingLevel(text string) int {
	// Very short text (< 5 words) → Level 1 (chapter)
	words := strings.Fields(text)
	if len(words) <= 2 {
		return 1
	}

	// 3-5 words → Level 2 (section)
	if len(words) <= 5 {
		return 2
	}

	// 6+ words → Level 3 (subsection)
	return 3
}

// BuildHierarchy organizes flat TOC entries into a proper hierarchy
// based on heading levels (1=top, 2=sub, 3=subsub, etc.)
func (t *TableOfContents) BuildHierarchy() *TableOfContents {
	return buildHierarchySimple(t)
}

// buildHierarchySimple uses a cleaner algorithm to build hierarchy
func buildHierarchySimple(t *TableOfContents) *TableOfContents {
	if len(t.Entries) == 0 {
		return t
	}

	result := &TableOfContents{
		Entries: []TOCEntry{},
		Source:  t.Source,
	}

	// Recursive function to build hierarchy
	var entries = t.Entries

	// Find all root-level entries (level 1)
	for i := 0; i < len(entries); i++ {
		if entries[i].Level == 1 {
			// Create new entry
			newEntry := TOCEntry{
				Title:    entries[i].Title,
				Page:     entries[i].Page,
				Level:    entries[i].Level,
				Children: []TOCEntry{},
			}

			// Find children of this entry (next entries until we hit another level 1)
			childStart := i + 1
			childEnd := childStart
			for childEnd < len(entries) && entries[childEnd].Level > 1 {
				childEnd++
			}

			// Add children
			for j := childStart; j < childEnd; j++ {
				child := TOCEntry{
					Title:    entries[j].Title,
					Page:     entries[j].Page,
					Level:    entries[j].Level,
					Children: []TOCEntry{},
				}
				newEntry.Children = append(newEntry.Children, child)
			}

			result.Entries = append(result.Entries, newEntry)
		}
	}

	return result
}

// FindEntryByPage returns all TOC entries that point to a given page
func (t *TableOfContents) FindEntryByPage(pageNum int) []TOCEntry {
	var results []TOCEntry

	var search func([]TOCEntry, int)
	search = func(entries []TOCEntry, pageNum int) {
		for _, entry := range entries {
			if entry.Page == pageNum {
				results = append(results, entry)
			}
			if len(entry.Children) > 0 {
				search(entry.Children, pageNum)
			}
		}
	}

	search(t.Entries, pageNum)
	return results
}

// GetPageForEntry returns the target page number for a TOC entry
func (e *TOCEntry) GetPageForEntry() int {
	return e.Page
}

// FormatTOC returns a formatted string representation of the TOC
func (t *TableOfContents) FormatTOC() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("Table of Contents (from %s)\n\n", t.Source))

	var format func([]TOCEntry, int)
	format = func(entries []TOCEntry, depth int) {
		for _, entry := range entries {
			indent := strings.Repeat("  ", depth)
			buf.WriteString(fmt.Sprintf("%s• %s (page %d)\n", indent, entry.Title, entry.Page))
			if len(entry.Children) > 0 {
				format(entry.Children, depth+1)
			}
		}
	}

	format(t.Entries, 0)
	return buf.String()
}

// Helper function
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RegexHeadingDetector detects headings using regex patterns
// Useful for PDFs with structured formatting
type RegexHeadingDetector struct {
	patterns []*regexp.Regexp
}

// NewRegexHeadingDetector creates a detector with common heading patterns
func NewRegexHeadingDetector() *RegexHeadingDetector {
	return &RegexHeadingDetector{
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`^Chapter\s+\d+`),
			regexp.MustCompile(`^[A-Z][A-Za-z\s]{2,}$`),
			regexp.MustCompile(`^\d+\.\s+[A-Z]`),
		},
	}
}

// MatchesPattern returns true if text matches any heading pattern
func (d *RegexHeadingDetector) MatchesPattern(text string) bool {
	for _, pattern := range d.patterns {
		if pattern.MatchString(text) {
			return true
		}
	}
	return false
}
