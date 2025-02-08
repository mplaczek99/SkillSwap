package middleware_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "github.com/mplaczek99/SkillSwap/config"
    "github.com/mplaczek99/SkillSwap/middleware"
    "github.com/mplaczek99/SkillSwap/utils"
)

// TestJWTAuthMiddleware_NoHeader and TestJWTAuthMiddleware_InvalidFormat already exist.
// Let's add a test for a valid token and an expired token.

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
    // Generate a valid token
    config.LoadConfig()
    token, err := utils.GenerateJWT(123, "testuser@example.com")
    if err != nil {
        t.Fatalf("failed to generate token: %v", err)
    }

    router := gin.Default()
    router.Use(middleware.JWTAuthMiddleware())
    router.GET("/protected", func(c *gin.Context) {
        c.String(http.StatusOK, "Protected Resource")
    })

    req, _ := http.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", w.Code)
    }
}

func TestJWTAuthMiddleware_ExpiredToken(t *testing.T) {
    // Manually create an expired token
    config.LoadConfig()
    claims := utils.JWTClaims{
        UserID: 999,
        Email:  "expired@example.com",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // already expired
            IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
        },
    }

    expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenStr, err := expiredToken.SignedString([]byte(config.AppConfig.JWTSecret))
    if err != nil {
        t.Fatalf("failed to sign token: %v", err)
    }

    router := gin.Default()
    router.Use(middleware.JWTAuthMiddleware())
    router.GET("/protected", func(c *gin.Context) {
        c.String(http.StatusOK, "Protected Resource")
    })

    req, _ := http.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+tokenStr)

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusUnauthorized {
        t.Errorf("expected 401 for expired token, got %d", w.Code)
    }
    if !containsString(w.Body.String(), "Invalid token:") {
        t.Errorf("response should contain 'Invalid token:', got %s", w.Body.String())
    }
}

func containsString(haystack, needle string) bool {
    return len(haystack) >= len(needle) && // quick check
        (func() bool {
            return (haystack[:len(needle)] == needle) || (len(haystack) > len(needle) && containsString(haystack[1:], needle))
        })()
}

