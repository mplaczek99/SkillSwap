package repositories_test

import (
    "testing"

    "github.com/mplaczek99/SkillSwap/models"
    "github.com/mplaczek99/SkillSwap/repositories"
)

func TestInsertUser(t *testing.T) {
    userToInsert := &models.User{
        Email:    "newuser@example.com",
        Password: "hashed_password",
        Name:     "New User",
    }

    insertedUser, err := repositories.InsertUser(userToInsert)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    if insertedUser.ID == 0 {
        t.Errorf("expected a dummy ID to be set, got ID = %d", insertedUser.ID)
    }

    if insertedUser.Email != userToInsert.Email {
        t.Errorf("expected email %s, got %s", userToInsert.Email, insertedUser.Email)
    }
}

func TestGetUserByEmail(t *testing.T) {
    // Known email
    knownEmail := "test@example.com"
    user, err := repositories.GetUserByEmail(knownEmail)
    if err != nil {
        t.Fatalf("expected no error for email %s, got %v", knownEmail, err)
    }
    if user.Email != knownEmail {
        t.Errorf("expected email %s, got %s", knownEmail, user.Email)
    }

    // Unknown email
    unknownEmail := "nobody@example.com"
    _, err = repositories.GetUserByEmail(unknownEmail)
    if err == nil {
        t.Errorf("expected an error for unknown email %s, but got none", unknownEmail)
    }

    // Empty email
    emptyEmail := ""
    _, err = repositories.GetUserByEmail(emptyEmail)
    if err == nil {
        t.Errorf("expected error for empty email, got nil")
    }
}

