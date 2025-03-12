package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/controllers"
	"github.com/mplaczek99/SkillSwap/middleware"
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
			protected.GET("/videos", controllers.GetVideosList)

			// New schedule endpoints.
			protected.POST("/schedule", controllers.CreateSchedule)
			protected.GET("/schedule", controllers.GetSchedules)

			// Transactions endpoint
			protected.GET("/transactions", controllers.GetTransactions)

			// Job endpoints
			protected.GET("/jobs", controllers.GetJobs)
			protected.GET("/jobs/:id", controllers.GetJob)
			protected.POST("/jobs", controllers.CreateJob)
			protected.PUT("/jobs/:id", controllers.UpdateJob)
			protected.DELETE("/jobs/:id", controllers.DeleteJob)
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
