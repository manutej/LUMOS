package pdf

import (
	"sync"
)

// ImagePageCache implements LRU caching for page images
// Similar to LRUCache but stores slices of PageImage
type ImagePageCache struct {
	mu    sync.RWMutex
	cache map[int][]PageImage // pageNum -> images
	order []int               // LRU order (oldest first)
	maxSize int
}

// NewImagePageCache creates a new image page cache with max size
func NewImagePageCache(maxPages int) *ImagePageCache {
	if maxPages <= 0 {
		maxPages = 5 // Default
	}
	return &ImagePageCache{
		cache:   make(map[int][]PageImage),
		order:   make([]int, 0, maxPages),
		maxSize: maxPages,
	}
}

// Get retrieves images for a page from cache
// Returns (images, true) if cached, (nil, false) otherwise
func (c *ImagePageCache) Get(pageNum int) ([]PageImage, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	images, exists := c.cache[pageNum]
	return images, exists
}

// Put stores images for a page in cache
// Evicts oldest page if cache is full
func (c *ImagePageCache) Put(pageNum int, images []PageImage) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If page already cached, update and move to end
	if _, exists := c.cache[pageNum]; exists {
		c.cache[pageNum] = images
		// Move to end of order
		for i, p := range c.order {
			if p == pageNum {
				c.order = append(c.order[:i], c.order[i+1:]...)
				break
			}
		}
		c.order = append(c.order, pageNum)
		return
	}

	// If cache full, evict oldest
	if len(c.cache) >= c.maxSize {
		oldest := c.order[0]
		delete(c.cache, oldest)
		c.order = c.order[1:]
	}

	// Add new page
	c.cache[pageNum] = images
	c.order = append(c.order, pageNum)
}

// Clear removes all cached images
func (c *ImagePageCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[int][]PageImage)
	c.order = make([]int, 0, c.maxSize)
}

// Stats returns cache statistics
type ImageCacheStats struct {
	CachedPages int
	MaxPages    int
}

// Stats returns current cache statistics
func (c *ImagePageCache) Stats() ImageCacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return ImageCacheStats{
		CachedPages: len(c.cache),
		MaxPages:    c.maxSize,
	}
}
