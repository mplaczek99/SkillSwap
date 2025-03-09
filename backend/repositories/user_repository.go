package repositories

import (
	"github.com/mplaczek99/SkillSwap/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SearchUsers searches for users by name or email containing the search term
func (r *UserRepository) SearchUsers(searchTerm string) ([]models.User, error) {
	var users []models.User

	// Use ILIKE for case-insensitive search in PostgreSQL
	// Or you can use LOWER() function with LIKE for more database compatibility
	err := r.DB.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?",
		"%"+searchTerm+"%", "%"+searchTerm+"%").Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByID gets a user by their ID
func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
