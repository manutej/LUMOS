package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luxor/lumos/pkg/config"
	"github.com/luxor/lumos/pkg/pdf"
)

// Model represents the main application state
type Model struct {
	// Document
	document *pdf.Document
	cache    *pdf.LRUCache

	// UI State
	currentPage       int
	theme             config.Theme
	styles            config.Styles
	keyHandler        *KeyHandler
	showHelp          bool
	searchActive      bool
	searchQuery       string
	searchResults     []pdf.SearchResult
	currentMatch      int
	activePaneIdx     int
	themeIndex        int // For cycling through themes
	clipboard         string // Store copied text
	showCopyFeedback  bool // Show copy confirmation

	// Phase 2: Table of Contents
	tocPane      *TOCPane
	showTOC      bool
	tocLoaded    bool

	// Phase 2: Advanced Search
	searchOptionsPane    *SearchOptionsPane
	searchHistoryPane    *SearchHistoryPane
	searchHistoryManager *pdf.SearchHistoryManager
	advancedSearchResults []pdf.SearchResultAdvanced
	showSearchOptions    bool
	showSearchHistory    bool

	// Phase 2.3: Bookmarks & Configuration
	cfg           *config.Config
	bookmarkPane  *BookmarkPane
	showBookmarks bool
	docPath       string // Full path to current document

	// Viewport
	viewport    viewport.Model
	metadataView viewport.Model
	searchView  viewport.Model

	// Dimensions
	width  int
	height int
}

// NewModel creates a new application model
func NewModel(document *pdf.Document) *Model {
	cache := pdf.NewLRUCache(5)

	// Load persistent configuration
	cfg, _ := config.LoadConfig()
	theme := config.GetTheme(cfg.UI.Theme)
	styles := config.NewStyles(theme)

	// Initialize viewport
	vp := viewport.New(80, 20)
	vp.Style = styles.Background

	// Initialize TOC pane
	tocPane := NewTOCPane(80, 20)

	// Initialize advanced search components
	searchHistoryManager := pdf.NewSearchHistoryManager(50)
	searchOptionsPane := NewSearchOptionsPane(80, 10)
	searchHistoryPane := NewSearchHistoryPane(80, 10, searchHistoryManager)

	// Initialize bookmark pane
	bookmarkPane := NewBookmarkPane(80, 10)

	m := &Model{
		document:             document,
		cache:                cache,
		currentPage:          1,
		theme:                theme,
		styles:               styles,
		keyHandler:           NewKeyHandler(),
		showHelp:             false,
		activePaneIdx:        1, // Start in viewer pane
		viewport:             vp,
		tocPane:              tocPane,
		showTOC:              false,
		tocLoaded:            false,
		searchOptionsPane:    searchOptionsPane,
		searchHistoryPane:    searchHistoryPane,
		searchHistoryManager: searchHistoryManager,
		advancedSearchResults: []pdf.SearchResultAdvanced{},
		showSearchOptions:    false,
		showSearchHistory:    false,
		cfg:                  cfg,
		bookmarkPane:         bookmarkPane,
		showBookmarks:        false,
	}

	return m
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	// Load the initial page
	return LoadPageCmd(m.document, m.currentPage)
}

// LoadTOC loads the table of contents from the document
func (m *Model) LoadTOC() tea.Cmd {
	return func() tea.Msg {
		if m.document != nil {
			toc, err := m.document.ExtractTableOfContents()
			if err == nil && toc != nil {
				return TOCLoadedMsg{
					TOC: toc,
				}
			}
		}
		return TOCLoadedMsg{
			TOC: &pdf.TableOfContents{
				Entries: []pdf.TOCEntry{},
				Source:  "none",
			},
		}
	}
}

// ToggleTOC toggles the table of contents display
func (m *Model) ToggleTOC() {
	if !m.tocLoaded {
		// Load TOC on first toggle
		m.LoadTOC()
	}
	m.showTOC = !m.showTOC
}

