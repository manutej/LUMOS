package ui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luxor/lumos/pkg/pdf"
)

// TestNewTOCPane creates a new TOC pane
func TestNewTOCPane(t *testing.T) {
	pane := NewTOCPane(80, 24)

	if pane == nil {
		t.Error("Failed to create TOC pane")
	}
	if pane.width != 80 {
		t.Errorf("Width mismatch: got %d, want %d", pane.width, 80)
	}
	if pane.height != 24 {
		t.Errorf("Height mismatch: got %d, want %d", pane.height, 24)
	}
	if pane.selectedIdx != 0 {
		t.Errorf("Initial index should be 0, got %d", pane.selectedIdx)
	}
}

// TestTOCPane_SetTableOfContents sets TOC data
func TestTOCPane_SetTableOfContents(t *testing.T) {
	pane := NewTOCPane(80, 24)
	toc := &pdf.TableOfContents{
		Entries: []pdf.TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
			{Title: "Chapter 2", Page: 10, Level: 1},
		},
		Source: "metadata",
	}

	pane.SetTableOfContents(toc)

	if pane.toc != toc {
		t.Error("TOC not set correctly")
	}
	if len(pane.flatEntries) != 2 {
		t.Errorf("Expected 2 flat entries, got %d", len(pane.flatEntries))
	}
}

// TestTOCPane_FlattenTOC flattens hierarchical TOC
func TestTOCPane_FlattenTOC(t *testing.T) {
	pane := NewTOCPane(80, 24)
	entries := []pdf.TOCEntry{
		{
			Title: "Chapter 1",
			Page:  1,
			Level: 1,
			Children: []pdf.TOCEntry{
				{Title: "Section 1.1", Page: 2, Level: 2},
				{Title: "Section 1.2", Page: 5, Level: 2},
			},
		},
		{Title: "Chapter 2", Page: 10, Level: 1},
	}

	flat := pane.flattenTOC(entries)

	// Should have 4 entries when flattened (1 chapter, 2 sections, 1 chapter)
	if len(flat) != 4 {
		t.Errorf("Expected 4 flattened entries, got %d", len(flat))
	}

	// Check order
	if flat[0].Title != "Chapter 1" {
		t.Errorf("First entry should be Chapter 1, got %s", flat[0].Title)
	}
	if flat[1].Title != "Section 1.1" {
		t.Errorf("Second entry should be Section 1.1, got %s", flat[1].Title)
	}
	if flat[3].Title != "Chapter 2" {
		t.Errorf("Last entry should be Chapter 2, got %s", flat[3].Title)
	}
}

// TestTOCPane_MoveUp navigates up in TOC
func TestTOCPane_MoveUp(t *testing.T) {
	pane := NewTOCPane(80, 24)
	toc := &pdf.TableOfContents{
		Entries: []pdf.TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
			{Title: "Chapter 2", Page: 10, Level: 1},
			{Title: "Chapter 3", Page: 20, Level: 1},
		},
		Source: "metadata",
	}
	pane.SetTableOfContents(toc)

	// Start at index 2
	pane.selectedIdx = 2
	pane.MoveUp()

	if pane.selectedIdx != 1 {
		t.Errorf("Expected index 1 after move up, got %d", pane.selectedIdx)
	}

	// Can't move up past 0
	pane.selectedIdx = 0
	pane.MoveUp()
	if pane.selectedIdx != 0 {
		t.Errorf("Should stay at 0 when at top, got %d", pane.selectedIdx)
	}
}

// TestTOCPane_MoveDown navigates down in TOC
func TestTOCPane_MoveDown(t *testing.T) {
	pane := NewTOCPane(80, 24)
	toc := &pdf.TableOfContents{
		Entries: []pdf.TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
			{Title: "Chapter 2", Page: 10, Level: 1},
			{Title: "Chapter 3", Page: 20, Level: 1},
		},
		Source: "metadata",
	}
	pane.SetTableOfContents(toc)

	// Start at index 0
	pane.selectedIdx = 0
	pane.MoveDown()

	if pane.selectedIdx != 1 {
		t.Errorf("Expected index 1 after move down, got %d", pane.selectedIdx)
	}

	// Can't move down past last
	pane.selectedIdx = 2
	pane.MoveDown()
	if pane.selectedIdx != 2 {
		t.Errorf("Should stay at last index, got %d", pane.selectedIdx)
	}
}

// TestTOCPane_GetSelectedEntry returns current selection
func TestTOCPane_GetSelectedEntry(t *testing.T) {
	pane := NewTOCPane(80, 24)
	toc := &pdf.TableOfContents{
		Entries: []pdf.TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
			{Title: "Chapter 2", Page: 10, Level: 1},
		},
		Source: "metadata",
	}
	pane.SetTableOfContents(toc)

	pane.selectedIdx = 1
	entry := pane.GetSelectedEntry()

	if entry == nil {
		t.Error("Selected entry is nil")
	}
	if entry.Title != "Chapter 2" {
		t.Errorf("Expected Chapter 2, got %s", entry.Title)
	}
}

