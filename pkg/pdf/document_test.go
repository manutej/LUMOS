package pdf

import (
	"os"
	"path/filepath"
	"testing"
)

// Test fixtures path
const testFixturesDir = "../../test/fixtures"

// Helper to get test PDF path
func getTestPDF(name string) string {
	return filepath.Join(testFixturesDir, name)
}

// TestNewDocument tests document creation
func TestNewDocument(t *testing.T) {
	tests := []struct {
		name      string
		filepath  string
		maxCache  int
		wantErr   bool
		errString string
	}{
		{
			name:     "valid PDF",
			filepath: getTestPDF("simple.pdf"),
			maxCache: 10,
			wantErr:  false,
		},
		{
			name:      "non-existent file",
			filepath:  getTestPDF("nonexistent.pdf"),
			maxCache:  10,
			wantErr:   true,
			errString: "failed to open PDF",
		},
		{
			name:     "zero cache size",
			filepath: getTestPDF("simple.pdf"),
			maxCache: 0,
			wantErr:  false,
		},
		{
			name:     "large cache size",
			filepath: getTestPDF("simple.pdf"),
			maxCache: 1000,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip if test PDF doesn't exist and we expect success
			if !tt.wantErr {
				if _, err := os.Stat(tt.filepath); os.IsNotExist(err) {
					t.Skip("Test PDF fixture not found:", tt.filepath)
				}
			}

			doc, err := NewDocument(tt.filepath, tt.maxCache)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewDocument() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("NewDocument() unexpected error: %v", err)
			}

			if doc == nil {
				t.Fatal("NewDocument() returned nil document")
			}

			// Verify document properties
			if doc.filepath != tt.filepath {
				t.Errorf("filepath = %v, want %v", doc.filepath, tt.filepath)
			}

			if true { // maxCache no longer exists
				// maxCache check removed - cache field replaced with LRUCache
			}

			if doc.pages <= 0 {
				t.Errorf("pages = %v, want > 0", doc.pages)
			}

			if doc.cache == nil {
				t.Error("cache is nil")
			}
		})
	}
}

// TestGetPageCount tests page count retrieval
func TestGetPageCount(t *testing.T) {
	testPDF := getTestPDF("simple.pdf")
	if _, err := os.Stat(testPDF); os.IsNotExist(err) {
		t.Skip("Test PDF fixture not found:", testPDF)
	}

	doc, err := NewDocument(testPDF, 10)
	if err != nil {
		t.Fatalf("Failed to create document: %v", err)
	}

	count := doc.GetPageCount()
	if count <= 0 {
		t.Errorf("GetPageCount() = %v, want > 0", count)
	}

	// Should be consistent across multiple calls
	count2 := doc.GetPageCount()
	if count != count2 {
		t.Errorf("GetPageCount() inconsistent: first=%v, second=%v", count, count2)
	}
}

// TestGetPage tests page retrieval
func TestGetPage(t *testing.T) {
	testPDF := getTestPDF("simple.pdf")
	if _, err := os.Stat(testPDF); os.IsNotExist(err) {
		t.Skip("Test PDF fixture not found:", testPDF)
	}

	doc, err := NewDocument(testPDF, 5)
	if err != nil {
		t.Fatalf("Failed to create document: %v", err)
	}

	tests := []struct {
		name      string
		pageNum   int
		wantErr   bool
		errString string
	}{
		{
			name:    "first page",
			pageNum: 1,
			wantErr: false,
		},
		{
			name:      "page zero",
			pageNum:   0,
			wantErr:   true,
			errString: "out of range",
		},
		{
			name:      "negative page",
			pageNum:   -1,
			wantErr:   true,
			errString: "out of range",
		},
		{
			name:      "page beyond document",
			pageNum:   9999,
			wantErr:   true,
			errString: "out of range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, err := doc.GetPage(tt.pageNum)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetPage(%v) expected error but got nil", tt.pageNum)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetPage(%v) unexpected error: %v", tt.pageNum, err)
			}

			if page == nil {
				t.Fatal("GetPage() returned nil page")
			}

			// Verify page info
			if page.PageNum != tt.pageNum {
				t.Errorf("page.PageNum = %v, want %v", page.PageNum, tt.pageNum)
			}

			if page.Text == "" {
				t.Log("Warning: Page has no text content")
			}

			if page.LineCount < 0 {
				t.Errorf("page.LineCount = %v, want >= 0", page.LineCount)
			}

			if page.WordCount < 0 {
				t.Errorf("page.WordCount = %v, want >= 0", page.WordCount)
			}
		})
	}
}