// Update handles messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowResize(msg)

	case tea.KeyMsg:
		cmd = m.handleKeyPress(msg)

	case PageLoadedMsg:
		m.handlePageLoaded(msg)

	case ScrollMsg:
		m.viewport.LineDown(msg.Amount)

	case NavigateMsg:
		cmd = m.handleNavigation(msg)

	case ThemeChangeMsg:
		m.changeTheme(msg.Theme)

	case ToggleHelpMsg:
		m.showHelp = !m.showHelp

	case ToggleTOCMsg:
		if !m.tocLoaded {
			cmd = m.LoadTOC()
		} else {
			m.showTOC = !m.showTOC
		}

	case TOCLoadedMsg:
		m.tocPane.SetTableOfContents(msg.TOC)
		m.tocLoaded = true
		m.showTOC = true

	case SearchMsg:
		cmd = m.handleSearch(msg)

	case ToggleSearchOptionsMsg:
		m.showSearchOptions = !m.showSearchOptions
		if m.showSearchOptions {
			m.searchOptionsPane.Show()
		} else {
			m.searchOptionsPane.Hide()
		}

	case ToggleSearchHistoryMsg:
		m.showSearchHistory = !m.showSearchHistory
		if m.showSearchHistory {
			m.searchHistoryPane.Show()
		} else {
			m.searchHistoryPane.Hide()
		}

	case ToggleBookmarkMsg:
		// Add or remove bookmark on current page
		if m.cfg.HasBookmark(m.docPath, m.currentPage) {
			m.cfg.RemoveBookmark(m.docPath, m.currentPage)
		} else {
			m.cfg.AddBookmark(m.docPath, m.currentPage, "")
		}
		m.cfg.Save()
		// Update bookmark pane with latest bookmarks
		m.bookmarkPane.SetBookmarks(m.cfg.GetBookmarks(m.docPath))

	case ToggleBookmarkListMsg:
		m.showBookmarks = !m.showBookmarks
		if m.showBookmarks {
			m.bookmarkPane.SetBookmarks(m.cfg.GetBookmarks(m.docPath))
			m.bookmarkPane.Show()
		} else {
			m.bookmarkPane.Hide()
		}
	}

	return m, cmd
}

// View renders the UI
func (m *Model) View() string {
	if m.showHelp {
		return m.renderHelp()
	}

	// Calculate pane widths
	metadataWidth := m.width / 5
	viewerWidth := (m.width / 10) * 6
	searchWidth := m.width - metadataWidth - viewerWidth

	// Render panes
	metadataPane := m.renderMetadataPane(metadataWidth, m.height - 2)
	viewerPane := m.renderViewerPane(viewerWidth, m.height - 2)
	searchPane := m.renderSearchPane(searchWidth, m.height - 2)

	// Combine panes horizontally
	paneContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		metadataPane,
		viewerPane,
		searchPane,
	)

	// Add status bar
	statusBar := m.renderStatusBar()

	// Combine all
	return lipgloss.JoinVertical(
		lipgloss.Left,
		paneContent,
		statusBar,
	)
}

// Handlers

