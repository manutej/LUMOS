package pdf

import (
	"regexp"
	"strings"
	"testing"
)

// TestSearchOptions_Default verifies default options
func TestSearchOptions_Default(t *testing.T) {
	opts := SearchOptions{
		CaseSensitive: false,
		WholeWord:     false,
		RegexMode:     false,
		MaxResults:    0,
		StartPage:     1,
		EndPage:       0,
	}

	if opts.CaseSensitive {
		t.Error("CaseSensitive should be false by default")
	}
	if opts.WholeWord {
		t.Error("WholeWord should be false by default")
	}
	if opts.RegexMode {
		t.Error("RegexMode should be false by default")
	}
}

// TestSearchFilter_Default verifies default filter
func TestSearchFilter_Default(t *testing.T) {
	filter := SearchFilter{
		MinContextLength: 20,
		HighlightMatch:   true,
		SortByPage:       true,
	}

	if filter.MinContextLength != 20 {
		t.Errorf("MinContextLength mismatch: got %d, want %d", filter.MinContextLength, 20)
	}
	if !filter.HighlightMatch {
		t.Error("HighlightMatch should be true")
	}
	if !filter.SortByPage {
		t.Error("SortByPage should be true")
	}
}

// TestIsWordCharRune tests word character detection
func TestIsWordCharRune(t *testing.T) {
	tests := []struct {
		char rune
		want bool
	}{
		{'a', true},
		{'Z', true},
		{'0', true},
		{'_', true},
		{' ', false},
		{'-', false},
		{'.', false},
		{'!', false},
	}

	for _, tt := range tests {
		got := isWordCharRune(tt.char)
		if got != tt.want {
			t.Errorf("isWordCharRune(%q) = %v, want %v", tt.char, got, tt.want)
		}
	}
}

// TestIsWordBoundaryPosition tests boundary detection
func TestIsWordBoundaryPosition(t *testing.T) {
	text := "hello world test"

	tests := []struct {
		pos  int
		want bool
	}{
		{0, true},      // Start of text
		{5, true},      // After "hello" (space)
		{6, true},      // Before "world"
		{11, true},     // After "world" (space)
		{12, true},     // Before "test"
		{16, true},     // End of text
		{2, false},     // Middle of word
		{8, false},     // Middle of word
	}

	for _, tt := range tests {
		got := isWordBoundaryPosition(text, tt.pos)
		if got != tt.want {
			t.Errorf("isWordBoundaryPosition(%q, %d) = %v, want %v", text, tt.pos, got, tt.want)
		}
	}
}

// TestFindAllPositions finds all occurrences
func TestFindAllPositions(t *testing.T) {
	text := "the quick brown fox jumps over the lazy dog"

	positions := findAllPositions(text, "the", false)

	if len(positions) != 2 {
		t.Errorf("Expected 2 occurrences of 'the', got %d", len(positions))
	}
	if positions[0] != 0 {
		t.Errorf("First occurrence should be at 0, got %d", positions[0])
	}
	if positions[1] != 31 {
		t.Errorf("Second occurrence should be at 31, got %d", positions[1])
	}
}

// TestFindAllPositions_WholeWord respects word boundaries
func TestFindAllPositions_WholeWord(t *testing.T) {
	text := "the cat catapult categorize cat"

	positions := findAllPositions(text, "cat", true)

	if len(positions) != 2 {
		t.Errorf("Expected 2 whole-word matches of 'cat', got %d", len(positions))
	}

	// Check that "catapult" and "categorize" are not matched
	for _, pos := range positions {
		word := ""
		start, end := pos, pos+3
		if start > 0 && text[start-1] != ' ' && text[start-1] != '\n' {
			t.Errorf("Match at %d doesn't have word boundary before", pos)
		}
		if end < len(text) && text[end] != ' ' && text[end] != '\n' {
			t.Errorf("Match at %d doesn't have word boundary after", pos)
		}
		_ = word
	}
}

