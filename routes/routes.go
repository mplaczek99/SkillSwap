package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/controllers"
	"github.com/mplaczek99/SkillSwap/middleware"
)

// SetupRoutes registers all API endpoints and their middleware.
func SetupRoutes(router *gin.Engine) {
	// Public routes (no authentication required)
	public := router.Group("/api")
	{
		public.POST("/login", controllers.Login)
		public.POST("/register", controllers.RegisterUser)
	}

	// Protected routes (authentication required)
	protected := router.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/user/:id", controllers.GetUser)
		protected.POST("/skills", controllers.AddSkill)
		protected.GET("/skills", controllers.GetSkills)
		// Add additional protected routes as needed
	}
}

