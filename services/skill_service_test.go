package services_test

import (
    "testing"

    "github.com/mplaczek99/SkillSwap/models"
    "github.com/mplaczek99/SkillSwap/services"
)

func TestCreateSkill(t *testing.T) {
    tests := []struct {
        name          string
        skill         models.Skill
        expectError   bool
        errorContains string
    }{
        {
            name: "Valid Skill",
            skill: models.Skill{
                Name:        "Cooking",
                Description: "Italian cuisine",
                UserID:      1,
            },
            expectError: false,
        },
        {
            name: "Missing Name",
            skill: models.Skill{
                Description: "No name here",
                UserID:      2,
            },
            expectError:   true,
            errorContains: "skill name is required",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := services.CreateSkill(&tt.skill)
            if tt.expectError && err == nil {
                t.Errorf("expected an error, got nil")
            }
            if !tt.expectError && err != nil {
                t.Errorf("expected no error, got %v", err)
            }
        })
    }
}

func TestGetAllSkills(t *testing.T) {
    skills, err := services.GetAllSkills()
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if len(skills) == 0 {
        t.Errorf("expected dummy skills, got zero")
    }
}

