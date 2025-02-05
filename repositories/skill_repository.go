package repositories

import (
	"github.com/mplaczek99/SkillSwap/config"
	"github.com/mplaczek99/SkillSwap/models"
)

// InsertSkill saves a new skill record to the database.
func InsertSkill(skill *models.Skill) (*models.Skill, error) {
	result := config.DB.Create(skill)
	return skill, result.Error
}

// GetAllSkills retrieves all skill records from the database.
func GetAllSkills() ([]models.Skill, error) {
	var skills []models.Skill
	result := config.DB.Find(&skills)
	return skills, result.Error
}

