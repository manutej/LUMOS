package ui

import "github.com/charmbracelet/bubbletea"

// KeyHandler handles vim-style keybindings
type KeyHandler struct {
	Mode KeyMode
}

// KeyMode represents the current key handling mode
type KeyMode int

const (
	KeyModeNormal KeyMode = iota
	KeyModeSearch
	KeyModeCommand
)

// NewKeyHandler creates a new key handler
func NewKeyHandler() *KeyHandler {
	return &KeyHandler{
		Mode: KeyModeNormal,
	}
}

// HandleKey processes a key press
func (kh *KeyHandler) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch kh.Mode {
	case KeyModeNormal:
		return kh.handleNormalKey(msg)
	case KeyModeSearch:
		return kh.handleSearchKey(msg)
	case KeyModeCommand:
		return kh.handleCommandKey(msg)
	}
	return nil
}

// handleNormalKey processes keys in normal mode
func (kh *KeyHandler) handleNormalKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	// Navigation
	case "j", "down":
		return ScrollDown
	case "k", "up":
		return ScrollUp
	case "d":
		return HalfPageDown
	case "u":
		return HalfPageUp
	case "g":
		return GoToTop
	case "G", "shift+g":
		return GoToBottom

	// Page navigation
	case "ctrl+n":
		return NextPage
	case "ctrl+p":
		return PrevPage

	// Search
	case "/":
		return EnterSearch
	case "ctrl+\\", "ctrl+/":
		return ToggleSearchOptions
	case "ctrl+h":
		return ToggleSearchHistory

	// UI
	case "tab":
		return CyclePane
	case "shift+tab":
		return CyclePanePrev
	case "1":
		return ToggleDarkMode
	case "2":
		return ToggleLightMode
	case "3":
		return CycleDarkTheme
	case "ctrl+t":
		return ToggleTOC
	case "m":
		return ToggleBookmark
	case "'", "`":
		return ToggleBookmarkList
	case "?":
		return ToggleHelp

	// Exit
	case "q", "ctrl+c":
		return Exit

	// Commands
	case ":":
		return EnterCommand
	}

	return nil
}

// handleSearchKey processes keys in search mode
func (kh *KeyHandler) handleSearchKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter":
		return ExecuteSearch
	case "escape", "ctrl+c":
		return ExitSearch
	case "n":
		return NextMatch
	case "N", "shift+n":
		return PrevMatch
	}

	return nil
}

// handleCommandKey processes keys in command mode
func (kh *KeyHandler) handleCommandKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter":
		return ExecuteCommand
	case "escape", "ctrl+c":
		return ExitCommand
	}

	return nil
}

// Command functions for key bindings

