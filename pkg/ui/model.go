package ui

import (
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
	currentPage    int
	theme          config.Theme
	styles         config.Styles
	keyHandler     *KeyHandler
	showHelp       bool
	searchActive   bool
	searchQuery    string
	searchResults  []pdf.SearchResult
	currentMatch   int
	activePaneIdx  int

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
	theme := config.DarkTheme
	styles := config.NewStyles(theme)

	m := &Model{
		document:   document,
		cache:      cache,
		currentPage: 1,
		theme:      theme,
		styles:     styles,
		keyHandler: NewKeyHandler(),
		showHelp:   false,
		activePaneIdx: 1, // Start in viewer pane
	}

	return m
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return LoadPageCmd(m.document, m.currentPage)
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

	case SearchMsg:
		cmd = m.handleSearch(msg)
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
		case "/":
			m.keyHandler.Mode = KeyModeSearch
			m.searchActive = true
			return nil
		case "tab":
			m.activePaneIdx = (m.activePaneIdx + 1) % 3
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
		case "ctrl+n":
			return m.goToNextPage()
		case "ctrl+p":
			return m.goToPreviousPage()
		}
	} else if m.keyHandler.Mode == KeyModeSearch {
		switch msg.String() {
		case "escape":
			m.keyHandler.Mode = KeyModeNormal
			m.searchActive = false
			return nil
		case "enter":
			// Execute search
			return m.executeSearch()
		case "backspace":
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			}
			return nil
		default:
			m.searchQuery += msg.String()
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

// Rendering

func (m *Model) renderMetadataPane(width, height int) string {
	content := "📄 Document\n\n"
	meta := m.document.GetMetadata()
	content += "Pages: " + string(rune(m.document.GetPageCount())) + "\n"
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
	title := m.styles.PaneTitle.Render("📖 Viewer - Page " + string(rune(m.currentPage)))
	paneStyle := m.styles.PaneBorder.Width(width).Height(height)
	return paneStyle.Render(title + "\n" + m.viewport.View())
}

func (m *Model) renderSearchPane(width, height int) string {
	title := m.styles.PaneTitle.Render("🔍 Search")
	content := "Results: " + string(rune(len(m.searchResults))) + "\n"
	if m.searchActive {
		content += "Query: " + m.searchQuery + "\n"
	}

	paneStyle := m.styles.PaneBorder.Width(width).Height(height)
	return paneStyle.Render(title + "\n" + content)
}

func (m *Model) renderStatusBar() string {
	status := "Page " + string(rune(m.currentPage)) + "/" + string(rune(m.document.GetPageCount()))
	status += " | Theme: " + m.theme.Name
	status += " | [?] Help [q] Quit"

	return m.styles.StatusBar.Width(m.width).Render(status)
}

func (m *Model) renderHelp() string {
	helpText := "LUMOS - PDF Dark Mode Reader\n\n"
	helpText += "Navigation:\n"
	helpText += "  j/k or ↑/↓  - Scroll\n"
	helpText += "  d/u         - Half page\n"
	helpText += "  gg/G        - Top/bottom\n"
	helpText += "  Ctrl+N/P    - Next/prev page\n\n"
	helpText += "Search:\n"
	helpText += "  /           - Search\n"
	helpText += "  n/N         - Next/prev match\n\n"
	helpText += "UI:\n"
	helpText += "  Tab         - Cycle panes\n"
	helpText += "  1/2         - Dark/light mode\n"
	helpText += "  ?           - Toggle help\n"
	helpText += "  q           - Quit\n"

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
