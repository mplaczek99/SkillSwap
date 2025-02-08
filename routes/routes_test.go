package routes_test

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/mplaczek99/SkillSwap/routes"
)

func TestPublicRoutes(t *testing.T) {
    router := gin.New()
    routes.SetupRoutes(router)

    // Test /api/login
    reqBody := `{"email":"test@example.com","password":"somepassword"}`
    req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer([]byte(reqBody)))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK && w.Code != http.StatusUnauthorized {
        t.Errorf("expected 200 or 401, got %d", w.Code)
    }

    // Test /api/register
    reqBody = `{"email":"new@ex.com","password":"1234"}`
    req, _ = http.NewRequest("POST", "/api/register", bytes.NewBuffer([]byte(reqBody)))
    req.Header.Set("Content-Type", "application/json")
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusCreated && w.Code != http.StatusBadRequest {
        t.Errorf("expected 201 or 400, got %d", w.Code)
    }
}

func TestProtectedRoutes(t *testing.T) {
    router := gin.New()
    routes.SetupRoutes(router)

    // We need a token to test protected routes properly. For now,
    // send no token and expect a 401.

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/user/1", nil)
    router.ServeHTTP(w, req)
    if w.Code != http.StatusUnauthorized {
        t.Errorf("expected 401 for missing token, got %d", w.Code)
    }

    // Similarly, you can test /api/skills, etc.
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/api/skills", nil)
    router.ServeHTTP(w, req)
    if w.Code != http.StatusUnauthorized {
        t.Errorf("expected 401 for missing token, got %d", w.Code)
    }

    // If you want to test a valid token scenario, you'd generate a token
    // (via your AuthService or a direct call to utils.GenerateJWT) and attach
    // it as "Authorization: Bearer <token>" in the header, then check for success.
}

func contains(body, substr string) bool {
    return strings.Contains(body, substr)
}

