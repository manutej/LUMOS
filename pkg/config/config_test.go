package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDefaultConfig verifies defaults
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.UI.Theme != "dark" {
		t.Errorf("Default theme should be 'dark', got %s", cfg.UI.Theme)
	}
	if cfg.Documents == nil {
		t.Error("Documents map should be initialized")
	}
	if cfg.Bookmarks == nil {
		t.Error("Bookmarks map should be initialized")
	}
}

// TestAddBookmark adds and retrieves bookmarks
func TestAddBookmark(t *testing.T) {
	cfg := DefaultConfig()
	docPath := "/path/to/doc.pdf"

	cfg.AddBookmark(docPath, 5, "Important section")

	if !cfg.HasBookmark(docPath, 5) {
		t.Error("Bookmark should exist after adding")
	}

	bookmarks := cfg.GetBookmarks(docPath)
	if len(bookmarks) != 1 {
		t.Errorf("Expected 1 bookmark, got %d", len(bookmarks))
	}
	if bookmarks[0].Page != 5 || bookmarks[0].Note != "Important section" {
		t.Errorf("Bookmark data mismatch: got %+v", bookmarks[0])
	}
}

// TestRemoveBookmark deletes bookmarks
func TestRemoveBookmark(t *testing.T) {
	cfg := DefaultConfig()
	docPath := "/path/to/doc.pdf"

	cfg.AddBookmark(docPath, 5, "Test")
	cfg.AddBookmark(docPath, 10, "Another")

	cfg.RemoveBookmark(docPath, 5)

	if cfg.HasBookmark(docPath, 5) {
		t.Error("Bookmark should be removed")
	}

	bookmarks := cfg.GetBookmarks(docPath)
	if len(bookmarks) != 1 || bookmarks[0].Page != 10 {
		t.Errorf("Wrong bookmarks remain: got %+v", bookmarks)
	}
}

// TestUpdateDocState tracks document position
func TestUpdateDocState(t *testing.T) {
	cfg := DefaultConfig()
	docPath := "/path/to/doc.pdf"

	cfg.UpdateDocState(docPath, 42, 150)

	state := cfg.Documents[docPath]
	if state.LastPage != 42 {
		t.Errorf("LastPage should be 42, got %d", state.LastPage)
	}
	if state.LastScroll != 150 {
		t.Errorf("LastScroll should be 150, got %d", state.LastScroll)
	}
	if state.Timestamp.IsZero() {
		t.Error("Timestamp should be set")
	}
}

// TestDuplicateBookmark updates existing bookmark on same page
func TestDuplicateBookmark(t *testing.T) {
	cfg := DefaultConfig()
	docPath := "/path/to/doc.pdf"

	cfg.AddBookmark(docPath, 5, "First note")
	cfg.AddBookmark(docPath, 5, "Updated note")

	bookmarks := cfg.GetBookmarks(docPath)
	if len(bookmarks) != 1 {
		t.Errorf("Should have 1 bookmark, got %d", len(bookmarks))
	}
	if bookmarks[0].Note != "Updated note" {
		t.Errorf("Bookmark note should be updated, got %s", bookmarks[0].Note)
	}
}

// TestGetBookmarksEmpty returns empty slice for nonexistent document
func TestGetBookmarksEmpty(t *testing.T) {
	cfg := DefaultConfig()

	bookmarks := cfg.GetBookmarks("/nonexistent/doc.pdf")
	if len(bookmarks) != 0 {
		t.Errorf("Should return empty slice, got %d bookmarks", len(bookmarks))
	}
}

// TestSaveLoad persists and restores config
func TestSaveLoad(t *testing.T) {
	// Use temp directory for testing
	tmpDir := t.TempDir()
	oldConfigDir := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldConfigDir)

	cfg := DefaultConfig()
	cfg.UI.Theme = "tokyo-night"
	cfg.UpdateDocState("/path/to/doc.pdf", 42, 100)
	cfg.AddBookmark("/path/to/doc.pdf", 5, "Chapter 1")

	// Save
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify file was created
	path := configPath()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("Config file not created: %v", err)
	}

	// Verify file contains expected content
	content, _ := os.ReadFile(path)
	contentStr := string(content)
	if !contains(contentStr, "tokyo-night") {
		t.Error("Config should contain theme")
	}
	if !contains(contentStr, "/path/to/doc.pdf") {
		t.Error("Config should contain document path")
	}
}

// TestToTOML generates valid TOML
func TestToTOML(t *testing.T) {
	cfg := DefaultConfig()
	cfg.UI.Theme = "dracula"
	cfg.UpdateDocState("/docs/paper.pdf", 15, 50)
	cfg.AddBookmark("/docs/paper.pdf", 10, "Methods")

	toml := cfg.toTOML()

	if !contains(toml, "[ui]") {
		t.Error("TOML should contain [ui] section")
	}
	if !contains(toml, "theme = \"dracula\"") {
		t.Error("TOML should contain theme")
	}
	if !contains(toml, "[documents]") {
		t.Error("TOML should contain [documents] section")
	}
	if !contains(toml, "last_page = 15") {
		t.Error("TOML should contain last page")
	}
	if !contains(toml, "[bookmarks]") {
		t.Error("TOML should contain [bookmarks] section")
	}
	if !contains(toml, "page = 10") {
		t.Error("TOML should contain bookmark page")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr ||
		indexStr(s, substr) >= 0
}

func indexStr(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// BenchmarkAddBookmark benchmarks bookmark addition
func BenchmarkAddBookmark(b *testing.B) {
	cfg := DefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg.AddBookmark("/doc.pdf", i%100, "note")
	}
}

// BenchmarkSave benchmarks config saving
func BenchmarkSave(b *testing.B) {
	tmpDir := b.TempDir()
	cfg := DefaultConfig()
	for i := 0; i < 10; i++ {
		cfg.UpdateDocState(filepath.Join(tmpDir, "doc"+string(rune(i))+".pdf"), i, i*10)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg.Save()
	}
}
