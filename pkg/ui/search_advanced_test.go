package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luxor/lumos/pkg/pdf"
)

// TestNewSearchOptionsPane creates new options pane
func TestNewSearchOptionsPane(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)

	if pane == nil {
		t.Error("Failed to create search options pane")
	}
	if pane.width != 80 {
		t.Errorf("Width mismatch: got %d, want %d", pane.width, 80)
	}
	if pane.height != 24 {
		t.Errorf("Height mismatch: got %d, want %d", pane.height, 24)
	}
	if pane.IsVisible() {
		t.Error("Pane should not be visible initially")
	}
}

// TestSearchOptionsPane_GetOptions returns current options
func TestSearchOptionsPane_GetOptions(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)
	opts := pane.GetOptions()

	if opts.CaseSensitive {
		t.Error("CaseSensitive should be false by default")
	}
	if opts.WholeWord {
		t.Error("WholeWord should be false by default")
	}
	if opts.RegexMode {
		t.Error("RegexMode should be false by default")
	}
}

// TestSearchOptionsPane_ToggleCaseSensitive toggles case sensitivity
func TestSearchOptionsPane_ToggleCaseSensitive(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)

	if pane.options.CaseSensitive {
		t.Error("Should start with case-insensitive")
	}

	pane.ToggleCaseSensitive()
	if !pane.options.CaseSensitive {
		t.Error("Should be case-sensitive after toggle")
	}

	pane.ToggleCaseSensitive()
	if pane.options.CaseSensitive {
		t.Error("Should be case-insensitive after second toggle")
	}
}

// TestSearchOptionsPane_ToggleWholeWord toggles whole-word matching
func TestSearchOptionsPane_ToggleWholeWord(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)

	if pane.options.WholeWord {
		t.Error("Should start without whole-word")
	}

	pane.ToggleWholeWord()
	if !pane.options.WholeWord {
		t.Error("Should have whole-word after toggle")
	}
}

// TestSearchOptionsPane_ToggleRegexMode toggles regex
func TestSearchOptionsPane_ToggleRegexMode(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)

	if pane.options.RegexMode {
		t.Error("Should start without regex")
	}

	pane.ToggleRegexMode()
	if !pane.options.RegexMode {
		t.Error("Should have regex after toggle")
	}

	// Regex should disable whole-word
	pane.options.WholeWord = true
	pane.ToggleRegexMode()
	pane.ToggleRegexMode()
	if pane.options.WholeWord {
		t.Error("Whole-word should be disabled when regex is enabled")
	}
}

// TestSearchOptionsPane_SetQuery sets search query
func TestSearchOptionsPane_SetQuery(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)

	pane.SetQuery("test")
	if pane.searchQuery != "test" {
		t.Errorf("Query should be 'test', got %s", pane.searchQuery)
	}
}

// TestSearchOptionsPane_Show shows pane
func TestSearchOptionsPane_Show(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)

	if pane.IsVisible() {
		t.Error("Pane should not be visible initially")
	}

	pane.Show()
	if !pane.IsVisible() {
		t.Error("Pane should be visible after Show()")
	}
}

// TestSearchOptionsPane_Hide hides pane
func TestSearchOptionsPane_Hide(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)
	pane.Show()

	pane.Hide()
	if pane.IsVisible() {
		t.Error("Pane should not be visible after Hide()")
	}
}

// TestSearchOptionsPane_MoveUp navigates up
func TestSearchOptionsPane_MoveUp(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)
	pane.selectedOptionIdx = 1

	pane.MoveUp()
	if pane.selectedOptionIdx != 0 {
		t.Errorf("Should move to index 0, got %d", pane.selectedOptionIdx)
	}

	pane.MoveUp()
	if pane.selectedOptionIdx != 0 {
		t.Errorf("Should stay at index 0, got %d", pane.selectedOptionIdx)
	}
}

// TestSearchOptionsPane_MoveDown navigates down
func TestSearchOptionsPane_MoveDown(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)
	pane.selectedOptionIdx = 0

	pane.MoveDown()
	if pane.selectedOptionIdx != 1 {
		t.Errorf("Should move to index 1, got %d", pane.selectedOptionIdx)
	}

	pane.MoveDown()
	if pane.selectedOptionIdx != 2 {
		t.Errorf("Should move to index 2, got %d", pane.selectedOptionIdx)
	}

	pane.MoveDown()
	if pane.selectedOptionIdx != 2 {
		t.Errorf("Should stay at index 2, got %d", pane.selectedOptionIdx)
	}
}

