package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luxor/lumos/pkg/pdf"
)

// TOCPane manages the table of contents display and navigation
type TOCPane struct {
	toc           *pdf.TableOfContents
	selectedIdx   int
	flatEntries   []pdf.TOCEntry // Flattened TOC for easier navigation
	viewport      viewport.Model
	width         int
	height        int
	selectedEntry *pdf.TOCEntry
}

// NewTOCPane creates a new TOC pane
func NewTOCPane(width, height int) *TOCPane {
	v := viewport.New(width, height)
	return &TOCPane{
		toc:         &pdf.TableOfContents{Entries: []pdf.TOCEntry{}, Source: "none"},
		viewport:    v,
		width:       width,
		height:      height,
		selectedIdx: 0,
	}
}

// SetTableOfContents sets the TOC to display
func (tp *TOCPane) SetTableOfContents(toc *pdf.TableOfContents) {
	tp.toc = toc
	tp.flatEntries = tp.flattenTOC(toc.Entries)
	tp.selectedIdx = 0
	tp.updateViewport()
}

// flattenTOC converts hierarchical TOC to flat list for easier navigation
func (tp *TOCPane) flattenTOC(entries []pdf.TOCEntry) []pdf.TOCEntry {
	var flat []pdf.TOCEntry

	var flatten func([]pdf.TOCEntry)
	flatten = func(entries []pdf.TOCEntry) {
		for _, entry := range entries {
			flat = append(flat, entry)
			if len(entry.Children) > 0 {
				flatten(entry.Children)
			}
		}
	}

	flatten(entries)
	return flat
}

// MoveUp moves selection up in TOC
func (tp *TOCPane) MoveUp() {
	if tp.selectedIdx > 0 {
		tp.selectedIdx--
		tp.updateViewport()
	}
}

// MoveDown moves selection down in TOC
func (tp *TOCPane) MoveDown() {
	if tp.selectedIdx < len(tp.flatEntries)-1 {
		tp.selectedIdx++
		tp.updateViewport()
	}
}

// MovePageUp moves selection up by page
func (tp *TOCPane) MovePageUp() {
	newIdx := tp.selectedIdx - (tp.height / 2)
	if newIdx < 0 {
		newIdx = 0
	}
	tp.selectedIdx = newIdx
	tp.updateViewport()
}

// MovePageDown moves selection down by page
func (tp *TOCPane) MovePageDown() {
	newIdx := tp.selectedIdx + (tp.height / 2)
	if newIdx >= len(tp.flatEntries) {
		newIdx = len(tp.flatEntries) - 1
	}
	tp.selectedIdx = newIdx
	tp.updateViewport()
}

// GetSelectedEntry returns the currently selected TOC entry
func (tp *TOCPane) GetSelectedEntry() *pdf.TOCEntry {
	if tp.selectedIdx < 0 || tp.selectedIdx >= len(tp.flatEntries) {
		return nil
	}
	return &tp.flatEntries[tp.selectedIdx]
}

// GetSelectedPage returns the page number of the selected entry
func (tp *TOCPane) GetSelectedPage() int {
	entry := tp.GetSelectedEntry()
	if entry != nil {
		return entry.Page
	}
	return 1 // Default to first page
}

// updateViewport updates the viewport content
func (tp *TOCPane) updateViewport() {
	content := tp.renderTOC()
	tp.viewport.SetContent(content)
}

// renderTOC renders the TOC with proper formatting and highlighting
func (tp *TOCPane) renderTOC() string {
	var buf strings.Builder

	if tp.toc == nil || len(tp.flatEntries) == 0 {
		buf.WriteString("No table of contents available")
		return buf.String()
	}

	buf.WriteString(fmt.Sprintf("Table of Contents (%s)\n", tp.toc.Source))
	buf.WriteString(strings.Repeat("─", tp.width) + "\n")

	for i, entry := range tp.flatEntries {
		line := tp.formatTOCEntry(entry, i)
		buf.WriteString(line + "\n")
	}

	return buf.String()
}

