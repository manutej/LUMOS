package pdf

import (
	"fmt"
	"sync"
	"testing"
)

// TestNewLRUCache tests creating a new LRU cache
func TestNewLRUCache(t *testing.T) {
	tests := []struct {
		name        string
		maxSize     int
		wantMaxSize int
	}{
		{
			name:        "normal size",
			maxSize:     10,
			wantMaxSize: 10,
		},
		{
			name:        "zero size - uses default",
			maxSize:     0,
			wantMaxSize: 5,
		},
		{
			name:        "negative size - uses default",
			maxSize:     -5,
			wantMaxSize: 5,
		},
		{
			name:        "large size",
			maxSize:     1000,
			wantMaxSize: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLRUCache(tt.maxSize)

			if cache == nil {
				t.Fatal("NewLRUCache() returned nil")
			}

			if cache.maxSize != tt.wantMaxSize {
				t.Errorf("maxSize = %v, want %v", cache.maxSize, tt.wantMaxSize)
			}

			if cache.size != 0 {
				t.Errorf("initial size = %v, want 0", cache.size)
			}

			if cache.cache == nil {
				t.Error("cache map is nil")
			}

			if cache.list == nil {
				t.Error("list is nil")
			}
		})
	}
}

// TestLRUCache_PutAndGet tests basic put and get operations
func TestLRUCache_PutAndGet(t *testing.T) {
	cache := NewLRUCache(3)

	// Test put and get
	cache.Put(1, "page1")
	data, found := cache.Get(1)

	if !found {
		t.Error("Get(1) not found, expected to find it")
	}

	if data != "page1" {
		t.Errorf("Get(1) = %q, want %q", data, "page1")
	}

	// Test get non-existent
	_, found = cache.Get(999)
	if found {
		t.Error("Get(999) found, expected not to find it")
	}
}

// TestLRUCache_Update tests updating existing entry
func TestLRUCache_Update(t *testing.T) {
	cache := NewLRUCache(3)

	// Put initial value
	cache.Put(1, "original")

	// Update with new value
	cache.Put(1, "updated")

	data, found := cache.Get(1)
	if !found {
		t.Fatal("Get(1) not found")
	}

	if data != "updated" {
		t.Errorf("Get(1) = %q, want %q", data, "updated")
	}

	// Size should still be 1
	if cache.size != 1 {
		t.Errorf("size = %v, want 1", cache.size)
	}
}

// TestLRUCache_Eviction tests LRU eviction
func TestLRUCache_Eviction(t *testing.T) {
	cache := NewLRUCache(3)

	// Fill cache to capacity
	cache.Put(1, "page1")
	cache.Put(2, "page2")
	cache.Put(3, "page3")

	// Verify all are present
	for i := 1; i <= 3; i++ {
		if _, found := cache.Get(i); !found {
			t.Errorf("Get(%d) not found before eviction", i)
		}
	}

	// Add one more - should evict page 1 (least recently used)
	cache.Put(4, "page4")

	// Page 1 should be evicted
	if _, found := cache.Get(1); found {
		t.Error("Page 1 should have been evicted")
	}

	// Pages 2, 3, 4 should be present
	for i := 2; i <= 4; i++ {
		if _, found := cache.Get(i); !found {
			t.Errorf("Get(%d) not found after eviction", i)
		}
	}

	// Size should be at max
	if cache.size != 3 {
		t.Errorf("size = %v, want 3", cache.size)
	}
}

// TestLRUCache_AccessUpdatesRecency tests that accessing an item updates recency
func TestLRUCache_AccessUpdatesRecency(t *testing.T) {
	cache := NewLRUCache(3)

	// Fill cache
	cache.Put(1, "page1")
	cache.Put(2, "page2")
	cache.Put(3, "page3")

	// Access page 1 to make it most recently used
	cache.Get(1)

	// Add page 4 - should evict page 2 (now least recently used)
	cache.Put(4, "page4")

	// Page 1 should still be present (was accessed)
	if _, found := cache.Get(1); !found {
		t.Error("Page 1 should still be present (was recently accessed)")
	}

	// Page 2 should be evicted
	if _, found := cache.Get(2); found {
		t.Error("Page 2 should have been evicted")
	}

	// Pages 3 and 4 should be present
	if _, found := cache.Get(3); !found {
		t.Error("Page 3 should be present")
	}
	if _, found := cache.Get(4); !found {
		t.Error("Page 4 should be present")
	}
}

