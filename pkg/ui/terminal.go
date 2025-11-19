package ui

import (
	"os"
	"strings"
)

// TerminalCapabilities describes what image rendering a terminal supports
type TerminalCapabilities struct {
	SupportsKitty     bool
	SupportsITerm2    bool
	SupportsSIXEL     bool
	SupportsHalfblock bool // Unicode halfblocks (always true)
	PreferredFormat   string
}

// DetectTerminal detects the current terminal and its graphics capabilities
// Based on environment variables and terminal emulator detection
func DetectTerminal() TerminalCapabilities {
	term := os.Getenv("TERM")
	termProgram := os.Getenv("TERM_PROGRAM")
	kittyWindowID := os.Getenv("KITTY_WINDOW_ID")
	itermSessionID := os.Getenv("ITERM_SESSION_ID")

	caps := TerminalCapabilities{
		SupportsHalfblock: true, // Always available as fallback
	}

	// Check for Kitty terminal
	// Kitty sets TERM=xterm-kitty and/or KITTY_WINDOW_ID
	if strings.Contains(term, "kitty") || kittyWindowID != "" {
		caps.SupportsKitty = true
		caps.PreferredFormat = "kitty"
		return caps
	}

	// Check for iTerm2 (macOS)
	if termProgram == "iTerm.app" || itermSessionID != "" {
		caps.SupportsITerm2 = true
		caps.PreferredFormat = "iterm2"
		return caps
	}

	// Check for WezTerm (supports multiple protocols)
	if strings.Contains(termProgram, "WezTerm") || strings.Contains(term, "wezterm") {
		caps.SupportsKitty = true // WezTerm supports Kitty protocol
		caps.SupportsITerm2 = true
		caps.SupportsSIXEL = true
		caps.PreferredFormat = "kitty"
		return caps
	}

	// Check for SIXEL support
	// Xterm, mlterm, and others may support SIXEL
	if strings.Contains(term, "xterm") || strings.Contains(term, "mlterm") {
		caps.SupportsSIXEL = true
		caps.PreferredFormat = "sixel"
		return caps
	}

	// Default fallback: halfblock only
	caps.PreferredFormat = "halfblock"
	return caps
}

// SupportsGraphics returns true if terminal supports any graphics protocol
func (tc *TerminalCapabilities) SupportsGraphics() bool {
	return tc.SupportsKitty || tc.SupportsITerm2 || tc.SupportsSIXEL
}

// GetFallbackChain returns ordered list of formats to try
// First working format in chain will be used
func (tc *TerminalCapabilities) GetFallbackChain() []string {
	chain := []string{}

	// Add preferred first
	if tc.SupportsKitty {
		chain = append(chain, "kitty")
	}
	if tc.SupportsITerm2 {
		chain = append(chain, "iterm2")
	}
	if tc.SupportsSIXEL {
		chain = append(chain, "sixel")
	}

	// Always end with halfblock fallback
	chain = append(chain, "halfblock")

	// If no specific format detected, start with halfblock
	if len(chain) == 1 {
		chain = []string{"halfblock"}
	}

	return chain
}

// IsModernTerminal returns true if terminal supports modern graphics
func (tc *TerminalCapabilities) IsModernTerminal() bool {
	return tc.SupportsKitty || tc.SupportsITerm2 || tc.SupportsSIXEL
}

// ImageRenderingMode specifies how to render images
type ImageRenderingMode string

const (
	// ImageRenderingEnabled - render all images
	ImageRenderingEnabled ImageRenderingMode = "enabled"

	// ImageRenderingText - render only as text placeholders
	ImageRenderingText ImageRenderingMode = "text"

	// ImageRenderingDisabled - don't show images at all
	ImageRenderingDisabled ImageRenderingMode = "disabled"
)

// ImageRenderConfig describes how to render images for a specific terminal
type ImageRenderConfig struct {
	Mode          ImageRenderingMode
	Format        string // "kitty", "iterm2", "sixel", "halfblock", "text"
	MaxWidth      int    // Terminal chars (width)
	MaxHeight     int    // Terminal lines (height)
	CacheImages   bool   // Cache rendered images
	ShowAltText   bool   // Show title/alt text for images
}

// DefaultImageRenderConfig returns sensible defaults
func DefaultImageRenderConfig(width int, height int) ImageRenderConfig {
	return ImageRenderConfig{
		Mode:        ImageRenderingEnabled,
		Format:      "halfblock", // Safe default
		MaxWidth:    width,
		MaxHeight:   height / 2, // Don't use full height
		CacheImages: true,
		ShowAltText: true,
	}
}

// GetImageRenderConfig returns config based on terminal capabilities
func GetImageRenderConfig(width int, height int) ImageRenderConfig {
	caps := DetectTerminal()
	cfg := DefaultImageRenderConfig(width, height)

	if caps.SupportsGraphics() {
		cfg.Format = caps.PreferredFormat
	}

	return cfg
}
