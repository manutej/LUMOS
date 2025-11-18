package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luxor/lumos/pkg/config"
)

// TestNewBookmarkPane creates new pane
func TestNewBookmarkPane(t *testing.T) {
	pane := NewBookmarkPane(80, 24)

	if pane == nil {
		t.Error("Failed to create bookmark pane")
	}
	if pane.visible {
		t.Error("Pane should not be visible initially")
	}
	if len(pane.bookmarks) != 0 {
		t.Error("Bookmarks should be empty initially")
	}
}

// TestSetBookmarks populates bookmarks
func TestSetBookmarks(t *testing.T) {
	pane := NewBookmarkPane(80, 24)
	bookmarks := []config.Bookmark{
		{Page: 5, Note: "Chapter 1"},
		{Page: 10, Note: "Methods"},
		{Page: 20, Note: "Results"},
	}

	pane.SetBookmarks(bookmarks)

	if len(pane.bookmarks) != 3 {
		t.Errorf("Expected 3 bookmarks, got %d", len(pane.bookmarks))
	}
	if pane.selectedIdx != 0 {
		t.Error("Selection should reset to 0")
	}
}

// TestMoveUp navigates up
func TestBookmarkPane_MoveUp(t *testing.T) {
	pane := NewBookmarkPane(80, 24)
	pane.SetBookmarks([]config.Bookmark{
		{Page: 5, Note: "First"},
		{Page: 10, Note: "Second"},
	})
	pane.selectedIdx = 1

	pane.MoveUp()
	if pane.selectedIdx != 0 {
		t.Errorf("Should move to 0, got %d", pane.selectedIdx)
	}
}

// TestMoveDown navigates down
func TestBookmarkPane_MoveDown(t *testing.T) {
	pane := NewBookmarkPane(80, 24)
	pane.SetBookmarks([]config.Bookmark{
		{Page: 5, Note: "First"},
		{Page: 10, Note: "Second"},
	})

	pane.MoveDown()
	if pane.selectedIdx != 1 {
		t.Errorf("Should move to 1, got %d", pane.selectedIdx)
	}
}

// TestGetSelectedBookmark returns current selection
func TestGetSelectedBookmark(t *testing.T) {
	pane := NewBookmarkPane(80, 24)
	bookmarks := []config.Bookmark{
		{Page: 5, Note: "First"},
		{Page: 10, Note: "Second"},
	}
	pane.SetBookmarks(bookmarks)
	pane.selectedIdx = 1

	bm := pane.GetSelectedBookmark()
	if bm == nil || bm.Page != 10 {
		t.Errorf("Selected bookmark should be page 10, got %+v", bm)
	}
}

// TestGetSelectedPage returns page number
func TestGetSelectedPage(t *testing.T) {
	pane := NewBookmarkPane(80, 24)
	pane.SetBookmarks([]config.Bookmark{
		{Page: 5, Note: "First"},
		{Page: 20, Note: "Second"},
	})
	pane.selectedIdx = 1

	page := pane.GetSelectedPage()
	if page != 20 {
		t.Errorf("Expected page 20, got %d", page)
	}
}

// TestShow makes pane visible
func TestBookmarkPane_Show(t *testing.T) {
	pane := NewBookmarkPane(80, 24)

	pane.Show()
	if !pane.IsVisible() {
		t.Error("Pane should be visible after Show()")
	}
}

// TestHide hides pane
func TestBookmarkPane_Hide(t *testing.T) {
	pane := NewBookmarkPane(80, 24)
	pane.Show()

	pane.Hide()
	if pane.IsVisible() {
		t.Error("Pane should not be visible after Hide()")
	}
}

// TestView renders pane
func TestBookmarkPane_View(t *testing.T) {
	pane := NewBookmarkPane(80, 24)
	pane.SetBookmarks([]config.Bookmark{
		{Page: 5, Note: "Chapter 1"},
		{Page: 10, Note: "Methods"},
	})
	pane.Show()

	view := pane.View()
	if !strings.Contains(view, "Page") {
		t.Error("View should contain 'Page'")
	}
	if !strings.Contains(view, "Chapter 1") {
		t.Error("View should contain bookmark note")
	}
}

// TestViewEmpty handles no bookmarks
func TestBookmarkPane_ViewEmpty(t *testing.T) {
	pane := NewBookmarkPane(80, 24)
	pane.Show()

	view := pane.View()
	if !strings.Contains(view, "No bookmarks") {
		t.Error("Empty view should say 'No bookmarks'")
	}
}

// TestHandleKey processes keyboard
func TestBookmarkPane_HandleKey(t *testing.T) {
	pane := NewBookmarkPane(80, 24)
	pane.SetBookmarks([]config.Bookmark{
		{Page: 5, Note: "First"},
		{Page: 10, Note: "Second"},
	})

	// Test down arrow
	msg := tea.KeyMsg{Type: tea.KeyDown}
	pane.HandleKey(msg)
	if pane.selectedIdx != 1 {
		t.Errorf("Down should move to 1, got %d", pane.selectedIdx)
	}

	// Test up arrow
	msg = tea.KeyMsg{Type: tea.KeyUp}
	pane.HandleKey(msg)
	if pane.selectedIdx != 0 {
		t.Errorf("Up should move to 0, got %d", pane.selectedIdx)
	}
}

// BenchmarkView benchmarks rendering
func BenchmarkBookmarkPane_View(b *testing.B) {
	pane := NewBookmarkPane(80, 24)
	bookmarks := []config.Bookmark{}
	for i := 0; i < 50; i++ {
		bookmarks = append(bookmarks, config.Bookmark{Page: i + 1, Note: "Bookmark"})
	}
	pane.SetBookmarks(bookmarks)
	pane.Show()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pane.View()
	}
}
