package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luxor/lumos/pkg/pdf"
)

// SearchOptionsPane manages advanced search settings UI
type SearchOptionsPane struct {
	options          pdf.SearchOptions
	selectedOptionIdx int
	width            int
	height           int
	visible          bool
	searchQuery      string
}

// SearchOption represents a toggleable search option
type SearchOption struct {
	Name        string
	Description string
	Value       bool
	Key         string // Keyboard shortcut
}

// NewSearchOptionsPane creates a new search options pane
func NewSearchOptionsPane(width, height int) *SearchOptionsPane {
	return &SearchOptionsPane{
		options: pdf.SearchOptions{
			CaseSensitive: false,
			WholeWord:     false,
			RegexMode:     false,
			MaxResults:    0,
			StartPage:     1,
			EndPage:       0,
		},
		selectedOptionIdx: 0,
		width:            width,
		height:           height,
		visible:          false,
		searchQuery:      "",
	}
}

// GetOptions returns current search options
func (sop *SearchOptionsPane) GetOptions() pdf.SearchOptions {
	return sop.options
}

// SetQuery sets the current search query
func (sop *SearchOptionsPane) SetQuery(query string) {
	sop.searchQuery = query
}

// ToggleCaseSensitive toggles case-sensitive matching
func (sop *SearchOptionsPane) ToggleCaseSensitive() {
	sop.options.CaseSensitive = !sop.options.CaseSensitive
}

// ToggleWholeWord toggles whole-word matching
func (sop *SearchOptionsPane) ToggleWholeWord() {
	sop.options.WholeWord = !sop.options.WholeWord
}

// ToggleRegexMode toggles regex matching
func (sop *SearchOptionsPane) ToggleRegexMode() {
	sop.options.RegexMode = !sop.options.RegexMode
	// Disable whole-word when regex is enabled
	if sop.options.RegexMode {
		sop.options.WholeWord = false
	}
}

// SetMaxResults sets maximum results limit
func (sop *SearchOptionsPane) SetMaxResults(max int) {
	sop.options.MaxResults = max
}

// SetPageRange sets the search page range
func (sop *SearchOptionsPane) SetPageRange(start, end int) {
	sop.options.StartPage = start
	sop.options.EndPage = end
}

// Show makes the pane visible
func (sop *SearchOptionsPane) Show() {
	sop.visible = true
	sop.selectedOptionIdx = 0
}

// Hide hides the pane
func (sop *SearchOptionsPane) Hide() {
	sop.visible = false
}

// IsVisible returns whether pane is shown
func (sop *SearchOptionsPane) IsVisible() bool {
	return sop.visible
}

// MoveUp moves selection up
func (sop *SearchOptionsPane) MoveUp() {
	if sop.selectedOptionIdx > 0 {
		sop.selectedOptionIdx--
	}
}

// MoveDown moves selection down
func (sop *SearchOptionsPane) MoveDown() {
	optionCount := 3 // case-sensitive, whole-word, regex
	if sop.selectedOptionIdx < optionCount-1 {
		sop.selectedOptionIdx++
	}
}

// SelectOption toggles the currently selected option
func (sop *SearchOptionsPane) SelectOption() {
	switch sop.selectedOptionIdx {
	case 0:
		sop.ToggleCaseSensitive()
	case 1:
		sop.ToggleWholeWord()
	case 2:
		sop.ToggleRegexMode()
	}
}

// SetSize updates pane dimensions
func (sop *SearchOptionsPane) SetSize(width, height int) {
	sop.width = width
	sop.height = height
}

// View renders the search options pane
func (sop *SearchOptionsPane) View() string {
	if !sop.visible {
		return ""
	}

	options := []SearchOption{
		{
			Name:        "Case Sensitive",
			Description: "Match case (Ctrl+S)",
			Value:       sop.options.CaseSensitive,
			Key:         "c",
		},
		{
			Name:        "Whole Word",
			Description: "Match whole words only (Ctrl+W)",
			Value:       sop.options.WholeWord,
			Key:         "w",
		},
		{
			Name:        "Regex Mode",
			Description: "Use regex patterns (Ctrl+R)",
			Value:       sop.options.RegexMode,
			Key:         "r",
		},
	}

	var content strings.Builder
	content.WriteString(fmt.Sprintf("Search Query: %s\n", sop.searchQuery))
	content.WriteString(strings.Repeat("─", sop.width) + "\n")
	content.WriteString("Search Options:\n\n")

	for i, opt := range options {
		marker := " "
		if i == sop.selectedOptionIdx {
			marker = "›"
		}

		checkbox := "☐"
		if opt.Value {
			checkbox = "☑"
		}

		line := fmt.Sprintf("%s %s %s - %s", marker, checkbox, opt.Name, opt.Description)
		if i == sop.selectedOptionIdx {
			// Highlight selected option
			content.WriteString(fmt.Sprintf("\033[7m%s\033[0m\n", line))
		} else {
			content.WriteString(line + "\n")
		}
	}

	content.WriteString("\n")
	content.WriteString("Controls:\n")
	content.WriteString("  j/↓: Next option\n")
	content.WriteString("  k/↑: Previous option\n")
	content.WriteString("  Space/Enter: Toggle option\n")
	content.WriteString("  Esc: Close options\n")

	return lipgloss.NewStyle().
		Width(sop.width).
		Height(sop.height).
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Render(content.String())
}

