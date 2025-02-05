# models Directory

This directory defines the data models (structs) for the SkillSwap application. These models represent the entities stored in the database.

## Files

- **user.go**:  
  Defines the `User` model with fields such as ID, name, email, password (omitted from JSON responses), bio, and creation timestamp.

- **skill.go**:  
  Defines the `Skill` model representing skills that a user can offer or request, including fields for name, description, and user association.

- **transaction.go**:  
  Defines the `Transaction` model for tracking SkillPoint exchanges between users.

## Usage

These models are used throughout the application (in controllers, services, and repositories) to ensure consistent data handling.

