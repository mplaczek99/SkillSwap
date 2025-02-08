package utils_test

import (
    "testing"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "github.com/mplaczek99/SkillSwap/config"
    "github.com/mplaczek99/SkillSwap/utils"
)

func TestGenerateJWTAndValidateJWT(t *testing.T) {
    config.LoadConfig()

    tokenString, err := utils.GenerateJWT(123, "test@example.com")
    if err != nil {
        t.Fatalf("GenerateJWT failed: %v", err)
    }
    if tokenString == "" {
        t.Errorf("expected a token string, got empty")
    }

    claims, err := utils.ValidateJWT(tokenString)
    if err != nil {
        t.Fatalf("ValidateJWT failed: %v", err)
    }
    if claims == nil {
        t.Fatal("expected valid claims, got nil")
    }
    if claims.UserID != 123 || claims.Email != "test@example.com" {
        t.Errorf("claims do not match expected values, got %+v", claims)
    }
}

func TestValidateJWT_Expired(t *testing.T) {
    config.LoadConfig()

    // Create an expired token manually
    expiredClaims := utils.JWTClaims{
        UserID: 999,
        Email:  "expired@example.com",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
        },
    }
    expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
    tokenStr, err := expiredToken.SignedString([]byte(config.AppConfig.JWTSecret))
    if err != nil {
        t.Fatalf("failed to create expired token: %v", err)
    }

    _, err = utils.ValidateJWT(tokenStr)
    if err == nil {
        t.Fatal("expected error for expired token, got nil")
    }
}

