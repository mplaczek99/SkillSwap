# services Directory

This directory contains the business logic layer for the application.

## Files

- **auth_service.go**:  
  Handles authentication logic, including verifying credentials and generating JWT tokens.

- **user_service.go**:  
  Processes user registration, password hashing, and user data retrieval.

- **skill_service.go**:  
  Manages operations related to skill creation and retrieval.

## How It Works

Services receive input from controllers, enforce business rules, and delegate database operations to the repositories. They return processed data or errors to the controllers.

