package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/middleware"
)

func TestAdminMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Admin Role Allowed", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			// Simulate auth middleware setting role
			c.Set("role", "Admin")
			c.Next()
		})
		router.Use(middleware.AdminMiddleware())

		router.GET("/admin", func(c *gin.Context) {
			c.String(http.StatusOK, "Admin access granted")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for Admin role, got %d", w.Code)
		}
	})

	t.Run("Non-Admin Role Forbidden", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			// Simulate auth middleware setting role
			c.Set("role", "User")
			c.Next()
		})
		router.Use(middleware.AdminMiddleware())

		router.GET("/admin", func(c *gin.Context) {
			c.String(http.StatusOK, "Admin access granted")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 for non-Admin role, got %d", w.Code)
		}
	})

	t.Run("Missing Role", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.AdminMiddleware())

		router.GET("/admin", func(c *gin.Context) {
			c.String(http.StatusOK, "Admin access granted")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 for missing role, got %d", w.Code)
		}
	})
}
