package repositories

import (
    "errors"
    "github.com/mplaczek99/SkillSwap/models"
)

// InsertUser returns the user with a dummy ID.
func InsertUser(user *models.User) (*models.User, error) {
    user.ID = 1
    return user, nil
}

// GetUserByEmail returns a dummy user for a specific email.
func GetUserByEmail(email string) (*models.User, error) {
    // Return a dummy user if the email matches; otherwise, return an error.
    if email == "test@example.com" {
        return &models.User{
            ID:       1,
            Email:    email,
            Password: "$2a$10$dummyhashedpassword", // Dummy hashed password
            Name:     "Dummy User",
        }, nil
    }
    return nil, errors.New("user not found")
}

// GetUserByID returns a dummy user.
func GetUserByID(id string) (*models.User, error) {
    return &models.User{
        ID:    1,
        Email: "test@example.com",
        Name:  "Dummy User",
    }, nil
}