// HandleKey processes keyboard input
func (sop *SearchOptionsPane) HandleKey(key tea.KeyMsg) {
	switch key.Type {
	case tea.KeyUp:
		sop.MoveUp()
	case tea.KeyDown:
		sop.MoveDown()
	case tea.KeySpace, tea.KeyEnter:
		sop.SelectOption()
	case tea.KeyEsc:
		sop.Hide()
	}

	// Handle rune keys for vim navigation
	if key.Type == tea.KeyRunes && len(key.Runes) > 0 {
		switch key.Runes[0] {
		case 'k':
			sop.MoveUp()
		case 'j':
			sop.MoveDown()
		case ' ':
			sop.SelectOption()
		case 'c':
			sop.ToggleCaseSensitive()
		case 'w':
			sop.ToggleWholeWord()
		case 'r':
			sop.ToggleRegexMode()
		}
	}
}

// GetStatusLine returns a concise status of current options
func (sop *SearchOptionsPane) GetStatusLine() string {
	var flags []string
	if sop.options.CaseSensitive {
		flags = append(flags, "Aa")
	}
	if sop.options.WholeWord {
		flags = append(flags, "W")
	}
	if sop.options.RegexMode {
		flags = append(flags, ".*")
	}

	if len(flags) == 0 {
		return ""
	}

	return fmt.Sprintf("[%s]", strings.Join(flags, "|"))
}

// SearchHistoryPane displays search history
type SearchHistoryPane struct {
	history       *pdf.SearchHistoryManager
	selectedIdx   int
	visible       bool
	width         int
	height        int
}

// NewSearchHistoryPane creates a new history pane
func NewSearchHistoryPane(width, height int, manager *pdf.SearchHistoryManager) *SearchHistoryPane {
	return &SearchHistoryPane{
		history:     manager,
		selectedIdx: 0,
		visible:     false,
		width:       width,
		height:      height,
	}
}

// Show makes the pane visible
func (shp *SearchHistoryPane) Show() {
	shp.visible = true
	shp.selectedIdx = 0
}

// Hide hides the pane
func (shp *SearchHistoryPane) Hide() {
	shp.visible = false
}

// IsVisible returns whether pane is shown
func (shp *SearchHistoryPane) IsVisible() bool {
	return shp.visible
}

// MoveUp moves selection up
func (shp *SearchHistoryPane) MoveUp() {
	if shp.selectedIdx > 0 {
		shp.selectedIdx--
	}
}

// MoveDown moves selection down
func (shp *SearchHistoryPane) MoveDown() {
	if shp.selectedIdx < shp.history.Size()-1 {
		shp.selectedIdx++
	}
}

// GetSelectedQuery returns the currently selected search query
func (shp *SearchHistoryPane) GetSelectedQuery() string {
	if shp.history.Size() == 0 {
		return ""
	}
	history := shp.history.GetHistory()
	if shp.selectedIdx < len(history) {
		return history[shp.selectedIdx].Query
	}
	return ""
}

// SetSize updates pane dimensions
func (shp *SearchHistoryPane) SetSize(width, height int) {
	shp.width = width
	shp.height = height
}

// View renders the history pane
func (shp *SearchHistoryPane) View() string {
	if !shp.visible || shp.history.Size() == 0 {
		return lipgloss.NewStyle().
			Width(shp.width).
			Height(shp.height).
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			Render("No search history")
	}

	var content strings.Builder
	content.WriteString("Search History:\n")
	content.WriteString(strings.Repeat("─", shp.width) + "\n")

	history := shp.history.GetHistory()
	for i, entry := range history {
		marker := " "
		if i == shp.selectedIdx {
			marker = "›"
		}

		line := fmt.Sprintf("%s %s (%d results)", marker, entry.Query, entry.ResultCount)
		if i == shp.selectedIdx {
			content.WriteString(fmt.Sprintf("\033[7m%s\033[0m\n", line))
		} else {
			content.WriteString(line + "\n")
		}
	}

	return lipgloss.NewStyle().
		Width(shp.width).
		Height(shp.height).
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Render(content.String())
}

// HandleKey processes keyboard input
func (shp *SearchHistoryPane) HandleKey(key tea.KeyMsg) {
	switch key.Type {
	case tea.KeyUp:
		shp.MoveUp()
	case tea.KeyDown:
		shp.MoveDown()
	case tea.KeyEsc:
		shp.Hide()
	}

	// Vim navigation
	if key.Type == tea.KeyRunes && len(key.Runes) > 0 {
		switch key.Runes[0] {
		case 'k':
			shp.MoveUp()
		case 'j':
			shp.MoveDown()
		}
	}
}
