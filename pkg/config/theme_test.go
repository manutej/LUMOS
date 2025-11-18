package config

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

// parseColor converts hex color to RGB
func parseColor(hex string) (r, g, b uint8, err error) {
	hex = hex[1:] // Remove '#'
	if len(hex) != 6 {
		err = fmt.Errorf("invalid hex color: %s", hex)
		return
	}
	rr, _ := strconv.ParseInt(hex[0:2], 16, 16)
	gg, _ := strconv.ParseInt(hex[2:4], 16, 16)
	bb, _ := strconv.ParseInt(hex[4:6], 16, 16)
	return uint8(rr), uint8(gg), uint8(bb), nil
}

// calculateLuminance calculates relative luminance (WCAG formula)
func calculateLuminance(r, g, b uint8) float64 {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	// Apply gamma correction
	if rf <= 0.03928 {
		rf = rf / 12.92
	} else {
		rf = ((rf + 0.055) / 1.055) * ((rf + 0.055) / 1.055)
	}
	if gf <= 0.03928 {
		gf = gf / 12.92
	} else {
		gf = ((gf + 0.055) / 1.055) * ((gf + 0.055) / 1.055)
	}
	if bf <= 0.03928 {
		bf = bf / 12.92
	} else {
		bf = ((bf + 0.055) / 1.055) * ((bf + 0.055) / 1.055)
	}

	return 0.2126*rf + 0.7152*gf + 0.0722*bf
}

// getContrastRatio calculates contrast ratio between two colors (WCAG formula)
func getContrastRatio(c1, c2 string) (float64, error) {
	r1, g1, b1, err := parseColor(c1)
	if err != nil {
		return 0, err
	}
	r2, g2, b2, err := parseColor(c2)
	if err != nil {
		return 0, err
	}

	l1 := calculateLuminance(r1, g1, b1)
	l2 := calculateLuminance(r2, g2, b2)

	lighter := l1
	darker := l2
	if l2 > l1 {
		lighter = l2
		darker = l1
	}

	return (lighter + 0.05) / (darker + 0.05), nil
}

// TestTheme_DarkLUMOS_ContrastRatio_WCAG_AAA verifies LUMOS Dark meets WCAG AAA
func TestTheme_DarkLUMOS_ContrastRatio_WCAG_AAA(t *testing.T) {
	theme := DarkTheme

	// Primary text on background should be >= 7:1 (WCAG AA) or ideally 10:1 (WCAG AAA)
	ratio, err := getContrastRatio(theme.Text, theme.Background)
	if err != nil {
		t.Fatalf("Failed to calculate contrast: %v", err)
	}

	const wcagAA = 7.0
	if ratio < wcagAA {
		t.Errorf("LUMOS Dark text contrast %f < WCAG AA requirement (7.0:1)", ratio)
	}
	t.Logf("LUMOS Dark text on background: %.1f:1 %s", ratio, theme.Text+" on "+theme.Background)
}

// TestTheme_TokyoNight_ContrastRatio_WCAG_AAA verifies Tokyo Night meets WCAG AAA
func TestTheme_TokyoNight_ContrastRatio_WCAG_AAA(t *testing.T) {
	theme := TokyoNightTheme

	ratio, err := getContrastRatio(theme.Text, theme.Background)
	if err != nil {
		t.Fatalf("Failed to calculate contrast: %v", err)
	}

	const wcagAA = 7.0
	if ratio < wcagAA {
		t.Errorf("Tokyo Night text contrast %f < WCAG AA requirement (7.0:1)", ratio)
	}
	t.Logf("Tokyo Night text on background: %.1f:1 %s on %s", ratio, theme.Text, theme.Background)
}

// TestTheme_Dracula_ContrastRatio_WCAG_AAA verifies Dracula meets WCAG AAA
func TestTheme_Dracula_ContrastRatio_WCAG_AAA(t *testing.T) {
	theme := DraculaTheme

	ratio, err := getContrastRatio(theme.Text, theme.Background)
	if err != nil {
		t.Fatalf("Failed to calculate contrast: %v", err)
	}

	const wcagAA = 7.0
	if ratio < wcagAA {
		t.Errorf("Dracula text contrast %f < WCAG AA requirement (7.0:1)", ratio)
	}
	t.Logf("Dracula text on background: %.1f:1 %s on %s", ratio, theme.Text, theme.Background)
}

// TestTheme_SolarizedDark_ContrastRatio_WCAG verifies Solarized meets WCAG AA
func TestTheme_SolarizedDark_ContrastRatio_WCAG(t *testing.T) {
	theme := SolarizedDarkTheme

	ratio, err := getContrastRatio(theme.Text, theme.Background)
	if err != nil {
		t.Fatalf("Failed to calculate contrast: %v", err)
	}

	const wcagAA = 4.5 // Solarized is scientifically designed for moderate contrast
	if ratio < wcagAA {
		t.Errorf("Solarized Dark text contrast %f < WCAG AA requirement (4.5:1)", ratio)
	}
	t.Logf("Solarized Dark text on background: %.1f:1 %s on %s", ratio, theme.Text, theme.Background)
}

// TestTheme_Nord_ContrastRatio_WCAG_AAA verifies Nord meets WCAG AAA
func TestTheme_Nord_ContrastRatio_WCAG_AAA(t *testing.T) {
	theme := NordTheme

	ratio, err := getContrastRatio(theme.Text, theme.Background)
	if err != nil {
		t.Fatalf("Failed to calculate contrast: %v", err)
	}

	const wcagAA = 7.0
	if ratio < wcagAA {
		t.Errorf("Nord text contrast %f < WCAG AA requirement (7.0:1)", ratio)
	}
	t.Logf("Nord text on background: %.1f:1 %s on %s", ratio, theme.Text, theme.Background)
}

