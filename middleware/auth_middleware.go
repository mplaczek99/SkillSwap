package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/utils"
)

// AuthMiddleware validates JWT and extracts user information
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			ctx.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		// Set user details in context
		ctx.Set("user_id", claims.UserID)
		ctx.Set("role", claims.Role)
		ctx.Set("email", claims.Email) // Add email to context

		// Add email to response headers (Optional)
		ctx.Header("X-User-Email", claims.Email)

		ctx.Next()
	}
}
