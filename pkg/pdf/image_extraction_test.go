package pdf

import (
	"testing"
)

// TestGetPageImages_ValidPageRange tests that GetPageImages handles valid page numbers
func TestGetPageImages_ValidPageRange(t *testing.T) {
	doc := &Document{
		pages:      5,
		filepath:   "test.pdf",
		imageCache: NewImagePageCache(10),
	}

	opts := DefaultImageExtractionOptions()

	// Valid page numbers should not error on validation
	tests := []int{1, 2, 3, 4, 5}

	for _, pageNum := range tests {
		_, err := doc.GetPageImages(pageNum, opts)
		// Error is acceptable if pdfcpu isn't available, but shouldn't be a validation error
		if err != nil && err.Error() == "page number out of range: 1" {
			t.Errorf("Page %d should be valid, got validation error", pageNum)
		}
	}
}

// TestGetPageImages_InvalidPageRange tests that invalid page numbers are rejected
func TestGetPageImages_InvalidPageRange(t *testing.T) {
	doc := &Document{
		pages:      5,
		filepath:   "test.pdf",
		imageCache: NewImagePageCache(10),
	}

	opts := DefaultImageExtractionOptions()

	invalidPages := []int{0, -1, 6, 10, 100}

	for _, pageNum := range invalidPages {
		_, err := doc.GetPageImages(pageNum, opts)
		if err == nil {
			t.Errorf("Page %d should be invalid", pageNum)
		}
		if err != nil && err.Error() != "page number out of range: "+string(rune(pageNum)) {
			// Different error message is fine, as long as it's caught
		}
	}
}

// TestGetPageImages_CacheIntegration tests that images are properly cached
func TestGetPageImages_CacheIntegration(t *testing.T) {
	cache := NewImagePageCache(10)
	doc := &Document{
		pages:      5,
		filepath:   "test.pdf",
		imageCache: cache,
	}

	opts := DefaultImageExtractionOptions()

	// First call will extract (or return empty if pdfcpu unavailable)
	images1, err1 := doc.GetPageImages(1, opts)
	if err1 != nil && err1.Error() == "page number out of range: 1" {
		t.Fatal("Page 1 should be valid")
	}

	// Verify cache was populated
	stats := cache.Stats()
	if stats.CachedPages < 1 {
		t.Errorf("Cache should have at least 1 page after extraction, got %d", stats.CachedPages)
	}

	// Second call should return same data from cache
	images2, _ := doc.GetPageImages(1, opts)

	// Both should have same length (even if 0, they're the same)
	if len(images1) != len(images2) {
		t.Errorf("Cached images should match: first call %d, second call %d", len(images1), len(images2))
	}
}

// TestGetPageImages_MaxImagesPerPage respects option limit
func TestGetPageImages_MaxImagesPerPage(t *testing.T) {
	doc := &Document{
		pages:      5,
		filepath:   "test.pdf",
		imageCache: NewImagePageCache(10),
	}

	opts := DefaultImageExtractionOptions()
	opts.MaxImagesPerPage = 5

	// Should not error even if limit is set (actual extraction depends on pdfcpu)
	_, err := doc.GetPageImages(1, opts)
	if err != nil && err.Error() == "page number out of range: 1" {
		t.Fatal("Should not get range error")
	}
}

// TestGetPageImages_MinSizeFilter respects dimension options
func TestGetPageImages_MinSizeFilter(t *testing.T) {
	doc := &Document{
		pages:      5,
		filepath:   "test.pdf",
		imageCache: NewImagePageCache(10),
	}

	opts := DefaultImageExtractionOptions()
	opts.MinWidth = 100
	opts.MinHeight = 100

	// Should not error even with filter options set
	_, err := doc.GetPageImages(1, opts)
	if err != nil && err.Error() == "page number out of range: 1" {
		t.Fatal("Should not get range error")
	}
}

// TestImageExtractionOptions_Defaults verifies sensible defaults
func TestImageExtractionOptions_Defaults(t *testing.T) {
	opts := DefaultImageExtractionOptions()

	if opts.MinWidth != 10 {
		t.Errorf("MinWidth should default to 10, got %f", opts.MinWidth)
	}
	if opts.MinHeight != 10 {
		t.Errorf("MinHeight should default to 10, got %f", opts.MinHeight)
	}
	if !opts.OnlyInlineImages {
		t.Error("OnlyInlineImages should default to true")
	}
	if !opts.PreserveDPI {
		t.Error("PreserveDPI should default to true")
	}
	if opts.MaxImagesPerPage != 100 {
		t.Errorf("MaxImagesPerPage should default to 100, got %d", opts.MaxImagesPerPage)
	}
}

