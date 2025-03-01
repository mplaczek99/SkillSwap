package utils

import (
	"github.com/gin-gonic/gin"
)

// JSONError sends a JSON-formatted error response and aborts the context.
func JSONError(c *gin.Context, status int, errMsg string) {
	c.JSON(status, gin.H{"error": errMsg})
	c.Abort()
}
