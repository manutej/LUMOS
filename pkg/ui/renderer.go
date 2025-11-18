package ui

import (
	"fmt"
	"image"
	"strings"
)

// ImageRenderer handles rendering images for terminal display
type ImageRenderer struct {
	config ImageRenderConfig
}

// NewImageRenderer creates a new image renderer with config
func NewImageRenderer(cfg ImageRenderConfig) *ImageRenderer {
	return &ImageRenderer{config: cfg}
}

// RenderImage converts a terminal image to displayable string
// Format depends on terminal capabilities
func (r *ImageRenderer) RenderImage(img image.Image, altText string) string {
	switch r.config.Mode {
	case ImageRenderingDisabled:
		return ""
	case ImageRenderingText:
		return r.renderAsText(img, altText)
	case ImageRenderingEnabled:
		// Try format, fallback to text if not supported
		return r.renderAsGraphics(img, altText)
	default:
		return r.renderAsText(img, altText)
	}
}

// renderAsGraphics renders image using terminal graphics protocol
func (r *ImageRenderer) renderAsGraphics(img image.Image, altText string) string {
	switch r.config.Format {
	case "kitty":
		return r.renderAsKitty(img, altText)
	case "iterm2":
		return r.renderAsITerm2(img, altText)
	case "sixel":
		return r.renderAsSIXEL(img, altText)
	case "halfblock":
		return r.renderAsHalfblocks(img, altText)
	default:
		return r.renderAsText(img, altText)
	}
}

// renderAsKitty renders image using Kitty graphics protocol
// Returns placeholder - actual implementation requires image encoding
func (r *ImageRenderer) renderAsKitty(img image.Image, altText string) string {
	// TODO: Implement Kitty graphics protocol
	// This requires:
	// 1. Encoding image to base64 PNG
	// 2. Creating Kitty escape sequence
	// 3. Properly resizing to fit terminal width
	// For now, return text fallback

	return r.renderAsText(img, altText)
}

// renderAsITerm2 renders image using iTerm2 inline images
// Returns placeholder - actual implementation requires image encoding
func (r *ImageRenderer) renderAsITerm2(img image.Image, altText string) string {
	// TODO: Implement iTerm2 inline image protocol
	// This requires:
	// 1. Encoding image to base64
	// 2. Creating iTerm2 escape sequence
	// 3. Properly sizing for display
	// For now, return text fallback

	return r.renderAsText(img, altText)
}

// renderAsSIXEL renders image using SIXEL graphics
// Returns placeholder - actual implementation is complex
func (r *ImageRenderer) renderAsSIXEL(img image.Image, altText string) string {
	// SIXEL is a legacy protocol but still supported by some terminals
	// Implementation would require sophisticated image-to-sixel conversion
	// For MVP, use text fallback

	return r.renderAsText(img, altText)
}

// renderAsHalfblocks renders image using Unicode halfblock characters
// This works on all terminals and provides reasonable quality
func (r *ImageRenderer) renderAsHalfblocks(img image.Image, altText string) string {
	// TODO: Implement halfblock rendering
	// This requires:
	// 1. Resizing image to fit terminal
	// 2. Sampling pixels from image
	// 3. Converting to RGB values
	// 4. Finding closest halfblock character
	// For now, return text fallback

	return r.renderAsText(img, altText)
}

// renderAsText renders image as ASCII/text placeholder
// Works on all terminals
func (r *ImageRenderer) renderAsText(img image.Image, altText string) string {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create text placeholder showing image info
	lines := []string{
		strings.Repeat("━", 40),
	}

	if altText != "" {
		lines = append(lines, fmt.Sprintf("[Image: %s]", altText))
	} else {
		lines = append(lines, fmt.Sprintf("[Image: %dx%d]", width, height))
	}

	lines = append(lines, fmt.Sprintf("(%d × %d pixels)", width, height))
	lines = append(lines, strings.Repeat("━", 40))

	return strings.Join(lines, "\n")
}

// CalculateScaledSize calculates how to scale image to fit terminal
func (r *ImageRenderer) CalculateScaledSize(srcWidth, srcHeight int) (dstWidth, dstHeight int) {
	// Calculate aspect ratio
	aspectRatio := float64(srcWidth) / float64(srcHeight)

	// Start with max width
	dstWidth = r.config.MaxWidth
	dstHeight = int(float64(dstWidth) / aspectRatio)

	// If height exceeded, recalculate from height
	if dstHeight > r.config.MaxHeight {
		dstHeight = r.config.MaxHeight
		dstWidth = int(float64(dstHeight) * aspectRatio)
	}

	// Ensure minimum size
	if dstWidth < 5 {
		dstWidth = 5
	}
	if dstHeight < 2 {
		dstHeight = 2
	}

	return dstWidth, dstHeight
}

// FormatLabel returns a formatted label for an image
func (r *ImageRenderer) FormatLabel(title string, width int, height int) string {
	if title != "" {
		return fmt.Sprintf("Image: %s (%dx%d)", title, width, height)
	}
	return fmt.Sprintf("Image (%dx%d)", width, height)
}

// SupportedFormats returns list of supported image formats for current terminal
func (r *ImageRenderer) SupportedFormats() []string {
	caps := DetectTerminal()
	return caps.GetFallbackChain()
}
