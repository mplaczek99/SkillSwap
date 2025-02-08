package services_test

import (
    "testing"

    "github.com/mplaczek99/SkillSwap/models"
    "github.com/mplaczek99/SkillSwap/services"
)

func TestCreateUser_TableDriven(t *testing.T) {
    tests := []struct {
        name          string
        user          models.User
        expectError   bool
        errorContains string
    }{
        {
            name: "Valid user",
            user: models.User{
                Email:    "valid@domain.com",
                Password: "somepassword",
            },
            expectError: false,
        },
        {
            name: "Missing Email",
            user: models.User{
                Password: "somepassword",
            },
            expectError:   true,
            errorContains: "required",
        },
        {
            name: "Missing Password",
            user: models.User{
                Email: "missingpass@domain.com",
            },
            expectError:   true,
            errorContains: "required",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := services.CreateUser(&tt.user)
            if tt.expectError && err == nil {
                t.Errorf("expected an error, got nil")
            }
            if !tt.expectError && err != nil {
                t.Errorf("expected no error, got %v", err)
            }
            if tt.errorContains != "" && err != nil && !containsSubstring(err.Error(), tt.errorContains) {
                t.Errorf("error message does not contain %q. Actual: %v", tt.errorContains, err)
            }
        })
    }
}

func TestGetUserByID(t *testing.T) {
    user, err := services.GetUserByID("1")
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if user == nil || user.Email != "test@example.com" {
        t.Errorf("expected test@example.com, got %v", user)
    }

    user, err = services.GetUserByID("999")
    if err == nil {
        t.Error("expected error for user not found, got nil")
    }
}

func containsSubstring(str, substr string) bool {
    return len(str) >= len(substr) && // quick check
        (func() bool {
            for i := 0; i+len(substr) <= len(str); i++ {
                if str[i:i+len(substr)] == substr {
                    return true
                }
            }
            return false
        })()
}

