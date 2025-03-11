package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/utils"
)

// Define a simple token cache structure
type tokenCacheEntry struct {
	claims    *utils.Claims
	expiresAt time.Time
}

// Create a global token cache with mutex for thread safety
var (
	tokenCache      = make(map[string]tokenCacheEntry)
	tokenCacheMutex sync.RWMutex
	// How long to cache tokens - shorter than token expiry time (e.g., 10 minutes)
	cacheTTL = 10 * time.Minute
)

// cleanExpiredTokens periodically removes expired tokens from the cache
func cleanExpiredTokens() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		removeExpiredTokens()
	}
}

// removeExpiredTokens is a helper function to clean up expired cache entries
func removeExpiredTokens() {
	now := time.Now()
	tokenCacheMutex.Lock()
	defer tokenCacheMutex.Unlock()

	for token, entry := range tokenCache {
		if now.After(entry.expiresAt) {
			delete(tokenCache, token)
		}
	}
}

// Initialize the cache cleaner when the package is imported
func init() {
	go cleanExpiredTokens()
}

// AuthMiddleware validates the Authorization header, extracts the token (handling the "Bearer" prefix),
// verifies it, and then sets user details in the context.
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			ctx.Abort()
			return
		}

		// Check for "Bearer " prefix.
		const bearerPrefix = "Bearer "
		tokenString := authHeader
		if strings.HasPrefix(authHeader, bearerPrefix) {
			tokenString = strings.TrimPrefix(authHeader, bearerPrefix)
		}

		// First, check if we have this token in the cache
		tokenCacheMutex.RLock()
		cacheEntry, found := tokenCache[tokenString]
		tokenCacheMutex.RUnlock()

		var claims *utils.Claims
		var err error

		if found && time.Now().Before(cacheEntry.expiresAt) {
			// Cache hit - use cached claims
			claims = cacheEntry.claims
		} else {
			// Cache miss or expired entry - validate token
			claims, err = utils.ValidateToken(tokenString)
			if err != nil {
				utils.Error("Token validation failed: " + err.Error())
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				ctx.Abort()
				return
			}

			// Cache the successfully validated token
			// Only cache if the token expiry is further than our cache TTL
			tokenExpiry := time.Unix(int64(claims.ExpiresAt.Time.Unix()), 0)
			cacheExpiryTime := time.Now().Add(cacheTTL)

			// Use the earlier of token expiry or cache TTL as cache entry expiry
			entryExpiresAt := cacheExpiryTime
			if tokenExpiry.Before(cacheExpiryTime) {
				entryExpiresAt = tokenExpiry
			}

			tokenCacheMutex.Lock()
			tokenCache[tokenString] = tokenCacheEntry{
				claims:    claims,
				expiresAt: entryExpiresAt,
			}
			tokenCacheMutex.Unlock()
		}

		// Set user details in context.
		ctx.Set("user_id", claims.UserID)
		ctx.Set("role", claims.Role)
		ctx.Set("email", claims.Email)
		ctx.Header("X-User-Email", claims.Email)

		ctx.Next()
	}
}