// TestTOCPane_GetSelectedPage returns page of selected entry
func TestTOCPane_GetSelectedPage(t *testing.T) {
	pane := NewTOCPane(80, 24)
	toc := &pdf.TableOfContents{
		Entries: []pdf.TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
			{Title: "Chapter 2", Page: 10, Level: 1},
		},
		Source: "metadata",
	}
	pane.SetTableOfContents(toc)

	pane.selectedIdx = 1
	page := pane.GetSelectedPage()

	if page != 10 {
		t.Errorf("Expected page 10, got %d", page)
	}
}

// TestTOCPane_JumpToPage navigates to page
func TestTOCPane_JumpToPage(t *testing.T) {
	pane := NewTOCPane(80, 24)
	toc := &pdf.TableOfContents{
		Entries: []pdf.TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
			{Title: "Chapter 2", Page: 10, Level: 1},
			{Title: "Chapter 3", Page: 20, Level: 1},
		},
		Source: "metadata",
	}
	pane.SetTableOfContents(toc)

	pane.JumpToPage(10)

	if pane.selectedIdx != 1 {
		t.Errorf("Expected to jump to index 1 for page 10, got %d", pane.selectedIdx)
	}

	pane.JumpToPage(20)
	if pane.selectedIdx != 2 {
		t.Errorf("Expected to jump to index 2 for page 20, got %d", pane.selectedIdx)
	}
}

// TestTOCPane_SearchTOC searches entries
func TestTOCPane_SearchTOC(t *testing.T) {
	pane := NewTOCPane(80, 24)
	toc := &pdf.TableOfContents{
		Entries: []pdf.TOCEntry{
			{Title: "Chapter One Introduction", Page: 1, Level: 1},
			{Title: "Chapter Two Advanced Topics", Page: 10, Level: 1},
			{Title: "Appendix A Reference", Page: 50, Level: 1},
		},
		Source: "metadata",
	}
	pane.SetTableOfContents(toc)

	results := pane.SearchTOC("Chapter")

	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'Chapter', got %d", len(results))
	}

	results = pane.SearchTOC("Advanced")
	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'Advanced', got %d", len(results))
	}

	results = pane.SearchTOC("NotFound")
	if len(results) != 0 {
		t.Errorf("Expected 0 results for 'NotFound', got %d", len(results))
	}
}

// TestTOCPane_IsEmpty checks if TOC is empty
func TestTOCPane_IsEmpty(t *testing.T) {
	pane := NewTOCPane(80, 24)

	if !pane.IsEmpty() {
		t.Error("New pane should be empty")
	}

	toc := &pdf.TableOfContents{
		Entries: []pdf.TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
		},
		Source: "metadata",
	}
	pane.SetTableOfContents(toc)

	if pane.IsEmpty() {
		t.Error("Pane should not be empty after setting TOC")
	}
}

// TestTOCPane_GetEntryCount returns entry count
func TestTOCPane_GetEntryCount(t *testing.T) {
	pane := NewTOCPane(80, 24)

	if pane.GetEntryCount() != 0 {
		t.Errorf("Expected 0 entries, got %d", pane.GetEntryCount())
	}

	toc := &pdf.TableOfContents{
		Entries: []pdf.TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
			{Title: "Chapter 2", Page: 10, Level: 1},
		},
		Source: "metadata",
	}
	pane.SetTableOfContents(toc)

	if pane.GetEntryCount() != 2 {
		t.Errorf("Expected 2 entries, got %d", pane.GetEntryCount())
	}
}

// TestTOCPane_GetSourceType returns TOC source
func TestTOCPane_GetSourceType(t *testing.T) {
	pane := NewTOCPane(80, 24)

	if pane.GetSourceType() != "none" {
		t.Errorf("Expected 'none' source, got %s", pane.GetSourceType())
	}

	toc := &pdf.TableOfContents{
		Entries: []pdf.TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
		},
		Source: "metadata",
	}
	pane.SetTableOfContents(toc)

	if pane.GetSourceType() != "metadata" {
		t.Errorf("Expected 'metadata' source, got %s", pane.GetSourceType())
	}
}

// TestTOCPane_SetSize updates pane dimensions
func TestTOCPane_SetSize(t *testing.T) {
	pane := NewTOCPane(80, 24)

	pane.SetSize(120, 30)

	if pane.width != 120 {
		t.Errorf("Width should be 120, got %d", pane.width)
	}
	if pane.height != 30 {
		t.Errorf("Height should be 30, got %d", pane.height)
	}
}

