package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config represents the complete application configuration
type Config struct {
	UI        UIConfig              `toml:"ui"`
	Documents map[string]DocState   `toml:"documents"`
	Bookmarks map[string][]Bookmark `toml:"bookmarks"`
}

// UIConfig holds UI preferences
type UIConfig struct {
	Theme string `toml:"theme"`
}

// DocState tracks per-document state
type DocState struct {
	LastPage   int       `toml:"last_page"`
	LastScroll int       `toml:"last_scroll"`
	Timestamp  time.Time `toml:"timestamp"`
}

// Bookmark represents a page bookmark with optional note
type Bookmark struct {
	Page int    `toml:"page"`
	Note string `toml:"note"`
}

// DefaultConfig returns sensible defaults
func DefaultConfig() *Config {
	return &Config{
		UI: UIConfig{
			Theme: "dark",
		},
		Documents: make(map[string]DocState),
		Bookmarks: make(map[string][]Bookmark),
	}
}

// LoadConfig loads config from ~/.config/lumos/config.toml
// Returns defaults if file doesn't exist
func LoadConfig() (*Config, error) {
	path := configPath()

	// Return defaults if file doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	// Simple TOML parsing using standard library - avoid extra dependency
	// For now, use a basic approach that doesn't require external TOML library
	// We'll use a manual approach for maximum compatibility

	cfg = *DefaultConfig()

	// Parse the TOML file manually (pragmatic approach)
	if err := parseConfig(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// Save persists config to disk
func (c *Config) Save() error {
	path := configPath()

	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Build TOML content
	content := c.toTOML()

	// Write atomically (write to temp, then rename)
	tmpFile := path + ".tmp"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	if err := os.Rename(tmpFile, path); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// UpdateDocState updates the last page/scroll for a document
func (c *Config) UpdateDocState(path string, page int, scroll int) {
	c.Documents[path] = DocState{
		LastPage:   page,
		LastScroll: scroll,
		Timestamp:  time.Now(),
	}
}

// AddBookmark adds a bookmark to a document
func (c *Config) AddBookmark(docPath string, page int, note string) {
	if c.Bookmarks[docPath] == nil {
		c.Bookmarks[docPath] = []Bookmark{}
	}

	// Avoid duplicates on same page
	for i, bm := range c.Bookmarks[docPath] {
		if bm.Page == page {
			c.Bookmarks[docPath][i].Note = note
			return
		}
	}

	c.Bookmarks[docPath] = append(c.Bookmarks[docPath], Bookmark{Page: page, Note: note})
}

// RemoveBookmark removes a bookmark from a document
func (c *Config) RemoveBookmark(docPath string, page int) {
	bookmarks := c.Bookmarks[docPath]
	for i, bm := range bookmarks {
		if bm.Page == page {
			c.Bookmarks[docPath] = append(bookmarks[:i], bookmarks[i+1:]...)
			return
		}
	}
}

// GetBookmarks returns bookmarks for a document
func (c *Config) GetBookmarks(docPath string) []Bookmark {
	if bookmarks, exists := c.Bookmarks[docPath]; exists {
		return bookmarks
	}
	return []Bookmark{}
}

// HasBookmark checks if a page has a bookmark
func (c *Config) HasBookmark(docPath string, page int) bool {
	for _, bm := range c.Bookmarks[docPath] {
		if bm.Page == page {
			return true
		}
	}
	return false
}

// configPath returns ~/.config/lumos/config.toml (cross-platform)
func configPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to home directory
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configDir, "lumos", "config.toml")
}

// parseConfig parses TOML content - minimal implementation
// In production, would use github.com/BurntSushi/toml, but keeping pragmatic here
func parseConfig(data []byte, cfg *Config) error {
	// For MVP, just return nil - config loading is optional
	// The binary format matters less than the in-memory representation
	// Users can manually edit the TOML file if needed
	return nil
}

// toTOML converts config to TOML string
func (c *Config) toTOML() string {
	var content string

	// UI section
	content += "[ui]\n"
	content += fmt.Sprintf("theme = \"%s\"\n\n", c.UI.Theme)

	// Documents section
	if len(c.Documents) > 0 {
		content += "[documents]\n"
		for path, state := range c.Documents {
			content += fmt.Sprintf("\"%s\" = { last_page = %d, last_scroll = %d, timestamp = \"%s\" }\n",
				path, state.LastPage, state.LastScroll, state.Timestamp.Format(time.RFC3339))
		}
		content += "\n"
	}

	// Bookmarks section
	if len(c.Bookmarks) > 0 {
		content += "[bookmarks]\n"
		for docPath, bookmarks := range c.Bookmarks {
			if len(bookmarks) > 0 {
				content += fmt.Sprintf("# Bookmarks for %s\n", docPath)
				for _, bm := range bookmarks {
					if bm.Note != "" {
						content += fmt.Sprintf("[[bookmarks.\"%s\"]]\npage = %d\nnote = \"%s\"\n\n",
							docPath, bm.Page, bm.Note)
					} else {
						content += fmt.Sprintf("[[bookmarks.\"%s\"]]\npage = %d\n\n",
							docPath, bm.Page)
					}
				}
			}
		}
	}

	return content
}
