package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luxor/lumos/pkg/config"
)

// BookmarkPane displays and manages bookmarks for the current document
type BookmarkPane struct {
	bookmarks   []config.Bookmark
	selectedIdx int
	visible     bool
	width       int
	height      int
}

// NewBookmarkPane creates a new bookmark pane
func NewBookmarkPane(width, height int) *BookmarkPane {
	return &BookmarkPane{
		bookmarks:   []config.Bookmark{},
		selectedIdx: 0,
		visible:     false,
		width:       width,
		height:      height,
	}
}

// SetBookmarks sets the list of bookmarks to display
func (bp *BookmarkPane) SetBookmarks(bookmarks []config.Bookmark) {
	bp.bookmarks = bookmarks
	bp.selectedIdx = 0
	if len(bookmarks) == 0 {
		bp.selectedIdx = 0
	}
}

// Show makes the pane visible
func (bp *BookmarkPane) Show() {
	bp.visible = true
	bp.selectedIdx = 0
}

// Hide hides the pane
func (bp *BookmarkPane) Hide() {
	bp.visible = false
}

// IsVisible returns if pane is shown
func (bp *BookmarkPane) IsVisible() bool {
	return bp.visible
}

// MoveUp moves selection up
func (bp *BookmarkPane) MoveUp() {
	if bp.selectedIdx > 0 {
		bp.selectedIdx--
	}
}

// MoveDown moves selection down
func (bp *BookmarkPane) MoveDown() {
	if bp.selectedIdx < len(bp.bookmarks)-1 {
		bp.selectedIdx++
	}
}

// GetSelectedBookmark returns the currently selected bookmark
func (bp *BookmarkPane) GetSelectedBookmark() *config.Bookmark {
	if bp.selectedIdx < 0 || bp.selectedIdx >= len(bp.bookmarks) {
		return nil
	}
	return &bp.bookmarks[bp.selectedIdx]
}

// GetSelectedPage returns the page number of selected bookmark
func (bp *BookmarkPane) GetSelectedPage() int {
	if bm := bp.GetSelectedBookmark(); bm != nil {
		return bm.Page
	}
	return 1
}

// SetSize updates pane dimensions
func (bp *BookmarkPane) SetSize(width, height int) {
	bp.width = width
	bp.height = height
}

// View renders the bookmark pane
func (bp *BookmarkPane) View() string {
	if !bp.visible {
		return ""
	}

	if len(bp.bookmarks) == 0 {
		return lipgloss.NewStyle().
			Width(bp.width).
			Height(bp.height).
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			Render("No bookmarks")
	}

	var content strings.Builder
	content.WriteString("Bookmarks:\n")
	content.WriteString(strings.Repeat("─", bp.width) + "\n")

	for i, bookmark := range bp.bookmarks {
		marker := " "
		if i == bp.selectedIdx {
			marker = "›"
		}

		note := ""
		if bookmark.Note != "" {
			note = fmt.Sprintf(" - %s", bookmark.Note)
		}

		line := fmt.Sprintf("%s Page %3d%s", marker, bookmark.Page, note)
		if i == bp.selectedIdx {
			content.WriteString(fmt.Sprintf("\033[7m%s\033[0m\n", line))
		} else {
			content.WriteString(line + "\n")
		}
	}

	return lipgloss.NewStyle().
		Width(bp.width).
		Height(bp.height).
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Render(content.String())
}

// HandleKey processes keyboard input
func (bp *BookmarkPane) HandleKey(key tea.KeyMsg) {
	switch key.Type {
	case tea.KeyUp:
		bp.MoveUp()
	case tea.KeyDown:
		bp.MoveDown()
	case tea.KeyEsc:
		bp.Hide()
	}

	// Vim navigation
	if key.Type == tea.KeyRunes && len(key.Runes) > 0 {
		switch key.Runes[0] {
		case 'k':
			bp.MoveUp()
		case 'j':
			bp.MoveDown()
		}
	}
}

// BookmarkAddedMsg signals that a bookmark was added
type BookmarkAddedMsg struct {
	Page int
	Note string
}

// BookmarkRemovedMsg signals that a bookmark was removed
type BookmarkRemovedMsg struct {
	Page int
}

// JumpToBookmarkMsg signals to jump to a bookmark
type JumpToBookmarkMsg struct {
	Page int
}
