package pdf

import (
	"reflect"
	"strings"
	"testing"
)

// TestNewTextSearch tests creating a new text search
func TestNewTextSearch(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{"simple query", "hello"},
		{"empty query", ""},
		{"multi-word query", "hello world"},
		{"special characters", "test@#$%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTextSearch(tt.query)

			if ts == nil {
				t.Fatal("NewTextSearch() returned nil")
			}

			if ts.query != tt.query {
				t.Errorf("query = %v, want %v", ts.query, tt.query)
			}

			if ts.current != 0 {
				t.Errorf("current = %v, want 0", ts.current)
			}

			if ts.results != nil {
				t.Errorf("results should be nil initially, got %v", ts.results)
			}
		})
	}
}

// TestTextSearch_GetResultCount tests result count
func TestTextSearch_GetResultCount(t *testing.T) {
	ts := NewTextSearch("test")

	// Initially zero
	if count := ts.GetResultCount(); count != 0 {
		t.Errorf("initial count = %v, want 0", count)
	}

	// Add some mock results
	ts.results = []SearchResult{
		{PageNum: 1, MatchText: "test"},
		{PageNum: 2, MatchText: "test"},
	}

	if count := ts.GetResultCount(); count != 2 {
		t.Errorf("count = %v, want 2", count)
	}
}

// TestTextSearch_Navigation tests result navigation
func TestTextSearch_Navigation(t *testing.T) {
	ts := NewTextSearch("test")

	// Setup test results
	ts.results = []SearchResult{
		{PageNum: 1, MatchText: "test1"},
		{PageNum: 2, MatchText: "test2"},
		{PageNum: 3, MatchText: "test3"},
	}

	// Test initial state
	if ts.GetCurrentIndex() != 0 {
		t.Errorf("initial index = %v, want 0", ts.GetCurrentIndex())
	}

	// Test NextResult
	if !ts.NextResult() {
		t.Error("NextResult() should succeed")
	}
	if ts.GetCurrentIndex() != 1 {
		t.Errorf("after next, index = %v, want 1", ts.GetCurrentIndex())
	}

	// Test PreviousResult
	if !ts.PreviousResult() {
		t.Error("PreviousResult() should succeed")
	}
	if ts.GetCurrentIndex() != 0 {
		t.Errorf("after previous, index = %v, want 0", ts.GetCurrentIndex())
	}

	// Test PreviousResult at boundary
	if ts.PreviousResult() {
		t.Error("PreviousResult() should fail at start boundary")
	}

	// Go to end
	ts.JumpToResult(2)

	// Test NextResult at boundary
	if ts.NextResult() {
		t.Error("NextResult() should fail at end boundary")
	}
}

// TestTextSearch_JumpToResult tests jumping to specific result
func TestTextSearch_JumpToResult(t *testing.T) {
	ts := NewTextSearch("test")
	ts.results = []SearchResult{
		{PageNum: 1},
		{PageNum: 2},
		{PageNum: 3},
	}

	tests := []struct {
		name    string
		index   int
		wantOK  bool
	}{
		{"jump to first", 0, true},
		{"jump to middle", 1, true},
		{"jump to last", 2, true},
		{"jump to negative", -1, false},
		{"jump beyond range", 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := ts.JumpToResult(tt.index)
			if ok != tt.wantOK {
				t.Errorf("JumpToResult(%v) = %v, want %v", tt.index, ok, tt.wantOK)
			}

			if tt.wantOK && ts.GetCurrentIndex() != tt.index {
				t.Errorf("after jump, index = %v, want %v", ts.GetCurrentIndex(), tt.index)
			}
		})
	}
}

// TestTextSearch_Reset tests resetting search
func TestTextSearch_Reset(t *testing.T) {
	ts := NewTextSearch("test")
	ts.results = []SearchResult{{PageNum: 1}, {PageNum: 2}}
	ts.current = 1

	ts.Reset()

	if ts.GetCurrentIndex() != 0 {
		t.Errorf("after reset, index = %v, want 0", ts.GetCurrentIndex())
	}
}