var (
	ScrollDown = func() tea.Msg {
		return ScrollMsg{Direction: "down", Amount: 1}
	}

	ScrollUp = func() tea.Msg {
		return ScrollMsg{Direction: "up", Amount: 1}
	}

	HalfPageDown = func() tea.Msg {
		return ScrollMsg{Direction: "down", Amount: 10}
	}

	HalfPageUp = func() tea.Msg {
		return ScrollMsg{Direction: "up", Amount: 10}
	}

	GoToTop = func() tea.Msg {
		return NavigateMsg{Type: "first_page"}
	}

	GoToBottom = func() tea.Msg {
		return NavigateMsg{Type: "last_page"}
	}

	NextPage = func() tea.Msg {
		return NavigateMsg{Type: "next_page"}
	}

	PrevPage = func() tea.Msg {
		return NavigateMsg{Type: "prev_page"}
	}

	EnterSearch = func() tea.Msg {
		return ModeChangeMsg{Mode: KeyModeSearch}
	}

	ExitSearch = func() tea.Msg {
		return ModeChangeMsg{Mode: KeyModeNormal}
	}

	ExecuteSearch = func() tea.Msg {
		return SearchMsg{}
	}

	NextMatch = func() tea.Msg {
		return SearchMsg{Direction: "next"}
	}

	PrevMatch = func() tea.Msg {
		return SearchMsg{Direction: "prev"}
	}

	EnterCommand = func() tea.Msg {
		return ModeChangeMsg{Mode: KeyModeCommand}
	}

	ExitCommand = func() tea.Msg {
		return ModeChangeMsg{Mode: KeyModeNormal}
	}

	ExecuteCommand = func() tea.Msg {
		return CommandMsg{}
	}

	CyclePane = func() tea.Msg {
		return PaneChangeMsg{Direction: "next"}
	}

	CyclePanePrev = func() tea.Msg {
		return PaneChangeMsg{Direction: "prev"}
	}

	ToggleDarkMode = func() tea.Msg {
		return ThemeChangeMsg{Theme: "dark"}
	}

	ToggleLightMode = func() tea.Msg {
		return ThemeChangeMsg{Theme: "light"}
	}

	CycleDarkTheme = func() tea.Msg {
		return ThemeChangeMsg{Theme: "cycle"}
	}

	ToggleTOC = func() tea.Msg {
		return ToggleTOCMsg{}
	}

	ToggleSearchOptions = func() tea.Msg {
		return ToggleSearchOptionsMsg{}
	}

	ToggleSearchHistory = func() tea.Msg {
		return ToggleSearchHistoryMsg{}
	}

	ToggleBookmark = func() tea.Msg {
		return ToggleBookmarkMsg{}
	}

	ToggleBookmarkList = func() tea.Msg {
		return ToggleBookmarkListMsg{}
	}

	ToggleHelp = func() tea.Msg {
		return ToggleHelpMsg{}
	}

	Exit = func() tea.Msg {
		return tea.QuitMsg{}
	}
)

// UI Messages

type ScrollMsg struct {
	Direction string // "up" or "down"
	Amount    int
}

type NavigateMsg struct {
	Type string // "first_page", "last_page", "next_page", "prev_page"
}

type SearchMsg struct {
	Direction string // "next" or "prev"
	Query     string
}

type ModeChangeMsg struct {
	Mode KeyMode
}

type CommandMsg struct {
	Command string
}

type PaneChangeMsg struct {
	Direction string // "next" or "prev"
}

type ThemeChangeMsg struct {
	Theme string // "dark", "light", or "cycle"
}

type ToggleHelpMsg struct{}

type ToggleTOCMsg struct{}

type ToggleSearchOptionsMsg struct{}

type ToggleSearchHistoryMsg struct{}

type ToggleBookmarkMsg struct{}

type ToggleBookmarkListMsg struct{}

// VimKeybindingReference provides a reference of all keybindings
var VimKeybindingReference = map[string]string{
	// Navigation - Line scrolling
	"j/↓":           "Scroll down one line",
	"k/↑":           "Scroll up one line",
	"h/←":           "Scroll left (wide PDFs)",
	"l/→":           "Scroll right (wide PDFs)",

	// Navigation - Page scrolling
	"d":             "Scroll down half page",
	"u":             "Scroll up half page",
	"Ctrl+F":        "Scroll down full page",
	"Ctrl+B":        "Scroll up full page",

	// Navigation - Page jumps
	"g":             "Go to first page",
	"G":             "Go to last page",
	"Ctrl+N":        "Go to next page",
	"Ctrl+P":        "Go to previous page",

	// Search & Copy
	"/":             "Start search",
	"Ctrl+\\":       "Toggle search options (regex, case, whole-word)",
	"Ctrl+H":        "Toggle search history",
	"n":             "Go to next match",
	"N":             "Go to previous match",
	"y":             "Copy current page text",
	"Esc":           "Exit search",

	// Themes (Dark mode - core feature)
	"1":             "Switch to dark theme",
	"2":             "Switch to light theme",
	"3":             "Cycle through dark themes",

	// Bookmarks
	"m":             "Add/toggle bookmark on current page",
	"'/backtick":    "Show bookmark list",

	// UI Controls
	"Tab":           "Cycle through panes (forward)",
	"Shift+Tab":     "Cycle through panes (backward)",
	"Ctrl+T":        "Toggle table of contents",
	"?":             "Toggle help screen",

	// General
	"q/Ctrl+C":      "Quit application",
}
