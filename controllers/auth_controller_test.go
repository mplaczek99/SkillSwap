package controllers_test

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/mplaczek99/SkillSwap/controllers"
)

func TestLogin(t *testing.T) {
    router := gin.Default()
    router.POST("/login", controllers.Login)

    tests := []struct {
        name           string
        requestBody    string
        wantStatusCode int
        wantContains   string
    }{
        {
            name:           "Valid Credentials",
            requestBody:    `{"email":"test@example.com","password":"somepassword"}`,
            wantStatusCode: http.StatusOK, // or http.StatusUnauthorized if dummy logic fails
            wantContains:   `"token"`,
        },
        {
            name:           "Missing Fields",
            requestBody:    `{"email":""}`,
            wantStatusCode: http.StatusBadRequest,
            wantContains:   `"error"`,
        },
        {
            name:           "Invalid Email Format",
            requestBody:    `{"email":"notanemail","password":"secret"}`,
            wantStatusCode: http.StatusBadRequest,
            wantContains:   `"error"`,
        },
        {
            name:           "User not found",
            requestBody:    `{"email":"nope@example.com","password":"secret"}`,
            wantStatusCode: http.StatusUnauthorized,
            wantContains:   `"error"`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(tt.requestBody)))
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