// TestCaseSensitiveMatch tests case-sensitive matching
func TestCaseSensitiveMatch(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		query string
		want  []int
	}{
		{
			name:  "simple match",
			text:  "hello world hello",
			query: "hello",
			want:  []int{0, 12},
		},
		{
			name:  "case sensitive - no match",
			text:  "Hello world",
			query: "hello",
			want:  []int{},
		},
		{
			name:  "case sensitive - exact match",
			text:  "Hello world",
			query: "Hello",
			want:  []int{0},
		},
		{
			name:  "no match",
			text:  "hello world",
			query: "foo",
			want:  []int{},
		},
		{
			name:  "overlapping matches",
			text:  "aaa",
			query: "aa",
			want:  []int{0, 1},
		},
		{
			name:  "empty query",
			text:  "hello",
			query: "",
			want:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CaseSensitiveMatch(tt.text, tt.query)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CaseSensitiveMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCaseInsensitiveMatch tests case-insensitive matching
func TestCaseInsensitiveMatch(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		query string
		want  []int
	}{
		{
			name:  "mixed case",
			text:  "Hello World hello",
			query: "hello",
			want:  []int{0, 12},
		},
		{
			name:  "uppercase query",
			text:  "hello world",
			query: "HELLO",
			want:  []int{0},
		},
		{
			name:  "mixed case query",
			text:  "hello world WORLD",
			query: "WoRlD",
			want:  []int{6, 12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CaseInsensitiveMatch(tt.text, tt.query)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CaseInsensitiveMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestWordMatch tests whole-word matching
func TestWordMatch(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		query string
		want  []int
	}{
		{
			name:  "exact word",
			text:  "hello world",
			query: "hello",
			want:  []int{0},
		},
		{
			name:  "word in middle",
			text:  "the hello world",
			query: "hello",
			want:  []int{4},
		},
		{
			name:  "partial word - no match",
			text:  "helloworld",
			query: "hello",
			want:  []int{},
		},
		{
			name:  "word at end",
			text:  "say hello",
			query: "hello",
			want:  []int{4},
		},
		{
			name:  "multiple word matches",
			text:  "hello there hello",
			query: "hello",
			want:  []int{0, 12},
		},
		{
			name:  "word with punctuation",
			text:  "hello, world!",
			query: "hello",
			want:  []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WordMatch(tt.text, tt.query)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WordMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestExtractContext tests context extraction
func TestExtractContext(t *testing.T) {
	text := "This is a sample text for testing context extraction."
	//         0         10        20        30        40        50

	tests := []struct {
		name         string
		matchPos     int
		matchLen     int
		contextChars int
		wantBefore   string
		wantMatch    string
		wantAfter    string
	}{
		{
			name:         "middle of text",
			matchPos:     10,
			matchLen:     6,
			contextChars: 5,
			wantBefore:   "is a ",
			wantMatch:    "sample",
			wantAfter:    " text",
		},
		{
			name:         "start of text",
			matchPos:     0,
			matchLen:     4,
			contextChars: 5,
			wantBefore:   "",
			wantMatch:    "This",
			wantAfter:    " is a",
		},
		{
			name:         "end of text",
			matchPos:     46,
			matchLen:     11,
			contextChars: 5,
			wantBefore:   " extr",
			wantMatch:    "action.",  // truncated from "extraction." due to text length
			wantAfter:    "",
		},
		{
			name:         "zero context",
			matchPos:     10,
			matchLen:     6,
			contextChars: 0,
			wantBefore:   "",
			wantMatch:    "sample",
			wantAfter:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before, match, after := ExtractContext(text, tt.matchPos, tt.matchLen, tt.contextChars)

			if before != tt.wantBefore {
				t.Errorf("before = %q, want %q", before, tt.wantBefore)
			}
			if match != tt.wantMatch {
				t.Errorf("match = %q, want %q", match, tt.wantMatch)
			}
			if after != tt.wantAfter {
				t.Errorf("after = %q, want %q", after, tt.wantAfter)
			}
		})
	}
}

// TestHighlightMatches tests match highlighting
func TestHighlightMatches(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		matches  []int
		queryLen int
		contains string // substring that should be in result
	}{
		{
			name:     "no matches",
			text:     "hello world",
			matches:  []int{},
			queryLen: 5,
			contains: "hello world",
		},
		{
			name:     "single match",
			text:     "hello world",
			matches:  []int{0},
			queryLen: 5,
			contains: "\x1b[43m", // ANSI color code
		},
		{
			name:     "multiple matches",
			text:     "hello hello",
			matches:  []int{0, 6},
			queryLen: 5,
			contains: "\x1b[43m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HighlightMatches(tt.text, tt.matches, tt.queryLen)

			if !strings.Contains(got, tt.contains) {
				t.Errorf("HighlightMatches() result doesn't contain %q", tt.contains)
			}

			// Should contain the original text
			if len(tt.matches) == 0 && got != tt.text {
				t.Errorf("HighlightMatches() with no matches should return original text")
			}
		})
	}
}

// TestTextToLines tests converting text to line information
func TestTextToLines(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		wantLines int
	}{
		{
			name:      "single line",
			text:      "hello",
			wantLines: 1,
		},
		{
			name:      "two lines",
			text:      "hello\nworld",
			wantLines: 2,
		},
		{
			name:      "three lines",
			text:      "line1\nline2\nline3",
			wantLines: 3,
		},
		{
			name:      "empty text",
			text:      "",
			wantLines: 0,
		},
		{
			name:      "trailing newline",
			text:      "hello\n",
			wantLines: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := TextToLines(tt.text)

			if len(lines) != tt.wantLines {
				t.Errorf("TextToLines() returned %v lines, want %v", len(lines), tt.wantLines)
			}

			// Verify line numbers are sequential
			for i, line := range lines {
				expectedLineNum := i + 1
				if line.LineNum != expectedLineNum {
					t.Errorf("line[%d].LineNum = %v, want %v", i, line.LineNum, expectedLineNum)
				}

				// Verify positions are valid
				if line.StartPos < 0 {
					t.Errorf("line[%d].StartPos = %v, want >= 0", i, line.StartPos)
				}
				if line.EndPos < line.StartPos {
					t.Errorf("line[%d].EndPos = %v, want >= StartPos(%v)", i, line.EndPos, line.StartPos)
				}
			}
		})
	}
}

// TestFindMatchOnLine tests finding matches on a line
func TestFindMatchOnLine(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		query         string
		caseSensitive bool
		want          int
	}{
		{
			name:          "case sensitive match",
			text:          "hello world",
			query:         "hello",
			caseSensitive: true,
			want:          0,
		},
		{
			name:          "case sensitive no match",
			text:          "Hello world",
			query:         "hello",
			caseSensitive: true,
			want:          -1,
		},
		{
			name:          "case insensitive match",
			text:          "Hello world",
			query:         "hello",
			caseSensitive: false,
			want:          0,
		},
		{
			name:          "no match",
			text:          "hello world",
			query:         "foo",
			caseSensitive: false,
			want:          -1,
		},
		{
			name:          "match in middle",
			text:          "the hello world",
			query:         "hello",
			caseSensitive: false,
			want:          4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMatchOnLine(tt.text, tt.query, tt.caseSensitive)
			if got != tt.want {
				t.Errorf("FindMatchOnLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTextSearch_GetCurrentResult tests getting current result
func TestTextSearch_GetCurrentResult(t *testing.T) {
	ts := NewTextSearch("test")

	// No results initially
	if result := ts.GetCurrentResult(); result != nil {
		t.Error("GetCurrentResult() should return nil when no results")
	}

	// Add results
	ts.results = []SearchResult{
		{PageNum: 1, MatchText: "test1"},
		{PageNum: 2, MatchText: "test2"},
	}

	// Get first result
	result := ts.GetCurrentResult()
	if result == nil {
		t.Fatal("GetCurrentResult() returned nil with results")
	}
	if result.PageNum != 1 {
		t.Errorf("current result PageNum = %v, want 1", result.PageNum)
	}

	// Move to next and get result
	ts.NextResult()
	result = ts.GetCurrentResult()
	if result == nil {
		t.Fatal("GetCurrentResult() returned nil after NextResult")
	}
	if result.PageNum != 2 {
		t.Errorf("current result PageNum = %v, want 2", result.PageNum)
	}
}

// TestTextSearch_GetResults tests getting all results
func TestTextSearch_GetResults(t *testing.T) {
	ts := NewTextSearch("test")

	// Initially empty
	results := ts.GetResults()
	if results != nil {
		t.Error("GetResults() should return nil initially")
	}

	// Add results
	expectedResults := []SearchResult{
		{PageNum: 1, MatchText: "test1"},
		{PageNum: 2, MatchText: "test2"},
	}
	ts.results = expectedResults

	results = ts.GetResults()
	if !reflect.DeepEqual(results, expectedResults) {
		t.Errorf("GetResults() = %v, want %v", results, expectedResults)
	}
}

// Benchmark tests
func BenchmarkCaseSensitiveMatch(b *testing.B) {
	text := strings.Repeat("hello world ", 1000)
	query := "world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CaseSensitiveMatch(text, query)
	}
}

func BenchmarkCaseInsensitiveMatch(b *testing.B) {
	text := strings.Repeat("Hello World ", 1000)
	query := "world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CaseInsensitiveMatch(text, query)
	}
}

func BenchmarkWordMatch(b *testing.B) {
	text := strings.Repeat("hello world testing ", 1000)
	query := "world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WordMatch(text, query)
	}
}

func BenchmarkHighlightMatches(b *testing.B) {
	text := strings.Repeat("hello world ", 1000)
	matches := []int{0, 12, 24, 36, 48}
	queryLen := 5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HighlightMatches(text, matches, queryLen)
	}
}
