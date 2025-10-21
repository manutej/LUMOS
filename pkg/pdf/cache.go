package pdf

import (
	"container/list"
	"sync"
)

// LRUCache implements a Least Recently Used cache for PDF pages
type LRUCache struct {
	maxSize   int
	size      int
	cache     map[int]*CacheNode
	list      *list.List
	mu        sync.RWMutex
	hits      int
	misses    int
}

// CacheNode represents a node in the cache
type CacheNode struct {
	pageNum int
	data    string
	element *list.Element
}

// NewLRUCache creates a new LRU cache
func NewLRUCache(maxSize int) *LRUCache {
	if maxSize <= 0 {
		maxSize = 5 // Default to 5 pages
	}

	return &LRUCache{
		maxSize: maxSize,
		size:    0,
		cache:   make(map[int]*CacheNode),
		list:    list.New(),
	}
}

// Get retrieves a page from cache
func (c *LRUCache) Get(pageNum int) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.cache[pageNum]
	if !exists {
		c.misses++
		return "", false
	}

	// Move to front (most recently used)
	c.list.MoveToFront(node.element)
	c.hits++
	return node.data, true
}

// Put stores a page in cache
func (c *LRUCache) Put(pageNum int, data string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If already exists, update and move to front
	if node, exists := c.cache[pageNum]; exists {
		node.data = data
		c.list.MoveToFront(node.element)
		return
	}

	// If cache is full, remove least recently used
	if c.size >= c.maxSize {
		c.evict()
	}

	// Add new node to front
	element := c.list.PushFront(&CacheNode{
		pageNum: pageNum,
		data:    data,
	})

	c.cache[pageNum] = &CacheNode{
		pageNum: pageNum,
		data:    data,
		element: element,
	}
	c.size++
}

// evict removes the least recently used item
func (c *LRUCache) evict() {
	back := c.list.Back()
	if back == nil {
		return
	}

	node := back.Value.(*CacheNode)
	c.list.Remove(back)
	delete(c.cache, node.pageNum)
	c.size--
}

// Clear clears all items from cache
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[int]*CacheNode)
	c.list = list.New()
	c.size = 0
}

// Stats returns cache statistics
func (c *LRUCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return CacheStats{
		CachedPages: c.size,
		MaxSize:     c.maxSize,
	}
}

// HitRate returns the cache hit rate
func (c *LRUCache) HitRate() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	if total == 0 {
		return 0
	}

	return float64(c.hits) / float64(total)
}

// Reset resets cache statistics
func (c *LRUCache) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.hits = 0
	c.misses = 0
}

// GetStats returns detailed cache statistics
func (c *LRUCache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	hitRate := 0.0
	if total > 0 {
		hitRate = float64(c.hits) / float64(total)
	}

	return map[string]interface{}{
		"size":       c.size,
		"max_size":   c.maxSize,
		"hits":       c.hits,
		"misses":     c.misses,
		"total":      total,
		"hit_rate":   hitRate,
	}
}
