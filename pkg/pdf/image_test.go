package pdf

import (
	"image"
	"image/color"
	"testing"
)

// TestDefaultImageExtractionOptions verifies default settings
func TestDefaultImageExtractionOptions(t *testing.T) {
	opts := DefaultImageExtractionOptions()

	if opts.MinWidth != 10 {
		t.Errorf("MinWidth should be 10, got %f", opts.MinWidth)
	}
	if opts.MinHeight != 10 {
		t.Errorf("MinHeight should be 10, got %f", opts.MinHeight)
	}
	if !opts.OnlyInlineImages {
		t.Error("OnlyInlineImages should be true")
	}
	if !opts.PreserveDPI {
		t.Error("PreserveDPI should be true")
	}
	if opts.MaxImagesPerPage != 100 {
		t.Errorf("MaxImagesPerPage should be 100, got %d", opts.MaxImagesPerPage)
	}
}

// TestPageImage_Structure verifies PageImage fields
func TestPageImage_Structure(t *testing.T) {
	// Create a simple test image
	img := &image.RGBA{
		Rect: image.Rect(0, 0, 100, 100),
		Pix:  make([]uint8, 100*100*4),
	}

	pageImg := PageImage{
		Data:   img,
		X:      10.5,
		Y:      20.5,
		Width:  100.0,
		Height: 100.0,
		Index:  0,
		Format: "PNG",
		Title:  "Test Image",
	}

	if pageImg.X != 10.5 {
		t.Errorf("X mismatch: got %f", pageImg.X)
	}
	if pageImg.Y != 20.5 {
		t.Errorf("Y mismatch: got %f", pageImg.Y)
	}
	if pageImg.Width != 100.0 {
		t.Errorf("Width mismatch: got %f", pageImg.Width)
	}
	if pageImg.Height != 100.0 {
		t.Errorf("Height mismatch: got %f", pageImg.Height)
	}
	if pageImg.Format != "PNG" {
		t.Errorf("Format mismatch: got %s", pageImg.Format)
	}
	if pageImg.Title != "Test Image" {
		t.Errorf("Title mismatch: got %s", pageImg.Title)
	}
}

// TestRenderSize_Structure verifies RenderSize fields
func TestRenderSize_Structure(t *testing.T) {
	size := RenderSize{
		Width:  80,
		Height: 24,
	}

	if size.Width != 80 {
		t.Errorf("Width should be 80, got %d", size.Width)
	}
	if size.Height != 24 {
		t.Errorf("Height should be 24, got %d", size.Height)
	}
}

// TestImageInfo_Structure verifies ImageInfo fields
func TestImageInfo_Structure(t *testing.T) {
	info := ImageInfo{
		PageNum:   1,
		ImageCount: 3,
		TotalSize: 150000,
		ImageDims: []image.Rectangle{
			image.Rect(0, 0, 800, 600),
			image.Rect(0, 0, 300, 200),
			image.Rect(0, 0, 150, 100),
		},
		HasLargeImg: true,
		HasSmallImg: true,
	}

	if info.PageNum != 1 {
		t.Errorf("PageNum should be 1, got %d", info.PageNum)
	}
	if info.ImageCount != 3 {
		t.Errorf("ImageCount should be 3, got %d", info.ImageCount)
	}
	if info.TotalSize != 150000 {
		t.Errorf("TotalSize should be 150000, got %d", info.TotalSize)
	}
	if !info.HasLargeImg {
		t.Error("HasLargeImg should be true")
	}
	if !info.HasSmallImg {
		t.Error("HasSmallImg should be true")
	}
}

// BenchmarkDefaultImageExtractionOptions benchmarks option creation
func BenchmarkDefaultImageExtractionOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = DefaultImageExtractionOptions()
	}
}

// BenchmarkPageImage_Creation benchmarks PageImage creation
func BenchmarkPageImage_Creation(b *testing.B) {
	img := &image.RGBA{
		Rect: image.Rect(0, 0, 100, 100),
		Pix:  make([]uint8, 100*100*4),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PageImage{
			Data:   img,
			X:      10.0,
			Y:      20.0,
			Width:  100.0,
			Height: 100.0,
			Index:  i,
			Format: "PNG",
			Title:  "Benchmark Image",
		}
	}
}

// Helper function to create a test image
func createTestImage(width, height int) image.Image {
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)

	// Fill with a simple gradient
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(128)
			img.SetRGBA(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}

// TestCreateTestImage verifies helper function
func TestCreateTestImage(t *testing.T) {
	img := createTestImage(100, 100)

	if img.Bounds().Dx() != 100 {
		t.Errorf("Image width should be 100, got %d", img.Bounds().Dx())
	}
	if img.Bounds().Dy() != 100 {
		t.Errorf("Image height should be 100, got %d", img.Bounds().Dy())
	}
}

// BenchmarkCreateTestImage benchmarks test image creation
func BenchmarkCreateTestImage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = createTestImage(800, 600)
	}
}
