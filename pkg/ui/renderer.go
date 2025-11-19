package ui

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
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
// Encodes image as PNG and sends via Kitty escape sequence
func (r *ImageRenderer) renderAsKitty(img image.Image, altText string) string {
	// Resize image to fit terminal
	scaledWidth, scaledHeight := r.CalculateScaledSize(img.Bounds().Dx(), img.Bounds().Dy())
	resized := r.resizeImage(img, scaledWidth, scaledHeight)

	// Encode as PNG
	pngData, err := r.encodeImageAsPNG(resized)
	if err != nil {
		return r.renderAsText(img, altText)
	}

	// Base64 encode for transmission
	encoded := base64.StdEncoding.EncodeToString(pngData)

	// Kitty graphics protocol escape sequence
	// Format: \033_Ga=T,f=100,s=WIDTH,v=HEIGHT,c=COLS,r=ROWS;\<base64>\033\\
	numChunks := (len(encoded) + 4094) / 4095 // Split into 4095-char chunks

	var output strings.Builder
	chunkSize := 4095

	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(encoded) {
			end = len(encoded)
		}
		chunk := encoded[start:end]

		// Transmission mode T, format PNG f=100
		output.WriteString("\033_Ga=T,f=100,s=")
		output.WriteString(fmt.Sprintf("%d,v=%d", resized.Bounds().Dx(), resized.Bounds().Dy()))

		if i == numChunks-1 {
			output.WriteString(",m=0") // Last chunk
		} else {
			output.WriteString(",m=1") // More chunks coming
		}

		output.WriteString(";")
		output.WriteString(chunk)
		output.WriteString("\033\\")
		if i < numChunks-1 {
			output.WriteString("\n")
		}
	}

	// Add alt text below image
	if altText != "" {
		output.WriteString(fmt.Sprintf("\n[%s]", altText))
	}

	return output.String()
}

// renderAsITerm2 renders image using iTerm2 inline images
// Encodes image as PNG and sends via iTerm2 escape sequence
func (r *ImageRenderer) renderAsITerm2(img image.Image, altText string) string {
	// Resize image to fit terminal
	scaledWidth, scaledHeight := r.CalculateScaledSize(img.Bounds().Dx(), img.Bounds().Dy())
	resized := r.resizeImage(img, scaledWidth, scaledHeight)

	// Encode as PNG
	pngData, err := r.encodeImageAsPNG(resized)
	if err != nil {
		return r.renderAsText(img, altText)
	}

	// Base64 encode for transmission
	encoded := base64.StdEncoding.EncodeToString(pngData)

	// iTerm2 inline image protocol escape sequence
	// Format: \033]1337;File=name=NAME;size=SIZE;width=WIDTHpx;height=HEIGHTpx;inline=1:BASE64\a
	output := fmt.Sprintf(
		"\033]1337;File=width=%dpx;height=%dpx;inline=1:%s\a",
		scaledWidth, scaledHeight, encoded,
	)

	// Add alt text below image
	if altText != "" {
		output += fmt.Sprintf("\n[%s]", altText)
	}

	return output
}

// renderAsSIXEL renders image using SIXEL graphics
// Simplified SIXEL implementation with color reduction
func (r *ImageRenderer) renderAsSIXEL(img image.Image, altText string) string {
	// Resize image to fit terminal
	scaledWidth, scaledHeight := r.CalculateScaledSize(img.Bounds().Dx(), img.Bounds().Dy())
	resized := r.resizeImage(img, scaledWidth, scaledHeight)

	// Start SIXEL sequence
	output := "\033Pq"

	// Convert image to SIXEL (simplified - 16 color palette)
	// For production, would use full color palette
	sixelData := r.imageToSIXEL(resized)
	output += sixelData

	// End SIXEL sequence
	output += "\033\\"

	// Add alt text below image
	if altText != "" {
		output += fmt.Sprintf("\n[%s]", altText)
	}

	return output
}

