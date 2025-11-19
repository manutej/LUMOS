package pdf

import (
	"fmt"
	"regexp"
	"strings"
)

// SearchOptions contains advanced search configuration
type SearchOptions struct {
	CaseSensitive bool   // Perform case-sensitive matching
	WholeWord     bool   // Match only whole words
	RegexMode     bool   // Treat query as regular expression
	MaxResults    int    // Maximum number of results (0 = unlimited)
	StartPage     int    // Start search from this page (1-indexed)
	EndPage       int    // End search at this page (0 = to end)
}

// SearchFilter filters results based on criteria
type SearchFilter struct {
	MinContextLength int // Minimum characters before/after match
	HighlightMatch   bool // Include highlighting markers
	SortByPage       bool // Sort results by page number
}

// SearchResultAdvanced extends SearchResult with additional metadata
type SearchResultAdvanced struct {
	SearchResult
	MatchCount   int    // Number of matches on this page
	Relevance    float64 // Relevance score (0-1)
	PreviewText  string  // Preview of context around match
}

// AdvancedSearch performs complex search with options
func (d *Document) AdvancedSearch(query string, opts SearchOptions) ([]SearchResultAdvanced, error) {
	if query == "" {
		return []SearchResultAdvanced{}, nil
	}

	// Validate options
	if opts.MaxResults < 0 {
		opts.MaxResults = 0
	}
	if opts.StartPage < 1 {
		opts.StartPage = 1
	}
	if opts.EndPage == 0 || opts.EndPage > d.pages {
		opts.EndPage = d.pages
	}

	// Compile regex if needed
	var regex *regexp.Regexp
	if opts.RegexMode {
		var err error
		flags := ""
		if !opts.CaseSensitive {
			flags = "(?i)"
		}
		regex, err = regexp.Compile(flags + query)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern: %w", err)
		}
	}

	var results []SearchResultAdvanced
	resultCount := 0

	for pageNum := opts.StartPage; pageNum <= opts.EndPage; pageNum++ {
		if opts.MaxResults > 0 && resultCount >= opts.MaxResults {
			break
		}

		page, err := d.GetPage(pageNum)
		if err != nil {
			continue
		}

		pageMatches := findAdvancedMatches(page.Text, query, opts, regex)
		for _, match := range pageMatches {
			if opts.MaxResults > 0 && resultCount >= opts.MaxResults {
				break
			}

			result := SearchResultAdvanced{
				SearchResult: SearchResult{
					PageNum:       pageNum,
					LineNum:       match.LineNum,
					ColumnNum:     match.ColumnNum,
					MatchText:     match.Text,
					ContextBefore: match.ContextBefore,
					ContextAfter:  match.ContextAfter,
				},
				MatchCount:  len(pageMatches),
				Relevance:   calculateRelevance(match.Text, query),
				PreviewText: formatPreview(match.ContextBefore, match.Text, match.ContextAfter),
			}
			results = append(results, result)
			resultCount++
		}
	}

	return results, nil
}

// findAdvancedMatches finds matches in text with advanced options
func findAdvancedMatches(text, query string, opts SearchOptions, regex *regexp.Regexp) []Match {
	var matches []Match

	lines := TextToLines(text)
	searchText := text
	if !opts.CaseSensitive && !opts.RegexMode {
		searchText = strings.ToLower(text)
		query = strings.ToLower(query)
	}

	if opts.RegexMode && regex != nil {
		// Regex matching
		allMatches := regex.FindAllStringIndex(searchText, -1)
		for _, matchIdx := range allMatches {
			start, end := matchIdx[0], matchIdx[1]
			matchText := text[start:end]

			// Get line and column info
			lineNum, colNum := getLineAndColumn(text, start)

			match := Match{
				Text:           matchText,
				LineNum:        lineNum,
				ColumnNum:      colNum,
				ContextBefore:  extractContextAdvanced(text, start, -40),
				ContextAfter:   extractContextAdvanced(text, end, 40),
			}
			matches = append(matches, match)
		}
	} else {
		// Standard matching with whole-word option
		for i := 0; i < len(lines); i++ {
			line := lines[i]
			lineSearchText := line.Text
			if !opts.CaseSensitive {
				lineSearchText = strings.ToLower(line.Text)
			}

			positions := findAllPositions(lineSearchText, query, opts.WholeWord)
			for _, pos := range positions {
				match := Match{
					Text:           query,
					LineNum:        i + 1,
					ColumnNum:      pos,
					ContextBefore:  extractLineContext(line.Text, pos, -30),
					ContextAfter:   extractLineContext(line.Text, pos+len(query), 30),
				}
				matches = append(matches, match)
			}
		}
	}

	return matches
}

