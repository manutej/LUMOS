package ui

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

// createTestImage creates a simple test image with specified dimensions
func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Create a gradient pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Gradient from black to white
			brightness := uint8((x * 255) / width)
			c := color.RGBA{brightness, brightness, brightness, 255}
			img.SetRGBA(x, y, c)
		}
	}

	return img
}

// createColorTestImage creates an image with different colors
func createColorTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	quadWidth := width / 2
	quadHeight := height / 2

	// Top-left: Red
	for y := 0; y < quadHeight; y++ {
		for x := 0; x < quadWidth; x++ {
			img.SetRGBA(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	// Top-right: Green
	for y := 0; y < quadHeight; y++ {
		for x := quadWidth; x < width; x++ {
			img.SetRGBA(x, y, color.RGBA{0, 255, 0, 255})
		}
	}

	// Bottom-left: Blue
	for y := quadHeight; y < height; y++ {
		for x := 0; x < quadWidth; x++ {
			img.SetRGBA(x, y, color.RGBA{0, 0, 255, 255})
		}
	}

	// Bottom-right: White
	for y := quadHeight; y < height; y++ {
		for x := quadWidth; x < width; x++ {
			img.SetRGBA(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	return img
}

// TestNewImageRenderer creates a new image renderer
func TestNewImageRenderer(t *testing.T) {
	cfg := ImageRenderConfig{
		Mode:      ImageRenderingEnabled,
		Format:    "halfblock",
		MaxWidth:  80,
		MaxHeight: 24,
	}

	renderer := NewImageRenderer(cfg)
	if renderer == nil {
		t.Error("NewImageRenderer should return non-nil renderer")
	}

	if renderer.config.Mode != ImageRenderingEnabled {
		t.Errorf("Renderer mode should be enabled, got %v", renderer.config.Mode)
	}
}

// TestRenderImage_Disabled tests rendering with disabled mode
func TestRenderImage_Disabled(t *testing.T) {
	cfg := ImageRenderConfig{
		Mode: ImageRenderingDisabled,
	}

	renderer := NewImageRenderer(cfg)
	img := createTestImage(10, 10)

	result := renderer.RenderImage(img, "test")
	if result != "" {
		t.Errorf("Disabled rendering should return empty string, got %q", result)
	}
}

// TestRenderImage_TextMode tests rendering in text mode
func TestRenderImage_TextMode(t *testing.T) {
	cfg := ImageRenderConfig{
		Mode:      ImageRenderingText,
		MaxWidth:  80,
		MaxHeight: 24,
	}

	renderer := NewImageRenderer(cfg)
	img := createTestImage(100, 100)

	result := renderer.RenderImage(img, "test image")
	if !strings.Contains(result, "[Image: test image]") {
		t.Errorf("Text rendering should contain image label, got %q", result)
	}
	if !strings.Contains(result, "━") {
		t.Error("Text rendering should contain border characters")
	}
}

// TestRenderImage_TextMode_NoAltText tests text rendering without alt text
func TestRenderImage_TextMode_NoAltText(t *testing.T) {
	cfg := ImageRenderConfig{
		Mode: ImageRenderingText,
	}

	renderer := NewImageRenderer(cfg)
	img := createTestImage(50, 50)

	result := renderer.RenderImage(img, "")
	if !strings.Contains(result, "50 × 50") {
		t.Errorf("Should show dimensions, got %q", result)
	}
}

// TestRenderAsText verifies text rendering format
func TestRenderAsText(t *testing.T) {
	cfg := ImageRenderConfig{
		Mode: ImageRenderingText,
	}
	renderer := NewImageRenderer(cfg)
	img := createTestImage(64, 32)

	result := renderer.renderAsText(img, "test")
	lines := strings.Split(result, "\n")

	if len(lines) < 3 {
		t.Errorf("Text output should have multiple lines, got %d", len(lines))
	}
}

// TestCalculateScaledSize tests image scaling calculation
func TestCalculateScaledSize(t *testing.T) {
	cfg := ImageRenderConfig{
		MaxWidth:  80,
		MaxHeight: 24,
	}
	renderer := NewImageRenderer(cfg)

	tests := []struct {
		srcWidth, srcHeight int
		maxWidth, maxHeight int
		name                string
	}{
		{100, 100, 80, 24, "square image"},
		{200, 100, 80, 24, "wide image"},
		{100, 200, 80, 24, "tall image"},
		{10, 10, 80, 24, "small image"},
		{1, 1, 80, 24, "tiny image"},
	}

	for _, tt := range tests {
		renderer.config.MaxWidth = tt.maxWidth
		renderer.config.MaxHeight = tt.maxHeight

		width, height := renderer.CalculateScaledSize(tt.srcWidth, tt.srcHeight)

		if width <= 0 || height <= 0 {
			t.Errorf("%s: Scaled size should be positive, got %dx%d", tt.name, width, height)
		}

		if width > tt.maxWidth {
			t.Errorf("%s: Scaled width %d exceeds max %d", tt.name, width, tt.maxWidth)
		}

		if height > tt.maxHeight {
			t.Errorf("%s: Scaled height %d exceeds max %d", tt.name, height, tt.maxHeight)
		}
	}
}

// TestCalculateScaledSize_MinimumSize ensures minimum dimensions
func TestCalculateScaledSize_MinimumSize(t *testing.T) {
	cfg := ImageRenderConfig{
		MaxWidth:  5,
		MaxHeight: 2,
	}
	renderer := NewImageRenderer(cfg)

	width, height := renderer.CalculateScaledSize(100, 100)

	if width < 5 {
		t.Errorf("Scaled width should be at least 5, got %d", width)
	}

	if height < 2 {
		t.Errorf("Scaled height should be at least 2, got %d", height)
	}
}

// TestResizeImage tests image resizing
func TestResizeImage(t *testing.T) {
	renderer := NewImageRenderer(ImageRenderConfig{})
	original := createTestImage(100, 100)

	resized := renderer.resizeImage(original, 50, 50)

	if resized.Bounds().Dx() != 50 {
		t.Errorf("Resized width should be 50, got %d", resized.Bounds().Dx())
	}
	if resized.Bounds().Dy() != 50 {
		t.Errorf("Resized height should be 50, got %d", resized.Bounds().Dy())
	}
}

// TestResizeImage_ExactDimensions verifies exact resizing
func TestResizeImage_ExactDimensions(t *testing.T) {
	renderer := NewImageRenderer(ImageRenderConfig{})
	original := createTestImage(200, 100) // 2:1 aspect ratio

	// resizeImage resizes to exact dimensions (scaling without maintaining aspect)
	// To maintain aspect ratio, use CalculateScaledSize first
	resized := renderer.resizeImage(original, 100, 100)

	// Function should resize to exact requested dimensions
	expectedWidth := 100
	expectedHeight := 100

	actualWidth := resized.Bounds().Dx()
	actualHeight := resized.Bounds().Dy()

	if actualWidth != expectedWidth || actualHeight != expectedHeight {
		t.Errorf("Exact resize failed: expected %dx%d, got %dx%d",
			expectedWidth, expectedHeight, actualWidth, actualHeight)
	}
}

// TestResizeImage_InvalidDimensions handles invalid sizes
func TestResizeImage_InvalidDimensions(t *testing.T) {
	renderer := NewImageRenderer(ImageRenderConfig{})
	original := createTestImage(100, 100)

	// Zero dimensions should return original or fallback
	resized := renderer.resizeImage(original, 0, 0)
	if resized == nil {
		t.Error("Resizing with zero dimensions should not return nil")
	}
}

// TestRenderAsHalfblocks tests halfblock rendering
func TestRenderAsHalfblocks(t *testing.T) {
	cfg := ImageRenderConfig{
		Mode:      ImageRenderingEnabled,
		Format:    "halfblock",
		MaxWidth:  20,
		MaxHeight: 10,
	}
	renderer := NewImageRenderer(cfg)
	img := createTestImage(20, 10)

	result := renderer.renderAsHalfblocks(img, "test")

	if result == "" {
		t.Error("Halfblock rendering should not be empty")
	}

	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) == 0 {
		t.Error("Halfblock output should have lines")
	}

	// Should contain some halfblock characters or spaces
	// Last line may contain [test] label, so exclude from halfblock check
	imageLines := lines
	if len(lines) > 0 && strings.HasPrefix(lines[len(lines)-1], "[") {
		imageLines = lines[:len(lines)-1]
	}

	if len(imageLines) == 0 {
		t.Error("Should have at least one line of image data")
	}

	// First line should contain halfblock characters or spaces
	for _, line := range imageLines {
		if line != "" {
			hasValidChars := strings.ContainsAny(line, " ▀▄█")
			if !hasValidChars {
				t.Errorf("Line contains unexpected characters: %q", line)
			}
		}
	}
}

// TestGetPixelBrightness tests brightness calculation
func TestGetPixelBrightness(t *testing.T) {
	renderer := NewImageRenderer(ImageRenderConfig{})

	// Create image with known colors
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	// Black pixel
	img.SetRGBA(0, 0, color.RGBA{0, 0, 0, 255})
	// White pixel
	img.SetRGBA(1, 0, color.RGBA{255, 255, 255, 255})
	// Gray pixel
	img.SetRGBA(2, 0, color.RGBA{128, 128, 128, 255})

	bounds := img.Bounds()

	blackBrightness := renderer.getPixelBrightness(img, bounds, 0, 0)
	whiteBrightness := renderer.getPixelBrightness(img, bounds, 1, 0)
	grayBrightness := renderer.getPixelBrightness(img, bounds, 2, 0)

	if blackBrightness >= 50 {
		t.Errorf("Black brightness should be low, got %d", blackBrightness)
	}

	if whiteBrightness <= 200 {
		t.Errorf("White brightness should be high, got %d", whiteBrightness)
	}

	if grayBrightness < 50 || grayBrightness > 200 {
		t.Errorf("Gray brightness should be mid-range, got %d", grayBrightness)
	}
}

// TestGetPixelBrightness_OutOfBounds handles out of bounds access
func TestGetPixelBrightness_OutOfBounds(t *testing.T) {
	renderer := NewImageRenderer(ImageRenderConfig{})
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	bounds := img.Bounds()

	// Out of bounds should return 0
	brightness := renderer.getPixelBrightness(img, bounds, 20, 20)
	if brightness != 0 {
		t.Errorf("Out of bounds brightness should be 0, got %d", brightness)
	}
}

// TestRenderAsKitty tests Kitty protocol rendering
func TestRenderAsKitty(t *testing.T) {
	cfg := ImageRenderConfig{
		Mode:      ImageRenderingEnabled,
		Format:    "kitty",
		MaxWidth:  80,
		MaxHeight: 24,
	}
	renderer := NewImageRenderer(cfg)
	img := createTestImage(50, 50)

	result := renderer.renderAsKitty(img, "test")

	if result == "" {
		t.Error("Kitty rendering should not be empty")
	}

	// Should contain Kitty escape sequence
	if !strings.Contains(result, "\033_Ga=T") {
		t.Error("Kitty output should contain escape sequence")
	}

	if strings.Contains(result, "[Image: test]") {
		// If fallback to text, that's also acceptable
		t.Log("Kitty rendering fell back to text (acceptable for testing)")
	}
}

// TestRenderAsITerm2 tests iTerm2 protocol rendering
func TestRenderAsITerm2(t *testing.T) {
	cfg := ImageRenderConfig{
		Mode:      ImageRenderingEnabled,
		Format:    "iterm2",
		MaxWidth:  80,
		MaxHeight: 24,
	}
	renderer := NewImageRenderer(cfg)
	img := createTestImage(50, 50)

	result := renderer.renderAsITerm2(img, "test")

	if result == "" {
		t.Error("iTerm2 rendering should not be empty")
	}

	// Should contain iTerm2 escape sequence or fallback to text
	if !strings.Contains(result, "\033]1337") && !strings.Contains(result, "[Image") {
		t.Error("iTerm2 output should contain escape sequence or text fallback")
	}
}

// TestRenderAsSIXEL tests SIXEL protocol rendering
func TestRenderAsSIXEL(t *testing.T) {
	cfg := ImageRenderConfig{
		Mode:      ImageRenderingEnabled,
		Format:    "sixel",
		MaxWidth:  80,
		MaxHeight: 24,
	}
	renderer := NewImageRenderer(cfg)
	img := createTestImage(50, 50)

	result := renderer.renderAsSIXEL(img, "test")

	if result == "" {
		t.Error("SIXEL rendering should not be empty")
	}

	// Should contain SIXEL escape sequence
	if !strings.Contains(result, "\033Pq") {
		t.Error("SIXEL output should contain escape sequence")
	}
}

// TestFormatLabel tests image label formatting
func TestFormatLabel(t *testing.T) {
	renderer := NewImageRenderer(ImageRenderConfig{})

	label1 := renderer.FormatLabel("my_image.jpg", 640, 480)
	if !strings.Contains(label1, "my_image.jpg") || !strings.Contains(label1, "640x480") {
		t.Errorf("Label formatting failed: %q", label1)
	}

	label2 := renderer.FormatLabel("", 800, 600)
	if !strings.Contains(label2, "800x600") {
		t.Errorf("Empty title label failed: %q", label2)
	}
}

// TestSupportedFormats returns supported formats
func TestSupportedFormats(t *testing.T) {
	renderer := NewImageRenderer(ImageRenderConfig{})
	formats := renderer.SupportedFormats()

	if len(formats) == 0 {
		t.Error("SupportedFormats should return at least one format")
	}

	// Formats should be non-empty strings
	for _, format := range formats {
		if format == "" {
			t.Error("Format should not be empty string")
		}
	}

	// Should contain some rendering methods (format names vary by terminal)
	if len(formats) < 1 {
		t.Error("Should have at least one format in fallback chain")
	}
}

// BenchmarkRenderAsHalfblocks benchmarks halfblock rendering
func BenchmarkRenderAsHalfblocks(b *testing.B) {
	cfg := ImageRenderConfig{
		Mode:      ImageRenderingEnabled,
		Format:    "halfblock",
		MaxWidth:  80,
		MaxHeight: 24,
	}
	renderer := NewImageRenderer(cfg)
	img := createTestImage(80, 24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer.renderAsHalfblocks(img, "test")
	}
}

// BenchmarkRenderAsText benchmarks text rendering
func BenchmarkRenderAsText(b *testing.B) {
	renderer := NewImageRenderer(ImageRenderConfig{})
	img := createTestImage(100, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer.renderAsText(img, "test image")
	}
}

// BenchmarkResizeImage benchmarks image resizing
func BenchmarkResizeImage(b *testing.B) {
	renderer := NewImageRenderer(ImageRenderConfig{})
	img := createTestImage(400, 400)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer.resizeImage(img, 80, 24)
	}
}

// BenchmarkCalculateScaledSize benchmarks size calculation
func BenchmarkCalculateScaledSize(b *testing.B) {
	cfg := ImageRenderConfig{
		MaxWidth:  80,
		MaxHeight: 24,
	}
	renderer := NewImageRenderer(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer.CalculateScaledSize(640, 480)
	}
}
