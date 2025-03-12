package middleware

import (
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/utils"
)

// Constants for optimal performance tuning
const (
	// Number of shards should be power of 2, matching CPU cores for optimal distribution
	shardCount         = 16
	shardCountMinusOne = shardCount - 1
	shardMask          = shardCount - 1 // For fast modulo with AND
	baseShardCapacity  = 128            // Per-shard capacity
	cleanupInterval    = 2 * time.Minute
	tokenCacheTTL      = 10 * time.Minute
	bearerPrefix       = "Bearer "
	bearerPrefixLen    = len(bearerPrefix)
	estimatedTokenLen  = 256 // Pre-allocation size for token strings
	maxCleanBatchSize  = 100 // Max items to clean per operation
)

// ShardedCache is a high-performance thread-safe token cache
// using multiple LRU caches with sharding to reduce lock contention
type ShardedCache struct {
	shards [shardCount]*LRUCache
	// Quick stats counters
	hits      atomic.Uint64
	misses    atomic.Uint64
	evictions atomic.Uint64
}

// LRUCache implements a fast, thread-safe LRU cache for a single shard
type LRUCache struct {
	capacity int
	mu       sync.RWMutex
	items    map[string]*lruItem
	head     *lruNode
	tail     *lruNode
	// Pool for nodes to reduce GC pressure
	nodePool sync.Pool
}

type lruItem struct {
	claims    *utils.Claims
	expiresAt int64 // Unix timestamp for faster comparison
	node      *lruNode
}

type lruNode struct {
	key  string
	prev *lruNode
	next *lruNode
}

// NewShardedCache creates a sharded LRU cache optimized for concurrent access
func NewShardedCache() *ShardedCache {
	cache := &ShardedCache{}

	// Calculate optimal shard capacity based on available cores
	cpus := runtime.NumCPU()
	optimalShards := shardCount
	if cpus < shardCount {
		optimalShards = cpus
	}

	// Initialize shards - distribute capacity based on cores
	perShardCapacity := baseShardCapacity
	if optimalShards < shardCount {
		// Redistribute capacity if fewer cores than shards
		perShardCapacity = (baseShardCapacity * shardCount) / optimalShards
	}

	for i := 0; i < shardCount; i++ {
		cache.shards[i] = newLRUCache(perShardCapacity)
	}

	return cache
}

// newLRUCache creates a new LRU cache with the given capacity
func newLRUCache(capacity int) *LRUCache {
	cache := &LRUCache{
		capacity: capacity,
		items:    make(map[string]*lruItem, capacity),
		head:     &lruNode{key: "head"},
		tail:     &lruNode{key: "tail"},
		nodePool: sync.Pool{
			New: func() interface{} {
				return &lruNode{}
			},
		},
	}
	cache.head.next = cache.tail
	cache.tail.prev = cache.head
	return cache
}

// getShard returns the appropriate shard for a given key using a fast modulo operation
// This uses bit masking instead of modulo for better performance
func (c *ShardedCache) getShard(key string) *LRUCache {
	// Fast hash function for strings - FNV-1a variant
	h := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		h ^= uint32(key[i])
		h *= 16777619
	}

	// Fast modulo with bit mask (requires power of 2 shards)
	return c.shards[h&shardMask]
}

// Get retrieves an item from the cache by key, returning the claims, expiration time, and existence
func (c *ShardedCache) Get(key string) (*utils.Claims, time.Time, bool) {
	shard := c.getShard(key)
	claims, expiry, found := shard.Get(key)

	// Update stats
	if found {
		c.hits.Add(1)
	} else {
		c.misses.Add(1)
	}

	return claims, expiry, found
}

// Put adds or updates an item in the cache
func (c *ShardedCache) Put(key string, claims *utils.Claims, expiresAt time.Time) {
	shard := c.getShard(key)
	if evicted := shard.Put(key, claims, expiresAt); evicted {
		c.evictions.Add(1)
	}
}

// CleanAllExpired removes expired entries from all shards
func (c *ShardedCache) CleanAllExpired() int {
	total := 0
	now := time.Now()
	for i := 0; i < shardCount; i++ {
		total += c.shards[i].CleanExpired(now)
	}
	return total
}