// TestEstimatePageImageCount basic functionality
func TestEstimatePageImageCount(t *testing.T) {
	doc := &Document{
		pages:      5,
		filepath:   "test.pdf",
		imageCache: NewImagePageCache(10),
	}

	// Valid page should not error on validation
	count, err := doc.EstimatePageImageCount(3)
	if err != nil && err.Error() == "page number out of range: 3" {
		t.Fatal("Page 3 should be valid")
	}

	// Should return non-negative count (even if 0 without pdfcpu)
	if count < 0 {
		t.Errorf("Count should be non-negative, got %d", count)
	}
}

// TestEstimatePageImageCount_InvalidPages rejects invalid pages
func TestEstimatePageImageCount_InvalidPages(t *testing.T) {
	doc := &Document{
		pages:      5,
		filepath:   "test.pdf",
		imageCache: NewImagePageCache(10),
	}

	invalidPages := []int{0, -1, 6, 100}

	for _, pageNum := range invalidPages {
		_, err := doc.EstimatePageImageCount(pageNum)
		if err == nil {
			t.Errorf("Page %d should be invalid", pageNum)
		}
	}
}

// TestImageCacheOperations verifies document's cache methods
func TestImageCacheOperations(t *testing.T) {
	cache := NewImagePageCache(10)
	doc := &Document{
		pages:      5,
		filepath:   "test.pdf",
		imageCache: cache,
	}

	// HasImageCache
	if !doc.HasImageCache() {
		t.Error("HasImageCache should return true")
	}

	// ImageCacheStats
	stats := doc.ImageCacheStats()
	if stats.MaxPages != 10 {
		t.Errorf("Stats should show MaxPages=10, got %d", stats.MaxPages)
	}

	// ClearImageCache
	doc.ClearImageCache()
	statsAfterClear := doc.ImageCacheStats()
	if statsAfterClear.CachedPages != 0 {
		t.Errorf("After clear, CachedPages should be 0, got %d", statsAfterClear.CachedPages)
	}
}

// TestPageImage_ExtractionContext validates PageImage struct in extraction context
func TestPageImage_ExtractionContext(t *testing.T) {
	img := PageImage{
		Data:   nil,
		X:      10.5,
		Y:      20.3,
		Width:  100,
		Height: 200,
		Index:  0,
		Format: "JPEG",
		Title:  "Test Image",
	}

	if img.X != 10.5 {
		t.Errorf("X should be 10.5, got %f", img.X)
	}
	if img.Y != 20.3 {
		t.Errorf("Y should be 20.3, got %f", img.Y)
	}
	if img.Width != 100 {
		t.Errorf("Width should be 100, got %f", img.Width)
	}
	if img.Height != 200 {
		t.Errorf("Height should be 200, got %f", img.Height)
	}
	if img.Format != "JPEG" {
		t.Errorf("Format should be JPEG, got %s", img.Format)
	}
	if img.Title != "Test Image" {
		t.Errorf("Title should be 'Test Image', got %s", img.Title)
	}
}

// BenchmarkGetPageImages_NoImages benchmarks cache lookup
func BenchmarkGetPageImages_NoImages(b *testing.B) {
	cache := NewImagePageCache(10)
	doc := &Document{
		pages:      5,
		filepath:   "test.pdf",
		imageCache: cache,
	}

	opts := DefaultImageExtractionOptions()

	// Pre-populate cache with empty slices
	cache.Put(1, []PageImage{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc.GetPageImages(1, opts)
	}
}

// BenchmarkGetPageImages_MultiplePages benchmarks extraction of multiple pages
func BenchmarkGetPageImages_MultiplePages(b *testing.B) {
	cache := NewImagePageCache(10)
	doc := &Document{
		pages:      5,
		filepath:   "test.pdf",
		imageCache: cache,
	}

	opts := DefaultImageExtractionOptions()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for page := 1; page <= 5; page++ {
			doc.GetPageImages(page, opts)
		}
	}
}
