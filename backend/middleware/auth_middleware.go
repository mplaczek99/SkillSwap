package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/utils"
)

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

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.Error("Token validation failed: " + err.Error())
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		// Set user details in context.
		ctx.Set("user_id", claims.UserID)
		ctx.Set("role", claims.Role)
		ctx.Set("email", claims.Email)
		ctx.Header("X-User-Email", claims.Email)

		ctx.Next()
	}
}