func (m *Model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height

	// Resize viewport
	m.viewport.Width = msg.Width / 2
	m.viewport.Height = msg.Height - 2
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) tea.Cmd {
	if m.keyHandler.Mode == KeyModeNormal {
		switch msg.String() {
		case "q", "ctrl+c":
			return tea.Quit
		case "?":
			m.showHelp = !m.showHelp
			return nil
		case "1":
			m.changeTheme("dark")
			return nil
		case "2":
			m.changeTheme("light")
			return nil
		case "3":
			// Cycle through dark themes
			m.cycleDarkTheme()
			return nil
		case "/":
			m.keyHandler.Mode = KeyModeSearch
			m.searchActive = true
			return nil
		case "y":
			// Copy current page text to clipboard
			m.copyCurrentPage()
			return nil
		case "tab":
			m.activePaneIdx = (m.activePaneIdx + 1) % 3
			return nil
		case "shift+tab":
			m.activePaneIdx = (m.activePaneIdx - 1 + 3) % 3
			return nil
		case "n":
			if len(m.searchResults) > 0 {
				m.currentMatch = (m.currentMatch + 1) % len(m.searchResults)
				m.jumpToSearchResult()
				return nil
			}
		case "N":
			if len(m.searchResults) > 0 {
				m.currentMatch--
				if m.currentMatch < 0 {
					m.currentMatch = len(m.searchResults) - 1
				}
				m.jumpToSearchResult()
				return nil
			}
		case "j", "down":
			m.viewport.LineDown(1)
			return nil
		case "k", "up":
			m.viewport.LineUp(1)
			return nil
		case "d":
			m.viewport.LineDown(10)
			return nil
		case "u":
			m.viewport.LineUp(10)
			return nil
		case "g":
			return m.goToFirstPage()
		case "G":
			return m.goToLastPage()
		case "ctrl+f":
			m.viewport.LineDown(m.viewport.Height)
			return nil
		case "ctrl+b":
			m.viewport.LineUp(m.viewport.Height)
			return nil
		case "ctrl+n":
			return m.goToNextPage()
		case "ctrl+p":
			return m.goToPreviousPage()
		}
	} else if m.keyHandler.Mode == KeyModeSearch {
		switch msg.Type {
		case tea.KeyEscape, tea.KeyCtrlC:
			m.keyHandler.Mode = KeyModeNormal
			m.searchActive = false
			return nil
		case tea.KeyEnter:
			// Execute search
			return m.executeSearch()
		case tea.KeyBackspace:
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			}
			return nil
		case tea.KeyRunes:
			m.searchQuery += string(msg.Runes)
			return nil
		}
	}

	return nil
}

func (m *Model) handlePageLoaded(msg PageLoadedMsg) {
	m.viewport.SetContent(msg.Content)
}

func (m *Model) handleNavigation(msg NavigateMsg) tea.Cmd {
	switch msg.Type {
	case "first_page":
		return m.goToFirstPage()
	case "last_page":
		return m.goToLastPage()
	case "next_page":
		return m.goToNextPage()
	case "prev_page":
		return m.goToPreviousPage()
	}
	return nil
}

func (m *Model) handleSearch(msg SearchMsg) tea.Cmd {
	if msg.Direction == "next" {
		if len(m.searchResults) > 0 {
			m.currentMatch = (m.currentMatch + 1) % len(m.searchResults)
			m.jumpToSearchResult()
		}
	}
	return nil
}

// Navigation helpers

func (m *Model) goToFirstPage() tea.Cmd {
	m.currentPage = 1
	return LoadPageCmd(m.document, m.currentPage)
}

func (m *Model) goToLastPage() tea.Cmd {
	m.currentPage = m.document.GetPageCount()
	return LoadPageCmd(m.document, m.currentPage)
}

func (m *Model) goToNextPage() tea.Cmd {
	if m.currentPage < m.document.GetPageCount() {
		m.currentPage++
		return LoadPageCmd(m.document, m.currentPage)
	}
	return nil
}

func (m *Model) goToPreviousPage() tea.Cmd {
	if m.currentPage > 1 {
		m.currentPage--
		return LoadPageCmd(m.document, m.currentPage)
	}
	return nil
}

// Theme

func (m *Model) changeTheme(themeName string) {
	m.theme = config.GetTheme(themeName)
	m.styles = config.NewStyles(m.theme)

	// Update theme index for cycling
	for i, t := range config.AvailableThemes {
		if t.Name == m.theme.Name {
			m.themeIndex = i
			break
		}
	}
}

func (m *Model) cycleDarkTheme() {
	m.themeIndex = (m.themeIndex + 1) % len(config.AvailableThemes)
	m.theme = config.AvailableThemes[m.themeIndex]
	m.styles = config.NewStyles(m.theme)
}

// Search

func (m *Model) executeSearch() tea.Cmd {
	// This will be implemented in search handling
	return nil
}

func (m *Model) jumpToSearchResult() {
	if len(m.searchResults) > 0 && m.currentMatch < len(m.searchResults) {
		result := m.searchResults[m.currentMatch]
		m.currentPage = result.PageNum
		// Load page and scroll to result
		_ = LoadPageCmd(m.document, m.currentPage)
	}
}

// Copy - "y" key functionality

