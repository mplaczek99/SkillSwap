# repositories Directory

This directory implements the data access layer for the application. It abstracts away direct database interactions from the rest of the code.

## Files

- **user_repository.go**:  
  Contains functions to create, retrieve, and query user records.

- **skill_repository.go**:  
  Contains functions to insert and query skill records.

- **transaction_repository.go**:  
  Contains functions to record and retrieve transactions (SkillPoint exchanges).

## How It Works

Repositories use the GORM ORM (configured in the `config` package) to perform CRUD operations. They return errors to be handled in the services layer.