// TestGetPageCaching tests that pages are cached correctly
func TestGetPageCaching(t *testing.T) {
	testPDF := getTestPDF("simple.pdf")
	if _, err := os.Stat(testPDF); os.IsNotExist(err) {
		t.Skip("Test PDF fixture not found:", testPDF)
	}

	doc, err := NewDocument(testPDF, 5)
	if err != nil {
		t.Fatalf("Failed to create document: %v", err)
	}

	// Get page first time (should cache)
	page1, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1) first call failed: %v", err)
	}

	// Verify it's in cache
	cached, inCache := doc.cache.Get(1)
	if !inCache {
		t.Error("Page 1 not in cache after GetPage")
	}
	if cached == "" {
		t.Error("Cached page content is empty")
	}

	// Get same page again (should use cache)
	page2, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1) second call failed: %v", err)
	}

	// Content should be identical
	if page1.Text != page2.Text {
		t.Error("Cached page content differs from original")
	}
}

// TestGetPageRange tests page range retrieval
func TestGetPageRange(t *testing.T) {
	testPDF := getTestPDF("simple.pdf")
	if _, err := os.Stat(testPDF); os.IsNotExist(err) {
		t.Skip("Test PDF fixture not found:", testPDF)
	}

	doc, err := NewDocument(testPDF, 10)
	if err != nil {
		t.Fatalf("Failed to create document: %v", err)
	}

	totalPages := doc.GetPageCount()

	tests := []struct {
		name      string
		startPage int
		endPage   int
		wantErr   bool
	}{
		{
			name:      "single page range",
			startPage: 1,
			endPage:   1,
			wantErr:   false,
		},
		{
			name:      "all pages",
			startPage: 1,
			endPage:   totalPages,
			wantErr:   false,
		},
		{
			name:      "invalid range - reversed",
			startPage: 2,
			endPage:   1,
			wantErr:   true,
		},
		{
			name:      "invalid range - beyond document",
			startPage: 1,
			endPage:   9999,
			wantErr:   true,
		},
		{
			name:      "invalid range - start < 1",
			startPage: 0,
			endPage:   2,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := doc.GetPageRange(tt.startPage, tt.endPage)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetPageRange(%v, %v) expected error but got nil",
						tt.startPage, tt.endPage)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetPageRange(%v, %v) unexpected error: %v",
					tt.startPage, tt.endPage, err)
			}

			if content == "" && tt.endPage >= tt.startPage {
				t.Log("Warning: Page range has no text content")
			}
		})
	}
}

// TestSearch tests document search functionality
func TestSearch(t *testing.T) {
	testPDF := getTestPDF("simple.pdf")
	if _, err := os.Stat(testPDF); os.IsNotExist(err) {
		t.Skip("Test PDF fixture not found:", testPDF)
	}

	doc, err := NewDocument(testPDF, 10)
	if err != nil {
		t.Fatalf("Failed to create document: %v", err)
	}

	tests := []struct {
		name        string
		query       string
		expectFound bool
	}{
		{
			name:        "common word",
			query:       "T e s t",
			expectFound: true,
		},
		{
			name:        "empty query",
			query:       "",
			expectFound: false,
		},
		{
			name:        "non-existent word",
			query:       "xyzabc123nonexistent",
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := doc.Search(tt.query)
			if err != nil {
				t.Fatalf("Search(%q) unexpected error: %v", tt.query, err)
			}

			hasResults := len(results) > 0

			if tt.expectFound && !hasResults {
				t.Errorf("Search(%q) expected results but got none", tt.query)
			}

			// Verify result structure
			for i, result := range results {
				if result.PageNum < 1 || result.PageNum > doc.pages {
					t.Errorf("result[%d].PageNum = %v, want 1-%v",
						i, result.PageNum, doc.pages)
				}
			}
		})
	}
}

