package pdf

import (
	"strings"
	"unicode"
)

// TextSearch provides full-text search functionality
type TextSearch struct {
	query   string
	results []SearchResult
	current int
}

// NewTextSearch creates a new text search
func NewTextSearch(query string) *TextSearch {
	return &TextSearch{
		query:   query,
		current: 0,
	}
}

// Execute performs the search across a document
func (ts *TextSearch) Execute(doc *Document) error {
	results, err := doc.Search(ts.query)
	if err != nil {
		return err
	}
	ts.results = results
	ts.current = 0
	return nil
}

// GetResults returns all search results
func (ts *TextSearch) GetResults() []SearchResult {
	return ts.results
}

// GetResultCount returns the number of results
func (ts *TextSearch) GetResultCount() int {
	return len(ts.results)
}

// GetCurrentResult returns the current search result
func (ts *TextSearch) GetCurrentResult() *SearchResult {
	if ts.current >= 0 && ts.current < len(ts.results) {
		return &ts.results[ts.current]
	}
	return nil
}

// NextResult moves to the next result
func (ts *TextSearch) NextResult() bool {
	if ts.current+1 < len(ts.results) {
		ts.current++
		return true
	}
	return false
}

// PreviousResult moves to the previous result
func (ts *TextSearch) PreviousResult() bool {
	if ts.current-1 >= 0 {
		ts.current--
		return true
	}
	return false
}

// JumpToResult jumps to a specific result index
func (ts *TextSearch) JumpToResult(index int) bool {
	if index >= 0 && index < len(ts.results) {
		ts.current = index
		return true
	}
	return false
}

// GetCurrentIndex returns the current result index
func (ts *TextSearch) GetCurrentIndex() int {
	return ts.current
}

// Reset resets the search
func (ts *TextSearch) Reset() {
	ts.current = 0
}

// CaseSensitiveMatch performs case-sensitive text matching
func CaseSensitiveMatch(text, query string) []int {
	var matches []int
	start := 0

	for {
		idx := strings.Index(text[start:], query)
		if idx == -1 {
			break
		}
		matches = append(matches, start+idx)
		start = start + idx + 1
	}

	return matches
}

// CaseInsensitiveMatch performs case-insensitive text matching
func CaseInsensitiveMatch(text, query string) []int {
	return CaseSensitiveMatch(strings.ToLower(text), strings.ToLower(query))
}

// WordMatch performs whole-word matching
func WordMatch(text, query string) []int {
	text = strings.ToLower(text)
	query = strings.ToLower(query)
	var matches []int

	start := 0
	for {
		idx := strings.Index(text[start:], query)
		if idx == -1 {
			break
		}

		globalIdx := start + idx

		// Check if it's a whole word
		if isWordBoundary(text, globalIdx) && isWordBoundary(text, globalIdx+len(query)) {
			matches = append(matches, globalIdx)
		}

		start = globalIdx + 1
	}

	return matches
}

// isWordBoundary checks if a position is at a word boundary
func isWordBoundary(text string, pos int) bool {
	if pos <= 0 || pos >= len(text) {
		return pos == 0 || pos == len(text)
	}

	before := text[pos-1]
	after := text[pos]

	beforeIsLetter := unicode.IsLetter(rune(before))
	afterIsLetter := unicode.IsLetter(rune(after))

	return beforeIsLetter != afterIsLetter
}

// ExtractContext extracts context around a match
func ExtractContext(text string, matchPos, matchLen, contextChars int) (before, match, after string) {
	start := matchPos - contextChars
	if start < 0 {
		start = 0
	}

	end := matchPos + matchLen + contextChars
	if end > len(text) {
		end = len(text)
	}

	before = text[start:matchPos]
	match = text[matchPos : matchPos+matchLen]
	after = text[matchPos+matchLen : end]

	return
}

// HighlightMatches returns text with matches highlighted (for terminal output)
func HighlightMatches(text string, matches []int, queryLen int) string {
	if len(matches) == 0 {
		return text
	}

	// Build result with markers
	result := ""
	lastEnd := 0

	for _, matchStart := range matches {
		result += text[lastEnd:matchStart]
		result += "\x1b[43m" + text[matchStart:matchStart+queryLen] + "\x1b[0m" // Yellow highlight
		lastEnd = matchStart + queryLen
	}

	result += text[lastEnd:]
	return result
}

// SplitIntoLines splits text into lines with positions
type LineInfo struct {
	LineNum   int
	Text      string
	StartPos  int
	EndPos    int
}

// TextToLines converts text into line information
func TextToLines(text string) []LineInfo {
	var lines []LineInfo
	startPos := 0
	lineNum := 1

	for _, char := range text {
		if char == '\n' {
			endPos := startPos + len(text[startPos:strings.Index(text[startPos:], "\n")])
			lines = append(lines, LineInfo{
				LineNum:  lineNum,
				Text:     text[startPos:endPos],
				StartPos: startPos,
				EndPos:   endPos,
			})
			startPos = endPos + 1
			lineNum++
		}
	}

	// Add last line
	if startPos < len(text) {
		lines = append(lines, LineInfo{
			LineNum:  lineNum,
			Text:     text[startPos:],
			StartPos: startPos,
			EndPos:   len(text),
		})
	}

	return lines
}

// FindMatchOnLine finds a match on a specific line
func FindMatchOnLine(text, query string, caseSensitive bool) int {
	if !caseSensitive {
		text = strings.ToLower(text)
		query = strings.ToLower(query)
	}
	return strings.Index(text, query)
}