// Stats returns cache statistics
func (c *ShardedCache) Stats() map[string]interface{} {
	var totalSize, totalCap int
	for i := 0; i < shardCount; i++ {
		size, cap := c.shards[i].Stats()
		totalSize += size
		totalCap += cap
	}

	hitCount := c.hits.Load()
	missCount := c.misses.Load()
	total := hitCount + missCount
	hitRate := 0.0
	if total > 0 {
		hitRate = float64(hitCount) / float64(total) * 100.0
	}

	return map[string]interface{}{
		"size":      totalSize,
		"capacity":  totalCap,
		"hit_rate":  hitRate,
		"hits":      hitCount,
		"misses":    missCount,
		"evictions": c.evictions.Load(),
	}
}

// Get retrieves an item from the cache by key
func (c *LRUCache) Get(key string) (*utils.Claims, time.Time, bool) {
	// Fast path - try read lock first
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		return nil, time.Time{}, false
	}

	// Check if item has expired - using int64 comparison is faster than time.Time
	now := time.Now().Unix()
	if now > item.expiresAt {
		// Item is expired - remove it under write lock
		c.mu.Lock()
		// Need to check again after acquiring write lock
		if item, stillExists := c.items[key]; stillExists {
			if now > item.expiresAt {
				c.removeItem(key, item.node)
			}
		}
		c.mu.Unlock()
		return nil, time.Time{}, false
	}

	// Update item position in the cache (move to front) - requires write lock
	c.mu.Lock()
	c.moveToFront(item.node)
	c.mu.Unlock()

	// Convert int64 back to time.Time
	expiry := time.Unix(item.expiresAt, 0)
	return item.claims, expiry, true
}

// Put adds an item to the cache, returns true if an item was evicted
func (c *LRUCache) Put(key string, claims *utils.Claims, expiresAt time.Time) bool {
	expiryUnix := expiresAt.Unix()
	evicted := false

	c.mu.Lock()
	defer c.mu.Unlock()

	// If key already exists, update it and move to front
	if item, found := c.items[key]; found {
		item.claims = claims
		item.expiresAt = expiryUnix
		c.moveToFront(item.node)
		return evicted
	}

	// If cache is full, remove least recently used item (tail)
	if len(c.items) >= c.capacity {
		c.removeItem(c.tail.prev.key, c.tail.prev)
		evicted = true
	}

	// Create new node from pool and add to front
	newNode := c.nodePool.Get().(*lruNode)
	newNode.key = key

	newItem := &lruItem{
		claims:    claims,
		expiresAt: expiryUnix,
		node:      newNode,
	}

	c.items[key] = newItem
	c.addToFront(newNode)
	return evicted
}

// CleanExpired removes all expired items from the cache
func (c *LRUCache) CleanExpired(now time.Time) int {
	nowUnix := now.Unix()
	c.mu.Lock()
	defer c.mu.Unlock()

	var removed int
	// We'll check from LRU end (tail) first as they're more likely to be expired
	node := c.tail.prev

	// Limit how many items we check in one pass to avoid long lock times
	checked := 0
	for node != c.head && checked < maxCleanBatchSize {
		prevNode := node.prev // Save before potential removal
		if item, found := c.items[node.key]; found && nowUnix > item.expiresAt {
			c.removeItem(node.key, node)
			removed++
		}
		node = prevNode
		checked++
	}
	return removed
}

// Stats returns size and capacity of this cache
func (c *LRUCache) Stats() (int, int) {
	c.mu.RLock()
	size := len(c.items)
	c.mu.RUnlock()
	return size, c.capacity
}

// addToFront adds a node to the front of the linked list (must be called under lock)
func (c *LRUCache) addToFront(node *lruNode) {
	node.prev = c.head
	node.next = c.head.next
	c.head.next.prev = node
	c.head.next = node
}

// moveToFront moves a node to the front of the linked list (must be called under lock)
func (c *LRUCache) moveToFront(node *lruNode) {
	// Remove from current position
	node.prev.next = node.next
	node.next.prev = node.prev

	// Add to front
	c.addToFront(node)
}

// removeItem removes a node from the linked list and the map (must be called under lock)
func (c *LRUCache) removeItem(key string, node *lruNode) {
	// Remove from the linked list
	node.prev.next = node.next
	node.next.prev = node.prev

	// Remove from the map
	delete(c.items, key)

	// Reset node and return to pool
	node.prev = nil
	node.next = nil
	node.key = ""
	c.nodePool.Put(node)
}