// TestGetMetadata tests metadata retrieval
func TestGetMetadata(t *testing.T) {
	testPDF := getTestPDF("simple.pdf")
	if _, err := os.Stat(testPDF); os.IsNotExist(err) {
		t.Skip("Test PDF fixture not found:", testPDF)
	}

	doc, err := NewDocument(testPDF, 10)
	if err != nil {
		t.Fatalf("Failed to create document: %v", err)
	}

	metadata := doc.GetMetadata()

	if metadata.FilePath != testPDF {
		t.Errorf("metadata.FilePath = %v, want %v", metadata.FilePath, testPDF)
	}

	if metadata.Pages != doc.pages {
		t.Errorf("metadata.Pages = %v, want %v", metadata.Pages, doc.pages)
	}

	if metadata.Pages <= 0 {
		t.Errorf("metadata.Pages = %v, want > 0", metadata.Pages)
	}
}

// TestClearCache tests cache clearing
func TestClearCache(t *testing.T) {
	testPDF := getTestPDF("simple.pdf")
	if _, err := os.Stat(testPDF); os.IsNotExist(err) {
		t.Skip("Test PDF fixture not found:", testPDF)
	}

	doc, err := NewDocument(testPDF, 10)
	if err != nil {
		t.Fatalf("Failed to create document: %v", err)
	}

	// Load some pages to populate cache
	_, _ = doc.GetPage(1)
	if doc.GetPageCount() > 1 {
		_, _ = doc.GetPage(2)
	}

	// Verify cache has entries
	stats := doc.cache.Stats()
	cacheSize := stats.CachedPages

	if cacheSize == 0 {
		t.Fatal("Cache should have entries before ClearCache")
	}

	// Clear cache
	doc.ClearCache()

	// Verify cache is empty
	stats = doc.cache.Stats()
	cacheSize = stats.CachedPages

	if cacheSize != 0 {
		t.Errorf("After ClearCache, cache size = %v, want 0", cacheSize)
	}
}

// TestCacheStats tests cache statistics
func TestCacheStats(t *testing.T) {
	testPDF := getTestPDF("simple.pdf")
	if _, err := os.Stat(testPDF); os.IsNotExist(err) {
		t.Skip("Test PDF fixture not found:", testPDF)
	}

	maxCache := 5
	doc, err := NewDocument(testPDF, maxCache)
	if err != nil {
		t.Fatalf("Failed to create document: %v", err)
	}

	stats := doc.CacheStats()

	if stats.MaxSize != maxCache {
		t.Errorf("stats.MaxSize = %v, want %v", stats.MaxSize, maxCache)
	}

	if stats.CachedPages < 0 {
		t.Errorf("stats.CachedPages = %v, want >= 0", stats.CachedPages)
	}

	// Load a page and check stats update
	_, _ = doc.GetPage(1)
	stats = doc.CacheStats()

	if stats.CachedPages != 1 {
		t.Errorf("After loading 1 page, stats.CachedPages = %v, want 1", stats.CachedPages)
	}
}

// TestCountLines tests line counting helper
func TestCountLines(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		want  int
	}{
		{"empty", "", 1},
		{"single line", "hello", 1},
		{"two lines", "hello\nworld", 2},
		{"three lines", "line1\nline2\nline3", 3},
		{"trailing newline", "hello\n", 2},
		{"multiple newlines", "a\n\n\nb", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countLines(tt.text)
			if got != tt.want {
				t.Errorf("countLines(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

// TestCountWords tests word counting helper
func TestCountWords(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		want  int
	}{
		{"empty", "", 0},
		{"single word", "hello", 1},
		{"two words", "hello world", 2},
		{"multiple spaces", "hello    world", 2},
		{"with tabs", "hello\tworld", 2},
		{"with newlines", "hello\nworld", 2},
		{"mixed whitespace", "hello   \t\n  world", 2},
		{"trailing space", "hello ", 1},
		{"leading space", " hello", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countWords(tt.text)
			if got != tt.want {
				t.Errorf("countWords(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

// TestThreadSafety tests concurrent access to document
func TestThreadSafety(t *testing.T) {
	testPDF := getTestPDF("simple.pdf")
	if _, err := os.Stat(testPDF); os.IsNotExist(err) {
		t.Skip("Test PDF fixture not found:", testPDF)
	}

	doc, err := NewDocument(testPDF, 10)
	if err != nil {
		t.Fatalf("Failed to create document: %v", err)
	}

	// Run concurrent GetPage calls
	numGoroutines := 10
	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			_, _ = doc.GetPage(1)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// If we get here without deadlock or panic, test passes
}