// TestLRUCache_Clear tests clearing the cache
func TestLRUCache_Clear(t *testing.T) {
	cache := NewLRUCache(5)

	// Add some entries
	for i := 1; i <= 5; i++ {
		cache.Put(i, fmt.Sprintf("page%d", i))
	}

	// Verify entries exist
	if cache.size != 5 {
		t.Errorf("size before clear = %v, want 5", cache.size)
	}

	// Clear cache
	cache.Clear()

	// Verify cache is empty
	if cache.size != 0 {
		t.Errorf("size after clear = %v, want 0", cache.size)
	}

	// Verify entries are gone
	for i := 1; i <= 5; i++ {
		if _, found := cache.Get(i); found {
			t.Errorf("Page %d should not be found after clear", i)
		}
	}

	// Verify we can add new entries after clear
	cache.Put(10, "new page")
	if _, found := cache.Get(10); !found {
		t.Error("Should be able to add entries after clear")
	}
}

// TestLRUCache_Stats tests cache statistics
func TestLRUCache_Stats(t *testing.T) {
	maxSize := 5
	cache := NewLRUCache(maxSize)

	// Empty cache
	stats := cache.Stats()
	if stats.CachedPages != 0 {
		t.Errorf("initial CachedPages = %v, want 0", stats.CachedPages)
	}
	if stats.MaxSize != maxSize {
		t.Errorf("MaxSize = %v, want %v", stats.MaxSize, maxSize)
	}

	// Add some pages
	for i := 1; i <= 3; i++ {
		cache.Put(i, fmt.Sprintf("page%d", i))
	}

	stats = cache.Stats()
	if stats.CachedPages != 3 {
		t.Errorf("CachedPages = %v, want 3", stats.CachedPages)
	}

	// Fill to capacity
	for i := 4; i <= 5; i++ {
		cache.Put(i, fmt.Sprintf("page%d", i))
	}

	stats = cache.Stats()
	if stats.CachedPages != maxSize {
		t.Errorf("CachedPages = %v, want %v", stats.CachedPages, maxSize)
	}
}

// TestLRUCache_HitRate tests cache hit rate calculation
func TestLRUCache_HitRate(t *testing.T) {
	cache := NewLRUCache(5)

	// Initially zero
	if rate := cache.HitRate(); rate != 0 {
		t.Errorf("initial HitRate = %v, want 0", rate)
	}

	// Add pages
	cache.Put(1, "page1")
	cache.Put(2, "page2")

	// 2 hits, 0 misses = 100%
	cache.Get(1)
	cache.Get(2)

	rate := cache.HitRate()
	if rate != 1.0 {
		t.Errorf("HitRate after 2 hits = %v, want 1.0", rate)
	}

	// 2 hits, 1 miss = 66.6%
	cache.Get(999) // miss

	rate = cache.HitRate()
	expected := 2.0 / 3.0
	if rate != expected {
		t.Errorf("HitRate after 2 hits, 1 miss = %v, want %v", rate, expected)
	}
}

// TestLRUCache_Reset tests resetting statistics
func TestLRUCache_Reset(t *testing.T) {
	cache := NewLRUCache(5)

	// Generate some statistics
	cache.Put(1, "page1")
	cache.Get(1)  // hit
	cache.Get(2)  // miss

	// Verify stats exist
	if cache.hits == 0 || cache.misses == 0 {
		t.Fatal("Setup: expected hits and misses")
	}

	// Reset
	cache.Reset()

	// Verify stats reset
	if cache.hits != 0 {
		t.Errorf("hits after reset = %v, want 0", cache.hits)
	}
	if cache.misses != 0 {
		t.Errorf("misses after reset = %v, want 0", cache.misses)
	}

	// Hit rate should be 0
	if rate := cache.HitRate(); rate != 0 {
		t.Errorf("HitRate after reset = %v, want 0", rate)
	}
}