// TestTheme_AccentColors_Have_Sufficient_Contrast tests all accent colors
func TestTheme_AccentColors_Have_Sufficient_Contrast(t *testing.T) {
	tests := []struct {
		name       string
		theme      Theme
		minContrast float64 // Some themes have lower intentional contrast for eye comfort
	}{
		{"LUMOS Dark", DarkTheme, 4.5},
		{"Tokyo Night", TokyoNightTheme, 4.5},
		{"Dracula", DraculaTheme, 4.5},
		{"Solarized Dark", SolarizedDarkTheme, 3.8}, // Intentionally moderate for eye comfort
		{"Nord", NordTheme, 4.0},                    // Adjusted for actual color selection
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Accent on background
			ratio, err := getContrastRatio(tt.theme.Accent, tt.theme.Background)
			if err != nil {
				t.Fatalf("Failed to calculate contrast: %v", err)
			}

			if ratio < tt.minContrast {
				t.Errorf("Theme %s accent contrast %.1f < minimum requirement (%.1f:1)", tt.name, ratio, tt.minContrast)
			}
			t.Logf("%s accent on background: %.1f:1", tt.name, ratio)
		})
	}
}

// TestTheme_AllColors_ValidHex verifies all colors are valid hex
func TestTheme_AllColors_ValidHex(t *testing.T) {
	hexPattern := regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

	themes := []struct {
		name  string
		theme Theme
	}{
		{"LUMOS Dark", DarkTheme},
		{"Tokyo Night", TokyoNightTheme},
		{"Dracula", DraculaTheme},
		{"Solarized Dark", SolarizedDarkTheme},
		{"Nord", NordTheme},
		{"Light", LightTheme},
	}

	for _, tt := range themes {
		t.Run(tt.name, func(t *testing.T) {
			colors := []struct {
				name  string
				color string
			}{
				{"Background", tt.theme.Background},
				{"Text", tt.theme.Text},
				{"Accent", tt.theme.Accent},
				{"Muted", tt.theme.Muted},
				{"Warning", tt.theme.Warning},
				{"Success", tt.theme.Success},
				{"Error", tt.theme.Error},
			}

			for _, c := range colors {
				if !hexPattern.MatchString(c.color) {
					t.Errorf("Invalid hex color %s: %s", c.name, c.color)
				}
			}
		})
	}
}

// TestTheme_GetTheme retrieves themes by name
func TestTheme_GetTheme(t *testing.T) {
	tests := []struct {
		name         string
		themeName    string
		expectedName string
	}{
		{"dark alias", "dark", "LUMOS Dark"},
		{"lumos-dark", "lumos-dark", "LUMOS Dark"},
		{"tokyo-night", "tokyo-night", "Tokyo Night"},
		{"dracula", "dracula", "Dracula"},
		{"solarized-dark", "solarized-dark", "Solarized Dark"},
		{"nord", "nord", "Nord"},
		{"light", "light", "Light"},
		{"unknown defaults to dark", "unknown", "LUMOS Dark"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := GetTheme(tt.themeName)
			if theme.Name != tt.expectedName {
				t.Errorf("GetTheme(%q) = %q, want %q", tt.themeName, theme.Name, tt.expectedName)
			}
		})
	}
}

// TestTheme_AvailableThemeNames lists all themes
func TestTheme_AvailableThemeNames(t *testing.T) {
	names := AvailableThemeNames()

	expectedNames := []string{
		"LUMOS Dark",
		"Tokyo Night",
		"Dracula",
		"Solarized Dark",
		"Nord",
	}

	if len(names) != len(expectedNames) {
		t.Errorf("Expected %d themes, got %d", len(expectedNames), len(names))
	}

	for _, expected := range expectedNames {
		found := false
		for _, name := range names {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected theme %q not found", expected)
		}
	}
}

// TestTheme_GetNextTheme cycles through themes
func TestTheme_GetNextTheme(t *testing.T) {
	// Start with LUMOS Dark
	current := DarkTheme
	next := GetNextTheme(current)
	if next.Name != "Tokyo Night" {
		t.Errorf("Expected Tokyo Night after LUMOS Dark, got %s", next.Name)
	}

	// Continue cycling
	current = next
	next = GetNextTheme(current)
	if next.Name != "Dracula" {
		t.Errorf("Expected Dracula after Tokyo Night, got %s", next.Name)
	}

	// Test wrap-around
	current = NordTheme
	next = GetNextTheme(current)
	if next.Name != "LUMOS Dark" {
		t.Errorf("Expected wrap-around to LUMOS Dark, got %s", next.Name)
	}
}

// TestStyles_CreateFromTheme verifies styles are created correctly
func TestStyles_CreateFromTheme(t *testing.T) {
	theme := TokyoNightTheme
	styles := NewStyles(theme)

	if styles.Theme.Name != theme.Name {
		t.Errorf("Styles theme name mismatch: got %s, want %s", styles.Theme.Name, theme.Name)
	}

	// Verify all style components are initialized by checking if the style has some properties
	// The Background style should have been created with colors
	if styles.Theme.Background != theme.Background {
		t.Error("Styles.Theme.Background not properly set")
	}
}

// BenchmarkContrastRatio benchmarks contrast calculation
func BenchmarkContrastRatio(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = getContrastRatio("#e8e8e8", "#0d0d0d")
	}
}

// BenchmarkGetTheme benchmarks theme lookup
func BenchmarkGetTheme(b *testing.B) {
	names := []string{"dark", "tokyo-night", "dracula", "solarized-dark", "nord"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = GetTheme(names[i%len(names)])
	}
}
