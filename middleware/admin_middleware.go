package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware ensures that only users with an "Admin" role can access the route.
func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists || role != "Admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden, admin access required"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