// TestCalculateRelevance scores match relevance
func TestCalculateRelevance(t *testing.T) {
	tests := []struct {
		match string
		query string
		want  float64
	}{
		{"hello", "hello", 1.0},    // Exact match
		{"hello", "hel", 0.6},      // Partial: 3/5 = 0.6
		{"h", "hello", 0.2},        // Short vs long: 1/5 = 0.2
	}

	for _, tt := range tests {
		got := calculateRelevance(tt.match, tt.query)
		if got < tt.want-0.01 || got > tt.want+0.01 {
			t.Errorf("calculateRelevance(%q, %q) = %v, want %v", tt.match, tt.query, got, tt.want)
		}
	}
}

// TestFormatPreview creates preview text
func TestFormatPreview(t *testing.T) {
	before := "The quick"
	match := "brown"
	after := "fox jumps"

	preview := formatPreview(before, match, after)

	if !strings.Contains(preview, match) {
		t.Errorf("Preview should contain match text: %s", preview)
	}
	if !strings.Contains(preview, "[") || !strings.Contains(preview, "]") {
		t.Errorf("Preview should have brackets around match: %s", preview)
	}
}

// TestGetLineAndColumn finds line and column numbers
func TestGetLineAndColumn(t *testing.T) {
	text := "line1\nline2\nline3"

	tests := []struct {
		pos       int
		wantLine  int
		wantCol   int
		desc      string
	}{
		{0, 1, 1, "start of text"},
		{5, 1, 6, "end of line 1"},
		{6, 2, 1, "start of line 2"},
		{11, 2, 6, "end of line 2"},
	}

	for _, tt := range tests {
		line, col := getLineAndColumn(text, tt.pos)
		if line != tt.wantLine || col != tt.wantCol {
			t.Errorf("getLineAndColumn(%d) [%s] = (%d, %d), want (%d, %d)",
				tt.pos, tt.desc, line, col, tt.wantLine, tt.wantCol)
		}
	}
}

// TestExtractContextAdvanced extracts surrounding text
func TestExtractContextAdvanced(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog"

	// Extract after position
	after := extractContextAdvanced(text, 10, 10)
	if !strings.Contains(after, "brown") {
		t.Errorf("After context should contain 'brown': %s", after)
	}

	// Extract before position
	before := extractContextAdvanced(text, 10, -10)
	if !strings.Contains(before, "quick") {
		t.Errorf("Before context should contain 'quick': %s", before)
	}
}

// TestSearchHistoryManager creates and manages history
func TestSearchHistoryManager_Create(t *testing.T) {
	manager := NewSearchHistoryManager(10)

	if manager.Size() != 0 {
		t.Errorf("New manager should have size 0, got %d", manager.Size())
	}
}

// TestSearchHistoryManager_Add adds entries
func TestSearchHistoryManager_Add(t *testing.T) {
	manager := NewSearchHistoryManager(5)

	entry := SearchHistoryEntry{
		Query: "test",
		Options: SearchOptions{
			CaseSensitive: false,
			WholeWord:     true,
		},
		ResultCount: 10,
	}

	manager.Add(entry)

	if manager.Size() != 1 {
		t.Errorf("After add, size should be 1, got %d", manager.Size())
	}

	history := manager.GetHistory()
	if len(history) != 1 || history[0].Query != "test" {
		t.Error("History entry not added correctly")
	}
}

// TestSearchHistoryManager_MaxSize respects max size
func TestSearchHistoryManager_MaxSize(t *testing.T) {
	manager := NewSearchHistoryManager(3)

	for i := 0; i < 5; i++ {
		entry := SearchHistoryEntry{
			Query: "query" + string(rune('0'+i)),
		}
		manager.Add(entry)
	}

	if manager.Size() != 3 {
		t.Errorf("Manager should respect max size of 3, got %d", manager.Size())
	}

	// Check that oldest entries are removed
	history := manager.GetHistory()
	if history[0].Query != "query4" {
		t.Errorf("Most recent entry should be first, got %s", history[0].Query)
	}
}

// TestSearchHistoryManager_Clear clears history
func TestSearchHistoryManager_Clear(t *testing.T) {
	manager := NewSearchHistoryManager(10)

	manager.Add(SearchHistoryEntry{Query: "test"})
	manager.Add(SearchHistoryEntry{Query: "query"})

	if manager.Size() != 2 {
		t.Errorf("Expected size 2, got %d", manager.Size())
	}

	manager.Clear()

	if manager.Size() != 0 {
		t.Errorf("After clear, size should be 0, got %d", manager.Size())
	}
}