// TestSearchOptionsPane_SelectOption toggles selected option
func TestSearchOptionsPane_SelectOption(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)

	// Toggle case-sensitive (index 0)
	if pane.options.CaseSensitive {
		t.Error("Should start without case-sensitive")
	}
	pane.SelectOption()
	if !pane.options.CaseSensitive {
		t.Error("Case-sensitive should be enabled")
	}

	// Move to whole-word (index 1)
	pane.MoveDown()
	pane.SelectOption()
	if !pane.options.WholeWord {
		t.Error("Whole-word should be enabled")
	}

	// Move to regex (index 2)
	pane.MoveDown()
	pane.SelectOption()
	if !pane.options.RegexMode {
		t.Error("Regex should be enabled")
	}
}

// TestSearchOptionsPane_SetPageRange sets page range
func TestSearchOptionsPane_SetPageRange(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)

	pane.SetPageRange(5, 15)
	if pane.options.StartPage != 5 {
		t.Errorf("StartPage should be 5, got %d", pane.options.StartPage)
	}
	if pane.options.EndPage != 15 {
		t.Errorf("EndPage should be 15, got %d", pane.options.EndPage)
	}
}

// TestSearchOptionsPane_SetSize updates dimensions
func TestSearchOptionsPane_SetSize(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)

	pane.SetSize(120, 30)
	if pane.width != 120 {
		t.Errorf("Width should be 120, got %d", pane.width)
	}
	if pane.height != 30 {
		t.Errorf("Height should be 30, got %d", pane.height)
	}
}

// TestSearchOptionsPane_View renders pane
func TestSearchOptionsPane_View(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)

	// View should be empty when not visible
	view := pane.View()
	if view != "" {
		t.Error("View should be empty when pane is not visible")
	}

	// View should contain options when visible
	pane.Show()
	view = pane.View()
	if !strings.Contains(view, "Case Sensitive") {
		t.Error("View should contain 'Case Sensitive'")
	}
	if !strings.Contains(view, "Whole Word") {
		t.Error("View should contain 'Whole Word'")
	}
	if !strings.Contains(view, "Regex") {
		t.Error("View should contain 'Regex'")
	}
}

// TestSearchOptionsPane_GetStatusLine returns option flags
func TestSearchOptionsPane_GetStatusLine(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)

	// Empty status when all options off
	status := pane.GetStatusLine()
	if status != "" {
		t.Errorf("Status should be empty, got %s", status)
	}

	// Show flags when options enabled
	pane.ToggleCaseSensitive()
	status = pane.GetStatusLine()
	if !strings.Contains(status, "Aa") {
		t.Errorf("Status should contain 'Aa', got %s", status)
	}

	pane.ToggleWholeWord()
	status = pane.GetStatusLine()
	if !strings.Contains(status, "W") {
		t.Errorf("Status should contain 'W', got %s", status)
	}

	pane.ToggleRegexMode()
	status = pane.GetStatusLine()
	if !strings.Contains(status, ".*") {
		t.Errorf("Status should contain '.*', got %s", status)
	}
}

// TestSearchOptionsPane_HandleKey processes keyboard
func TestSearchOptionsPane_HandleKey(t *testing.T) {
	pane := NewSearchOptionsPane(80, 24)

	// Test down key
	msg := tea.KeyMsg{Type: tea.KeyDown}
	pane.HandleKey(msg)
	if pane.selectedOptionIdx != 1 {
		t.Errorf("Should move to index 1, got %d", pane.selectedOptionIdx)
	}

	// Test up key
	msg = tea.KeyMsg{Type: tea.KeyUp}
	pane.HandleKey(msg)
	if pane.selectedOptionIdx != 0 {
		t.Errorf("Should move to index 0, got %d", pane.selectedOptionIdx)
	}

	// Test space to toggle
	msg = tea.KeyMsg{Type: tea.KeySpace}
	pane.HandleKey(msg)
	if !pane.options.CaseSensitive {
		t.Error("Space should toggle option")
	}
}

