// File: ./routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/controllers"
	"github.com/mplaczek99/SkillSwap/middleware"
)

func SetupRoutes(router *gin.Engine, authController *controllers.AuthController) {
	// Create an API group for all API routes
	api := router.Group("/api")
	{
		// Auth routes group
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Protected endpoints for any authenticated user
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/protected", func(ctx *gin.Context) {
				ctx.JSON(200, gin.H{"message": "You are authenticated"})
			})
		}

		// Admin-only endpoints
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			admin.GET("/dashboard", func(ctx *gin.Context) {
				ctx.JSON(200, gin.H{"message": "Welcome Admin"})
			})
		}
	}
}
