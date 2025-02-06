package repositories

import (
	"github.com/mplaczek99/SkillSwap/models"
)

// InsertSkill assigns a dummy ID to the skill.
func InsertSkill(skill *models.Skill) (*models.Skill, error) {
	skill.ID = 1
	return skill, nil
}

// GetAllSkills returns a slice of dummy skills.
func GetAllSkills() ([]models.Skill, error) {
	dummySkill := models.Skill{
		ID:          1,
		Name:        "Dummy Skill",
		Description: "This is a dummy skill",
		UserID:      1,
	}
	return []models.Skill{dummySkill}, nil
}
