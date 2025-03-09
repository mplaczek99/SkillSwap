package repositories

import (
	"strings"

	"github.com/mplaczek99/SkillSwap/models"
	"gorm.io/gorm"
)

// SkillRepository handles database operations for skills
type SkillRepository struct {
	DB *gorm.DB
}

// NewSkillRepository creates a new instance of SkillRepository
func NewSkillRepository(db *gorm.DB) *SkillRepository {
	return &SkillRepository{DB: db}
}

// InsertSkill saves a new skill to the database
func (r *SkillRepository) InsertSkill(skill *models.Skill) (*models.Skill, error) {
	if err := r.DB.Create(skill).Error; err != nil {
		return nil, err
	}
	return skill, nil
}

// GetAllSkills returns all skills from the database
func (r *SkillRepository) GetAllSkills() ([]models.Skill, error) {
	var skills []models.Skill
	if err := r.DB.Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}

// SearchSkills searches for skills by name or description containing the search term
func (r *SkillRepository) SearchSkills(searchTerm string) ([]models.Skill, error) {
	var skills []models.Skill

	// Use ILIKE for case-insensitive search in PostgreSQL
	// Or you can use LOWER() function with LIKE for more database compatibility
	err := r.DB.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?",
		"%"+searchTerm+"%", "%"+searchTerm+"%").Find(&skills).Error

	if err != nil {
		return nil, err
	}

	return skills, nil
}

// For backward compatibility with existing code, provide these as standalone functions

func InsertSkill(skill *models.Skill) (*models.Skill, error) {
	// In a real app with a database connection, this would be implemented properly
	// For now, this maintains compatibility with existing code
	skill.ID = 1 // This is a temporary implementation
	return skill, nil
}

func GetAllSkills() ([]models.Skill, error) {
	// In a real app with a database connection, this would query the database
	// For now, this maintains compatibility with existing code
	dummySkill := models.Skill{
		ID:          1,
		Name:        "Dummy Skill",
		Description: "This is a dummy skill",
		UserID:      1,
	}
	return []models.Skill{dummySkill}, nil
}

// SearchSkills as a standalone function for compatibility
func SearchSkills(searchTerm string) ([]models.Skill, error) {
	// This is a temporary implementation
	// In a real app, you would get a DB connection and perform the query

	// For now, return dummy skills that match the search term
	skills := []models.Skill{
		{
			ID:          1,
			Name:        "Programming",
			Description: "Learn to code with various languages",
			UserID:      1,
		},
		{
			ID:          2,
			Name:        "Language Learning",
			Description: "Master new languages quickly",
			UserID:      2,
		},
		{
			ID:          3,
			Name:        "Cooking",
			Description: "Learn to cook delicious meals",
			UserID:      1,
		},
	}

	var results []models.Skill
	for _, skill := range skills {
		if strings.Contains(strings.ToLower(skill.Name), searchTerm) ||
			strings.Contains(strings.ToLower(skill.Description), searchTerm) {
			results = append(results, skill)
		}
	}

	return results, nil
}