func (m *Model) copyCurrentPage() {
	pageNum := m.currentPage
	pageContent := m.viewport.View()

	// Store in clipboard (in-memory for now, can integrate with system clipboard)
	m.clipboard = fmt.Sprintf("=== Page %d ===\n%s", pageNum, pageContent)

	// Show feedback (briefly display copy confirmation)
	m.showCopyFeedback = true
}

// GetClipboard returns the currently copied text
func (m *Model) GetClipboard() string {
	return m.clipboard
}

// Rendering

func (m *Model) renderMetadataPane(width, height int) string {
	content := "📄 Document\n\n"
	meta := m.document.GetMetadata()
	content += fmt.Sprintf("Pages: %d\n", m.document.GetPageCount())
	if meta.Title != "" {
		content += "Title: " + meta.Title + "\n"
	}
	if meta.Author != "" {
		content += "Author: " + meta.Author + "\n"
	}

	paneStyle := m.styles.PaneBorder.Width(width).Height(height)
	return paneStyle.Render(content)
}

func (m *Model) renderViewerPane(width, height int) string {
	title := m.styles.PaneTitle.Render(fmt.Sprintf("📖 Viewer - Page %d", m.currentPage))
	paneStyle := m.styles.PaneBorder.Width(width).Height(height)
	return paneStyle.Render(title + "\n" + m.viewport.View())
}

func (m *Model) renderSearchPane(width, height int) string {
	title := m.styles.PaneTitle.Render("🔍 Search")
	content := fmt.Sprintf("Results: %d\n", len(m.searchResults))
	if m.searchActive {
		content += "Query: " + m.searchQuery + "\n"
	}

	paneStyle := m.styles.PaneBorder.Width(width).Height(height)
	return paneStyle.Render(title + "\n" + content)
}

func (m *Model) renderStatusBar() string {
	status := fmt.Sprintf("Page %d/%d", m.currentPage, m.document.GetPageCount())
	status += " | Theme: " + m.theme.Name

	if m.showCopyFeedback {
		status += " | [✓] Copied!"
	}

	status += " | [?] Help [q] Quit"

	return m.styles.StatusBar.Width(m.width).Render(status)
}

func (m *Model) renderHelp() string {
	helpText := "LUMOS - Dark Mode PDF Reader for Developers\n\n"
	helpText += "NAVIGATION\n"
	helpText += "  j/k or ↑/↓  - Scroll line up/down\n"
	helpText += "  d/u         - Scroll half page up/down\n"
	helpText += "  Ctrl+F/B    - Full page down/up\n"
	helpText += "  gg/G        - Go to first/last page\n"
	helpText += "  Ctrl+N/P    - Next/previous page\n\n"
	helpText += "SEARCH & COPY\n"
	helpText += "  /           - Start search\n"
	helpText += "  n/N         - Next/previous match\n"
	helpText += "  y           - Copy current page\n"
	helpText += "  Esc         - Exit search\n\n"
	helpText += "THEMES (Professional Dark Modes)\n"
	helpText += "  1           - Switch to dark theme\n"
	helpText += "  2           - Switch to light theme\n"
	helpText += "  3           - Cycle through dark themes\n"
	helpText += "    Available: LUMOS Dark, Tokyo Night, Dracula, Solarized, Nord\n\n"
	helpText += "UI CONTROLS\n"
	helpText += "  Tab/Shift+Tab - Cycle panes forward/backward\n"
	helpText += "  ?           - Toggle this help screen\n"
	helpText += "  q/Ctrl+C    - Quit\n\n"
	helpText += "Press ? to close this help"

	return m.styles.Background.Width(m.width).Height(m.height).Render(helpText)
}

// Command definitions

type PageLoadedMsg struct {
	Content string
}

func LoadPageCmd(doc *pdf.Document, pageNum int) tea.Cmd {
	return func() tea.Msg {
		page, err := doc.GetPage(pageNum)
		if err != nil {
			return PageLoadedMsg{Content: "Error loading page: " + err.Error()}
		}
		return PageLoadedMsg{Content: page.Text}
	}
}
