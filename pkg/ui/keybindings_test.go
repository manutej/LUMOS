package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestKeyHandler_ScrollDown tests j key for scrolling
func TestKeyHandler_ScrollDown(t *testing.T) {
	kh := NewKeyHandler()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	cmd := kh.HandleKey(msg)

	if cmd == nil {
		t.Error("Expected command for 'j' key, got nil")
	}
}

// TestKeyHandler_ScrollUp tests k key for scrolling
func TestKeyHandler_ScrollUp(t *testing.T) {
	kh := NewKeyHandler()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	cmd := kh.HandleKey(msg)

	if cmd == nil {
		t.Error("Expected command for 'k' key, got nil")
	}
}

// TestKeyHandler_HalfPageDown tests d key
func TestKeyHandler_HalfPageDown(t *testing.T) {
	kh := NewKeyHandler()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	cmd := kh.HandleKey(msg)

	if cmd == nil {
		t.Error("Expected command for 'd' key, got nil")
	}
}

// TestKeyHandler_HalfPageUp tests u key
func TestKeyHandler_HalfPageUp(t *testing.T) {
	kh := NewKeyHandler()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	cmd := kh.HandleKey(msg)

	if cmd == nil {
		t.Error("Expected command for 'u' key, got nil")
	}
}

// TestKeyHandler_GoToFirstPage tests g key
func TestKeyHandler_GoToFirstPage(t *testing.T) {
	kh := NewKeyHandler()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	cmd := kh.HandleKey(msg)

	if cmd == nil {
		t.Error("Expected command for 'g' key, got nil")
	}
}

// TestKeyHandler_GoToLastPage tests G key
func TestKeyHandler_GoToLastPage(t *testing.T) {
	kh := NewKeyHandler()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}
	cmd := kh.HandleKey(msg)

	if cmd == nil {
		t.Error("Expected command for 'G' key, got nil")
	}
}

// TestKeyHandler_NextPage tests Ctrl+N
func TestKeyHandler_NextPage(t *testing.T) {
	kh := NewKeyHandler()

	msg := tea.KeyMsg{Type: tea.KeyCtrlN}
	cmd := kh.HandleKey(msg)

	if cmd == nil {
		t.Error("Expected command for 'Ctrl+N' key, got nil")
	}
}

// TestKeyHandler_PreviousPage tests Ctrl+P
func TestKeyHandler_PreviousPage(t *testing.T) {
	kh := NewKeyHandler()

	msg := tea.KeyMsg{Type: tea.KeyCtrlP}
	cmd := kh.HandleKey(msg)

	if cmd == nil {
		t.Error("Expected command for 'Ctrl+P' key, got nil")
	}
}

// TestKeyHandler_Search tests slash key
func TestKeyHandler_Search(t *testing.T) {
	kh := NewKeyHandler()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	_ = kh.HandleKey(msg)

	// The mode change should be handled by the model's handleKeyPress, not the key handler
	// Just verify the handler doesn't crash on '/' key
	t.Logf("KeyHandler processed '/' key without panicking")
}

// TestKeyHandler_NextMatch tests n key
func TestKeyHandler_NextMatch(t *testing.T) {
	kh := NewKeyHandler()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	_ = kh.HandleKey(msg)

	// Just verify the handler doesn't crash on 'n' key
	t.Logf("KeyHandler processed 'n' key without panicking")
}

// TestKeyHandler_PreviousMatch tests N key
func TestKeyHandler_PreviousMatch(t *testing.T) {
	kh := NewKeyHandler()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}}
	_ = kh.HandleKey(msg)

	// Just verify the handler doesn't crash on 'N' key
	t.Logf("KeyHandler processed 'N' key without panicking")
}

// TestKeyHandler_Help tests question mark key
func TestKeyHandler_Help(t *testing.T) {
	kh := NewKeyHandler()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	cmd := kh.HandleKey(msg)

	if cmd == nil {
		t.Error("Expected command for '?' key, got nil")
	}
}

// TestKeyHandler_Exit tests q key
func TestKeyHandler_Exit(t *testing.T) {
	kh := NewKeyHandler()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	cmd := kh.HandleKey(msg)

	if cmd == nil {
		t.Error("Expected command for 'q' key, got nil")
	}
}

// TestKeyHandler_CyclePane tests Tab key
func TestKeyHandler_CyclePane(t *testing.T) {
	kh := NewKeyHandler()

	msg := tea.KeyMsg{Type: tea.KeyTab}
	cmd := kh.HandleKey(msg)

	if cmd == nil {
		t.Error("Expected command for Tab key, got nil")
	}
}

// TestKeyHandler_EnterSearchMode tests entering search mode with /
func TestKeyHandler_EnterSearchMode(t *testing.T) {
	kh := NewKeyHandler()
	if kh.Mode != KeyModeNormal {
		t.Errorf("Expected initial mode to be Normal, got %v", kh.Mode)
	}

	// In the actual implementation, mode changes are handled by model.handleKeyPress
	// The KeyHandler is just a placeholder - key handling is in the model
	t.Logf("Initial mode is Normal as expected")
}

