package controllers_test

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/mplaczek99/SkillSwap/controllers"
)

func TestRegisterUser(t *testing.T) {
    router := gin.Default()
    router.POST("/register", controllers.RegisterUser)

    tests := []struct {
        name           string
        requestBody    string
        wantStatusCode int
        wantContains   string
    }{
        {
            name:           "Valid User",
            requestBody:    `{"name":"John Doe","email":"john@example.com","password":"pass1234"}`,
            wantStatusCode: http.StatusCreated,
            wantContains:   `"email":"john@example.com"`,
        },
        {
            name:           "Missing Password",
            requestBody:    `{"name":"John Doe","email":"john@example.com"}`,
            wantStatusCode: http.StatusBadRequest,
            wantContains:   `"error"`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(tt.requestBody)))
            req.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            if w.Code != tt.wantStatusCode {
                t.Errorf("[%s] expected status %d, got %d", tt.name, tt.wantStatusCode, w.Code)
            }
            if tt.wantContains != "" && !bytes.Contains(w.Body.Bytes(), []byte(tt.wantContains)) {
                t.Errorf("[%s] response body does not contain %q. Body: %s", tt.name, tt.wantContains, w.Body.String())
            }
        })
    }
}

func TestGetUser(t *testing.T) {
    router := gin.Default()
    // We'll define /user/:id for test
    router.GET("/user/:id", controllers.GetUser)

    tests := []struct {
        name           string
        url            string
        wantStatusCode int
        wantContains   string
    }{
        {
            name:           "Valid User ID",
            url:            "/user/1",
            wantStatusCode: http.StatusOK,
            wantContains:   `"email":"test@example.com"`,
        },
        {
            name:           "Invalid User ID",
            url:            "/user/999",
            wantStatusCode: http.StatusNotFound,
            wantContains:   `"error":"User not found"`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, _ := http.NewRequest("GET", tt.url, nil)
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            if w.Code != tt.wantStatusCode {
                t.Errorf("[%s] expected status %d, got %d", tt.name, tt.wantStatusCode, w.Code)
            }
            if tt.wantContains != "" && !bytes.Contains(w.Body.Bytes(), []byte(tt.wantContains)) {
                t.Errorf("[%s] body does not contain %q. Body: %s", tt.name, tt.wantContains, w.Body.String())
            }
        })
    }
}