// TestLRUCache_GetStats tests detailed statistics
func TestLRUCache_GetStats(t *testing.T) {
	cache := NewLRUCache(5)

	// Add some entries and generate stats
	cache.Put(1, "page1")
	cache.Put(2, "page2")
	cache.Get(1)  // hit
	cache.Get(2)  // hit
	cache.Get(3)  // miss

	stats := cache.GetStats()

	// Verify all stats fields
	if stats["size"].(int) != 2 {
		t.Errorf("stats[size] = %v, want 2", stats["size"])
	}

	if stats["max_size"].(int) != 5 {
		t.Errorf("stats[max_size] = %v, want 5", stats["max_size"])
	}

	if stats["hits"].(int) != 2 {
		t.Errorf("stats[hits] = %v, want 2", stats["hits"])
	}

	if stats["misses"].(int) != 1 {
		t.Errorf("stats[misses] = %v, want 1", stats["misses"])
	}

	if stats["total"].(int) != 3 {
		t.Errorf("stats[total] = %v, want 3", stats["total"])
	}

	expectedRate := 2.0 / 3.0
	if stats["hit_rate"].(float64) != expectedRate {
		t.Errorf("stats[hit_rate] = %v, want %v", stats["hit_rate"], expectedRate)
	}
}

// TestLRUCache_ThreadSafety tests concurrent access
func TestLRUCache_ThreadSafety(t *testing.T) {
	cache := NewLRUCache(100)
	numGoroutines := 50
	opsPerGoroutine := 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Concurrent puts
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				pageNum := (id * opsPerGoroutine) + j
				cache.Put(pageNum, fmt.Sprintf("page%d", pageNum))
			}
		}(i)
	}

	wg.Wait()

	// Concurrent gets
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				pageNum := (id * opsPerGoroutine) + j
				cache.Get(pageNum)
			}
		}(i)
	}

	wg.Wait()

	// If we reach here without deadlock or panic, test passes
	if cache.size > cache.maxSize {
		t.Errorf("cache size = %v exceeds maxSize = %v", cache.size, cache.maxSize)
	}
}

// TestLRUCache_ConcurrentPutGet tests concurrent put and get
func TestLRUCache_ConcurrentPutGet(t *testing.T) {
	cache := NewLRUCache(50)
	numGoroutines := 20

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)

	// Writers
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				cache.Put(j, fmt.Sprintf("writer%d-page%d", id, j))
			}
		}(i)
	}

	// Readers
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				cache.Get(j)
			}
		}(i)
	}

	wg.Wait()

	// Verify cache is in valid state
	stats := cache.Stats()
	if stats.CachedPages > cache.maxSize {
		t.Errorf("CachedPages = %v exceeds maxSize = %v", stats.CachedPages, cache.maxSize)
	}
}

// TestLRUCache_LargeData tests cache with large data
func TestLRUCache_LargeData(t *testing.T) {
	cache := NewLRUCache(10)

	// Create large page data
	largeData := string(make([]byte, 1024*1024)) // 1MB

	// Add large pages
	for i := 1; i <= 10; i++ {
		cache.Put(i, largeData)
	}

	// Verify all can be retrieved
	for i := 1; i <= 10; i++ {
		data, found := cache.Get(i)
		if !found {
			t.Errorf("Large page %d not found", i)
		}
		if len(data) != len(largeData) {
			t.Errorf("Large page %d data length = %v, want %v", i, len(data), len(largeData))
		}
	}
}

// Benchmark tests
func BenchmarkLRUCache_Put(b *testing.B) {
	cache := NewLRUCache(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Put(i%1000, fmt.Sprintf("page%d", i))
	}
}

func BenchmarkLRUCache_Get_Hit(b *testing.B) {
	cache := NewLRUCache(1000)

	// Populate cache
	for i := 0; i < 1000; i++ {
		cache.Put(i, fmt.Sprintf("page%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(i % 1000)
	}
}

func BenchmarkLRUCache_Get_Miss(b *testing.B) {
	cache := NewLRUCache(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(i + 1000) // Always miss
	}
}

func BenchmarkLRUCache_PutEvict(b *testing.B) {
	cache := NewLRUCache(100)

	// Fill cache first
	for i := 0; i < 100; i++ {
		cache.Put(i, fmt.Sprintf("page%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Put(i+100, fmt.Sprintf("page%d", i+100))
	}
}

func BenchmarkLRUCache_Concurrent(b *testing.B) {
	cache := NewLRUCache(1000)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				cache.Put(i%1000, fmt.Sprintf("page%d", i))
			} else {
				cache.Get(i % 1000)
			}
			i++
		}
	})
}