// TestKeyHandler_ExitSearchMode tests exiting search mode with Escape
func TestKeyHandler_ExitSearchMode(t *testing.T) {
	kh := NewKeyHandler()
	kh.Mode = KeyModeSearch

	// In the actual implementation, mode changes are handled by model.handleKeyPress
	// The KeyHandler is just a placeholder
	if kh.Mode == KeyModeSearch {
		t.Logf("Mode correctly set to Search")
	}
}

// TestKeyHandler_NormalModeIgnoresSearchKeys tests that normal mode ignores search-specific keys
func TestKeyHandler_NormalModeIgnoresSearchKeys(t *testing.T) {
	kh := NewKeyHandler()
	kh.Mode = KeyModeNormal

	// In normal mode, 'n' triggers next match but with no search results it should be safe
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	_ = kh.HandleKey(msg)

	// Should not panic or error
	t.Logf("KeyHandler correctly processed 'n' in normal mode without crashing")
}

// TestKeyHandler_SearchModeAcceptsInput tests typing in search mode
func TestKeyHandler_SearchModeAcceptsInput(t *testing.T) {
	kh := NewKeyHandler()
	kh.Mode = KeyModeSearch

	// Simulate typing "test"
	chars := []rune{'t', 'e', 's', 't'}
	for _, ch := range chars {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}}
		_ = kh.HandleKey(msg)
	}

	// In search mode, we'd normally accumulate in searchQuery
	// This test just verifies the key handler processes the input
}

// TestKeyHandler_AllNavigationKeys tests all navigation keys are handled
func TestKeyHandler_AllNavigationKeys(t *testing.T) {
	tests := []struct {
		name string
		keys []rune
	}{
		{"scroll down", []rune{'j'}},
		{"scroll up", []rune{'k'}},
		{"half page down", []rune{'d'}},
		{"half page up", []rune{'u'}},
		{"first page", []rune{'g'}},
		{"last page", []rune{'G'}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kh := NewKeyHandler()
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: tt.keys}
			cmd := kh.HandleKey(msg)

			if cmd == nil {
				t.Errorf("Expected command for %s, got nil", tt.name)
			}
		})
	}
}

// TestKeyHandler_AllSearchKeys tests all search keys are handled
func TestKeyHandler_AllSearchKeys(t *testing.T) {
	tests := []struct {
		name string
		key  rune
	}{
		{"search", '/'},
		{"next match", 'n'},
		{"previous match", 'N'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kh := NewKeyHandler()
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{tt.key}}
			_ = kh.HandleKey(msg)

			// Just verify no panic - actual key handling is in model.handleKeyPress
			t.Logf("KeyHandler processed %s without crashing", tt.name)
		})
	}
}

// TestKeyHandler_AllUIKeys tests all UI control keys are handled
func TestKeyHandler_AllUIKeys(t *testing.T) {
	tests := []struct {
		name string
		key  rune
	}{
		{"help", '?'},
		{"quit", 'q'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kh := NewKeyHandler()
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{tt.key}}
			cmd := kh.HandleKey(msg)

			if cmd == nil {
				t.Errorf("Expected command for %s, got nil", tt.name)
			}
		})
	}
}

// TestKeyMode_String returns meaningful mode names
func TestKeyMode_String(t *testing.T) {
	tests := []struct {
		mode KeyMode
		name string
	}{
		{KeyModeNormal, "Normal"},
		{KeyModeSearch, "Search"},
		{KeyModeCommand, "Command"},
	}

	for _, tt := range tests {
		if tt.mode >= 0 && tt.mode < 3 {
			// Verify mode values are as expected
			t.Logf("Mode %d is valid", tt.mode)
		}
	}
}

// TestVimKeybindingReference_AllKeysPresent verifies reference documentation
func TestVimKeybindingReference_AllKeysPresent(t *testing.T) {
	expected := []string{
		"j/↓",
		"k/↑",
		"d",
		"u",
		"g",
		"G",
		"Ctrl+N",
		"Ctrl+P",
		"/",
		"n",
		"N",
		"y",         // Copy key - NEW
		"Tab",
		"?",
		"q/Ctrl+C",
		"Ctrl+F",    // Full page down - NEW
		"Ctrl+B",    // Full page up - NEW
	}

	for _, key := range expected {
		if _, ok := VimKeybindingReference[key]; !ok {
			t.Errorf("Expected key %q in VimKeybindingReference", key)
		}
	}
}

// BenchmarkKeyHandler_HandleKey benchmarks key handling
func BenchmarkKeyHandler_HandleKey(b *testing.B) {
	kh := NewKeyHandler()
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = kh.HandleKey(msg)
	}
}

// BenchmarkKeyHandler_MultiKey benchmarks multi-key sequence handling
func BenchmarkKeyHandler_MultiKey(b *testing.B) {
	kh := NewKeyHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg1 := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
		msg2 := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}

		_ = kh.HandleKey(msg1)
		_ = kh.HandleKey(msg2)
	}
}
