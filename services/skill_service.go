package services

import (
	"errors"
	"time"

	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/repositories"
)

// CreateSkill creates a new skill record.
func CreateSkill(skill *models.Skill) (*models.Skill, error) {
	if skill.Name == "" {
		return nil, errors.New("skill name is required")
	}
	skill.CreatedAt = time.Now()
	return repositories.InsertSkill(skill)
}

// GetAllSkills returns all skill records.
func GetAllSkills() ([]models.Skill, error) {
	return repositories.GetAllSkills()
}
