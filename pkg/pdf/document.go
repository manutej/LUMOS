package pdf

import (
	"fmt"
	"sync"

	"github.com/ledongthuc/pdf"
)

// Document represents a PDF document with cached pages
type Document struct {
	filepath string
	reader   *pdf.Reader
	pages    int

	// Caching
	pageCache map[int]string
	cacheMu   sync.RWMutex
	maxCache  int // LRU cache size

	// Metadata
	title      string
	author     string
	subject    string
	creator    string
	totalWords int
}

// PageInfo contains extracted text and metadata for a page
type PageInfo struct {
	PageNum    int
	Text       string
	LineCount  int
	WordCount  int
	HasImages  bool
	HasTables  bool
}

// NewDocument creates a new PDF document from a file path
func NewDocument(filepath string, maxCachePages int) (*Document, error) {
	f, r, err := pdf.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	pages := r.NumPage()
	if pages == 0 {
		return nil, fmt.Errorf("PDF has no pages")
	}

	doc := &Document{
		filepath:  filepath,
		pages:     pages,
		pageCache: make(map[int]string),
		maxCache:  maxCachePages,
	}

	// Extract metadata
	// Note: Metadata extraction depends on PDF structure
	// We'll implement basic metadata extraction later
	doc.title = filepath
	doc.author = ""
	doc.subject = ""
	doc.creator = ""

	return doc, nil
}

// GetPageCount returns the total number of pages
func (d *Document) GetPageCount() int {
	return d.pages
}

// GetPage retrieves text content from a specific page
func (d *Document) GetPage(pageNum int) (*PageInfo, error) {
	if pageNum < 1 || pageNum > d.pages {
		return nil, fmt.Errorf("page number out of range: %d", pageNum)
	}

	// Check cache first
	d.cacheMu.RLock()
	if cached, exists := d.pageCache[pageNum]; exists {
		d.cacheMu.RUnlock()
		return d.createPageInfo(pageNum, cached), nil
	}
	d.cacheMu.RUnlock()

	// Extract text from page
	f, r, err := pdf.Open(d.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to reopen PDF: %w", err)
	}
	defer f.Close()

	page := r.Page(pageNum)
	if page.V.IsNull() {
		return nil, fmt.Errorf("page %d is empty or null", pageNum)
	}

	// Extract plain text from page
	// NOTE: PDFs often split text into small chunks (sometimes single characters)
	// We concatenate without adding spaces and let natural spacing in the PDF show through
	content := ""
	texts := page.Content().Text
	for _, text := range texts {
		content += text.S
	}

	// Cache the result
	d.cacheMu.Lock()
	d.pageCache[pageNum] = content
	d.cacheMu.Unlock()

	return d.createPageInfo(pageNum, content), nil
}

// GetPageRange retrieves text from a range of pages
func (d *Document) GetPageRange(startPage, endPage int) (string, error) {
	if startPage < 1 || endPage > d.pages || startPage > endPage {
		return "", fmt.Errorf("invalid page range: %d-%d", startPage, endPage)
	}

	var content string
	for pageNum := startPage; pageNum <= endPage; pageNum++ {
		page, err := d.GetPage(pageNum)
		if err != nil {
			return "", err
		}
		content += page.Text + "\n"
	}

	return content, nil
}

// Search finds all occurrences of a search term
func (d *Document) Search(query string) ([]SearchResult, error) {
	var results []SearchResult

	for pageNum := 1; pageNum <= d.pages; pageNum++ {
		page, err := d.GetPage(pageNum)
		if err != nil {
			continue
		}

		matches := findMatches(page.Text, query)
		for _, match := range matches {
			results = append(results, SearchResult{
				PageNum:      pageNum,
				LineNum:      match.LineNum,
				ColumnNum:    match.ColumnNum,
				MatchText:    match.Text,
				ContextBefore: match.ContextBefore,
				ContextAfter:  match.ContextAfter,
			})
		}
	}

	return results, nil
}

// GetMetadata returns document metadata
func (d *Document) GetMetadata() Metadata {
	return Metadata{
		FilePath: d.filepath,
		Pages:    d.pages,
		Title:    d.title,
		Author:   d.author,
		Subject:  d.subject,
		Creator:  d.creator,
	}
}

// ClearCache clears all cached pages
func (d *Document) ClearCache() {
	d.cacheMu.Lock()
	defer d.cacheMu.Unlock()
	d.pageCache = make(map[int]string)
}

// CacheStats returns cache statistics
func (d *Document) CacheStats() CacheStats {
	d.cacheMu.RLock()
	defer d.cacheMu.RUnlock()

	return CacheStats{
		CachedPages: len(d.pageCache),
		MaxSize:     d.maxCache,
	}
}

// Helper functions

func (d *Document) createPageInfo(pageNum int, text string) *PageInfo {
	lineCount := countLines(text)
	wordCount := countWords(text)

	return &PageInfo{
		PageNum:   pageNum,
		Text:      text,
		LineCount: lineCount,
		WordCount: wordCount,
		// TODO: Detect images and tables
		HasImages: false,
		HasTables: false,
	}
}

func countLines(text string) int {
	count := 0
	for _, char := range text {
		if char == '\n' {
			count++
		}
	}
	return count + 1 // +1 for first line
}

func countWords(text string) int {
	count := 0
	inWord := false

	for _, char := range text {
		if char == ' ' || char == '\n' || char == '\t' {
			if inWord {
				count++
				inWord = false
			}
		} else {
			inWord = true
		}
	}

	if inWord {
		count++
	}

	return count
}

// SearchResult represents a single search result
type SearchResult struct {
	PageNum        int
	LineNum        int
	ColumnNum      int
	MatchText      string
	ContextBefore  string
	ContextAfter   string
}

// Metadata contains document metadata
type Metadata struct {
	FilePath string
	Pages    int
	Title    string
	Author   string
	Subject  string
	Creator  string
}

// CacheStats contains cache statistics
type CacheStats struct {
	CachedPages int
	MaxSize     int
}

// Match represents a text match within a page
type Match struct {
	Text          string
	LineNum       int
	ColumnNum     int
	ContextBefore string
	ContextAfter  string
}

// findMatches finds all case-insensitive matches of a query in text
func findMatches(text, query string) []Match {
	if query == "" {
		return []Match{}
	}

	var matches []Match
	lines := TextToLines(text)

	for _, line := range lines {
		// Case-insensitive search
		positions := CaseInsensitiveMatch(line.Text, query)

		for _, pos := range positions {
			before, match, after := ExtractContext(line.Text, pos, len(query), 30)

			matches = append(matches, Match{
				Text:          match,
				LineNum:       line.LineNum,
				ColumnNum:     pos + 1, // 1-indexed
				ContextBefore: before,
				ContextAfter:  after,
			})
		}
	}

	return matches
}