// findAllPositions finds all positions of a string in text
func findAllPositions(text, substr string, wholeWord bool) []int {
	var positions []int
	start := 0

	for {
		pos := strings.Index(text[start:], substr)
		if pos == -1 {
			break
		}

		absolutePos := start + pos

		// Check whole word boundary if needed
		if wholeWord {
			if !isWordBoundaryPosition(text, absolutePos) ||
				!isWordBoundaryPosition(text, absolutePos+len(substr)) {
				start = absolutePos + 1
				continue
			}
		}

		positions = append(positions, absolutePos)
		start = absolutePos + 1
	}

	return positions
}

// isWordBoundaryPosition checks if position is at a word boundary
func isWordBoundaryPosition(text string, pos int) bool {
	if pos < 0 || pos > len(text) {
		return false
	}

	// Start of text or after whitespace/punctuation
	if pos == 0 {
		return true
	}
	prevChar := rune(text[pos-1])
	if !isWordCharRune(prevChar) {
		return true
	}

	// End of text or before whitespace/punctuation
	if pos == len(text) {
		return true
	}
	nextChar := rune(text[pos])
	if !isWordCharRune(nextChar) {
		return true
	}

	return false
}

// isWordCharRune checks if a rune is part of a word
func isWordCharRune(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '_'
}

// getLineAndColumn returns line and column numbers for a position in text
func getLineAndColumn(text string, pos int) (int, int) {
	line := 1
	col := 1
	for i := 0; i < pos && i < len(text); i++ {
		if text[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return line, col
}

// extractLineContext extracts context from a single line
func extractLineContext(line string, pos, length int) string {
	if length > 0 {
		// Extract after position
		end := pos + length
		if end > len(line) {
			end = len(line)
		}
		return strings.TrimSpace(line[pos:end])
	} else {
		// Extract before position
		start := pos + length // length is negative
		if start < 0 {
			start = 0
		}
		return strings.TrimSpace(line[start:pos])
	}
}

// extractContextAdvanced extracts context around a position in text
func extractContextAdvanced(text string, pos int, length int) string {
	if length > 0 {
		// Extract after position
		end := pos + length
		if end > len(text) {
			end = len(text)
		}
		return strings.TrimSpace(text[pos:end])
	} else {
		// Extract before position
		start := pos + length // length is negative
		if start < 0 {
			start = 0
		}
		return strings.TrimSpace(text[start:pos])
	}
}

// calculateRelevance calculates relevance score (0-1)
func calculateRelevance(matchText, query string) float64 {
	// Exact match = 1.0
	if matchText == query {
		return 1.0
	}

	// Partial match scored by length ratio
	if len(matchText) == 0 {
		return 0.0
	}

	ratio := float64(len(query)) / float64(len(matchText))
	if ratio > 1.0 {
		ratio = 1.0 / ratio
	}
	return ratio
}

// formatPreview formats a search result preview
func formatPreview(before, match, after string) string {
	const maxLen = 60
	preview := fmt.Sprintf("...%s [%s] %s...", before, match, after)
	if len(preview) > maxLen {
		preview = preview[:maxLen] + "..."
	}
	return preview
}

// SearchHistoryManager tracks search history
type SearchHistoryManager struct {
	history []SearchHistoryEntry
	maxSize int
}

// SearchHistoryEntry represents a single search in history
type SearchHistoryEntry struct {
	Query     string
	Options   SearchOptions
	ResultCount int
	Timestamp int64 // Unix timestamp
}

// NewSearchHistoryManager creates a new history manager
func NewSearchHistoryManager(maxSize int) *SearchHistoryManager {
	return &SearchHistoryManager{
		history: []SearchHistoryEntry{},
		maxSize: maxSize,
	}
}

// Add adds a search to history
func (shm *SearchHistoryManager) Add(entry SearchHistoryEntry) {
	shm.history = append([]SearchHistoryEntry{entry}, shm.history...)
	if len(shm.history) > shm.maxSize {
		shm.history = shm.history[:shm.maxSize]
	}
}

// GetHistory returns search history
func (shm *SearchHistoryManager) GetHistory() []SearchHistoryEntry {
	return shm.history
}

// Clear clears search history
func (shm *SearchHistoryManager) Clear() {
	shm.history = []SearchHistoryEntry{}
}

// Size returns number of items in history
func (shm *SearchHistoryManager) Size() int {
	return len(shm.history)
}