// Global token cache with optimal size and sharding
var (
	tokenCache = NewShardedCache()
)

// startCacheCleanup begins the background cleanup process
func startCacheCleanup() {
	ticker := time.NewTicker(cleanupInterval)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			// This runs every 2 minutes to clean expired entries
			tokenCache.CleanAllExpired()
		}
	}()
}

// Initialize the cache cleaner when the package is imported
func init() {
	// Start cleanup in a separate goroutine
	startCacheCleanup()
}

// extractTokenFromHeader efficiently extracts a JWT token from the Authorization header
// using zero-allocation techniques where possible
func extractTokenFromHeader(authHeader string) string {
	if len(authHeader) == 0 {
		return ""
	}

	// Check if the authorization header has the Bearer prefix
	if len(authHeader) > bearerPrefixLen && authHeader[:bearerPrefixLen] == bearerPrefix {
		// Use safe string slicing instead of reflect.StringHeader manipulation
		return authHeader[bearerPrefixLen:]
	}

	// Return the header as-is if it doesn't have the Bearer prefix
	return authHeader
}

// preflightCheck performs fast validation of token format before expensive crypto operations
// Returns true if the token passes basic format checks
func preflightCheck(token string) bool {
	// Length check
	if len(token) < 10 {
		return false
	}

	// Structure check (3 parts separated by dots)
	dotCount := 0
	hasInvalidChar := false

	// Manual scan is faster than strings.Count or strings.Split
	for i := 0; i < len(token); i++ {
		if token[i] == '.' {
			dotCount++
		} else if token[i] < 32 || token[i] > 126 {
			// ASCII printable range check
			hasInvalidChar = true
			break
		}
	}

	return dotCount == 2 && !hasInvalidChar
}

// AuthMiddleware validates the Authorization header, extracts the token,
// verifies it, and then sets user details in the context.
//
//go:noinline
func AuthMiddleware() gin.HandlerFunc {
	// Pre-allocate common error responses
	missingTokenResponse := gin.H{"error": "missing token"}
	invalidTokenResponse := gin.H{"error": "invalid token"}
	expiredTokenResponse := gin.H{"error": "token has expired"}

	return func(ctx *gin.Context) {
		// Extract token from header - optimized for zero allocations
		tokenString := extractTokenFromHeader(ctx.GetHeader("Authorization"))
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, missingTokenResponse)
			ctx.Abort()
			return
		}

		// Fast preflight check before trying more expensive operations
		if !preflightCheck(tokenString) {
			ctx.JSON(http.StatusUnauthorized, invalidTokenResponse)
			ctx.Abort()
			return
		}

		// Try to get claims from cache first
		claims, _, found := tokenCache.Get(tokenString)

		if !found {
			// Cache miss - validate token
			var err error
			claims, err = utils.ValidateToken(tokenString)
			if err != nil {
				// Use pre-allocated responses for common errors
				if err == utils.ErrExpiredToken {
					ctx.JSON(http.StatusUnauthorized, expiredTokenResponse)
				} else {
					ctx.JSON(http.StatusUnauthorized, invalidTokenResponse)
				}
				ctx.Abort()
				return
			}

			// Calculate optimal cache TTL
			tokenExpiry := time.Unix(claims.ExpiresAt.Time.Unix(), 0)
			cacheExpiry := time.Now().Add(tokenCacheTTL)
			if tokenExpiry.Before(cacheExpiry) {
				cacheExpiry = tokenExpiry
			}

			// Cache the validated token
			tokenCache.Put(tokenString, claims, cacheExpiry)
		}

		// Set user details in context using a batch operation for efficiency
		// This avoids multiple calls to c.Set which each acquire internal locks
		ctxValues := map[string]interface{}{
			"user_id": claims.UserID,
			"role":    claims.Role,
			"email":   claims.Email,
		}

		for k, v := range ctxValues {
			ctx.Set(k, v)
		}

		// Continue with request handling
		ctx.Next()
	}
}

// GetCacheStats returns statistics about the token cache
func GetCacheStats() map[string]interface{} {
	return tokenCache.Stats()
}

// End of file