// TestNewSearchHistoryPane creates new history pane
func TestNewSearchHistoryPane(t *testing.T) {
	manager := pdf.NewSearchHistoryManager(10)
	pane := NewSearchHistoryPane(80, 24, manager)

	if pane == nil {
		t.Error("Failed to create search history pane")
	}
	if pane.IsVisible() {
		t.Error("Pane should not be visible initially")
	}
}

// TestSearchHistoryPane_Show shows pane
func TestSearchHistoryPane_Show(t *testing.T) {
	manager := pdf.NewSearchHistoryManager(10)
	pane := NewSearchHistoryPane(80, 24, manager)

	pane.Show()
	if !pane.IsVisible() {
		t.Error("Pane should be visible after Show()")
	}
}

// TestSearchHistoryPane_Hide hides pane
func TestSearchHistoryPane_Hide(t *testing.T) {
	manager := pdf.NewSearchHistoryManager(10)
	pane := NewSearchHistoryPane(80, 24, manager)
	pane.Show()

	pane.Hide()
	if pane.IsVisible() {
		t.Error("Pane should not be visible after Hide()")
	}
}

// TestSearchHistoryPane_GetSelectedQuery returns selected query
func TestSearchHistoryPane_GetSelectedQuery(t *testing.T) {
	manager := pdf.NewSearchHistoryManager(10)
	manager.Add(pdf.SearchHistoryEntry{Query: "test1"})
	manager.Add(pdf.SearchHistoryEntry{Query: "test2"})

	pane := NewSearchHistoryPane(80, 24, manager)

	// Most recent is at index 0
	query := pane.GetSelectedQuery()
	if query != "test2" {
		t.Errorf("Selected query should be 'test2', got %s", query)
	}

	pane.MoveDown()
	query = pane.GetSelectedQuery()
	if query != "test1" {
		t.Errorf("Selected query should be 'test1' after move, got %s", query)
	}
}

// TestSearchHistoryPane_View renders history
func TestSearchHistoryPane_View(t *testing.T) {
	manager := pdf.NewSearchHistoryManager(10)
	manager.Add(pdf.SearchHistoryEntry{Query: "test", ResultCount: 5})

	pane := NewSearchHistoryPane(80, 24, manager)
	pane.Show()

	view := pane.View()
	if !strings.Contains(view, "test") {
		t.Error("View should contain search query")
	}
	if !strings.Contains(view, "5 results") {
		t.Error("View should contain result count")
	}
}

// TestSearchHistoryPane_HandleKey processes keyboard
func TestSearchHistoryPane_HandleKey(t *testing.T) {
	manager := pdf.NewSearchHistoryManager(10)
	manager.Add(pdf.SearchHistoryEntry{Query: "test1"})
	manager.Add(pdf.SearchHistoryEntry{Query: "test2"})

	pane := NewSearchHistoryPane(80, 24, manager)

	// Test down key
	msg := tea.KeyMsg{Type: tea.KeyDown}
	pane.HandleKey(msg)
	if pane.selectedIdx != 1 {
		t.Errorf("Should move to index 1, got %d", pane.selectedIdx)
	}

	// Test up key
	msg = tea.KeyMsg{Type: tea.KeyUp}
	pane.HandleKey(msg)
	if pane.selectedIdx != 0 {
		t.Errorf("Should move to index 0, got %d", pane.selectedIdx)
	}
}

// BenchmarkSearchOptionsPane_SelectOption benchmarks option selection
func BenchmarkSearchOptionsPane_SelectOption(b *testing.B) {
	pane := NewSearchOptionsPane(80, 24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pane.SelectOption()
	}
}

// BenchmarkSearchOptionsPane_GetStatusLine benchmarks status line rendering
func BenchmarkSearchOptionsPane_GetStatusLine(b *testing.B) {
	pane := NewSearchOptionsPane(80, 24)
	pane.ToggleCaseSensitive()
	pane.ToggleWholeWord()
	pane.ToggleRegexMode()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pane.GetStatusLine()
	}
}

// BenchmarkSearchHistoryPane_GetSelectedQuery benchmarks query retrieval
func BenchmarkSearchHistoryPane_GetSelectedQuery(b *testing.B) {
	manager := pdf.NewSearchHistoryManager(100)
	for i := 0; i < 50; i++ {
		manager.Add(pdf.SearchHistoryEntry{Query: "query"})
	}
	pane := NewSearchHistoryPane(80, 24, manager)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pane.GetSelectedQuery()
	}
}
