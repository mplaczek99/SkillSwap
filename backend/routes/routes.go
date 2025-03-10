package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/controllers"
	"github.com/mplaczek99/SkillSwap/middleware"
	"net/http"
)

func SetupRoutes(router *gin.Engine, authController *controllers.AuthController) {
	// Serve static files at the router level, not inside the API group
	router.StaticFS("/uploads", http.Dir("./uploads"))

	// Create an API group for all API routes.
	api := router.Group("/api")
	{
		// Auth routes.
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Search endpoint.
		api.GET("/search", controllers.Search)

		// Protected endpoints.
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/protected", func(ctx *gin.Context) {
				ctx.JSON(200, gin.H{"message": "You are authenticated"})
			})

			// Video upload endpoint.
			protected.POST("/videos/upload", controllers.VideoUpload)
			protected.GET("/videos", controllers.GetVideosList) // Add this new endpoint

			// New schedule endpoints.
			protected.POST("/schedule", controllers.CreateSchedule)
			protected.GET("/schedule", controllers.GetSchedules)
		}

		// Admin endpoints.
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			admin.GET("/dashboard", func(ctx *gin.Context) {
				ctx.JSON(200, gin.H{"message": "Welcome Admin"})
			})
		}
	}
}
