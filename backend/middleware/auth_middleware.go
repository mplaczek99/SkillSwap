package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/utils"
)

// Simple cache for storing validated tokens
type TokenCache struct {
	mu    sync.RWMutex
	cache map[string]*CacheItem
}

type CacheItem struct {
	claims    *utils.Claims
	expiresAt time.Time
}

// Global token cache
var tokenCache = &TokenCache{
	cache: make(map[string]*CacheItem),
}

// Get retrieves a token from the cache
// Get retrieves a token from the cache
func (c *TokenCache) Get(token string) (*utils.Claims, bool) {
	c.mu.RLock()
	item, found := c.cache[token]
	c.mu.RUnlock()

	if !found {
		return nil, false
	}

	// Check if expired, but don't delete here
	if time.Now().After(item.expiresAt) {
		// Let the periodic cleanup handle deletion
		return nil, false
	}

	return item.claims, true
}

// Set adds a token to the cache
func (c *TokenCache) Set(token string, claims *utils.Claims, expiry time.Time) {
	c.mu.Lock()
	c.cache[token] = &CacheItem{
		claims:    claims,
		expiresAt: expiry,
	}
	c.mu.Unlock()
}

// CleanExpired removes expired tokens more efficiently
func (c *TokenCache) CleanExpired() {
	// First phase: identify expired tokens with read lock
	var expiredTokens []string

	c.mu.RLock()
	now := time.Now()
	for token, item := range c.cache {
		if now.After(item.expiresAt) {
			expiredTokens = append(expiredTokens, token)
		}
	}
	c.mu.RUnlock()

	// Second phase: remove expired tokens with write lock
	// Only if there are tokens to remove
	if len(expiredTokens) > 0 {
		c.mu.Lock()
		// Get fresh timestamp for accurate expiration check
		now = time.Now()
		for _, token := range expiredTokens {
			// Double-check expiration again under write lock
			// in case another goroutine updated the token
			if item, found := c.cache[token]; found && now.After(item.expiresAt) {
				delete(c.cache, token)
			}
		}
		c.mu.Unlock()
	}
}

// Start periodic cleanup
func init() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			tokenCache.CleanExpired()
		}
	}()
}

// Extract token from Authorization header
func extractToken(authHeader string) string {
	if len(authHeader) == 0 {
		return ""
	}

	// Check if the header starts with "Bearer "
	const bearerPrefix = "Bearer "
	if len(authHeader) > len(bearerPrefix) && strings.HasPrefix(authHeader, bearerPrefix) {
		return authHeader[len(bearerPrefix):]
	}

	// Only accept properly formatted Bearer tokens
	return ""
}

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from header
		tokenString := extractToken(c.GetHeader("Authorization"))
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		// Check token cache first
		claims, found := tokenCache.Get(tokenString)

		if !found {
			// Validate the token
			var err error
			claims, err = utils.ValidateToken(tokenString)
			if err != nil {
				if err == utils.ErrExpiredToken {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "token has expired"})
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				}
				c.Abort()
				return
			}

			// Calculate expiry time for cache
			tokenExpiry := time.Unix(claims.ExpiresAt.Time.Unix(), 0)
			cacheExpiry := time.Now().Add(10 * time.Minute) // 10 minute cache
			if tokenExpiry.Before(cacheExpiry) {
				cacheExpiry = tokenExpiry
			}

			// Cache the validated token
			tokenCache.Set(tokenString, claims, cacheExpiry)
		}

		// Set user details in context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)

		c.Next()
	}
}
