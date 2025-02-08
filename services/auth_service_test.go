package services_test

import (
    "testing"

    "github.com/mplaczek99/SkillSwap/services"
)

func TestAuthenticateUser(t *testing.T) {
    // Known user "test@example.com" with dummy hashed password
    token, err := services.AuthenticateUser("test@example.com", "somepassword")
    if err != nil {
        t.Errorf("expected no error for valid user, got %v", err)
    }
    if token == "" {
        t.Errorf("expected a token, got empty string")
    }

    // Unknown user
    _, err = services.AuthenticateUser("nope@example.com", "password")
    if err == nil {
        t.Errorf("expected error for unknown user, got nil")
    }
}

