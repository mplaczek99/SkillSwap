package repositories_test

import (
    "testing"

    "github.com/mplaczek99/SkillSwap/models"
    "github.com/mplaczek99/SkillSwap/repositories"
)

func TestInsertSkill(t *testing.T) {
    skill := &models.Skill{
        Name:        "Test Skill",
        Description: "Some description",
        UserID:      1,
    }
    created, err := repositories.InsertSkill(skill)
    if err != nil {
        t.Errorf("InsertSkill returned error: %v", err)
    }
    if created.ID == 0 {
        t.Errorf("expected skill ID to be set, got 0")
    }
}

func TestGetAllSkills(t *testing.T) {
    skills, err := repositories.GetAllSkills()
    if err != nil {
        t.Errorf("GetAllSkills returned error: %v", err)
    }
    if len(skills) == 0 {
        t.Error("expected at least one dummy skill, got none")
    }
}