// TestTOCPane_View renders pane
func TestTOCPane_View(t *testing.T) {
	pane := NewTOCPane(80, 24)
	view := pane.View()

	if view == "" {
		t.Error("View should not be empty")
	}
	if !strings.Contains(view, "No TOC available") {
		t.Error("View should indicate no TOC when empty")
	}

	toc := &pdf.TableOfContents{
		Entries: []pdf.TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
		},
		Source: "metadata",
	}
	pane.SetTableOfContents(toc)
	view = pane.View()

	if !strings.Contains(view, "Chapter 1") {
		t.Error("View should contain TOC entry")
	}
}

// TestTOCPane_HandleKey processes keyboard
func TestTOCPane_HandleKey(t *testing.T) {
	pane := NewTOCPane(80, 24)
	toc := &pdf.TableOfContents{
		Entries: []pdf.TOCEntry{
			{Title: "Chapter 1", Page: 1, Level: 1},
			{Title: "Chapter 2", Page: 10, Level: 1},
			{Title: "Chapter 3", Page: 20, Level: 1},
		},
		Source: "metadata",
	}
	pane.SetTableOfContents(toc)

	// Test down arrow key
	msg := tea.KeyMsg{Type: tea.KeyDown}
	pane.HandleKey(msg)
	if pane.selectedIdx != 1 {
		t.Errorf("Down arrow should move to 1, got %d", pane.selectedIdx)
	}

	// Test up arrow key
	msg = tea.KeyMsg{Type: tea.KeyUp}
	pane.HandleKey(msg)
	if pane.selectedIdx != 0 {
		t.Errorf("Up arrow should move to 0, got %d", pane.selectedIdx)
	}

	// Test 'j' key (vim down)
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	pane.HandleKey(msg)
	if pane.selectedIdx != 1 {
		t.Errorf("'j' should move down, got %d", pane.selectedIdx)
	}

	// Test 'k' key (vim up)
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	pane.HandleKey(msg)
	if pane.selectedIdx != 0 {
		t.Errorf("'k' should move up, got %d", pane.selectedIdx)
	}
}

// TestTOCPane_MovePageUp moves by page
func TestTOCPane_MovePageUp(t *testing.T) {
	pane := NewTOCPane(80, 10) // Small height
	entries := make([]pdf.TOCEntry, 20)
	for i := 0; i < 20; i++ {
		entries[i] = pdf.TOCEntry{
			Title: fmt.Sprintf("Chapter %d", i+1),
			Page:  i + 1,
			Level: 1,
		}
	}
	toc := &pdf.TableOfContents{Entries: entries, Source: "metadata"}
	pane.SetTableOfContents(toc)

	pane.selectedIdx = 15
	pane.MovePageUp()

	if pane.selectedIdx >= 15 {
		t.Errorf("Page up should decrease selection, got %d", pane.selectedIdx)
	}
}

// TestTOCPane_MovePageDown moves by page
func TestTOCPane_MovePageDown(t *testing.T) {
	pane := NewTOCPane(80, 10) // Small height
	entries := make([]pdf.TOCEntry, 20)
	for i := 0; i < 20; i++ {
		entries[i] = pdf.TOCEntry{
			Title: fmt.Sprintf("Chapter %d", i+1),
			Page:  i + 1,
			Level: 1,
		}
	}
	toc := &pdf.TableOfContents{Entries: entries, Source: "metadata"}
	pane.SetTableOfContents(toc)

	pane.selectedIdx = 5
	pane.MovePageDown()

	if pane.selectedIdx <= 5 {
		t.Errorf("Page down should increase selection, got %d", pane.selectedIdx)
	}
}

// BenchmarkTOCPane_FormatEntry benchmarks entry formatting
func BenchmarkTOCPane_FormatEntry(b *testing.B) {
	pane := NewTOCPane(80, 24)
	entry := pdf.TOCEntry{
		Title: "Chapter One Advanced Topics in Computer Science",
		Page:  42,
		Level: 1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pane.formatTOCEntry(entry, i%10)
	}
}

// BenchmarkTOCPane_Flatten benchmarks TOC flattening
func BenchmarkTOCPane_Flatten(b *testing.B) {
	pane := NewTOCPane(80, 24)
	entries := []pdf.TOCEntry{
		{
			Title: "Chapter 1",
			Page:  1,
			Level: 1,
			Children: []pdf.TOCEntry{
				{Title: "Section 1.1", Page: 2, Level: 2},
				{Title: "Section 1.2", Page: 5, Level: 2},
			},
		},
		{Title: "Chapter 2", Page: 10, Level: 1},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pane.flattenTOC(entries)
	}
}

// Ensure code compiles
var _ = fmt.Sprintf
