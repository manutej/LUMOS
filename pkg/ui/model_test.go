package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luxor/lumos/pkg/pdf"
)

func TestModelInit(t *testing.T) {
	// Create a test document
	doc, err := pdf.NewDocument("../../test/fixtures/simple.pdf", 5)
	if err != nil {
		t.Fatalf("Failed to load test PDF: %v", err)
	}

	model := NewModel(doc)

	// Test initialization
	cmd := model.Init()
	if cmd == nil {
		t.Error("Init() should return a command to load the first page")
	}

	// Test initial state
	if model.currentPage != 1 {
		t.Errorf("Expected currentPage = 1, got %d", model.currentPage)
	}

	if model.theme.Name != "dark" {
		t.Errorf("Expected dark theme by default, got %s", model.theme.Name)
	}

	if model.showHelp {
		t.Error("Help should be hidden by default")
	}

	if model.searchActive {
		t.Error("Search should be inactive by default")
	}
}

func TestUpdate_QuitMessage(t *testing.T) {
	doc, err := pdf.NewDocument("../../test/fixtures/simple.pdf", 5)
	if err != nil {
		t.Fatalf("Failed to load test PDF: %v", err)
	}

	model := NewModel(doc)

	// Test quit with 'q'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(msg)

	// Should return quit command
	if cmd == nil {
		t.Error("Expected quit command, got nil")
	}
}

func TestUpdate_ToggleHelp(t *testing.T) {
	doc, err := pdf.NewDocument("../../test/fixtures/simple.pdf", 5)
	if err != nil {
		t.Fatalf("Failed to load test PDF: %v", err)
	}

	model := NewModel(doc)

	// Initially help is hidden
	if model.showHelp {
		t.Error("Help should be hidden initially")
	}

	// Press '?'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	_, _ = model.Update(msg)

	// Help should now be shown
	if !model.showHelp {
		t.Error("Help should be shown after pressing '?'")
	}

	// Press '?' again
	_, _ = model.Update(msg)

	// Help should be hidden again
	if model.showHelp {
		t.Error("Help should be hidden after pressing '?' again")
	}
}

func TestUpdate_ThemeChange(t *testing.T) {
	doc, err := pdf.NewDocument("../../test/fixtures/simple.pdf", 5)
	if err != nil {
		t.Fatalf("Failed to load test PDF: %v", err)
	}

	model := NewModel(doc)

	// Initially dark theme
	if model.theme.Name != "dark" {
		t.Errorf("Expected dark theme initially, got %s", model.theme.Name)
	}

	// Press '2' for light mode
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}}
	_, _ = model.Update(msg)

	// Should be light theme now
	if model.theme.Name != "light" {
		t.Errorf("Expected light theme after pressing '2', got %s", model.theme.Name)
	}

	// Press '1' for dark mode
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	_, _ = model.Update(msg)

	// Should be dark theme again
	if model.theme.Name != "dark" {
		t.Errorf("Expected dark theme after pressing '1', got %s", model.theme.Name)
	}
}

func TestUpdate_SearchMode(t *testing.T) {
	doc, err := pdf.NewDocument("../../test/fixtures/simple.pdf", 5)
	if err != nil {
		t.Fatalf("Failed to load test PDF: %v", err)
	}

	model := NewModel(doc)

	// Initially in normal mode
	if model.keyHandler.Mode != KeyModeNormal {
		t.Error("Should start in normal mode")
	}

	// Press '/' to enter search mode
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	_, _ = model.Update(msg)

	// Should be in search mode
	if model.keyHandler.Mode != KeyModeSearch {
		t.Error("Should be in search mode after pressing '/'")
	}

	if !model.searchActive {
		t.Error("Search should be active")
	}

	// Press escape to exit search mode
	msg = tea.KeyMsg{Type: tea.KeyEscape}
	_, _ = model.Update(msg)

	// Should be back in normal mode
	if model.keyHandler.Mode != KeyModeNormal {
		t.Error("Should be back in normal mode after escape")
	}

	if model.searchActive {
		t.Error("Search should be inactive")
	}
}

func TestUpdate_PageNavigation(t *testing.T) {
	doc, err := pdf.NewDocument("../../test/fixtures/multipage.pdf", 5)
	if err != nil {
		t.Fatalf("Failed to load test PDF: %v", err)
	}

	model := NewModel(doc)
	totalPages := doc.GetPageCount()

	// Test next page (Ctrl+N)
	initialPage := model.currentPage
	msg := tea.KeyMsg{Type: tea.KeyCtrlN}
	_, _ = model.Update(msg)

	if model.currentPage != initialPage+1 {
		t.Errorf("Expected page %d after next, got %d", initialPage+1, model.currentPage)
	}

	// Test previous page (Ctrl+P)
	msg = tea.KeyMsg{Type: tea.KeyCtrlP}
	_, _ = model.Update(msg)

	if model.currentPage != initialPage {
		t.Errorf("Expected page %d after previous, got %d", initialPage, model.currentPage)
	}

	// Test go to first page (g)
	model.currentPage = 5 // Set to middle page
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	_, _ = model.Update(msg)

	if model.currentPage != 1 {
		t.Errorf("Expected page 1 after 'g', got %d", model.currentPage)
	}

	// Test go to last page (G)
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}
	_, _ = model.Update(msg)

	if model.currentPage != totalPages {
		t.Errorf("Expected page %d after 'G', got %d", totalPages, model.currentPage)
	}
}

func TestUpdate_WindowResize(t *testing.T) {
	doc, err := pdf.NewDocument("../../test/fixtures/simple.pdf", 5)
	if err != nil {
		t.Fatalf("Failed to load test PDF: %v", err)
	}

	model := NewModel(doc)

	// Send window resize message
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	_, _ = model.Update(msg)

	// Check that dimensions were updated
	if model.width != 120 {
		t.Errorf("Expected width = 120, got %d", model.width)
	}

	if model.height != 40 {
		t.Errorf("Expected height = 40, got %d", model.height)
	}
}

func TestView_Rendering(t *testing.T) {
	doc, err := pdf.NewDocument("../../test/fixtures/simple.pdf", 5)
	if err != nil {
		t.Fatalf("Failed to load test PDF: %v", err)
	}

	model := NewModel(doc)
	model.width = 120
	model.height = 40

	// Test normal view rendering
	view := model.View()
	if view == "" {
		t.Error("View should return non-empty string")
	}

	// Test help view rendering
	model.showHelp = true
	helpView := model.View()
	if helpView == "" {
		t.Error("Help view should return non-empty string")
	}

	// Help view should contain keyboard shortcuts
	if !contains(helpView, "Navigation") {
		t.Error("Help view should contain navigation instructions")
	}
}

func TestNavigationBounds(t *testing.T) {
	doc, err := pdf.NewDocument("../../test/fixtures/simple.pdf", 5)
	if err != nil {
		t.Fatalf("Failed to load test PDF: %v", err)
	}

	model := NewModel(doc)
	model.currentPage = 1

	// Try to go to previous page from first page
	msg := tea.KeyMsg{Type: tea.KeyCtrlP}
	_, _ = model.Update(msg)

	// Should still be on page 1
	if model.currentPage != 1 {
		t.Errorf("Should not go below page 1, got %d", model.currentPage)
	}

	// Go to last page
	totalPages := doc.GetPageCount()
	model.currentPage = totalPages

	// Try to go to next page from last page
	msg = tea.KeyMsg{Type: tea.KeyCtrlN}
	_, _ = model.Update(msg)

	// Should still be on last page
	if model.currentPage != totalPages {
		t.Errorf("Should not go above page %d, got %d", totalPages, model.currentPage)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