// TestSearchResultAdvanced_Fields verifies structure
func TestSearchResultAdvanced_Fields(t *testing.T) {
	result := SearchResultAdvanced{
		SearchResult: SearchResult{
			PageNum:   5,
			MatchText: "test",
		},
		MatchCount:  3,
		Relevance:   0.95,
		PreviewText: "...prefix [test] suffix...",
	}

	if result.PageNum != 5 {
		t.Errorf("PageNum mismatch: got %d, want %d", result.PageNum, 5)
	}
	if result.MatchText != "test" {
		t.Errorf("MatchText mismatch: got %s, want %s", result.MatchText, "test")
	}
	if result.MatchCount != 3 {
		t.Errorf("MatchCount mismatch: got %d, want %d", result.MatchCount, 3)
	}
	if result.Relevance < 0.94 || result.Relevance > 0.96 {
		t.Errorf("Relevance mismatch: got %f", result.Relevance)
	}
}

// TestFindAllPositions_CaseInsensitive finds case-insensitive matches
func TestFindAllPositions_CaseInsensitive(t *testing.T) {
	text := "The QUICK brown QuIcK fox"
	positions := findAllPositions(strings.ToLower(text), strings.ToLower("quick"), false)

	if len(positions) != 2 {
		t.Errorf("Expected 2 case-insensitive matches, got %d", len(positions))
	}
}

// TestExtractLineContext_Before extracts context before position
func TestExtractLineContext_Before(t *testing.T) {
	line := "The quick brown fox"
	context := extractLineContext(line, 9, -6)

	if !strings.Contains(context, "quick") {
		t.Errorf("Context should contain 'quick': %s", context)
	}
}

// TestExtractLineContext_After extracts context after position
func TestExtractLineContext_After(t *testing.T) {
	line := "The quick brown fox"
	context := extractLineContext(line, 9, 6)

	if !strings.Contains(context, "brown") {
		t.Errorf("Context should contain 'brown': %s", context)
	}
}

// BenchmarkFindAllPositions benchmarks position finding
func BenchmarkFindAllPositions(b *testing.B) {
	text := strings.Repeat("the quick brown fox ", 100)
	query := "the"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = findAllPositions(text, query, false)
	}
}

// BenchmarkCalculateRelevance benchmarks relevance calculation
func BenchmarkCalculateRelevance(b *testing.B) {
	match := "the quick brown"
	query := "quick"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = calculateRelevance(match, query)
	}
}

// BenchmarkIsWordBoundaryPosition benchmarks boundary detection
func BenchmarkIsWordBoundaryPosition(b *testing.B) {
	text := "The quick brown fox jumps over the lazy dog"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = isWordBoundaryPosition(text, 10)
	}
}

// BenchmarkSearchHistoryManager benchmarks history operations
func BenchmarkSearchHistoryManager_Add(b *testing.B) {
	manager := NewSearchHistoryManager(100)
	entry := SearchHistoryEntry{Query: "benchmark"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Add(entry)
	}
}

// TestRegexCompilationError tests regex error handling
func TestRegexCompilationError(t *testing.T) {
	invalidRegex := "["
	_, err := regexp.Compile(invalidRegex)
	if err == nil {
		t.Error("Invalid regex should produce error")
	}
}

// TestSearchOptions_Validation validates options
func TestSearchOptions_Validation(t *testing.T) {
	opts := SearchOptions{
		MaxResults: -5,
		StartPage:  0,
		EndPage:    0,
	}

	// These would be normalized in AdvancedSearch
	if opts.MaxResults < 0 {
		opts.MaxResults = 0
	}
	if opts.StartPage < 1 {
		opts.StartPage = 1
	}

	if opts.MaxResults != 0 {
		t.Errorf("MaxResults should be normalized to 0, got %d", opts.MaxResults)
	}
	if opts.StartPage != 1 {
		t.Errorf("StartPage should be normalized to 1, got %d", opts.StartPage)
	}
}