// renderAsHalfblocks renders image using Unicode halfblock characters
// Uses Unicode block elements to represent 2x2 pixel groups
func (r *ImageRenderer) renderAsHalfblocks(img image.Image, altText string) string {
	// Resize image to fit terminal (half height for aspect ratio with monospace)
	scaledWidth, scaledHeight := r.CalculateScaledSize(img.Bounds().Dx(), img.Bounds().Dy())
	scaledHeight = scaledHeight / 2 // Monospace chars are ~2x tall
	if scaledHeight < 1 {
		scaledHeight = 1
	}
	resized := r.resizeImage(img, scaledWidth, scaledHeight)

	// Convert to halfblock representation
	return r.imageToHalfblocks(resized, altText)
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

// Helper Methods

// encodeImageAsPNG encodes an image to PNG bytes
func (r *ImageRenderer) encodeImageAsPNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// resizeImage resizes an image to the specified dimensions
func (r *ImageRenderer) resizeImage(src image.Image, width, height int) image.Image {
	if width <= 0 || height <= 0 {
		return src
	}

	// Create new image
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	// Simple nearest-neighbor resize
	// For better quality, would use more sophisticated algorithm
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Map destination pixel to source pixel
			srcX := (x * srcWidth) / width
			srcY := (y * srcHeight) / height
			dst.Set(x, y, src.At(srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY))
		}
	}

	return dst
}

// imageToSIXEL converts an image to SIXEL format (simplified)
func (r *ImageRenderer) imageToSIXEL(img image.Image) string {
	// SIXEL format: each character represents a 6-bit vertical strip
	// This is a simplified implementation using 16 colors
	// Full implementation would support 256+ colors

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Pad height to multiple of 6
	paddedHeight := ((height + 5) / 6) * 6

	var output strings.Builder

	// Process image in 6-pixel-high strips
	for stripY := 0; stripY < paddedHeight; stripY += 6 {
		for x := 0; x < width; x++ {
			// Create 6-bit value for this vertical strip
			bits := 0
			for bit := 0; bit < 6; bit++ {
				y := stripY + bit
				if y < height {
					// Sample pixel and determine if it should be set
					r32, g32, b32, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
					brightness := (r32 + g32 + b32) / 3
					if brightness > 32768 { // > 50% brightness
						bits |= (1 << uint(bit))
					}
				}
			}

			// SIXEL character is '?' + bits
			if bits >= 0 && bits < 64 {
				output.WriteByte(byte(bits + 63))
			}
		}
		// End of line
		if stripY+6 < paddedHeight {
			output.WriteString("-")
		}
	}

	return output.String()
}

// imageToHalfblocks converts an image to Unicode halfblock characters
func (r *ImageRenderer) imageToHalfblocks(img image.Image, altText string) string {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Halfblock characters: top/bottom combinations
	halfblocks := []string{
		" ",      // 0000 - empty
		"▄",      // 0001 - bottom
		"▀",      // 0010 - top
		"█",      // 0011 - full
	}

	var output strings.Builder

	for y := 0; y < height; y += 2 {
		for x := 0; x < width; x++ {
			// Get top and bottom pixels
			topBrightness := r.getPixelBrightness(img, bounds, x, y)
			botBrightness := 0
			if y+1 < height {
				botBrightness = r.getPixelBrightness(img, bounds, x, y+1)
			}

			// Determine which halfblock to use
			topSet := topBrightness > 127
			botSet := botBrightness > 127
			index := 0
			if botSet {
				index |= 1
			}
			if topSet {
				index |= 2
			}

			output.WriteString(halfblocks[index])
		}
		output.WriteString("\n")
	}

	// Add alt text if provided
	if altText != "" {
		output.WriteString(fmt.Sprintf("[%s]\n", altText))
	}

	return output.String()
}

// getPixelBrightness returns brightness of pixel (0-255)
func (r *ImageRenderer) getPixelBrightness(img image.Image, bounds image.Rectangle, x, y int) int {
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return 0
	}

	r32, g32, b32, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
	// RGBA returns 16-bit values, convert to 8-bit
	r8 := int(r32 >> 8)
	g8 := int(g32 >> 8)
	b8 := int(b32 >> 8)

	// Calculate brightness using luminance formula
	brightness := (77*r8 + 150*g8 + 29*b8) / 256
	if brightness > 255 {
		brightness = 255
	}
	return brightness
}
