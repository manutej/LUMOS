package pdf

import (
	"testing"
)

// TestNewImagePageCache creates a new cache
func TestNewImagePageCache(t *testing.T) {
	cache := NewImagePageCache(5)

	if cache == nil {
		t.Error("NewImagePageCache should not return nil")
	}

	stats := cache.Stats()
	if stats.CachedPages != 0 {
		t.Errorf("New cache should have 0 cached pages, got %d", stats.CachedPages)
	}
	if stats.MaxPages != 5 {
		t.Errorf("MaxPages should be 5, got %d", stats.MaxPages)
	}
}

// TestImagePageCache_Put stores images
func TestImagePageCache_Put(t *testing.T) {
	cache := NewImagePageCache(5)

	images := []PageImage{
		{Index: 0, Format: "PNG"},
		{Index: 1, Format: "JPEG"},
	}

	cache.Put(1, images)

	stats := cache.Stats()
	if stats.CachedPages != 1 {
		t.Errorf("After Put, should have 1 cached page, got %d", stats.CachedPages)
	}
}

// TestImagePageCache_Get retrieves images
func TestImagePageCache_Get(t *testing.T) {
	cache := NewImagePageCache(5)

	images := []PageImage{
		{Index: 0, Format: "PNG", Title: "Image1"},
		{Index: 1, Format: "JPEG", Title: "Image2"},
	}

	cache.Put(1, images)

	retrieved, exists := cache.Get(1)
	if !exists {
		t.Error("Get should return exists=true for cached page")
	}

	if len(retrieved) != 2 {
		t.Errorf("Retrieved images should have length 2, got %d", len(retrieved))
	}

	if retrieved[0].Title != "Image1" {
		t.Errorf("First image title should be 'Image1', got %s", retrieved[0].Title)
	}
}

// TestImagePageCache_Get_NotCached handles missing pages
func TestImagePageCache_Get_NotCached(t *testing.T) {
	cache := NewImagePageCache(5)

	_, exists := cache.Get(999)
	if exists {
		t.Error("Get should return exists=false for uncached page")
	}
}

// TestImagePageCache_Clear empties cache
func TestImagePageCache_Clear(t *testing.T) {
	cache := NewImagePageCache(5)

	images := []PageImage{{Index: 0, Format: "PNG"}}
	cache.Put(1, images)
	cache.Put(2, images)

	stats := cache.Stats()
	if stats.CachedPages != 2 {
		t.Errorf("Before clear, should have 2 cached pages, got %d", stats.CachedPages)
	}

	cache.Clear()

	stats = cache.Stats()
	if stats.CachedPages != 0 {
		t.Errorf("After clear, should have 0 cached pages, got %d", stats.CachedPages)
	}

	_, exists := cache.Get(1)
	if exists {
		t.Error("After clear, Get should return exists=false")
	}
}

// TestImagePageCache_LRU_Eviction tests LRU eviction policy
func TestImagePageCache_LRU_Eviction(t *testing.T) {
	cache := NewImagePageCache(3) // Small cache for testing

	images1 := []PageImage{{Index: 0, Title: "Page1"}}
	images2 := []PageImage{{Index: 1, Title: "Page2"}}
	images3 := []PageImage{{Index: 2, Title: "Page3"}}
	images4 := []PageImage{{Index: 3, Title: "Page4"}}

	// Fill cache
	cache.Put(1, images1)
	cache.Put(2, images2)
	cache.Put(3, images3)

	stats := cache.Stats()
	if stats.CachedPages != 3 {
		t.Errorf("Cache should have 3 pages, got %d", stats.CachedPages)
	}

	// Add 4th page, should evict oldest (page 1)
	cache.Put(4, images4)

	stats = cache.Stats()
	if stats.CachedPages != 3 {
		t.Errorf("Cache should still have 3 pages, got %d", stats.CachedPages)
	}

	// Page 1 should be evicted
	_, exists := cache.Get(1)
	if exists {
		t.Error("Page 1 should be evicted from cache")
	}

	// Other pages should still exist
	_, exists = cache.Get(2)
	if !exists {
		t.Error("Page 2 should still be in cache")
	}
	_, exists = cache.Get(3)
	if !exists {
		t.Error("Page 3 should still be in cache")
	}
	_, exists = cache.Get(4)
	if !exists {
		t.Error("Page 4 should be in cache")
	}
}

// TestImagePageCache_Stats returns correct stats
func TestImagePageCache_Stats(t *testing.T) {
	cache := NewImagePageCache(10)

	images := []PageImage{{Index: 0, Format: "PNG"}}
	cache.Put(1, images)
	cache.Put(2, images)
	cache.Put(3, images)

	stats := cache.Stats()

	if stats.CachedPages != 3 {
		t.Errorf("Stats.CachedPages should be 3, got %d", stats.CachedPages)
	}
	if stats.MaxPages != 10 {
		t.Errorf("Stats.MaxPages should be 10, got %d", stats.MaxPages)
	}
}

// TestImagePageCache_Update existing entry
func TestImagePageCache_Update(t *testing.T) {
	cache := NewImagePageCache(5)

	images1 := []PageImage{{Index: 0, Title: "V1"}}
	images2 := []PageImage{{Index: 0, Title: "V2"}, {Index: 1, Title: "V2b"}}

	cache.Put(1, images1)
	cache.Put(1, images2) // Update with new images

	retrieved, exists := cache.Get(1)
	if !exists {
		t.Error("Page should still be in cache after update")
	}

	if len(retrieved) != 2 {
		t.Errorf("Updated images should have length 2, got %d", len(retrieved))
	}

	if retrieved[0].Title != "V2" {
		t.Errorf("First image should be updated to 'V2', got %s", retrieved[0].Title)
	}
}

// TestImagePageCache_EmptyCache handles empty images
func TestImagePageCache_EmptyCache(t *testing.T) {
	cache := NewImagePageCache(5)

	emptyImages := []PageImage{}
	cache.Put(1, emptyImages)

	retrieved, exists := cache.Get(1)
	if !exists {
		t.Error("Empty image list should still be cached")
	}

	if len(retrieved) != 0 {
		t.Errorf("Retrieved empty list should have length 0, got %d", len(retrieved))
	}
}

// BenchmarkImagePageCache_Put benchmarks cache insertion
func BenchmarkImagePageCache_Put(b *testing.B) {
	cache := NewImagePageCache(100)
	images := []PageImage{
		{Index: 0, Format: "PNG"},
		{Index: 1, Format: "JPEG"},
		{Index: 2, Format: "PNG"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Put(i%100, images)
	}
}

// BenchmarkImagePageCache_Get benchmarks cache retrieval
func BenchmarkImagePageCache_Get(b *testing.B) {
	cache := NewImagePageCache(100)
	images := []PageImage{
		{Index: 0, Format: "PNG"},
		{Index: 1, Format: "JPEG"},
	}

	// Pre-populate cache
	for i := 0; i < 100; i++ {
		cache.Put(i, images)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(i % 100)
	}
}

// BenchmarkImagePageCache_Stats benchmarks stats retrieval
func BenchmarkImagePageCache_Stats(b *testing.B) {
	cache := NewImagePageCache(100)
	images := []PageImage{{Index: 0, Format: "PNG"}}

	for i := 0; i < 50; i++ {
		cache.Put(i, images)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cache.Stats()
	}
}