// formatTOCEntry formats a single TOC entry with indentation and selection indicator
func (tp *TOCPane) formatTOCEntry(entry pdf.TOCEntry, idx int) string {
	// Indentation based on level
	indent := strings.Repeat("  ", entry.Level-1)

	// Selection indicator
	marker := " "
	if idx == tp.selectedIdx {
		marker = "›"
	}

	// Format: > Chapter Title ........................... Page 42
	title := entry.Title
	if len(title) > tp.width-30 {
		title = title[:tp.width-30]
	}

	pageNum := fmt.Sprintf("p.%d", entry.Page)
	dots := tp.width - len(indent) - 2 - len(title) - len(pageNum) - 2
	if dots < 1 {
		dots = 1
	}

	line := fmt.Sprintf("%s%s %s %s %s",
		indent, marker, title, strings.Repeat(".", dots), pageNum)

	// Highlight selected entry
	if idx == tp.selectedIdx {
		return fmt.Sprintf("\033[7m%s\033[0m", line) // Reverse video for selection
	}

	return line
}

// SetSize updates pane size
func (tp *TOCPane) SetSize(width, height int) {
	tp.width = width
	tp.height = height
	tp.viewport.Width = width
	tp.viewport.Height = height
	tp.updateViewport()
}

// View renders the TOC pane
func (tp *TOCPane) View() string {
	if tp.toc == nil || len(tp.flatEntries) == 0 {
		return lipgloss.NewStyle().
			Width(tp.width).
			Height(tp.height).
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			Render("No TOC available")
	}

	return lipgloss.NewStyle().
		Width(tp.width).
		Height(tp.height).
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Render(tp.viewport.View())
}

// HandleKey processes keyboard input
func (tp *TOCPane) HandleKey(key tea.KeyMsg) {
	switch key.Type {
	case tea.KeyUp:
		tp.MoveUp()
	case tea.KeyDown:
		tp.MoveDown()
	case tea.KeyHome:
		tp.selectedIdx = 0
		tp.updateViewport()
	case tea.KeyEnd:
		if len(tp.flatEntries) > 0 {
			tp.selectedIdx = len(tp.flatEntries) - 1
			tp.updateViewport()
		}
	}

	// Handle rune keys for vim navigation
	if key.Type == tea.KeyRunes && len(key.Runes) > 0 {
		switch key.Runes[0] {
		case 'k':
			tp.MoveUp()
		case 'j':
			tp.MoveDown()
		case 'u':
			tp.MovePageUp()
		case 'd':
			tp.MovePageDown()
		case 'g':
			tp.selectedIdx = 0
			tp.updateViewport()
		case 'G':
			if len(tp.flatEntries) > 0 {
				tp.selectedIdx = len(tp.flatEntries) - 1
				tp.updateViewport()
			}
		}
	}
}

// TOCPaneMsg represents a TOC pane event
type TOCPaneMsg struct {
	Page int // Page number selected in TOC
}

// TOCChangedMsg is sent when TOC is updated
type TOCChangedMsg struct {
	TOC *pdf.TableOfContents
}

// SearchTOC searches the TOC for entries matching a query
func (tp *TOCPane) SearchTOC(query string) []pdf.TOCEntry {
	var results []pdf.TOCEntry

	query = strings.ToLower(query)
	for _, entry := range tp.flatEntries {
		if strings.Contains(strings.ToLower(entry.Title), query) {
			results = append(results, entry)
		}
	}

	return results
}

// JumpToPage selects the TOC entry that corresponds to a page
func (tp *TOCPane) JumpToPage(pageNum int) {
	// Find entries for this page
	for i, entry := range tp.flatEntries {
		if entry.Page == pageNum {
			tp.selectedIdx = i
			tp.updateViewport()
			return
		}
	}
}

// IsEmpty returns true if TOC has no entries
func (tp *TOCPane) IsEmpty() bool {
	return len(tp.flatEntries) == 0
}

// GetEntryCount returns the number of entries in the TOC
func (tp *TOCPane) GetEntryCount() int {
	return len(tp.flatEntries)
}

// GetSourceType returns the source of the TOC (metadata, headings, none)
func (tp *TOCPane) GetSourceType() string {
	if tp.toc != nil {
		return tp.toc.Source
	}
	return "none"
}
