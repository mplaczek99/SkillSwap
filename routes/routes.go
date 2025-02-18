package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/controllers"
	"github.com/mplaczek99/SkillSwap/middleware"
)

func SetupRoutes(router *gin.Engine, authController *controllers.AuthController) {
	// Public endpoints
	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)

	// Protected endpoints for any authenticated user
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.GET("/protected", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "You are authenticated"})
	})

	// Admin-only endpoints
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	admin.GET("/dashboard", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Welcome Admin"})
	})
}
